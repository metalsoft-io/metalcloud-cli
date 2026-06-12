package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// serverItem satisfies both serverRaw (ServerList reads raw body) and the full
// sdk.Server required fields (ServerGet calls GetServerInfo which unmarshals).
var serverItem = map[string]interface{}{
	"serverId":          1.0,
	"siteId":            1.0,
	"serverTypeId":      1.0,
	"serverUUID":        "uuid-0001",
	"serialNumber":      "SN-001",
	"managementAddress": "10.0.0.1",
	"vendor":            "Dell",
	"model":             "R740",
	"serverStatus":      "registered",
	"revision":          1.0,
	// Required by sdk.Server (used by ServerGet / GetServerInfo).
	// bdkDebug, requiresReRegister, supportsFcProvisioning are float32 in the SDK.
	"datacenterName":                 "dc-01",
	"bdkDebug":                       0.0,
	"requiresReRegister":             0.0,
	"serverClass":                    "bigdata",
	"administrationState":            "active",
	"serverDhcpStatus":               "none",
	"supportsFcProvisioning":         0.0,
	"serverCreatedTimestamp":         "2024-01-01T00:00:00Z",
	"powerStatus":                    "on",
	"powerStatusLastUpdateTimestamp": "2024-01-01T00:00:00Z",
}

// serverTypeItem satisfies the ServerType required fields.
var serverTypeItem = map[string]interface{}{
	"id":                       1.0,
	"ramGbytes":                128.0,
	"processorCount":           2.0,
	"processorCoreMhz":         2400.0,
	"processorCoreCount":       16.0,
	"name":                     "Standard",
	"label":                    "standard",
	"networkTotalCapacityMbps": 10000.0,
	"networkInterfaceCount":    2.0,
	"networkInterfaceSpeeds":   []float64{10000},
	"processorNames":           []string{"Intel Xeon"},
	"diskCount":                4.0,
	"serverClass":              "bigdata",
}

func newServerTestServer() *httptest.Server {
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		// server list — ServerList reads raw body so wrap in {"data":[...]} only
		mux.HandleFunc("/api/v2/servers", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": []interface{}{serverItem},
			})
		})
		// server get — ServerGet calls GetServerInfo returning a Server object
		mux.HandleFunc("/api/v2/servers/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(serverItem)
		})
		// server-type list
		mux.HandleFunc("/api/v2/server-types", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(serverTypeItem))
		})
		// server-type get
		mux.HandleFunc("/api/v2/server-types/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(serverTypeItem)
		})
	}))
}

// --- server list ---

func TestServerList_HappyPath(t *testing.T) {
	srv := newServerTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "server", "list")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "serverId") {
		t.Fatalf("expected serverId in output, got: %s", out)
	}
}

func TestServerList_Alias(t *testing.T) {
	srv := newServerTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "server", "ls")
	if err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
	if !strings.Contains(out, "serverId") {
		t.Fatalf("expected serverId in output, got: %s", out)
	}
}

func TestServerList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "server", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- server get ---

func TestServerGet_HappyPath(t *testing.T) {
	srv := newServerTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "server", "get", "1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "serverId") {
		t.Fatalf("expected serverId in output, got: %s", out)
	}
}

func TestServerGet_NoArgs(t *testing.T) {
	srv := newServerTestServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "server", "get"); err == nil {
		t.Fatal("expected error when no args given to server get")
	}
}

// --- server-type list ---

func TestServerTypeList_HappyPath(t *testing.T) {
	srv := newServerTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "server-type", "list")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, `"id"`) {
		t.Fatalf("expected id in output, got: %s", out)
	}
}

func TestServerTypeList_Alias(t *testing.T) {
	srv := newServerTestServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "server-type", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestServerTypeList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "server-type", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- server-type get ---

func TestServerTypeGet_HappyPath(t *testing.T) {
	srv := newServerTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "server-type", "get", "1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "standard") {
		t.Fatalf("expected 'standard' in output, got: %s", out)
	}
}

func TestServerTypeGet_NoArgs(t *testing.T) {
	srv := newServerTestServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "server-type", "get"); err == nil {
		t.Fatal("expected error when no args given to server-type get")
	}
}

// --- server update ---

func TestServerUpdate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/servers/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(serverItem)
		})
	}))
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "server-update-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"vendor":"Dell"}`)
	f.Close()

	if _, execErr := runCLI(t, srv, "server", "update", "1", "--config-source", f.Name()); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

// --- server delete ---

func TestServerDelete(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/servers/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(serverItem)
		})
	}))
	defer srv.Close()

	if _, execErr := runCLI(t, srv, "server", "delete", "1"); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestServerList_Formats(t *testing.T) {
	srv := newServerTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "server", "list")
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

func TestServerTypeList_Formats(t *testing.T) {
	// serverTypeRaw.ProcessorNames is *string (a comma-joined string), but
	// sdk.ServerType.ProcessorNames is []string (used by the GET path).
	// The shared serverTypeItem uses []string to satisfy the GET path, so we
	// use a separate inline fixture with a string value for the list path.
	stListItem := map[string]interface{}{
		"id": 1.0, "ramGbytes": 128.0, "processorCount": 2.0,
		"processorCoreMhz": 2400.0, "processorCoreCount": 16.0,
		"name": "Standard", "label": "standard",
		"networkTotalCapacityMbps": 10000.0, "networkInterfaceCount": 2.0,
		"networkInterfaceSpeeds": []interface{}{10000.0},
		"processorNames":          "Intel Xeon",
		"diskCount":               4.0, "serverClass": "bigdata",
	}
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/server-types", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(stListItem))
		})
	}))
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "server-type", "list")
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
