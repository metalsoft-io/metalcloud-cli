package file_share

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

func makeFileShare(id int) map[string]any {
	return map[string]any{
		"id": id, "label": "fs-1", "sizeGB": float64(100),
		"infrastructureId": float64(123), "serviceStatus": "active",
		"revision": float64(1),
		"createdTimestamp": "2024-01-01T00:00:00Z",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
		"infrastructure": map[string]any{"id": float64(123)},
		"config": map[string]any{
			"revision": float64(1), "sizeGB": float64(100),
			"updatedTimestamp": "2024-01-01T00:00:00Z",
			"label": "fs-1",
			"deployType": "deploy", "deployStatus": "not_started",
		},
		"meta": map[string]any{"name": "fs-1"},
		"links": []any{},
	}
}

func makeSnapshot(name string) map[string]any {
	return map[string]any{
		"name":             name,
		"createdTimestamp": "2024-01-01T00:00:00Z",
	}
}

func TestFileShareList_HappyPath(t *testing.T) {
	body := map[string]any{
		"data": []any{makeFileShare(1)},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                 infraHandler(),
		"/api/v2/infrastructures/123/file-shares": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := FileShareList(ctx, "123", nil); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestFileShareList_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                 infraHandler(),
		"/api/v2/infrastructures/123/file-shares": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := FileShareList(ctx, "123", nil); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestFileShareList_Empty(t *testing.T) {
	body := map[string]any{
		"data": []any{},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                 infraHandler(),
		"/api/v2/infrastructures/123/file-shares": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := FileShareList(ctx, "123", nil); err != nil {
		t.Fatalf("expected nil error on empty, got: %v", err)
	}
}

func TestFileShareList_Pagination(t *testing.T) {
	page1 := make([]any, 100)
	page2 := make([]any, 100)
	page3 := make([]any, 5)
	for i := 0; i < 100; i++ {
		page1[i] = makeFileShare(i + 1)
		page2[i] = makeFileShare(i + 101)
	}
	for i := 0; i < 5; i++ {
		page3[i] = makeFileShare(i + 201)
	}

	callCount := 0
	pages := [][]any{page1, page2, page3}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/file-shares": func(w http.ResponseWriter, r *http.Request) {
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
	if err := FileShareList(ctx, "123", nil); err != nil {
		t.Fatalf("pagination: expected nil error, got: %v", err)
	}
}

func TestFileShareGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                   infraHandler(),
		"/api/v2/infrastructures/123/file-shares/1": testutils.JSONHandler(200, makeFileShare(1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := FileShareGet(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestFileShareGet_404(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                     infraHandler(),
		"/api/v2/infrastructures/123/file-shares/999": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := FileShareGet(ctx, "123", "999"); err == nil {
		t.Fatal("expected error on 404, got nil")
	}
}

func TestFileShareCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                 infraHandler(),
		"/api/v2/infrastructures/123/file-shares": testutils.JSONHandler(201, makeFileShare(2)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"label":"fs-2","sizeGB":200,"storagePoolId":1}`)
	if err := FileShareCreate(ctx, "123", config); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestFileShareDelete_HappyPath(t *testing.T) {
	// FileShareDelete calls getFileShareIdAndRevision (GET) then DELETE on the same path.
	getHandler := testutils.JSONHandler(200, makeFileShare(1))
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/file-shares/1": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			getHandler(w, r)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := FileShareDelete(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// FileShareSnapshotList is NOT paginated at the SDK level — single call, array response.
func TestFileShareSnapshotList_HappyPath(t *testing.T) {
	body := []any{makeSnapshot("snap-1"), makeSnapshot("snap-2")}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                              infraHandler(),
		"/api/v2/infrastructures/123/file-shares/1/snapshots": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := FileShareSnapshotList(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestFileShareSnapshotList_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                              infraHandler(),
		"/api/v2/infrastructures/123/file-shares/1/snapshots": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := FileShareSnapshotList(ctx, "123", "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestFileShareUpdateConfig_HappyPath(t *testing.T) {
	// FileShareUpdateConfig calls getFileShareIdAndRevision (GET on /file-shares/1) then PATCH on /file-shares/1/config.
	// FileShareConfiguration requires: revision, label, sizeGB, deployType, deployStatus, updatedTimestamp.
	configBody := map[string]any{
		"revision": float64(1), "label": "fs-updated", "sizeGB": float64(100),
		"deployType": "deploy", "deployStatus": "not_started",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                             infraHandler(),
		"/api/v2/infrastructures/123/file-shares/1":          testutils.JSONHandler(200, makeFileShare(1)),
		"/api/v2/infrastructures/123/file-shares/1/config":   testutils.JSONHandler(200, configBody),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"label":"fs-updated"}`)
	if err := FileShareUpdateConfig(ctx, "123", "1", config); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestFileShareUpdateConfig_GetError(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                   infraHandler(),
		"/api/v2/infrastructures/123/file-shares/1": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"label":"fs-updated"}`)
	if err := FileShareUpdateConfig(ctx, "123", "1", config); err == nil {
		t.Fatal("expected error when GET fails, got nil")
	}
}

func TestFileShareUpdateMeta_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                        infraHandler(),
		"/api/v2/infrastructures/123/file-shares/1/meta": testutils.JSONHandler(200, makeFileShare(1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"name":"new-name"}`)
	if err := FileShareUpdateMeta(ctx, "123", "1", config); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestFileShareGetHosts_HappyPath(t *testing.T) {
	// FileShareHosts has an "instanceGroup" field of type FileShareHostType.
	hosts := map[string]any{
		"instanceGroup": map[string]any{
			"willBeConnected":    []any{},
			"connected":          []any{"1"},
			"willBeDisconnected": []any{},
			"disconnected":       []any{},
		},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                         infraHandler(),
		"/api/v2/infrastructures/123/file-shares/1/hosts": testutils.JSONHandler(200, hosts),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := FileShareGetHosts(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestFileShareGetHosts_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                         infraHandler(),
		"/api/v2/infrastructures/123/file-shares/1/hosts": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := FileShareGetHosts(ctx, "123", "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestFileShareUpdateHosts_HappyPath(t *testing.T) {
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
		"/api/v2/infrastructures/123/file-shares/1/actions/modify-instance-array-hosts-bulk": testutils.JSONHandler(200, hostsResp),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	// FileShareHostsModifyBulk requires "fileShareHostBulkOperations" key.
	config := []byte(`{"fileShareHostBulkOperations":[]}`)
	if err := FileShareUpdateHosts(ctx, "123", "1", config); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestFileShareGetConfigInfo_HappyPath(t *testing.T) {
	// FileShareConfiguration requires: revision, label, sizeGB, deployType, deployStatus, updatedTimestamp.
	configInfo := map[string]any{
		"revision":         float64(1),
		"label":            "fs-1",
		"sizeGB":           float64(100),
		"deployType":       "deploy",
		"deployStatus":     "not_started",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                          infraHandler(),
		"/api/v2/infrastructures/123/file-shares/1/config": testutils.JSONHandler(200, configInfo),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := FileShareGetConfigInfo(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestFileShareGetConfigInfo_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                          infraHandler(),
		"/api/v2/infrastructures/123/file-shares/1/config": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := FileShareGetConfigInfo(ctx, "123", "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

