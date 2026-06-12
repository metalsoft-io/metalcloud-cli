package firmware_policy

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func TestFirmwarePolicyList_HappyPath(t *testing.T) {
	page1 := []map[string]any{
		policyFixture(1),
		policyFixture(2),
	}
	srv := testutils.MultiPageServer("/api/v2/firmware/policies", []any{page1})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwarePolicyList(ctx); err != nil {
		t.Fatalf("FirmwarePolicyList() unexpected error: %v", err)
	}
}

func TestFirmwarePolicyList_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/policies": testutils.ErrorHandler(500, "server error"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwarePolicyList(ctx); err == nil {
		t.Fatal("FirmwarePolicyList() expected error, got nil")
	}
}

func TestFirmwarePolicyList_Empty(t *testing.T) {
	srv := testutils.MultiPageServer("/api/v2/firmware/policies", []any{[]map[string]any{}})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwarePolicyList(ctx); err != nil {
		t.Fatalf("FirmwarePolicyList() unexpected error on empty: %v", err)
	}
}

func TestFirmwarePolicyList_Pagination(t *testing.T) {
	makeItems := func(start, count int) []map[string]any {
		items := make([]map[string]any, count)
		for i := range items {
			items[i] = policyFixture(start + i)
		}
		return items
	}

	page1 := makeItems(1, 100)
	page2 := makeItems(101, 100)
	page3 := makeItems(201, 5)

	srv := testutils.MultiPageServer("/api/v2/firmware/policies", []any{page1, page2, page3})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwarePolicyList(ctx); err != nil {
		t.Fatalf("FirmwarePolicyList() pagination error: %v", err)
	}
}

func TestFirmwarePolicyGet_HappyPath(t *testing.T) {
	policy := policyFixture(3)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/policies/3": testutils.JSONHandler(200, policy),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwarePolicyGet(ctx, "3"); err != nil {
		t.Fatalf("FirmwarePolicyGet() unexpected error: %v", err)
	}
}

func TestFirmwarePolicyGet_NotFound(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/policies/99": testutils.ErrorHandler(404, "not found"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwarePolicyGet(ctx, "99"); err == nil {
		t.Fatal("FirmwarePolicyGet() expected error for not found, got nil")
	}
}

func TestFirmwarePolicyGet_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := FirmwarePolicyGet(ctx, "not-a-number"); err == nil {
		t.Fatal("FirmwarePolicyGet() expected error for invalid ID, got nil")
	}
}

func policyFixture(id int) map[string]any {
	return map[string]any{
		"id":               id,
		"label":            "test-policy",
		"status":           "active",
		"action":           "upgrade",
		"createdTimestamp": "2024-01-01T00:00:00Z",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	}
}

func TestFirmwarePolicyCreate_HappyPath(t *testing.T) {
	created := policyFixture(11)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/policies": testutils.JSONHandler(201, created),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	config := []byte(`{"label":"new-policy","action":"upgrade","rules":[]}`)
	if err := FirmwarePolicyCreate(ctx, config); err != nil {
		t.Fatalf("FirmwarePolicyCreate() unexpected error: %v", err)
	}
}

func TestFirmwarePolicyCreate_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/policies": testutils.ErrorHandler(400, "bad request"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	config := []byte(`{"label":"new-policy"}`)
	if err := FirmwarePolicyCreate(ctx, config); err == nil {
		t.Fatal("FirmwarePolicyCreate() expected error, got nil")
	}
}

func TestFirmwarePolicyDelete_HappyPath(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/policies/6": testutils.RawHandler(204, ""),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwarePolicyDelete(ctx, "6"); err != nil {
		t.Fatalf("FirmwarePolicyDelete() unexpected error: %v", err)
	}
}

func TestFirmwarePolicyDelete_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/policies/6": testutils.ErrorHandler(500, "server error"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwarePolicyDelete(ctx, "6"); err == nil {
		t.Fatal("FirmwarePolicyDelete() expected error, got nil")
	}
}

func TestGetGlobalFirmwareConfiguration_HappyPath(t *testing.T) {
	config := map[string]any{"activated": true, "upgradeStartTime": "2024-01-01T00:00:00Z"}
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/configuration": testutils.JSONHandler(200, config),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := GetGlobalFirmwareConfiguration(ctx); err != nil {
		t.Fatalf("GetGlobalFirmwareConfiguration() unexpected error: %v", err)
	}
}

