package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

const serverGetJSON = `{
	"serverId": 42,
	"siteId": 1,
	"serverTypeId": 3,
	"serverUUID": "uuid-42",
	"serialNumber": "SN-042",
	"managementAddress": "10.0.0.42",
	"vendor": "Dell",
	"model": "R740",
	"serverStatus": "active",
	"revision": 7,
	"datacenterName": "dc1",
	"bdkDebug": 0,
	"requiresReRegister": 0,
	"serverClass": "bigdata",
	"administrationState": "active",
	"serverDhcpStatus": "active",
	"supportsFcProvisioning": 0,
	"serverCreatedTimestamp": "2024-01-01T00:00:00Z",
	"powerStatus": "on",
	"powerStatusLastUpdateTimestamp": "2024-01-01T00:00:00Z",
	"links": []
}`

func TestServerGet_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/42": testutils.RawHandler(http.StatusOK, serverGetJSON),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerGet(ctx, "42", false); err != nil {
		t.Errorf("ServerGet: expected nil error, got: %v", err)
	}
}

func TestServerGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerGet(ctx, "not-a-number", false); err == nil {
		t.Error("ServerGet with invalid id: expected error, got nil")
	}
}

func TestServerGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerGet(ctx, "99", false); err == nil {
		t.Error("ServerGet with 404: expected error, got nil")
	}
}

func TestServerArchive_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerArchive(ctx, "bad"); err == nil {
		t.Error("ServerArchive with invalid id: expected error, got nil")
	}
}

func TestServerArchive_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/5": testutils.RawHandler(http.StatusOK, serverGetJSON),
		"/api/v2/servers/5/actions/archive": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			http.NotFound(w, r)
		},
	})
	defer ts.Close()

	// Override the server response to return id=5
	const s5JSON = `{
		"serverId": 5,
		"siteId": 1,
		"serverTypeId": 3,
		"serverUUID": "uuid-5",
		"serialNumber": "SN-5",
		"managementAddress": "10.0.0.5",
		"vendor": "Dell",
		"model": "R640",
		"serverStatus": "active",
		"revision": 3,
		"datacenterName": "dc1",
		"bdkDebug": 0,
		"requiresReRegister": 0,
		"serverClass": "bigdata",
		"administrationState": "active",
		"serverDhcpStatus": "active",
		"supportsFcProvisioning": 0,
		"serverCreatedTimestamp": "2024-01-01T00:00:00Z",
		"powerStatus": "on",
		"powerStatusLastUpdateTimestamp": "2024-01-01T00:00:00Z",
		"links": []
	}`
	ts2 := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/5": testutils.RawHandler(http.StatusOK, s5JSON),
		"/api/v2/servers/5/actions/archive": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts2.Close()

	ctx := setupTestContext(ts2.URL)
	if err := ServerArchive(ctx, "5"); err != nil {
		t.Errorf("ServerArchive: expected nil error, got: %v", err)
	}
}

func TestServerDelete_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerDelete(ctx, "bad"); err == nil {
		t.Error("ServerDelete with invalid id: expected error, got nil")
	}
}

func TestServerDelete_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerDelete(ctx, "99"); err == nil {
		t.Error("ServerDelete with 404 get: expected error, got nil")
	}
}

func TestServerPower_InvalidAction(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerPower(ctx, "1", "explode"); err == nil {
		t.Error("ServerPower with invalid action: expected error, got nil")
	}
}

func TestServerPower_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerPower(ctx, "bad", "on"); err == nil {
		t.Error("ServerPower with invalid id: expected error, got nil")
	}
}

func TestServerUpdate_InvalidId(t *testing.T) {
	ctx := setupTestContext("http://localhost")
	if err := ServerUpdate(ctx, "bad-id", []byte(`{}`)); err == nil {
		t.Error("ServerUpdate with invalid id: expected error, got nil")
	}
}

