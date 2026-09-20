package endpoint_instance

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

// ---------------------------------------------------------------------------
// Fixtures - every map satisfies the required properties of its SDK model.
// ---------------------------------------------------------------------------

var infrastructureItem = map[string]any{
	"id": 123, "revision": 1, "label": "prod-env", "serviceStatus": "active",
	"datacenterName": "dc1", "siteId": 1, "designIsLocked": 0, "config": map[string]any{},
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
}

var endpointInstanceItem = map[string]any{
	"id": 42, "revision": 7, "label": "ei-1", "infrastructureId": 123, "groupId": 12,
	"endpointId": 5, "hostname": "ei-1.example.com", "serviceStatus": "active",
	"meta":             map[string]any{"tags": []any{"prod"}},
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

var endpointInstanceConfigItem = map[string]any{
	"revision": 3, "label": "ei-1", "updatedTimestamp": "2024-01-02T00:00:00Z",
	"groupId": 12, "endpointId": 5, "deployType": "create", "deployStatus": "not_started",
}

var endpointInstanceGroupItem = map[string]any{
	"id": 12, "revision": 9, "label": "eig-1", "infrastructureId": 123,
	"endpointGroupName": "web", "serviceStatus": "active",
	"meta":             map[string]any{"tags": []any{"prod"}},
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

var endpointInstanceGroupConfigItem = map[string]any{
	"revision": 4, "label": "eig-1", "updatedTimestamp": "2024-01-02T00:00:00Z",
	"endpointGroupName": "web", "deployType": "create", "deployStatus": "not_started",
}

var networkEndpointGroupItem = map[string]any{
	"id": "77", "name": "eig-1-neg", "siteId": 1, "revision": "2",
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

var networkConnectionItem = map[string]any{
	"id": "5", "tagged": true, "accessMode": "l2", "mtu": 1500,
	"providesDefaultRoute": false,
}

var aclItem = map[string]any{
	"id": "3", "ruleType": "ipv4", "direction": "in", "sequence": 10,
	"forwardingAction": "allow", "enforcementPoint": "svi",
	"sourceAddress": "10.0.0.0/24", "destinationAddress": "10.0.1.0/24",
	"endpointGroupId": "77", "logicalNetworkId": "7",
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

// ---------------------------------------------------------------------------
// Request recorder
// ---------------------------------------------------------------------------

type recorder struct {
	mu   sync.Mutex
	reqs []recordedRequest
}

type recordedRequest struct {
	Method string
	Path   string
	Query  string
	Header http.Header
	Body   string
}

func (r *recorder) wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		r.mu.Lock()
		r.reqs = append(r.reqs, recordedRequest{
			Method: req.Method, Path: req.URL.Path, Query: req.URL.RawQuery,
			Header: req.Header.Clone(), Body: string(body),
		})
		r.mu.Unlock()
		next(w, req)
	}
}

func (r *recorder) last() recordedRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.reqs[len(r.reqs)-1]
}

func (r *recorder) paths() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	paths := make([]string, 0, len(r.reqs))
	for _, req := range r.reqs {
		paths = append(paths, req.Path)
	}
	return paths
}

// methodSwitch serves listHandler for GET and writeHandler for every other verb.
func methodSwitch(listHandler, writeHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			listHandler(w, r)
			return
		}
		writeHandler(w, r)
	}
}

