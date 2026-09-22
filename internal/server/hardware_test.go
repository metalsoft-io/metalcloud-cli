package server

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

// hwRecorder captures the requests received by the mock server so tests can
// assert on method, path, headers and body.
type hwRecorder struct {
	mu   sync.Mutex
	reqs []hwRecordedRequest
}

type hwRecordedRequest struct {
	Method string
	Path   string
	Query  string
	Header http.Header
	Body   string
}

func (r *hwRecorder) wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		r.mu.Lock()
		r.reqs = append(r.reqs, hwRecordedRequest{
			Method: req.Method,
			Path:   req.URL.Path,
			Query:  req.URL.RawQuery,
			Header: req.Header.Clone(),
			Body:   string(body),
		})
		r.mu.Unlock()
		next(w, req)
	}
}

func (r *hwRecorder) last() hwRecordedRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.reqs) == 0 {
		return hwRecordedRequest{}
	}
	return r.reqs[len(r.reqs)-1]
}

func hwServer(t *testing.T, rec *hwRecorder, routes map[string]http.HandlerFunc) *httptest.Server {
	t.Helper()
	wrapped := make(map[string]http.HandlerFunc, len(routes))
	for path, handler := range routes {
		wrapped[path] = rec.wrap(handler)
	}
	return testutils.NewTestServer(wrapped)
}

const hwDriftItem = `{
	"id": 9,
	"serverId": 1,
	"snapshotId": "abc123",
	"configurationDrift": "- foo\n+ bar",
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"acknowledgedBy": 3,
	"acknowledgeTimestamp": "2024-01-02T00:00:00Z",
	"links": []
}`

const hwSnapshotItem = `{
	"oid": "5b17fbca",
	"message": "caller-cron",
	"timestamp": "2024-01-01T00:00:00Z",
	"kind": "cron"
}`

const hwServerItem = `{
	"serverId": 1,
	"siteId": 1,
	"serverTypeId": 3,
	"serverUUID": "uuid-1",
	"serialNumber": "SN-001",
	"managementAddress": "10.0.0.1",
	"vendor": "Dell",
	"model": "R740",
	"serverStatus": "available",
	"revision": 7,
	"links": []
}`

