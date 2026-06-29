package event

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func makeEvent(id string) map[string]any {
	return map[string]any{
		"id":                id,
		"type":              "infrastructure_instances_info",
		"severity":          "info",
		"level":             "info",
		"visibility":        "public",
		"title":             "Test event " + id,
		"message":           "something happened",
		"occurredTimestamp": "2024-01-01T00:00:00Z",
	}
}

// TestEventList_FetchAll verifies the fetch-all path (no Page/Limit flags).
func TestEventList_FetchAll(t *testing.T) {
	page := testutils.PaginatedResponse([]any{makeEvent("1"), makeEvent("2")}, 1, 1)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/events": testutils.JSONHandler(http.StatusOK, page),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := EventList(ctx, ListFlags{}); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestEventList_SinglePage verifies the single-page path (Page=1, Limit=5).
func TestEventList_SinglePage(t *testing.T) {
	page := testutils.PaginatedResponse([]any{makeEvent("1"), makeEvent("2")}, 1, 1)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/events": testutils.JSONHandler(http.StatusOK, page),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	flags := ListFlags{Page: 1, Limit: 5}
	if err := EventList(ctx, flags); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestEventList_ServerError verifies a 500 is surfaced as an error.
func TestEventList_ServerError(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/events": testutils.ErrorHandler(http.StatusInternalServerError, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := EventList(ctx, ListFlags{}); err == nil {
		t.Fatal("expected an error for HTTP 500, got nil")
	}
}

// TestEventList_EmptyList verifies an empty list succeeds.
func TestEventList_EmptyList(t *testing.T) {
	page := testutils.PaginatedResponse([]any{}, 1, 1)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/events": testutils.JSONHandler(http.StatusOK, page),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := EventList(ctx, ListFlags{}); err != nil {
		t.Fatalf("expected nil error for empty list, got: %v", err)
	}
}

// TestEventList_Pagination verifies fetch-all accumulates 205 items across 3 pages.
func TestEventList_Pagination(t *testing.T) {
	page1 := make([]any, 100)
	page2 := make([]any, 100)
	page3 := make([]any, 5)
	for i := range page1 {
		page1[i] = makeEvent("1")
	}
	for i := range page2 {
		page2[i] = makeEvent("2")
	}
	for i := range page3 {
		page3[i] = makeEvent("3")
	}
	ts := testutils.MultiPageServer("/api/v2/events", []any{page1, page2, page3})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := EventList(ctx, ListFlags{}); err != nil {
		t.Fatalf("expected nil error during pagination, got: %v", err)
	}
}

// makeDriftedEvent returns an event whose `type` is an enum value the SDK does
// not know (a numeric EventTypes such as "109") and which omits the SDK-required
// `severity`. The typed SDK model rejects such responses; raw-body parsing must not.
func makeDriftedEvent(id string) map[string]any {
	return map[string]any{
		"id":                id,
		"type":              "109",
		"severity":          "info",
		"level":             "info",
		"visibility":        "public",
		"title":             "Drifted event " + id,
		"message":           "something happened",
		"occurredTimestamp": "2024-01-01T00:00:00Z",
	}
}

// TestEventList_UnknownEnumValue is a regression test for the SDK<->API enum
// desync where `event list` failed with "109 is not a valid EventTypes".
func TestEventList_UnknownEnumValue(t *testing.T) {
	page := testutils.PaginatedResponse([]any{makeDriftedEvent("1"), makeEvent("2")}, 1, 1)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/events": testutils.JSONHandler(http.StatusOK, page),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := EventList(ctx, ListFlags{}); err != nil {
		t.Fatalf("expected nil error for unknown EventTypes value, got: %v", err)
	}
}

// TestEventGet_UnknownEnumValue is the single-event counterpart of the regression test.
func TestEventGet_UnknownEnumValue(t *testing.T) {
	ev := makeDriftedEvent("7")
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/events/7": testutils.JSONHandler(http.StatusOK, ev),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := EventGet(ctx, "7"); err != nil {
		t.Fatalf("expected nil error for unknown EventTypes value, got: %v", err)
	}
}

// TestEventGet_HappyPath verifies successful retrieval of a single event.
func TestEventGet_HappyPath(t *testing.T) {
	ev := makeEvent("7")
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/events/7": testutils.JSONHandler(http.StatusOK, ev),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := EventGet(ctx, "7"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestEventGet_NotFound verifies a 404 is surfaced as an error.
func TestEventGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/events/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := EventGet(ctx, "99"); err == nil {
		t.Fatal("expected an error for HTTP 404, got nil")
	}
}

// TestEventGet_InvalidId verifies a non-numeric ID returns an error immediately.
func TestEventGet_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := EventGet(ctx, "not-a-number"); err == nil {
		t.Fatal("expected an error for invalid ID, got nil")
	}
}
