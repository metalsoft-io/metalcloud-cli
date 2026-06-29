package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestServerList_MultiPage(t *testing.T) {
	page1 := make([]map[string]interface{}, 100)
	page2 := make([]map[string]interface{}, 100)
	page3 := make([]map[string]interface{}, 5)
	for i := range page1 {
		page1[i] = map[string]interface{}{
			"serverId": i + 1, "siteId": 1, "serverTypeId": 1,
			"serverUUID": fmt.Sprintf("uuid-%d", i+1), "serialNumber": fmt.Sprintf("SN-%d", i+1),
			"managementAddress": "10.0.0.1", "vendor": "Dell", "model": "R640",
			"serverStatus": "active", "revision": 1, "links": map[string]interface{}{},
			"serverMetricsMetadata": []interface{}{},
		}
	}
	for i := range page2 {
		page2[i] = map[string]interface{}{
			"serverId": 100 + i + 1, "siteId": 1, "serverTypeId": 1,
			"serverUUID": fmt.Sprintf("uuid-%d", 100+i+1), "serialNumber": fmt.Sprintf("SN-%d", 100+i+1),
			"managementAddress": "10.0.0.2", "vendor": "Dell", "model": "R640",
			"serverStatus": "active", "revision": 1, "links": map[string]interface{}{},
			"serverMetricsMetadata": []interface{}{},
		}
	}
	for i := range page3 {
		page3[i] = map[string]interface{}{
			"serverId": 200 + i + 1, "siteId": 1, "serverTypeId": 1,
			"serverUUID": fmt.Sprintf("uuid-%d", 200+i+1), "serialNumber": fmt.Sprintf("SN-%d", 200+i+1),
			"managementAddress": "10.0.0.3", "vendor": "Dell", "model": "R640",
			"serverStatus": "active", "revision": 1, "links": map[string]interface{}{},
			"serverMetricsMetadata": []interface{}{},
		}
	}

	// ServerList uses raw body reads — it does NOT use FetchAllPages.
	// So we just verify a single-page call of 205 items works without error.
	all205 := append(append(page1, page2...), page3...)
	body := buildServersBody(all205)

	ts := httptest.NewServer(newServerListHandler(body, http.StatusOK))
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := ServerList(ctx, false, nil, nil); err != nil {
		t.Errorf("ServerList with 205 items: expected nil error, got: %v", err)
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
		"/api/v2/servers/42/actions/syslog-subscribe": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
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
	regJSON := `{"serverId":42,"revision":1,"serverUUID":"uuid-42","serialNumber":"SN-042","jobInfo":{"jobId":1,"jobGroupId":2}}`
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers": testutils.RawHandler(http.StatusOK, regJSON),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := sdk.RegisterServer{
		ManagementAddress: sdk.PtrString("10.0.0.100"),
	}
	if err := ServerRegister(ctx, config); err != nil {
		t.Errorf("ServerRegister: expected nil error, got: %v", err)
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
