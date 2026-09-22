package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var eiInfrastructureItem = map[string]interface{}{
	"id": 123.0, "revision": 1.0, "label": "prod-env", "serviceStatus": "active",
	"datacenterName": "dc1", "siteId": 1.0, "designIsLocked": 0.0,
	"config":           map[string]interface{}{},
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
}

var eiItem = map[string]interface{}{
	"id": 42.0, "revision": 7.0, "label": "ei-1", "infrastructureId": 123.0, "groupId": 12.0,
	"endpointId": 5.0, "serviceStatus": "active",
	"meta":             map[string]interface{}{"tags": []interface{}{"prod"}},
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

var eiConfigItem = map[string]interface{}{
	"revision": 3.0, "label": "ei-1", "updatedTimestamp": "2024-01-02T00:00:00Z",
	"groupId": 12.0, "endpointId": 5.0, "deployType": "create", "deployStatus": "not_started",
}

var eigItem = map[string]interface{}{
	"id": 12.0, "revision": 9.0, "label": "eig-1", "infrastructureId": 123.0,
	"endpointGroupName": "web", "serviceStatus": "active",
	"meta":             map[string]interface{}{"tags": []interface{}{"prod"}},
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

var eigConfigItem = map[string]interface{}{
	"revision": 4.0, "label": "eig-1", "updatedTimestamp": "2024-01-02T00:00:00Z",
	"endpointGroupName": "web", "deployType": "create", "deployStatus": "not_started",
}

var eigNetworkItem = map[string]interface{}{
	"id": "77", "name": "eig-1-neg", "siteId": 1.0, "revision": "2",
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

var eigConnectionItem = map[string]interface{}{
	"id": "5", "tagged": true, "accessMode": "l2", "mtu": 1500.0,
}

var eigACLItem = map[string]interface{}{
	"id": "3", "ruleType": "ipv4", "direction": "in", "sequence": 10.0,
	"forwardingAction": "allow", "enforcementPoint": "svi",
	"endpointGroupId": "77", "logicalNetworkId": "7",
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

// eiCapture records the body of the last write request served by the mock.
type eiCapture struct {
	body   string
	method string
	path   string
	ifMate string
}

func eiNewTestServer(t *testing.T, capture *eiCapture) *httptest.Server {
	t.Helper()

	record := func(r *http.Request) {
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		raw, _ := json.Marshal(body)
		capture.body = string(raw)
		capture.method = r.Method
		capture.path = r.URL.Path
		capture.ifMate = r.Header.Get("If-Match")
	}

	listOrWrite := func(list interface{}, created interface{}) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				jsonResponse(w, http.StatusOK, list)
				return
			}
			record(r)
			jsonResponse(w, http.StatusCreated, created)
		}
	}

	getOrWrite := func(object interface{}) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				jsonResponse(w, http.StatusOK, object)
				return
			}
			record(r)
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			jsonResponse(w, http.StatusOK, object)
		}
	}

	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/infrastructures", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(eiInfrastructureItem))
		})

		mux.HandleFunc("/api/v2/endpoint-instances", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(eiItem))
		})
		mux.HandleFunc("/api/v2/infrastructures/123/endpoint-instances",
			listOrWrite(paginatedList(eiItem), eiItem))
		mux.HandleFunc("/api/v2/endpoint-instances/42", getOrWrite(eiItem))
		mux.HandleFunc("/api/v2/endpoint-instances/42/config", getOrWrite(eiConfigItem))
		mux.HandleFunc("/api/v2/endpoint-instances/42/meta", func(w http.ResponseWriter, r *http.Request) {
			record(r)
			w.WriteHeader(http.StatusNoContent)
		})

		mux.HandleFunc("/api/v2/infrastructures/123/endpoint-instance-groups",
			listOrWrite(paginatedList(eigItem), eigItem))
		mux.HandleFunc("/api/v2/endpoint-instance-groups/12", getOrWrite(eigItem))
		mux.HandleFunc("/api/v2/endpoint-instance-groups/12/config", getOrWrite(eigConfigItem))
		mux.HandleFunc("/api/v2/endpoint-instance-groups/12/meta", func(w http.ResponseWriter, r *http.Request) {
			record(r)
			w.WriteHeader(http.StatusNoContent)
		})
		mux.HandleFunc("/api/v2/endpoint-instance-groups/12/endpoint-instances", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(eiItem))
		})
		mux.HandleFunc("/api/v2/endpoint-instance-groups/12/config/networking", getOrWrite(eigNetworkItem))
		mux.HandleFunc("/api/v2/endpoint-instance-groups/12/config/networking/connections",
			listOrWrite(map[string]interface{}{"data": []interface{}{eigConnectionItem}}, eigConnectionItem))
		mux.HandleFunc("/api/v2/endpoint-instance-groups/12/config/networking/connections/5", getOrWrite(eigConnectionItem))
		mux.HandleFunc("/api/v2/endpoint-instance-groups/12/config/networking/connections/5/security/rules",
			listOrWrite([]interface{}{eigACLItem}, eigACLItem))
		mux.HandleFunc("/api/v2/endpoint-instance-groups/12/config/networking/connections/5/security/rules/3", getOrWrite(eigACLItem))
	})

	return httptest.NewServer(mux)
}

