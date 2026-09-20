package server_instance

import (
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

var siInfrastructureItem = map[string]any{
	"id": 123, "revision": 1, "label": "prod-env", "serviceStatus": "active",
	"datacenterName": "dc1", "siteId": 1, "designIsLocked": 0, "config": map[string]any{},
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
}

var siInstanceItem = map[string]any{
	"id": 42, "revision": 7, "label": "si-1", "infrastructureId": 123, "groupId": 12,
	"serviceStatus": "active", "isVmInstance": 0, "isEndpointInstance": 0,
	"meta":             map[string]any{"tags": []any{"prod"}},
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

var siInstanceConfigItem = map[string]any{
	"revision": 3, "label": "si-1", "updatedTimestamp": "2024-01-02T00:00:00Z",
	"groupId": 12, "serverTypeId": 5, "serverId": 100, "osTemplateId": 1,
	"hostname": "si-1.example.com", "deployType": "create", "deployStatus": "not_started",
}

var siDriveItem = map[string]any{
	"id": 8, "revision": 1, "label": "drive-1", "groupId": 4, "sizeMb": 40960,
	"storageType": "iscsi_ssd", "infrastructureId": 123,
	"infrastructure": map[string]any{"id": 123, "label": "prod-env"},
	"serviceStatus":  "active", "storageUpdatedTimestamp": "2024-01-02T00:00:00Z",
	"provisioningProtocol": "iscsi", "meta": map[string]any{"name": "drive-1"},
	"config": map[string]any{
		"revision": 1, "label": "drive-1", "groupId": 4, "sizeMb": 40960,
		"storageType": "iscsi_ssd", "updatedTimestamp": "2024-01-02T00:00:00Z",
		"deployType": "create", "deployStatus": "not_started",
	},
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

var siDriveGroupItem = map[string]any{
	"id": 4, "revision": 1, "label": "drive-group-1", "infrastructureId": 123,
	"driveSizeMbDefault": 40960, "expandWithServerInstanceGroup": 1,
	"storageType": "iscsi_ssd", "serviceStatus": "active", "allocationAffinity": "same_storage",
	"meta": map[string]any{"name": "drive-group-1"},
	"config": map[string]any{
		"revision": 1, "label": "drive-group-1", "infrastructureId": 123,
		"driveSizeMbDefault": 40960, "expandWithServerInstanceGroup": 1,
		"storageType": "iscsi_ssd", "updatedTimestamp": "2024-01-02T00:00:00Z",
		"deployType": "create", "deployStatus": "not_started",
	},
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

var siInterfaceItem = map[string]any{
	"id": 7, "revision": 2, "label": "iface-1", "infrastructureId": 123,
	"instanceId": 42, "index": 0, "capacityMbps": 10000, "dirtyBit": false,
	"serviceStatus": "active", "networkId": 9,
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
	"config": map[string]any{
		"revision": 5, "label": "iface-1", "updatedTimestamp": "2024-01-02T00:00:00Z",
		"instanceId": 42, "index": 0, "capacityMbps": 10000,
		"deployType": "create", "deployStatus": "not_started",
	},
}

var siInterfaceConfigItem = map[string]any{
	"revision": 5, "label": "iface-1", "updatedTimestamp": "2024-01-02T00:00:00Z",
	"instanceId": 42, "index": 0, "capacityMbps": 10000,
	"deployType": "create", "deployStatus": "not_started",
}

var sigGroupItem = map[string]any{
	"id": 12, "revision": 9, "label": "sig-1", "infrastructureId": 123,
	"instanceCount": 1, "defaultServerTypeId": 5,
	"ipAllocateAuto": 1, "ipv4SubnetCreateAuto": 1,
	"processorCount": 2, "processorCoreCount": 8, "processorCoreMhz": 2400,
	"diskCount": 2, "diskSizeMbytes": 102400, "diskTypes": []any{},
	"virtualInterfacesEnabled": 0, "serviceStatus": "active",
	"isVmGroup": 0, "isEndpointInstanceGroup": 0, "meta": map[string]any{},
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

var sigGroupConfigItem = map[string]any{
	"revision": 4, "label": "sig-1", "updatedTimestamp": "2024-01-02T00:00:00Z",
	"instanceCount": 1, "defaultServerTypeId": 5,
	"ipAllocateAuto": 1, "ipv4SubnetCreateAuto": 1,
	"processorCount": 2, "processorCoreCount": 8, "processorCoreMhz": 2400,
	"diskCount": 2, "diskSizeMbytes": 102400, "diskTypes": []any{},
	"virtualInterfacesEnabled": 0,
	"deployType":               "create", "deployStatus": "not_started",
}

var sigInterfaceItem = map[string]any{
	"id": 7, "revision": 2, "label": "sig-iface-1", "infrastructureId": 123,
	"groupId": 12, "index": 0, "serviceStatus": "active", "networkId": 9,
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

var sigNetworkEndpointGroupItem = map[string]any{
	"id": "77", "name": "sig-1-neg", "siteId": 1, "revision": "2",
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

var sigNetworkConnectionItem = map[string]any{
	"id": "5", "tagged": true, "accessMode": "l2", "mtu": 1500,
	"providesDefaultRoute": false,
}

var sigACLItem = map[string]any{
	"id": "3", "ruleType": "ipv4", "direction": "in", "sequence": 10,
	"forwardingAction": "allow", "enforcementPoint": "svi",
	"sourceAddress": "10.0.0.0/24", "destinationAddress": "10.0.1.0/24",
	"endpointGroupId": "77", "logicalNetworkId": "7",
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
}

var siStatisticsItem = map[string]any{
	"serverStatus": map[string]any{"available": 3, "used": 1},
	"site":         map[string]any{"dc1": 4},
}

// siVariablesItem satisfies sdk.ServerInstanceContextVariables.
var siVariablesItem = map[string]any{
	"site":       map[string]any{"id": 1, "revision": 1, "slug": "dc1", "name": "dc1"},
	"siteConfig": map[string]any{},
	"server": map[string]any{
		"serverId": 100, "revision": 1, "siteId": 1, "datacenterName": "dc1",
		"bdkDebug": 0, "requiresReRegister": 0, "serverClass": "bigdata",
		"serverStatus": "available", "administrationState": "active",
		"serverDhcpStatus": "quarantine", "supportsFcProvisioning": 0,
		"serverCreatedTimestamp": "2024-01-01T00:00:00Z", "powerStatus": "on",
		"powerStatusLastUpdateTimestamp": "2024-01-02T00:00:00Z",
	},
	"serverInstance": map[string]any{
		"id": 42, "revision": 7, "label": "si-1",
		"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
		"infrastructureId": 123, "groupId": 12, "serviceStatus": "active",
		"isVmInstance": 0, "isEndpointInstance": 0,
	},
	"serverInstanceGroup": map[string]any{
		"id": 12, "revision": 9, "label": "sig-1",
		"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-02T00:00:00Z",
		"infrastructureId": 123, "instanceCount": 1, "defaultServerTypeId": 5,
		"ipAllocateAuto": 1, "ipv4SubnetCreateAuto": 1,
		"processorCount": 2, "processorCoreCount": 8, "processorCoreMhz": 2400,
		"diskCount": 2, "diskSizeMbytes": 102400, "diskTypes": []any{},
		"virtualInterfacesEnabled": 0, "serviceStatus": "active",
		"isVmGroup": 0, "isEndpointInstanceGroup": 0,
	},
	"infrastructure": map[string]any{
		"id": 123, "revision": 1, "label": "prod-env", "serviceStatus": "active",
		"datacenterName": "dc1", "siteId": 1, "designIsLocked": 0,
		"config":           map[string]any{},
		"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
	},
	"driveGroups": []any{}, "drives": []any{}, "fileShares": []any{},
	"buckets": []any{}, "sharedDrives": []any{}, "network": []any{},
}

// siOSInstallationDataItem satisfies sdk.ServerInstanceContextOSInstallationData,
// whose sub-models differ from the variables ones despite the same key set.
var siOSInstallationDataItem = map[string]any{
	"site": map[string]any{"id": 1, "slug": "dc1", "name": "dc1"},
	"siteConfig": map[string]any{
		"repoURL": "http://repo.example.com", "DNSServers": []any{"8.8.8.8"},
		"NTPServers": []any{"pool.ntp.org"},
	},
	"server": map[string]any{"serverId": 100, "passwordEncrypted": "secret"},
	"serverInstance": map[string]any{
		"id": 42, "label": "si-1", "isVmInstance": 0, "isEndpointInstance": 0,
		"provisionInstanceDnsRecords": false,
	},
	"serverInstanceGroup": map[string]any{
		"id": 12, "label": "sig-1", "isVmGroup": 0, "isEndpointInstanceGroup": 0,
	},
	"infrastructure": map[string]any{"label": "prod-env"},
	"driveGroups":    []any{}, "drives": []any{}, "fileShares": []any{},
	"buckets": []any{}, "sharedDrives": []any{}, "network": []any{},
}

// ---------------------------------------------------------------------------
// Request recorder
// ---------------------------------------------------------------------------

type siRecordedRequest struct {
	Method string
	Path   string
	Query  string
	Header http.Header
	Body   string
}

type siRecorder struct {
	mu   sync.Mutex
	reqs []siRecordedRequest
}

func (r *siRecorder) wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		r.mu.Lock()
		r.reqs = append(r.reqs, siRecordedRequest{
			Method: req.Method, Path: req.URL.Path, Query: req.URL.RawQuery,
			Header: req.Header.Clone(), Body: string(body),
		})
		r.mu.Unlock()
		next(w, req)
	}
}

func (r *siRecorder) last() siRecordedRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.reqs[len(r.reqs)-1]
}

// find returns the last recorded request for method and path.
func (r *siRecorder) find(method, path string) (siRecordedRequest, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := len(r.reqs) - 1; i >= 0; i-- {
		if r.reqs[i].Method == method && r.reqs[i].Path == path {
			return r.reqs[i], true
		}
	}
	return siRecordedRequest{}, false
}

// siMethodSwitch serves readHandler for GET and writeHandler for every other verb.
func siMethodSwitch(readHandler, writeHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			readHandler(w, r)
			return
		}
		writeHandler(w, r)
	}
}

// siDeleteOr serves NoContent for DELETE and okHandler otherwise.
func siDeleteOr(okHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			testutils.NoContentHandler()(w, r)
			return
		}
		okHandler(w, r)
	}
}

