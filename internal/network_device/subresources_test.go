package network_device

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

// --- fixtures ---------------------------------------------------------------

// ndIfaceItem satisfies every required property of sdk.NetworkEquipmentInterface.
var ndIfaceItem = map[string]any{
	"interfaceId": 42, "switchId": 1, "kind": "loopback", "interfaceName": "Loopback1",
	"interfaceDescription": "mgmt", "mtu": 9216, "enabled": true,
	"serviceStatus": "active", "pendingDelete": false,
	"config": map[string]any{"revision": 3},
	"ipv4":   map[string]any{"addresses": []any{}},
	"ipv6":   map[string]any{"addresses": []any{}},
}

// ndPortConfigItem satisfies sdk.NetworkEquipmentInterfaceConfig.
var ndPortConfigItem = map[string]any{"revision": 3, "description": "mgmt", "mtu": 9216, "enabled": true}

// ndPortIpItem satisfies sdk.NetworkEquipmentInterfaceIp.
var ndPortIpItem = map[string]any{
	"id": 7, "interfaceId": 42, "kind": "ipv4", "address": "10.0.0.1",
	"prefixLength": 32, "serviceStatus": "active", "pendingDelete": false,
}

// ndBreakoutItem satisfies sdk.NetworkEquipmentBreakout.
var ndBreakoutItem = map[string]any{
	"id": 5, "portName": "Ethernet0", "serviceStatus": "active", "pendingDelete": false,
	"revision":       2,
	"breakoutGroups": []any{map[string]any{"numberOfInterfaces": 4, "speed": "25G"}},
	"config": map[string]any{
		"revision":       9,
		"breakoutGroups": []any{map[string]any{"numberOfInterfaces": 2, "speed": "50G"}},
	},
}

// ndBreakoutConfigItem satisfies sdk.NetworkEquipmentBreakoutConfigDto.
var ndBreakoutConfigItem = map[string]any{
	"revision":       9,
	"breakoutGroups": []any{map[string]any{"numberOfInterfaces": 2, "speed": "50G"}},
}

// ndVirtualFunctionItem satisfies sdk.NetworkDeviceInterfaceVirtualFunction.
var ndVirtualFunctionItem = map[string]any{
	"id": 3, "networkDeviceId": 1, "interfaceId": 42, "name": "vf0", "index": 0, "status": "active",
}

// ndDriftItem satisfies sdk.NetworkDeviceDriftHistory.
var ndDriftItem = map[string]any{
	"id": 9, "networkDeviceId": 1, "snapshotId": "abc123",
	"configurationDrift": "- hostname a\n+ hostname b",
	"createdTimestamp":   "2024-01-01T00:00:00Z",
}

// ndSnapshotItem satisfies sdk.NetworkDeviceSnapshot.
var ndSnapshotItem = map[string]any{"oid": "abc123", "message": "backup", "timestamp": "2024-01-01T00:00:00Z", "kind": "backup"}

// ndVendorItem satisfies sdk.NetworkDeviceVendors.
var ndVendorItem = map[string]any{
	"id": 3, "revision": 4, "kind": "sonic_enterprise",
	"oidGroups": []any{map[string]any{
		"name": "interfaces",
		"oids": map[string]any{"1.3.6.1": "ifInOctets"},
		"mapping": map[string]any{
			"metric": "interface_in_octets",
			"type":   "integer",
		},
	}},
	"healthCheckRules":      map[string]any{"rules": []any{}},
	"optionalFilesToBackup": map[string]any{"files": map[string]any{"default": []any{"/etc/config"}}},
}

// ndHealthSummaryItem satisfies sdk.NetworkDeviceHealthSummaryDto.
var ndHealthSummaryItem = map[string]any{
	"lastUpdated": "2024-01-01T00:00:00Z", "overallSeverity": "warning", "trendDirection": "stable",
	"suspectedRootCauses": []any{"optics"}, "keyFindings": "one flapping port",
	"notableEvents": []any{}, "detectedIssues": []any{},
}

// ndLivePortsItem satisfies sdk.NetworkDevicePorts.
var ndLivePortsItem = map[string]any{
	"switch_id": "1",
	"ports": []any{map[string]any{
		"port_name": "Ethernet0", "enabled": true, "active": true, "link_speed": 25000,
		"link_duplex": "full", "utilization_in": 0.1, "utilization_out": 0.2,
		"counters": map[string]any{
			"octets_in": 1, "octets_out": 2, "packets_in": 3,
			"packets_out": 4, "errors_in": 0, "errors_out": 0,
		},
	}},
}

