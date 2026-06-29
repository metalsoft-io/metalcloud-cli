package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// roleItem matches the SDK Role struct: id and permissions are string/[]string.
var roleItem = map[string]interface{}{
	"id": "1", "name": "admin", "label": "Admin",
	"type":           "system",
	"permissions":    []interface{}{},
	"quotaProfileId": nil,
}

func newRoleTestServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		// GetRoles returns a paginated list: {"data": [...], "meta": {...}}
		mux.HandleFunc("/api/v2/roles", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(roleItem))
		})
		mux.HandleFunc("/api/v2/roles/admin", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(roleItem)
		})
	})
	return httptest.NewServer(mux)
}

func TestRoleList(t *testing.T) {
	srv := newRoleTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "role", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "admin") {
		t.Errorf("expected output to contain 'admin', got: %s", out)
	}
}

func TestRoleListAlias(t *testing.T) {
	srv := newRoleTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "roles", "ls")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "admin") {
		t.Errorf("expected output to contain 'admin', got: %s", out)
	}
}

func TestRoleGet(t *testing.T) {
	srv := newRoleTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "role", "get", "admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "admin") {
		t.Errorf("expected output to contain 'admin', got: %s", out)
	}
}

func TestRoleGetRequiresArg(t *testing.T) {
	srv := newRoleTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "role", "get")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}

func TestRoleHelp(t *testing.T) {
	out, err := runCLI(t, nil, "role", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "role") {
		t.Errorf("expected help output to contain 'role', got: %s", out)
	}
}

func TestRoleCreate(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/roles", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				jsonResponse(w, http.StatusOK, roleItem)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(roleItem))
		})
		mux.HandleFunc("/api/v2/roles/admin", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(roleItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "role-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"label":"Admin","permissions":[]}`)
	f.Close()

	_, execErr := runCLI(t, srv, "role", "create", "--config-source", f.Name())
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestRoleUpdate(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/roles", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(roleItem))
		})
		mux.HandleFunc("/api/v2/roles/admin", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPatch {
				jsonResponse(w, http.StatusOK, roleItem)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(roleItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "role-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"label":"Admin Updated","permissions":[]}`)
	f.Close()

	_, execErr := runCLI(t, srv, "role", "update", "admin", "--config-source", f.Name())
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestRoleDelete(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/roles", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(roleItem))
		})
		mux.HandleFunc("/api/v2/roles/admin", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(roleItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	_, execErr := runCLI(t, srv, "role", "delete", "admin")
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestRoleList_Formats(t *testing.T) {
	srv := newRoleTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "role", "list")
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
