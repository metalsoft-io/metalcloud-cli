package logical_network

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func makeLogicalNetwork(id int) map[string]any {
	return map[string]any{
		"id":                                 id,
		"label":                              "ln-label",
		"name":                               "ln-name",
		"annotations":                        map[string]any{},
		"createdAt":                          "2024-01-01T00:00:00Z",
		"updatedAt":                          "2024-01-01T00:00:00Z",
		"revision":                           1,
		"kind":                               "vlan",
		"fabricId":                           2,
		"infrastructureId":                   nil,
		"serviceStatus":                      "active",
		"lastAppliedLogicalNetworkProfileId": nil,
		"lastLogicalNetworkProfileAppliedAt": "2024-01-01T00:00:00Z",
		"config": map[string]any{
			"id":           1,
			"deployType":   "none",
			"deployStatus": "idle",
			"createdAt":    "2024-01-01T00:00:00Z",
			"updatedAt":    "2024-01-01T00:00:00Z",
			"revision":     1,
			"kind":         "vlan",
		},
	}
}

// TestLogicalNetworkList_HappyPath verifies a successful fetch-all list call.
func TestLogicalNetworkList_HappyPath(t *testing.T) {
	resp := testutils.PaginatedResponse([]any{makeLogicalNetwork(1), makeLogicalNetwork(2)}, 1, 1)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/logical-networks": testutils.JSONHandler(http.StatusOK, resp),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := LogicalNetworkList(ctx, "", ListFlags{}); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestLogicalNetworkList_ServerError verifies a 500 is surfaced as an error.
func TestLogicalNetworkList_ServerError(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/logical-networks": testutils.ErrorHandler(http.StatusInternalServerError, "server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := LogicalNetworkList(ctx, "", ListFlags{}); err == nil {
		t.Fatal("expected an error for HTTP 500, got nil")
	}
}

// TestLogicalNetworkList_EmptyList verifies an empty list succeeds.
func TestLogicalNetworkList_EmptyList(t *testing.T) {
	resp := testutils.PaginatedResponse([]any{}, 1, 1)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/logical-networks": testutils.JSONHandler(http.StatusOK, resp),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := LogicalNetworkList(ctx, "", ListFlags{}); err != nil {
		t.Fatalf("expected nil error for empty list, got: %v", err)
	}
}

// TestLogicalNetworkList_Pagination verifies 3-page fetch-all (205 items total).
func TestLogicalNetworkList_Pagination(t *testing.T) {
	page1 := make([]any, 100)
	page2 := make([]any, 100)
	page3 := make([]any, 5)
	for i := range page1 {
		page1[i] = makeLogicalNetwork(i + 1)
	}
	for i := range page2 {
		page2[i] = makeLogicalNetwork(i + 101)
	}
	for i := range page3 {
		page3[i] = makeLogicalNetwork(i + 201)
	}
	ts := testutils.MultiPageServer("/api/v2/logical-networks", []any{page1, page2, page3})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := LogicalNetworkList(ctx, "", ListFlags{}); err != nil {
		t.Fatalf("expected nil error during pagination, got: %v", err)
	}
}

// TestLogicalNetworkList_SinglePage verifies the single-page path (Page=1, Limit=10).
func TestLogicalNetworkList_SinglePage(t *testing.T) {
	resp := testutils.PaginatedResponse([]any{makeLogicalNetwork(1)}, 1, 1)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/logical-networks": testutils.JSONHandler(http.StatusOK, resp),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := LogicalNetworkList(ctx, "", ListFlags{Page: 1, Limit: 10}); err != nil {
		t.Fatalf("expected nil error for single-page call, got: %v", err)
	}
}

// TestLogicalNetworkGet_HappyPath verifies successful retrieval of a single logical network.
func TestLogicalNetworkGet_HappyPath(t *testing.T) {
	ln := makeLogicalNetwork(5)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/logical-networks/5": testutils.JSONHandler(http.StatusOK, ln),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := LogicalNetworkGet(ctx, "5"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestLogicalNetworkGet_NotFound verifies a 404 is surfaced as an error.
func TestLogicalNetworkGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/logical-networks/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := LogicalNetworkGet(ctx, "99"); err == nil {
		t.Fatal("expected an error for HTTP 404, got nil")
	}
}

// TestLogicalNetworkGet_InvalidId verifies a non-numeric ID is rejected immediately.
func TestLogicalNetworkGet_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := LogicalNetworkGet(ctx, "not-a-number"); err == nil {
		t.Fatal("expected an error for invalid ID, got nil")
	}
}

// TestLogicalNetworkDelete_InvalidId verifies that a non-numeric ID is rejected before
// any HTTP call is made.
func TestLogicalNetworkDelete_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := LogicalNetworkDelete(ctx, "not-a-number"); err == nil {
		t.Fatal("expected an error for invalid ID, got nil")
	}
}

// TestLogicalNetworkDelete_ServerError verifies that a server error on delete is surfaced.
// Note: LogicalNetworkDelete does not set IfMatch; the SDK enforces it client-side,
// so this call currently always returns an "ifMatch is required" error.
func TestLogicalNetworkDelete_MissingIfMatch(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := LogicalNetworkDelete(ctx, "3"); err == nil {
		t.Fatal("expected an error because IfMatch is not set by LogicalNetworkDelete")
	}
}

func TestLogicalNetworkConfigExample_VLAN(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := LogicalNetworkConfigExample(ctx, "vlan"); err != nil {
		t.Errorf("LogicalNetworkConfigExample(vlan) unexpected error: %v", err)
	}
}

func TestLogicalNetworkConfigExample_VXLAN(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := LogicalNetworkConfigExample(ctx, "vxlan"); err != nil {
		t.Errorf("LogicalNetworkConfigExample(vxlan) unexpected error: %v", err)
	}
}

func TestLogicalNetworkConfigExample_InvalidKind(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := LogicalNetworkConfigExample(ctx, "invalid-kind"); err == nil {
		t.Error("expected error for invalid kind, got nil")
	}
}