func TestGetGlobalFirmwareConfiguration_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/configuration": testutils.ErrorHandler(500, "server error"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := GetGlobalFirmwareConfiguration(ctx); err == nil {
		t.Fatal("GetGlobalFirmwareConfiguration() expected error, got nil")
	}
}

func TestUpdateGlobalFirmwareConfiguration_HappyPath(t *testing.T) {
	updated := map[string]any{"activated": true}
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/configuration": testutils.JSONHandler(200, updated),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	config := []byte(`{"activated":true}`)
	if err := UpdateGlobalFirmwareConfiguration(ctx, config); err != nil {
		t.Fatalf("UpdateGlobalFirmwareConfiguration() unexpected error: %v", err)
	}
}

func TestUpdateGlobalFirmwareConfiguration_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/configuration": testutils.ErrorHandler(400, "bad request"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	config := []byte(`{"activated":true}`)
	if err := UpdateGlobalFirmwareConfiguration(ctx, config); err == nil {
		t.Fatal("UpdateGlobalFirmwareConfiguration() expected error, got nil")
	}
}

func TestFirmwarePolicyApplyWithGroups_HappyPath(t *testing.T) {
	result := map[string]any{"scheduled": []any{}, "failedToSchedule": []any{}}
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/policies/actions/apply-with-server-instance-groups": testutils.JSONHandler(200, result),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwarePolicyApplyWithGroups(ctx); err != nil {
		t.Fatalf("FirmwarePolicyApplyWithGroups() unexpected error: %v", err)
	}
}

func TestFirmwarePolicyApplyWithGroups_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/policies/actions/apply-with-server-instance-groups": testutils.ErrorHandler(500, "server error"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwarePolicyApplyWithGroups(ctx); err == nil {
		t.Fatal("FirmwarePolicyApplyWithGroups() expected error, got nil")
	}
}

func TestFirmwarePolicyApplyWithoutGroups_HappyPath(t *testing.T) {
	result := map[string]any{"scheduled": []any{}, "failedToSchedule": []any{}}
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/policies/actions/apply-without-server-instance-groups": testutils.JSONHandler(200, result),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwarePolicyApplyWithoutGroups(ctx); err != nil {
		t.Fatalf("FirmwarePolicyApplyWithoutGroups() unexpected error: %v", err)
	}
}

func TestFirmwarePolicyApplyWithoutGroups_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/policies/actions/apply-without-server-instance-groups": testutils.ErrorHandler(500, "server error"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwarePolicyApplyWithoutGroups(ctx); err == nil {
		t.Fatal("FirmwarePolicyApplyWithoutGroups() expected error, got nil")
	}
}

func TestFirmwarePolicyGenerateAudit_HappyPath(t *testing.T) {
	audit := map[string]any{"audit": map[string]any{}}
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/policies/7/actions/generate-audit": testutils.JSONHandler(200, audit),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwarePolicyGenerateAudit(ctx, "7"); err != nil {
		t.Fatalf("FirmwarePolicyGenerateAudit() unexpected error: %v", err)
	}
}

func TestFirmwarePolicyGenerateAudit_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := FirmwarePolicyGenerateAudit(ctx, "bad-id"); err == nil {
		t.Fatal("FirmwarePolicyGenerateAudit() expected error for invalid id, got nil")
	}
}

func TestFirmwarePolicyGenerateAudit_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/policies/7/actions/generate-audit": testutils.ErrorHandler(500, "server error"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwarePolicyGenerateAudit(ctx, "7"); err == nil {
		t.Fatal("FirmwarePolicyGenerateAudit() expected error, got nil")
	}
}

func TestFirmwarePolicyUpdate_HappyPath(t *testing.T) {
	updated := policyFixture(5)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/policies/5": testutils.JSONHandler(200, updated),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwarePolicyUpdate(ctx, "5", []byte(`{"label":"updated"}`)); err != nil {
		t.Fatalf("FirmwarePolicyUpdate() unexpected error: %v", err)
	}
}

func TestFirmwarePolicyUpdate_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := FirmwarePolicyUpdate(ctx, "bad-id", []byte(`{}`)); err == nil {
		t.Fatal("FirmwarePolicyUpdate() expected error for invalid id, got nil")
	}
}

func TestFirmwarePolicyConfigExample(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := FirmwarePolicyConfigExample(ctx); err != nil {
		t.Errorf("FirmwarePolicyConfigExample() unexpected error: %v", err)
	}
}
