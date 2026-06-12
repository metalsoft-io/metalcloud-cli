package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// server_default_credentials_test.go covers create/delete/update — list and get
// are already covered in server_creds_test.go.

func newServerDefaultCredentialsWriteServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/servers/default-credentials/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Method {
			case http.MethodDelete:
				w.WriteHeader(http.StatusNoContent)
			case http.MethodPatch, http.MethodPut:
				_ = json.NewEncoder(w).Encode(serverDefaultCredentialsItem)
			default:
				_ = json.NewEncoder(w).Encode(serverDefaultCredentialsItem)
			}
		})
		mux.HandleFunc("/api/v2/servers/default-credentials", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Method {
			case http.MethodPost:
				_ = json.NewEncoder(w).Encode(serverDefaultCredentialsItem)
			default:
				_ = json.NewEncoder(w).Encode(paginatedList(serverDefaultCredentialsItem))
			}
		})
	})
	return httptest.NewServer(mux)
}

func TestServerDefaultCredentialsCreate(t *testing.T) {
	srv := newServerDefaultCredentialsWriteServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "server-default-credentials", "create",
		"--site-id", "1",
		"--serial", "SN123456",
		"--mac", "aa:bb:cc:dd:ee:ff",
		"--username", "admin",
		"--password", "secret",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServerDefaultCredentialsDelete(t *testing.T) {
	srv := newServerDefaultCredentialsWriteServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "server-default-credentials", "delete", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServerDefaultCredentialsDeleteRequiresArg(t *testing.T) {
	srv := newServerDefaultCredentialsWriteServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "server-default-credentials", "delete")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}
