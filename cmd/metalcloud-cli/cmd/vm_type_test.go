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

// --- vm-type vms without datacenterName ---

// TestVmTypeVms_NoDatacenterName guards against SDK schema drift: the SDK VM
// model marks `datacenterName` required but the real API omits it (regression:
// "no value given for required property datacenterName").
func TestVmTypeVms_NoDatacenterName(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/vm-types/1/vms", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(map[string]interface{}{
				"id": 65, "name": "vm-65", "siteId": 2081, "powerState": "off",
				"infrastructureId": 1, "userId": 1, "instanceId": 1, "vmInstanceId": 1,
				"host": "host-1", "hosts": []interface{}{"host-1"},
				"cpuCores": 4, "ramGB": 8, "diskSizeGB": 100, "typeId": 1, "poolId": 1,
				"administrationState": "active",
				"powerStateLastUpdatedTimestamp": "2024-01-01T00:00:00Z",
				"createdTimestamp":               "2024-01-01T00:00:00Z",
				"allocationTimestamp":            "2024-01-01T00:00:00Z",
				"disks":                          []interface{}{},
			}))
		})
	}))
	defer srv.Close()

	out, err := runCLI(t, srv, "vm-type", "vms", "1", "-f", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "vm-65") {
		t.Errorf("expected output to contain vm-65, got: %s", out)
	}
}
