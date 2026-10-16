package event

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/spf13/viper"
)

func makeFullEvent() map[string]any {
	return map[string]any{
		"id":                        "42",
		"type":                      "infrastructure_instances_info",
		"level":                     "error",
		"visibility":                "public",
		"title":                     "Server wiped",
		"message":                   "Server 236 was wiped and returned to the pool",
		"occurredTimestamp":         "2024-01-01T00:00:00Z",
		"userId":                    "11",
		"userIdAuthenticated":       "12",
		"userEmailAuthenticated":    "ops@example.com",
		"remoteIpAddress":           "10.1.2.3",
		"infrastructureId":          "100",
		"serverId":                  "236",
		"networkDeviceId":           "7",
		"networkDeviceControllerId": "8",
		"storageId":                 "9",
		"vmPoolId":                  "10",
		"siteId":                    "2",
		"jobId":                     "5000",
		"context":                   map[string]any{"reason": "decommission"},
	}
}

func TestEventList_CsvIsAnAuditExport(t *testing.T) {
	page := testutils.PaginatedResponse([]any{makeFullEvent()}, 1, 1)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/events": testutils.JSONHandler(200, page),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)

	header, rows := renderCSV(t, func() {
		if err := EventList(ctx, ListFlags{}); err != nil {
			t.Fatalf("EventList() unexpected error: %v", err)
		}
	})

	expect := map[string]string{
		"ID":                           "42",
		"Level":                        "error",
		"Title":                        "Server wiped",
		"Message":                      "Server 236 was wiped and returned to the pool",
		"User ID":                      "11",
		"Authenticated User ID":        "12",
		"Infrastructure ID":            "100",
		"Server ID":                    "236",
		"Network Device ID":            "7",
		"Network Device Controller ID": "8",
		"Storage ID":                   "9",
		"VM Pool ID":                   "10",
		"Site ID":                      "2",
		"Job ID":                       "5000",
	}
	for column, want := range expect {
		if got := rows[0][columnIndex(t, header, column)]; got != want {
			t.Errorf("%s: expected %q, got %q", column, want, got)
		}
	}
}

func TestEventList_OptionalIdsAreEmptyWhenAbsent(t *testing.T) {
	page := testutils.PaginatedResponse([]any{makeLevelEvent("1", "info")}, 1, 1)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/events": testutils.JSONHandler(200, page),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)

	header, rows := renderCSV(t, func() {
		if err := EventList(ctx, ListFlags{}); err != nil {
			t.Fatalf("EventList() unexpected error: %v", err)
		}
	})

	for _, column := range []string{"User ID", "Authenticated User ID", "Infrastructure ID", "Server ID",
		"Network Device ID", "Network Device Controller ID", "Storage ID", "VM Pool ID", "Site ID", "Job ID"} {
		if got := rows[0][columnIndex(t, header, column)]; got != "" {
			t.Errorf("%s: expected empty for an event without it, got %q", column, got)
		}
	}
}

func TestEventList_NumericIdsAreTolerated(t *testing.T) {
	event := makeLevelEvent("1", "info")
	event["serverId"] = 236
	event["userId"] = 11

	page := testutils.PaginatedResponse([]any{event}, 1, 1)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/events": testutils.JSONHandler(200, page),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)

	header, rows := renderCSV(t, func() {
		if err := EventList(ctx, ListFlags{}); err != nil {
			t.Fatalf("EventList() unexpected error: %v", err)
		}
	})

	if got := rows[0][columnIndex(t, header, "Server ID")]; got != "236" {
		t.Errorf("Server ID: expected 236, got %q", got)
	}
	if got := rows[0][columnIndex(t, header, "User ID")]; got != "11" {
		t.Errorf("User ID: expected 11, got %q", got)
	}
}

func TestEventList_JsonKeepsFullRecord(t *testing.T) {
	page := testutils.PaginatedResponse([]any{makeFullEvent()}, 1, 1)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/events": testutils.JSONHandler(200, page),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)

	out := testutils.CaptureStdout(t, func() {
		if err := EventList(ctx, ListFlags{}); err != nil {
			t.Fatalf("EventList() unexpected error: %v", err)
		}
	})

	var records []map[string]any
	if err := json.Unmarshal([]byte(out), &records); err != nil || len(records) != 1 {
		t.Fatalf("unexpected json output %q (%v)", out, err)
	}
	for _, field := range []string{"context", "userEmailAuthenticated", "remoteIpAddress", "message"} {
		if records[0][field] == nil {
			t.Errorf("json output lost %q", field)
		}
	}
}

