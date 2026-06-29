package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var secretItem = map[string]interface{}{
	"id": 1.0, "userIdOwner": 1.0, "name": "test-secret",
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
}

func newSecretTestServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/secrets", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(secretItem))
		})
		mux.HandleFunc("/api/v2/secrets/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(secretItem)
		})
	})
	return httptest.NewServer(mux)
}

func TestSecretList(t *testing.T) {
	srv := newSecretTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "secret", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-secret") {
		t.Errorf("expected output to contain 'test-secret', got: %s", out)
	}
}

func TestSecretGet(t *testing.T) {
	srv := newSecretTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "secret", "get", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-secret") {
		t.Errorf("expected output to contain 'test-secret', got: %s", out)
	}
}

func TestSecretGetRequiresArg(t *testing.T) {
	srv := newSecretTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "secret", "get")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}

func TestSecretHelp(t *testing.T) {
	out, err := runCLI(t, nil, "secret", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "secret") {
		t.Errorf("expected help output to contain 'secret', got: %s", out)
	}
}

func TestSecretCreate(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/secrets", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				jsonResponse(w, http.StatusOK, secretItem)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(secretItem))
		})
		mux.HandleFunc("/api/v2/secrets/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(secretItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "secret-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"name":"test-secret","value":"s3cr3t"}`)
	f.Close()

	_, execErr := runCLI(t, srv, "secret", "create", "--config-source", f.Name())
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestSecretUpdate(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/secrets", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(secretItem))
		})
		mux.HandleFunc("/api/v2/secrets/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPatch {
				jsonResponse(w, http.StatusOK, secretItem)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(secretItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "secret-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"value":"updated-secret-value"}`)
	f.Close()

	_, execErr := runCLI(t, srv, "secret", "update", "1", "--config-source", f.Name())
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestSecretDelete(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/secrets", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(secretItem))
		})
		mux.HandleFunc("/api/v2/secrets/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(secretItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	_, execErr := runCLI(t, srv, "secret", "delete", "1")
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestSecretList_Formats(t *testing.T) {
	srv := newSecretTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "secret", "list")
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