func TestEndpointInstanceCommandAliases(t *testing.T) {
	capture := &eiCapture{}
	srv := eiNewTestServer(t, capture)
	defer srv.Close()

	for _, name := range []string{"endpoint-instance", "ei"} {
		out, err := runCLI(t, srv, name, "list")
		if err != nil {
			t.Fatalf("%s list: unexpected error: %v", name, err)
		}
		if !strings.Contains(out, "ei-1") {
			t.Errorf("%s list: expected 'ei-1' in output, got: %s", name, out)
		}
	}

	for _, name := range []string{"endpoint-instance-group", "eig"} {
		out, err := runCLI(t, srv, name, "list", "prod-env")
		if err != nil {
			t.Fatalf("%s list: unexpected error: %v", name, err)
		}
		if !strings.Contains(out, "eig-1") {
			t.Errorf("%s list: expected 'eig-1' in output, got: %s", name, out)
		}
	}
}

func TestEndpointInstanceListByInfrastructureFlag(t *testing.T) {
	capture := &eiCapture{}
	srv := eiNewTestServer(t, capture)
	defer srv.Close()

	out, err := runCLI(t, srv, "ei", "list", "--infrastructure", "prod-env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "ei-1") {
		t.Errorf("expected 'ei-1' in output, got: %s", out)
	}
}

func TestEndpointInstanceGetAndConfig(t *testing.T) {
	capture := &eiCapture{}
	srv := eiNewTestServer(t, capture)
	defer srv.Close()

	out, err := runCLI(t, srv, "ei", "get", "42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"id":42`) {
		t.Errorf("expected id in output, got: %s", out)
	}

	out, err = runCLI(t, srv, "ei", "config", "42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"deployType":"create"`) {
		t.Errorf("expected deploy type in output, got: %s", out)
	}

	if _, err := runCLI(t, srv, "ei", "get"); err == nil {
		t.Error("expected error when no argument is given")
	}
}