func TestServerDriftList(t *testing.T) {
	rec := &hwRecorder{}
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/1/drift-history": testutils.RawHandler(http.StatusOK,
			`{"data":[`+hwDriftItem+`],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	out := testutils.CaptureStdout(t, func() {
		if err := ServerDriftList(ctx, "1"); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	if !strings.Contains(out, "abc123") {
		t.Errorf("expected snapshot id in output, got: %s", out)
	}
	if rec.last().Path != "/api/v2/servers/1/drift-history" {
		t.Errorf("unexpected path: %s", rec.last().Path)
	}
}

func TestServerDriftListInvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://127.0.0.1:1")
	if err := ServerDriftList(ctx, "not-a-number"); err == nil {
		t.Fatal("expected error for an invalid server ID")
	}
}

func TestServerDriftGet(t *testing.T) {
	rec := &hwRecorder{}
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/1/drift-history/9": testutils.RawHandler(http.StatusOK, hwDriftItem),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	out := testutils.CaptureStdout(t, func() {
		if err := ServerDriftGet(ctx, "1", "9"); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	if !strings.Contains(out, `"id":9`) {
		t.Errorf("expected drift id in output, got: %s", out)
	}
}

func TestServerDriftAcknowledge(t *testing.T) {
	rec := &hwRecorder{}
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/1/drift-history/9/actions/acknowledge": testutils.RawHandler(http.StatusOK, hwDriftItem),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	testutils.CaptureStdout(t, func() {
		if err := ServerDriftAcknowledge(ctx, "1", "9"); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	if rec.last().Method != http.MethodPost {
		t.Errorf("acknowledge should be a POST, got %s", rec.last().Method)
	}
}

func TestServerSnapshotList(t *testing.T) {
	rec := &hwRecorder{}
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/1/snapshots": testutils.RawHandler(http.StatusOK,
			`{"data":[`+hwSnapshotItem+`],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	out := testutils.CaptureStdout(t, func() {
		if err := ServerSnapshotList(ctx, "1", "cron"); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	if !strings.Contains(out, "5b17fbca") {
		t.Errorf("expected snapshot oid in output, got: %s", out)
	}
	if !strings.Contains(rec.last().Query, "kind=cron") {
		t.Errorf("expected the kind filter to be sent, got query: %s", rec.last().Query)
	}
}

func TestServerSyncTargetSnapshot(t *testing.T) {
	rec := &hwRecorder{}
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/1/actions/sync-target-snapshot-with-latest-snapshot": testutils.NoContentHandler(),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	if err := ServerSyncTargetSnapshot(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if rec.last().Method != http.MethodPost {
		t.Errorf("sync-target-snapshot should be a POST, got %s", rec.last().Method)
	}
}

func TestServerStatistics(t *testing.T) {
	rec := &hwRecorder{}
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/statistics": testutils.RawHandler(http.StatusOK,
			`{"serverStatus":{"available":3,"used":1},"site":{"dc-1":4}}`),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	out := testutils.CaptureStdout(t, func() {
		if err := ServerStatistics(ctx); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	if !strings.Contains(out, "available") {
		t.Errorf("expected the status counts in the output, got: %s", out)
	}
}

func TestServerStatisticsHttpError(t *testing.T) {
	rec := &hwRecorder{}
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/statistics": testutils.ErrorHandler(http.StatusInternalServerError, "boom"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ServerStatistics(ctx); err == nil {
		t.Fatal("expected an error for HTTP 500")
	}
}

func TestServerHardwareRescan(t *testing.T) {
	rec := &hwRecorder{}
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/1": testutils.RawHandler(http.StatusOK, hwServerItem),
		"/api/v2/servers/1/actions/hardware-rescan": testutils.RawHandler(http.StatusOK,
			`{"serverId":1,"revision":8,"jobInfo":{"jobId":42,"jobGroupId":7}}`),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	testutils.CaptureStdout(t, func() {
		if err := ServerHardwareRescan(ctx, "1", true); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	last := rec.last()
	if last.Header.Get("If-Match") != "7" {
		t.Errorf("expected the server revision as If-Match, got %q", last.Header.Get("If-Match"))
	}
	if !strings.Contains(last.Body, `"rebootAllowed":true`) {
		t.Errorf("expected the rescan body to be sent, got %q", last.Body)
	}
}

func TestServerRegisterProduction(t *testing.T) {
	rec := &hwRecorder{}
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/actions/register-production": testutils.RawHandler(http.StatusOK,
			`{"serverId":1,"revision":1,"serverUUID":"uuid-1","serialNumber":"SN-001","jobInfo":{"jobId":42,"jobGroupId":7}}`),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	out := testutils.CaptureStdout(t, func() {
		err := ServerRegisterProduction(ctx, sdk.RegisterProductionServer{
			SiteId:   1,
			Settings: sdk.RegisterProductionServerSettings{InfrastructureId: 100},
		})
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	if !strings.Contains(out, "uuid-1") {
		t.Errorf("expected the server UUID in the output, got: %s", out)
	}
	if !strings.Contains(rec.last().Body, `"infrastructureId":100`) {
		t.Errorf("expected the settings to be sent, got %q", rec.last().Body)
	}
}

func TestServerConnectInterface(t *testing.T) {
	rec := &hwRecorder{}
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/1/actions/connect-interface": testutils.NoContentHandler(),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	err := ServerConnectInterface(ctx, "1", sdk.ServerConnectInterface{
		ServerInterfaceId:     2,
		NetworkDevicePortId:   "Ethernet0",
		NetworkDeviceHostname: "leaf-01",
	})
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	// The body must always be sent: an omitted body serializes as a literal
	// null, which the API rejects with 400.
	if !strings.Contains(rec.last().Body, `"serverInterfaceId":2`) {
		t.Errorf("expected the connect body to be sent, got %q", rec.last().Body)
	}
}

func TestServerSetInterfacesDefaultFabric(t *testing.T) {
	rec := &hwRecorder{}
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/1/actions/set-interfaces-default-fabric": testutils.NoContentHandler(),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	fabricConfig := sdk.ServerInterfacesDefaultFabric{ServerInterfaceIds: []int64{1, 2}}
	fabricConfig.DefaultFabricId.Set(sdk.PtrInt64(10))

	if err := ServerSetInterfacesDefaultFabric(ctx, "1", fabricConfig); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if !strings.Contains(rec.last().Body, `"defaultFabricId":10`) {
		t.Errorf("expected the fabric id to be sent, got %q", rec.last().Body)
	}

	// A cleared fabric must serialize as an explicit null, not be dropped.
	fabricConfig.DefaultFabricId.Set(nil)
	if err := ServerSetInterfacesDefaultFabric(ctx, "1", fabricConfig); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if !strings.Contains(rec.last().Body, `"defaultFabricId":null`) {
		t.Errorf("expected a null fabric id, got %q", rec.last().Body)
	}
}

func TestServerSetInterfacesRedundancyGroup(t *testing.T) {
	rec := &hwRecorder{}
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/1/actions/set-interfaces-redundancy-group": testutils.NoContentHandler(),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	redundancyConfig := sdk.ServerInterfacesRedundancyGroup{ServerInterfaceIds: []int64{1, 2}}
	redundancyConfig.RedundancyGroupIndex.Set(sdk.PtrFloat32(1))

	if err := ServerSetInterfacesRedundancyGroup(ctx, "1", redundancyConfig); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if !strings.Contains(rec.last().Body, `"redundancyGroupIndex":1`) {
		t.Errorf("expected the redundancy group index to be sent, got %q", rec.last().Body)
	}
}

func TestServerImportUnmanaged(t *testing.T) {
	rec := &hwRecorder{}
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/unmanaged/import": testutils.RawHandler(http.StatusOK, hwServerItem),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	config := []byte(`{
		"siteId": 1,
		"serverTypeId": 2,
		"serverInterfaces": [
			{"serverInterfaceMacAddress":"AA:BB:CC:DD:EE:00","identifierString":"eth0","networkEquipmentInterfaceIdentifierString":"Ethernet0"}
		]
	}`)

	out := testutils.CaptureStdout(t, func() {
		if err := ServerImportUnmanaged(ctx, config); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	if !strings.Contains(out, "SN-001") {
		t.Errorf("expected the imported server in the output, got: %s", out)
	}
	if !strings.Contains(rec.last().Body, `"serverTypeId":2`) {
		t.Errorf("expected the import body to be sent, got %q", rec.last().Body)
	}
}

func TestServerImportUnmanagedHttpError(t *testing.T) {
	rec := &hwRecorder{}
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/unmanaged/import": testutils.ErrorHandler(http.StatusBadRequest, "bad request"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	err := ServerImportUnmanaged(ctx, []byte(`{"siteId":1,"serverTypeId":2,"serverInterfaces":[]}`))
	if err == nil {
		t.Fatal("expected an error for HTTP 400")
	}
}

func TestServerFirmwareUpgradeBatch(t *testing.T) {
	rec := &hwRecorder{}
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/firmware/actions/batch-upgrade": testutils.RawHandler(http.StatusOK,
			`{"successful":{"1":{"jobId":42,"jobGroupId":7}},"failed":{"2":"no baseline"}}`),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	out := testutils.CaptureStdout(t, func() {
		if err := ServerFirmwareUpgradeBatch(ctx, []string{"1", "2"}); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	if !strings.Contains(out, "no baseline") {
		t.Errorf("expected the failed entry in the output, got: %s", out)
	}
	if !strings.Contains(rec.last().Body, `"serverIds":[1,2]`) {
		t.Errorf("expected the server ids to be sent, got %q", rec.last().Body)
	}
}

func TestServerFirmwareScheduleUpgradeBatch(t *testing.T) {
	rec := &hwRecorder{}
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/1/firmware/actions/batch-schedule-upgrade": testutils.RawHandler(http.StatusOK,
			`{"failed":{}}`),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	testutils.CaptureStdout(t, func() {
		err := ServerFirmwareScheduleUpgradeBatch(ctx, []string{"1", "2"}, "2024-01-01T10:00:00Z", true)
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	last := rec.last()
	if last.Path != "/api/v2/servers/1/firmware/actions/batch-schedule-upgrade" {
		t.Errorf("the first server id should address the endpoint, got %s", last.Path)
	}
	if !strings.Contains(last.Body, `"serverIds":[1,2]`) ||
		!strings.Contains(last.Body, `"scheduleUpdateTimestamp":"2024-01-01T10:00:00Z"`) ||
		!strings.Contains(last.Body, `"confirmationRequired":true`) {
		t.Errorf("unexpected schedule body: %q", last.Body)
	}
}

func TestServerFirmwareScheduleUpgradeBatchNoServers(t *testing.T) {
	ctx := testutils.SetupTestContext("http://127.0.0.1:1")
	if err := ServerFirmwareScheduleUpgradeBatch(ctx, nil, "", false); err == nil {
		t.Fatal("expected an error when no server ID is given")
	}
}

func TestServerConfigExample(t *testing.T) {
	ctx := testutils.SetupTestContext("http://127.0.0.1:1")

	for _, kind := range ServerConfigExampleKinds {
		out := testutils.CaptureStdout(t, func() {
			if err := ServerConfigExample(ctx, kind); err != nil {
				t.Fatalf("config example %s: %v", kind, err)
			}
		})
		if !strings.Contains(out, "{") {
			t.Errorf("config example %s produced no output: %s", kind, out)
		}
	}

	if err := ServerConfigExample(ctx, "no-such-kind"); err == nil {
		t.Fatal("expected an error for an unknown configuration kind")
	}
}
