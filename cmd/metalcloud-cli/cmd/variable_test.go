package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var variableItem = map[string]interface{}{
	"id": 1.0, "userIdOwner": 1.0, "name": "test-var",
	"value":            map[string]interface{}{"data": "test-value"},
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
}

func newVariableTestServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/variables", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(variableItem))
		})
		mux.HandleFunc("/api/v2/variables/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(variableItem)
		})
	})
	return httptest.NewServer(mux)
}

func TestVariableList(t *testing.T) {
	srv := newVariableTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "variable", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-var") {
		t.Errorf("expected output to contain 'test-var', got: %s", out)
	}
}

func TestVariableGet(t *testing.T) {
	srv := newVariableTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "variable", "get", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-var") {
		t.Errorf("expected output to contain 'test-var', got: %s", out)
	}
}

func TestVariableGetRequiresArg(t *testing.T) {
	srv := newVariableTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "variable", "get")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}

func TestVariableHelp(t *testing.T) {
	out, err := runCLI(t, nil, "variable", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "variable") {
		t.Errorf("expected help output to contain 'variable', got: %s", out)
	}
}

func TestVariableCreate(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/variables", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				jsonResponse(w, http.StatusOK, variableItem)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(variableItem))
		})
		mux.HandleFunc("/api/v2/variables/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(variableItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "variable-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"name":"test-var","value":{"data":"test-value"}}`)
	f.Close()

	_, execErr := runCLI(t, srv, "variable", "create", "--config-source", f.Name())
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestVariableUpdate(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/variables", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(variableItem))
		})
		mux.HandleFunc("/api/v2/variables/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPatch {
				jsonResponse(w, http.StatusOK, variableItem)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(variableItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "variable-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"name":"test-var-updated","value":{"data":"new-value"}}`)
	f.Close()

	_, execErr := runCLI(t, srv, "variable", "update", "1", "--config-source", f.Name())
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestVariableDelete(t *testing.T) {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/variables", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(variableItem))
		})
		mux.HandleFunc("/api/v2/variables/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(variableItem)
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	_, execErr := runCLI(t, srv, "variable", "delete", "1")
	if execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestVariableList_Formats(t *testing.T) {
	srv := newVariableTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "variable", "list")
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
