package vm_pool

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

func vmPoolImportPayload() sdk.VMPoolImportVMs {
	return sdk.VMPoolImportVMs{}
}

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func makeVMPool(id int) map[string]any {
	return map[string]any{
		"id": id, "name": "pool-1", "datacenterName": "dc1",
		"managementHost": "vcenter.example.com", "managementPort": float64(443),
		"type": "vmware", "status": "active", "siteId": float64(1),
		"networkFabricId":  float64(1),
		"createdTimestamp": "2024-01-01T00:00:00Z",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
		"tags":             []any{}, "gpus": []any{},
		"links": []any{},
	}
}

func TestVMPoolList_HappyPath(t *testing.T) {
	body := map[string]any{
		"data": []any{makeVMPool(1)},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolList(ctx, nil); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolList_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolList(ctx, nil); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestVMPoolList_Empty(t *testing.T) {
	body := map[string]any{
		"data": []any{},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolList(ctx, nil); err != nil {
		t.Fatalf("expected nil error on empty, got: %v", err)
	}
}

func TestVMPoolList_Pagination(t *testing.T) {
	page1 := make([]any, 100)
	page2 := make([]any, 100)
	page3 := make([]any, 5)
	for i := 0; i < 100; i++ {
		page1[i] = makeVMPool(i + 1)
		page2[i] = makeVMPool(i + 101)
	}
	for i := 0; i < 5; i++ {
		page3[i] = makeVMPool(i + 201)
	}

	ts := testutils.MultiPageServer("/api/v2/vm-pools", []any{page1, page2, page3})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolList(ctx, nil); err != nil {
		t.Fatalf("pagination: expected nil error, got: %v", err)
	}
}

func TestVMPoolGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools":   testutils.JSONHandler(200, map[string]any{"data": []any{makeVMPool(1)}, "meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100}}),
		"/api/v2/vm-pools/1": testutils.JSONHandler(200, makeVMPool(1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGet(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolGet_404(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools":     testutils.ErrorHandler(404, "not found"),
		"/api/v2/vm-pools/999": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGet(ctx, "999"); err == nil {
		t.Fatal("expected error on 404, got nil")
	}
}

func TestVMPoolCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools": testutils.JSONHandler(201, makeVMPool(2)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"name":"pool-2","datacenterName":"dc1","managementHost":"vcenter2.example.com","managementPort":443,"type":"vmware","siteId":1,"networkFabricId":1}`)
	if err := VMPoolCreate(ctx, config); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolDelete_HappyPath(t *testing.T) {
	// VMPoolDelete is a direct DELETE — no pre-GET needed.
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1": testutils.RawHandler(204, ""),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolDelete(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolDelete_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolDelete(ctx, "bad"); err == nil {
		t.Fatal("expected error for invalid id, got nil")
	}
}

func TestVMPoolGetCredentials_HappyPath(t *testing.T) {
	creds := map[string]any{
		"username": "admin",
		"password": "secret",
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/credentials": testutils.JSONHandler(200, creds),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetCredentials(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolGetCredentials_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetCredentials(ctx, "bad"); err == nil {
		t.Fatal("expected error for invalid id, got nil")
	}
}

func TestVMPoolGetCredentials_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/credentials": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetCredentials(ctx, "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestVMPoolGetClusterHosts_HappyPath(t *testing.T) {
	body := map[string]any{
		"data": []any{},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/cluster-hosts": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetClusterHosts(ctx, "1", 0, 0); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolGetClusterHosts_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetClusterHosts(ctx, "bad", 0, 0); err == nil {
		t.Fatal("expected error for invalid id, got nil")
	}
}

func TestVMPoolGetClusterHosts_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/cluster-hosts": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetClusterHosts(ctx, "1", 0, 0); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestVMPoolImportVMs_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/actions/import-vms": testutils.RawHandler(204, ""),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	sdk := vmPoolImportPayload()
	if err := VMPoolImportVMs(ctx, "1", sdk); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolImportVMs_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	sdk := vmPoolImportPayload()
	if err := VMPoolImportVMs(ctx, "bad", sdk); err == nil {
		t.Fatal("expected error for invalid id, got nil")
	}
}

func TestVMPoolImportVMs_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/actions/import-vms": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	sdk := vmPoolImportPayload()
	if err := VMPoolImportVMs(ctx, "1", sdk); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestVMPoolGetClusterHostInterfaces_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/cluster-hosts/2/interfaces": testutils.JSONHandler(200, []any{}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetClusterHostInterfaces(ctx, "1", "2"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolGetClusterHostInterfaces_InvalidPoolId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetClusterHostInterfaces(ctx, "bad", "2"); err == nil {
		t.Fatal("expected error for invalid pool id, got nil")
	}
}

