package container_type

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

var containerTypeItem = map[string]any{
	"id": 42, "name": "small", "label": "small-ct", "displayName": "Small",
	"cpuCores": 4, "ramGB": 16, "isExperimental": 0,
	"forUnmanagedContainersOnly": 0, "tags": []any{"prod"},
}

var containerItem = map[string]any{
	"id": 100, "name": "container-100", "siteId": 1, "infrastructureId": 1234,
	"userId": 7, "instanceId": 5678, "containerInstanceId": 5678,
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
		"/api/v2/container-types": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				testutils.JSONHandler(201, containerTypeItem)(w, r)
				return
			}
			testutils.JSONHandler(200, testutils.PaginatedResponse([]any{containerTypeItem}, 1, 1))(w, r)
		},
		"/api/v2/container-types/42": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				testutils.NoContentHandler()(w, r)
				return
			}
			testutils.JSONHandler(200, containerTypeItem)(w, r)
		},
		"/api/v2/container-types/42/containers": testutils.JSONHandler(200, testutils.PaginatedResponse([]any{containerItem}, 1, 1)),
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

func TestContainerTypeList(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error {
		return ContainerTypeList(ctx, ContainerTypeFilters{Label: []string{"small-ct"}})
	})

	if !strings.Contains(out, "small-ct") {
		t.Fatalf("expected the container type label in the output, got %s", out)
	}
	if got := rec.last(); got.Method != http.MethodGet || got.Path != "/api/v2/container-types" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
}

func TestContainerTypeGet(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error { return ContainerTypeGet(ctx, "42") })

	if !strings.Contains(out, "small-ct") {
		t.Fatalf("expected the container type label in the output, got %s", out)
	}
	if got := rec.last(); got.Method != http.MethodGet || got.Path != "/api/v2/container-types/42" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
}

func TestContainerTypeGetInvalidId(t *testing.T) {
	if err := ContainerTypeGet(nil, "not-a-number"); err == nil {
		t.Fatal("expected an error for a non numeric container type ID")
	}
}

func TestContainerTypeConfigExample(t *testing.T) {
	out := run(t, func() error { return ContainerTypeConfigExample(nil) })

	if !strings.Contains(out, "cpuCores") {
		t.Fatalf("expected the example to contain cpuCores, got %s", out)
	}
}

func TestContainerTypeCreate(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	config := []byte(`{"name":"small","cpuCores":4,"ramGB":16,"label":"small-ct"}`)

	run(t, func() error { return ContainerTypeCreate(ctx, config) })

	got := rec.last()
	if got.Method != http.MethodPost || got.Path != "/api/v2/container-types" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
	if !strings.Contains(got.Body, `"name":"small"`) {
		t.Fatalf("unexpected request body %s", got.Body)
	}
}

func TestContainerTypeCreateInvalidConfig(t *testing.T) {
	if err := ContainerTypeCreate(nil, []byte("not json")); err == nil {
		t.Fatal("expected an error for an invalid configuration")
	}
}

func TestContainerTypeUpdate(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	run(t, func() error { return ContainerTypeUpdate(ctx, "42", []byte(`{"label":"renamed"}`)) })

	got := rec.last()
	if got.Method != http.MethodPatch || got.Path != "/api/v2/container-types/42" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
	if !strings.Contains(got.Body, `"label":"renamed"`) {
		t.Fatalf("unexpected request body %s", got.Body)
	}
}

func TestContainerTypeDelete(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	if err := ContainerTypeDelete(ctx, "42"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := rec.last(); got.Method != http.MethodDelete || got.Path != "/api/v2/container-types/42" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
}

func TestContainerTypeContainers(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error { return ContainerTypeContainers(ctx, "42") })

	if !strings.Contains(out, "container-100") {
		t.Fatalf("expected the container name in the output, got %s", out)
	}
	if got := rec.last(); got.Path != "/api/v2/container-types/42/containers" {
		t.Fatalf("unexpected request path %s", got.Path)
	}
}

func TestContainerTypeErrorResponse(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/container-types/42": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	if err := ContainerTypeGet(ctx, "42"); err == nil {
		t.Fatal("expected an error for a 404 response")
	}
}
