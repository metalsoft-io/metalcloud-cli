package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// infraItem is a minimal Infrastructure JSON for ID=1, label="test-infra".
const infraItem = `{"id":1,"label":"test-infra","serviceStatus":"active","revision":1,"datacenterName":"dc1","siteId":1,"designIsLocked":0,"createdTimestamp":"2024-01-01T00:00:00Z","updatedTimestamp":"2024-01-01T00:00:00Z","config":{"deployStatus":"not_started","label":"test-infra","deployType":"soft"}}`

// infraListBody is the paginated list response for GET /api/v2/infrastructures.
// GetInfrastructureByIdOrLabel calls GET /api/v2/infrastructures?search=<id_or_label>.
const infraListBody = `{"data":[` + infraItem + `],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`

// newInfraMux returns a ServeMux with /api/v2/user and /api/v2/infrastructures
// pre-registered, plus any extra routes from the extra callback.
// Used by infra-scoped commands that call GetInfrastructureByIdOrLabel.
func newInfraMux(extra func(mux *http.ServeMux)) *http.ServeMux {
	return newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/infrastructures", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(infraListBody))
		})
		if extra != nil {
			extra(mux)
		}
	})
}

var serverInstanceItem = map[string]interface{}{
	"id":               1.0,
	"revision":         1.0,
	"label":            "si-1",
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
	"infrastructureId": 1.0,
	"groupId":          1.0,
	"serviceStatus":    "active",
	"isVmInstance":       0.0,
	"isEndpointInstance": 0.0,
	"meta":               map[string]interface{}{"name": "si-1"},
}

var serverInstanceGroupItem = map[string]interface{}{
	"id":                      1.0,
	"revision":                1.0,
	"label":                   "sig-1",
	"createdTimestamp":        "2024-01-01T00:00:00Z",
	"updatedTimestamp":        "2024-01-01T00:00:00Z",
	"infrastructureId":        1.0,
	"instanceCount":           1.0,
	"defaultServerTypeId":     0.0,
	"ipAllocateAuto":          1.0,
	"ipv4SubnetCreateAuto":    1.0,
	"processorCount":          1.0,
	"processorCoreCount":      1.0,
	"processorCoreMhz":        1000.0,
	"diskCount":               1.0,
	"diskSizeMbytes":          10240.0,
	"diskTypes":               []interface{}{},
	"virtualInterfacesEnabled": 0.0,
	"serviceStatus":           "active",
	"isVmGroup":               0.0,
	"isEndpointInstanceGroup": 0.0,
	"meta":                    map[string]interface{}{},
}

func newServerInstanceTestServer() *httptest.Server {
	mux := newInfraMux(func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/infrastructures/1/server-instances", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(serverInstanceItem))
		})
		mux.HandleFunc("/api/v2/infrastructures/1/server-instance-groups", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(serverInstanceGroupItem))
		})
	})
	return httptest.NewServer(mux)
}

func TestServerInstanceList(t *testing.T) {
	srv := newServerInstanceTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "server-instance", "list", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServerInstanceGroupList(t *testing.T) {
	srv := newServerInstanceTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "server-instance-group", "list", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServerInstanceListRequiresArg(t *testing.T) {
	srv := newServerInstanceTestServer()
	defer srv.Close()

	_, err := runCLI(t, srv, "server-instance", "list")
	if err == nil {
		t.Fatal("expected error when no arg provided, got nil")
	}
}

func TestServerInstanceList_Formats(t *testing.T) {
	srv := newServerInstanceTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "server-instance", "list", "1")
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

func TestServerInstanceGroupList_Formats(t *testing.T) {
	srv := newServerInstanceTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "server-instance-group", "list", "1")
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
