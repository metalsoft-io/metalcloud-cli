package container_instance

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

const ts0 = "2024-01-01T00:00:00Z"

var infrastructureItem = map[string]any{
	"id": 1234, "label": "test-infra", "revision": 1, "serviceStatus": "active",
	"datacenterName": "dc1", "siteId": 1, "designIsLocked": 0,
	"createdTimestamp": ts0, "updatedTimestamp": ts0,
	"config": map[string]any{},
}

var containerInstanceConfigItem = map[string]any{
	"revision": 3, "label": "ci-1", "typeId": 42,
	"deployType": "create", "deployStatus": "not_started",
	"diskSizeGB": 40, "ramGB": 16, "cpuCores": 4,
	"updatedTimestamp": ts0,
}

var containerInstanceItem = map[string]any{
	"label": "ci-1", "typeId": 42, "containerId": 100,
	"diskSizeGB": 40, "ramGB": 16, "cpuCores": 4,
	"updatedTimestamp": ts0, "id": 5678, "revision": 9,
	"groupId": 77, "infrastructureId": 1234,
	"infrastructure": map[string]any{"id": 1234},
	"serviceStatus":  "active", "instanceType": "container",
	"createdTimestamp": ts0,
	"config":           containerInstanceConfigItem,
	"meta":             map[string]any{},
}

// containerInstanceVariablesItem mirrors what the API actually answers for a
// container instance: `server` is null, which the strict SDK model rejects
// ("no value given for required property serverId"), so these two endpoints are
// read raw.
var containerInstanceVariablesItem = map[string]any{
	"site":       map[string]any{"id": 1, "revision": 1, "slug": "dc1", "name": "DC1"},
	"siteConfig": map[string]any{},
	"server":     nil,
	"serverInstance": map[string]any{
		"id": 5678, "revision": 1, "label": "ci-1", "infrastructureId": 1234,
	},
	"serverInstanceGroup": map[string]any{"id": 77, "label": "cig-1"},
	"infrastructure":      infrastructureItem,
	"driveGroups":         []any{}, "drives": []any{}, "fileShares": []any{},
	"buckets": []any{}, "sharedDrives": []any{}, "network": []any{},
	"variables": map[string]any{"foo": "bar"},
}

var containerInstanceOsInstallationDataItem = map[string]any{
	"site":       map[string]any{"id": 1, "slug": "dc1", "name": "DC1"},
	"siteConfig": map[string]any{"repoURL": "https://repo", "DNSServers": []any{"1.1.1.1"}, "NTPServers": []any{"pool.ntp.org"}},
	"server":     nil,
	"serverInstance": map[string]any{
		"id": 5678, "label": "ci-1",
	},
	"serverInstanceGroup": map[string]any{"id": 77, "label": "cig-1"},
	"infrastructure":      map[string]any{"label": "test-infra"},
	"driveGroups":         []any{}, "drives": []any{}, "fileShares": []any{},
	"buckets": []any{}, "sharedDrives": []any{}, "network": []any{},
}

