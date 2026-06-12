package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var accountItem = map[string]interface{}{
	"id": 1.0, "name": "acme", "revision": 1.0,
	"limits": map[string]interface{}{},
	"config": map[string]interface{}{"revision": 1.0, "name": "acme"},
}

func newAccountTestServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(accountItem))
		})
		mux.HandleFunc("/api/v2/accounts/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(accountItem)
		})
	})
	return httptest.NewServer(mux)
}

func TestAccountList(t *testing.T) {
	srv := newAccountTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "account", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "acme") {
		t.Errorf("expected output to contain 'acme', got: %s", out)
	}
}

func TestAccountListAlias(t *testing.T) {
	srv := newAccountTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "accounts", "ls")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "acme") {
		t.Errorf("expected output to contain 'acme', got: %s", out)
	}
}

func TestAccountGet(t *testing.T) {
	srv := newAccountTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "account", "get", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "acme") {
		t.Errorf("expected output to contain 'acme', got: %s", out)
	}
}

func TestAccountGetRequiresArg(t *testing.T) {
	srv := newAccountTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "account", "get")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}

func TestAccountHelp(t *testing.T) {
	out, err := runCLI(t, nil, "account", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "account") {
		t.Errorf("expected help output to contain 'account', got: %s", out)
	}
}

func TestAccountCreate(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				jsonResponse(w, http.StatusOK, accountItem)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(accountItem))
		})
		mux.HandleFunc("/api/v2/accounts/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(accountItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "account-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"name":"acme"}`)
	f.Close()

	_, execErr := runCLI(t, srv, "account", "create", "--config-source", f.Name())
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestAccountUpdate(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(accountItem))
		})
		mux.HandleFunc("/api/v2/accounts/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(accountItem)
		})
		mux.HandleFunc("/api/v2/accounts/1/config", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPatch {
				jsonResponse(w, http.StatusOK, accountItem)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(accountItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "account-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"name":"acme-updated"}`)
	f.Close()

	_, execErr := runCLI(t, srv, "account", "update", "1", "--config-source", f.Name())
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestAccountList_Formats(t *testing.T) {
	srv := newAccountTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "account", "list")
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
