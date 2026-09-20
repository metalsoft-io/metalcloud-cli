package route_domain

import (
	"encoding/json"
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

// routeDomainItem is the route domain entity. Its revision (2) deliberately
// differs from its config's (9) so the tests can prove that the config
// endpoints use the config's own revision for If-Match.
var routeDomainItem = map[string]any{
	"id": 2, "label": "tenant1", "name": "tenant1", "kind": "evpn_l3vpn", "revision": 2,
	"annotations": map[string]any{}, "serviceStatus": "ordered",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
}

var routeDomainConfigItem = map[string]any{
	"id": 2, "deployType": "create", "deployStatus": "not_started",
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
	"revision": 9, "kind": "evpn_l3vpn",
	"vrfAllocationStrategies":    []any{},
	"l3VlanAllocationStrategies": []any{},
	"l3VniAllocationStrategies":  []any{},
	"autoRouteDistinguisher":     false,
	"autoRouteTarget":            false,
}

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

func newRouteDomainServer(t *testing.T, rec *recorder) *httptest.Server {
	t.Helper()
	routes := map[string]http.HandlerFunc{
		"/api/v2/route-domains/2":        testutils.JSONHandler(http.StatusOK, routeDomainItem),
		"/api/v2/route-domains/2/config": testutils.JSONHandler(http.StatusOK, routeDomainConfigItem),
	}
	wrapped := make(map[string]http.HandlerFunc, len(routes))
	for path, handler := range routes {
		wrapped[path] = rec.wrap(handler)
	}
	return testutils.NewTestServer(wrapped)
}

func TestRouteDomainConfigGet(t *testing.T) {
	rec := &recorder{}
	ts := newRouteDomainServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	out := testutils.CaptureStdout(t, func() {
		if err := RouteDomainConfigGet(ctx, "2"); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if !strings.Contains(out, `"revision":9`) {
		t.Errorf("get-config output missing the config revision: %s", out)
	}
	if rec.last().Path != "/api/v2/route-domains/2/config" {
		t.Errorf("unexpected path %s", rec.last().Path)
	}
}

func TestRouteDomainConfigUpdate_UsesConfigRevision(t *testing.T) {
	rec := &recorder{}
	ts := newRouteDomainServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	testutils.CaptureStdout(t, func() {
		if err := RouteDomainConfigUpdate(ctx, "2", []byte(`{"autoRouteTarget":true}`)); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	req := rec.last()
	if req.Method != http.MethodPatch || req.Path != "/api/v2/route-domains/2/config" {
		t.Fatalf("unexpected request %s %s", req.Method, req.Path)
	}
	// The config's revision (9), not the route domain's (2).
	if got := req.Header.Get("If-Match"); got != "9" {
		t.Errorf("expected If-Match 9 (the config revision), got %q", got)
	}
	var body map[string]any
	_ = json.Unmarshal([]byte(req.Body), &body)
	if body["autoRouteTarget"] != true {
		t.Errorf("update body should carry autoRouteTarget, got %s", req.Body)
	}
}

func TestRouteDomainConfigGet_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := RouteDomainConfigGet(ctx, "not-a-number"); err == nil {
		t.Fatal("expected an error for an invalid route domain ID, got nil")
	}
}

func TestRouteDomainConfigGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/route-domains/99/config": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := RouteDomainConfigGet(ctx, "99"); err == nil {
		t.Fatal("expected an error for HTTP 404, got nil")
	}
}
