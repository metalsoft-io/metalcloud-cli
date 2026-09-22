package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// server_type_test.go covers the server-type write commands, the statistics
// report and the configuration examples; list and get are covered in
// server_test.go.

// stypeStatisticsItem satisfies every required property of
// sdk.ServerTypeStatisticsBatch and sdk.ServerTypeUtilizationReport.
var stypeStatisticsItem = map[string]interface{}{
	"serverTypeIdToServerCount":       map[string]interface{}{"1": 3, "2": 0},
	"serverTypeIdToServerInformation": map[string]interface{}{},
	"utilizationReport": map[string]interface{}{
		"groupByServerRamGb":       map[string]interface{}{},
		"groupByServerTypeName":    map[string]interface{}{},
		"groupByServerProductName": map[string]interface{}{},
		"groupByUserId":            map[string]interface{}{},
	},
}

func stypeWriteServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/server-types/statistics", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, stypeStatisticsItem)
		})
		mux.HandleFunc("/api/v2/server-types/actions/clean-unused", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
		mux.HandleFunc("/api/v2/server-types/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(serverTypeItem)
		})
		mux.HandleFunc("/api/v2/server-types", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(serverTypeItem)
		})
	})
	return httptest.NewServer(mux)
}

func stypeConfigFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "server-type-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestServerTypeCreate(t *testing.T) {
	srv := stypeWriteServer(t)
	defer srv.Close()

	path := stypeConfigFile(t, `{"name":"Standard","label":"standard","ramGbytes":128,
		"processorCount":2,"processorCoreMhz":2400,"processorCoreCount":16,
		"processorNames":["Intel Xeon"],"networkInterfaceSpeeds":[10000],
		"diskCount":4,"serverClass":"bigdata"}`)

	if _, err := runCLI(t, srv, "server-type", "create", "--config-source", path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServerTypeCreateRequiresConfigSource(t *testing.T) {
	srv := stypeWriteServer(t)
	defer srv.Close()

	if _, err := runCLI(t, srv, "server-type", "create"); err == nil {
		t.Fatal("expected an error when --config-source is missing")
	}
}

func TestServerTypeUpdate(t *testing.T) {
	srv := stypeWriteServer(t)
	defer srv.Close()

	path := stypeConfigFile(t, `{"label":"standard","description":"updated"}`)

	if _, err := runCLI(t, srv, "server-type", "update", "1", "--config-source", path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServerTypeUpdateRequiresArg(t *testing.T) {
	srv := stypeWriteServer(t)
	defer srv.Close()

	path := stypeConfigFile(t, `{"label":"standard"}`)

	if _, err := runCLI(t, srv, "server-type", "update", "--config-source", path); err == nil {
		t.Fatal("expected an error when no server type ID is given")
	}
}

func TestServerTypeDelete(t *testing.T) {
	srv := stypeWriteServer(t)
	defer srv.Close()

	if _, err := runCLI(t, srv, "server-type", "delete", "1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServerTypeCleanUnused(t *testing.T) {
	srv := stypeWriteServer(t)
	defer srv.Close()

	if _, err := runCLI(t, srv, "server-type", "clean-unused"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServerTypeStatistics(t *testing.T) {
	srv := stypeWriteServer(t)
	defer srv.Close()

	out, err := runCLI(t, srv, "server-type", "statistics", "1", "2", "--site-id", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "serverTypeIdToServerCount") {
		t.Fatalf("expected the counts in the output, got: %s", out)
	}
}

func TestServerTypeStatisticsRequiresSiteId(t *testing.T) {
	srv := stypeWriteServer(t)
	defer srv.Close()

	if _, err := runCLI(t, srv, "server-type", "statistics"); err == nil {
		t.Fatal("expected an error when --site-id is missing")
	}
}

func TestServerTypeConfigExamples(t *testing.T) {
	srv := stypeWriteServer(t)
	defer srv.Close()

	out, err := runCLI(t, srv, "server-type", "config-example")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "serverClass") {
		t.Fatalf("expected a create example, got: %s", out)
	}

	out, err = runCLI(t, srv, "server-type", "update-config-example")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "label") {
		t.Fatalf("expected an update example, got: %s", out)
	}
}
