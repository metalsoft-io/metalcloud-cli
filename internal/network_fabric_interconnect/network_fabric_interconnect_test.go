package network_fabric_interconnect

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

var interconnectItem = map[string]any{
	"id": "12", "interconnectType": "dci-evpn", "label": "dc1-dc2",
	"name": "DC1 to DC2", "bgpConfigurationTemplateId": 3, "revision": "7",
	"status":           "draft",
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
}

var linkItem = map[string]any{
	"id": 3, "interconnectId": 12, "status": "draft", "fabricId": 7, "networkEquipmentId": 45,
	"revision": "2", "createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
}

var fabricItem = map[string]any{
	"id": "7", "name": "dc1-fabric", "siteId": 1, "revision": "1",
	"fabricConfiguration": map[string]any{"fabricType": "ethernet"},
	"createdTimestamp":    "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
}

var networkDeviceItem = testutils.NetworkDeviceFixture("45", 1, "border-leaf-01")

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

func newServer(t *testing.T, rec *recorder) *httptest.Server {
	t.Helper()
	routes := map[string]http.HandlerFunc{
		"/api/v2/network-fabric-interconnects":    testutils.JSONHandler(200, testutils.PaginatedResponse([]any{interconnectItem}, 1, 1)),
		"/api/v2/network-fabric-interconnects/12": testutils.JSONHandler(200, interconnectItem),
		"/api/v2/network-fabric-interconnects/12/links": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				testutils.JSONHandler(201, linkItem)(w, r)
				return
			}
			testutils.JSONHandler(200, testutils.PaginatedResponse([]any{linkItem}, 1, 1))(w, r)
		},
		"/api/v2/network-fabric-interconnects/12/links/3":                  testutils.JSONHandler(200, linkItem),
		"/api/v2/network-fabric-interconnects/12/fabrics":                  testutils.JSONHandler(200, map[string]any{"data": []any{fabricItem}}),
		"/api/v2/network-fabric-interconnects/12/fabrics-available":        testutils.JSONHandler(200, testutils.PaginatedResponse([]any{fabricItem}, 1, 1)),
		"/api/v2/network-fabric-interconnects/12/actions/deploy":           testutils.JSONHandler(200, map[string]any{"jobId": 99, "jobGroupId": 5}),
		"/api/v2/network-fabric-interconnects/12/actions/detach":           testutils.JSONHandler(200, map[string]any{"jobId": 98, "jobGroupId": 5}),
		"/api/v2/network-fabric-interconnects/12/actions/activate-links":   testutils.JSONHandler(200, map[string]any{"jobId": 97, "jobGroupId": 5}),
		"/api/v2/network-fabric-interconnects/12/actions/deactivate-links": testutils.JSONHandler(200, map[string]any{"jobId": 96, "jobGroupId": 5}),
		"/api/v2/network-fabric-interconnects/12/actions/accept-deploy":    testutils.NoContentHandler(),
		"/api/v2/network-fabric-interconnects/12/actions/reject-deploy":    testutils.NoContentHandler(),
		"/api/v2/network-fabric-interconnects/12/actions/deployment-check": testutils.JSONHandler(200, []any{map[string]any{
			"linkId": 3, "networkDeviceId": 45, "canActivate": true,
			"resolvedVariables": map[string]any{"router_id": "10.0.0.1", "internal_vtep_ip_address": "10.0.1.1", "external_vtep_ip_address": "10.0.2.1", "asn": 65000},
			"errors":            []any{},
		}}),
		"/api/v2/network-fabric-interconnects/12/deployment-info": testutils.JSONHandler(200, map[string]any{"status": "draft"}),
		"/api/v2/network-fabric-interconnects/template/dci-evpn": testutils.JSONHandler(200, map[string]any{
			"addGlobalConfigBgpTemplate": "a", "addNeighborBgpTemplate": "b",
			"removeGlobalConfigBgpTemplate": "c", "removeNeighborBgpTemplate": "d",
		}),
		"/api/v2/network-fabrics/7":  testutils.JSONHandler(200, fabricItem),
		"/api/v2/network-fabrics":    testutils.JSONHandler(200, testutils.PaginatedResponse([]any{fabricItem}, 1, 1)),
		"/api/v2/network-devices/45": testutils.JSONHandler(200, networkDeviceItem),
		"/api/v2/network-devices":    testutils.JSONHandler(200, testutils.PaginatedResponse([]any{networkDeviceItem}, 1, 1)),
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

