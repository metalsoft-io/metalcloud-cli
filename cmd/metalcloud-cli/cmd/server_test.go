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
		"processorNames":         []interface{}{"Intel Xeon"},
		"diskCount":              4.0, "serverClass": "bigdata",
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

// --- hardware / drift / onboarding commands ---

// srvDriftItem satisfies the required properties of sdk.ServerDriftHistory.
var srvDriftItem = map[string]interface{}{
	"id":                 9.0,
	"serverId":           1.0,
	"snapshotId":         "abc123",
	"configurationDrift": "- foo\n+ bar",
	"createdTimestamp":   "2024-01-01T00:00:00Z",
}

// srvSnapshotItem satisfies the required properties of sdk.ServerSnapshot.
var srvSnapshotItem = map[string]interface{}{
	"oid":       "5b17fbca",
	"message":   "caller-cron",
	"timestamp": "2024-01-01T00:00:00Z",
	"kind":      "cron",
}

// srvJobResponse satisfies sdk.HardwareRescanServerResponse and
// sdk.RegisterServerResponse (both require serverId and revision).
var srvJobResponse = map[string]interface{}{
	"serverId":     1.0,
	"revision":     8.0,
	"serverUUID":   "uuid-0001",
	"serialNumber": "SN-001",
	"jobInfo":      map[string]interface{}{"jobId": 42.0, "jobGroupId": 7.0},
}

func newServerHardwareTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/servers/1", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, serverItem)
		})
		mux.HandleFunc("/api/v2/servers/statistics", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{
				"serverStatus": map[string]interface{}{"available": 3, "used": 1},
				"site":         map[string]interface{}{"dc-1": 4},
			})
		})
		mux.HandleFunc("/api/v2/servers/1/drift-history", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(srvDriftItem))
		})
		mux.HandleFunc("/api/v2/servers/1/drift-history/9", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, srvDriftItem)
		})
		mux.HandleFunc("/api/v2/servers/1/drift-history/9/actions/acknowledge", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, srvDriftItem)
		})
		mux.HandleFunc("/api/v2/servers/1/snapshots", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(srvSnapshotItem))
		})
		mux.HandleFunc("/api/v2/servers/1/actions/sync-target-snapshot-with-latest-snapshot", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
		mux.HandleFunc("/api/v2/servers/1/actions/hardware-rescan", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, srvJobResponse)
		})
		mux.HandleFunc("/api/v2/servers/actions/register-production", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, srvJobResponse)
		})
		mux.HandleFunc("/api/v2/servers/unmanaged/import", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, serverItem)
		})
		for _, action := range []string{"connect-interface", "set-interfaces-default-fabric", "set-interfaces-redundancy-group"} {
			mux.HandleFunc("/api/v2/servers/1/actions/"+action, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})
		}
		mux.HandleFunc("/api/v2/servers/firmware/actions/batch-upgrade", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{
				"successful": map[string]interface{}{"1": map[string]interface{}{"jobId": 42.0, "jobGroupId": 7.0}},
				"failed":     map[string]interface{}{},
			})
		})
		mux.HandleFunc("/api/v2/servers/1/firmware/actions/batch-schedule-upgrade", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{"failed": map[string]interface{}{}})
		})
	}))
}

func srvConfigFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "server-cmd-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestServerDrift(t *testing.T) {
	srv := newServerHardwareTestServer(t)
	defer srv.Close()

	out, err := runCLI(t, srv, "server", "drift", "list", "1")
	if err != nil {
		t.Fatalf("drift list: %v", err)
	}
	if !strings.Contains(out, "abc123") {
		t.Fatalf("expected the drift entry in the output, got: %s", out)
	}

	if _, err := runCLI(t, srv, "server", "drift", "get", "1", "9"); err != nil {
		t.Fatalf("drift get: %v", err)
	}
	if _, err := runCLI(t, srv, "server", "drift", "acknowledge", "1", "9"); err != nil {
		t.Fatalf("drift acknowledge: %v", err)
	}
	if _, err := runCLI(t, srv, "server", "drift", "get", "1"); err == nil {
		t.Fatal("expected an error when the drift ID is missing")
	}
}

func TestServerSnapshotsAndSyncTargetSnapshot(t *testing.T) {
	srv := newServerHardwareTestServer(t)
	defer srv.Close()

	out, err := runCLI(t, srv, "server", "snapshots", "1", "--kind", "cron")
	if err != nil {
		t.Fatalf("snapshots: %v", err)
	}
	if !strings.Contains(out, "5b17fbca") {
		t.Fatalf("expected the snapshot in the output, got: %s", out)
	}

	if _, err := runCLI(t, srv, "server", "sync-target-snapshot", "1"); err != nil {
		t.Fatalf("sync-target-snapshot: %v", err)
	}
}

