package server_instance

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestServerInstanceGet_Success(t *testing.T) {
	const siJSON = `{"id":42,"revision":1,"label":"si-42","infrastructureId":123,"groupId":10,"serviceStatus":"active","isVmInstance":0,"isEndpointInstance":0,"createdTimestamp":"2024-01-01T00:00:00Z","updatedTimestamp":"2024-01-01T00:00:00Z","meta":{},"links":[]}`
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/server-instances/42": testutils.RawHandler(http.StatusOK, siJSON),
	})
	defer ts.Close()
	ctx := setupTestContext(ts.URL)
	if err := ServerInstanceGet(ctx, "42"); err != nil {
		t.Errorf("ServerInstanceGet: expected nil error, got: %v", err)
	}
}

func TestServerInstanceGet_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/server-instances/42": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := setupTestContext(ts.URL)
	if err := ServerInstanceGet(ctx, "42"); err == nil {
		t.Error("ServerInstanceGet with 500: expected error, got nil")
	}
}

func TestServerInstanceConfig_Success(t *testing.T) {
	const cfgJSON = `{"label":"si-42","groupId":10,"serverTypeId":5,"serverId":100,"osTemplateId":1,"hostname":"host1","deployType":"deploy","deployStatus":"finished","revision":1,"updatedTimestamp":"2024-01-01T00:00:00Z"}`
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/server-instances/42/config": testutils.RawHandler(http.StatusOK, cfgJSON),
	})
	defer ts.Close()
	ctx := setupTestContext(ts.URL)
	if err := ServerInstanceConfig(ctx, "42"); err != nil {
		t.Errorf("ServerInstanceConfig: expected nil error, got: %v", err)
	}
}

func TestServerInstanceConfig_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/server-instances/42/config": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := setupTestContext(ts.URL)
	if err := ServerInstanceConfig(ctx, "42"); err == nil {
		t.Error("ServerInstanceConfig with 500: expected error, got nil")
	}
}

func TestServerInstanceCredentials_Success(t *testing.T) {
	const credJSON = `{"username":"admin","initialPassword":"pass123","publicSshKey":"ssh-rsa AAAA"}`
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/server-instances/42/credentials": testutils.RawHandler(http.StatusOK, credJSON),
	})
	defer ts.Close()
	ctx := setupTestContext(ts.URL)
	if err := ServerInstanceCredentials(ctx, "42"); err != nil {
		t.Errorf("ServerInstanceCredentials: expected nil error, got: %v", err)
	}
}

func TestServerInstanceCredentials_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/server-instances/42/credentials": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := setupTestContext(ts.URL)
	if err := ServerInstanceCredentials(ctx, "42"); err == nil {
		t.Error("ServerInstanceCredentials with 500: expected error, got nil")
	}
}

func TestServerInstanceGroupUpdate_Success(t *testing.T) {
	const sigConfigJSON = `{"label":"sig-1","instanceCount":1,"defaultServerTypeId":5,"ipAllocateAuto":1,"ipv4SubnetCreateAuto":1,"processorCount":2,"processorCoreCount":8,"processorCoreMhz":2400,"diskCount":2,"diskSizeMbytes":102400,"diskTypes":[],"virtualInterfacesEnabled":0,"deployType":"deploy","deployStatus":"finished","revision":1,"updatedTimestamp":"2024-01-01T00:00:00Z"}`
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/server-instance-groups/10/config": testutils.RawHandler(http.StatusOK, sigConfigJSON),
	})
	defer ts.Close()
	ctx := setupTestContext(ts.URL)
	if err := ServerInstanceGroupUpdate(ctx, "10", "new-label", 2, 0); err != nil {
		t.Errorf("ServerInstanceGroupUpdate: expected nil error, got: %v", err)
	}
}

func TestServerInstanceGroupUpdate_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/server-instance-groups/10/config": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := setupTestContext(ts.URL)
	if err := ServerInstanceGroupUpdate(ctx, "10", "new-label", 2, 0); err == nil {
		t.Error("ServerInstanceGroupUpdate with 500: expected error, got nil")
	}
}

func TestServerInstanceGroupDelete_Success(t *testing.T) {
	const sigJSON = `{"id":10,"revision":1,"label":"sig-1","infrastructureId":123,"instanceCount":1,"defaultServerTypeId":5,"ipAllocateAuto":1,"ipv4SubnetCreateAuto":1,"processorCount":2,"processorCoreCount":8,"processorCoreMhz":2400,"diskCount":2,"diskSizeMbytes":102400,"diskTypes":[],"virtualInterfacesEnabled":0,"serviceStatus":"active","isVmGroup":0,"isEndpointInstanceGroup":0,"meta":{},"createdTimestamp":"2024-01-01T00:00:00Z","updatedTimestamp":"2024-01-01T00:00:00Z","links":[]}`
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/server-instance-groups/10": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			testutils.RawHandler(http.StatusOK, sigJSON)(w, r)
		},
	})
	defer ts.Close()
	ctx := setupTestContext(ts.URL)
	if err := ServerInstanceGroupDelete(ctx, "10"); err != nil {
		t.Errorf("ServerInstanceGroupDelete: expected nil error, got: %v", err)
	}
}

func TestServerInstanceGroupDelete_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/server-instance-groups/10": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := setupTestContext(ts.URL)
	if err := ServerInstanceGroupDelete(ctx, "10"); err == nil {
		t.Error("ServerInstanceGroupDelete with 500: expected error, got nil")
	}
}

func TestServerInstanceGroupNetworkList_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/server-instance-groups/10/config/networking/connections": testutils.RawHandler(http.StatusOK, `{"data":[],"links":[]}`),
	})
	defer ts.Close()
	ctx := setupTestContext(ts.URL)
	if err := ServerInstanceGroupNetworkList(ctx, "10"); err != nil {
		t.Errorf("ServerInstanceGroupNetworkList: expected nil error, got: %v", err)
	}
}

func TestServerInstanceGroupNetworkList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/server-instance-groups/10/config/networking/connections": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := setupTestContext(ts.URL)
	if err := ServerInstanceGroupNetworkList(ctx, "10"); err == nil {
		t.Error("ServerInstanceGroupNetworkList with 500: expected error, got nil")
	}
}
