package fabric

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

// ---- fixtures -------------------------------------------------------------

var subResFabricItem = map[string]any{
	"id": "1", "name": "fabric-1", "revision": "1", "siteId": 10,
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
	"fabricConfiguration": map[string]any{"fabricType": "ethernet"},
}

var subResBgpSessionItem = map[string]any{
	"id": 7, "networkFabricId": 1, "name": "session-7", "status": "active",
	"bgpNumbering": "numbered", "bgpLinkConfiguration": "active",
	"networkFabricLinkId": 3, "revision": "2",
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
	"customVariables": map[string]any{"example_variable": "example_value"},
}

var subResLinkAggregationItem = map[string]any{
	"id": 4, "networkFabricId": 1, "name": "lag-4", "type": "lag", "status": "active",
	"revision": "3", "createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
}

var subResLinkItem = map[string]any{
	"id": 3, "networkFabricId": 1, "linkType": "leaf_spine", "status": "active",
	"revision": "1", "createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
}

var subResInterconnectItem = map[string]any{
	"id": "12", "interconnectType": "dci-evpn", "label": "dc1-dc2", "name": "DC1 to DC2",
	"bgpConfigurationTemplateId": 3, "revision": "7", "status": "draft",
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
}

// ---- recorder -------------------------------------------------------------

type subResRecorder struct {
	mu   sync.Mutex
	reqs []subResRequest
}

type subResRequest struct {
	Method string
	Path   string
	Body   string
}

func (r *subResRecorder) wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		r.mu.Lock()
		r.reqs = append(r.reqs, subResRequest{Method: req.Method, Path: req.URL.Path, Body: string(body)})
		r.mu.Unlock()
		next(w, req)
	}
}

func (r *subResRecorder) last() subResRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.reqs[len(r.reqs)-1]
}

func newSubResServer(t *testing.T, rec *subResRecorder) *httptest.Server {
	t.Helper()
	routes := map[string]http.HandlerFunc{
		"/api/v2/network-fabrics": testutils.JSONHandler(200, testutils.PaginatedResponse([]any{subResFabricItem}, 1, 1)),
		"/api/v2/network-fabrics/1": methodRouter(map[string]http.HandlerFunc{
			http.MethodGet:    testutils.JSONHandler(200, subResFabricItem),
			http.MethodDelete: testutils.NoContentHandler(),
		}),
		"/api/v2/network-fabrics/1/actions/accept-deploy": testutils.NoContentHandler(),
		"/api/v2/network-fabrics/1/actions/reject-deploy": testutils.NoContentHandler(),
		"/api/v2/network-fabrics/1/links/3":               testutils.JSONHandler(200, subResLinkItem),
		"/api/v2/network-fabrics/1/network-fabric-interconnects": testutils.JSONHandler(200,
			testutils.PaginatedResponse([]any{subResInterconnectItem}, 1, 1)),
		"/api/v2/network-fabrics/1/bgp-sessions": methodRouter(map[string]http.HandlerFunc{
			http.MethodGet:  testutils.JSONHandler(200, testutils.PaginatedResponse([]any{subResBgpSessionItem}, 1, 1)),
			http.MethodPost: testutils.JSONHandler(201, subResBgpSessionItem),
		}),
		"/api/v2/network-fabrics/1/bgp-sessions/7": methodRouter(map[string]http.HandlerFunc{
			http.MethodGet:    testutils.JSONHandler(200, subResBgpSessionItem),
			http.MethodPatch:  testutils.JSONHandler(200, subResBgpSessionItem),
			http.MethodPut:    testutils.JSONHandler(200, subResBgpSessionItem),
			http.MethodDelete: testutils.NoContentHandler(),
		}),
		"/api/v2/network-fabrics/1/link-aggregations": methodRouter(map[string]http.HandlerFunc{
			http.MethodGet:  testutils.JSONHandler(200, testutils.PaginatedResponse([]any{subResLinkAggregationItem}, 1, 1)),
			http.MethodPost: testutils.JSONHandler(201, subResLinkAggregationItem),
		}),
		"/api/v2/network-fabrics/1/link-aggregations/4": methodRouter(map[string]http.HandlerFunc{
			http.MethodGet:    testutils.JSONHandler(200, subResLinkAggregationItem),
			http.MethodPatch:  testutils.JSONHandler(200, subResLinkAggregationItem),
			http.MethodPut:    testutils.JSONHandler(200, subResLinkAggregationItem),
			http.MethodDelete: testutils.NoContentHandler(),
		}),
	}
	wrapped := make(map[string]http.HandlerFunc, len(routes))
	for path, handler := range routes {
		wrapped[path] = rec.wrap(handler)
	}
	return testutils.NewTestServer(wrapped)
}