// ---------------------------------------------------------------------------
// Test server
// ---------------------------------------------------------------------------

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
	base := "/api/v2/infrastructures/1234"
	routes := map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.JSONHandler(200, testutils.PaginatedResponse([]any{infrastructureItem}, 1, 1)),

		base + "/container-instances": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				testutils.JSONHandler(201, containerInstanceItem)(w, r)
				return
			}
			testutils.JSONHandler(200, testutils.PaginatedResponse([]any{containerInstanceItem}, 1, 1))(w, r)
		},
		base + "/container-instances/5678": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				testutils.NoContentHandler()(w, r)
				return
			}
			testutils.JSONHandler(200, containerInstanceItem)(w, r)
		},
		base + "/container-instances/5678/config":                testutils.JSONHandler(200, containerInstanceConfigItem),
		base + "/container-instances/5678/meta":                  testutils.JSONHandler(200, containerInstanceItem),
		base + "/container-instances/5678/credentials":           testutils.JSONHandler(200, map[string]any{"username": "root", "initialPassword": "pass"}),
		base + "/container-instances/5678/variables":             testutils.JSONHandler(200, containerInstanceVariablesItem),
		base + "/container-instances/5678/os-installation-data":  testutils.JSONHandler(200, containerInstanceOsInstallationDataItem),
		base + "/container-instances/5678/power-status":          testutils.JSONHandler(200, "on"),
		base + "/container-instances/5678/start":                 testutils.NoContentHandler(),
		base + "/container-instances/5678/shutdown":              testutils.NoContentHandler(),
		base + "/container-instances/5678/reboot":                testutils.NoContentHandler(),
		base + "/container-instances/5678/actions/apply-type/42": testutils.JSONHandler(200, containerInstanceItem),

		base + "/container-instance-groups": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				testutils.JSONHandler(201, containerInstanceGroupItem)(w, r)
				return
			}
			testutils.JSONHandler(200, testutils.PaginatedResponse([]any{containerInstanceGroupItem}, 1, 1))(w, r)
		},
		base + "/container-instance-groups/77": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				testutils.NoContentHandler()(w, r)
				return
			}
			testutils.JSONHandler(200, containerInstanceGroupItem)(w, r)
		},
		base + "/container-instance-groups/77/config":                testutils.JSONHandler(200, containerInstanceGroupConfigItem),
		base + "/container-instance-groups/77/meta":                  testutils.JSONHandler(200, containerInstanceGroupItem),
		base + "/container-instance-groups/77/container-instances":   testutils.JSONHandler(200, testutils.PaginatedResponse([]any{containerInstanceItem}, 1, 1)),
		base + "/container-instance-groups/77/interfaces":            testutils.JSONHandler(200, testutils.PaginatedResponse([]any{containerInstanceGroupInterfaceItem}, 1, 1)),
		base + "/container-instance-groups/77/interfaces/3":          testutils.JSONHandler(200, containerInstanceGroupInterfaceItem),
		base + "/container-instance-groups/77/actions/apply-type/42": testutils.JSONHandler(200, containerInstanceGroupItem),
		base + "/container-instance-groups/77/config/networking":     testutils.JSONHandler(200, networkEndpointGroupItem),
		base + "/container-instance-groups/77/config/networking/connections": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				testutils.JSONHandler(201, networkConnectionItem)(w, r)
				return
			}
			testutils.JSONHandler(200, map[string]any{"data": []any{networkConnectionItem}})(w, r)
		},
		base + "/container-instance-groups/77/config/networking/connections/5": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				testutils.NoContentHandler()(w, r)
				return
			}
			testutils.JSONHandler(200, networkConnectionItem)(w, r)
		},
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

// ---------------------------------------------------------------------------
// Container instance tests
// ---------------------------------------------------------------------------

func TestContainerInstanceList(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error { return ContainerInstanceList(ctx, "test-infra") })

	if !strings.Contains(out, "ci-1") {
		t.Fatalf("expected the container instance label in the output, got %s", out)
	}
	if got := rec.last(); got.Path != "/api/v2/infrastructures/1234/container-instances" {
		t.Fatalf("unexpected request path %s", got.Path)
	}
}

func TestContainerInstanceGet(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error { return ContainerInstanceGet(ctx, "1234", "5678") })

	if !strings.Contains(out, "ci-1") {
		t.Fatalf("expected the container instance label in the output, got %s", out)
	}
}

func TestContainerInstanceGetInvalidId(t *testing.T) {
	if err := ContainerInstanceGet(nil, "1234", "nope"); err == nil {
		t.Fatal("expected an error for a non numeric container instance ID")
	}
}

func TestContainerInstanceConfigExample(t *testing.T) {
	out := run(t, func() error { return ContainerInstanceConfigExample(nil) })

	if !strings.Contains(out, "typeId") {
		t.Fatalf("expected the example to contain typeId, got %s", out)
	}
}

func TestContainerInstanceCreateFromConfig(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	run(t, func() error {
		return ContainerInstanceCreate(ctx, "test-infra", []byte(`{"typeId":42,"groupId":77,"diskSizeGB":40}`))
	})

	got := rec.last()
	if got.Method != http.MethodPost || got.Path != "/api/v2/infrastructures/1234/container-instances" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
	if !strings.Contains(got.Body, `"typeId":42`) {
		t.Fatalf("unexpected request body %s", got.Body)
	}
}

func TestContainerInstanceCreateFromFlags(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	run(t, func() error {
		return ContainerInstanceCreateFromFlags(ctx, "1234", "42", "77", "40", []string{"prod"})
	})

	got := rec.last()
	if got.Method != http.MethodPost {
		t.Fatalf("unexpected request method %s", got.Method)
	}
	if !strings.Contains(got.Body, `"groupId":77`) || !strings.Contains(got.Body, `"diskSizeGB":40`) {
		t.Fatalf("unexpected request body %s", got.Body)
	}
}

func TestContainerInstanceDeleteSendsIfMatch(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	if err := ContainerInstanceDelete(ctx, "1234", "5678"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := rec.last()
	if got.Method != http.MethodDelete {
		t.Fatalf("unexpected request method %s", got.Method)
	}
	if got.Header.Get("If-Match") != "9" {
		t.Fatalf("expected the instance revision as If-Match, got %q", got.Header.Get("If-Match"))
	}
}

func TestContainerInstanceGetConfig(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error { return ContainerInstanceGetConfig(ctx, "1234", "5678") })

	if !strings.Contains(out, "ci-1") {
		t.Fatalf("expected the configuration label in the output, got %s", out)
	}
	if got := rec.last(); got.Path != "/api/v2/infrastructures/1234/container-instances/5678/config" {
		t.Fatalf("unexpected request path %s", got.Path)
	}
}

