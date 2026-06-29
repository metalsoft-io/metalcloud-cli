package storage

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func storageItem(id int, name string) map[string]interface{} {
	return map[string]interface{}{
		"id":             id,
		"siteId":         1,
		"datacenterName": "dc1",
		"driver":         "netapp",
		"technologies":   []interface{}{"block"},
		"type":           "san",
		"name":           name,
		"status":         "active",
		"operationMode":  "normal",
		"managementHost": "storage.host",
		"subnetType":     "wan",
		"revision":       1,
	}
}

func TestStorageList_HappyPath(t *testing.T) {
	items := []interface{}{storageItem(1, "stor-one"), storageItem(2, "stor-two")}
	ts := testutils.MultiPageServer("/api/v2/storages", []interface{}{items})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageList(ctx, nil); err != nil {
		t.Errorf("StorageList: expected nil error, got: %v", err)
	}
}

func TestStorageList_Empty(t *testing.T) {
	ts := testutils.MultiPageServer("/api/v2/storages", []interface{}{[]interface{}{}})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageList(ctx, nil); err != nil {
		t.Errorf("StorageList empty: expected nil error, got: %v", err)
	}
}

func TestStorageList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageList(ctx, nil); err == nil {
		t.Error("StorageList with 500: expected error, got nil")
	}
}

func TestStorageList_MultiPage(t *testing.T) {
	page1 := make([]interface{}, 100)
	page2 := make([]interface{}, 100)
	page3 := make([]interface{}, 5)
	for i := range page1 {
		page1[i] = storageItem(i+1, "stor-p1")
	}
	for i := range page2 {
		page2[i] = storageItem(100+i+1, "stor-p2")
	}
	for i := range page3 {
		page3[i] = storageItem(200+i+1, "stor-p3")
	}

	ts := testutils.MultiPageServer("/api/v2/storages", []interface{}{page1, page2, page3})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageList(ctx, nil); err != nil {
		t.Errorf("StorageList multi-page: expected nil error, got: %v", err)
	}
}

func TestStorageGet_Success(t *testing.T) {
	item := storageItem(5, "my-storage")
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages/5": testutils.JSONHandler(http.StatusOK, item),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGet(ctx, "5"); err != nil {
		t.Errorf("StorageGet: expected nil error, got: %v", err)
	}
}

func TestStorageGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGet(ctx, "not-a-number"); err == nil {
		t.Error("StorageGet with invalid id: expected error, got nil")
	}
}

func TestStorageGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGet(ctx, "99"); err == nil {
		t.Error("StorageGet with 404: expected error, got nil")
	}
}

func TestStorageDelete_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageDelete(ctx, "bad"); err == nil {
		t.Error("StorageDelete with invalid id: expected error, got nil")
	}
}

func TestStorageDelete_Success(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages/3": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			http.NotFound(w, r)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageDelete(ctx, "3"); err != nil {
		t.Errorf("StorageDelete: expected nil error, got: %v", err)
	}
}

func TestStorageDelete_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageDelete(ctx, "99"); err == nil {
		t.Error("StorageDelete with 404: expected error, got nil")
	}
}

func TestStorageCreate_HappyPath(t *testing.T) {
	item := storageItem(10, "new-storage")
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages": testutils.JSONHandler(http.StatusCreated, item),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"siteId":1,"driver":"netapp","technologies":["block"],"name":"new-storage","managementHost":"storage.host","subnetType":"wan"}`)
	if err := StorageCreate(ctx, config); err != nil {
		t.Errorf("StorageCreate: expected nil error, got: %v", err)
	}
}

func TestStorageCreate_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"siteId":1,"driver":"netapp","technologies":["block"],"name":"new-storage","managementHost":"storage.host","subnetType":"wan"}`)
	if err := StorageCreate(ctx, config); err == nil {
		t.Error("StorageCreate with 500: expected error, got nil")
	}
}

func TestStorageGetCredentials_Success(t *testing.T) {
	creds := map[string]any{
		"username": "admin",
		"password": "secret",
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages/5/credentials": testutils.JSONHandler(http.StatusOK, creds),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGetCredentials(ctx, "5"); err != nil {
		t.Errorf("StorageGetCredentials: expected nil error, got: %v", err)
	}
}

func TestStorageGetCredentials_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGetCredentials(ctx, "bad"); err == nil {
		t.Error("StorageGetCredentials with invalid id: expected error, got nil")
	}
}

func TestStorageGetCredentials_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages/5/credentials": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGetCredentials(ctx, "5"); err == nil {
		t.Error("StorageGetCredentials with 404: expected error, got nil")
	}
}

func TestStorageGetDrives_Success(t *testing.T) {
	body := map[string]any{
		"data": []any{},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages/5/drives": testutils.JSONHandler(http.StatusOK, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGetDrives(ctx, "5", 0, 0); err != nil {
		t.Errorf("StorageGetDrives: expected nil error, got: %v", err)
	}
}

func TestStorageGetDrives_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGetDrives(ctx, "bad", 0, 0); err == nil {
		t.Error("StorageGetDrives with invalid id: expected error, got nil")
	}
}

func TestStorageGetDrives_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages/5/drives": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGetDrives(ctx, "5", 0, 0); err == nil {
		t.Error("StorageGetDrives with 500: expected error, got nil")
	}
}

func TestStorageGetFileShares_Success(t *testing.T) {
	body := map[string]any{
		"data": []any{},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages/5/file-shares": testutils.JSONHandler(http.StatusOK, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGetFileShares(ctx, "5", 0, 0); err != nil {
		t.Errorf("StorageGetFileShares: expected nil error, got: %v", err)
	}
}

func TestStorageGetFileShares_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGetFileShares(ctx, "bad", 0, 0); err == nil {
		t.Error("StorageGetFileShares with invalid id: expected error, got nil")
	}
}

func TestStorageGetFileShares_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages/5/file-shares": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGetFileShares(ctx, "5", 0, 0); err == nil {
		t.Error("StorageGetFileShares with 500: expected error, got nil")
	}
}

func TestStorageGetBuckets_Success(t *testing.T) {
	body := map[string]any{
		"data": []any{},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages/5/buckets": testutils.JSONHandler(http.StatusOK, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGetBuckets(ctx, "5", 0, 0); err != nil {
		t.Errorf("StorageGetBuckets: expected nil error, got: %v", err)
	}
}

func TestStorageGetBuckets_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGetBuckets(ctx, "bad", 0, 0); err == nil {
		t.Error("StorageGetBuckets with invalid id: expected error, got nil")
	}
}

func TestStorageGetBuckets_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages/5/buckets": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGetBuckets(ctx, "5", 0, 0); err == nil {
		t.Error("StorageGetBuckets with 500: expected error, got nil")
	}
}

func TestStorageConfigExample(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageConfigExample(ctx); err != nil {
		t.Errorf("StorageConfigExample: expected nil error, got: %v", err)
	}
}
