package firmware_baseline

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func TestFirmwareBaselineList_HappyPath(t *testing.T) {
	page1 := []map[string]any{
		{"id": 1, "name": "Baseline A"},
		{"id": 2, "name": "Baseline B"},
	}
	srv := testutils.MultiPageServer("/api/v2/firmware/baseline", []any{page1})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareBaselineList(ctx); err != nil {
		t.Fatalf("FirmwareBaselineList() unexpected error: %v", err)
	}
}

func TestFirmwareBaselineList_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/baseline": testutils.ErrorHandler(500, "server error"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareBaselineList(ctx); err == nil {
		t.Fatal("FirmwareBaselineList() expected error, got nil")
	}
}

func TestFirmwareBaselineList_Empty(t *testing.T) {
	srv := testutils.MultiPageServer("/api/v2/firmware/baseline", []any{[]map[string]any{}})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareBaselineList(ctx); err != nil {
		t.Fatalf("FirmwareBaselineList() unexpected error on empty: %v", err)
	}
}

func TestFirmwareBaselineList_Pagination(t *testing.T) {
	makeItems := func(start, count int) []map[string]any {
		items := make([]map[string]any, count)
		for i := range items {
			items[i] = map[string]any{"id": start + i, "name": "baseline"}
		}
		return items
	}

	page1 := makeItems(1, 100)
	page2 := makeItems(101, 100)
	page3 := makeItems(201, 5)

	srv := testutils.MultiPageServer("/api/v2/firmware/baseline", []any{page1, page2, page3})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareBaselineList(ctx); err != nil {
		t.Fatalf("FirmwareBaselineList() pagination error: %v", err)
	}
}

func TestFirmwareBaselineGet_HappyPath(t *testing.T) {
	baseline := map[string]any{"id": 5, "name": "Prod Baseline"}
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/baseline/5": testutils.JSONHandler(200, baseline),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareBaselineGet(ctx, "5"); err != nil {
		t.Fatalf("FirmwareBaselineGet() unexpected error: %v", err)
	}
}

func TestFirmwareBaselineGet_NotFound(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/baseline/99": testutils.ErrorHandler(404, "not found"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareBaselineGet(ctx, "99"); err == nil {
		t.Fatal("FirmwareBaselineGet() expected error for not found, got nil")
	}
}

func TestFirmwareBaselineGet_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := FirmwareBaselineGet(ctx, "not-a-number"); err == nil {
		t.Fatal("FirmwareBaselineGet() expected error for invalid ID, got nil")
	}
}

func TestFirmwareBaselineCreate_HappyPath(t *testing.T) {
	created := map[string]any{"id": 10, "name": "New Baseline"}
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/baseline": testutils.JSONHandler(201, created),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	config := []byte(`{"name":"New Baseline","catalog":["cat-1"]}`)
	if err := FirmwareBaselineCreate(ctx, config); err != nil {
		t.Fatalf("FirmwareBaselineCreate() unexpected error: %v", err)
	}
}

func TestFirmwareBaselineCreate_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/baseline": testutils.ErrorHandler(400, "bad request"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	config := []byte(`{"name":"New Baseline"}`)
	if err := FirmwareBaselineCreate(ctx, config); err == nil {
		t.Fatal("FirmwareBaselineCreate() expected error, got nil")
	}
}

func TestFirmwareBaselineDelete_HappyPath(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/baseline/3": testutils.RawHandler(204, ""),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareBaselineDelete(ctx, "3"); err != nil {
		t.Fatalf("FirmwareBaselineDelete() unexpected error: %v", err)
	}
}

func TestFirmwareBaselineDelete_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/baseline/3": testutils.ErrorHandler(500, "server error"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareBaselineDelete(ctx, "3"); err == nil {
		t.Fatal("FirmwareBaselineDelete() expected error, got nil")
	}
}

func TestFirmwareBaselineConfigExample(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := FirmwareBaselineConfigExample(ctx); err != nil {
		t.Errorf("FirmwareBaselineConfigExample() unexpected error: %v", err)
	}
}
