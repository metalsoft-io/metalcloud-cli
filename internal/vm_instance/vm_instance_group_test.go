package vm_instance

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/spf13/viper"
)

func init() {
	viper.Set(formatter.ConfigFormat, "json")
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

func makeVMGroup(id int) map[string]any {
	return map[string]any{
		"id": id, "label": "vmg-1",
		"infrastructureId": float64(123), "serviceStatus": "active",
		"diskSizeGB": float64(50), "revision": float64(1),
		"instanceCount": float64(2),
		"createdTimestamp": "2024-01-01T00:00:00Z",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
		"infrastructure": map[string]any{"id": float64(123)},
		"config": map[string]any{
			"revision": float64(1), "label": "vmg-1",
			"updatedTimestamp": "2024-01-01T00:00:00Z",
			"deployType": "deploy", "deployStatus": "not_started",
		},
		"meta": map[string]any{},
		"links": []any{},
	}
}

func TestVMInstanceGroupList_HappyPath(t *testing.T) {
	body := map[string]any{
		"data": []any{makeVMGroup(1)},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                          infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups":  testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupList(ctx, "123"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGroupList_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                         infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupList(ctx, "123"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestVMInstanceGroupList_Empty(t *testing.T) {
	body := map[string]any{
		"data": []any{},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                         infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupList(ctx, "123"); err != nil {
		t.Fatalf("expected nil error on empty, got: %v", err)
	}
}

func TestVMInstanceGroupGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                           infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/1": testutils.JSONHandler(200, makeVMGroup(1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupGet(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMInstanceGroupGet_404(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                             infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/999": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupGet(ctx, "123", "999"); err == nil {
		t.Fatal("expected error on 404, got nil")
	}
}

func TestVMInstanceGroupDelete_HappyPath(t *testing.T) {
	// VMInstanceGroupDelete calls getVmInstanceGroupIdAndRevision (GET) then DELETE — same path.
	getHandler := testutils.JSONHandler(200, makeVMGroup(1))
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": infraHandler(),
		"/api/v2/infrastructures/123/vm-instance-groups/1": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			getHandler(w, r)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMInstanceGroupDelete(ctx, "123", "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}
