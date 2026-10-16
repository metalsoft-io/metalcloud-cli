package server

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/spf13/viper"
)

// serverWithOperationalFieldsJSON is a server whose lifecycle status is
// "available" while its health is "failed" (the case the table could not show).
const serverWithOperationalFieldsJSON = `{
	"serverId": 236,
	"siteId": 2,
	"serverTypeId": 3,
	"serverUUID": "uuid-236",
	"serialNumber": "SN-236",
	"managementAddress": "10.0.0.236",
	"vendor": "Dell",
	"model": "R640",
	"serverStatus": "available",
	"revision": 1,
	"links": [],
	"serverMetricsMetadata": [],
	"healthState": "failed",
	"healthDetails": ["PSU 1 failed", "Fan 3 degraded"],
	"healthLastCheckedTimestamp": "2026-10-15T10:00:00Z",
	"healthMonitoringEnabled": true,
	"firmwareBaselineComplianceStatus": "non_compliant",
	"targetFirmwareBaselineId": 12,
	"currentFirmwareBaselineId": 9,
	"driftDetectionSyncStatus": "out_of_sync",
	"firmwareBackupEnabled": true,
	"firmwareDriftDetectionEnabled": false,
	"tags": ["rack-12", "row-A"]
}`

var operationalColumns = []string{
	"Health", "Health details", "Health checked", "FW compliance", "Target FW baseline",
	"Current FW baseline", "Drift sync", "FW backup", "FW drift detection", "Health monitoring", "Tags",
}

// captureStdout runs f with os.Stdout redirected and returns what was written.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}

	old := os.Stdout
	os.Stdout = w

	done := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()

	f()

	w.Close()
	os.Stdout = old

	return <-done
}

// renderCSV runs f with the CSV output format and returns the header and the first row as a map.
func renderCSV(t *testing.T, f func()) map[string]string {
	t.Helper()

	viper.Set(formatter.ConfigFormat, "csv")
	defer viper.Set(formatter.ConfigFormat, "json")

	out := captureStdout(t, f)

	records, err := csv.NewReader(strings.NewReader(out)).ReadAll()
	if err != nil {
		t.Fatalf("output is not valid CSV: %v\n%s", err, out)
	}
	if len(records) < 2 {
		t.Fatalf("expected a header and one row, got %q", out)
	}

	row := map[string]string{}
	for i, name := range records[0] {
		row[name] = records[1][i]
	}
	return row
}

func assertOperationalRow(t *testing.T, row map[string]string) {
	t.Helper()

	for _, column := range operationalColumns {
		if _, ok := row[column]; !ok {
			t.Errorf("missing column %q", column)
		}
	}

	expect := map[string]string{
		"Status":              "available",
		"Health":              "failed",
		"Health details":      "PSU 1 failed, Fan 3 degraded",
		"FW compliance":       "non_compliant",
		"Target FW baseline":  "12",
		"Current FW baseline": "9",
		"Drift sync":          "out_of_sync",
		"FW backup":           "true",
		"FW drift detection":  "false",
		"Health monitoring":   "true",
		"Tags":                "rack-12, row-A",
	}
	for column, want := range expect {
		if got := row[column]; got != want {
			t.Errorf("%s: expected %q, got %q", column, want, got)
		}
	}
	if row["Health checked"] == "" {
		t.Error("Health checked: expected a timestamp, got empty")
	}
}

func TestServerList_ShowsOperationalFields(t *testing.T) {
	body := `{"data":[` + serverWithOperationalFieldsJSON + `]}`
	ts := httptest.NewServer(newServerListHandler(body, http.StatusOK))
	defer ts.Close()

	row := renderCSV(t, func() {
		if err := ServerList(setupTestContext(ts.URL), ListFlags{}); err != nil {
			t.Fatalf("ServerList: unexpected error: %v", err)
		}
	})

	assertOperationalRow(t, row)
}

func TestServerGet_ShowsOperationalFields(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/servers/236") {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, serverWithOperationalFieldsJSON)
			return
		}
		http.NotFound(w, r)
	}))
	defer ts.Close()

	row := renderCSV(t, func() {
		if err := ServerGet(setupTestContext(ts.URL), "236", false); err != nil {
			t.Fatalf("ServerGet: unexpected error: %v", err)
		}
	})

	assertOperationalRow(t, row)
}

func TestServerList_OperationalFieldsAbsent(t *testing.T) {
	ts := httptest.NewServer(newServerListHandler(serverListLinksArrayJSON, http.StatusOK))
	defer ts.Close()

	row := renderCSV(t, func() {
		if err := ServerList(setupTestContext(ts.URL), ListFlags{}); err != nil {
			t.Fatalf("ServerList: unexpected error: %v", err)
		}
	})

	for _, column := range operationalColumns {
		got, ok := row[column]
		if !ok {
			t.Errorf("missing column %q", column)
		} else if got != "" {
			t.Errorf("%s: expected empty for a server without the field, got %q", column, got)
		}
	}
	if row["Status"] != "active" {
		t.Errorf("Status: expected active, got %q", row["Status"])
	}
}

func TestServerList_JsonOutputKeepsFullRecord(t *testing.T) {
	body := `{"data":[` + serverWithOperationalFieldsJSON + `]}`
	ts := httptest.NewServer(newServerListHandler(body, http.StatusOK))
	defer ts.Close()

	out := captureStdout(t, func() {
		if err := ServerList(setupTestContext(ts.URL), ListFlags{}); err != nil {
			t.Fatalf("ServerList: unexpected error: %v", err)
		}
	})

	for _, field := range []string{`"healthState":"failed"`, `"serverMetricsMetadata":[]`, `"tags":["rack-12","row-A"]`} {
		if !strings.Contains(strings.ReplaceAll(out, " ", ""), field) {
			t.Errorf("json output lost %s: %s", field, out)
		}
	}
}