func TestServerUpdate_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/42": testutils.RawHandler(http.StatusOK, serverGetJSON),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerUpdate(ctx, "42", []byte(`{}`)); err != nil {
		t.Errorf("ServerUpdate: expected nil error, got: %v", err)
	}
}

func TestServerPower_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/42":                   testutils.RawHandler(http.StatusOK, serverGetJSON),
		"/api/v2/servers/42/actions/set-power": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerPower(ctx, "42", "on"); err != nil {
		t.Errorf("ServerPower on: expected nil error, got: %v", err)
	}
}

func TestServerPowerStatus_InvalidId(t *testing.T) {
	ctx := setupTestContext("http://localhost")
	if err := ServerPowerStatus(ctx, "bad"); err == nil {
		t.Error("ServerPowerStatus with invalid id: expected error, got nil")
	}
}

func TestServerPowerStatus_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/42/actions/get-power": testutils.RawHandler(http.StatusOK, `{"powerStatus":"on"}`),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerPowerStatus(ctx, "42"); err != nil {
		t.Errorf("ServerPowerStatus: expected nil error, got: %v", err)
	}
}

func TestServerGetCredentials_InvalidId(t *testing.T) {
	ctx := setupTestContext("http://localhost")
	if err := ServerGet(ctx, "bad-id", true); err == nil {
		t.Error("ServerGet with credentials and invalid id: expected error, got nil")
	}
}

func TestServerGetCredentials_Success(t *testing.T) {
	credJSON := `{"username":"admin","password":"secret"}`
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/42":             testutils.RawHandler(http.StatusOK, serverGetJSON),
		"/api/v2/servers/42/credentials": testutils.RawHandler(http.StatusOK, credJSON),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerGet(ctx, "42", true); err != nil {
		t.Errorf("ServerGet with credentials: expected nil error, got: %v", err)
	}
}

func TestServerCapabilities_InvalidId(t *testing.T) {
	ctx := setupTestContext("http://localhost")
	if err := ServerCapabilities(ctx, "bad"); err == nil {
		t.Error("ServerCapabilities with invalid id: expected error, got nil")
	}
}

func TestServerCapabilities_Success(t *testing.T) {
	capJSON := `{"firmwareUpgradeSupported":true,"firmwareUpgradeApplyOnRebootSupported":false,"virtualMediaDeviceCount":0,"vncEnabled":false}`
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/42/capabilities": testutils.RawHandler(http.StatusOK, capJSON),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerCapabilities(ctx, "42"); err != nil {
		t.Errorf("ServerCapabilities: expected nil error, got: %v", err)
	}
}

func TestServerFactoryReset_InvalidId(t *testing.T) {
	ctx := setupTestContext("http://localhost")
	if err := ServerFactoryReset(ctx, "bad"); err == nil {
		t.Error("ServerFactoryReset with invalid id: expected error, got nil")
	}
}

func TestServerFactoryReset_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/42":                       testutils.RawHandler(http.StatusOK, serverGetJSON),
		"/api/v2/servers/42/actions/factory-reset": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerFactoryReset(ctx, "42"); err != nil {
		t.Errorf("ServerFactoryReset: expected nil error, got: %v", err)
	}
}

func TestServerEnableSnmp_InvalidId(t *testing.T) {
	ctx := setupTestContext("http://localhost")
	if err := ServerEnableSnmp(ctx, "bad"); err == nil {
		t.Error("ServerEnableSnmp with invalid id: expected error, got nil")
	}
}

func TestServerEnableSnmp_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/42":                     testutils.RawHandler(http.StatusOK, serverGetJSON),
		"/api/v2/servers/42/actions/enable-snmp": testutils.RawHandler(http.StatusOK, `1`),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerEnableSnmp(ctx, "42"); err != nil {
		t.Errorf("ServerEnableSnmp: expected nil error, got: %v", err)
	}
}

