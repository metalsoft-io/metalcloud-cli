package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// vm_instance_test.go covers:
//   vm-instance list --infrastructure 1
//   vm-instance get 1 1
//   vm-instance start 1 1
//   vm-instance shutdown 1 1
//
// vm-pool list, vm-type list, and vm-instance-group list are in vm_test.go.

var vmInstanceItem = map[string]interface{}{
	"id":               1.0,
	"revision":         1.0,
	"label":            "test-vm-instance",
	"infrastructureId": 1.0,
	"infrastructure":   map[string]interface{}{"id": 1.0},
	"groupId":          1.0,
	"typeId":           1.0,
	"diskSizeGB":       20.0,
	"ramGB":            8.0,
	"cpuCores":         4.0,
	"serviceStatus":    "active",
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
	"config": map[string]interface{}{
		"revision":         1.0,
		"label":            "test-vm-instance",
		"typeId":           1.0,
		"deployType":       "soft",
		"deployStatus":     "not_started",
		"diskSizeGB":       20.0,
		"ramGB":            8.0,
		"cpuCores":         4.0,
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	},
	"meta": map[string]interface{}{},
}

func newVMInstanceTestServer() *httptest.Server {
	mux := newInfraMux(func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instances/1/start", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		})
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instances/1/shutdown", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		})
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instances/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(vmInstanceItem)
		})
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instances", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(vmInstanceItem))
		})
	})
	return httptest.NewServer(mux)
}

func TestVMInstanceList(t *testing.T) {
	srv := newVMInstanceTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance", "list", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMInstanceGet(t *testing.T) {
	srv := newVMInstanceTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance", "get", "1", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMInstanceListRequiresArg(t *testing.T) {
	srv := newVMInstanceTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance", "list")
	if err == nil {
		t.Fatal("expected error when no infra arg provided, got nil")
	}
}

func TestVMInstanceStart(t *testing.T) {
	srv := newVMInstanceTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance", "start", "1", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMInstanceShutdown(t *testing.T) {
	srv := newVMInstanceTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance", "shutdown", "1", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVMInstanceList_Formats(t *testing.T) {
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			srv := newVMInstanceTestServer()
			defer srv.Close()
			out, err := runCLIFormat(t, srv, format, "vm-instance", "list", "1")
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
