package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Required: id, name, cpuCores, ramGB
func vmTypeFixture(id int) map[string]interface{} {
	return map[string]interface{}{
		"id": id, "name": "Standard-4", "cpuCores": 4, "ramGB": 8,
		"links": []interface{}{},
	}
}

func TestVmTypeList_HappyPath(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/vm-types", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(vmTypeFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "vm-type", "list"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestVmTypeList_WithPagination(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/vm-types", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(vmTypeFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "vm-type", "list", "--page", "1", "--limit", "10"); err != nil {
		t.Fatalf("expected no error with --page and --limit, got: %v", err)
	}
}

func TestVmTypeList_Alias(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/vm-types", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(vmTypeFixture(1)))
		})
	}))
	defer srv.Close()

	if _, err := runCLI(t, srv, "vm-type", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestVmTypeList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "vm-type", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

func TestVmTypeList_Formats(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/vm-types", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(vmTypeFixture(1)))
		})
	}))
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "vm-type", "list")
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
