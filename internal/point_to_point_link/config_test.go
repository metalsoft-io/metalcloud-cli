package point_to_point_link

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

// linkItem is the link object; its revision (3) guards the link itself.
var ptpLinkItem = map[string]any{
	"id": 42, "label": "leaf01-to-spine01", "description": "uplink",
	"routingActivation": "default", "serviceStatus": "active", "revision": 3,
}

// configItem is the link's config object; its own revision (11) guards the
// config sub-resources - deliberately different from the link's.
var ptpConfigItem = map[string]any{
	"id": 42, "deployType": "staged", "deployStatus": "not_started",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
	"revision": 11, "mtu": 9216,
	"ipv4": map[string]any{
		"subnetAllocationStrategies": []any{
			map[string]any{"id": 1, "kind": "manual", "scope": map[string]any{"kind": "global"}, "subnetId": 7},
		},
		"staticRoutes": []any{map[string]any{"id": 5, "destinationPrefix": "10.0.0.0/24"}},
	},
}

var ptpRouteItem = map[string]any{"id": 5, "destinationPrefix": "10.0.0.0/24"}

type ptpRecorder struct {
	mu   sync.Mutex
	reqs []ptpRequest
}

type ptpRequest struct {
	Method  string
	Path    string
	IfMatch string
	Body    string
}

func (r *ptpRecorder) wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		r.mu.Lock()
		r.reqs = append(r.reqs, ptpRequest{
			Method: req.Method, Path: req.URL.Path,
			IfMatch: req.Header.Get("If-Match"), Body: string(body),
		})
		r.mu.Unlock()
		next(w, req)
	}
}

func (r *ptpRecorder) last() ptpRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.reqs[len(r.reqs)-1]
}

func newPtpServer(t *testing.T, rec *ptpRecorder) *httptest.Server {
	t.Helper()
	routes := map[string]http.HandlerFunc{
		"/api/v2/point-to-point-links/42": ptpMethodRouter(map[string]http.HandlerFunc{
			http.MethodGet:   testutils.JSONHandler(200, ptpLinkItem),
			http.MethodPatch: testutils.JSONHandler(200, ptpLinkItem),
		}),
		"/api/v2/point-to-point-links/42/config": ptpMethodRouter(map[string]http.HandlerFunc{
			http.MethodGet:   testutils.JSONHandler(200, ptpConfigItem),
			http.MethodPatch: testutils.JSONHandler(200, ptpConfigItem),
		}),
		"/api/v2/point-to-point-links/42/config/ipv4/static-routes": ptpMethodRouter(map[string]http.HandlerFunc{
			http.MethodGet:  testutils.JSONHandler(200, []any{ptpRouteItem}),
			http.MethodPost: testutils.JSONHandler(201, ptpRouteItem),
		}),
		"/api/v2/point-to-point-links/42/config/ipv4/static-routes/5": ptpMethodRouter(map[string]http.HandlerFunc{
			http.MethodGet:    testutils.JSONHandler(200, ptpRouteItem),
			http.MethodDelete: testutils.NoContentHandler(),
		}),
		"/api/v2/point-to-point-links/42/config/ipv6/static-routes": testutils.JSONHandler(200,
			map[string]any{"data": []any{ptpRouteItem}}),
	}
	wrapped := make(map[string]http.HandlerFunc, len(routes))
	for path, handler := range routes {
		wrapped[path] = rec.wrap(handler)
	}
	return testutils.NewTestServer(wrapped)
}