// methodRouter dispatches one path to a handler per HTTP method.
func methodRouter(handlers map[string]http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler, ok := handlers[r.Method]
		if !ok {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}

func runSubRes(t *testing.T, fn func() error) string {
	t.Helper()
	var out string
	var err error
	out = testutils.CaptureStdout(t, func() { err = fn() })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return out
}

// ---- fabric-level operations ---------------------------------------------

func TestFabricDeployDecisionsAndDelete(t *testing.T) {
	rec := &subResRecorder{}
	srv := newSubResServer(t, rec)
	defer srv.Close()
	ctx := setupTestContext(srv.URL)

	if err := FabricAcceptDeploy(ctx, "1"); err != nil {
		t.Fatalf("accept-deploy: %v", err)
	}
	if got := rec.last(); got.Method != http.MethodPost || got.Path != "/api/v2/network-fabrics/1/actions/accept-deploy" {
		t.Errorf("accept-deploy hit %s %s", got.Method, got.Path)
	}

	if err := FabricRejectDeploy(ctx, "fabric-1"); err != nil {
		t.Fatalf("reject-deploy: %v", err)
	}
	if got := rec.last(); got.Method != http.MethodPost || got.Path != "/api/v2/network-fabrics/1/actions/reject-deploy" {
		t.Errorf("reject-deploy hit %s %s", got.Method, got.Path)
	}

	if err := FabricDelete(ctx, "1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if got := rec.last(); got.Method != http.MethodDelete || got.Path != "/api/v2/network-fabrics/1" {
		t.Errorf("delete hit %s %s", got.Method, got.Path)
	}
}

func TestFabricLinkGet(t *testing.T) {
	rec := &subResRecorder{}
	srv := newSubResServer(t, rec)
	defer srv.Close()
	ctx := setupTestContext(srv.URL)

	out := runSubRes(t, func() error { return FabricLinkGet(ctx, "1", "3") })
	if !strings.Contains(out, `"linkType":"leaf_spine"`) {
		t.Errorf("get-link output missing link: %s", out)
	}
	if got := rec.last(); got.Path != "/api/v2/network-fabrics/1/links/3" {
		t.Errorf("get-link hit %s", got.Path)
	}

	if err := FabricLinkGet(ctx, "1", "not-a-number"); err == nil {
		t.Error("expected an error for a non-numeric link id")
	}
}

func TestFabricInterconnectsGet(t *testing.T) {
	rec := &subResRecorder{}
	srv := newSubResServer(t, rec)
	defer srv.Close()
	ctx := setupTestContext(srv.URL)

	out := runSubRes(t, func() error { return FabricInterconnectsGet(ctx, "1", nil) })
	if !strings.Contains(out, "dc1-dc2") {
		t.Errorf("get-interconnects output missing interconnect: %s", out)
	}
	if got := rec.last(); got.Path != "/api/v2/network-fabrics/1/network-fabric-interconnects" {
		t.Errorf("get-interconnects hit %s", got.Path)
	}
}

// ---- BGP sessions ---------------------------------------------------------

func TestBgpSessionCrud(t *testing.T) {
	rec := &subResRecorder{}
	srv := newSubResServer(t, rec)
	defer srv.Close()
	ctx := setupTestContext(srv.URL)

	out := runSubRes(t, func() error { return BgpSessionList(ctx, "1") })
	if !strings.Contains(out, `"name":"session-7"`) {
		t.Errorf("list output missing session: %s", out)
	}

	out = runSubRes(t, func() error { return BgpSessionGet(ctx, "fabric-1", "7") })
	if !strings.Contains(out, `"id":7`) {
		t.Errorf("get output missing session: %s", out)
	}
	if got := rec.last(); got.Path != "/api/v2/network-fabrics/1/bgp-sessions/7" {
		t.Errorf("get hit %s", got.Path)
	}

	// The config body is forwarded to the API with camelCase keys.
	config := []byte(`{"bgpNumbering":"numbered","bgpLinkConfiguration":"active","linkId":3}`)
	runSubRes(t, func() error { return BgpSessionCreate(ctx, "1", config) })
	created := rec.last()
	if created.Method != http.MethodPost || created.Path != "/api/v2/network-fabrics/1/bgp-sessions" {
		t.Errorf("create hit %s %s", created.Method, created.Path)
	}
	for _, want := range []string{`"bgpNumbering":"numbered"`, `"bgpLinkConfiguration":"active"`, `"linkId":3`} {
		if !strings.Contains(created.Body, want) {
			t.Errorf("create body missing %s: %s", want, created.Body)
		}
	}

	runSubRes(t, func() error {
		return BgpSessionUpdate(ctx, "1", "7", []byte(`{"customVariables":{"asn":65000}}`))
	})
	updated := rec.last()
	if updated.Path != "/api/v2/network-fabrics/1/bgp-sessions/7" {
		t.Errorf("update hit %s", updated.Path)
	}
	if !strings.Contains(updated.Body, `"customVariables"`) {
		t.Errorf("update body missing custom variables: %s", updated.Body)
	}

	if err := BgpSessionDelete(ctx, "1", "7"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if got := rec.last(); got.Method != http.MethodDelete {
		t.Errorf("delete used %s", got.Method)
	}

	if err := BgpSessionGet(ctx, "1", "abc"); err == nil {
		t.Error("expected an error for a non-numeric BGP session id")
	}
}

func TestBgpSessionConfigExample(t *testing.T) {
	out := testutils.CaptureStdout(t, func() {
		if err := BgpSessionConfigExample(nil); err != nil {
			t.Fatalf("config-example: %v", err)
		}
	})
	if !strings.Contains(out, "bgpNumbering") || !strings.Contains(out, "bgpLinkConfiguration") {
		t.Errorf("config example missing required fields: %s", out)
	}
}

// ---- link aggregations ----------------------------------------------------

func TestLinkAggregationCrud(t *testing.T) {
	rec := &subResRecorder{}
	srv := newSubResServer(t, rec)
	defer srv.Close()
	ctx := setupTestContext(srv.URL)

	out := runSubRes(t, func() error { return LinkAggregationList(ctx, "1") })
	if !strings.Contains(out, `"name":"lag-4"`) {
		t.Errorf("list output missing aggregation: %s", out)
	}

	out = runSubRes(t, func() error { return LinkAggregationGet(ctx, "1", "4") })
	if !strings.Contains(out, `"type":"lag"`) {
		t.Errorf("get output missing aggregation: %s", out)
	}

	runSubRes(t, func() error {
		return LinkAggregationCreate(ctx, "1", []byte(`{"type":"lag","linkIds":[1,2]}`))
	})
	created := rec.last()
	if created.Method != http.MethodPost || created.Path != "/api/v2/network-fabrics/1/link-aggregations" {
		t.Errorf("create hit %s %s", created.Method, created.Path)
	}
	for _, want := range []string{`"type":"lag"`, `"linkIds":[1,2]`} {
		if !strings.Contains(created.Body, want) {
			t.Errorf("create body missing %s: %s", want, created.Body)
		}
	}

	runSubRes(t, func() error {
		return LinkAggregationUpdate(ctx, "1", "4", []byte(`{"linkIds":[1,2,3]}`))
	})
	if got := rec.last(); !strings.Contains(got.Body, `"linkIds":[1,2,3]`) {
		t.Errorf("update body missing link ids: %s", got.Body)
	}

	if err := LinkAggregationDelete(ctx, "1", "4"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if got := rec.last(); got.Method != http.MethodDelete || got.Path != "/api/v2/network-fabrics/1/link-aggregations/4" {
		t.Errorf("delete hit %s %s", got.Method, got.Path)
	}
}

func TestLinkAggregationConfigExample(t *testing.T) {
	out := testutils.CaptureStdout(t, func() {
		if err := LinkAggregationConfigExample(nil); err != nil {
			t.Fatalf("config-example: %v", err)
		}
	})
	if !strings.Contains(out, "linkIds") {
		t.Errorf("config example missing linkIds: %s", out)
	}
}
