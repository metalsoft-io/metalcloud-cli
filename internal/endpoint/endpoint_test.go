package endpoint

import (
	"context"
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func makeEndpoint(id string, name string) map[string]any {
	return map[string]any{
		"id":               id,
		"revision":         "1",
		"siteId":           1,
		"name":             name,
		"label":            "ep-" + id,
		"createdTimestamp": "2024-01-01T00:00:00Z",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	}
}

func makeEndpointInterface(id float32) map[string]any {
	return map[string]any{
		"id":                         id,
		"revision":                   "1",
		"networkDeviceId":            10,
		"networkDeviceInterfaceId":   20,
		"networkDeviceInterfaceName": "eth0",
		"createdTimestamp":           "2024-01-01T00:00:00Z",
		"updatedTimestamp":           "2024-01-01T00:00:00Z",
	}
}

// TestEndpointList_HappyPath verifies a successful list call returns without error.
func TestEndpointList_HappyPath(t *testing.T) {
	page1 := testutils.PaginatedResponse([]any{makeEndpoint("1", "ep-one"), makeEndpoint("2", "ep-two")}, 1, 1)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/endpoints": testutils.JSONHandler(http.StatusOK, page1),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := EndpointList(ctx, nil, nil); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestEndpointList_ServerError verifies a 500 response is surfaced as an error.
func TestEndpointList_ServerError(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/endpoints": testutils.ErrorHandler(http.StatusInternalServerError, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := EndpointList(ctx, nil, nil); err == nil {
		t.Fatal("expected an error for HTTP 500, got nil")
	}
}

// TestEndpointList_EmptyList verifies an empty list succeeds.
func TestEndpointList_EmptyList(t *testing.T) {
	page := testutils.PaginatedResponse([]any{}, 1, 1)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/endpoints": testutils.JSONHandler(http.StatusOK, page),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := EndpointList(ctx, nil, nil); err != nil {
		t.Fatalf("expected nil error for empty list, got: %v", err)
	}
}

// TestEndpointList_Pagination verifies FetchAllPages fetches all 3 pages (205 items total).
// MultiPageServer wraps each page in PaginatedResponse itself, so pass raw item slices.
func TestEndpointList_Pagination(t *testing.T) {
	page1 := make([]any, 100)
	page2 := make([]any, 100)
	page3 := make([]any, 5)
	for i := range page1 {
		page1[i] = makeEndpoint("1", "ep")
	}
	for i := range page2 {
		page2[i] = makeEndpoint("2", "ep")
	}
	for i := range page3 {
		page3[i] = makeEndpoint("3", "ep")
	}
	ts := testutils.MultiPageServer("/api/v2/endpoints", []any{page1, page2, page3})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := EndpointList(ctx, nil, nil); err != nil {
		t.Fatalf("expected nil error during pagination, got: %v", err)
	}
}

// TestEndpointGet_HappyPath verifies successful retrieval of a single endpoint.
func TestEndpointGet_HappyPath(t *testing.T) {
	ep := makeEndpoint("42", "test-ep")
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/endpoints/42": testutils.JSONHandler(http.StatusOK, ep),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := EndpointGet(ctx, "42"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestEndpointGet_NotFound verifies a 404 response is surfaced as an error.
func TestEndpointGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/endpoints/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := EndpointGet(ctx, "99"); err == nil {
		t.Fatal("expected an error for HTTP 404, got nil")
	}
}

// TestEndpointGet_InvalidId verifies a non-numeric ID returns an error immediately.
func TestEndpointGet_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := EndpointGet(ctx, "not-a-number"); err == nil {
		t.Fatal("expected an error for invalid ID, got nil")
	}
}

// TestEndpointCreate_HappyPath verifies a successful endpoint creation.
func TestEndpointCreate_HappyPath(t *testing.T) {
	ep := makeEndpoint("1", "new-ep")
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/endpoints": testutils.JSONHandler(http.StatusCreated, ep),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	payload := sdk.CreateEndpoint{
		Name:   "new-ep",
		Label:  "new-ep",
		SiteId: 1,
	}
	if err := EndpointCreate(ctx, payload); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestEndpointDelete_HappyPath verifies a successful endpoint deletion.
func TestEndpointDelete_HappyPath(t *testing.T) {
	ep := makeEndpoint("5", "to-delete")
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/endpoints/5": func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				testutils.JSONHandler(http.StatusOK, ep)(w, r)
			case http.MethodDelete:
				w.WriteHeader(http.StatusNoContent)
			default:
				http.NotFound(w, r)
			}
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := EndpointDelete(ctx, "5"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestEndpointCreateBulk_NumericIds verifies bulk creation when interfaces are
// specified by numeric networkDeviceInterfaceId (no device lookup needed).
func TestEndpointCreateBulk_NumericIds(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/endpoints/actions/bulk-create": testutils.JSONHandler(http.StatusCreated, map[string]any{}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`[
		{"name": "h08", "label": "h08", "siteId": 1, "endpointInterfaces": [{"networkDeviceInterfaceId": 111}]},
		{"name": "h24", "label": "h24", "siteId": 1, "endpointInterfaces": [{"networkDeviceInterfaceId": 222}]}
	]`)

	if err := EndpointCreateBulk(ctx, config); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestEndpointCreateBulk_Empty verifies an empty list is rejected.
func TestEndpointCreateBulk_Empty(t *testing.T) {
	ctx := testutils.SetupTestContext("http://unused")
	if err := EndpointCreateBulk(ctx, []byte(`[]`)); err == nil {
		t.Fatal("expected an error for an empty endpoint list, got nil")
	}
}

// TestInterfaceResolver_NumericId verifies the numeric id fast-path returns the
// id verbatim without any device lookup.
func TestInterfaceResolver_NumericId(t *testing.T) {
	r := newInterfaceResolver()
	id := int64(4242)
	got, err := r.resolve(context.Background(), EndpointInterfaceInput{NetworkDeviceInterfaceId: &id})
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if got != id {
		t.Fatalf("expected %d, got %d", id, got)
	}
}

// TestInterfaceResolver_MissingFields verifies an interface with neither a
// numeric id nor a device/interface label pair is rejected.
func TestInterfaceResolver_MissingFields(t *testing.T) {
	r := newInterfaceResolver()

	iface := "swp1"
	cases := []EndpointInterfaceInput{
		{},                          // nothing
		{NetworkDevice: strPtr("")}, // device only, empty
		{Interface: &iface},         // interface only, no device
	}
	for i, tc := range cases {
		if _, err := r.resolve(context.Background(), tc); err == nil {
			t.Fatalf("case %d: expected an error, got nil", i)
		}
	}
}

func strPtr(s string) *string { return &s }

// TestEndpointInterfaceList_HappyPath verifies listing interfaces for an endpoint.
func TestEndpointInterfaceList_HappyPath(t *testing.T) {
	ifaces := testutils.PaginatedResponse([]any{makeEndpointInterface(1), makeEndpointInterface(2)}, 1, 1)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/endpoints/10/interfaces": testutils.JSONHandler(http.StatusOK, ifaces),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := EndpointInterfaceList(ctx, "10"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}