func TestContainerInstanceUpdateConfigUsesConfigRevision(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	run(t, func() error {
		return ContainerInstanceUpdateConfig(ctx, "1234", "5678", []byte(`{"label":"renamed"}`))
	})

	got := rec.last()
	if got.Method != http.MethodPatch || got.Path != "/api/v2/infrastructures/1234/container-instances/5678/config" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
	// The config sub-resource is guarded by the config revision (3), not by the
	// container instance revision (9).
	if got.Header.Get("If-Match") != "3" {
		t.Fatalf("expected the config revision as If-Match, got %q", got.Header.Get("If-Match"))
	}
	if !strings.Contains(got.Body, `"label":"renamed"`) {
		t.Fatalf("unexpected request body %s", got.Body)
	}
}

func TestContainerInstanceUpdateMeta(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	run(t, func() error {
		return ContainerInstanceUpdateMeta(ctx, "1234", "5678", []byte(`{"tags":["prod"]}`))
	})

	got := rec.last()
	if got.Method != http.MethodPatch || got.Path != "/api/v2/infrastructures/1234/container-instances/5678/meta" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
	if !strings.Contains(got.Body, `"tags":["prod"]`) {
		t.Fatalf("unexpected request body %s", got.Body)
	}
}

func TestContainerInstanceCredentials(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error { return ContainerInstanceCredentials(ctx, "1234", "5678") })

	if !strings.Contains(out, "root") {
		t.Fatalf("expected the username in the output, got %s", out)
	}
}

func TestContainerInstanceVariables(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error { return ContainerInstanceVariables(ctx, "1234", "5678", "OSAsset") })

	if !strings.Contains(out, "DC1") {
		t.Fatalf("expected the site name in the output, got %s", out)
	}
	if got := rec.last(); got.Path != "/api/v2/infrastructures/1234/container-instances/5678/variables" {
		t.Fatalf("unexpected request path %s", got.Path)
	}
}

func TestContainerInstanceVariablesInvalidUsage(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	if err := ContainerInstanceVariables(ctx, "1234", "5678", "NotAUsage"); err == nil {
		t.Fatal("expected an error for an invalid usage type")
	}
}

func TestContainerInstanceOSInstallationData(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error {
		return ContainerInstanceOSInstallationData(ctx, "1234", "5678", "OSAsset", true)
	})

	if !strings.Contains(out, "https://repo") {
		t.Fatalf("expected the repo URL in the output, got %s", out)
	}
	if got := rec.last(); got.Path != "/api/v2/infrastructures/1234/container-instances/5678/os-installation-data" {
		t.Fatalf("unexpected request path %s", got.Path)
	}
}

func TestContainerInstancePowerStatus(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	out := run(t, func() error { return ContainerInstancePowerStatus(ctx, "1234", "5678") })

	if !strings.Contains(out, "on") {
		t.Fatalf("expected the power state in the output, got %s", out)
	}
}

func TestContainerInstancePowerControl(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	for _, action := range []string{"start", "shutdown", "reboot"} {
		if err := ContainerInstancePowerControl(ctx, "1234", "5678", action); err != nil {
			t.Fatalf("unexpected error for %s: %v", action, err)
		}

		got := rec.last()
		if got.Method != http.MethodPost || !strings.HasSuffix(got.Path, "/"+action) {
			t.Fatalf("unexpected request %s %s", got.Method, got.Path)
		}
	}
}

func TestContainerInstancePowerControlUnsupportedAction(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	if err := ContainerInstancePowerControl(ctx, "1234", "5678", "explode"); err == nil {
		t.Fatal("expected an error for an unsupported power action")
	}
}

func TestContainerInstanceApplyTypeSendsIfMatch(t *testing.T) {
	rec := &recorder{}
	ts := newServer(t, rec)
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	run(t, func() error { return ContainerInstanceApplyType(ctx, "1234", "5678", "42") })

	got := rec.last()
	if got.Method != http.MethodPost || got.Path != "/api/v2/infrastructures/1234/container-instances/5678/actions/apply-type/42" {
		t.Fatalf("unexpected request %s %s", got.Method, got.Path)
	}
	if got.Header.Get("If-Match") != "9" {
		t.Fatalf("expected the instance revision as If-Match, got %q", got.Header.Get("If-Match"))
	}
}

func TestContainerInstanceUnknownInfrastructure(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.JSONHandler(200, testutils.PaginatedResponse([]any{}, 1, 1)),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)

	if err := ContainerInstanceList(ctx, "missing"); err == nil {
		t.Fatal("expected an error for an unknown infrastructure")
	}
}
