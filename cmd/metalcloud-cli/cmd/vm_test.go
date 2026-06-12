package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var vmPoolItem = map[string]interface{}{
	"id":             1.0,
	"siteId":         1.0,
	"datacenterName": "dc1",
	"managementHost": "vcenter.example.com",
	"managementPort": 443.0,
	"name":           "test-vm-pool",
	"type":           "vmware",
	"status":         "active",
	"networkFabricId": 1.0,
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
}

var vmTypeItem = map[string]interface{}{
	"id":     1.0,
	"name":   "test-vm-type",
	"cpuCores": 4.0,
	"ramGB":  8.0,
}

var vmInstanceGroupItem = map[string]interface{}{
	"id":               1.0,
	"revision":         1.0,
	"label":            "test-vmig",
	"infrastructureId": 1.0,
	"infrastructure":   map[string]interface{}{"id": 1.0},
	"serviceStatus":    "active",
	"diskSizeGB":       20.0,
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
	"config": map[string]interface{}{
		"revision":         1.0,
		"label":            "test-vmig",
		"deployType":       "soft",
		"deployStatus":     "not_started",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	},
	"meta": map[string]interface{}{},
}

func newVMTestServer() *httptest.Server {
	mux := newInfraMux(func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/vm-pools", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(vmPoolItem))
		})
		mux.HandleFunc("/api/v2/vm-types", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(vmTypeItem))
		})
		mux.HandleFunc("/api/v2/infrastructures/1/vm-instance-groups", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(vmInstanceGroupItem))
		})
	})
	return httptest.NewServer(mux)
}

func TestVMPoolList(t *testing.T) {
	srv := newVMTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "vm-pool", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-vm-pool") {
		t.Errorf("expected output to contain 'test-vm-pool', got: %s", out)
	}
}

func TestVMTypeList(t *testing.T) {
	srv := newVMTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "vm-type", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-vm-type") {
		t.Errorf("expected output to contain 'test-vm-type', got: %s", out)
	}
}

func TestVMInstanceGroupList(t *testing.T) {
	srv := newVMTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "vm-instance-group", "list", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "test-vmig") {
		t.Errorf("expected output to contain 'test-vmig', got: %s", out)
	}
}

func TestVMInstanceGroupListRequiresArg(t *testing.T) {
	srv := newVMTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "vm-instance-group", "list")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}

func TestVMPoolList_Formats(t *testing.T) {
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			srv := newVMTestServer()
			defer srv.Close()
			out, err := runCLIFormat(t, srv, format, "vm-pool", "list")
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

func TestVMTypeList_Formats(t *testing.T) {
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			srv := newVMTestServer()
			defer srv.Close()
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

func TestVMInstanceGroupList_Formats(t *testing.T) {
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			srv := newVMTestServer()
			defer srv.Close()
			out, err := runCLIFormat(t, srv, format, "vm-instance-group", "list", "1")
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
