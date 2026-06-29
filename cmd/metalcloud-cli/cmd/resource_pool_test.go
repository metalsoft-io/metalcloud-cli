package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var resourcePoolItem = map[string]interface{}{
	"resourcePoolId":               1.0,
	"resourcePoolLabel":            "test-pool",
	"resourcePoolDescription":      "Test pool",
	"resourcePoolCreatedTimestamp": "2024-01-01T00:00:00Z",
	"resourcePoolUpdatedTimestamp": "2024-01-01T00:00:00Z",
	"statistics": map[string]interface{}{
		"users":       0.0,
		"servers":     0.0,
		"subnetPools": 0.0,
	},
}

func newResourcePoolTestServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/resource-pools/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resourcePoolItem)
		})
		mux.HandleFunc("/api/v2/resource-pools", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(resourcePoolItem))
		})
	})
	return httptest.NewServer(mux)
}

func TestResourcePoolList(t *testing.T) {
	srv := newResourcePoolTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "resource-pool", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-pool") {
		t.Errorf("expected output to contain 'test-pool', got: %s", out)
	}
}

func TestResourcePoolListPaginated(t *testing.T) {
	srv := newResourcePoolTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "resource-pool", "list", "--page", "1", "--limit", "5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-pool") {
		t.Errorf("expected output to contain 'test-pool', got: %s", out)
	}
}

func TestResourcePoolGet(t *testing.T) {
	srv := newResourcePoolTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "resource-pool", "get", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-pool") {
		t.Errorf("expected output to contain 'test-pool', got: %s", out)
	}
}

func TestResourcePoolGetRequiresArg(t *testing.T) {
	srv := newResourcePoolTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "resource-pool", "get")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}

// --- resource-pool create ---

func TestResourcePoolCreate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/resource-pools", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resourcePoolItem)
		})
	}))
	defer srv.Close()

	if _, execErr := runCLI(t, srv, "resource-pool", "create", "--label", "test-pool", "--description", "Test pool"); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

// --- resource-pool delete ---

func TestResourcePoolDelete(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/resource-pools/1", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
	}))
	defer srv.Close()

	if _, execErr := runCLI(t, srv, "resource-pool", "delete", "1"); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestResourcePoolList_Formats(t *testing.T) {
	srv := newResourcePoolTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "resource-pool", "list")
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