func newSIServer(t *testing.T, rec *siRecorder) *httptest.Server {
	t.Helper()
	routes := map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.JSONHandler(200, testutils.PaginatedResponse([]any{siInfrastructureItem}, 1, 1)),

		"/api/v2/server-instances": testutils.JSONHandler(200, testutils.PaginatedResponse([]any{siInstanceItem}, 1, 1)),
		"/api/v2/infrastructures/123/server-instances": siMethodSwitch(
			testutils.JSONHandler(200, testutils.PaginatedResponse([]any{siInstanceItem}, 1, 1)),
			testutils.JSONHandler(201, siInstanceItem),
		),
		"/api/v2/server-instances/statistics": testutils.JSONHandler(200, siStatisticsItem),
		"/api/v2/server-instances/42": siMethodSwitch(
			testutils.JSONHandler(200, siInstanceItem),
			testutils.NoContentHandler(),
		),
		"/api/v2/server-instances/42/config": siMethodSwitch(
			testutils.JSONHandler(200, siInstanceConfigItem),
			testutils.JSONHandler(200, siInstanceConfigItem),
		),
		"/api/v2/server-instances/42/meta":                 testutils.NoContentHandler(),
		"/api/v2/server-instances/42/actions/reset":        testutils.NoContentHandler(),
		"/api/v2/server-instances/42/drives":               testutils.JSONHandler(200, map[string]any{"data": []any{siDriveItem}}),
		"/api/v2/server-instances/42/interfaces":           testutils.JSONHandler(200, testutils.PaginatedResponse([]any{siInterfaceItem}, 1, 1)),
		"/api/v2/server-instances/42/interfaces/7":         testutils.JSONHandler(200, siInterfaceItem),
		"/api/v2/server-instances/42/interfaces/7/config":  testutils.JSONHandler(200, siInterfaceConfigItem),
		"/api/v2/server-instances/42/os-installation-data": testutils.JSONHandler(200, siOSInstallationDataItem),
		"/api/v2/server-instances/42/variables":            testutils.JSONHandler(200, siVariablesItem),
		"/api/v2/infrastructures/123/actions/power-get":    testutils.JSONHandler(200, map[string]any{"42": "on"}),
		"/api/v2/infrastructures/123/actions/power-set":    testutils.NoContentHandler(),
		"/api/v2/server-instance-groups/12":                siMethodSwitch(testutils.JSONHandler(200, sigGroupItem), testutils.NoContentHandler()),
		"/api/v2/server-instance-groups/12/config":         testutils.JSONHandler(200, sigGroupConfigItem),
		"/api/v2/server-instance-groups/12/meta":           testutils.NoContentHandler(),
		"/api/v2/server-instance-groups/12/drive-groups":   testutils.JSONHandler(200, map[string]any{"data": []any{siDriveGroupItem}}),
		"/api/v2/server-instance-groups/12/interfaces":     testutils.JSONHandler(200, testutils.PaginatedResponse([]any{sigInterfaceItem}, 1, 1)),
		"/api/v2/server-instance-groups/12/interfaces/7":   testutils.JSONHandler(200, sigInterfaceItem),
		"/api/v2/server-instance-groups/12/config/networking": siMethodSwitch(
			testutils.JSONHandler(200, sigNetworkEndpointGroupItem),
			testutils.JSONHandler(200, sigNetworkEndpointGroupItem),
		),
		"/api/v2/server-instance-groups/12/config/networking/connections": siMethodSwitch(
			testutils.JSONHandler(200, map[string]any{"data": []any{sigNetworkConnectionItem}}),
			testutils.JSONHandler(201, sigNetworkConnectionItem),
		),
		"/api/v2/server-instance-groups/12/config/networking/connections/5": siMethodSwitch(
			testutils.JSONHandler(200, sigNetworkConnectionItem),
			siDeleteOr(testutils.JSONHandler(200, sigNetworkConnectionItem)),
		),
		"/api/v2/server-instance-groups/12/config/networking/connections/5/security/rules": siMethodSwitch(
			testutils.JSONHandler(200, []any{sigACLItem}),
			testutils.JSONHandler(201, sigACLItem),
		),
		"/api/v2/server-instance-groups/12/config/networking/connections/5/security/rules/3": siMethodSwitch(
			testutils.JSONHandler(200, sigACLItem),
			siDeleteOr(testutils.JSONHandler(200, sigACLItem)),
		),
	}

	wrapped := make(map[string]http.HandlerFunc, len(routes))
	for path, handler := range routes {
		wrapped[path] = rec.wrap(handler)
	}
	return testutils.NewTestServer(wrapped)
}

