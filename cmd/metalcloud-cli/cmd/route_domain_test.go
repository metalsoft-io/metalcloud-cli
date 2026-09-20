package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// routeDomainItem is the route domain entity; routeDomainConfigItem is its
// config sub-resource. The revisions differ (2 vs 9) so the config commands'
// If-Match assertions are meaningful.
var routeDomainItem = map[string]interface{}{
	"id": 2, "label": "tenant1", "name": "tenant1", "kind": "evpn_l3vpn", "revision": 2,
	"annotations": map[string]interface{}{}, "serviceStatus": "ordered",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
}

var routeDomainConfigItem = map[string]interface{}{
	"id": 2, "deployType": "create", "deployStatus": "not_started",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
	"revision": 9, "kind": "evpn_l3vpn",
	"vrfAllocationStrategies":    []interface{}{},
	"l3VlanAllocationStrategies": []interface{}{},
	"l3VniAllocationStrategies":  []interface{}{},
	"autoRouteDistinguisher":     false,
	"autoRouteTarget":            false,
}

func newRouteDomainConfigServer(lastIfMatch *string, lastBody *string) *httptest.Server {
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/route-domains/2", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, routeDomainItem)
		})
		mux.HandleFunc("/api/v2/route-domains/2/config", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPatch {
				*lastIfMatch = r.Header.Get("If-Match")
				var body map[string]interface{}
				_ = json.NewDecoder(r.Body).Decode(&body)
				raw, _ := json.Marshal(body)
				*lastBody = string(raw)
			}
			jsonResponse(w, http.StatusOK, routeDomainConfigItem)
		})
	}))
}

func TestRouteDomainGetConfig(t *testing.T) {
	var ifMatch, body string
	srv := newRouteDomainConfigServer(&ifMatch, &body)
	defer srv.Close()

	out, err := runCLI(t, srv, "route-domain", "get-config", "2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"revision":9`) {
		t.Errorf("expected the config revision in output, got: %s", out)
	}

	if _, err := runCLI(t, srv, "route-domain", "get-config"); err == nil {
		t.Error("expected an error when no route domain id is given")
	}
}

func TestRouteDomainUpdateConfig(t *testing.T) {
	var ifMatch, body string
	srv := newRouteDomainConfigServer(&ifMatch, &body)
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "rd-config-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"autoRouteTarget":true}`)
	f.Close()

	if _, execErr := runCLI(t, srv, "route-domain", "update-config", "2", "--config-source", f.Name()); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
	if ifMatch != "9" {
		t.Errorf("expected If-Match 9 (the config revision), got %q", ifMatch)
	}
	if !strings.Contains(body, `"autoRouteTarget":true`) {
		t.Errorf("update-config body missing autoRouteTarget: %s", body)
	}
}

func TestRouteDomainUpdateConfig_RequiresConfigSource(t *testing.T) {
	var ifMatch, body string
	srv := newRouteDomainConfigServer(&ifMatch, &body)
	defer srv.Close()

	if _, err := runCLI(t, srv, "route-domain", "update-config", "2"); err == nil {
		t.Fatal("expected an error when --config-source is missing")
	}
}
