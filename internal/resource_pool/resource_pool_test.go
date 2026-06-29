package resource_pool

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func poolItem(id int, label string) map[string]interface{} {
	return map[string]interface{}{
		"resourcePoolId":               id,
		"resourcePoolLabel":            label,
		"resourcePoolDescription":      "test pool",
		"resourcePoolCreatedTimestamp": "2024-01-01T00:00:00Z",
		"resourcePoolUpdatedTimestamp": "2024-01-01T00:00:00Z",
		"statistics": map[string]interface{}{
			"users":       0,
			"servers":     0,
			"subnetPools": 0,
		},
	}
}

func TestResourcePoolList_HappyPath(t *testing.T) {
	items := []interface{}{poolItem(1, "pool-one"), poolItem(2, "pool-two")}
	ts := testutils.MultiPageServer("/api/v2/resource-pools", []interface{}{items})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	// page=0, limit=0 → fetch-all path
	if err := ResourcePoolList(ctx, 0, 0, ""); err != nil {
		t.Errorf("ResourcePoolList: expected nil error, got: %v", err)
	}
}

func TestResourcePoolList_Empty(t *testing.T) {
	ts := testutils.MultiPageServer("/api/v2/resource-pools", []interface{}{[]interface{}{}})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolList(ctx, 0, 0, ""); err != nil {
		t.Errorf("ResourcePoolList empty: expected nil error, got: %v", err)
	}
}

func TestResourcePoolList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolList(ctx, 0, 0, ""); err == nil {
		t.Error("ResourcePoolList with 500: expected error, got nil")
	}
}

func TestResourcePoolList_ExplicitPage(t *testing.T) {
	items := []interface{}{poolItem(1, "pool-one"), poolItem(2, "pool-two")}
	resp := map[string]interface{}{
		"data": items,
		"meta": testutils.PaginatedMeta(1, 1, 5),
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools": testutils.JSONHandler(http.StatusOK, resp),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	// page=1, limit=5 → explicit pagination path
	if err := ResourcePoolList(ctx, 1, 5, ""); err != nil {
		t.Errorf("ResourcePoolList explicit page: expected nil error, got: %v", err)
	}
}

func TestResourcePoolList_MultiPage(t *testing.T) {
	page1 := make([]interface{}, 100)
	page2 := make([]interface{}, 100)
	page3 := make([]interface{}, 5)
	for i := range page1 {
		page1[i] = poolItem(i+1, "pool-p1")
	}
	for i := range page2 {
		page2[i] = poolItem(100+i+1, "pool-p2")
	}
	for i := range page3 {
		page3[i] = poolItem(200+i+1, "pool-p3")
	}

	ts := testutils.MultiPageServer("/api/v2/resource-pools", []interface{}{page1, page2, page3})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolList(ctx, 0, 0, ""); err != nil {
		t.Errorf("ResourcePoolList multi-page (205): expected nil error, got: %v", err)
	}
}

func TestResourcePoolGet_Success(t *testing.T) {
	item := poolItem(3, "my-pool")
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/3": testutils.JSONHandler(http.StatusOK, item),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolGet(ctx, "3"); err != nil {
		t.Errorf("ResourcePoolGet: expected nil error, got: %v", err)
	}
}

func TestResourcePoolGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolGet(ctx, "bad"); err == nil {
		t.Error("ResourcePoolGet with invalid id: expected error, got nil")
	}
}

func TestResourcePoolGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolGet(ctx, "99"); err == nil {
		t.Error("ResourcePoolGet with 404: expected error, got nil")
	}
}

func TestResourcePoolCreate_Success(t *testing.T) {
	item := poolItem(10, "new-pool")
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools": testutils.JSONHandler(http.StatusOK, item),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolCreate(ctx, "new-pool", "test description"); err != nil {
		t.Errorf("ResourcePoolCreate: expected nil error, got: %v", err)
	}
}

func TestResourcePoolDelete_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolDelete(ctx, "bad"); err == nil {
		t.Error("ResourcePoolDelete with invalid id: expected error, got nil")
	}
}

func TestResourcePoolDelete_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/5": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			http.NotFound(w, r)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolDelete(ctx, "5"); err != nil {
		t.Errorf("ResourcePoolDelete: expected nil error, got: %v", err)
	}
}

func TestResourcePoolDelete_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolDelete(ctx, "99"); err == nil {
		t.Error("ResourcePoolDelete with 404: expected error, got nil")
	}
}

func TestResourcePoolCreate_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools": testutils.ErrorHandler(http.StatusBadRequest, "bad request"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolCreate(ctx, "bad-pool", "desc"); err == nil {
		t.Error("ResourcePoolCreate with 400: expected error, got nil")
	}
}

func TestResourcePoolGetUsers_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/5/users": testutils.RawHandler(http.StatusOK, `[]`),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolGetUsers(ctx, "5"); err != nil {
		t.Errorf("ResourcePoolGetUsers: expected nil error, got: %v", err)
	}
}

func TestResourcePoolGetUsers_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolGetUsers(ctx, "bad"); err == nil {
		t.Error("ResourcePoolGetUsers with invalid id: expected error, got nil")
	}
}

func TestResourcePoolGetUsers_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/5/users": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolGetUsers(ctx, "5"); err == nil {
		t.Error("ResourcePoolGetUsers with 500: expected error, got nil")
	}
}

func TestResourcePoolAddUser_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/user/10/pool/5": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolAddUser(ctx, "5", "10"); err != nil {
		t.Errorf("ResourcePoolAddUser: expected nil error, got: %v", err)
	}
}

func TestResourcePoolAddUser_InvalidPoolId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolAddUser(ctx, "bad", "10"); err == nil {
		t.Error("ResourcePoolAddUser with invalid pool id: expected error, got nil")
	}
}

