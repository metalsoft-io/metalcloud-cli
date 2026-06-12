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
		"networkFabricId": float64(1),
		"createdTimestamp": "2024-01-01T00:00:00Z",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
		"tags": []any{}, "gpus": []any{},
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