func newServer(t *testing.T, rec *recorder) *httptest.Server {
	t.Helper()
	routes := map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.JSONHandler(200, testutils.PaginatedResponse([]any{infrastructureItem}, 1, 1)),

		"/api/v2/endpoint-instances": testutils.JSONHandler(200, testutils.PaginatedResponse([]any{endpointInstanceItem}, 1, 1)),
		"/api/v2/infrastructures/123/endpoint-instances": methodSwitch(
			testutils.JSONHandler(200, testutils.PaginatedResponse([]any{endpointInstanceItem}, 1, 1)),
			testutils.JSONHandler(201, endpointInstanceItem),
		),
		"/api/v2/endpoint-instances/42": methodSwitch(
			testutils.JSONHandler(200, endpointInstanceItem),
			testutils.NoContentHandler(),
		),
		"/api/v2/endpoint-instances/42/config": testutils.JSONHandler(200, endpointInstanceConfigItem),
		"/api/v2/endpoint-instances/42/meta":   testutils.NoContentHandler(),

		"/api/v2/infrastructures/123/endpoint-instance-groups": methodSwitch(
			testutils.JSONHandler(200, testutils.PaginatedResponse([]any{endpointInstanceGroupItem}, 1, 1)),
			testutils.JSONHandler(201, endpointInstanceGroupItem),
		),
		"/api/v2/endpoint-instance-groups/12": methodSwitch(
			testutils.JSONHandler(200, endpointInstanceGroupItem),
			testutils.NoContentHandler(),
		),
		"/api/v2/endpoint-instance-groups/12/config":             testutils.JSONHandler(200, endpointInstanceGroupConfigItem),
		"/api/v2/endpoint-instance-groups/12/meta":               testutils.NoContentHandler(),
		"/api/v2/endpoint-instance-groups/12/endpoint-instances": testutils.JSONHandler(200, testutils.PaginatedResponse([]any{endpointInstanceItem}, 1, 1)),
		"/api/v2/endpoint-instance-groups/12/config/networking":  testutils.JSONHandler(200, networkEndpointGroupItem),
		"/api/v2/endpoint-instance-groups/12/config/networking/connections": methodSwitch(
			testutils.JSONHandler(200, map[string]any{"data": []any{networkConnectionItem}}),
			testutils.JSONHandler(201, networkConnectionItem),
		),
		"/api/v2/endpoint-instance-groups/12/config/networking/connections/5": methodSwitch(
			testutils.JSONHandler(200, networkConnectionItem),
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete {
					testutils.NoContentHandler()(w, r)
					return
				}
				testutils.JSONHandler(200, networkConnectionItem)(w, r)
			},
		),
		"/api/v2/endpoint-instance-groups/12/config/networking/connections/5/security/rules": methodSwitch(
			testutils.JSONHandler(200, []any{aclItem}),
			testutils.JSONHandler(201, aclItem),
		),
		"/api/v2/endpoint-instance-groups/12/config/networking/connections/5/security/rules/3": methodSwitch(
			testutils.JSONHandler(200, aclItem),
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodDelete {
					testutils.NoContentHandler()(w, r)
					return
				}
				testutils.JSONHandler(200, aclItem)(w, r)
			},
		),
	}

	wrapped := make(map[string]http.HandlerFunc, len(routes))
	for path, handler := range routes {
		wrapped[path] = rec.wrap(handler)
	}
	return testutils.NewTestServer(wrapped)
}

func run(t *testing.T, fn func() error) string {
	t.Helper()
	var out string
	var err error
	out = testutils.CaptureStdout(t, func() { err = fn() })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return out
}

// ---------------------------------------------------------------------------
// Endpoint instance tests
// ---------------------------------------------------------------------------

func TestEndpointInstanceListGlobalAndByInfrastructure(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error {
		return EndpointInstanceList(ctx, "", EndpointInstanceFilters{ServiceStatus: []string{"active"}})
	})
	if !strings.Contains(out, `"label":"ei-1"`) {
		t.Errorf("global list output missing instance: %s", out)
	}
	if rec.last().Path != "/api/v2/endpoint-instances" {
		t.Errorf("global list should hit /endpoint-instances, got %s", rec.last().Path)
	}
	if !strings.Contains(rec.last().Query, "serviceStatus=active") {
		t.Errorf("service status filter should reach the query string: %s", rec.last().Query)
	}

	out = run(t, func() error {
		return EndpointInstanceList(ctx, "prod-env", EndpointInstanceFilters{GroupId: []string{"12"}})
	})
	if !strings.Contains(out, `"label":"ei-1"`) {
		t.Errorf("scoped list output missing instance: %s", out)
	}
	if rec.last().Path != "/api/v2/infrastructures/123/endpoint-instances" {
		t.Errorf("scoped list should hit the infrastructure endpoint, got %s", rec.last().Path)
	}
}

func TestEndpointInstanceGetAndConfig(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return EndpointInstanceGet(ctx, "42") })
	if !strings.Contains(out, `"id":42`) {
		t.Errorf("get output missing id: %s", out)
	}

	out = run(t, func() error { return EndpointInstanceConfigGet(ctx, "42") })
	if !strings.Contains(out, `"deployStatus":"not_started"`) {
		t.Errorf("config output missing deploy status: %s", out)
	}
	if rec.last().Path != "/api/v2/endpoint-instances/42/config" {
		t.Errorf("config should hit the /config endpoint, got %s", rec.last().Path)
	}
}

func TestEndpointInstanceGetInvalidId(t *testing.T) {
	srv := newServer(t, &recorder{})
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if err := EndpointInstanceGet(ctx, "not-a-number"); err == nil || !strings.Contains(err.Error(), "invalid endpoint instance ID") {
		t.Fatalf("expected invalid id error, got %v", err)
	}
}

