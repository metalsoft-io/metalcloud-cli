package cmd

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var allocStratItem = map[string]interface{}{
	"id": 3.0, "kind": "manual", "vlanId": 150.0, "granularityLevel": "fabric",
	"scope":     map[string]interface{}{"kind": "fabric", "resourceId": 3.0},
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
}

func newAllocStratTestServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/logical-networks/12", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{"id": 12.0, "revision": 9.0})
		})
		mux.HandleFunc("/api/v2/logical-networks/12/config", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{"id": 12.0, "revision": 5.0})
		})
		mux.HandleFunc("/api/v2/logical-networks/12/config/vlan/vlan-allocation-strategies", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(allocStratItem))
		})
		mux.HandleFunc("/api/v2/route-domains/2/config/vrf-allocation-strategies/7", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, allocStratItem)
		})
		mux.HandleFunc("/api/v2/route-domains/2/config", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{"id": 2.0, "revision": 1.0})
		})
	})
	return httptest.NewServer(mux)
}

func TestAllocationStrategyRegisteredOnAllParents(t *testing.T) {
	for _, parent := range []string{"logical-network", "logical-network-profile", "route-domain", "point-to-point-link"} {
		out, err := runCLI(t, nil, parent, "allocation-strategy", "--help")
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", parent, err)
		}
		for _, sub := range []string{"list", "get", "config-example", "add", "replace", "remove"} {
			if !strings.Contains(out, sub) {
				t.Errorf("%s allocation-strategy help missing %q", parent, sub)
			}
		}
	}
}

func TestAllocationStrategyListAndGetViaAliases(t *testing.T) {
	srv := newAllocStratTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "logical-network", "strategy", "list", "12", "vlan")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"vlanId":150`) {
		t.Errorf("expected strategy in list output, got: %s", out)
	}

	out, err = runCLI(t, srv, "route-domain", "strategies", "get", "2", "vrf", "7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"kind":"manual"`) {
		t.Errorf("expected strategy in get output, got: %s", out)
	}
}

func TestAllocationStrategyValidation(t *testing.T) {
	srv := newAllocStratTestServer()
	defer srv.Close()

	// Family not supported by this parent.
	if _, err := runCLI(t, srv, "route-domain", "strategy", "list", "2", "vlan"); err == nil {
		t.Error("expected error for vlan family on route-domain")
	}
	// add requires --config-source.
	if _, err := runCLI(t, srv, "logical-network", "strategy", "add", "12", "vlan"); err == nil {
		t.Error("expected error when --config-source is missing")
	}
	// Wrong arity.
	if _, err := runCLI(t, srv, "logical-network", "strategy", "get", "12", "vlan"); err == nil {
		t.Error("expected error for missing strategy_id")
	}
}

func TestAllocationStrategyConfigExample(t *testing.T) {
	srv := newAllocStratTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "point-to-point-link", "strategy", "config-example", "ipv4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, kind := range []string{`"kind":"auto"`, `"kind":"manual"`, `"kind":"unnumbered"`} {
		if !strings.Contains(out, kind) {
			t.Errorf("config example missing %s: %s", kind, out)
		}
	}
}