func TestInterconnectListAndGet(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return InterconnectList(ctx, []string{"draft"}) })
	if !strings.Contains(out, "dc1-dc2") {
		t.Errorf("list output missing interconnect: %s", out)
	}

	out = run(t, func() error { return InterconnectGet(ctx, "12") })
	if !strings.Contains(out, `"label":"dc1-dc2"`) {
		t.Errorf("get output missing label: %s", out)
	}
	if rec.last().Path != "/api/v2/network-fabric-interconnects/12" {
		t.Errorf("get by numeric id should hit the by-id endpoint, got %s", rec.last().Path)
	}

	out = run(t, func() error { return InterconnectGet(ctx, "dc1-dc2") })
	if !strings.Contains(out, `"id":"12"`) {
		t.Errorf("get by label should resolve to id 12: %s", out)
	}
}

func TestInterconnectGetUnknownLabel(t *testing.T) {
	srv := newServer(t, &recorder{})
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if err := InterconnectGet(ctx, "does-not-exist"); err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestInterconnectUpdateSendsIfMatch(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	run(t, func() error { return InterconnectUpdate(ctx, "dc1-dc2", []byte(`{"description":"updated"}`)) })

	req := rec.last()
	if req.Method != http.MethodPatch || req.Path != "/api/v2/network-fabric-interconnects/12" {
		t.Fatalf("unexpected request %s %s", req.Method, req.Path)
	}
	if got := req.Header.Get("If-Match"); got != "7" {
		t.Errorf("expected If-Match 7, got %q", got)
	}
	var body map[string]any
	_ = json.Unmarshal([]byte(req.Body), &body)
	if body["description"] != "updated" {
		t.Errorf("config should be sent as JSON body, got %s", req.Body)
	}
}

func TestInterconnectDeployOptions(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return InterconnectDeploy(ctx, "12", true) })
	if !strings.Contains(out, `"jobId":99`) {
		t.Errorf("deploy should print the job info: %s", out)
	}
	if !strings.Contains(rec.last().Body, `"requireConfirmation":true`) {
		t.Errorf("deploy body should carry requireConfirmation: %s", rec.last().Body)
	}
}

func TestInterconnectAcceptRejectDetach(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	run(t, func() error { return InterconnectAcceptDeploy(ctx, "12") })
	if rec.last().Path != "/api/v2/network-fabric-interconnects/12/actions/accept-deploy" {
		t.Errorf("unexpected accept path %s", rec.last().Path)
	}
	run(t, func() error { return InterconnectRejectDeploy(ctx, "12") })
	if rec.last().Path != "/api/v2/network-fabric-interconnects/12/actions/reject-deploy" {
		t.Errorf("unexpected reject path %s", rec.last().Path)
	}
	out := run(t, func() error { return InterconnectDetach(ctx, "12") })
	if !strings.Contains(out, `"jobId":98`) {
		t.Errorf("detach should print the job info: %s", out)
	}
}

func TestInterconnectDeploymentCheckAndInfo(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return InterconnectDeploymentCheck(ctx, "12", []string{"3"}) })
	if !strings.Contains(out, `"canActivate":true`) {
		t.Errorf("check output missing validation: %s", out)
	}
	if !strings.Contains(rec.last().Body, `"linkIds":[3]`) {
		t.Errorf("check body should carry link ids: %s", rec.last().Body)
	}

	if err := InterconnectDeploymentCheck(ctx, "12", []string{"x"}); err == nil {
		t.Error("expected error for non-numeric link id")
	}

	out = run(t, func() error { return InterconnectDeploymentInfo(ctx, "12") })
	if !strings.Contains(out, `"status":"draft"`) {
		t.Errorf("info output missing status: %s", out)
	}
}