var ndJobInfo = map[string]any{"jobId": 11, "jobGroupId": 5}

var ndDeviceItem = testutils.NetworkDeviceFixture("1", 1, "sw-01")

// --- recorder ---------------------------------------------------------------

// ndRecorder captures the requests the mock server received so tests can assert
// on method, path, headers and body.
type ndRecorder struct {
	mu   sync.Mutex
	reqs []ndRecordedRequest
}

type ndRecordedRequest struct {
	Method string
	Path   string
	Query  string
	Header http.Header
	Body   string
}

func (r *ndRecorder) wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		r.mu.Lock()
		r.reqs = append(r.reqs, ndRecordedRequest{
			Method: req.Method, Path: req.URL.Path, Query: req.URL.RawQuery,
			Header: req.Header.Clone(), Body: string(body),
		})
		r.mu.Unlock()
		next(w, req)
	}
}

func (r *ndRecorder) last() ndRecordedRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.reqs[len(r.reqs)-1]
}

// ndNewServer builds a mock API server covering every network device
// sub-resource endpoint exercised by these tests.
func ndNewServer(t *testing.T, rec *ndRecorder) *httptest.Server {
	t.Helper()

	routes := map[string]http.HandlerFunc{
		"/api/v2/network-devices":   testutils.JSONHandler(200, testutils.PaginatedResponse([]any{ndDeviceItem}, 1, 1)),
		"/api/v2/network-devices/1": testutils.JSONHandler(200, ndDeviceItem),

		// ports
		"/api/v2/network-devices/1/ports": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				testutils.JSONHandler(201, ndIfaceItem)(w, r)
				return
			}
			testutils.JSONHandler(200, testutils.PaginatedResponse([]any{ndIfaceItem}, 1, 1))(w, r)
		},
		"/api/v2/network-devices/1/ports/42": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				testutils.NoContentHandler()(w, r)
				return
			}
			testutils.JSONHandler(200, ndIfaceItem)(w, r)
		},
		"/api/v2/network-devices/1/ports/42/config":              testutils.JSONHandler(200, ndPortConfigItem),
		"/api/v2/network-devices/1/ports/42/virtual-functions":   testutils.JSONHandler(200, testutils.PaginatedResponse([]any{ndVirtualFunctionItem}, 1, 1)),
		"/api/v2/network-devices/1/ports/42/virtual-functions/3": testutils.JSONHandler(200, ndVirtualFunctionItem),
		"/api/v2/network-devices/1/ports/42/ipv4/addresses": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPut {
				testutils.JSONHandler(200, []any{ndPortIpItem})(w, r)
				return
			}
			testutils.JSONHandler(200, []any{ndPortIpItem})(w, r)
		},
		"/api/v2/network-devices/1/ports/42/ipv4/addresses/7": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				testutils.NoContentHandler()(w, r)
				return
			}
			testutils.JSONHandler(200, ndPortIpItem)(w, r)
		},
		"/api/v2/network-devices/1/actions/ports": testutils.JSONHandler(200, ndLivePortsItem),

		// breakouts
		"/api/v2/network-devices/1/breakouts": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				testutils.JSONHandler(201, ndBreakoutItem)(w, r)
				return
			}
			testutils.JSONHandler(200, []any{ndBreakoutItem})(w, r)
		},
		"/api/v2/network-devices/1/breakouts/5": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				testutils.NoContentHandler()(w, r)
				return
			}
			testutils.JSONHandler(200, ndBreakoutItem)(w, r)
		},
		"/api/v2/network-devices/1/breakouts/5/config": testutils.JSONHandler(200, ndBreakoutConfigItem),

		// virtual functions
		"/api/v2/network-devices/1/virtual-functions":   testutils.JSONHandler(200, testutils.PaginatedResponse([]any{ndVirtualFunctionItem}, 1, 1)),
		"/api/v2/network-devices/1/virtual-functions/3": testutils.JSONHandler(200, ndVirtualFunctionItem),

		// secrets
		"/api/v2/network-devices/1/secrets": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				testutils.NoContentHandler()(w, r)
				return
			}
			testutils.JSONHandler(200, map[string]any{"names": []any{"enable_password"}})(w, r)
		},
		"/api/v2/network-devices/1/secrets/enable_password":             testutils.NoContentHandler(),
		"/api/v2/network-devices/1/secrets/enable_password/credentials": testutils.JSONHandler(200, map[string]any{"secretValue": "s3cr3t"}),

		// monitoring
		"/api/v2/network-devices/1/actions/snmp-monitoring-subscribe":       testutils.NoContentHandler(),
		"/api/v2/network-devices/1/actions/snmp-monitoring-unsubscribe":     testutils.NoContentHandler(),
		"/api/v2/network-devices/actions/snmp-monitoring-subscribe/batch":   testutils.NoContentHandler(),
		"/api/v2/network-devices/actions/snmp-monitoring-unsubscribe/batch": testutils.NoContentHandler(),
		"/api/v2/network-devices/1/actions/snmp-service-enable":             testutils.JSONHandler(200, ndJobInfo),
		"/api/v2/network-devices/1/actions/snmp-service-disable":            testutils.JSONHandler(200, ndJobInfo),
		"/api/v2/network-devices/1/actions/syslog-unsubscribe":              testutils.JSONHandler(200, ndJobInfo),
		"/api/v2/network-devices/snmp-monitoring/agent-info/batch":          testutils.JSONHandler(200, map[string]any{"allocationInfo": map[string]any{"1": "agent-a"}}),
		"/api/v2/network-devices/1/health-summary":                          testutils.JSONHandler(200, ndHealthSummaryItem),
		"/api/v2/network-devices/health-monitoring-filter":                  testutils.NoContentHandler(),
		"/api/v2/network-devices/statistics":                                testutils.JSONHandler(200, map[string]any{"deviceCount": 2, "portCount": 64}),

		// drift & snapshots
		"/api/v2/network-devices/1/drift-history":                                     testutils.JSONHandler(200, testutils.PaginatedResponse([]any{ndDriftItem}, 1, 1)),
		"/api/v2/network-devices/1/drift-history/9":                                   testutils.JSONHandler(200, ndDriftItem),
		"/api/v2/network-devices/1/drift-history/9/actions/acknowledge":               testutils.JSONHandler(200, ndDriftItem),
		"/api/v2/network-devices/1/snapshots":                                         testutils.JSONHandler(200, testutils.PaginatedResponse([]any{ndSnapshotItem}, 1, 1)),
		"/api/v2/network-devices/1/actions/sync-target-snapshot-with-latest-snapshot": testutils.NoContentHandler(),

		// lifecycle
		"/api/v2/network-devices/1/replace":                         testutils.JSONHandler(200, map[string]any{"status": "ok"}),
		"/api/v2/network-devices/1/re-provision":                    testutils.JSONHandler(200, ndJobInfo),
		"/api/v2/network-devices/1/actions/return-to-planned":       testutils.JSONHandler(200, ndDeviceItem),
		"/api/v2/network-devices/1/actions/revert-failed-state":     testutils.JSONHandler(200, ndDeviceItem),
		"/api/v2/network-devices/1/actions/start-registration":      testutils.JSONHandler(200, ndDeviceItem),
		"/api/v2/network-devices/1/actions/mark-installation-ready": testutils.JSONHandler(200, ndDeviceItem),
		"/api/v2/network-devices/1/actions/run-extension":           testutils.JSONHandler(200, ndJobInfo),

		// vendors & drivers
		"/api/v2/network-devices/vendors":              testutils.JSONHandler(200, testutils.PaginatedResponse([]any{ndVendorItem}, 1, 1)),
		"/api/v2/network-devices/vendors/3":            testutils.JSONHandler(200, ndVendorItem),
		"/api/v2/network-devices/drivers/capabilities": testutils.JSONHandler(200, []any{map[string]any{"driver": "sonic_enterprise", "activationActions": []any{"start-registration"}}}),
	}

	if rec != nil {
		for path, handler := range routes {
			routes[path] = rec.wrap(handler)
		}
	}

	return testutils.NewTestServer(routes)
}

