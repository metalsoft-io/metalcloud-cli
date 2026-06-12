package logical_network_profile

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func makeLogicalNetworkProfile(id int) map[string]any {
	return map[string]any{
		"id":          id,
		"label":       "lnp-label",
		"name":        "lnp-name",
		"annotations": map[string]any{},
		"createdAt":   "2024-01-01T00:00:00Z",
		"updatedAt":   "2024-01-01T00:00:00Z",
		"revision":    1,
		"kind":        "vlan",
		"fabricId":    2,
	}
}

// TestLogicalNetworkProfileList_HappyPath verifies a successful list call.
func TestLogicalNetworkProfileList_HappyPath(t *testing.T) {
	resp := testutils.PaginatedResponse([]any{makeLogicalNetworkProfile(1), makeLogicalNetworkProfile(2)}, 1, 1)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/logical-network-profiles": testutils.JSONHandler(http.StatusOK, resp),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := LogicalNetworkProfileList(ctx, ListFlags{}); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestLogicalNetworkProfileList_ServerError verifies a 500 is surfaced as an error.
func TestLogicalNetworkProfileList_ServerError(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/logical-network-profiles": testutils.ErrorHandler(http.StatusInternalServerError, "server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := LogicalNetworkProfileList(ctx, ListFlags{}); err == nil {
		t.Fatal("expected an error for HTTP 500, got nil")
	}
}

// TestLogicalNetworkProfileList_EmptyList verifies an empty list succeeds.
func TestLogicalNetworkProfileList_EmptyList(t *testing.T) {
	resp := testutils.PaginatedResponse([]any{}, 1, 1)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/logical-network-profiles": testutils.JSONHandler(http.StatusOK, resp),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := LogicalNetworkProfileList(ctx, ListFlags{}); err != nil {
		t.Fatalf("expected nil error for empty list, got: %v", err)
	}
}

// TestLogicalNetworkProfileList_Pagination verifies 3-page fetch-all (205 items total).
func TestLogicalNetworkProfileList_Pagination(t *testing.T) {
	page1 := make([]any, 100)
	page2 := make([]any, 100)
	page3 := make([]any, 5)
	for i := range page1 {
		page1[i] = makeLogicalNetworkProfile(i + 1)
	}
	for i := range page2 {
		page2[i] = makeLogicalNetworkProfile(i + 101)
	}
	for i := range page3 {
		page3[i] = makeLogicalNetworkProfile(i + 201)
	}
	ts := testutils.MultiPageServer("/api/v2/logical-network-profiles", []any{page1, page2, page3})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := LogicalNetworkProfileList(ctx, ListFlags{}); err != nil {
		t.Fatalf("expected nil error during pagination, got: %v", err)
	}
}

// TestLogicalNetworkProfileGet_HappyPath verifies successful retrieval of a single profile.
func TestLogicalNetworkProfileGet_HappyPath(t *testing.T) {
	profile := makeLogicalNetworkProfile(7)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/logical-network-profiles/7": testutils.JSONHandler(http.StatusOK, profile),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := LogicalNetworkProfileGet(ctx, "7"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestLogicalNetworkProfileGet_NotFound verifies a 404 is surfaced as an error.
func TestLogicalNetworkProfileGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/logical-network-profiles/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := LogicalNetworkProfileGet(ctx, "99"); err == nil {
		t.Fatal("expected an error for HTTP 404, got nil")
	}
}

// TestLogicalNetworkProfileGet_InvalidId verifies a non-numeric ID is rejected immediately.
func TestLogicalNetworkProfileGet_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := LogicalNetworkProfileGet(ctx, "bad-id"); err == nil {
		t.Fatal("expected an error for invalid ID, got nil")
	}
}

// TestLogicalNetworkProfileDelete_HappyPath verifies successful deletion of a profile.
func TestLogicalNetworkProfileDelete_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/logical-network-profiles/4": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			http.NotFound(w, r)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := LogicalNetworkProfileDelete(ctx, "4"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}