func TestVMPoolGetClusterHostInterfaces_InvalidHostId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetClusterHostInterfaces(ctx, "1", "bad"); err == nil {
		t.Fatal("expected error for invalid host id, got nil")
	}
}

func TestVMPoolGetClusterHostVMs_HappyPath(t *testing.T) {
	body := map[string]any{
		"data": []any{},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/cluster-hosts/2/vms": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetClusterHostVMs(ctx, "1", "2", 0, 0); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolGetClusterHostVMs_InvalidPoolId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetClusterHostVMs(ctx, "bad", "2", 0, 0); err == nil {
		t.Fatal("expected error for invalid pool id, got nil")
	}
}

// ---------------------------------------------------------------------------
// Fixtures for the update / sync / statistics / cluster host command set.
// Each map satisfies every requiredProperties entry of the matching SDK model.
// ---------------------------------------------------------------------------

func makeVMPoolHost(id int) map[string]any {
	return map[string]any{
		"id": id, "name": "host-1", "poolId": float64(1),
		"address": "192.168.0.1", "port": float64(8443),
		"status": "online", "healthStatus": "normal",
		"allowVMsToBeCreated": true, "allowContainersToBeCreated": true,
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	}
}

func makeVMPoolHostInterface(id int) map[string]any {
	return map[string]any{
		"id": id, "hostId": float64(2), "status": "managed",
		"name": "phy0", "macAddress": "14:18:77:4b:5b:05",
		"networkDevices": []any{},
	}
}

func makeVMPoolHostInterfaceNetworkDevice(id int) map[string]any {
	return map[string]any{
		"id": id, "hostInterfaceId": float64(3),
		"networkDeviceId": float64(4), "networkDeviceInterfaceName": "ethernet1/1/5",
	}
}

func makeVMPoolStatistics() map[string]any {
	return map[string]any{
		"totalRamGB": 67, "freeRamGB": 58, "usedRamGB": 9,
		"totalSpaceGB": 588, "usedSpaceGB": 17, "freeSpaceGB": 571,
		"gpuInfo": []any{map[string]any{"name": "G200eR2", "count": 1}},
	}
}

func TestVMPoolUpdate_HappyPath(t *testing.T) {
	var method string
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1": func(w http.ResponseWriter, r *http.Request) {
			method = r.Method
			testutils.JSONHandler(200, makeVMPool(1))(w, r)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolUpdate(ctx, "1", []byte(`{"description":"updated"}`)); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if method != http.MethodPatch {
		t.Errorf("expected PATCH, got %s", method)
	}
}

func TestVMPoolUpdate_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolUpdate(ctx, "bad", []byte(`{}`)); err == nil {
		t.Fatal("expected error for invalid id, got nil")
	}
}

func TestVMPoolUpdate_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolUpdate(ctx, "1", []byte(`{}`)); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestVMPoolUpdateConfigExample(t *testing.T) {
	ctx := testutils.SetupTestContext("http://unused")
	if err := VMPoolUpdateConfigExample(ctx); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolSync_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/actions/sync": testutils.JSONHandler(200, map[string]any{"jobId": 7, "jobGroupId": 3}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolSync(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolSync_NoBody(t *testing.T) {
	// A 204 leaves the typed return value nil; the command must not panic.
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/actions/sync": testutils.NoContentHandler(),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolSync(ctx, "1"); err != nil {
		t.Fatalf("expected nil error on 204, got: %v", err)
	}
}