func TestServerEnableSyslog_InvalidId(t *testing.T) {
	ctx := setupTestContext("http://localhost")
	if err := ServerEnableSyslog(ctx, "bad"); err == nil {
		t.Error("ServerEnableSyslog with invalid id: expected error, got nil")
	}
}

func TestServerEnableSyslog_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/42":                          testutils.RawHandler(http.StatusOK, serverGetJSON),
		"/api/v2/servers/42/actions/syslog-subscribe": testutils.RawHandler(http.StatusAccepted, `{"jobId":5,"jobGroupId":6}`),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerEnableSyslog(ctx, "42"); err != nil {
		t.Errorf("ServerEnableSyslog: expected nil error, got: %v", err)
	}
}

func TestServerUpdateIpmiCredentials_InvalidId(t *testing.T) {
	ctx := setupTestContext("http://localhost")
	if err := ServerUpdateIpmiCredentials(ctx, "bad", "user", "pass"); err == nil {
		t.Error("ServerUpdateIpmiCredentials with invalid id: expected error, got nil")
	}
}

func TestServerUpdateIpmiCredentials_Success(t *testing.T) {
	credJSON := `{"username":"admin","password":"newpass"}`
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/42":                              testutils.RawHandler(http.StatusOK, serverGetJSON),
		"/api/v2/servers/42/actions/update-ipmi-credentials": testutils.RawHandler(http.StatusOK, credJSON),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerUpdateIpmiCredentials(ctx, "42", "admin", "newpass"); err != nil {
		t.Errorf("ServerUpdateIpmiCredentials: expected nil error, got: %v", err)
	}
}

func TestServerReRegister_InvalidId(t *testing.T) {
	ctx := setupTestContext("http://localhost")
	if err := ServerReRegister(ctx, "bad"); err == nil {
		t.Error("ServerReRegister with invalid id: expected error, got nil")
	}
}

func TestServerReRegister_Success(t *testing.T) {
	reregJSON := `{"serverId":42,"revision":7,"jobInfo":{"jobId":1,"jobGroupId":2}}`
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/42":                      testutils.RawHandler(http.StatusOK, serverGetJSON),
		"/api/v2/servers/42/actions/re-register": testutils.RawHandler(http.StatusOK, reregJSON),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerReRegister(ctx, "42"); err != nil {
		t.Errorf("ServerReRegister: expected nil error, got: %v", err)
	}
}

func TestServerRegister_Success(t *testing.T) {
	var body map[string]any
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
				return
			}
			_ = json.NewDecoder(r.Body).Decode(&body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{"serverId":42,"revision":1}`)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := sdk.CreateServer{
		SiteId:            1,
		SerialNumber:      sdk.PtrString("SN-042"),
		ManagementAddress: sdk.PtrString("10.0.0.100"),
	}
	if err := ServerRegister(ctx, config, nil); err != nil {
		t.Fatalf("ServerRegister: expected nil error, got: %v", err)
	}
	if body["serialNumber"] != "SN-042" || body["managementAddress"] != "10.0.0.100" {
		t.Errorf("unexpected request body: %v", body)
	}
	if _, ok := body["serverStatus"]; ok {
		t.Errorf("serverStatus must be omitted when not set, got body: %v", body)
	}
}

func TestServerRegister_PlannedWithRackUpdate(t *testing.T) {
	var createBody, updateBody map[string]any
	var ifMatch string
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers": func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewDecoder(r.Body).Decode(&createBody)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{"serverId":42,"revision":1}`)
		},
		"/api/v2/servers/42": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPatch {
				http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
				return
			}
			ifMatch = r.Header.Get("If-Match")
			_ = json.NewDecoder(r.Body).Decode(&updateBody)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, serverGetJSON)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := sdk.CreateServer{
		SiteId:       1,
		SerialNumber: sdk.PtrString("SN-042"),
		ServerStatus: sdk.PtrString("planned"),
	}
	update := &sdk.UpdateServer{
		RackName:    sdk.PtrString("R01"),
		InventoryId: sdk.PtrString("INV-1"),
	}
	if err := ServerRegister(ctx, config, update); err != nil {
		t.Fatalf("ServerRegister: expected nil error, got: %v", err)
	}
	if createBody["serverStatus"] != "planned" {
		t.Errorf("expected serverStatus planned, got body: %v", createBody)
	}
	if _, ok := createBody["managementAddress"]; ok {
		t.Errorf("managementAddress must be omitted for a planned server, got body: %v", createBody)
	}
	if ifMatch != "1" {
		t.Errorf("expected If-Match 1 (revision from create), got %q", ifMatch)
	}
	if updateBody["rackName"] != "R01" || updateBody["inventoryId"] != "INV-1" {
		t.Errorf("unexpected update body: %v", updateBody)
	}
}