// --- ports ------------------------------------------------------------------

func TestNetworkDevicePortCommands(t *testing.T) {
	rec := &ndRecorder{}
	ts := ndNewServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	if err := NetworkDevicePortCreate(ctx, "1", []byte(`{"kind":"loopback","name":"Loopback1"}`)); err != nil {
		t.Fatalf("NetworkDevicePortCreate: %v", err)
	}
	if got := rec.last(); got.Method != http.MethodPost || !strings.Contains(got.Body, "Loopback1") {
		t.Errorf("unexpected create request: %+v", got)
	}

	if err := NetworkDevicePortGetConfig(ctx, "1", "42"); err != nil {
		t.Errorf("NetworkDevicePortGetConfig: %v", err)
	}

	if err := NetworkDevicePortsLiveStatus(ctx, "1"); err != nil {
		t.Errorf("NetworkDevicePortsLiveStatus: %v", err)
	}
	if got := rec.last(); got.Method != http.MethodPost || got.Path != "/api/v2/network-devices/1/actions/ports" {
		t.Errorf("unexpected live-status request: %+v", got)
	}

	// The delete reads the port first to derive the If-Match revision
	// (config.revision + 1).
	if err := NetworkDevicePortDelete(ctx, "1", "42"); err != nil {
		t.Fatalf("NetworkDevicePortDelete: %v", err)
	}
	if got := rec.last(); got.Method != http.MethodDelete || got.Header.Get("If-Match") != "4" {
		t.Errorf("expected DELETE with If-Match 4, got: %+v", got)
	}
}

