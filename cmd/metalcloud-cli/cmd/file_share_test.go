package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var fileShareItem = map[string]interface{}{
	"id":               1.0,
	"revision":         1.0,
	"label":            "test-fs",
	"sizeGB":           10.0,
	"infrastructureId": 1.0,
	"infrastructure":   map[string]interface{}{"id": 1.0},
	"serviceStatus":    "active",
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
	"config": map[string]interface{}{
		"revision":         1.0,
		"sizeGB":           10.0,
		"label":            "test-fs",
		"deployType":       "soft",
		"deployStatus":     "not_started",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	},
	"meta": map[string]interface{}{"name": "test-fs"},
}

func newFileShareTestServer() *httptest.Server {
	mux := newInfraMux(func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/infrastructures/1/file-shares", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(fileShareItem))
		})
	})
	return httptest.NewServer(mux)
}

func TestFileShareList(t *testing.T) {
	srv := newFileShareTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "file-share", "list", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-fs") {
		t.Errorf("expected output to contain 'test-fs', got: %s", out)
	}
}

func TestFileShareListRequiresArg(t *testing.T) {
	srv := newFileShareTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "file-share", "list")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}

func newFileShareWriteTestServer() *httptest.Server {
	mux := newInfraMux(func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/infrastructures/1/file-shares/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			_ = json.NewEncoder(w).Encode(fileShareItem)
		})
		mux.HandleFunc("/api/v2/infrastructures/1/file-shares", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Method == http.MethodPost {
				_ = json.NewEncoder(w).Encode(fileShareItem)
				return
			}
			_ = json.NewEncoder(w).Encode(paginatedList(fileShareItem))
		})
	})
	return httptest.NewServer(mux)
}

func TestFileShareCreate(t *testing.T) {
	srv := newFileShareWriteTestServer()
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "fs-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(`{"label":"test-fs","sizeGB":10,"storagePoolId":1}`)
	f.Close()

	_, err = runCLI(t, srv, "file-share", "create", "1", "--config-source", f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFileShareDelete(t *testing.T) {
	srv := newFileShareWriteTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "file-share", "delete", "1", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFileShareList_Formats(t *testing.T) {
	srv := newFileShareTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "file-share", "list", "1")
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
