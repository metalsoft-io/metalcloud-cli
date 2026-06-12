package vm_type

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func makeVMType(id int) map[string]any {
	return map[string]any{
		"id": id, "name": "standard-2cpu-4gb",
		"cpuCores": float64(2), "ramGB": float64(4),
		"gpuInfo": []any{}, "tags": []any{},
		"links": []any{},
	}
}

func TestVMTypeList_FetchAll(t *testing.T) {
	// limit=0, page=0 means fetch-all via pagination.
	body := map[string]any{
		"data": []any{makeVMType(1), makeVMType(2)},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-types": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMTypeList(ctx, 0, 0); err != nil {
		t.Fatalf("VMTypeList(0,0): expected nil error, got: %v", err)
	}
}

func TestVMTypeList_SinglePage(t *testing.T) {
	// limit=5, page=1 — single page request.
	body := map[string]any{
		"data": []any{makeVMType(1)},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 5},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-types": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMTypeList(ctx, 5, 1); err != nil {
		t.Fatalf("VMTypeList(5,1): expected nil error, got: %v", err)
	}
}

func TestVMTypeList_500(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-types": testutils.ErrorHandler(500, "internal server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMTypeList(ctx, 0, 0); err == nil {
		t.Fatal("expected error on 500, got nil")
	}
}

func TestVMTypeList_Empty(t *testing.T) {
	body := map[string]any{
		"data": []any{},
		"meta": map[string]any{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-types": testutils.JSONHandler(200, body),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMTypeList(ctx, 0, 0); err != nil {
		t.Fatalf("expected nil error on empty, got: %v", err)
	}
}

func TestVMTypeList_Pagination(t *testing.T) {
	page1 := make([]any, 100)
	page2 := make([]any, 100)
	page3 := make([]any, 5)
	for i := 0; i < 100; i++ {
		page1[i] = makeVMType(i + 1)
		page2[i] = makeVMType(i + 101)
	}
	for i := 0; i < 5; i++ {
		page3[i] = makeVMType(i + 201)
	}

	ts := testutils.MultiPageServer("/api/v2/vm-types", []any{page1, page2, page3})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMTypeList(ctx, 0, 0); err != nil {
		t.Fatalf("pagination: expected nil error, got: %v", err)
	}
}

func TestVMTypeGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-types/1": testutils.JSONHandler(200, makeVMType(1)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMTypeGet(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMTypeGet_404(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-types/999": testutils.ErrorHandler(404, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMTypeGet(ctx, "999"); err == nil {
		t.Fatal("expected error on 404, got nil")
	}
}

func TestVMTypeCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-types": testutils.JSONHandler(201, makeVMType(3)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"name":"large-4cpu-8gb","cpuCores":4,"ramGB":8}`)
	if err := VMTypeCreate(ctx, config); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestVMTypeDelete_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/vm-types/1": testutils.RawHandler(204, ""),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := VMTypeDelete(ctx, "1"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}