func TestNetworkDevicePortCreate_InvalidConfig(t *testing.T) {
	ts := ndNewServer(t, nil)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDevicePortCreate(ctx, "1", []byte(`not json`)); err == nil {
		t.Error("expected error for invalid config, got nil")
	}
}

func TestNetworkDevicePortIpCommands(t *testing.T) {
	rec := &ndRecorder{}
	ts := ndNewServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	if err := NetworkDevicePortIpList(ctx, "1", "42", "ipv4"); err != nil {
		t.Errorf("NetworkDevicePortIpList: %v", err)
	}
	if err := NetworkDevicePortIpGet(ctx, "1", "42", "ipv4", "7"); err != nil {
		t.Errorf("NetworkDevicePortIpGet: %v", err)
	}

	if err := NetworkDevicePortIpReplace(ctx, "1", "42", "ipv4", []byte(`{"ips":[{"address":"10.0.0.1","prefixLength":32}]}`)); err != nil {
		t.Fatalf("NetworkDevicePortIpReplace: %v", err)
	}
	if got := rec.last(); got.Method != http.MethodPut || got.Header.Get("If-Match") != "4" {
		t.Errorf("expected PUT with If-Match 4, got: %+v", got)
	}

	if err := NetworkDevicePortIpRemove(ctx, "1", "42", "ipv4", "7"); err != nil {
		t.Fatalf("NetworkDevicePortIpRemove: %v", err)
	}
	if got := rec.last(); got.Method != http.MethodDelete {
		t.Errorf("expected DELETE, got: %+v", got)
	}
}

func TestNetworkDevicePortIp_InvalidFamily(t *testing.T) {
	ts := ndNewServer(t, nil)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDevicePortIpList(ctx, "1", "42", "ipv5"); err == nil {
		t.Error("expected error for invalid address family, got nil")
	}
}

func TestNetworkDeviceVirtualFunctionCommands(t *testing.T) {
	ts := ndNewServer(t, nil)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	if err := NetworkDeviceVirtualFunctionList(ctx, "1"); err != nil {
		t.Errorf("NetworkDeviceVirtualFunctionList: %v", err)
	}
	if err := NetworkDeviceVirtualFunctionGet(ctx, "1", "3"); err != nil {
		t.Errorf("NetworkDeviceVirtualFunctionGet: %v", err)
	}
	if err := NetworkDevicePortVirtualFunctionList(ctx, "1", "42"); err != nil {
		t.Errorf("NetworkDevicePortVirtualFunctionList: %v", err)
	}
	if err := NetworkDevicePortVirtualFunctionGet(ctx, "1", "42", "3"); err != nil {
		t.Errorf("NetworkDevicePortVirtualFunctionGet: %v", err)
	}
}

// --- breakouts --------------------------------------------------------------