func TestServerRegister_InvalidStatus(t *testing.T) {
	ctx := setupTestContext("http://localhost")
	config := sdk.CreateServer{
		SiteId:       1,
		SerialNumber: sdk.PtrString("SN-042"),
		ServerStatus: sdk.PtrString("active"),
	}
	if err := ServerRegister(ctx, config, nil); err == nil {
		t.Error("ServerRegister with invalid status: expected error, got nil")
	}
}

func TestServerRegisterPlanned_Success(t *testing.T) {
	var ifMatch string
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/42": testutils.RawHandler(http.StatusOK, serverGetJSON),
		"/api/v2/servers/42/actions/register": func(w http.ResponseWriter, r *http.Request) {
			ifMatch = r.Header.Get("If-Match")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, `{"serverId":42,"revision":8,"jobInfo":{"jobId":1,"jobGroupId":2}}`)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerRegisterPlanned(ctx, "42"); err != nil {
		t.Fatalf("ServerRegisterPlanned: expected nil error, got: %v", err)
	}
	if ifMatch != "7" {
		t.Errorf("expected If-Match 7, got %q", ifMatch)
	}
}

func TestServerSetRegistrationReady_Success(t *testing.T) {
	var ifMatch string
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/42": testutils.RawHandler(http.StatusOK, serverGetJSON),
		"/api/v2/servers/42/actions/set-registration-ready": func(w http.ResponseWriter, r *http.Request) {
			ifMatch = r.Header.Get("If-Match")
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerSetRegistrationReady(ctx, "42"); err != nil {
		t.Fatalf("ServerSetRegistrationReady: expected nil error, got: %v", err)
	}
	if ifMatch != "7" {
		t.Errorf("expected If-Match 7, got %q", ifMatch)
	}
}

const serverReplacementCandidatesJSON = `{
	"instanceId": 10,
	"requirements": {
		"oldServerId": 42,
		"minRamGbytes": 256,
		"minCores": 32,
		"gpuModels": {},
		"diskSizesByType": {},
		"localDisksRequired": false,
		"minFibreChannelPorts": 0,
		"links": []
	},
	"candidates": [
		{"serverId": 43, "serialNumber": "SN-043", "eligible": true, "downgrade": false, "rank": 1, "skipReasons": [], "warnings": []},
		{"serverId": 44, "eligible": false, "downgrade": false, "skipReasons": ["not enough RAM", "wrong site"], "warnings": []}
	]
}`

func TestServerReplacementCandidates_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/42/replacement-candidates": testutils.RawHandler(http.StatusOK, serverReplacementCandidatesJSON),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerReplacementCandidates(ctx, "42"); err != nil {
		t.Errorf("ServerReplacementCandidates: expected nil error, got: %v", err)
	}
}

func TestServerReplace_Success(t *testing.T) {
	var body map[string]any
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/42/actions/replace": func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewDecoder(r.Body).Decode(&body)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprint(w, `{"instanceId":10,"replacedServerId":42,"replacementServerId":43,"jobGroupId":9}`)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerReplace(ctx, "42", "43", true); err != nil {
		t.Fatalf("ServerReplace: expected nil error, got: %v", err)
	}
	if body["replacementServerId"] != float64(43) || body["acceptDowngrade"] != true {
		t.Errorf("unexpected request body: %v", body)
	}
}

