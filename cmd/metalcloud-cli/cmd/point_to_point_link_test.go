package cmd

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var ptpCmdLinkItem = map[string]interface{}{
	"id": 42.0, "label": "leaf01-to-spine01", "description": "uplink",
	"routingActivation": "default", "serviceStatus": "active", "revision": 3.0,
}

// ptpCmdConfigItem carries its own revision (11), which is the one that guards
// the config sub-resources - the link's revision is 3.
var ptpCmdConfigItem = map[string]interface{}{
	"id": 42.0, "deployType": "staged", "deployStatus": "not_started",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
	"revision": 11.0, "mtu": 9216.0,
}

var ptpCmdRouteItem = map[string]interface{}{"id": 5.0, "destinationPrefix": "10.0.0.0/24"}

type ptpCmdRequest struct {
	Method  string
	Path    string
	IfMatch string
	Body    string
}

func newPtpCmdServer(last *ptpCmdRequest) *httptest.Server {
	record := func(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			*last = ptpCmdRequest{Method: r.Method, Path: r.URL.Path, IfMatch: r.Header.Get("If-Match"), Body: string(body)}
			handler(w, r)
		}
	}

	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/point-to-point-links/42", record(func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, ptpCmdLinkItem)
		}))
		mux.HandleFunc("/api/v2/point-to-point-links/42/config", record(func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, ptpCmdConfigItem)
		}))
		mux.HandleFunc("/api/v2/point-to-point-links/42/config/ipv4/static-routes", record(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				jsonResponse(w, http.StatusCreated, ptpCmdRouteItem)
				return
			}
			jsonResponse(w, http.StatusOK, []interface{}{ptpCmdRouteItem})
		}))
		mux.HandleFunc("/api/v2/point-to-point-links/42/config/ipv4/static-routes/5", record(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			jsonResponse(w, http.StatusOK, ptpCmdRouteItem)
		}))
	})
	return httptest.NewServer(mux)
}

func ptpWriteTempConfig(t *testing.T, body string) string {
	t.Helper()
	path := t.TempDir() + "/config.json"
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	return path
}

func TestPointToPointLinkUpdateCmd(t *testing.T) {
	var last ptpCmdRequest
	srv := newPtpCmdServer(&last)
	defer srv.Close()

	configFile := ptpWriteTempConfig(t, `{"description":"leaf01 uplink"}`)
	if _, err := runCLI(t, srv, "point-to-point-link", "update", "42", "--config-source", configFile); err != nil {
		t.Fatalf("update: %v", err)
	}
	if last.Method != http.MethodPatch || last.Path != "/api/v2/point-to-point-links/42" {
		t.Errorf("update hit %s %s", last.Method, last.Path)
	}
	if last.IfMatch != "3" {
		t.Errorf("update must send the link revision as If-Match, got %q", last.IfMatch)
	}

	if _, err := runCLI(t, srv, "p2p", "update", "42"); err == nil {
		t.Error("expected an error when --config-source is missing")
	}
}

func TestPointToPointLinkConfigCmds(t *testing.T) {
	var last ptpCmdRequest
	srv := newPtpCmdServer(&last)
	defer srv.Close()

	out, err := runCLI(t, srv, "point-to-point-link", "get-config", "42")
	if err != nil {
		t.Fatalf("get-config: %v", err)
	}
	if !strings.Contains(out, "9216") {
		t.Errorf("get-config output missing mtu: %s", out)
	}

	configFile := ptpWriteTempConfig(t, `{"mtu":1500}`)
	if _, err := runCLI(t, srv, "p2p", "update-config", "42", "--config-source", configFile); err != nil {
		t.Fatalf("update-config: %v", err)
	}
	if last.Method != http.MethodPatch || last.Path != "/api/v2/point-to-point-links/42/config" {
		t.Errorf("update-config hit %s %s", last.Method, last.Path)
	}
	if last.IfMatch != "11" {
		t.Errorf("update-config must send the config revision as If-Match, got %q", last.IfMatch)
	}
}

func TestPointToPointLinkStaticRouteCmds(t *testing.T) {
	var last ptpCmdRequest
	srv := newPtpCmdServer(&last)
	defer srv.Close()

	out, err := runCLI(t, srv, "point-to-point-link", "static-route", "list", "42", "ipv4")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(out, "10.0.0.0/24") {
		t.Errorf("list output missing route: %s", out)
	}

	if _, err := runCLI(t, srv, "p2p", "route", "get", "42", "ipv4", "5"); err != nil {
		t.Fatalf("get: %v", err)
	}

	if _, err := runCLI(t, srv, "p2p", "static-route", "add", "42", "ipv4", "10.0.0.0/24"); err != nil {
		t.Fatalf("add: %v", err)
	}
	if last.Method != http.MethodPost || last.IfMatch != "11" {
		t.Errorf("add sent %s with If-Match %q", last.Method, last.IfMatch)
	}
	if !strings.Contains(last.Body, `"destinationPrefix":"10.0.0.0/24"`) {
		t.Errorf("add body wrong: %s", last.Body)
	}

	if _, err := runCLI(t, srv, "p2p", "static-route", "rm", "42", "ipv4", "5"); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if last.Method != http.MethodDelete || last.IfMatch != "11" {
		t.Errorf("remove sent %s with If-Match %q", last.Method, last.IfMatch)
	}

	if _, err := runCLI(t, srv, "p2p", "static-route", "list", "42", "ipv5"); err == nil {
		t.Error("expected an error for an invalid address family")
	}
	if _, err := runCLI(t, srv, "p2p", "static-route", "add", "42", "ipv4"); err == nil {
		t.Error("expected an error when the destination prefix is missing")
	}
}

func TestPointToPointLinkHelpListsNewCommands(t *testing.T) {
	out, err := runCLI(t, nil, "point-to-point-link", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, sub := range []string{"update", "get-config", "update-config", "static-route", "allocation-strategy"} {
		if !strings.Contains(out, sub) {
			t.Errorf("help missing sub-command %q", sub)
		}
	}
}
