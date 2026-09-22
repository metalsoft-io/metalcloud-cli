package storage

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

// storIfaceRecorder captures the requests received by the mock server so the
// tests can assert on method, path, headers and body.
type storIfaceRecorder struct {
	mu   sync.Mutex
	reqs []storIfaceRequest
}

type storIfaceRequest struct {
	Method string
	Path   string
	Header http.Header
	Body   string
}

func (r *storIfaceRecorder) wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		r.mu.Lock()
		r.reqs = append(r.reqs, storIfaceRequest{
			Method: req.Method,
			Path:   req.URL.Path,
			Header: req.Header.Clone(),
			Body:   string(body),
		})
		r.mu.Unlock()
		next(w, req)
	}
}

func (r *storIfaceRecorder) last() storIfaceRequest {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.reqs[len(r.reqs)-1]
}

func storInterfaceItem(id int) map[string]interface{} {
	return map[string]interface{}{
		"id":            id,
		"revision":      7,
		"storageId":     1,
		"name":          "if-a",
		"nodeIds":       []interface{}{"node-1"},
		"protocols":     []interface{}{"iscsi"},
		"isUplink":      true,
		"useForDeploys": false,
	}
}

func storScopedAccessUserItem(id int) map[string]interface{} {
	return map[string]interface{}{
		"id":               id,
		"revision":         2,
		"storageId":        1,
		"infrastructureId": 10,
		"username":         "scoped-user",
		"createdTimestamp": "2024-01-01T00:00:00Z",
		"updatedTimestamp": "2024-01-02T00:00:00Z",
	}
}

func TestStorageGetInterfaces(t *testing.T) {
	ts := testutils.MultiPageServer("/api/v2/storages/1/interfaces",
		[]interface{}{[]interface{}{storInterfaceItem(3)}})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGetInterfaces(ctx, "1"); err != nil {
		t.Errorf("StorageGetInterfaces: expected nil error, got: %v", err)
	}
}

func TestStorageGetInterfaces_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://127.0.0.1:1")
	if err := StorageGetInterfaces(ctx, "abc"); err == nil {
		t.Error("StorageGetInterfaces with invalid ID: expected error, got nil")
	}
}

func TestStorageGetInterface(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages/1/interfaces/3": testutils.JSONHandler(http.StatusOK, storInterfaceItem(3)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGetInterface(ctx, "1", "3"); err != nil {
		t.Errorf("StorageGetInterface: expected nil error, got: %v", err)
	}
}

func TestStorageGetInterface_InvalidInterfaceId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://127.0.0.1:1")
	if err := StorageGetInterface(ctx, "1", "abc"); err == nil {
		t.Error("StorageGetInterface with invalid interface ID: expected error, got nil")
	}
}

func TestStorageUpdateInterface_SendsIfMatch(t *testing.T) {
	rec := &storIfaceRecorder{}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/storages/1/interfaces/3", rec.wrap(func(w http.ResponseWriter, r *http.Request) {
		testutils.JSONHandler(http.StatusOK, storInterfaceItem(3))(w, r)
	}))
	ts := httptest.NewServer(mux)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageUpdateInterface(ctx, "1", "3", []byte(`{"useForDeploys": true}`)); err != nil {
		t.Fatalf("StorageUpdateInterface: expected nil error, got: %v", err)
	}

	last := rec.last()
	if last.Method != http.MethodPatch {
		t.Errorf("expected PATCH, got %s", last.Method)
	}
	if got := last.Header.Get("If-Match"); got != "7" {
		t.Errorf("expected If-Match 7, got %q", got)
	}
	if strings.TrimSpace(last.Body) != `{"useForDeploys":true}` {
		t.Errorf("unexpected request body: %s", last.Body)
	}
}

func TestStorageUpdate_SendsIfMatchAndBody(t *testing.T) {
	rec := &storIfaceRecorder{}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/storages/1", rec.wrap(func(w http.ResponseWriter, r *http.Request) {
		testutils.JSONHandler(http.StatusOK, storageItem(1, "stor-one"))(w, r)
	}))
	ts := httptest.NewServer(mux)
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageUpdate(ctx, "1", []byte(`{"inMaintenance": 1}`)); err != nil {
		t.Fatalf("StorageUpdate: expected nil error, got: %v", err)
	}

	last := rec.last()
	if last.Method != http.MethodPatch {
		t.Errorf("expected PATCH, got %s", last.Method)
	}
	if got := last.Header.Get("If-Match"); got != "1" {
		t.Errorf("expected If-Match 1 (storage revision), got %q", got)
	}
	if strings.TrimSpace(last.Body) != `{"inMaintenance":1}` {
		t.Errorf("unexpected request body: %s", last.Body)
	}
}

func TestStorageUpdate_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages/1": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageUpdate(ctx, "1", []byte(`{}`)); err == nil {
		t.Error("StorageUpdate with HTTP 500: expected error, got nil")
	}
}

func TestStorageGetScopedAccessUsers(t *testing.T) {
	ts := testutils.MultiPageServer("/api/v2/storages/1/scoped-access-users",
		[]interface{}{[]interface{}{storScopedAccessUserItem(5)}})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGetScopedAccessUsers(ctx, "1"); err != nil {
		t.Errorf("StorageGetScopedAccessUsers: expected nil error, got: %v", err)
	}
}

func TestStorageGetScopedAccessUser(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages/1/scoped-access-users/5": testutils.JSONHandler(http.StatusOK, storScopedAccessUserItem(5)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGetScopedAccessUser(ctx, "1", "5"); err != nil {
		t.Errorf("StorageGetScopedAccessUser: expected nil error, got: %v", err)
	}
}

func TestStorageGetScopedAccessUserCredentials(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages/1/scoped-access-users/5/credentials": testutils.JSONHandler(http.StatusOK, map[string]interface{}{
			"username": "scoped-user",
			"password": "secret",
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGetScopedAccessUserCredentials(ctx, "1", "5"); err != nil {
		t.Errorf("StorageGetScopedAccessUserCredentials: expected nil error, got: %v", err)
	}
}

func TestStorageGetStatistics(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages/1/statistics": testutils.JSONHandler(http.StatusOK, map[string]interface{}{
			"totalSpaceGB": 100,
			"usedSpaceGB":  40,
			"freeSpaceGB":  60,
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StorageGetStatistics(ctx, "1"); err != nil {
		t.Errorf("StorageGetStatistics: expected nil error, got: %v", err)
	}
}

func TestStoragesGetStatistics(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/storages/statistics": testutils.JSONHandler(http.StatusOK, map[string]interface{}{
			"maintenanceCount":  1,
			"experimentalCount": 0,
			"lowSpaceCount":     2,
			"usedSpace":         400,
			"freeSpace":         600,
			"types":             map[string]interface{}{"netapp": 2, "pure": 1},
			"pendingCount":      0,
			"readyCount":        1,
			"activeCount":       3,
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := StoragesGetStatistics(ctx); err != nil {
		t.Errorf("StoragesGetStatistics: expected nil error, got: %v", err)
	}
}

func TestFormatStorageTypeCounts(t *testing.T) {
	if got := formatStorageTypeCounts(nil); got != "" {
		t.Errorf("expected empty string for nil map, got %q", got)
	}

	got := formatStorageTypeCounts(map[string]interface{}{"pure": 1, "netapp": 2})
	if got != "netapp: 2, pure: 1" {
		t.Errorf("expected sorted type counts, got %q", got)
	}
}
