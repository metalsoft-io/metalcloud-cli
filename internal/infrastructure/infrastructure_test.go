package infrastructure

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func infraItem(id int, label string) map[string]interface{} {
	return map[string]interface{}{
		"id":               id,
		"label":            label,
		"serviceStatus":    "active",
		"userIdOwner":      1,
		"siteId":           1,
		"datacenterName":   "dc1",
		"designIsLocked":   0,
		"revision":         1,
		"createdTimestamp": "2024-01-01T00:00:00Z",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
		"config": map[string]interface{}{
			"label":          label + "-config",
			"deployType":     "deploy",
			"deployStatus":   "finished",
			"datacenterName": "dc1",
			"siteId":         1,
			"revision":       1,
			"updatedTimestamp": "2024-01-01T00:00:00Z",
		},
	}
}

func TestInfrastructureList_HappyPath(t *testing.T) {
	items := []interface{}{infraItem(1, "infra-one"), infraItem(2, "infra-two")}
	ts := testutils.MultiPageServer("/api/v2/infrastructures", []interface{}{items})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureList(ctx, true, true, true); err != nil {
		t.Errorf("InfrastructureList: expected nil error, got: %v", err)
	}
}

func TestInfrastructureList_Empty(t *testing.T) {
	ts := testutils.MultiPageServer("/api/v2/infrastructures", []interface{}{[]interface{}{}})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureList(ctx, true, true, true); err != nil {
		t.Errorf("InfrastructureList empty: expected nil error, got: %v", err)
	}
}

func TestInfrastructureList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureList(ctx, true, true, true); err == nil {
		t.Error("InfrastructureList with 500: expected error, got nil")
	}
}

func TestInfrastructureList_MultiPage(t *testing.T) {
	page1 := make([]interface{}, 100)
	page2 := make([]interface{}, 100)
	page3 := make([]interface{}, 5)
	for i := range page1 {
		page1[i] = infraItem(i+1, "infra-p1")
	}
	for i := range page2 {
		page2[i] = infraItem(100+i+1, "infra-p2")
	}
	for i := range page3 {
		page3[i] = infraItem(200+i+1, "infra-p3")
	}

	ts := testutils.MultiPageServer("/api/v2/infrastructures", []interface{}{page1, page2, page3})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureList(ctx, true, true, true); err != nil {
		t.Errorf("InfrastructureList multi-page: expected nil error, got: %v", err)
	}
}

func TestInfrastructureGet_Success(t *testing.T) {
	item := infraItem(7, "my-infra")
	resp := map[string]interface{}{
		"data": []interface{}{item},
		"meta": testutils.PaginatedMeta(1, 1, 100),
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.JSONHandler(http.StatusOK, resp),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureGet(ctx, "my-infra"); err != nil {
		t.Errorf("InfrastructureGet: expected nil error, got: %v", err)
	}
}

func TestInfrastructureGet_NotFound(t *testing.T) {
	resp := map[string]interface{}{
		"data": []interface{}{},
		"meta": testutils.PaginatedMeta(1, 0, 100),
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.JSONHandler(http.StatusOK, resp),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureGet(ctx, "nonexistent"); err == nil {
		t.Error("InfrastructureGet not-found: expected error, got nil")
	}
}

func TestInfrastructureGet_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.ErrorHandler(http.StatusInternalServerError, "server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureGet(ctx, "any"); err == nil {
		t.Error("InfrastructureGet with 500: expected error, got nil")
	}
}

func TestInfrastructureCreate_Success(t *testing.T) {
	item := infraItem(20, "new-infra")
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.JSONHandler(http.StatusOK, item),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureCreate(ctx, "1", "new-infra"); err != nil {
		t.Errorf("InfrastructureCreate: expected nil error, got: %v", err)
	}
}

func TestInfrastructureCreate_InvalidSiteId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureCreate(ctx, "bad", "label"); err == nil {
		t.Error("InfrastructureCreate with invalid site id: expected error, got nil")
	}
}

func TestInfrastructureUpdate_Success(t *testing.T) {
	searchResp := map[string]interface{}{
		"data": []interface{}{infraItem(123, "my-infra")},
		"meta": testutils.PaginatedMeta(1, 1, 100),
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":        testutils.JSONHandler(http.StatusOK, searchResp),
		"/api/v2/infrastructures/123/config": testutils.JSONHandler(http.StatusOK, infraItem(123, "my-infra")),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureUpdate(ctx, "my-infra", "my-infra", ""); err != nil {
		t.Errorf("InfrastructureUpdate: expected nil error, got: %v", err)
	}
}

func TestInfrastructureUpdate_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureUpdate(ctx, "my-infra", "my-infra", ""); err == nil {
		t.Error("InfrastructureUpdate with 500: expected error, got nil")
	}
}

func TestInfrastructureDelete_Success(t *testing.T) {
	searchResp := map[string]interface{}{
		"data": []interface{}{infraItem(123, "my-infra")},
		"meta": testutils.PaginatedMeta(1, 1, 100),
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.JSONHandler(http.StatusOK, searchResp),
		"/api/v2/infrastructures/123": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureDelete(ctx, "my-infra"); err != nil {
		t.Errorf("InfrastructureDelete: expected nil error, got: %v", err)
	}
}

func TestInfrastructureDelete_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureDelete(ctx, "my-infra"); err == nil {
		t.Error("InfrastructureDelete with 500: expected error, got nil")
	}
}

