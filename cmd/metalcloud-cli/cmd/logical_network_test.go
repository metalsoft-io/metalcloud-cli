package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// makeLogicalNetworkFixture mirrors the fixture from internal/logical_network tests.
// config is a LogicalNetworkConfig with its own required fields.
func makeLogicalNetworkFixture(id int) map[string]interface{} {
	return map[string]interface{}{
		"id": id, "label": "ln-label", "name": "ln-name",
		"annotations": map[string]interface{}{},
		"createdAt":   "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
		"revision": 1, "kind": "vlan", "fabricId": 2,
		"infrastructureId":                   nil,
		"serviceStatus":                      "active",
		"lastAppliedLogicalNetworkProfileId": nil,
		"lastLogicalNetworkProfileAppliedAt": "2024-01-01T00:00:00Z",
		"config": map[string]interface{}{
			"id": 1, "deployType": "none", "deployStatus": "idle",
			"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
			"revision": 1, "kind": "vlan",
		},
	}
}

// Required: id, label, name, annotations, createdAt, updatedAt, revision, kind, fabricId
func makeLogicalNetworkProfileFixture(id int) map[string]interface{} {
	return map[string]interface{}{
		"id": id, "label": "cisco-profile", "name": "Cisco Profile",
		"annotations": map[string]interface{}{},
		"createdAt":   "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
		"revision": 1, "kind": "vlan", "fabricId": 1,
		"links": []interface{}{},
	}
}

// --- logical-network list ---

func TestLogicalNetworkList_HappyPath(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/logical-networks", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(makeLogicalNetworkFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "logical-network", "list"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestLogicalNetworkList_Alias(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/logical-networks", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(makeLogicalNetworkFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "logical-network", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestLogicalNetworkList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "logical-network", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- logical-network-profile list ---

func TestLogicalNetworkProfileList_HappyPath(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/logical-network-profiles", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(makeLogicalNetworkProfileFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "logical-network-profile", "list"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestLogicalNetworkProfileList_Alias(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/logical-network-profiles", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(makeLogicalNetworkProfileFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "logical-network-profile", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestLogicalNetworkProfileList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "logical-network-profile", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- format tests ---

func TestLogicalNetworkList_Formats(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/logical-networks", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(makeLogicalNetworkFixture(1)))
		})
	}))
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "logical-network", "list")
			if err != nil {
				t.Fatalf("format %s: %v", format, err)
			}
			if out == "" {
				t.Errorf("format %s: empty output", format)
			}
			if format == "json" && !json.Valid([]byte(out)) {
				t.Errorf("format json: invalid JSON: %s", out)
			}
			if format == "csv" && !strings.Contains(out, ",") {
				t.Errorf("format csv: no comma: %s", out)
			}
		})
	}
}

func TestLogicalNetworkProfileList_Formats(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/logical-network-profiles", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(makeLogicalNetworkProfileFixture(1)))
		})
	}))
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "logical-network-profile", "list")
			if err != nil {
				t.Fatalf("format %s: %v", format, err)
			}
			if out == "" {
				t.Errorf("format %s: empty output", format)
			}
			if format == "json" && !json.Valid([]byte(out)) {
				t.Errorf("format json: invalid JSON: %s", out)
			}
			if format == "csv" && !strings.Contains(out, ",") {
				t.Errorf("format csv: no comma: %s", out)
			}
		})
	}
}

// --- logical-network create ---

func TestLogicalNetworkCreate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/logical-networks", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(makeLogicalNetworkFixture(1))
		})
	}))
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "ln-create-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"kind":"vlan","fabricId":1,"label":"test-net"}`)
	f.Close()

	if _, execErr := runCLI(t, srv, "logical-network", "create", "vlan", "--config-source", f.Name()); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

// --- logical-network delete ---

func TestLogicalNetworkDelete(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/logical-networks/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			// GET — return fixture so delete can read the revision
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(makeLogicalNetworkFixture(1))
		})
	}))
	defer srv.Close()

	if _, execErr := runCLI(t, srv, "logical-network", "delete", "1"); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

// --- logical-network config, create-from-profile and attached resources ---

// lnConfigFixture is the config sub-resource. Its revision (7) differs from
// the logical network's (1) so the If-Match assertions are meaningful.
var lnConfigFixture = map[string]interface{}{
	"id": 1, "deployType": "create", "deployStatus": "not_started",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
	"revision": 7, "kind": "vlan", "mtu": nil,
}

var lnExternalConnectionFixture = map[string]interface{}{
	"id": "4", "label": "ext-conn", "name": "External connection", "fabricId": 2,
	"revision": "1", "createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
}

var lnInterconnectFixture = map[string]interface{}{
	"id": "9", "label": "ln-interconnect", "name": "LN interconnect", "revision": "1",
	"kind": "dci-evpn", "fabricInterconnectId": 12, "status": "active",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
}