func TestResourcePoolAddUser_InvalidUserId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolAddUser(ctx, "5", "bad"); err == nil {
		t.Error("ResourcePoolAddUser with invalid user id: expected error, got nil")
	}
}

func TestResourcePoolAddUser_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/user/10/pool/5": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolAddUser(ctx, "5", "10"); err == nil {
		t.Error("ResourcePoolAddUser with 500: expected error, got nil")
	}
}

func TestResourcePoolRemoveUser_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/user/10/pool/5": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolRemoveUser(ctx, "5", "10"); err != nil {
		t.Errorf("ResourcePoolRemoveUser: expected nil error, got: %v", err)
	}
}

func TestResourcePoolRemoveUser_InvalidPoolId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolRemoveUser(ctx, "bad", "10"); err == nil {
		t.Error("ResourcePoolRemoveUser with invalid pool id: expected error, got nil")
	}
}

func TestResourcePoolRemoveUser_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/user/10/pool/5": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolRemoveUser(ctx, "5", "10"); err == nil {
		t.Error("ResourcePoolRemoveUser with 500: expected error, got nil")
	}
}

func TestResourcePoolGetServers_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/5/servers": testutils.RawHandler(http.StatusOK, `[]`),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolGetServers(ctx, "5"); err != nil {
		t.Errorf("ResourcePoolGetServers: expected nil error, got: %v", err)
	}
}

func TestResourcePoolGetServers_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolGetServers(ctx, "bad"); err == nil {
		t.Error("ResourcePoolGetServers with invalid id: expected error, got nil")
	}
}

func TestResourcePoolGetServers_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/5/servers": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolGetServers(ctx, "5"); err == nil {
		t.Error("ResourcePoolGetServers with 500: expected error, got nil")
	}
}

func TestResourcePoolAddServer_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/5/server/7": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolAddServer(ctx, "5", "7"); err != nil {
		t.Errorf("ResourcePoolAddServer: expected nil error, got: %v", err)
	}
}

func TestResourcePoolAddServer_InvalidPoolId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolAddServer(ctx, "bad", "7"); err == nil {
		t.Error("ResourcePoolAddServer with invalid pool id: expected error, got nil")
	}
}

func TestResourcePoolAddServer_InvalidServerId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolAddServer(ctx, "5", "bad"); err == nil {
		t.Error("ResourcePoolAddServer with invalid server id: expected error, got nil")
	}
}

func TestResourcePoolAddServer_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/5/server/7": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolAddServer(ctx, "5", "7"); err == nil {
		t.Error("ResourcePoolAddServer with 500: expected error, got nil")
	}
}

func TestResourcePoolRemoveServer_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/5/server/7": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolRemoveServer(ctx, "5", "7"); err != nil {
		t.Errorf("ResourcePoolRemoveServer: expected nil error, got: %v", err)
	}
}

func TestResourcePoolRemoveServer_InvalidPoolId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolRemoveServer(ctx, "bad", "7"); err == nil {
		t.Error("ResourcePoolRemoveServer with invalid pool id: expected error, got nil")
	}
}

func TestResourcePoolRemoveServer_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/5/server/7": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolRemoveServer(ctx, "5", "7"); err == nil {
		t.Error("ResourcePoolRemoveServer with 500: expected error, got nil")
	}
}

func TestResourcePoolGetSubnetPools_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/5/subnet-pools": testutils.RawHandler(http.StatusOK, `[]`),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolGetSubnetPools(ctx, "5"); err != nil {
		t.Errorf("ResourcePoolGetSubnetPools: expected nil error, got: %v", err)
	}
}

func TestResourcePoolGetSubnetPools_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolGetSubnetPools(ctx, "bad"); err == nil {
		t.Error("ResourcePoolGetSubnetPools with invalid id: expected error, got nil")
	}
}

func TestResourcePoolGetSubnetPools_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/5/subnet-pools": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolGetSubnetPools(ctx, "5"); err == nil {
		t.Error("ResourcePoolGetSubnetPools with 500: expected error, got nil")
	}
}

func TestResourcePoolAddSubnetPool_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/5/subnet-pool/3": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolAddSubnetPool(ctx, "5", "3"); err != nil {
		t.Errorf("ResourcePoolAddSubnetPool: expected nil error, got: %v", err)
	}
}

func TestResourcePoolAddSubnetPool_InvalidPoolId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolAddSubnetPool(ctx, "bad", "3"); err == nil {
		t.Error("ResourcePoolAddSubnetPool with invalid pool id: expected error, got nil")
	}
}

func TestResourcePoolAddSubnetPool_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/5/subnet-pool/3": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolAddSubnetPool(ctx, "5", "3"); err == nil {
		t.Error("ResourcePoolAddSubnetPool with 500: expected error, got nil")
	}
}

func TestResourcePoolRemoveSubnetPool_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/5/subnet-pool/3": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolRemoveSubnetPool(ctx, "5", "3"); err != nil {
		t.Errorf("ResourcePoolRemoveSubnetPool: expected nil error, got: %v", err)
	}
}

func TestResourcePoolRemoveSubnetPool_InvalidPoolId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolRemoveSubnetPool(ctx, "bad", "3"); err == nil {
		t.Error("ResourcePoolRemoveSubnetPool with invalid pool id: expected error, got nil")
	}
}

func TestResourcePoolRemoveSubnetPool_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/resource-pools/5/subnet-pool/3": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := ResourcePoolRemoveSubnetPool(ctx, "5", "3"); err == nil {
		t.Error("ResourcePoolRemoveSubnetPool with 500: expected error, got nil")
	}
}
