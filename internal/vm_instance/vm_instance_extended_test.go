package vm_instance

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func makeVMInstance(id int) map[string]any {
	return map[string]any{
		"id":               id,
		"label":            "vm-1",
		"infrastructureId": float64(123),
		"groupId":          float64(10),
		"serviceStatus":    "active",
		"typeId":           float64(5),
		"diskSizeGB":       float64(100),
		"ramGB":            float64(16),
		"cpuCores":         float64(8),
		"revision":         float64(1),
		"createdTimestamp": "2024-01-01T00:00:00Z",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
		"links":            []any{},
		"infrastructure":   map[string]any{"id": float64(123)},
		"config": map[string]any{
			"revision":         float64(1),
			"label":            "vm-1",
			"typeId":           float64(5),
			"deployType":       "deploy",
			"deployStatus":     "not_started",
			"diskSizeGB":       float64(100),
			"ramGB":            float64(16),
			"cpuCores":         float64(8),
			"updatedTimestamp": "2024-01-01T00:00:00Z",
		},
		"meta": map[string]any{},
	}
}

// VMInstanceGet

func TestVMInstanceGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures/123/vm-instances/1": testutils.JSONHandler(200, makeVMInstance(1)),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceGet(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGet_InvalidInfraId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceGet(ctx, "bad", "1"); err == nil {
		t.Fatal("expected error for invalid infra id, got nil")
	}
}

func TestVMInstanceGet_InvalidVMId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceGet(ctx, "123", "bad"); err == nil {
		t.Fatal("expected error for invalid vm id, got nil")
	}
}

func TestVMInstanceGet_404(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures/123/vm-instances/999": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceGet(ctx, "123", "999"); err == nil {
		t.Fatal("expected error on 404, got nil")
	}
}

// VMInstanceGetConfig

func TestVMInstanceGetConfig_HappyPath(t *testing.T) {
	config := map[string]any{
		"revision":         float64(1),
		"label":            "vm-1",
		"typeId":           float64(5),
		"diskSizeGB":       float64(100),
		"ramGB":            float64(16),
		"cpuCores":         float64(8),
		"deployType":       "deploy",
		"deployStatus":     "not_started",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures/123/vm-instances/1/config": testutils.JSONHandler(200, config),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceGetConfig(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGetConfig_InvalidIds(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceGetConfig(ctx, "bad", "1"); err == nil {
		t.Fatal("expected error for invalid infra id, got nil")
	}
	if err := VMInstanceGetConfig(ctx, "123", "bad"); err == nil {
		t.Fatal("expected error for invalid vm id, got nil")
	}
}

func TestVMInstanceGetConfig_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures/123/vm-instances/1/config": testutils.ErrorHandler(500, "internal error"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceGetConfig(ctx, "123", "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

// VMInstanceGetPowerStatus

func TestVMInstanceGetPowerStatus_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures/123/vm-instances/1/power-status": testutils.RawHandler(200, `"on"`),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceGetPowerStatus(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGetPowerStatus_InvalidIds(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceGetPowerStatus(ctx, "bad", "1"); err == nil {
		t.Fatal("expected error for invalid infra id, got nil")
	}
	if err := VMInstanceGetPowerStatus(ctx, "123", "bad"); err == nil {
		t.Fatal("expected error for invalid vm id, got nil")
	}
}

func TestVMInstanceGetPowerStatus_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures/123/vm-instances/1/power-status": testutils.ErrorHandler(500, "internal error"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceGetPowerStatus(ctx, "123", "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

// VMInstancePowerControl — start, shutdown, reboot, invalid action

func TestVMInstanceStart_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures/123/vm-instances/1/start": testutils.RawHandler(204, ""),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstancePowerControl(ctx, "123", "1", "start"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceShutdown_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures/123/vm-instances/1/shutdown": testutils.RawHandler(204, ""),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstancePowerControl(ctx, "123", "1", "shutdown"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceReboot_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures/123/vm-instances/1/reboot": testutils.RawHandler(204, ""),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstancePowerControl(ctx, "123", "1", "reboot"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstancePowerControl_InvalidAction(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstancePowerControl(ctx, "123", "1", "explode"); err == nil {
		t.Fatal("expected error for unsupported action, got nil")
	}
}

func TestVMInstancePowerControl_InvalidIds(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstancePowerControl(ctx, "bad", "1", "start"); err == nil {
		t.Fatal("expected error for invalid infra id, got nil")
	}
	if err := VMInstancePowerControl(ctx, "123", "bad", "start"); err == nil {
		t.Fatal("expected error for invalid vm id, got nil")
	}
}

func TestVMInstancePowerControl_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures/123/vm-instances/1/start": testutils.ErrorHandler(500, "internal error"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstancePowerControl(ctx, "123", "1", "start"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

// VMInstanceGetCredentials — additional cases (invalid ids)

func TestVMInstanceGetCredentials_InvalidInfraId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceGetCredentials(ctx, "bad", "1"); err == nil {
		t.Fatal("expected error for invalid infra id, got nil")
	}
}

func TestVMInstanceGetCredentials_InvalidVMId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceGetCredentials(ctx, "123", "bad"); err == nil {
		t.Fatal("expected error for invalid vm id, got nil")
	}
}

func TestVMInstanceGetCredentials_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures/123/vm-instances/1/credentials": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceGetCredentials(ctx, "123", "1"); err == nil {
		t.Fatal("expected error on 404, got nil")
	}
}

// VMInstanceGroupCreate

func TestVMInstanceGroupCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures/123/vm-instance-groups": testutils.JSONHandler(201, makeVMGroup(2)),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceGroupCreate(ctx, "123", "5", "50", "2", ""); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGroupCreate_InvalidInfraId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceGroupCreate(ctx, "bad", "5", "50", "2", ""); err == nil {
		t.Fatal("expected error for invalid infra id, got nil")
	}
}

func TestVMInstanceGroupCreate_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures/123/vm-instance-groups": testutils.ErrorHandler(500, "internal error"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceGroupCreate(ctx, "123", "5", "50", "2", ""); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

// VMInstanceGroupUpdate — calls GET (for revision) then PATCH

func TestVMInstanceGroupUpdate_HappyPath(t *testing.T) {
	// VMInstanceGroupUpdate calls getVmInstanceGroupIdAndRevision (GET) then PATCH on /config subpath.
	// PATCH returns VMInstanceGroupConfiguration, not the group struct.
	groupConfig := map[string]any{
		"revision":         float64(1),
		"label":            "new-label",
		"deployType":       "deploy",
		"deployStatus":     "not_started",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures/123/vm-instance-groups/1":        testutils.JSONHandler(200, makeVMGroup(1)),
		"/api/v2/infrastructures/123/vm-instance-groups/1/config": testutils.JSONHandler(200, groupConfig),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceGroupUpdate(ctx, "123", "1", "new-label", nil); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGroupUpdate_GetError(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures/123/vm-instance-groups/1": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceGroupUpdate(ctx, "123", "1", "new-label", nil); err == nil {
		t.Fatal("expected error when GET fails, got nil")
	}
}

// VMInstanceGroupDelete — additional invalid id cases

func TestVMInstanceGroupDelete_InvalidInfraId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceGroupDelete(ctx, "bad", "1"); err == nil {
		t.Fatal("expected error for invalid infra id, got nil")
	}
}

func TestVMInstanceGroupDelete_404(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures/123/vm-instance-groups/1": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VMInstanceGroupDelete(ctx, "123", "1"); err == nil {
		t.Fatal("expected error on 404, got nil")
	}
}