func TestNetworkDeviceBreakoutCommands(t *testing.T) {
	rec := &ndRecorder{}
	ts := ndNewServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	if err := NetworkDeviceBreakoutList(ctx, "1"); err != nil {
		t.Errorf("NetworkDeviceBreakoutList: %v", err)
	}
	if err := NetworkDeviceBreakoutGet(ctx, "1", "5"); err != nil {
		t.Errorf("NetworkDeviceBreakoutGet: %v", err)
	}
	if err := NetworkDeviceBreakoutGetConfig(ctx, "1", "5"); err != nil {
		t.Errorf("NetworkDeviceBreakoutGetConfig: %v", err)
	}

	if err := NetworkDeviceBreakoutCreate(ctx, "1", []byte(`{"portName":"Ethernet0","groups":[{"numberOfInterfaces":4,"speed":"25G"}]}`)); err != nil {
		t.Fatalf("NetworkDeviceBreakoutCreate: %v", err)
	}
	if got := rec.last(); got.Method != http.MethodPost || !strings.Contains(got.Body, "Ethernet0") {
		t.Errorf("unexpected breakout create request: %+v", got)
	}

	// The config sub-resource carries its own revision (9), not the breakout's (2).
	if err := NetworkDeviceBreakoutUpdateConfig(ctx, "1", "5", []byte(`{"breakoutGroups":[{"numberOfInterfaces":2,"speed":"50G"}]}`)); err != nil {
		t.Fatalf("NetworkDeviceBreakoutUpdateConfig: %v", err)
	}
	if got := rec.last(); got.Method != http.MethodPatch || got.Header.Get("If-Match") != "9" {
		t.Errorf("expected PATCH with If-Match 9, got: %+v", got)
	}

	if err := NetworkDeviceBreakoutDelete(ctx, "1", "5"); err != nil {
		t.Fatalf("NetworkDeviceBreakoutDelete: %v", err)
	}
	if got := rec.last(); got.Method != http.MethodDelete {
		t.Errorf("expected DELETE, got: %+v", got)
	}
}

func TestFormatBreakoutGroups(t *testing.T) {
	breakout := sdk.NetworkEquipmentBreakout{
		Id:       5,
		PortName: "Ethernet0",
		BreakoutGroups: []sdk.BreakoutGroup{
			{NumberOfInterfaces: 4, Speed: *sdk.NewNullableString(sdk.PtrString("25G"))},
		},
		Config: sdk.NetworkEquipmentBreakoutConfigDto{
			Revision: 9,
			// A group without a forced speed splits the port by count only.
			BreakoutGroups: []sdk.BreakoutGroup{{NumberOfInterfaces: 2}},
		},
	}

	row := toBreakoutRow(breakout)
	if row.AppliedGroups != "4x25G" {
		t.Errorf("expected applied groups '4x25G', got %q", row.AppliedGroups)
	}
	if row.StagedGroups != "2x" {
		t.Errorf("expected staged groups '2x', got %q", row.StagedGroups)
	}
}

func TestNetworkDeviceBreakoutInvalidId(t *testing.T) {
	ts := ndNewServer(t, nil)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceBreakoutGet(ctx, "1", "not-a-number"); err == nil {
		t.Error("expected error for non-numeric breakout id, got nil")
	}
}

// --- secrets ----------------------------------------------------------------

func TestNetworkDeviceSecretCommands(t *testing.T) {
	rec := &ndRecorder{}
	ts := ndNewServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	if err := NetworkDeviceSecretList(ctx, "1"); err != nil {
		t.Errorf("NetworkDeviceSecretList: %v", err)
	}

	if err := NetworkDeviceSecretSet(ctx, "1", "enable_password", "s3cr3t"); err != nil {
		t.Fatalf("NetworkDeviceSecretSet: %v", err)
	}
	if got := rec.last(); got.Method != http.MethodPatch || !strings.Contains(got.Body, "s3cr3t") {
		t.Errorf("unexpected secret set request: %+v", got)
	}

	if err := NetworkDeviceSecretGetCredentials(ctx, "1", "enable_password"); err != nil {
		t.Errorf("NetworkDeviceSecretGetCredentials: %v", err)
	}
	if err := NetworkDeviceSecretRemove(ctx, "1", "enable_password"); err != nil {
		t.Errorf("NetworkDeviceSecretRemove: %v", err)
	}
	if err := NetworkDeviceSecretRemoveAll(ctx, "1"); err != nil {
		t.Errorf("NetworkDeviceSecretRemoveAll: %v", err)
	}
}

// --- monitoring -------------------------------------------------------------