func TestInterconnectTemplate(t *testing.T) {
	srv := newServer(t, &recorder{})
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return InterconnectTemplateGet(ctx, "dci-evpn") })
	if !strings.Contains(out, "addGlobalConfigBgpTemplate") {
		t.Errorf("template output missing fields: %s", out)
	}
	if err := InterconnectTemplateGet(ctx, "bogus"); err == nil {
		t.Error("expected error for invalid interconnect type")
	}
}

func TestInterconnectFabrics(t *testing.T) {
	srv := newServer(t, &recorder{})
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return InterconnectFabricsGet(ctx, "12") })
	if !strings.Contains(out, "dc1-fabric") {
		t.Errorf("fabrics output missing fabric: %s", out)
	}
	out = run(t, func() error { return InterconnectAvailableFabricsGet(ctx, "12") })
	if !strings.Contains(out, "dc1-fabric") {
		t.Errorf("available fabrics output missing fabric: %s", out)
	}
}

func TestInterconnectLinks(t *testing.T) {
	rec := &recorder{}
	srv := newServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := run(t, func() error { return InterconnectLinksGet(ctx, "12") })
	if !strings.Contains(out, `"networkEquipmentId":45`) {
		t.Errorf("links output missing link: %s", out)
	}

	run(t, func() error { return InterconnectLinkGet(ctx, "12", "3") })
	if rec.last().Path != "/api/v2/network-fabric-interconnects/12/links/3" {
		t.Errorf("unexpected link get path %s", rec.last().Path)
	}

	// add-link resolves the fabric by name and the device by identifier.
	run(t, func() error { return InterconnectLinkAdd(ctx, "dc1-dc2", "dc1-fabric", "border-leaf-01") })
	req := rec.last()
	if req.Method != http.MethodPost || req.Path != "/api/v2/network-fabric-interconnects/12/links" {
		t.Fatalf("unexpected add-link request %s %s", req.Method, req.Path)
	}
	if !strings.Contains(req.Body, `"fabricId":7`) || !strings.Contains(req.Body, `"networkEquipmentId":45`) {
		t.Errorf("add-link body should carry resolved ids: %s", req.Body)
	}

	// remove-link fetches the link first to obtain its revision.
	run(t, func() error { return InterconnectLinkRemove(ctx, "12", "3") })
	req = rec.last()
	if req.Method != http.MethodDelete || req.Header.Get("If-Match") != "2" {
		t.Errorf("remove-link should DELETE with If-Match 2, got %s %q", req.Method, req.Header.Get("If-Match"))
	}

	run(t, func() error { return InterconnectLinksActivate(ctx, "12", []string{"3", "4"}, false) })
	if !strings.Contains(rec.last().Body, `"linkIds":[3,4]`) {
		t.Errorf("activate body should carry link ids: %s", rec.last().Body)
	}
	run(t, func() error { return InterconnectLinksDeactivate(ctx, "12", []string{"3"}, true) })
	if !strings.Contains(rec.last().Body, `"requireConfirmation":true`) {
		t.Errorf("deactivate body should carry requireConfirmation: %s", rec.last().Body)
	}
}

func TestInterconnectConfigExample(t *testing.T) {
	out := run(t, func() error { return InterconnectConfigExample(testutils.SetupTestContext("http://unused")) })
	if !strings.Contains(out, `"interconnectType":"dci-evpn"`) {
		t.Errorf("config example missing type: %s", out)
	}

	viper.Set(formatter.ConfigFormat, "yaml")
	defer testutils.SetupTestFormat()
	out = run(t, func() error { return InterconnectConfigExample(testutils.SetupTestContext("http://unused")) })
	if !strings.Contains(out, "interconnectType: dci-evpn") {
		t.Errorf("yaml config example missing type: %s", out)
	}
}