func siRun(t *testing.T, fn func() error) string {
	t.Helper()
	var err error
	out := testutils.CaptureStdout(t, func() { err = fn() })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return out
}

// ---------------------------------------------------------------------------
// Server instance tests
// ---------------------------------------------------------------------------

func TestServerInstanceListGlobalAndByInfrastructure(t *testing.T) {
	rec := &siRecorder{}
	srv := newSIServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := siRun(t, func() error {
		return ServerInstanceList(ctx, "", ServerInstanceFilters{ServiceStatus: []string{"active"}})
	})
	if !strings.Contains(out, `"label":"si-1"`) {
		t.Errorf("global list output missing instance: %s", out)
	}
	if rec.last().Path != "/api/v2/server-instances" {
		t.Errorf("global list should hit /server-instances, got %s", rec.last().Path)
	}
	if !strings.Contains(rec.last().Query, "serviceStatus=active") {
		t.Errorf("service status filter should reach the query string: %s", rec.last().Query)
	}

	siRun(t, func() error {
		return ServerInstanceList(ctx, "prod-env", ServerInstanceFilters{GroupId: []string{"12"}})
	})
	if rec.last().Path != "/api/v2/infrastructures/123/server-instances" {
		t.Errorf("scoped list should hit the infrastructure collection, got %s", rec.last().Path)
	}
	if !strings.Contains(rec.last().Query, "groupId=12") {
		t.Errorf("group filter should reach the query string: %s", rec.last().Query)
	}
}