func TestInfrastructureDeploy_Success(t *testing.T) {
	searchResp := map[string]interface{}{
		"data": []interface{}{infraItem(123, "my-infra")},
		"meta": testutils.PaginatedMeta(1, 1, 100),
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                testutils.JSONHandler(http.StatusOK, searchResp),
		"/api/v2/infrastructures/123/actions/deploy": testutils.JSONHandler(http.StatusOK, infraItem(123, "my-infra")),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureDeploy(ctx, "my-infra", false, false, false, 60, false); err != nil {
		t.Errorf("InfrastructureDeploy: expected nil error, got: %v", err)
	}
}

func TestInfrastructureDeploy_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureDeploy(ctx, "my-infra", false, false, false, 60, false); err == nil {
		t.Error("InfrastructureDeploy with 500: expected error, got nil")
	}
}

func TestInfrastructureRevert_Success(t *testing.T) {
	searchResp := map[string]interface{}{
		"data": []interface{}{infraItem(123, "my-infra")},
		"meta": testutils.PaginatedMeta(1, 1, 100),
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                 testutils.JSONHandler(http.StatusOK, searchResp),
		"/api/v2/infrastructures/123/actions/revert": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureRevert(ctx, "my-infra"); err != nil {
		t.Errorf("InfrastructureRevert: expected nil error, got: %v", err)
	}
}

func TestInfrastructureRevert_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureRevert(ctx, "my-infra"); err == nil {
		t.Error("InfrastructureRevert with 500: expected error, got nil")
	}
}

func TestInfrastructureCancelDeploy_Success(t *testing.T) {
	searchResp := map[string]interface{}{
		"data": []interface{}{infraItem(123, "my-infra")},
		"meta": testutils.PaginatedMeta(1, 1, 100),
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                        testutils.JSONHandler(http.StatusOK, searchResp),
		"/api/v2/infrastructures/123/actions/cancel-deploy": testutils.JSONHandler(http.StatusOK, infraItem(123, "my-infra")),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureCancelDeploy(ctx, "my-infra"); err != nil {
		t.Errorf("InfrastructureCancelDeploy: expected nil error, got: %v", err)
	}
}

func TestInfrastructureCancelDeploy_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureCancelDeploy(ctx, "my-infra"); err == nil {
		t.Error("InfrastructureCancelDeploy with 500: expected error, got nil")
	}
}