func TestNetworkDeviceMonitoringCommands(t *testing.T) {
	rec := &ndRecorder{}
	ts := ndNewServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	if err := NetworkDeviceSnmpMonitoringEnable(ctx, "1"); err != nil {
		t.Errorf("NetworkDeviceSnmpMonitoringEnable: %v", err)
	}
	if err := NetworkDeviceSnmpMonitoringDisable(ctx, "1"); err != nil {
		t.Errorf("NetworkDeviceSnmpMonitoringDisable: %v", err)
	}

	// The batch body must always be sent; an omitted body would serialize as
	// a literal null and be rejected by the API.
	if err := NetworkDeviceSnmpMonitoringEnableBatch(ctx, []string{"1", "2"}, nil, nil); err != nil {
		t.Fatalf("NetworkDeviceSnmpMonitoringEnableBatch: %v", err)
	}
	if got := rec.last(); !strings.Contains(got.Body, `"networkDeviceIds":[1,2]`) {
		t.Errorf("expected batch body with device ids, got: %q", got.Body)
	}

	if err := NetworkDeviceSnmpMonitoringDisableBatch(ctx, nil, []string{"7"}, nil); err != nil {
		t.Fatalf("NetworkDeviceSnmpMonitoringDisableBatch: %v", err)
	}
	if got := rec.last(); !strings.Contains(got.Body, `"siteIds":[7]`) {
		t.Errorf("expected batch body with site ids, got: %q", got.Body)
	}

	if err := NetworkDeviceSnmpServiceEnable(ctx, "1", nil, nil, nil); err != nil {
		t.Fatalf("NetworkDeviceSnmpServiceEnable: %v", err)
	}
	if got := rec.last(); strings.TrimSpace(got.Body) == "null" {
		t.Error("expected an object body on snmp-service-enable, got null")
	}

	if err := NetworkDeviceSnmpServiceDisable(ctx, "1"); err != nil {
		t.Errorf("NetworkDeviceSnmpServiceDisable: %v", err)
	}
	if err := NetworkDeviceDisableSyslog(ctx, "1"); err != nil {
		t.Errorf("NetworkDeviceDisableSyslog: %v", err)
	}
	if err := NetworkDeviceSnmpAgentInfo(ctx, []string{"1"}, nil, nil); err != nil {
		t.Errorf("NetworkDeviceSnmpAgentInfo: %v", err)
	}
	if err := NetworkDeviceHealthSummary(ctx, "1"); err != nil {
		t.Errorf("NetworkDeviceHealthSummary: %v", err)
	}
	if err := NetworkDeviceStatisticsGet(ctx); err != nil {
		t.Errorf("NetworkDeviceStatisticsGet: %v", err)
	}

	if err := NetworkDeviceSetHealthMonitoringFilter(ctx, "socket-1", nil, []string{"1"}, nil, nil); err != nil {
		t.Fatalf("NetworkDeviceSetHealthMonitoringFilter: %v", err)
	}
	if got := rec.last(); !strings.Contains(got.Query, "socketId=socket-1") {
		t.Errorf("expected socketId in query, got: %q", got.Query)
	}
}

func TestNetworkDeviceSnmpMonitoringBatch_NoSelection(t *testing.T) {
	ts := ndNewServer(t, nil)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceSnmpMonitoringEnableBatch(ctx, nil, nil, nil); err == nil {
		t.Error("expected error when no devices, sites or fabrics are selected")
	}
}

// --- drift & snapshots ------------------------------------------------------

func TestNetworkDeviceDriftCommands(t *testing.T) {
	rec := &ndRecorder{}
	ts := ndNewServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	if err := NetworkDeviceDriftList(ctx, "1"); err != nil {
		t.Errorf("NetworkDeviceDriftList: %v", err)
	}
	if err := NetworkDeviceDriftGet(ctx, "1", "9"); err != nil {
		t.Errorf("NetworkDeviceDriftGet: %v", err)
	}
	if err := NetworkDeviceDriftAcknowledge(ctx, "1", "9"); err != nil {
		t.Fatalf("NetworkDeviceDriftAcknowledge: %v", err)
	}
	if got := rec.last(); got.Method != http.MethodPost {
		t.Errorf("expected POST on acknowledge, got: %+v", got)
	}

	if err := NetworkDeviceSnapshotList(ctx, "1", "backup"); err != nil {
		t.Fatalf("NetworkDeviceSnapshotList: %v", err)
	}
	if got := rec.last(); !strings.Contains(got.Query, "kind=backup") {
		t.Errorf("expected kind filter in query, got: %q", got.Query)
	}

	if err := NetworkDeviceSyncTargetSnapshot(ctx, "1"); err != nil {
		t.Errorf("NetworkDeviceSyncTargetSnapshot: %v", err)
	}
}

