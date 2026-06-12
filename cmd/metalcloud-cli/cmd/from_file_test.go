package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// ---------------------------------------------------------------------------
// subnet --config-source
// ---------------------------------------------------------------------------

func newSubnetTestServerWithCreate() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/subnets", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Method == http.MethodPost {
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(subnetItem)
				return
			}
			_ = json.NewEncoder(w).Encode(paginatedList(subnetItem))
		})
		mux.HandleFunc("/api/v2/subnets/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(subnetItem)
		})
	})
	return httptest.NewServer(mux)
}

func TestSubnetCreate_FromFile(t *testing.T) {
	srv := newSubnetTestServerWithCreate()
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "subnet-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"name":"ci-subnet","networkAddress":"192.168.254.0","prefixLength":24,"ipVersion":"ipv4","isPool":false}`)
	f.Close()

	_, execErr := runCLI(t, srv, "subnet", "create", "--config-source", f.Name())
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestSubnetCreate_FromFile_Missing(t *testing.T) {
	srv := newSubnetTestServerWithCreate()
	defer srv.Close()

	_, execErr := runCLI(t, srv, "subnet", "create", "--config-source", "/nonexistent/path/subnet.json")
	if execErr == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

// ---------------------------------------------------------------------------
// secret --config-source
// ---------------------------------------------------------------------------

var secretItemCmd = map[string]interface{}{
	"id": 1.0, "userIdOwner": 1.0, "name": "ci-secret",
	"valueEncrypted":   "enc",
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
}

func newSecretTestServerWithCreate() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/secrets", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Method == http.MethodPost {
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(secretItemCmd)
				return
			}
			_ = json.NewEncoder(w).Encode(paginatedList(secretItemCmd))
		})
		mux.HandleFunc("/api/v2/secrets/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(secretItemCmd)
		})
	})
	return httptest.NewServer(mux)
}

func TestSecretCreate_FromFile(t *testing.T) {
	srv := newSecretTestServerWithCreate()
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "secret-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"name":"ci-secret","value":"plain-text"}`)
	f.Close()

	_, execErr := runCLI(t, srv, "secret", "create", "--config-source", f.Name())
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestSecretCreate_FromFile_Missing(t *testing.T) {
	srv := newSecretTestServerWithCreate()
	defer srv.Close()

	_, execErr := runCLI(t, srv, "secret", "create", "--config-source", "/nonexistent/path/secret.json")
	if execErr == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

// ---------------------------------------------------------------------------
// variable --config-source
// ---------------------------------------------------------------------------

var variableItemCmd = map[string]interface{}{
	"id": 1.0, "name": "ci-variable",
	"userIdOwner":      1.0,
	"value":            map[string]interface{}{},
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
}

func newVariableTestServerWithCreate() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/variables", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Method == http.MethodPost {
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(variableItemCmd)
				return
			}
			_ = json.NewEncoder(w).Encode(paginatedList(variableItemCmd))
		})
		mux.HandleFunc("/api/v2/variables/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(variableItemCmd)
		})
	})
	return httptest.NewServer(mux)
}

func TestVariableCreate_FromFile(t *testing.T) {
	srv := newVariableTestServerWithCreate()
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "variable-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"name":"ci-variable","value":{"env":"test"}}`)
	f.Close()

	_, execErr := runCLI(t, srv, "variable", "create", "--config-source", f.Name())
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestVariableCreate_FromFile_Missing(t *testing.T) {
	srv := newVariableTestServerWithCreate()
	defer srv.Close()

	out, execErr := runCLI(t, srv, "variable", "create", "--config-source", "/nonexistent/path/variable.json")
	if execErr == nil {
		t.Fatal("expected error for missing file, got nil")
	}
	_ = out
}