func TestInfrastructureGetUsers_Success(t *testing.T) {
	searchResp := map[string]interface{}{
		"data": []interface{}{infraItem(123, "my-infra")},
		"meta": testutils.PaginatedMeta(1, 1, 100),
	}
	const userJSON = `{"id":1,"revision":1,"email":"test@test.com","displayName":"Test User","emailStatus":"verified","language":"en","brand":"default","isBrandManager":false,"lastLoginTimestamp":"2024-01-01T00:00:00Z","lastLoginType":"password","isBlocked":false,"passwordChangeRequired":false,"accessLevel":"admin","isBillable":false,"isTestingMode":false,"authenticatorMustChange":false,"authenticatorCreatedTimestamp":"2024-01-01T00:00:00Z","excludeFromReports":false,"isTestAccount":false,"isArchived":false,"isDatastorePublisher":false,"provider":"local","passwordLastChangedTimestamp":"2024-01-01T00:00:00Z","franchise":"","createdTimestamp":"2024-01-01T00:00:00Z","planType":"","isSuspended":false,"authenticatorEnabled":false,"config":{"revision":1,"displayName":"Test User","emailStatus":"verified","language":"en","brand":"default","isBrandManager":false,"lastLoginTimestamp":"2024-01-01T00:00:00Z","lastLoginType":"password","isBlocked":false,"passwordChangeRequired":false,"accessLevel":"admin","isBillable":false,"isTestingMode":false,"authenticatorMustChange":false,"authenticatorCreatedTimestamp":"2024-01-01T00:00:00Z","excludeFromReports":false,"isTestAccount":false,"isArchived":false,"isDatastorePublisher":false,"provider":"local","passwordLastChangedTimestamp":"2024-01-01T00:00:00Z","franchise":"","createdTimestamp":"2024-01-01T00:00:00Z","planType":"","isSuspended":false,"authenticatorEnabled":false},"meta":{}}`
	usersRaw := `{"data":[` + userJSON + `],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":          testutils.JSONHandler(http.StatusOK, searchResp),
		"/api/v2/infrastructures/123/users": testutils.RawHandler(http.StatusOK, usersRaw),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureGetUsers(ctx, "my-infra"); err != nil {
		t.Errorf("InfrastructureGetUsers: expected nil error, got: %v", err)
	}
}

func TestInfrastructureGetUsers_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureGetUsers(ctx, "my-infra"); err == nil {
		t.Error("InfrastructureGetUsers with 500: expected error, got nil")
	}
}

func TestInfrastructureAddUser_Success(t *testing.T) {
	searchResp := map[string]interface{}{
		"data": []interface{}{infraItem(123, "my-infra")},
		"meta": testutils.PaginatedMeta(1, 1, 100),
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":          testutils.JSONHandler(http.StatusOK, searchResp),
		"/api/v2/infrastructures/123/users": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureAddUser(ctx, "my-infra", "user@test.com", "false"); err != nil {
		t.Errorf("InfrastructureAddUser: expected nil error, got: %v", err)
	}
}

func TestInfrastructureAddUser_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureAddUser(ctx, "my-infra", "user@test.com", "false"); err == nil {
		t.Error("InfrastructureAddUser with 500: expected error, got nil")
	}
}

func TestInfrastructureRemoveUser_Success(t *testing.T) {
	searchResp := map[string]interface{}{
		"data": []interface{}{infraItem(123, "my-infra")},
		"meta": testutils.PaginatedMeta(1, 1, 100),
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":            testutils.JSONHandler(http.StatusOK, searchResp),
		"/api/v2/infrastructures/123/users/5": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureRemoveUser(ctx, "my-infra", "5"); err != nil {
		t.Errorf("InfrastructureRemoveUser: expected nil error, got: %v", err)
	}
}

func TestInfrastructureRemoveUser_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureRemoveUser(ctx, "my-infra", "5"); err == nil {
		t.Error("InfrastructureRemoveUser with 500: expected error, got nil")
	}
}

func TestInfrastructureGetStatistics_Success(t *testing.T) {
	searchResp := map[string]interface{}{
		"data": []interface{}{infraItem(123, "my-infra")},
		"meta": testutils.PaginatedMeta(1, 1, 100),
	}
	statsResp := map[string]interface{}{
		"groupId": 1, "jobsThrownError": 0, "jobsCompleted": 5,
		"serverTypesForUsage": []interface{}{}, "vmPoolsForUsage": []interface{}{}, "storagePoolsForUsage": []interface{}{},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":               testutils.JSONHandler(http.StatusOK, searchResp),
		"/api/v2/infrastructures/123/statistics": testutils.JSONHandler(http.StatusOK, statsResp),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureGetStatistics(ctx, "my-infra"); err != nil {
		t.Errorf("InfrastructureGetStatistics: expected nil error, got: %v", err)
	}
}

func TestInfrastructureGetStatistics_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureGetStatistics(ctx, "my-infra"); err == nil {
		t.Error("InfrastructureGetStatistics with 500: expected error, got nil")
	}
}

func TestInfrastructureGetUserLimits_Success(t *testing.T) {
	searchResp := map[string]interface{}{
		"data": []interface{}{infraItem(123, "my-infra")},
		"meta": testutils.PaginatedMeta(1, 1, 100),
	}
	limitsResp := map[string]interface{}{"computeNodesInstancesToProvisionLimit": 10, "drivesAttachedToInstancesLimit": 20, "infrastructuresLimit": 5}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                testutils.JSONHandler(http.StatusOK, searchResp),
		"/api/v2/infrastructures/123/user-limits": testutils.JSONHandler(http.StatusOK, limitsResp),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureGetUserLimits(ctx, "my-infra"); err != nil {
		t.Errorf("InfrastructureGetUserLimits: expected nil error, got: %v", err)
	}
}

func TestInfrastructureGetUserLimits_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := InfrastructureGetUserLimits(ctx, "my-infra"); err == nil {
		t.Error("InfrastructureGetUserLimits with 500: expected error, got nil")
	}
}

func TestGetInfrastructureByIdOrLabel_ById(t *testing.T) {
	items := []interface{}{infraItem(1, "infra-1")}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.JSONHandler(http.StatusOK, map[string]interface{}{
			"data": items,
			"meta": map[string]interface{}{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
		}),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	result, err := GetInfrastructureByIdOrLabel(ctx, "1")
	if err != nil {
		t.Fatalf("GetInfrastructureByIdOrLabel(\"1\") unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("GetInfrastructureByIdOrLabel(\"1\") returned nil")
	}
}

func TestGetInfrastructureByIdOrLabel_ByLabel(t *testing.T) {
	items := []interface{}{infraItem(2, "my-infra")}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.JSONHandler(http.StatusOK, map[string]interface{}{
			"data": items,
			"meta": map[string]interface{}{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
		}),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	result, err := GetInfrastructureByIdOrLabel(ctx, "my-infra")
	if err != nil {
		t.Fatalf("GetInfrastructureByIdOrLabel(\"my-infra\") unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("GetInfrastructureByIdOrLabel(\"my-infra\") returned nil")
	}
}

func TestGetInfrastructureByIdOrLabel_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.JSONHandler(http.StatusOK, map[string]interface{}{
			"data": []interface{}{},
			"meta": map[string]interface{}{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
		}),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if _, err := GetInfrastructureByIdOrLabel(ctx, "no-such"); err == nil {
		t.Error("GetInfrastructureByIdOrLabel: expected error for not-found, got nil")
	}
}
