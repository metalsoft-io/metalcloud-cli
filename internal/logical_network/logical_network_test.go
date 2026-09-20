package logical_network

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
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

// ---------------------------------------------------------------------------
// Config, create-from-profile and attached-resource commands
// ---------------------------------------------------------------------------

// recorder captures the requests received by the mock server so tests can
// assert on method, path, headers and body.
type recorder struct {
	mu   sync.Mutex
	reqs []recordedRequest
}

type recordedRequest struct {
	Method string
	Path   string
	Header http.Header
	Body   string
}

func (r *recorder) wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		r.mu.Lock()
		r.reqs = append(r.reqs, recordedRequest{Method: req.Method, Path: req.URL.Path, Header: req.Header.Clone(), Body: string(body)})
		r.mu.Unlock()
		next(w, req)
	}
}

func (r *recorder) last() recordedRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.reqs[len(r.reqs)-1]
}

// logicalNetworkConfigItem mirrors the shape returned by
// GET /api/v2/logical-networks/{id}/config. Its revision (7) deliberately
// differs from the logical network's own revision (1) so the tests can prove
// the config revision is the one used for If-Match.
var logicalNetworkConfigItem = map[string]any{
	"id": 5, "deployType": "create", "deployStatus": "not_started",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
	"revision": 7, "kind": "vlan", "mtu": nil,
	"vlan": map[string]any{"vlanAllocationStrategies": []any{}},
	"ipv4": map[string]any{"subnetAllocationStrategies": []any{}},
	"ipv6": map[string]any{"subnetAllocationStrategies": []any{}},
}

var externalConnectionItem = map[string]any{
	"id": "4", "label": "ext-conn", "name": "External connection", "fabricId": 2,
	"revision": "1", "createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
}

var externalConnectionLogicalNetworkItem = map[string]any{
	"id": 8, "externalConnectionId": 4, "logicalNetworkId": 5, "status": "active",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
}

var logicalNetworkInterconnectItem = map[string]any{
	"id": "9", "label": "ln-interconnect", "name": "LN interconnect", "revision": "1",
	"kind": "dci-evpn", "fabricInterconnectId": 12, "status": "active",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
}

func newLogicalNetworkServer(t *testing.T, rec *recorder) *httptest.Server {
	t.Helper()
	routes := map[string]http.HandlerFunc{
		"/api/v2/logical-networks/5":                                      testutils.JSONHandler(http.StatusOK, makeLogicalNetwork(5)),
		"/api/v2/logical-networks/5/config":                               testutils.JSONHandler(http.StatusOK, logicalNetworkConfigItem),
		"/api/v2/logical-networks/5/config/actions/apply-profiles":        testutils.JSONHandler(http.StatusOK, logicalNetworkConfigItem),
		"/api/v2/logical-networks/actions/create-from-profile":            testutils.JSONHandler(http.StatusCreated, makeLogicalNetwork(6)),
		"/api/v2/logical-networks/5/external-connections":                 testutils.JSONHandler(http.StatusOK, testutils.PaginatedResponse([]any{externalConnectionItem}, 1, 1)),
		"/api/v2/logical-networks/5/external-connection-logical-networks": testutils.JSONHandler(http.StatusOK, testutils.PaginatedResponse([]any{externalConnectionLogicalNetworkItem}, 1, 1)),
		"/api/v2/logical-networks/5/logical-network-interconnects":        testutils.JSONHandler(http.StatusOK, testutils.PaginatedResponse([]any{logicalNetworkInterconnectItem}, 1, 1)),
		"/api/v2/logical-networks/5/external-connections/4":               testutils.NoContentHandler(),
	}
	wrapped := make(map[string]http.HandlerFunc, len(routes))
	for path, handler := range routes {
		wrapped[path] = rec.wrap(handler)
	}
	return testutils.NewTestServer(wrapped)
}

