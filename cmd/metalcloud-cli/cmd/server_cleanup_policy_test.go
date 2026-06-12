package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// server_cleanup_policy_test.go covers:
//   server-cleanup-policy list
//   server-cleanup-policy get <id>
//   server-cleanup-policy create --label ...
//   server-cleanup-policy update <id> --label ...
//   server-cleanup-policy delete <id>

func newCleanupPolicyTestServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/servers/cleanup-policies/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(serverCleanupPolicyItem)
			case http.MethodPatch, http.MethodPut:
				_ = json.NewEncoder(w).Encode(serverCleanupPolicyItem)
			case http.MethodDelete:
				w.WriteHeader(http.StatusNoContent)
			}
		})
		mux.HandleFunc("/api/v2/servers/cleanup-policies", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(paginatedList(serverCleanupPolicyItem))
			case http.MethodPost:
				_ = json.NewEncoder(w).Encode(serverCleanupPolicyItem)
			}
		})
	})
	return httptest.NewServer(mux)
}

func TestServerCleanupPolicyGet(t *testing.T) {
	srv := newCleanupPolicyTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "server-cleanup-policy", "get", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "default-policy") {
		t.Errorf("expected output to contain 'default-policy', got: %s", out)
	}
}

func TestServerCleanupPolicyGetRequiresArg(t *testing.T) {
	srv := newCleanupPolicyTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "server-cleanup-policy", "get")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}

func TestServerCleanupPolicyCreate(t *testing.T) {
	srv := newCleanupPolicyTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "server-cleanup-policy", "create",
		"--label", "default-policy",
		"--raid-one-drive", "none",
		"--raid-two-drives", "none",
		"--raid-even-drives", "none",
		"--raid-odd-drives", "none",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServerCleanupPolicyUpdate(t *testing.T) {
	srv := newCleanupPolicyTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "server-cleanup-policy", "update", "1",
		"--label", "updated-policy",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServerCleanupPolicyDelete(t *testing.T) {
	srv := newCleanupPolicyTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "server-cleanup-policy", "delete", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServerCleanupPolicyDeleteRequiresArg(t *testing.T) {
	srv := newCleanupPolicyTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "server-cleanup-policy", "delete")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}