func TestServerStatisticsCmd(t *testing.T) {
	srv := newServerHardwareTestServer(t)
	defer srv.Close()

	out, err := runCLI(t, srv, "server", "statistics")
	if err != nil {
		t.Fatalf("statistics: %v", err)
	}
	if !strings.Contains(out, "serverStatus") {
		t.Fatalf("expected the statistics in the output, got: %s", out)
	}
}

func TestServerHardwareRescanCmd(t *testing.T) {
	srv := newServerHardwareTestServer(t)
	defer srv.Close()

	if _, err := runCLI(t, srv, "server", "hardware-rescan", "1", "--reboot-allowed"); err != nil {
		t.Fatalf("hardware-rescan: %v", err)
	}
	if _, err := runCLI(t, srv, "server", "hardware-rescan"); err == nil {
		t.Fatal("expected an error when the server ID is missing")
	}
}

func TestServerRegisterProductionCmd(t *testing.T) {
	srv := newServerHardwareTestServer(t)
	defer srv.Close()

	if _, err := runCLI(t, srv, "server", "register-production",
		"--site-id", "1", "--infrastructure-id", "100",
		"--management-address", "10.0.0.1", "--username", "admin", "--password", "secret"); err != nil {
		t.Fatalf("register-production with flags: %v", err)
	}

	path := srvConfigFile(t, `{"siteId":1,"settings":{"infrastructureId":100}}`)
	if _, err := runCLI(t, srv, "server", "register-production", "--config-source", path); err != nil {
		t.Fatalf("register-production with config: %v", err)
	}

	if _, err := runCLI(t, srv, "server", "register-production"); err == nil {
		t.Fatal("expected an error when neither --config-source nor --site-id is given")
	}
	if _, err := runCLI(t, srv, "server", "register-production", "--site-id", "1"); err == nil {
		t.Fatal("expected an error when --infrastructure-id is missing")
	}
}

func TestServerInterfaceCommands(t *testing.T) {
	srv := newServerHardwareTestServer(t)
	defer srv.Close()

	if _, err := runCLI(t, srv, "server", "connect-interface", "1",
		"--interface-id", "2", "--port-id", "Ethernet0", "--hostname", "leaf-01"); err != nil {
		t.Fatalf("connect-interface: %v", err)
	}
	if _, err := runCLI(t, srv, "server", "connect-interface", "1", "--interface-id", "2"); err == nil {
		t.Fatal("expected an error when --port-id and --hostname are missing")
	}

	if _, err := runCLI(t, srv, "server", "set-interfaces-default-fabric", "1",
		"--interface-id", "1", "--fabric-id", "10"); err != nil {
		t.Fatalf("set-interfaces-default-fabric: %v", err)
	}
	if _, err := runCLI(t, srv, "server", "set-interfaces-default-fabric", "1",
		"--interface-id", "1", "--fabric-id", "10", "--clear-fabric"); err == nil {
		t.Fatal("--fabric-id and --clear-fabric should be mutually exclusive")
	}

	if _, err := runCLI(t, srv, "server", "set-interfaces-redundancy-group", "1",
		"--interface-id", "1", "--group-index", "1"); err != nil {
		t.Fatalf("set-interfaces-redundancy-group: %v", err)
	}
}

func TestServerImportUnmanagedCmd(t *testing.T) {
	srv := newServerHardwareTestServer(t)
	defer srv.Close()

	path := srvConfigFile(t, `{"siteId":1,"serverTypeId":2,"serverInterfaces":[]}`)
	if _, err := runCLI(t, srv, "server", "import-unmanaged", "--config-source", path); err != nil {
		t.Fatalf("import-unmanaged: %v", err)
	}
	if _, err := runCLI(t, srv, "server", "import-unmanaged"); err == nil {
		t.Fatal("expected an error when --config-source is missing")
	}
}

func TestServerConfigExampleCmd(t *testing.T) {
	srv := newServerHardwareTestServer(t)
	defer srv.Close()

	out, err := runCLI(t, srv, "server", "config-example", "register-production")
	if err != nil {
		t.Fatalf("config-example: %v", err)
	}
	if !strings.Contains(out, "infrastructureId") {
		t.Fatalf("expected the production registration example, got: %s", out)
	}

	if _, err := runCLI(t, srv, "server", "config-example", "nope"); err == nil {
		t.Fatal("expected an error for an unknown configuration kind")
	}
}

func TestServerFirmwareBatchCommands(t *testing.T) {
	srv := newServerHardwareTestServer(t)
	defer srv.Close()

	if _, err := runCLI(t, srv, "server", "firmware", "upgrade-batch", "1", "2"); err != nil {
		t.Fatalf("upgrade-batch: %v", err)
	}
	if _, err := runCLI(t, srv, "server", "firmware", "upgrade-batch"); err == nil {
		t.Fatal("expected an error when no server ID is given")
	}

	if _, err := runCLI(t, srv, "server", "firmware", "schedule-upgrade-batch", "1", "2",
		"--schedule-timestamp", "2024-01-01T10:00:00Z", "--confirmation-required"); err != nil {
		t.Fatalf("schedule-upgrade-batch: %v", err)
	}
}