func TestVMPoolSync_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolSync(ctx, "bad"); err == nil {
		t.Fatal("expected error for invalid id, got nil")
	}
}

func TestVMPoolSync_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/actions/sync": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolSync(ctx, "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestVMPoolRefresh_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/actions/refresh-information": testutils.JSONHandler(200, makeVMPool(1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolRefresh(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolRefresh_NoBody(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/actions/refresh-information": testutils.NoContentHandler(),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolRefresh(ctx, "1"); err != nil {
		t.Fatalf("expected nil error on 204, got: %v", err)
	}
}

func TestVMPoolRefresh_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolRefresh(ctx, "bad"); err == nil {
		t.Fatal("expected error for invalid id, got nil")
	}
}

func TestVMPoolGetStatistics_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/statistics": testutils.JSONHandler(200, makeVMPoolStatistics()),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetStatistics(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolGetStatistics_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetStatistics(ctx, "bad"); err == nil {
		t.Fatal("expected error for invalid id, got nil")
	}
}

func TestVMPoolGetStatistics_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/statistics": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetStatistics(ctx, "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestVMPoolGetContainers_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/containers": testutils.JSONHandler(200, testutils.PaginatedResponse([]any{}, 1, 1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetContainers(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolGetContainers_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetContainers(ctx, "bad"); err == nil {
		t.Fatal("expected error for invalid id, got nil")
	}
}

func TestVMPoolGetClusterHost_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/cluster-hosts/2": testutils.JSONHandler(200, makeVMPoolHost(2)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetClusterHost(ctx, "1", "2"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolGetClusterHost_InvalidHostId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetClusterHost(ctx, "1", "bad"); err == nil {
		t.Fatal("expected error for invalid host id, got nil")
	}
}

func TestVMPoolUpdateClusterHost_HappyPath(t *testing.T) {
	var method string
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/cluster-hosts/2": func(w http.ResponseWriter, r *http.Request) {
			method = r.Method
			testutils.JSONHandler(200, makeVMPoolHost(2))(w, r)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolUpdateClusterHost(ctx, "1", "2", []byte(`{"allowVMsToBeCreated":false}`)); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if method != http.MethodPatch {
		t.Errorf("expected PATCH, got %s", method)
	}
}

func TestVMPoolUpdateClusterHost_InvalidHostId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolUpdateClusterHost(ctx, "1", "bad", []byte(`{}`)); err == nil {
		t.Fatal("expected error for invalid host id, got nil")
	}
}

