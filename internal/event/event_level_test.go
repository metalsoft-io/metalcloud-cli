package event

import (
	"encoding/csv"
	"net/http"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/spf13/viper"
)

// makeLevelEvent returns an event as the API now sends it: the level is in `level`
// and there is no `severity` field at all.
func makeLevelEvent(id, level string) map[string]any {
	return map[string]any{
		"id":                id,
		"type":              "infrastructure_instances_info",
		"level":             level,
		"visibility":        "public",
		"title":             "Event " + id,
		"message":           "something happened",
		"occurredTimestamp": "2024-01-01T00:00:00Z",
	}
}

// renderCSV runs fn with CSV output and returns the header and the data rows.
func renderCSV(t *testing.T, fn func()) ([]string, [][]string) {
	t.Helper()

	viper.Set(formatter.ConfigFormat, "csv")
	defer viper.Set(formatter.ConfigFormat, "json")

	out := testutils.CaptureStdout(t, fn)

	records, err := csv.NewReader(strings.NewReader(out)).ReadAll()
	if err != nil || len(records) < 2 {
		t.Fatalf("unexpected CSV output %q (%v)", out, err)
	}
	return records[0], records[1:]
}

func columnIndex(t *testing.T, header []string, name string) int {
	t.Helper()

	for i, h := range header {
		if h == name {
			return i
		}
	}
	t.Fatalf("column %q not found in %v", name, header)
	return -1
}

func TestEventList_ShowsLevelInTableOutput(t *testing.T) {
	page := testutils.PaginatedResponse([]any{
		makeLevelEvent("1", "info"),
		makeLevelEvent("2", "warn"),
		makeLevelEvent("3", "error"),
	}, 1, 1)
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

	level := columnIndex(t, header, "Level")
	for _, h := range header {
		if h == "Severity" {
			t.Error("the Severity column was renamed to Level")
		}
	}

	want := []string{"info", "warn", "error"}
	if len(rows) != len(want) {
		t.Fatalf("expected %d rows, got %d", len(want), len(rows))
	}
	for i, row := range rows {
		if row[level] != want[i] {
			t.Errorf("row %d: expected level %q, got %q", i, want[i], row[level])
		}
	}
}

func TestEventGet_ShowsLevel(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/events/5": testutils.JSONHandler(200, makeLevelEvent("5", "error")),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)

	header, rows := renderCSV(t, func() {
		if err := EventGet(ctx, "5"); err != nil {
			t.Fatalf("EventGet() unexpected error: %v", err)
		}
	})

	if got := rows[0][columnIndex(t, header, "Level")]; got != "error" {
		t.Errorf("expected level %q, got %q", "error", got)
	}
}
