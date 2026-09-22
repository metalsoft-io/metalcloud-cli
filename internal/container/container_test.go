package container

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

var containerItem = map[string]any{
	"id": 100, "name": "container-100", "siteId": 1, "infrastructureId": 1234,
	"userId": 7, "instanceId": 5678, "containerInstanceId": 5678,
	// The live API returns `hosts` as a plain string while sdk.Container declares
	// []string; the raw-body path must cope with it.
	"host": "host-1", "hosts": "host-1",
	"cpuCores": 4, "ramGB": 16, "diskSizeGB": 40,
	"typeId": 42, "poolId": 3, "administrationState": "active",
	"powerState": "on", "powerStateLastUpdatedTimestamp": "2024-01-01T00:00:00Z",
	"createdTimestamp": "2024-01-01T00:00:00Z", "allocationTimestamp": "2024-01-01T00:00:00Z",
	"disks": []any{},
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

func newServer(t *testing.T, rec *recorder) *httptest.Server {
	t.Helper()
	routes := map[string]http.HandlerFunc{
		"/api/v2/containers":                         testutils.JSONHandler(200, testutils.PaginatedResponse([]any{containerItem}, 1, 1)),
		"/api/v2/containers/100":                     testutils.JSONHandler(200, containerItem),
		"/api/v2/containers/100/power-status":        testutils.JSONHandler(200, "on"),
		"/api/v2/containers/100/remote-console-info": testutils.JSONHandler(200, map[string]any{"active_connections": 2}),
		"/api/v2/containers/100/start":               testutils.NoContentHandler(),
		"/api/v2/containers/100/shutdown":            testutils.NoContentHandler(),
		"/api/v2/containers/100/reboot":              testutils.NoContentHandler(),
	}
	wrapped := make(map[string]http.HandlerFunc, len(routes))
	for path, handler := range routes {
		wrapped[path] = rec.wrap(handler)
	}
	return testutils.NewTestServer(wrapped)
}

func run(t *testing.T, fn func() error) string {
	t.Helper()
	var err error
	out := testutils.CaptureStdout(t, func() { err = fn() })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return out
}

func TestContainerList(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error {
		return ContainerList(ctx, ContainerFilters{InfrastructureId: []string{"1234"}, TypeId: []string{"42"}})
	})

	if !strings.Contains(out, "container-100") {
		t.Fatalf("expected the container name in the output, got %s", out)
	}
	if got := rec.last(); got.Method != http.MethodGet || got.Path != "/api/v2/containers" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
}

func TestContainerGet(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error { return ContainerGet(ctx, "100") })

	if !strings.Contains(out, "container-100") {
		t.Fatalf("expected the container name in the output, got %s", out)
	}
}

func TestContainerGetInvalidId(t *testing.T) {
	if err := ContainerGet(nil, "abc"); err == nil {
		t.Fatal("expected an error for a non numeric container ID")
	}
}

func TestContainerUpdate(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	run(t, func() error { return ContainerUpdate(ctx, "100", []byte(`{"tags":["prod"]}`)) })

	got := rec.last()
	if got.Method != http.MethodPatch || got.Path != "/api/v2/containers/100" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
	if !strings.Contains(got.Body, `"tags":["prod"]`) {
		t.Fatalf("unexpected request body %s", got.Body)
	}
}

func TestContainerPowerStatus(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error { return ContainerPowerStatus(ctx, "100") })

	if !strings.Contains(out, "on") {
		t.Fatalf("expected the power state in the output, got %s", out)
	}
	if got := rec.last(); got.Path != "/api/v2/containers/100/power-status" {
		t.Fatalf("unexpected request path %s", got.Path)
	}
}

func TestContainerRemoteConsoleInfo(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error { return ContainerRemoteConsoleInfo(ctx, "100") })

	if !strings.Contains(out, "active_connections") {
		t.Fatalf("expected the console info in the output, got %s", out)
	}
}

func TestContainerPowerControl(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	for _, action := range []string{"start", "shutdown", "reboot"} {
		if err := ContainerPowerControl(ctx, "100", action); err != nil {
			t.Fatalf("unexpected error for %s: %v", action, err)
		}

		got := rec.last()
		if got.Method != http.MethodPost || got.Path != "/api/v2/containers/100/"+action {
			t.Fatalf("unexpected request %s %s", got.Method, got.Path)
		}
	}
}

func TestContainerPowerControlUnsupportedAction(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	if err := ContainerPowerControl(ctx, "100", "explode"); err == nil {
		t.Fatal("expected an error for an unsupported power action")
	}
}

func TestContainerErrorResponse(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/containers/100": testutils.ErrorHandler(403, "forbidden"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	if err := ContainerGet(ctx, "100"); err == nil {
		t.Fatal("expected an error for a 403 response")
	}
}
