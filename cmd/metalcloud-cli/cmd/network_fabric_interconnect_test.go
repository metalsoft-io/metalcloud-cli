package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var interconnectItem = map[string]interface{}{
	"id": "12", "interconnectType": "dci-evpn", "label": "dc1-dc2",
	"name": "DC1 to DC2", "bgpConfigurationTemplateId": 3.0, "revision": "7",
	"status":           "draft",
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
}

func newInterconnectTestServer(t *testing.T, lastBody *string) *httptest.Server {
	t.Helper()
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/network-fabric-interconnects", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				var body map[string]interface{}
				_ = json.NewDecoder(r.Body).Decode(&body)
				raw, _ := json.Marshal(body)
				*lastBody = string(raw)
				jsonResponse(w, http.StatusCreated, interconnectItem)
				return
			}
			jsonResponse(w, http.StatusOK, paginatedList(interconnectItem))
		})
		mux.HandleFunc("/api/v2/network-fabric-interconnects/12", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, interconnectItem)
		})
	})
	return httptest.NewServer(mux)
}

func TestNetworkFabricInterconnectListAndAliases(t *testing.T) {
	var body string
	srv := newInterconnectTestServer(t, &body)
	defer srv.Close()

	for _, name := range []string{"network-fabric-interconnect", "fabric-interconnect", "interconnect", "nfi"} {
		out, err := runCLI(t, srv, name, "list")
		if err != nil {
			t.Fatalf("%s list: unexpected error: %v", name, err)
		}
		if !strings.Contains(out, "dc1-dc2") {
			t.Errorf("%s list: expected output to contain 'dc1-dc2', got: %s", name, out)
		}
	}
}

func TestNetworkFabricInterconnectGet(t *testing.T) {
	var body string
	srv := newInterconnectTestServer(t, &body)
	defer srv.Close()

	out, err := runCLI(t, srv, "nfi", "get", "12")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"label":"dc1-dc2"`) {
		t.Errorf("expected label in output, got: %s", out)
	}

	if _, err := runCLI(t, srv, "nfi", "get"); err == nil {
		t.Error("expected error when no arg provided")
	}
}

func TestNetworkFabricInterconnectCreateFromFlags(t *testing.T) {
	var body string
	srv := newInterconnectTestServer(t, &body)
	defer srv.Close()

	_, err := runCLI(t, srv, "nfi", "create", "--label", "dc1-dc2", "--bgp-template-id", "3", "--name", "DC1 to DC2", "--transport-id", "9")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{`"label":"dc1-dc2"`, `"bgpConfigurationTemplateId":3`, `"interconnectType":"dci-evpn"`, `"name":"DC1 to DC2"`, `"transportId":9`} {
		if !strings.Contains(body, want) {
			t.Errorf("create body missing %s: %s", want, body)
		}
	}
	if strings.Contains(body, "description") {
		t.Errorf("unset optional flags must not be sent: %s", body)
	}
}

func TestNetworkFabricInterconnectCreateFlagValidation(t *testing.T) {
	var body string
	srv := newInterconnectTestServer(t, &body)
	defer srv.Close()

	// Neither --config-source nor --label.
	if _, err := runCLI(t, srv, "nfi", "create"); err == nil {
		t.Error("expected error when no source flag given")
	}
	// --label without --bgp-template-id.
	if _, err := runCLI(t, srv, "nfi", "create", "--label", "x"); err == nil {
		t.Error("expected error when --label given without --bgp-template-id")
	}
	// Invalid type.
	if _, err := runCLI(t, srv, "nfi", "create", "--label", "x", "--bgp-template-id", "1", "--type", "bogus"); err == nil {
		t.Error("expected error for invalid --type")
	}
}

func TestNetworkFabricInterconnectHelpListsCommands(t *testing.T) {
	out, err := runCLI(t, nil, "network-fabric-interconnect", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, sub := range []string{"list", "get", "create", "update", "delete", "deploy", "deployment-check", "accept-deploy", "reject-deploy", "detach", "add-link", "remove-link", "activate-links", "deactivate-links", "get-fabrics", "get-available-fabrics", "template"} {
		if !strings.Contains(out, sub) {
			t.Errorf("help missing sub-command %q", sub)
		}
	}
}