// --- lifecycle --------------------------------------------------------------

func TestNetworkDeviceLifecycleCommands(t *testing.T) {
	rec := &ndRecorder{}
	ts := ndNewServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	if err := NetworkDeviceReplace(ctx, "1", "1"); err != nil {
		t.Fatalf("NetworkDeviceReplace: %v", err)
	}
	if got := rec.last(); !strings.Contains(got.Body, `"newSwitchId":1`) {
		t.Errorf("expected replacement id in body, got: %q", got.Body)
	}

	if err := NetworkDeviceReProvision(ctx, "1", "full"); err != nil {
		t.Fatalf("NetworkDeviceReProvision: %v", err)
	}
	if got := rec.last(); !strings.Contains(got.Body, `"reprovisionType":"full"`) {
		t.Errorf("expected reprovision type in body, got: %q", got.Body)
	}

	// The optional body is always sent so the SDK does not post a literal null.
	if err := NetworkDeviceReturnToPlanned(ctx, "1", nil); err != nil {
		t.Fatalf("NetworkDeviceReturnToPlanned: %v", err)
	}
	if got := rec.last(); strings.TrimSpace(got.Body) == "null" || got.Header.Get("If-Match") != "1" {
		t.Errorf("expected object body and If-Match 1, got: %+v", got)
	}

	if err := NetworkDeviceRevertFailedState(ctx, "1"); err != nil {
		t.Errorf("NetworkDeviceRevertFailedState: %v", err)
	}
	if err := NetworkDeviceStartRegistration(ctx, "1"); err != nil {
		t.Errorf("NetworkDeviceStartRegistration: %v", err)
	}
	if err := NetworkDeviceMarkInstallationReady(ctx, "1"); err != nil {
		t.Errorf("NetworkDeviceMarkInstallationReady: %v", err)
	}

	if err := NetworkDeviceRunExtension(ctx, "1", []byte(`{"extensionId":7,"inputArguments":null}`)); err != nil {
		t.Fatalf("NetworkDeviceRunExtension: %v", err)
	}
	if got := rec.last(); !strings.Contains(got.Body, `"inputArguments":{}`) {
		t.Errorf("expected defaulted inputArguments in body, got: %q", got.Body)
	}
}

// --- vendors ----------------------------------------------------------------

func TestNetworkDeviceVendorCommands(t *testing.T) {
	rec := &ndRecorder{}
	ts := ndNewServer(t, rec)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	if err := NetworkDeviceVendorList(ctx, []string{"sonic_enterprise"}); err != nil {
		t.Errorf("NetworkDeviceVendorList: %v", err)
	}
	if err := NetworkDeviceVendorGet(ctx, "3"); err != nil {
		t.Errorf("NetworkDeviceVendorGet: %v", err)
	}

	if err := NetworkDeviceVendorUpdate(ctx, "3", []byte(`{"oidGroups":[]}`)); err != nil {
		t.Fatalf("NetworkDeviceVendorUpdate: %v", err)
	}
	if got := rec.last(); got.Method != http.MethodPatch || got.Header.Get("If-Match") != "4" {
		t.Errorf("expected PATCH with If-Match 4, got: %+v", got)
	}

	if err := NetworkDeviceDriverCapabilities(ctx, "sonic"); err != nil {
		t.Errorf("NetworkDeviceDriverCapabilities: %v", err)
	}
}

// --- config examples --------------------------------------------------------

func TestNetworkDeviceConfigExamples(t *testing.T) {
	ctx := testutils.SetupTestContext("")

	examples := map[string]func() error{
		"port":              func() error { return NetworkDevicePortConfigExample(ctx) },
		"port-ip":           func() error { return NetworkDevicePortIpReplaceConfigExample(ctx) },
		"breakout":          func() error { return NetworkDeviceBreakoutConfigExample(ctx) },
		"breakout-config":   func() error { return NetworkDeviceBreakoutUpdateConfigExample(ctx) },
		"vendor":            func() error { return NetworkDeviceVendorConfigExample(ctx) },
		"return-to-planned": func() error { return NetworkDeviceReturnToPlannedConfigExample(ctx) },
		"run-extension":     func() error { return NetworkDeviceRunExtensionConfigExample(ctx) },
	}

	for name, example := range examples {
		if err := example(); err != nil {
			t.Errorf("%s config example: %v", name, err)
		}
	}
}
