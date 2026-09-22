package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var extConnItem = map[string]interface{}{
	"id": "12", "label": "dc1-ext", "name": "DC1 external", "fabricId": 7.0,
	"revision":  "4",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
}

var extConnInterfaceItem = map[string]interface{}{
	"id": 5.0, "networkDeviceInterfaceId": 101.0, "networkDeviceInterfaceName": "Ethernet1/1",
	"networkDeviceId": 45.0, "revision": "2",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
}

var extConnLogicalNetworkItem = map[string]interface{}{
	"id": 3.0, "externalConnectionId": 12.0, "logicalNetworkId": 44.0, "status": "active",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
}

var extConnFabricItem = map[string]interface{}{
	"id": "7", "name": "dc1-fabric", "siteId": 1.0, "revision": "1",
	"fabricConfiguration": map[string]interface{}{"fabricType": "ethernet"},
	"createdTimestamp":    "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
}

func newExtConnTestServer(t *testing.T, lastBody *string) *httptest.Server {
	t.Helper()
	capture := func(r *http.Request) {
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		raw, _ := json.Marshal(body)
		*lastBody = string(raw)
	}

	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/external-connections", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				capture(r)
				jsonResponse(w, http.StatusCreated, extConnItem)
				return
			}
			jsonResponse(w, http.StatusOK, paginatedList(extConnItem))
		})
		mux.HandleFunc("/api/v2/external-connections/12", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if r.Method != http.MethodGet {
				capture(r)
			}
			jsonResponse(w, http.StatusOK, extConnItem)
		})
		mux.HandleFunc("/api/v2/external-connections/12/interfaces", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				capture(r)
				jsonResponse(w, http.StatusCreated, extConnInterfaceItem)
				return
			}
			jsonResponse(w, http.StatusOK, paginatedList(extConnInterfaceItem))
		})
		mux.HandleFunc("/api/v2/external-connections/12/interfaces/5", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if r.Method != http.MethodGet {
				capture(r)
			}
			jsonResponse(w, http.StatusOK, extConnInterfaceItem)
		})
		mux.HandleFunc("/api/v2/external-connections/12/logical-networks", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				capture(r)
				jsonResponse(w, http.StatusCreated, extConnLogicalNetworkItem)
				return
			}
			jsonResponse(w, http.StatusOK, paginatedList(extConnLogicalNetworkItem))
		})
		mux.HandleFunc("/api/v2/external-connections/12/logical-networks/3", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			jsonResponse(w, http.StatusOK, extConnLogicalNetworkItem)
		})
		mux.HandleFunc("/api/v2/network-fabrics/7", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, extConnFabricItem)
		})
		mux.HandleFunc("/api/v2/network-fabrics", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(extConnFabricItem))
		})
	})
	return httptest.NewServer(mux)
}

func TestExternalConnectionListAndAliases(t *testing.T) {
	var body string
	srv := newExtConnTestServer(t, &body)
	defer srv.Close()

	for _, name := range []string{"external-connection", "ext-conn", "external-connections"} {
		out, err := runCLI(t, srv, name, "list")
		if err != nil {
			t.Fatalf("%s list: unexpected error: %v", name, err)
		}
		if !strings.Contains(out, "dc1-ext") {
			t.Errorf("%s list: expected output to contain 'dc1-ext', got: %s", name, out)
		}
	}
}

func TestExternalConnectionGetCmd(t *testing.T) {
	var body string
	srv := newExtConnTestServer(t, &body)
	defer srv.Close()

	out, err := runCLI(t, srv, "ext-conn", "get", "12")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"label":"dc1-ext"`) {
		t.Errorf("expected label in output, got: %s", out)
	}

	if _, err := runCLI(t, srv, "ext-conn", "get"); err == nil {
		t.Error("expected error when no arg provided")
	}
}

func TestExternalConnectionCreateFromFlags(t *testing.T) {
	var body string
	srv := newExtConnTestServer(t, &body)
	defer srv.Close()

	_, err := runCLI(t, srv, "ext-conn", "create", "--label", "dc1-ext", "--name", "DC1 external",
		"--fabric", "dc1-fabric", "--interface-ids", "101,102")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{`"label":"dc1-ext"`, `"name":"DC1 external"`, `"fabricId":7`, `"networkDeviceInterfaceId":101`, `"networkDeviceInterfaceId":102`} {
		if !strings.Contains(body, want) {
			t.Errorf("create body missing %s: %s", want, body)
		}
	}
}

func TestExternalConnectionCreateFlagValidation(t *testing.T) {
	var body string
	srv := newExtConnTestServer(t, &body)
	defer srv.Close()

	// Neither --config-source nor --label.
	if _, err := runCLI(t, srv, "ext-conn", "create"); err == nil {
		t.Error("expected error when no source flag given")
	}
	// --label without --name/--fabric.
	if _, err := runCLI(t, srv, "ext-conn", "create", "--label", "x"); err == nil {
		t.Error("expected error when --label given without --name and --fabric")
	}
}

func TestExternalConnectionInterfaceCommands(t *testing.T) {
	var body string
	srv := newExtConnTestServer(t, &body)
	defer srv.Close()

	out, err := runCLI(t, srv, "ext-conn", "interface", "list", "12")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Ethernet1/1") {
		t.Errorf("interface list output wrong: %s", out)
	}

	if _, err := runCLI(t, srv, "ext-conn", "interfaces", "get", "12", "5"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := runCLI(t, srv, "ext-conn", "interface", "add", "12", "101"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(body, `"networkDeviceInterfaceId":101`) {
		t.Errorf("interface add body wrong: %s", body)
	}

	if _, err := runCLI(t, srv, "ext-conn", "interface", "remove", "12", "5"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExternalConnectionLogicalNetworkCommands(t *testing.T) {
	var body string
	srv := newExtConnTestServer(t, &body)
	defer srv.Close()

	out, err := runCLI(t, srv, "ext-conn", "logical-network", "list", "12")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"logicalNetworkId":44`) {
		t.Errorf("logical network list output wrong: %s", out)
	}

	if _, err := runCLI(t, srv, "ext-conn", "logical-network", "add", "12", "44"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(body, `"logicalNetworkId":44`) {
		t.Errorf("logical network add body wrong: %s", body)
	}

	if _, err := runCLI(t, srv, "ext-conn", "logical-network", "remove", "12", "3"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExternalConnectionHelpListsCommands(t *testing.T) {
	out, err := runCLI(t, nil, "external-connection", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, sub := range []string{"list", "get", "create", "update", "delete", "config-example", "interface", "logical-network", "get-network-device-interfaces"} {
		if !strings.Contains(out, sub) {
			t.Errorf("help missing sub-command %q", sub)
		}
	}
}