func ptpMethodRouter(handlers map[string]http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler, ok := handlers[r.Method]
		if !ok {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}

func runPtp(t *testing.T, fn func() error) string {
	t.Helper()
	var out string
	var err error
	out = testutils.CaptureStdout(t, func() { err = fn() })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return out
}

func TestPointToPointLinkUpdateUsesLinkRevision(t *testing.T) {
	rec := &ptpRecorder{}
	srv := newPtpServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := runPtp(t, func() error {
		return PointToPointLinkUpdate(ctx, "42", []byte(`{"description":"leaf01 uplink"}`))
	})
	if !strings.Contains(out, `"Label":"leaf01-to-spine01"`) && !strings.Contains(out, "leaf01-to-spine01") {
		t.Errorf("update output missing link: %s", out)
	}

	got := rec.last()
	if got.Method != http.MethodPatch || got.Path != "/api/v2/point-to-point-links/42" {
		t.Errorf("update hit %s %s", got.Method, got.Path)
	}
	if got.IfMatch != "3" {
		t.Errorf("update must carry the link revision as If-Match, got %q", got.IfMatch)
	}
	if !strings.Contains(got.Body, `"description":"leaf01 uplink"`) {
		t.Errorf("update body wrong: %s", got.Body)
	}
}

func TestPointToPointLinkConfigGetAndUpdate(t *testing.T) {
	rec := &ptpRecorder{}
	srv := newPtpServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := runPtp(t, func() error { return PointToPointLinkConfigGet(ctx, "42") })
	if !strings.Contains(out, `"mtu":9216`) {
		t.Errorf("get-config output missing mtu: %s", out)
	}
	if got := rec.last(); got.Path != "/api/v2/point-to-point-links/42/config" || got.Method != http.MethodGet {
		t.Errorf("get-config hit %s %s", got.Method, got.Path)
	}

	runPtp(t, func() error { return PointToPointLinkConfigUpdate(ctx, "42", []byte(`{"mtu":1500}`)) })
	got := rec.last()
	if got.Method != http.MethodPatch || got.Path != "/api/v2/point-to-point-links/42/config" {
		t.Errorf("update-config hit %s %s", got.Method, got.Path)
	}
	// The config object's revision (11) guards the write, not the link's (3).
	if got.IfMatch != "11" {
		t.Errorf("update-config must use the config revision, got %q", got.IfMatch)
	}
	if !strings.Contains(got.Body, `"mtu":1500`) {
		t.Errorf("update-config body wrong: %s", got.Body)
	}

	if err := PointToPointLinkConfigGet(ctx, "abc"); err == nil {
		t.Error("expected an error for a non-numeric link id")
	}
}

func TestStaticRouteListGetAddRemove(t *testing.T) {
	rec := &ptpRecorder{}
	srv := newPtpServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	// A bare array response.
	out := runPtp(t, func() error { return StaticRouteList(ctx, "42", "ipv4") })
	if !strings.Contains(out, "10.0.0.0/24") {
		t.Errorf("list output missing route: %s", out)
	}

	// A {data: [...]} envelope response for the other family.
	out = runPtp(t, func() error { return StaticRouteList(ctx, "42", "ipv6") })
	if !strings.Contains(out, "10.0.0.0/24") {
		t.Errorf("paginated list output missing route: %s", out)
	}

	out = runPtp(t, func() error { return StaticRouteGet(ctx, "42", "ipv4", "5") })
	if !strings.Contains(out, "10.0.0.0/24") {
		t.Errorf("get output missing route: %s", out)
	}

	runPtp(t, func() error { return StaticRouteAdd(ctx, "42", "ipv4", "10.0.0.0/24") })
	added := rec.last()
	if added.Method != http.MethodPost || added.Path != "/api/v2/point-to-point-links/42/config/ipv4/static-routes" {
		t.Errorf("add hit %s %s", added.Method, added.Path)
	}
	if added.IfMatch != "11" {
		t.Errorf("add must use the config revision, got %q", added.IfMatch)
	}
	if !strings.Contains(added.Body, `"destinationPrefix":"10.0.0.0/24"`) {
		t.Errorf("add body wrong: %s", added.Body)
	}

	if err := StaticRouteRemove(ctx, "42", "ipv4", "5"); err != nil {
		t.Fatalf("remove: %v", err)
	}
	removed := rec.last()
	if removed.Method != http.MethodDelete || removed.Path != "/api/v2/point-to-point-links/42/config/ipv4/static-routes/5" {
		t.Errorf("remove hit %s %s", removed.Method, removed.Path)
	}
	if removed.IfMatch != "11" {
		t.Errorf("remove must use the config revision, got %q", removed.IfMatch)
	}
}

func TestStaticRouteInvalidFamilyAndId(t *testing.T) {
	rec := &ptpRecorder{}
	srv := newPtpServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if err := StaticRouteList(ctx, "42", "ipv5"); err == nil || !strings.Contains(err.Error(), "invalid address family") {
		t.Errorf("expected an invalid family error, got %v", err)
	}
	if err := StaticRouteGet(ctx, "42", "ipv4", "xyz"); err == nil {
		t.Error("expected an error for a non-numeric route id")
	}
	if err := StaticRouteRemove(ctx, "42", "ipv6", "abc"); err == nil {
		t.Error("expected an error for a non-numeric route id")
	}
}

func TestPointToPointLinkConfigExamples(t *testing.T) {
	out := testutils.CaptureStdout(t, func() {
		if err := PointToPointLinkConfigExampleUpdate(nil); err != nil {
			t.Fatalf("update example: %v", err)
		}
		if err := PointToPointLinkConfigExampleConfig(nil); err != nil {
			t.Fatalf("config example: %v", err)
		}
	})
	if !strings.Contains(out, "label") || !strings.Contains(out, "mtu") {
		t.Errorf("examples missing expected fields: %s", out)
	}
}
