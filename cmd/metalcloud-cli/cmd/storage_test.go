package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// storageItem satisfies sdk.Storage required fields.
var storageItem = map[string]interface{}{
	"id":             1.0,
	"revision":       1.0,
	"siteId":         1.0,
	"datacenterName": "dc-01",
	"driver":         "iscsi_softlayer",
	"technologies":   []string{"iSCSI"},
	"status":         "active",
	"operationMode":  "rw",
	"name":           "storage-01",
	"managementHost": "10.0.0.1",
	"subnetType":     "oob",
}

func newStorageTestServer() *httptest.Server {
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		// storage list — FetchAllPages on GET /api/v2/storages
		mux.HandleFunc("/api/v2/storages", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(storageItem))
		})
		// storage get — GET /api/v2/storages/{storageId}
		mux.HandleFunc("/api/v2/storages/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(storageItem)
		})
	}))
}

// --- storage list ---

func TestStorageList_HappyPath(t *testing.T) {
	srv := newStorageTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "storage", "list")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "storage-01") {
		t.Fatalf("expected storage-01 in output, got: %s", out)
	}
}

func TestStorageList_Alias(t *testing.T) {
	srv := newStorageTestServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "storage", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestStorageList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "storage", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- storage get ---

func TestStorageGet_HappyPath(t *testing.T) {
	srv := newStorageTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "storage", "get", "1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "storage-01") {
		t.Fatalf("expected storage-01 in output, got: %s", out)
	}
}

func TestStorageGet_NoArgs(t *testing.T) {
	srv := newStorageTestServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "storage", "get"); err == nil {
		t.Fatal("expected error when no args given to storage get")
	}
}

// --- storage create ---

func TestStorageCreate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/storages", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(storageItem)
		})
	}))
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "storage-create-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"siteId":1,"driver":"iscsi_softlayer","technologies":["iSCSI"],"name":"storage-new","managementHost":"10.0.0.1","subnetType":"oob"}`)
	f.Close()

	if _, execErr := runCLI(t, srv, "storage", "create", "--config-source", f.Name()); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

// --- storage delete ---

func TestStorageDelete(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/storages/1", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
	}))
	defer srv.Close()

	if _, execErr := runCLI(t, srv, "storage", "delete", "1"); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestStorageList_Formats(t *testing.T) {
	srv := newStorageTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "storage", "list")
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
