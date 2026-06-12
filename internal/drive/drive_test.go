package drive

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func infraHandler() http.HandlerFunc {
	return testutils.JSONHandler(200, map[string]any{
		"data": []any{
			map[string]any{
				"id": float64(123), "label": "test-infra",
				"serviceStatus": "active", "revision": float64(1),
				"datacenterName": "dc1", "siteId": float64(1),
				"designIsLocked": float64(0),
				"createdTimestamp": "2024-01-01T00:00:00Z",
				"updatedTimestamp": "2024-01-01T00:00:00Z",
				"config": map[string]any{}, "links": []any{},
			},
		},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	})
}

func makeDrive(id int) map[string]any {
	return map[string]any{
		"id": id, "label": "drive-1", "sizeMb": float64(102400),
		"infrastructureId": float64(123), "serviceStatus": "active",
		"storageType": "iscsi", "revision": float64(1),
		"provisioningProtocol": "iscsi",
		"allocationAffinity": "none",
		"storageUpdatedTimestamp": "2024-01-01T00:00:00Z",
		"createdTimestamp": "2024-01-01T00:00:00Z",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
		"infrastructure": map[string]any{"id": float64(123)},
		"config": map[string]any{
			"revision": float64(1), "label": "drive-1",
			"groupId": float64(1), "sizeMb": float64(102400),
			"storageType": "iscsi",
			"updatedTimestamp": "2024-01-01T00:00:00Z",
			"deployType": "deploy", "deployStatus": "not_started",
		},
		"meta": map[string]any{"name": "drive-1"},
		"links": []any{},
	}
}

func makeSnapshot(name string) map[string]any {
	return map[string]any{
		"name":             name,
		"createdTimestamp": "2024-01-01T00:00:00Z",
	}
}