func TestServerInstanceCreate(t *testing.T) {
	rec := &siRecorder{}
	srv := newSIServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	siRun(t, func() error {
		return ServerInstanceCreate(ctx, "prod-env", sdk.ServerInstanceCreate{
			Label:   sdk.PtrString("si-1"),
			GroupId: sdk.PtrInt64(12),
		})
	})

	req, ok := rec.find(http.MethodPost, "/api/v2/infrastructures/123/server-instances")
	if !ok {
		t.Fatal("create should POST to the infrastructure server instance collection")
	}
	if !strings.Contains(req.Body, `"label":"si-1"`) {
		t.Errorf("create body should carry the label: %s", req.Body)
	}
}

func TestServerInstanceDeleteSendsIfMatch(t *testing.T) {
	rec := &siRecorder{}
	srv := newSIServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if err := ServerInstanceDelete(ctx, "42"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req, ok := rec.find(http.MethodDelete, "/api/v2/server-instances/42")
	if !ok {
		t.Fatal("delete should DELETE the server instance")
	}
	if got := req.Header.Get("If-Match"); got != "7" {
		t.Errorf("delete should send the instance revision as If-Match, got %q", got)
	}
}

func TestServerInstanceConfigUpdateUsesConfigRevision(t *testing.T) {
	rec := &siRecorder{}
	srv := newSIServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	siRun(t, func() error {
		return ServerInstanceConfigUpdate(ctx, "42", []byte(`{"label":"si-renamed"}`))
	})

	req, ok := rec.find(http.MethodPatch, "/api/v2/server-instances/42/config")
	if !ok {
		t.Fatal("update-config should PATCH the configuration")
	}
	// The configuration revision (3) guards writes below /config, not the
	// instance revision (7).
	if got := req.Header.Get("If-Match"); got != "3" {
		t.Errorf("update-config should send the config revision as If-Match, got %q", got)
	}
	if !strings.Contains(req.Body, `"label":"si-renamed"`) {
		t.Errorf("update-config body should carry the new label: %s", req.Body)
	}
}

func TestServerInstanceMetaUpdate(t *testing.T) {
	rec := &siRecorder{}
	srv := newSIServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if err := ServerInstanceMetaUpdate(ctx, "42", []byte(`{"tags":["prod"]}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req, ok := rec.find(http.MethodPatch, "/api/v2/server-instances/42/meta")
	if !ok {
		t.Fatal("update-meta should PATCH the metadata")
	}
	if !strings.Contains(req.Body, `"tags"`) {
		t.Errorf("update-meta body should carry the tags: %s", req.Body)
	}
}

func TestServerInstanceReset(t *testing.T) {
	rec := &siRecorder{}
	srv := newSIServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	siRun(t, func() error { return ServerInstanceReset(ctx, "42") })

	req, ok := rec.find(http.MethodPost, "/api/v2/server-instances/42/actions/reset")
	if !ok {
		t.Fatal("reset should POST to the reset action")
	}
	if got := req.Header.Get("If-Match"); got != "7" {
		t.Errorf("reset should send the instance revision as If-Match, got %q", got)
	}
}

func TestServerInstanceDrivesInterfacesAndVariables(t *testing.T) {
	rec := &siRecorder{}
	srv := newSIServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if out := siRun(t, func() error { return ServerInstanceDrives(ctx, "42") }); !strings.Contains(out, `"label":"drive-1"`) {
		t.Errorf("drives output missing drive: %s", out)
	}

	out := siRun(t, func() error {
		return ServerInstanceInterfaces(ctx, "42", ServerInstanceFilters{ServiceStatus: []string{"active"}})
	})
	if !strings.Contains(out, `"label":"iface-1"`) {
		t.Errorf("interfaces output missing interface: %s", out)
	}
	if !strings.Contains(rec.last().Query, "serviceStatus=active") {
		t.Errorf("interface filter should reach the query string: %s", rec.last().Query)
	}

	if out := siRun(t, func() error { return ServerInstanceInterfaceGet(ctx, "42", "7") }); !strings.Contains(out, `"id":7`) {
		t.Errorf("interface get output missing interface: %s", out)
	}

	siRun(t, func() error { return ServerInstanceVariables(ctx, "42", "AnsibleBundle") })
	if !strings.Contains(rec.last().Query, "usage=AnsibleBundle") {
		t.Errorf("usage should reach the query string: %s", rec.last().Query)
	}

	siRun(t, func() error { return ServerInstanceOSInstallationData(ctx, "42", "") })
	if rec.last().Path != "/api/v2/server-instances/42/os-installation-data" {
		t.Errorf("os-installation-data hit the wrong path: %s", rec.last().Path)
	}
	if rec.last().Query != "" {
		t.Errorf("no usage flag should mean no usage query parameter: %s", rec.last().Query)
	}
}

func TestServerInstanceInterfaceConfigUpdateUsesConfigRevision(t *testing.T) {
	rec := &siRecorder{}
	srv := newSIServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	siRun(t, func() error {
		return ServerInstanceInterfaceConfigUpdate(ctx, "42", "7", []byte(`{"networkId":9}`))
	})

	req, ok := rec.find(http.MethodPatch, "/api/v2/server-instances/42/interfaces/7/config")
	if !ok {
		t.Fatal("update-interface-config should PATCH the interface configuration")
	}
	// The interface configuration revision (5) guards the write, not the
	// interface revision (2).
	if got := req.Header.Get("If-Match"); got != "5" {
		t.Errorf("update-interface-config should send the config revision, got %q", got)
	}
}

func TestServerInstanceStatistics(t *testing.T) {
	rec := &siRecorder{}
	srv := newSIServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	out := siRun(t, func() error { return ServerInstanceStatistics(ctx) })
	if !strings.Contains(out, "available") {
		t.Errorf("statistics output missing server status counts: %s", out)
	}
	if rec.last().Path != "/api/v2/server-instances/statistics" {
		t.Errorf("statistics hit the wrong path: %s", rec.last().Path)
	}
}

func TestServerInstancePowerBatchAlwaysSendsBody(t *testing.T) {
	rec := &siRecorder{}
	srv := newSIServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	siRun(t, func() error {
		return ServerInstancePowerStatusBatch(ctx, "prod-env", []string{"42", "43"})
	})
	req, ok := rec.find(http.MethodPost, "/api/v2/infrastructures/123/actions/power-get")
	if !ok {
		t.Fatal("power-status-batch should POST to the power-get action")
	}
	if strings.TrimSpace(req.Body) != `["42","43"]` {
		t.Errorf("power-status-batch should send the instance ids as body, got %q", req.Body)
	}

	if err := ServerInstancePowerSetBatch(ctx, "prod-env", "off", []string{"42"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	req, ok = rec.find(http.MethodPost, "/api/v2/infrastructures/123/actions/power-set")
	if !ok {
		t.Fatal("power-set-batch should POST to the power-set action")
	}
	if !strings.Contains(req.Body, `"powerCommand":"off"`) || !strings.Contains(req.Body, `"instances":["42"]`) {
		t.Errorf("power-set-batch body should carry instances and command, got %q", req.Body)
	}
}

func TestServerInstancePowerBatchRejectsBadInput(t *testing.T) {
	rec := &siRecorder{}
	srv := newSIServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if err := ServerInstancePowerStatusBatch(ctx, "prod-env", nil); err == nil {
		t.Error("an empty instance list should be rejected")
	}
	if err := ServerInstancePowerSetBatch(ctx, "prod-env", "explode", []string{"42"}); err == nil {
		t.Error("an unknown power command should be rejected")
	}
	if err := ServerInstancePowerStatusBatch(ctx, "prod-env", []string{"not-a-number"}); err == nil {
		t.Error("a non numeric instance id should be rejected")
	}
}

// ---------------------------------------------------------------------------
// Server instance group tests
// ---------------------------------------------------------------------------

func TestServerInstanceGroupMetaDriveGroupsAndInterfaces(t *testing.T) {
	rec := &siRecorder{}
	srv := newSIServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if err := ServerInstanceGroupMetaUpdate(ctx, "12", []byte(`{"tags":["prod"]}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := rec.find(http.MethodPatch, "/api/v2/server-instance-groups/12/meta"); !ok {
		t.Error("group update-meta should PATCH the metadata")
	}

	if out := siRun(t, func() error { return ServerInstanceGroupDriveGroups(ctx, "12") }); !strings.Contains(out, `"label":"drive-group-1"`) {
		t.Errorf("drive-groups output missing drive group: %s", out)
	}

	out := siRun(t, func() error {
		return ServerInstanceGroupInterfaces(ctx, "12", ServerInstanceFilters{ServiceStatus: []string{"active"}})
	})
	if !strings.Contains(out, `"label":"sig-iface-1"`) {
		t.Errorf("group interfaces output missing interface: %s", out)
	}

	if out := siRun(t, func() error { return ServerInstanceGroupInterfaceGet(ctx, "12", "7") }); !strings.Contains(out, `"id":7`) {
		t.Errorf("group interface get output missing interface: %s", out)
	}
}

func TestServerInstanceGroupNetworkConfigurationAndReplace(t *testing.T) {
	rec := &siRecorder{}
	srv := newSIServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	if out := siRun(t, func() error { return ServerInstanceGroupNetworkConfiguration(ctx, "12") }); !strings.Contains(out, `"name":"sig-1-neg"`) {
		t.Errorf("network config output missing endpoint group: %s", out)
	}

	// Without a configuration the typed SDK call is used and sends no body.
	siRun(t, func() error { return ServerInstanceGroupNetworkReplace(ctx, "12", nil) })
	req, ok := rec.find(http.MethodPut, "/api/v2/server-instance-groups/12/config/networking")
	if !ok {
		t.Fatal("network replace should PUT the network configuration")
	}
	if strings.TrimSpace(req.Body) != "" {
		t.Errorf("bodyless replace should send no body, got %q", req.Body)
	}

	// With a configuration the raw request carries the group config revision.
	siRun(t, func() error {
		return ServerInstanceGroupNetworkReplace(ctx, "12", []byte(`{"name":"sig-1-neg"}`))
	})
	req, _ = rec.find(http.MethodPut, "/api/v2/server-instance-groups/12/config/networking")
	if got := req.Header.Get("If-Match"); got != "4" {
		t.Errorf("replace should send the group config revision as If-Match, got %q", got)
	}
	if !strings.Contains(req.Body, `"name":"sig-1-neg"`) {
		t.Errorf("replace body should be forwarded verbatim, got %q", req.Body)
	}
}

func TestServerInstanceGroupACLLifecycle(t *testing.T) {
	rec := &siRecorder{}
	srv := newSIServer(t, rec)
	defer srv.Close()
	ctx := testutils.SetupTestContext(srv.URL)

	// The SDK types the rules collection as a single object, so the listing
	// goes through the raw fallback.
	if out := siRun(t, func() error { return ServerInstanceGroupACLList(ctx, "12", "5") }); !strings.Contains(out, `"id":"3"`) {
		t.Errorf("acl list output missing rule: %s", out)
	}

	if out := siRun(t, func() error { return ServerInstanceGroupACLGet(ctx, "12", "5", "3") }); !strings.Contains(out, `"id":"3"`) {
		t.Errorf("acl get output missing rule: %s", out)
	}

	siRun(t, func() error { return ServerInstanceGroupACLConfigExample(ctx) })

	siRun(t, func() error {
		return ServerInstanceGroupACLAdd(ctx, "12", "5", sdk.CreateLogicalNetworkACL{
			RuleType:         sdk.ACLTYPE_IPV4,
			Direction:        sdk.ACLDIRECTION_IN,
			Sequence:         10,
			ForwardingAction: sdk.ACLFORWARDINGACTION_ALLOW,
			EnforcementPoint: sdk.ACLENFORCEMENTPOINT_SVI,
		})
	})
	req, ok := rec.find(http.MethodPost, "/api/v2/server-instance-groups/12/config/networking/connections/5/security/rules")
	if !ok {
		t.Fatal("acl add should POST to the rules collection")
	}
	if !strings.Contains(req.Body, `"ruleType":"ipv4"`) {
		t.Errorf("acl add body should carry the rule type: %s", req.Body)
	}

	siRun(t, func() error {
		return ServerInstanceGroupACLUpdate(ctx, "12", "5", "3", []byte(`{"forwardingAction":"deny"}`))
	})
	req, ok = rec.find(http.MethodPatch, "/api/v2/server-instance-groups/12/config/networking/connections/5/security/rules/3")
	if !ok {
		t.Fatal("acl update should PATCH the rule")
	}
	if !strings.Contains(req.Body, `"forwardingAction":"deny"`) {
		t.Errorf("acl update body should carry the new action: %s", req.Body)
	}

	if err := ServerInstanceGroupACLRemove(ctx, "12", "5", "3"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := rec.find(http.MethodDelete, "/api/v2/server-instance-groups/12/config/networking/connections/5/security/rules/3"); !ok {
		t.Error("acl remove should DELETE the rule")
	}
}

func TestServerInstanceIdValidation(t *testing.T) {
	if _, err := getServerInstanceId("abc"); err == nil {
		t.Error("a non numeric server instance id should be rejected")
	}
	if _, err := GetServerInstanceGroupId("abc"); err == nil {
		t.Error("a non numeric server instance group id should be rejected")
	}
}