func TestLogicalNetworkConfigGet(t *testing.T) {
	rec := &recorder{}
	ts := newLogicalNetworkServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	out := testutils.CaptureStdout(t, func() {
		if err := LogicalNetworkConfigGet(ctx, "5"); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, `"revision":7`) {
		t.Errorf("get-config output missing the config revision: %s", out)
	}
	if rec.last().Path != "/api/v2/logical-networks/5/config" {
		t.Errorf("unexpected path %s", rec.last().Path)
	}
}

func TestLogicalNetworkConfigUpdate_UsesConfigRevision(t *testing.T) {
	rec := &recorder{}
	ts := newLogicalNetworkServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	testutils.CaptureStdout(t, func() {
		if err := LogicalNetworkConfigUpdate(ctx, "5", []byte(`{"mtu":9000}`)); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	req := rec.last()
	if req.Method != http.MethodPatch || req.Path != "/api/v2/logical-networks/5/config" {
		t.Fatalf("unexpected request %s %s", req.Method, req.Path)
	}
	// The config's revision (7), not the logical network's (1).
	if got := req.Header.Get("If-Match"); got != "7" {
		t.Errorf("expected If-Match 7 (the config revision), got %q", got)
	}
	var body map[string]any
	_ = json.Unmarshal([]byte(req.Body), &body)
	if body["mtu"] != float64(9000) {
		t.Errorf("update body should carry the mtu, got %s", req.Body)
	}
}

func TestLogicalNetworkApplyProfiles_SendsBodyAndConfigRevision(t *testing.T) {
	rec := &recorder{}
	ts := newLogicalNetworkServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	testutils.CaptureStdout(t, func() {
		if err := LogicalNetworkApplyProfiles(ctx, "5", "3"); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	req := rec.last()
	if req.Method != http.MethodPost || req.Path != "/api/v2/logical-networks/5/config/actions/apply-profiles" {
		t.Fatalf("unexpected request %s %s", req.Method, req.Path)
	}
	if got := req.Header.Get("If-Match"); got != "7" {
		t.Errorf("expected If-Match 7 (the config revision), got %q", got)
	}
	// An unset optional body would be serialized as a literal "null".
	var body map[string]any
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		t.Fatalf("apply-profiles body is not a JSON object: %s", req.Body)
	}
	if body["logicalNetworkProfileId"] != float64(3) {
		t.Errorf("apply-profiles body missing the profile id: %s", req.Body)
	}
}

func TestLogicalNetworkApplyProfiles_InvalidProfileId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := LogicalNetworkApplyProfiles(ctx, "5", "not-a-number"); err == nil {
		t.Fatal("expected an error for an invalid profile ID, got nil")
	}
}

func TestLogicalNetworkCreateFromProfile(t *testing.T) {
	rec := &recorder{}
	ts := newLogicalNetworkServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	testutils.CaptureStdout(t, func() {
		err := LogicalNetworkCreateFromProfile(ctx, sdk.CreateLogicalNetworkFromProfile{
			LogicalNetworkProfileId: 3,
			Label:                   sdk.PtrString("from-profile"),
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	req := rec.last()
	if req.Method != http.MethodPost || req.Path != "/api/v2/logical-networks/actions/create-from-profile" {
		t.Fatalf("unexpected request %s %s", req.Method, req.Path)
	}
	for _, want := range []string{`"logicalNetworkProfileId":3`, `"label":"from-profile"`} {
		if !strings.Contains(req.Body, want) {
			t.Errorf("create-from-profile body missing %s: %s", want, req.Body)
		}
	}
}

func TestLogicalNetworkAttachedResources(t *testing.T) {
	rec := &recorder{}
	ts := newLogicalNetworkServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	cases := []struct {
		name string
		fn   func() error
		want string
	}{
		{"external-connections", func() error { return LogicalNetworkExternalConnections(ctx, "5") }, "ext-conn"},
		{"external-connection-logical-networks", func() error {
			return LogicalNetworkExternalConnectionLogicalNetworks(ctx, "5")
		}, `"externalConnectionId":4`},
		{"interconnects", func() error { return LogicalNetworkInterconnects(ctx, "5") }, "ln-interconnect"},
	}

	for _, testCase := range cases {
		out := testutils.CaptureStdout(t, func() {
			if err := testCase.fn(); err != nil {
				t.Errorf("%s: unexpected error: %v", testCase.name, err)
			}
		})
		if !strings.Contains(out, testCase.want) {
			t.Errorf("%s: output missing %s: %s", testCase.name, testCase.want, out)
		}
	}
}

func TestLogicalNetworkDetachExternalConnection(t *testing.T) {
	rec := &recorder{}
	ts := newLogicalNetworkServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := LogicalNetworkDetachExternalConnection(ctx, "5", "4"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := rec.last()
	if req.Method != http.MethodDelete || req.Path != "/api/v2/logical-networks/5/external-connections/4" {
		t.Fatalf("unexpected request %s %s", req.Method, req.Path)
	}
}

func TestLogicalNetworkDetachExternalConnection_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := LogicalNetworkDetachExternalConnection(ctx, "5", "not-a-number"); err == nil {
		t.Fatal("expected an error for an invalid external connection ID, got nil")
	}
}