func TestEndpointInstanceCreateResolvesInfrastructure(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	run(t, func() error {
		return EndpointInstanceCreate(ctx, "prod-env", sdkEndpointInstanceCreate(5, "ei-1"))
	})

	req := rec.last()
	if req.Method != http.MethodPost || req.Path != "/api/v2/infrastructures/123/endpoint-instances" {
		t.Fatalf("unexpected create request %s %s", req.Method, req.Path)
	}
	var body map[string]any
	_ = json.Unmarshal([]byte(req.Body), &body)
	if body["endpointId"] != float64(5) {
		t.Errorf("create body should carry the endpoint id: %s", req.Body)
	}
}

func TestEndpointInstanceDeleteSendsIfMatch(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if err := EndpointInstanceDelete(ctx, "42"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := rec.last()
	if req.Method != http.MethodDelete || req.Path != "/api/v2/endpoint-instances/42" {
		t.Fatalf("unexpected delete request %s %s", req.Method, req.Path)
	}
	if got := req.Header.Get("If-Match"); got != "7" {
		t.Errorf("expected If-Match 7 (entity revision), got %q", got)
	}
}

func TestEndpointInstanceConfigUpdateUsesConfigRevision(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	run(t, func() error {
		return EndpointInstanceConfigUpdate(ctx, "42", []byte(`{"label":"renamed"}`))
	})

	req := rec.last()
	if req.Method != http.MethodPatch || req.Path != "/api/v2/endpoint-instances/42/config" {
		t.Fatalf("unexpected update request %s %s", req.Method, req.Path)
	}
	if got := req.Header.Get("If-Match"); got != "3" {
		t.Errorf("expected If-Match 3 (config revision), got %q", got)
	}
	if !strings.Contains(req.Body, `"label":"renamed"`) {
		t.Errorf("update body should carry the config: %s", req.Body)
	}
}

func TestEndpointInstanceMetaUpdateHandlesNoContent(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if err := EndpointInstanceMetaUpdate(ctx, "42", []byte(`{"tags":["a","b"]}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := rec.last()
	if req.Path != "/api/v2/endpoint-instances/42/meta" {
		t.Fatalf("unexpected meta request path %s", req.Path)
	}
	if !strings.Contains(req.Body, `"tags":["a","b"]`) {
		t.Errorf("meta body should carry the tags: %s", req.Body)
	}
}

// ---------------------------------------------------------------------------
// Endpoint instance group tests
// ---------------------------------------------------------------------------

func TestEndpointInstanceGroupListAndGet(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error {
		return EndpointInstanceGroupList(ctx, "prod-env", EndpointInstanceGroupFilters{ServiceStatus: []string{"active"}})
	})
	if !strings.Contains(out, `"label":"eig-1"`) {
		t.Errorf("group list output missing group: %s", out)
	}
	if rec.last().Path != "/api/v2/infrastructures/123/endpoint-instance-groups" {
		t.Errorf("group list should hit the infrastructure endpoint, got %s", rec.last().Path)
	}

	out = run(t, func() error { return EndpointInstanceGroupGet(ctx, "12") })
	if !strings.Contains(out, `"id":12`) {
		t.Errorf("group get output missing id: %s", out)
	}
}

func TestEndpointInstanceGroupCreateAndDelete(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	run(t, func() error {
		return EndpointInstanceGroupCreate(ctx, "123", sdkEndpointInstanceGroupCreate("eig-1"))
	})
	req := rec.last()
	if req.Method != http.MethodPost || req.Path != "/api/v2/infrastructures/123/endpoint-instance-groups" {
		t.Fatalf("unexpected create request %s %s", req.Method, req.Path)
	}

	if err := EndpointInstanceGroupDelete(ctx, "12"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req = rec.last()
	if req.Method != http.MethodDelete || req.Header.Get("If-Match") != "9" {
		t.Errorf("delete should send the group revision as If-Match, got %s %q", req.Method, req.Header.Get("If-Match"))
	}
}

func TestEndpointInstanceGroupConfigUpdateUsesConfigRevision(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	run(t, func() error { return EndpointInstanceGroupConfigGet(ctx, "12") })

	run(t, func() error {
		return EndpointInstanceGroupConfigUpdate(ctx, "12", []byte(`{"label":"renamed"}`))
	})
	req := rec.last()
	if req.Method != http.MethodPatch || req.Path != "/api/v2/endpoint-instance-groups/12/config" {
		t.Fatalf("unexpected update request %s %s", req.Method, req.Path)
	}
	if got := req.Header.Get("If-Match"); got != "4" {
		t.Errorf("expected If-Match 4 (config revision), got %q", got)
	}

	if err := EndpointInstanceGroupMetaUpdate(ctx, "12", []byte(`{"tags":["x"]}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.last().Path != "/api/v2/endpoint-instance-groups/12/meta" {
		t.Errorf("unexpected meta path %s", rec.last().Path)
	}
}

func TestEndpointInstanceGroupInstances(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return EndpointInstanceGroupInstances(ctx, "12") })
	if !strings.Contains(out, `"label":"ei-1"`) {
		t.Errorf("instances output missing member: %s", out)
	}
	if rec.last().Path != "/api/v2/endpoint-instance-groups/12/endpoint-instances" {
		t.Errorf("unexpected instances path %s", rec.last().Path)
	}
}

// ---------------------------------------------------------------------------
// Network configuration tests
// ---------------------------------------------------------------------------

func TestEndpointInstanceGroupNetworkListAndReplace(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return EndpointInstanceGroupNetworkList(ctx, "12") })
	if !strings.Contains(out, `"name":"eig-1-neg"`) {
		t.Errorf("network list output missing network endpoint group: %s", out)
	}

	// Without a configuration the typed SDK call is used: PUT with no body.
	run(t, func() error { return EndpointInstanceGroupNetworkReplace(ctx, "12", nil) })
	req := rec.last()
	if req.Method != http.MethodPut || req.Path != "/api/v2/endpoint-instance-groups/12/config/networking" {
		t.Fatalf("unexpected replace request %s %s", req.Method, req.Path)
	}
	if req.Body != "" {
		t.Errorf("replace without a config should send no body, got %s", req.Body)
	}

	// With a configuration the raw request carries the body and the config revision.
	run(t, func() error {
		return EndpointInstanceGroupNetworkReplace(ctx, "12", []byte(`{"name":"new-neg"}`))
	})
	req = rec.last()
	if req.Method != http.MethodPut || req.Path != "/api/v2/endpoint-instance-groups/12/config/networking" {
		t.Fatalf("unexpected replace request %s %s", req.Method, req.Path)
	}
	if !strings.Contains(req.Body, `"name":"new-neg"`) {
		t.Errorf("replace body should carry the config: %s", req.Body)
	}
	if got := req.Header.Get("If-Match"); got != "4" {
		t.Errorf("expected If-Match 4 (group config revision), got %q", got)
	}
}

