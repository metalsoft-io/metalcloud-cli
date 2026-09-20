package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var negItem = map[string]interface{}{
	"id": "3", "name": "dc1-endpoint-group", "siteId": 1.0, "revision": "5",
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
}

var negLogicalNetworkItem = map[string]interface{}{
	"logicalNetworkId": "44", "tagged": true, "accessMode": "l2",
	"networkEndpointGroupId": "3", "mtu": 9000.0,
}

func newNegTestServer(t *testing.T, lastBody *string) *httptest.Server {
	t.Helper()
	capture := func(r *http.Request) {
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		raw, _ := json.Marshal(body)
		*lastBody = string(raw)
	}

	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/network-endpoint-groups", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				capture(r)
				jsonResponse(w, http.StatusCreated, negItem)
				return
			}
			jsonResponse(w, http.StatusOK, paginatedList(negItem))
		})
		mux.HandleFunc("/api/v2/network-endpoint-groups/3", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if r.Method != http.MethodGet {
				capture(r)
			}
			jsonResponse(w, http.StatusOK, negItem)
		})
		mux.HandleFunc("/api/v2/network-endpoint-groups/3/logical-networks", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				capture(r)
				w.WriteHeader(http.StatusNoContent)
				return
			}
			jsonResponse(w, http.StatusOK, map[string]interface{}{"data": []interface{}{negLogicalNetworkItem}})
		})
		mux.HandleFunc("/api/v2/network-endpoint-groups/3/logical-networks/44", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if r.Method != http.MethodGet {
				capture(r)
			}
			jsonResponse(w, http.StatusOK, negLogicalNetworkItem)
		})
	})
	return httptest.NewServer(mux)
}

func TestNetworkEndpointGroupListAndAliases(t *testing.T) {
	var body string
	srv := newNegTestServer(t, &body)
	defer srv.Close()

	for _, name := range []string{"network-endpoint-group", "neg"} {
		out, err := runCLI(t, srv, name, "list")
		if err != nil {
			t.Fatalf("%s list: unexpected error: %v", name, err)
		}
		if !strings.Contains(out, "dc1-endpoint-group") {
			t.Errorf("%s list: expected output to contain the group name, got: %s", name, out)
		}
	}
}

func TestNetworkEndpointGroupGetCmd(t *testing.T) {
	var body string
	srv := newNegTestServer(t, &body)
	defer srv.Close()

	out, err := runCLI(t, srv, "neg", "get", "3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"name":"dc1-endpoint-group"`) {
		t.Errorf("expected name in output, got: %s", out)
	}

	if _, err := runCLI(t, srv, "neg", "get"); err == nil {
		t.Error("expected error when no arg provided")
	}
}

func TestNetworkEndpointGroupCreateFromFlags(t *testing.T) {
	var body string
	srv := newNegTestServer(t, &body)
	defer srv.Close()

	if _, err := runCLI(t, srv, "neg", "create", "--name", "dc1-endpoint-group", "--site-id", "1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{`"name":"dc1-endpoint-group"`, `"siteId":1`} {
		if !strings.Contains(body, want) {
			t.Errorf("create body missing %s: %s", want, body)
		}
	}

	if _, err := runCLI(t, srv, "neg", "create"); err == nil {
		t.Error("expected error when neither --config-source nor --name given")
	}
}

func TestNetworkEndpointGroupLogicalNetworkCommands(t *testing.T) {
	var body string
	srv := newNegTestServer(t, &body)
	defer srv.Close()

	out, err := runCLI(t, srv, "neg", "logical-network", "list", "3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"logicalNetworkId":"44"`) {
		t.Errorf("logical network list output wrong: %s", out)
	}

	if _, err := runCLI(t, srv, "neg", "logical-networks", "get", "3", "44"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := runCLI(t, srv, "neg", "logical-network", "add", "3", "44", "--tagged", "--mtu", "9000"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{`"logicalNetworkId":"44"`, `"tagged":true`, `"accessMode":"l2"`, `"mtu":9000`} {
		if !strings.Contains(body, want) {
			t.Errorf("logical network add body missing %s: %s", want, body)
		}
	}

	if _, err := runCLI(t, srv, "neg", "logical-network", "update", "3", "44", "--mtu", "1500"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(body, `"mtu":1500`) {
		t.Errorf("logical network update body wrong: %s", body)
	}

	if _, err := runCLI(t, srv, "neg", "logical-network", "remove", "3", "44"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNetworkEndpointGroupHelpListsCommands(t *testing.T) {
	out, err := runCLI(t, nil, "network-endpoint-group", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, sub := range []string{"list", "get", "create", "update", "delete", "config-example", "logical-network"} {
		if !strings.Contains(out, sub) {
			t.Errorf("help missing sub-command %q", sub)
		}
	}
}
