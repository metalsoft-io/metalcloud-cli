package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var lniItem = map[string]interface{}{
	"id": "4", "revision": "6", "label": "dc1-dc2-ln", "name": "DC1 to DC2 network",
	"kind": "dci-evpn", "fabricInterconnectId": 12.0, "status": "active",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
}

var lniLinkItem = map[string]interface{}{
	"id": "9", "logicalNetworkId": 44.0, "logicalNetworkInterconnectId": 4.0, "status": "active",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
}

func newLniTestServer(t *testing.T, lastBody *string) *httptest.Server {
	t.Helper()
	capture := func(r *http.Request) {
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		raw, _ := json.Marshal(body)
		*lastBody = string(raw)
	}

	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/logical-network-interconnects", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				capture(r)
				jsonResponse(w, http.StatusCreated, lniItem)
				return
			}
			jsonResponse(w, http.StatusOK, paginatedList(lniItem))
		})
		mux.HandleFunc("/api/v2/logical-network-interconnects/4", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if r.Method != http.MethodGet {
				capture(r)
			}
			jsonResponse(w, http.StatusOK, lniItem)
		})
		mux.HandleFunc("/api/v2/logical-network-interconnects/4/links", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				capture(r)
				jsonResponse(w, http.StatusCreated, lniLinkItem)
				return
			}
			jsonResponse(w, http.StatusOK, paginatedList(lniLinkItem))
		})
		mux.HandleFunc("/api/v2/logical-network-interconnects/4/links/9", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			jsonResponse(w, http.StatusOK, lniLinkItem)
		})
	})
	return httptest.NewServer(mux)
}

func TestLogicalNetworkInterconnectListAndAliases(t *testing.T) {
	var body string
	srv := newLniTestServer(t, &body)
	defer srv.Close()

	for _, name := range []string{"logical-network-interconnect", "ln-interconnect", "lni"} {
		out, err := runCLI(t, srv, name, "list")
		if err != nil {
			t.Fatalf("%s list: unexpected error: %v", name, err)
		}
		if !strings.Contains(out, "dc1-dc2-ln") {
			t.Errorf("%s list: expected output to contain 'dc1-dc2-ln', got: %s", name, out)
		}
	}
}

func TestLogicalNetworkInterconnectGetCmd(t *testing.T) {
	var body string
	srv := newLniTestServer(t, &body)
	defer srv.Close()

	out, err := runCLI(t, srv, "lni", "get", "4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"label":"dc1-dc2-ln"`) {
		t.Errorf("expected label in output, got: %s", out)
	}

	if _, err := runCLI(t, srv, "lni", "get"); err == nil {
		t.Error("expected error when no arg provided")
	}
}

func TestLogicalNetworkInterconnectCreateFromFlags(t *testing.T) {
	var body string
	srv := newLniTestServer(t, &body)
	defer srv.Close()

	_, err := runCLI(t, srv, "lni", "create", "--label", "dc1-dc2-ln",
		"--name", "DC1 to DC2 network", "--fabric-interconnect-id", "12")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{`"label":"dc1-dc2-ln"`, `"name":"DC1 to DC2 network"`, `"fabricInterconnectId":12`, `"kind":"dci-evpn"`} {
		if !strings.Contains(body, want) {
			t.Errorf("create body missing %s: %s", want, body)
		}
	}
}

func TestLogicalNetworkInterconnectCreateFlagValidation(t *testing.T) {
	var body string
	srv := newLniTestServer(t, &body)
	defer srv.Close()

	if _, err := runCLI(t, srv, "lni", "create"); err == nil {
		t.Error("expected error when no source flag given")
	}
	if _, err := runCLI(t, srv, "lni", "create", "--label", "x"); err == nil {
		t.Error("expected error when --label given without --name and --fabric-interconnect-id")
	}
	if _, err := runCLI(t, srv, "lni", "create", "--label", "x", "--name", "y",
		"--fabric-interconnect-id", "12", "--kind", "bogus"); err == nil {
		t.Error("expected error for invalid --kind")
	}
}

func TestLogicalNetworkInterconnectLinkCommands(t *testing.T) {
	var body string
	srv := newLniTestServer(t, &body)
	defer srv.Close()

	out, err := runCLI(t, srv, "lni", "link", "list", "4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"logicalNetworkId":44`) {
		t.Errorf("link list output wrong: %s", out)
	}

	if _, err := runCLI(t, srv, "lni", "links", "get", "4", "9"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := runCLI(t, srv, "lni", "link", "add", "4", "44"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(body, `"logicalNetworkId":44`) {
		t.Errorf("link add body wrong: %s", body)
	}

	if _, err := runCLI(t, srv, "lni", "link", "remove", "4", "9"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLogicalNetworkInterconnectHelpListsCommands(t *testing.T) {
	out, err := runCLI(t, nil, "logical-network-interconnect", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, sub := range []string{"list", "get", "create", "update", "delete", "config-example", "link"} {
		if !strings.Contains(out, sub) {
			t.Errorf("help missing sub-command %q", sub)
		}
	}
}