func TestEndpointInstanceGroupNetworkConnections(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return EndpointInstanceGroupNetworkConnections(ctx, "12") })
	if !strings.Contains(out, `"accessMode":"l2"`) {
		t.Errorf("connections output missing connection: %s", out)
	}

	out = run(t, func() error { return EndpointInstanceGroupNetworkGet(ctx, "12", "5") })
	if !strings.Contains(out, `"id":"5"`) {
		t.Errorf("connection get output missing id: %s", out)
	}

	run(t, func() error {
		return EndpointInstanceGroupNetworkConnect(ctx, "12", sdkNetworkConnectionCreate("7"))
	})
	req := rec.last()
	if req.Method != http.MethodPost || req.Path != "/api/v2/endpoint-instance-groups/12/config/networking/connections" {
		t.Fatalf("unexpected connect request %s %s", req.Method, req.Path)
	}
	if !strings.Contains(req.Body, `"logicalNetworkId":"7"`) {
		t.Errorf("connect body should carry the logical network: %s", req.Body)
	}

	run(t, func() error {
		return EndpointInstanceGroupNetworkUpdate(ctx, "12", "5", sdkNetworkConnectionUpdate("l2"))
	})
	req = rec.last()
	if req.Method != http.MethodPatch || req.Path != "/api/v2/endpoint-instance-groups/12/config/networking/connections/5" {
		t.Fatalf("unexpected connection update request %s %s", req.Method, req.Path)
	}

	if err := EndpointInstanceGroupNetworkDisconnect(ctx, "12", "5"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.last().Method != http.MethodDelete {
		t.Errorf("disconnect should issue a DELETE, got %s", rec.last().Method)
	}
}

// ---------------------------------------------------------------------------
// Security rule (ACL) tests
// ---------------------------------------------------------------------------