func TestVMPoolUpdateClusterHostConfigExample(t *testing.T) {
	ctx := testutils.SetupTestContext("http://unused")
	if err := VMPoolUpdateClusterHostConfigExample(ctx); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolGetClusterHostStatistics_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/cluster-hosts/2/statistics": testutils.JSONHandler(200, makeVMPoolStatistics()),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetClusterHostStatistics(ctx, "1", "2"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolGetClusterHostStatistics_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/cluster-hosts/2/statistics": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetClusterHostStatistics(ctx, "1", "2"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestVMPoolGetClusterHostContainers_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/cluster-hosts/2/containers": testutils.JSONHandler(200, testutils.PaginatedResponse([]any{}, 1, 1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetClusterHostContainers(ctx, "1", "2"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolGetClusterHostInterface_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/cluster-hosts/2/interfaces/3": testutils.JSONHandler(200, makeVMPoolHostInterface(3)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetClusterHostInterface(ctx, "1", "2", "3"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolGetClusterHostInterface_InvalidInterfaceId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetClusterHostInterface(ctx, "1", "2", "bad"); err == nil {
		t.Fatal("expected error for invalid interface id, got nil")
	}
}

func TestVMPoolUpdateClusterHostInterface_HappyPath(t *testing.T) {
	var method string
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/cluster-hosts/2/interfaces/3": func(w http.ResponseWriter, r *http.Request) {
			method = r.Method
			testutils.JSONHandler(200, makeVMPoolHostInterface(3))(w, r)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolUpdateClusterHostInterface(ctx, "1", "2", "3", []byte(`{"status":"managed"}`)); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if method != http.MethodPatch {
		t.Errorf("expected PATCH, got %s", method)
	}
}

func TestVMPoolUpdateClusterHostInterface_InvalidConfig(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	// status is a required property of sdk.UpdateVMPoolHostInterface.
	if err := VMPoolUpdateClusterHostInterface(ctx, "1", "2", "3", []byte(`{}`)); err == nil {
		t.Fatal("expected error for config without status, got nil")
	}
}

func TestVMPoolUpdateClusterHostInterfaceConfigExample(t *testing.T) {
	ctx := testutils.SetupTestContext("http://unused")
	if err := VMPoolUpdateClusterHostInterfaceConfigExample(ctx); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolGetClusterHostInterfaceNetworkDevices_HappyPath(t *testing.T) {
	body := testutils.PaginatedResponse([]any{makeVMPoolHostInterfaceNetworkDevice(1)}, 1, 1)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/cluster-hosts/2/interfaces/3/network-devices": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetClusterHostInterfaceNetworkDevices(ctx, "1", "2", "3"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolGetClusterHostInterfaceNetworkDevices_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/cluster-hosts/2/interfaces/3/network-devices": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetClusterHostInterfaceNetworkDevices(ctx, "1", "2", "3"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestVMPoolGetClusterHostInterfaceNetworkDevice_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/cluster-hosts/2/interfaces/3/network-devices/4": testutils.JSONHandler(200, makeVMPoolHostInterfaceNetworkDevice(4)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetClusterHostInterfaceNetworkDevice(ctx, "1", "2", "3", "4"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolGetClusterHostInterfaceNetworkDevice_InvalidAssignmentId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolGetClusterHostInterfaceNetworkDevice(ctx, "1", "2", "3", "bad"); err == nil {
		t.Fatal("expected error for invalid assignment id, got nil")
	}
}

func TestVMPoolAddClusterHostInterfaceNetworkDevice_FromConfig(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/cluster-hosts/2/interfaces/3/network-devices": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			testutils.JSONHandler(201, makeVMPoolHostInterfaceNetworkDevice(5))(w, r)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"networkDeviceId":4,"networkDeviceInterfaceName":"ethernet1/1/5"}`)
	if err := VMPoolAddClusterHostInterfaceNetworkDevice(ctx, "1", "2", "3", "", "", config); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolAddClusterHostInterfaceNetworkDevice_ResolvesNetworkDevice(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/4": testutils.JSONHandler(200,
			testutils.NetworkDeviceFixture("4", 1, "leaf-1")),
		"/api/v2/vm-pools/1/cluster-hosts/2/interfaces/3/network-devices": testutils.JSONHandler(201,
			makeVMPoolHostInterfaceNetworkDevice(5)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolAddClusterHostInterfaceNetworkDevice(ctx, "1", "2", "3", "4", "ethernet1/1/5", nil); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolAddClusterHostInterfaceNetworkDevice_UnknownNetworkDevice(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices":      testutils.JSONHandler(200, testutils.PaginatedResponse([]any{}, 1, 1)),
		"/api/v2/network-devices/9999": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolAddClusterHostInterfaceNetworkDevice(ctx, "1", "2", "3", "9999", "ethernet1/1/5", nil); err == nil {
		t.Fatal("expected error for unknown network device, got nil")
	}
}

func TestVMPoolRemoveClusterHostInterfaceNetworkDevice_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/cluster-hosts/2/interfaces/3/network-devices/4": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", r.Method)
			}
			testutils.NoContentHandler()(w, r)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolRemoveClusterHostInterfaceNetworkDevice(ctx, "1", "2", "3", "4"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMPoolRemoveClusterHostInterfaceNetworkDevice_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-pools/1/cluster-hosts/2/interfaces/3/network-devices/4": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPoolRemoveClusterHostInterfaceNetworkDevice(ctx, "1", "2", "3", "4"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestVMPoolAddClusterHostInterfaceNetworkDeviceConfigExample(t *testing.T) {
	ctx := testutils.SetupTestContext("http://unused")
	if err := VMPoolAddClusterHostInterfaceNetworkDeviceConfigExample(ctx); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}
