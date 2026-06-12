package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// logicalNetworkProfileFixture satisfies sdk.LogicalNetworkProfile required fields.
func logicalNetworkProfileFixture(id int) map[string]interface{} {
	return map[string]interface{}{
		"id":          id,
		"label":       "vlan-profile-1",
		"name":        "VLAN Profile 1",
		"annotations": map[string]interface{}{},
		"createdAt":   "2024-01-01T00:00:00Z",
		"updatedAt":   "2024-01-01T00:00:00Z",
		"revision":    1,
		"kind":        "vlan",
		"fabricId":    1,
	}
}

func newLogicalNetworkProfileServer() *httptest.Server {
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/logical-network-profiles", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(logicalNetworkProfileFixture(2))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(logicalNetworkProfileFixture(1)))
		})
		mux.HandleFunc("/api/v2/logical-network-profiles/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if r.Method == http.MethodPatch || r.Method == http.MethodPut {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(logicalNetworkProfileFixture(1))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(logicalNetworkProfileFixture(1))
		})
	}))
}

// List tests are in logical_network_test.go (TestLogicalNetworkProfileList_*).

// --- logical-network-profile get ---

func TestLogicalNetworkProfileGet_HappyPath(t *testing.T) {
	srv := newLogicalNetworkProfileServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "logical-network-profile", "get", "1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "vlan-profile-1") {
		t.Fatalf("expected profile label in output, got: %s", out)
	}
}

func TestLogicalNetworkProfileGet_Alias(t *testing.T) {
	srv := newLogicalNetworkProfileServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "logical-network-profile", "show", "1"); err != nil {
		t.Fatalf("alias show: expected no error, got: %v", err)
	}
}

func TestLogicalNetworkProfileGet_NoArgs(t *testing.T) {
	srv := newLogicalNetworkProfileServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "logical-network-profile", "get"); err == nil {
		t.Fatal("expected error when no args given")
	}
}

// --- logical-network-profile create ---

func TestLogicalNetworkProfileCreate_HappyPath(t *testing.T) {
	srv := newLogicalNetworkProfileServer()
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "lnp-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(`{"label":"vlan-profile-new","name":"New","kind":"vlan","fabricId":1}`)
	f.Close()

	if _, err := runCLI(t, srv, "logical-network-profile", "create", "vlan", "--config-source", f.Name()); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestLogicalNetworkProfileCreate_MissingFlag(t *testing.T) {
	srv := newLogicalNetworkProfileServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "logical-network-profile", "create", "vlan"); err == nil {
		t.Fatal("expected error when --config-source is missing")
	}
}

// --- logical-network-profile delete ---

func TestLogicalNetworkProfileDelete_HappyPath(t *testing.T) {
	srv := newLogicalNetworkProfileServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "logical-network-profile", "delete", "1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestLogicalNetworkProfileDelete_NoArgs(t *testing.T) {
	srv := newLogicalNetworkProfileServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "logical-network-profile", "delete"); err == nil {
		t.Fatal("expected error when no args given")
	}
}