func TestEndpointInstanceGroupACLLifecycle(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return EndpointInstanceGroupACLList(ctx, "12", "5") })
	if !strings.Contains(out, `"sequence":10`) {
		t.Errorf("acl list output missing rule: %s", out)
	}
	if rec.last().Path != "/api/v2/endpoint-instance-groups/12/config/networking/connections/5/security/rules" {
		t.Errorf("unexpected acl list path %s", rec.last().Path)
	}

	out = run(t, func() error { return EndpointInstanceGroupACLGet(ctx, "12", "5", "3") })
	if !strings.Contains(out, `"id":"3"`) {
		t.Errorf("acl get output missing id: %s", out)
	}

	run(t, func() error {
		return EndpointInstanceGroupACLAdd(ctx, "12", "5", sdkACLCreate(20))
	})
	req := rec.last()
	if req.Method != http.MethodPost {
		t.Fatalf("acl add should POST, got %s", req.Method)
	}
	if !strings.Contains(req.Body, `"sequence":20`) {
		t.Errorf("acl add body should carry the rule: %s", req.Body)
	}

	run(t, func() error {
		return EndpointInstanceGroupACLUpdate(ctx, "12", "5", "3", []byte(`{"forwardingAction":"deny"}`))
	})
	req = rec.last()
	if req.Method != http.MethodPatch || !strings.Contains(req.Body, `"forwardingAction":"deny"`) {
		t.Errorf("unexpected acl update request %s %s", req.Method, req.Body)
	}

	if err := EndpointInstanceGroupACLRemove(ctx, "12", "5", "3"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.last().Method != http.MethodDelete {
		t.Errorf("acl remove should issue a DELETE, got %s", rec.last().Method)
	}
}

// TestEndpointInstanceGroupACLListSingleObject covers the SDK/API mismatch
// where the collection endpoint answers with a single object instead of a list.
func TestEndpointInstanceGroupACLListSingleObject(t *testing.T) {
	rec := &recorder{}
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/endpoint-instance-groups/12/config/networking/connections/5/security/rules": rec.wrap(
			testutils.JSONHandler(200, aclItem)),
	})
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return EndpointInstanceGroupACLList(ctx, "12", "5") })
	if !strings.Contains(out, `"id":"3"`) {
		t.Errorf("single object response should still render: %s", out)
	}
}

func TestEndpointInstanceListInfrastructureNotFound(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.JSONHandler(200, testutils.PaginatedResponse([]any{}, 1, 1)),
	})
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if err := EndpointInstanceList(ctx, "missing", EndpointInstanceFilters{}); err == nil {
		t.Fatal("expected an error when the infrastructure cannot be resolved")
	}
}

func TestEndpointInstanceApiError(t *testing.T) {
	rec := &recorder{}
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/endpoint-instances/42": rec.wrap(testutils.ErrorHandler(http.StatusInternalServerError, "boom")),
	})
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if err := EndpointInstanceGet(ctx, "42"); err == nil {
		t.Fatal("expected an error for HTTP 500")
	}
	if len(rec.paths()) != 1 {
		t.Errorf("expected exactly one request, got %v", rec.paths())
	}
}

// ---------------------------------------------------------------------------
// Payload builders
// ---------------------------------------------------------------------------

func sdkEndpointInstanceCreate(endpointId int64, label string) sdk.EndpointInstanceCreate {
	return sdk.EndpointInstanceCreate{
		Label:      sdk.PtrString(label),
		EndpointId: endpointId,
	}
}

func sdkEndpointInstanceGroupCreate(label string) sdk.EndpointInstanceGroupCreate {
	return sdk.EndpointInstanceGroupCreate{
		Label: sdk.PtrString(label),
	}
}

func sdkNetworkConnectionCreate(logicalNetworkId string) sdk.CreateEndpointInstanceGroupNetworkConnection {
	return sdk.CreateEndpointInstanceGroupNetworkConnection{
		LogicalNetworkId: logicalNetworkId,
		Tagged:           true,
		AccessMode:       sdk.NETWORKENDPOINTGROUPALLOWEDACCESSMODE_L2,
	}
}

func sdkNetworkConnectionUpdate(accessMode string) sdk.UpdateNetworkEndpointGroupLogicalNetwork {
	mode := sdk.NetworkEndpointGroupAllowedAccessMode(accessMode)
	return sdk.UpdateNetworkEndpointGroupLogicalNetwork{
		AccessMode: &mode,
	}
}

func sdkACLCreate(sequence int32) sdk.CreateLogicalNetworkACL {
	return sdk.CreateLogicalNetworkACL{
		RuleType:         sdk.ACLTYPE_IPV4,
		Direction:        sdk.ACLDIRECTION_IN,
		Sequence:         sequence,
		ForwardingAction: sdk.ACLFORWARDINGACTION_ALLOW,
		EnforcementPoint: sdk.ACLENFORCEMENTPOINT_SVI,
	}
}