func TestEventList_SendsUserAuthenticatedAndNetworkDeviceFilters(t *testing.T) {
	var mu sync.Mutex
	var queries []url.Values

	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/events": func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			queries = append(queries, r.URL.Query())
			mu.Unlock()
			testutils.JSONHandler(200, testutils.PaginatedResponse([]any{}, 1, 1))(w, r)
		},
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)

	err := EventList(ctx, ListFlags{
		FilterUserId:              []string{"11"},
		FilterUserIdAuthenticated: []string{"12"},
		FilterServerId:            []string{"236"},
		FilterNetworkDeviceId:     []string{"7"},
	})
	if err != nil {
		t.Fatalf("EventList() unexpected error: %v", err)
	}

	if len(queries) == 0 {
		t.Fatal("no request recorded")
	}

	for key, want := range map[string]string{
		"filter.userId":              "11",
		"filter.userIdAuthenticated": "12",
		"filter.serverId":            "236",
		"filter.networkDeviceId":     "7",
	} {
		if got := strings.Join(queries[0][key], ","); got != want {
			t.Errorf("%s: expected %q, got %q", key, want, got)
		}
	}
}

// renderMarkdownHeader runs fn with the markdown table format and returns the upper-cased header cells.
func renderMarkdownHeader(t *testing.T, fn func()) []string {
	t.Helper()

	viper.Set(formatter.ConfigFormat, "md")
	defer viper.Set(formatter.ConfigFormat, "json")

	out := testutils.CaptureStdout(t, fn)

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 || !strings.HasPrefix(lines[0], "|") {
		t.Fatalf("unexpected markdown output %q", out)
	}

	var header []string
	for _, cell := range strings.Split(strings.Trim(lines[0], "|"), "|") {
		header = append(header, strings.ToUpper(strings.TrimSpace(cell)))
	}
	return header
}

func hasColumn(header []string, name string) bool {
	for _, h := range header {
		if h == strings.ToUpper(name) {
			return true
		}
	}
	return false
}

func TestEventList_TableHidesEmptyOptionalColumns(t *testing.T) {
	onlyServer := makeLevelEvent("1", "info")
	onlyServer["serverId"] = "236"
	onlyServer["userId"] = "11"
	other := makeLevelEvent("2", "warn")
	other["serverId"] = "237"

	page := testutils.PaginatedResponse([]any{onlyServer, other}, 1, 1)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/events": testutils.JSONHandler(200, page),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)

	header := renderMarkdownHeader(t, func() {
		if err := EventList(ctx, ListFlags{}); err != nil {
			t.Fatalf("EventList() unexpected error: %v", err)
		}
	})

	for _, name := range []string{"ID", "Type", "Level", "Visibility", "Time", "Title", "Message", "Server ID", "User ID"} {
		if !hasColumn(header, name) {
			t.Errorf("expected column %q in %v", name, header)
		}
	}
	for _, name := range []string{"Authenticated User ID", "Infrastructure ID", "Network Device ID",
		"Network Device Controller ID", "Storage ID", "VM Pool ID", "Site ID", "Job ID"} {
		if hasColumn(header, name) {
			t.Errorf("column %q is empty for every event and should be hidden, got %v", name, header)
		}
	}
}

func TestEventList_TableShowsColumnWhenAnyEventHasIt(t *testing.T) {
	plain := makeLevelEvent("1", "info")
	withJob := makeLevelEvent("2", "info")
	withJob["jobId"] = "5000"

	page := testutils.PaginatedResponse([]any{plain, withJob}, 1, 1)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/events": testutils.JSONHandler(200, page),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)

	header := renderMarkdownHeader(t, func() {
		if err := EventList(ctx, ListFlags{}); err != nil {
			t.Fatalf("EventList() unexpected error: %v", err)
		}
	})

	if !hasColumn(header, "Job ID") {
		t.Errorf("one event has a job, expected the Job ID column in %v", header)
	}
}

func TestEventGet_TableShowsOnlyThePresentIds(t *testing.T) {
	event := makeLevelEvent("9", "error")
	event["storageId"] = "3"

	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/events/9": testutils.JSONHandler(200, event),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)

	header := renderMarkdownHeader(t, func() {
		if err := EventGet(ctx, "9"); err != nil {
			t.Fatalf("EventGet() unexpected error: %v", err)
		}
	})

	if !hasColumn(header, "Storage ID") || hasColumn(header, "Server ID") {
		t.Errorf("expected only Storage ID among the ID columns, got %v", header)
	}
}

func TestEventList_CsvKeepsAllColumnsWhateverTheData(t *testing.T) {
	page := testutils.PaginatedResponse([]any{makeLevelEvent("1", "info")}, 1, 1)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/events": testutils.JSONHandler(200, page),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)

	header, _ := renderCSV(t, func() {
		if err := EventList(ctx, ListFlags{}); err != nil {
			t.Fatalf("EventList() unexpected error: %v", err)
		}
	})

	for _, name := range []string{"User ID", "Authenticated User ID", "Infrastructure ID", "Server ID", "Network Device ID",
		"Network Device Controller ID", "Storage ID", "VM Pool ID", "Site ID", "Job ID"} {
		columnIndex(t, header, name)
	}
}