// newLogicalNetworkConfigServer records the If-Match header and body of the
// last mutating request so the config-revision rule can be asserted.
func newLogicalNetworkConfigServer(lastIfMatch *string, lastBody *string) *httptest.Server {
	record := func(r *http.Request) {
		*lastIfMatch = r.Header.Get("If-Match")
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		raw, _ := json.Marshal(body)
		*lastBody = string(raw)
	}

	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/logical-networks/1", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, makeLogicalNetworkFixture(1))
		})
		mux.HandleFunc("/api/v2/logical-networks/1/config", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPatch {
				record(r)
			}
			jsonResponse(w, http.StatusOK, lnConfigFixture)
		})
		mux.HandleFunc("/api/v2/logical-networks/1/config/actions/apply-profiles", func(w http.ResponseWriter, r *http.Request) {
			record(r)
			jsonResponse(w, http.StatusOK, lnConfigFixture)
		})
		mux.HandleFunc("/api/v2/logical-networks/actions/create-from-profile", func(w http.ResponseWriter, r *http.Request) {
			record(r)
			jsonResponse(w, http.StatusCreated, makeLogicalNetworkFixture(2))
		})
		mux.HandleFunc("/api/v2/logical-networks/1/external-connections", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(lnExternalConnectionFixture))
		})
		mux.HandleFunc("/api/v2/logical-networks/1/external-connection-logical-networks", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(map[string]interface{}{
				"id": 8, "externalConnectionId": 4, "logicalNetworkId": 1, "status": "active",
				"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
			}))
		})
		mux.HandleFunc("/api/v2/logical-networks/1/logical-network-interconnects", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(lnInterconnectFixture))
		})
		mux.HandleFunc("/api/v2/logical-networks/1/external-connections/4", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
	}))
}

func TestLogicalNetworkGetConfig(t *testing.T) {
	var ifMatch, body string
	srv := newLogicalNetworkConfigServer(&ifMatch, &body)
	defer srv.Close()

	out, err := runCLI(t, srv, "logical-network", "get-config", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"revision":7`) {
		t.Errorf("expected the config revision in output, got: %s", out)
	}
}

func TestLogicalNetworkUpdateConfig(t *testing.T) {
	var ifMatch, body string
	srv := newLogicalNetworkConfigServer(&ifMatch, &body)
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "ln-config-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"mtu":9000}`)
	f.Close()

	if _, execErr := runCLI(t, srv, "logical-network", "update-config", "1", "--config-source", f.Name()); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
	if ifMatch != "7" {
		t.Errorf("expected If-Match 7 (the config revision), got %q", ifMatch)
	}
	if !strings.Contains(body, `"mtu":9000`) {
		t.Errorf("update-config body missing the mtu: %s", body)
	}
}

func TestLogicalNetworkUpdateConfig_RequiresConfigSource(t *testing.T) {
	var ifMatch, body string
	srv := newLogicalNetworkConfigServer(&ifMatch, &body)
	defer srv.Close()

	if _, err := runCLI(t, srv, "logical-network", "update-config", "1"); err == nil {
		t.Fatal("expected an error when --config-source is missing")
	}
}

func TestLogicalNetworkApplyProfiles(t *testing.T) {
	var ifMatch, body string
	srv := newLogicalNetworkConfigServer(&ifMatch, &body)
	defer srv.Close()

	if _, err := runCLI(t, srv, "logical-network", "apply-profiles", "1", "3"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ifMatch != "7" {
		t.Errorf("expected If-Match 7 (the config revision), got %q", ifMatch)
	}
	if !strings.Contains(body, `"logicalNetworkProfileId":3`) {
		t.Errorf("apply-profiles body missing the profile id: %s", body)
	}

	if _, err := runCLI(t, srv, "logical-network", "apply-profiles", "1"); err == nil {
		t.Error("expected an error when the profile id is missing")
	}
}

func TestLogicalNetworkCreateFromProfile_Flags(t *testing.T) {
	var ifMatch, body string
	srv := newLogicalNetworkConfigServer(&ifMatch, &body)
	defer srv.Close()

	_, err := runCLI(t, srv, "logical-network", "create-from-profile",
		"--profile-id", "3", "--label", "from-profile", "--mtu", "9000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{`"logicalNetworkProfileId":3`, `"label":"from-profile"`, `"mtu":9000`} {
		if !strings.Contains(body, want) {
			t.Errorf("create-from-profile body missing %s: %s", want, body)
		}
	}
	if strings.Contains(body, "infrastructureId") {
		t.Errorf("unset optional flags must not be sent: %s", body)
	}
}

func TestLogicalNetworkCreateFromProfile_RequiresProfileIdOrConfigSource(t *testing.T) {
	var ifMatch, body string
	srv := newLogicalNetworkConfigServer(&ifMatch, &body)
	defer srv.Close()

	if _, err := runCLI(t, srv, "logical-network", "create-from-profile"); err == nil {
		t.Fatal("expected an error when neither --profile-id nor --config-source is given")
	}
}

func TestLogicalNetworkAttachedResourceCommands(t *testing.T) {
	var ifMatch, body string
	srv := newLogicalNetworkConfigServer(&ifMatch, &body)
	defer srv.Close()

	cases := []struct {
		command string
		want    string
	}{
		{"get-external-connections", "ext-conn"},
		{"get-external-connection-logical-networks", `"externalConnectionId":4`},
		{"get-interconnects", "ln-interconnect"},
	}

	for _, testCase := range cases {
		out, err := runCLI(t, srv, "logical-network", testCase.command, "1")
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", testCase.command, err)
		}
		if !strings.Contains(out, testCase.want) {
			t.Errorf("%s: output missing %s: %s", testCase.command, testCase.want, out)
		}
	}
}

func TestLogicalNetworkDetachExternalConnection(t *testing.T) {
	var ifMatch, body string
	srv := newLogicalNetworkConfigServer(&ifMatch, &body)
	defer srv.Close()

	if _, err := runCLI(t, srv, "logical-network", "detach-external-connection", "1", "4"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := runCLI(t, srv, "logical-network", "detach-external-connection", "1"); err == nil {
		t.Error("expected an error when the external connection id is missing")
	}
}
