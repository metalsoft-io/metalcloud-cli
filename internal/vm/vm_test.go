package vm

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

// makeVM builds a live-shaped VM payload. Note "hosts" is a plain string, the
// way the API actually answers — the strict sdk.VM model declares []string and
// would reject it, which is why every VM endpoint decodes raw bodies.
func makeVM(id int) map[string]any {
	return map[string]any{
		"id": id, "name": "ms-vm-1", "siteId": 1, "infrastructureId": 1664,
		"userId": 1, "instanceId": 1, "vmInstanceId": 1,
		"host": "instance-960.example.local", "hosts": "instance-960.example.local",
		"cpuCores": 1, "ramGB": 2, "diskSizeGB": 10,
		"typeId": 111, "poolId": 70,
		"administrationState": "managed", "powerState": "on",
		"powerStateLastUpdatedTimestamp": "2024-01-01T00:00:00Z",
		"createdTimestamp":               "2024-01-01T00:00:00Z",
		"allocationTimestamp":            "2024-01-01T00:00:00Z",
		"disks":                          []any{},
	}
}

func TestVMList_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vms": testutils.JSONHandler(200, testutils.PaginatedResponse([]any{makeVM(1)}, 1, 1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMList(ctx, VMFilters{}); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMList_Empty(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vms": testutils.JSONHandler(200, testutils.PaginatedResponse([]any{}, 1, 1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMList(ctx, VMFilters{}); err != nil {
		t.Fatalf("expected nil error on empty, got: %v", err)
	}
}

func TestVMList_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vms": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMList(ctx, VMFilters{}); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestVMList_Pagination(t *testing.T) {
	page1 := make([]any, 100)
	page2 := make([]any, 3)
	for i := range page1 {
		page1[i] = makeVM(i + 1)
	}
	for i := range page2 {
		page2[i] = makeVM(i + 101)
	}

	ts := testutils.MultiPageServer("/api/v2/vms", []any{page1, page2})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMList(ctx, VMFilters{}); err != nil {
		t.Fatalf("pagination: expected nil error, got: %v", err)
	}
}

func TestVMList_FiltersAreSent(t *testing.T) {
	var query string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.RawQuery
		testutils.JSONHandler(200, testutils.PaginatedResponse([]any{}, 1, 1))(w, r)
	}))
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	// A single filter value is passed through verbatim; several values are
	// expanded into the $eq / $or DSL by utils.ProcessFilterStringSlice.
	filters := VMFilters{PoolId: []string{"70"}, InfrastructureId: []string{"1", "-2"}}
	if err := VMList(ctx, filters); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	for _, want := range []string{
		"filter.poolId=70",
		"filter.infrastructureId=%24eq%3A1",
		"filter.infrastructureId=%24or%3A%24not%3A%24eq%3A2",
		"sortBy=id%3AASC",
	} {
		if !strings.Contains(query, want) {
			t.Errorf("expected query to contain %s, got: %s", want, query)
		}
	}
}

func TestVMGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vms/1": testutils.JSONHandler(200, makeVM(1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMGet(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMGet(ctx, "bad"); err == nil {
		t.Fatal("expected error for invalid id, got nil")
	}
}

func TestVMGet_404(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vms/999": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMGet(ctx, "999"); err == nil {
		t.Fatal("expected error on 404, got nil")
	}
}

func TestVMUpdate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vms/1": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPatch {
				t.Errorf("expected PATCH, got %s", r.Method)
			}
			testutils.JSONHandler(200, makeVM(1))(w, r)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMUpdate(ctx, "1", []byte(`{"cpuCores":2}`)); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMUpdate_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMUpdate(ctx, "bad", []byte(`{}`)); err == nil {
		t.Fatal("expected error for invalid id, got nil")
	}
}

func TestVMPowerStatus_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vms/1/power-status": testutils.RawHandler(200, `"on"`),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMPowerStatus(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMStart_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vms/1/start": testutils.NoContentHandler(),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMStart(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMShutdown_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vms/1/shutdown": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMShutdown(ctx, "1"); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestVMReboot_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vms/1/reboot": testutils.NoContentHandler(),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMReboot(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMRemoteConsoleInfo_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vms/1/remote-console-info": testutils.JSONHandler(200, map[string]any{
			"active_connections": 0,
			"consoleUrl":         "https://console.example.com/vm/1",
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMRemoteConsoleInfo(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}
