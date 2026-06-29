package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// driveItem is a SharedDrive fixture (the list endpoint returns SharedDrivePaginatedList).
var driveItem = map[string]interface{}{
	"id":                      1.0,
	"revision":                1.0,
	"label":                   "test-drive",
	"sizeMb":                  10240.0,
	"storageType":             "iscsi",
	"infrastructureId":        1.0,
	"infrastructure":          map[string]interface{}{"id": 1.0},
	"serviceStatus":           "active",
	"storageUpdatedTimestamp": "2024-01-01T00:00:00Z",
	"allocationAffinity":      "none",
	"provisioningProtocol":    "iscsi",
	"createdTimestamp":        "2024-01-01T00:00:00Z",
	"updatedTimestamp":        "2024-01-01T00:00:00Z",
	"config": map[string]interface{}{
		"revision":         1.0,
		"label":            "test-drive",
		"groupId":          1.0,
		"sizeMb":           10240.0,
		"storageType":      "iscsi",
		"deployType":       "soft",
		"deployStatus":     "not_started",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	},
	"meta": map[string]interface{}{"name": "test-drive"},
}

func newDriveTestServer() *httptest.Server {
	mux := newInfraMux(func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/infrastructures/1/drives/1/snapshots", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			// GetDriveSnapshots returns []SharedDriveSnapshot (plain array, not paginated).
			_ = json.NewEncoder(w).Encode([]interface{}{})
		})
		mux.HandleFunc("/api/v2/infrastructures/1/drives", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(driveItem))
		})
	})
	return httptest.NewServer(mux)
}

func TestDriveList(t *testing.T) {
	srv := newDriveTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "drive", "list", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-drive") {
		t.Errorf("expected output to contain 'test-drive', got: %s", out)
	}
}

func TestDriveSnapshotList(t *testing.T) {
	srv := newDriveTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "drive", "snapshot", "list", "1", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func newDriveWriteTestServer() *httptest.Server {
	mux := newInfraMux(func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/infrastructures/1/drives/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			_ = json.NewEncoder(w).Encode(driveItem)
		})
		mux.HandleFunc("/api/v2/infrastructures/1/drives", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Method == http.MethodPost {
				_ = json.NewEncoder(w).Encode(driveItem)
				return
			}
			_ = json.NewEncoder(w).Encode(paginatedList(driveItem))
		})
	})
	return httptest.NewServer(mux)
}

func TestDriveCreate(t *testing.T) {
	srv := newDriveWriteTestServer()
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "drive-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(`{"label":"test-drive","sizeMb":10240,"storageType":"iscsi","storagePoolId":1}`)
	f.Close()

	_, err = runCLI(t, srv, "drive", "create", "1", "--config-source", f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDriveDelete(t *testing.T) {
	srv := newDriveWriteTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "drive", "delete", "1", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDriveListRequiresArg(t *testing.T) {
	srv := newDriveTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "drive", "list")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}

func TestDriveList_Formats(t *testing.T) {
	srv := newDriveTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "drive", "list", "1")
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