func TestEndpointInstanceCreateFromFlags(t *testing.T) {
	capture := &eiCapture{}
	srv := eiNewTestServer(t, capture)
	defer srv.Close()

	if _, err := runCLI(t, srv, "ei", "create", "prod-env", "--endpoint-id", "5", "--label", "ei-1", "--group-id", "12"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{`"endpointId":5`, `"label":"ei-1"`, `"groupId":12`} {
		if !strings.Contains(capture.body, want) {
			t.Errorf("create body missing %s: %s", want, capture.body)
		}
	}
	if capture.path != "/api/v2/infrastructures/123/endpoint-instances" {
		t.Errorf("create should post to the infrastructure endpoint, got %s", capture.path)
	}
}

func TestEndpointInstanceCreateFlagValidation(t *testing.T) {
	capture := &eiCapture{}
	srv := eiNewTestServer(t, capture)
	defer srv.Close()

	if _, err := runCLI(t, srv, "ei", "create", "prod-env"); err == nil {
		t.Error("expected error when neither --config-source nor --endpoint-id is given")
	}
	if _, err := runCLI(t, srv, "ei", "create", "prod-env", "--config-source", "x.json", "--endpoint-id", "5"); err == nil {
		t.Error("expected error when --config-source and --endpoint-id are combined")
	}
}

func TestEndpointInstanceDeleteAndUpdates(t *testing.T) {
	capture := &eiCapture{}
	srv := eiNewTestServer(t, capture)
	defer srv.Close()

	if _, err := runCLI(t, srv, "ei", "delete", "42"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capture.method != http.MethodDelete || capture.ifMate != "7" {
		t.Errorf("delete should send If-Match 7, got %s %q", capture.method, capture.ifMate)
	}

	if _, err := runCLI(t, srv, "ei", "update-config", "42", "--config-source", eiWriteTempJSON(t, `{"label":"renamed"}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capture.ifMate != "3" {
		t.Errorf("update-config should send the config revision as If-Match, got %q", capture.ifMate)
	}
	if !strings.Contains(capture.body, `"label":"renamed"`) {
		t.Errorf("update-config body missing label: %s", capture.body)
	}

	if _, err := runCLI(t, srv, "ei", "update-meta", "42", "--config-source", eiWriteTempJSON(t, `{"tags":["a"]}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capture.path != "/api/v2/endpoint-instances/42/meta" {
		t.Errorf("update-meta should post to the meta endpoint, got %s", capture.path)
	}
}

func TestEndpointInstanceGroupCreateAndNetworkCommands(t *testing.T) {
	capture := &eiCapture{}
	srv := eiNewTestServer(t, capture)
	defer srv.Close()

	if _, err := runCLI(t, srv, "eig", "create", "prod-env", "--label", "eig-1", "--endpoint-group-name", "web"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(capture.body, `"endpointGroupName":"web"`) {
		t.Errorf("create body missing endpoint group name: %s", capture.body)
	}

	out, err := runCLI(t, srv, "eig", "instances", "12")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "ei-1") {
		t.Errorf("instances output missing member: %s", out)
	}

	out, err = runCLI(t, srv, "eig", "network", "list", "12")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "eig-1-neg") {
		t.Errorf("network list output missing network endpoint group: %s", out)
	}

	out, err = runCLI(t, srv, "eig", "network", "connections", "12")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"accessMode":"l2"`) {
		t.Errorf("connections output missing connection: %s", out)
	}

	if _, err := runCLI(t, srv, "eig", "network", "connect", "12", "--logical-network-id", "7", "--tagged", "true"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{`"logicalNetworkId":"7"`, `"tagged":true`, `"accessMode":"l2"`} {
		if !strings.Contains(capture.body, want) {
			t.Errorf("connect body missing %s: %s", want, capture.body)
		}
	}

	if _, err := runCLI(t, srv, "eig", "network", "update", "12", "5", "--access-mode", "l2"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capture.method != http.MethodPatch {
		t.Errorf("connection update should PATCH, got %s", capture.method)
	}

	if _, err := runCLI(t, srv, "eig", "network", "disconnect", "12", "5"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capture.method != http.MethodDelete {
		t.Errorf("disconnect should DELETE, got %s", capture.method)
	}
}

func TestEndpointInstanceGroupNetworkConnectFlagValidation(t *testing.T) {
	capture := &eiCapture{}
	srv := eiNewTestServer(t, capture)
	defer srv.Close()

	if _, err := runCLI(t, srv, "eig", "network", "connect", "12"); err == nil {
		t.Error("expected error when neither --config-source nor --logical-network-id is given")
	}
	if _, err := runCLI(t, srv, "eig", "network", "connect", "12", "--logical-network-id", "7", "--tagged", "maybe"); err == nil {
		t.Error("expected error for a non-boolean --tagged value")
	}
}

func TestEndpointInstanceGroupACLCommands(t *testing.T) {
	capture := &eiCapture{}
	srv := eiNewTestServer(t, capture)
	defer srv.Close()

	out, err := runCLI(t, srv, "eig", "network", "acl", "list", "12", "5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"ruleType":"ipv4"`) {
		t.Errorf("acl list output missing rule: %s", out)
	}

	out, err = runCLI(t, srv, "eig", "network", "acl", "get", "12", "5", "3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"id":"3"`) {
		t.Errorf("acl get output missing id: %s", out)
	}

	if _, err := runCLI(t, srv, "eig", "network", "acl", "add", "12", "5", "--rule-type", "ipv4", "--sequence", "20", "--source-address", "10.0.0.0/24"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{`"ruleType":"ipv4"`, `"sequence":20`, `"sourceAddress":"10.0.0.0/24"`, `"direction":"in"`, `"forwardingAction":"allow"`} {
		if !strings.Contains(capture.body, want) {
			t.Errorf("acl add body missing %s: %s", want, capture.body)
		}
	}

	if _, err := runCLI(t, srv, "eig", "network", "acl", "update", "12", "5", "3", "--config-source", eiWriteTempJSON(t, `{"forwardingAction":"deny"}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(capture.body, `"forwardingAction":"deny"`) {
		t.Errorf("acl update body missing action: %s", capture.body)
	}

	if _, err := runCLI(t, srv, "eig", "network", "acl", "remove", "12", "5", "3"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capture.method != http.MethodDelete {
		t.Errorf("acl remove should DELETE, got %s", capture.method)
	}
}

func TestEndpointInstanceConfigExamples(t *testing.T) {
	capture := &eiCapture{}
	srv := eiNewTestServer(t, capture)
	defer srv.Close()

	cases := []struct {
		args []string
		want string
	}{
		{[]string{"ei", "config-example"}, `"endpointId"`},
		{[]string{"eig", "config-example"}, `"endpointGroupName"`},
		{[]string{"eig", "network", "config-example"}, `"logicalNetworkId"`},
		{[]string{"eig", "network", "acl", "config-example"}, `"ruleType"`},
	}
	for _, tc := range cases {
		out, err := runCLI(t, srv, tc.args...)
		if err != nil {
			t.Fatalf("%v: unexpected error: %v", tc.args, err)
		}
		if !strings.Contains(out, tc.want) {
			t.Errorf("%v: expected %s in output, got: %s", tc.args, tc.want, out)
		}
	}
}

func TestEndpointInstanceHelpListsCommands(t *testing.T) {
	out, err := runCLI(t, nil, "endpoint-instance", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, sub := range []string{"list", "get", "create", "delete", "config", "update-config", "update-meta", "config-example"} {
		if !strings.Contains(out, sub) {
			t.Errorf("endpoint-instance help missing sub-command %q", sub)
		}
	}

	out, err = runCLI(t, nil, "endpoint-instance-group", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, sub := range []string{"list", "get", "create", "delete", "config", "update-config", "update-meta", "instances", "network"} {
		if !strings.Contains(out, sub) {
			t.Errorf("endpoint-instance-group help missing sub-command %q", sub)
		}
	}

	out, err = runCLI(t, nil, "endpoint-instance-group", "network", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, sub := range []string{"list", "replace", "connections", "connect", "update", "disconnect", "acl"} {
		if !strings.Contains(out, sub) {
			t.Errorf("network help missing sub-command %q", sub)
		}
	}
}

// eiWriteTempJSON writes content to a temporary file and returns its path, so
// that --config-source can be exercised without a pipe.
func eiWriteTempJSON(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return path
}