func TestDriveList_HappyPath(t *testing.T) {
	body := map[string]any{
		"data": []any{makeDrive(1)},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":            infraHandler(),
		"/api/v2/infrastructures/123/drives": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DriveList(ctx, "123", nil); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestDriveList_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":            infraHandler(),
		"/api/v2/infrastructures/123/drives": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DriveList(ctx, "123", nil); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestDriveList_Empty(t *testing.T) {
	body := map[string]any{
		"data": []any{},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":            infraHandler(),
		"/api/v2/infrastructures/123/drives": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DriveList(ctx, "123", nil); err != nil {
		t.Fatalf("expected nil error on empty, got: %v", err)
	}
}

func TestDriveList_Pagination(t *testing.T) {
	page1 := make([]any, 100)
	page2 := make([]any, 100)
	page3 := make([]any, 5)
	for i := 0; i < 100; i++ {
		page1[i] = makeDrive(i + 1)
		page2[i] = makeDrive(i + 101)
	}
	for i := 0; i < 5; i++ {
		page3[i] = makeDrive(i + 201)
	}

	callCount := 0
	pages := [][]any{page1, page2, page3}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/drives": func(w http.ResponseWriter, r *http.Request) {
			idx := callCount
			if idx >= len(pages) {
				idx = len(pages) - 1
			}
			callCount++
			resp := testutils.PaginatedResponse(pages[idx], int32(idx+1), int32(len(pages)))
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DriveList(ctx, "123", nil); err != nil {
		t.Fatalf("pagination: expected nil error, got: %v", err)
	}
}

func TestDriveGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":              infraHandler(),
		"/api/v2/infrastructures/123/drives/1": testutils.JSONHandler(200, makeDrive(1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DriveGet(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestDriveGet_404(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                infraHandler(),
		"/api/v2/infrastructures/123/drives/999": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DriveGet(ctx, "123", "999"); err == nil {
		t.Fatal("expected error on 404, got nil")
	}
}

func TestDriveCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":            infraHandler(),
		"/api/v2/infrastructures/123/drives": testutils.JSONHandler(201, makeDrive(2)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"label":"drive-2","sizeMb":51200,"storageType":"iscsi","storagePoolId":1}`)
	if err := DriveCreate(ctx, "123", config); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestDriveDelete_HappyPath(t *testing.T) {
	// DriveDelete calls getDriveIdAndRevision (GET) then DELETE on the same path.
	getHandler := testutils.JSONHandler(200, makeDrive(1))
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/drives/1": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			getHandler(w, r)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DriveDelete(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// DriveSnapshotList is NOT paginated at the SDK level — single call, array response.
func TestDriveSnapshotList_HappyPath(t *testing.T) {
	body := []any{makeSnapshot("snap-1"), makeSnapshot("snap-2")}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                        infraHandler(),
		"/api/v2/infrastructures/123/drives/1/snapshots": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DriveSnapshotList(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestDriveSnapshotList_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                        infraHandler(),
		"/api/v2/infrastructures/123/drives/1/snapshots": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DriveSnapshotList(ctx, "123", "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestDriveUpdateConfig_HappyPath(t *testing.T) {
	// DriveUpdateConfig calls getDriveIdAndRevision (GET on /drives/1) then PATCH on /drives/1/config.
	// SharedDriveConfiguration requires: revision, label, sizeMb, storageType, deployType, deployStatus, updatedTimestamp.
	configBody := map[string]any{
		"revision": float64(1), "label": "drive-updated",
		"sizeMb": float64(102400), "storageType": "iscsi",
		"deployType": "deploy", "deployStatus": "not_started",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                     infraHandler(),
		"/api/v2/infrastructures/123/drives/1":        testutils.JSONHandler(200, makeDrive(1)),
		"/api/v2/infrastructures/123/drives/1/config": testutils.JSONHandler(200, configBody),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"label":"drive-updated"}`)
	if err := DriveUpdateConfig(ctx, "123", "1", config); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestDriveUpdateConfig_GetError(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":              infraHandler(),
		"/api/v2/infrastructures/123/drives/1": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"label":"drive-updated"}`)
	if err := DriveUpdateConfig(ctx, "123", "1", config); err == nil {
		t.Fatal("expected error when GET fails, got nil")
	}
}

func TestDriveUpdateMeta_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                     infraHandler(),
		"/api/v2/infrastructures/123/drives/1/meta":   testutils.JSONHandler(200, makeDrive(1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"name":"new-name"}`)
	if err := DriveUpdateMeta(ctx, "123", "1", config); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestDriveGetHosts_HappyPath(t *testing.T) {
	// SharedDriveHosts has an "instanceGroup" field with willBeConnected/connected/willBeDisconnected arrays.
	hosts := map[string]any{
		"instanceGroup": map[string]any{
			"willBeConnected":    []any{},
			"connected":          []any{float64(1)},
			"willBeDisconnected": []any{},
			"disconnected":       []any{},
		},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                    infraHandler(),
		"/api/v2/infrastructures/123/drives/1/hosts": testutils.JSONHandler(200, hosts),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DriveGetHosts(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestDriveGetHosts_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                    infraHandler(),
		"/api/v2/infrastructures/123/drives/1/hosts": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DriveGetHosts(ctx, "123", "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestDriveUpdateHosts_HappyPath(t *testing.T) {
	hostsResp := map[string]any{
		"instanceGroup": map[string]any{
			"willBeConnected":    []any{},
			"connected":          []any{},
			"willBeDisconnected": []any{},
			"disconnected":       []any{},
		},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/drives/1/actions/modify-server-instance-group-hosts-bulk": testutils.JSONHandler(200, hostsResp),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	// SharedDriveHostsModifyBulk requires "sharedDriveHostBulkOperations" key.
	config := []byte(`{"sharedDriveHostBulkOperations":[]}`)
	if err := DriveUpdateHosts(ctx, "123", "1", config); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestDriveGetConfigInfo_HappyPath(t *testing.T) {
	// SharedDriveConfiguration requires: revision, label, sizeMb, storageType, deployType, deployStatus, updatedTimestamp.
	configInfo := map[string]any{
		"revision":         float64(1),
		"label":            "drive-1",
		"sizeMb":           float64(102400),
		"storageType":      "iscsi",
		"deployType":       "deploy",
		"deployStatus":     "not_started",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                      infraHandler(),
		"/api/v2/infrastructures/123/drives/1/config":  testutils.JSONHandler(200, configInfo),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DriveGetConfigInfo(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestDriveGetConfigInfo_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                     infraHandler(),
		"/api/v2/infrastructures/123/drives/1/config": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DriveGetConfigInfo(ctx, "123", "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}