func TestServerReplace_InvalidReplacementId(t *testing.T) {
	ctx := setupTestContext("http://localhost")
	if err := ServerReplace(ctx, "42", "bad", false); err == nil {
		t.Error("ServerReplace with invalid replacement id: expected error, got nil")
	}
}

func TestServerCleanupReplaced_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/42/actions/cleanup": testutils.RawHandler(http.StatusAccepted, `{"jobId":3,"jobGroupId":4}`),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerCleanupReplaced(ctx, "42"); err != nil {
		t.Errorf("ServerCleanupReplaced: expected nil error, got: %v", err)
	}
}

func TestServerFirmwareComponentsList_InvalidId(t *testing.T) {
	ctx := setupTestContext("http://localhost")
	if err := ServerFirmwareComponentsList(ctx, "bad"); err == nil {
		t.Error("ServerFirmwareComponentsList with invalid id: expected error, got nil")
	}
}

func TestServerDelete_Success(t *testing.T) {
	ts2 := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/42": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, serverGetJSON)
				return
			}
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			http.NotFound(w, r)
		},
	})
	defer ts2.Close()

	ctx := setupTestContext(ts2.URL)
	if err := ServerDelete(ctx, "42"); err != nil {
		t.Errorf("ServerDelete: expected nil error, got: %v", err)
	}
}

// buildServersBody constructs a minimal JSON server list body.
func buildServersBody(items []map[string]interface{}) string {
	var sb strings.Builder
	sb.WriteString(`{"data":[`)
	for i, item := range items {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf(
			`{"serverId":%v,"siteId":1,"serverTypeId":1,"serverUUID":%q,"serialNumber":%q,"managementAddress":"10.0.0.1","vendor":"Dell","model":"R640","serverStatus":"active","revision":1,"links":{},"serverMetricsMetadata":[]}`,
			item["serverId"], item["serverUUID"], item["serialNumber"],
		))
	}
	sb.WriteString(`]}`)
	return sb.String()
}

// serverGetWithDpuJSON mirrors a real API response for a server with a DPU: the API
// omits the password from dpuInfo entries, which the SDK's typed model rejects.
var serverGetWithDpuJSON = strings.Replace(serverGetJSON, `"links": []`, `"dpuInfo": [
		{
			"macAddressOnHostOS": "00:1a:ca:ff:ff:01",
			"managementAddress": "172.18.35.73",
			"username": "ubuntu"
		}
	],
	"links": []`, 1)

func TestServerGet_DpuWithoutPassword(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/42": testutils.RawHandler(http.StatusOK, serverGetWithDpuJSON),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerGet(ctx, "42", false); err != nil {
		t.Errorf("ServerGet with DPU server: expected nil error, got: %v", err)
	}
	if err := ServerUpdate(ctx, "42", []byte(`{}`)); err != nil {
		t.Errorf("ServerUpdate with DPU server: expected nil error, got: %v", err)
	}
}

// TestServerRaw_JsonRoundTrip checks that json/yaml output keeps every field of the
// API response, and that marshalling a serverRaw without a stored body still works.
func TestServerRaw_JsonRoundTrip(t *testing.T) {
	var server serverRaw
	if err := json.Unmarshal([]byte(serverGetWithDpuJSON), &server); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	out, err := json.Marshal(server)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(out, &decoded); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}
	if _, ok := decoded["dpuInfo"]; !ok {
		t.Error("expected dpuInfo to survive the round trip")
	}

	out, err = json.Marshal(serverRaw{ServerId: 42})
	if err != nil {
		t.Fatalf("marshal without body: %v", err)
	}
	if !strings.Contains(string(out), `"serverId":42`) {
		t.Errorf("marshal without body: got %s", out)
	}
}
