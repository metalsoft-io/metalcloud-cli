package template_asset

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func templateAssetFixture(id int) map[string]any {
	return map[string]any{
		"id":         id,
		"templateId": 10,
		"usage":      "build_source_image",
		"revision":   1,
		"createdBy":  1,
		"createdAt":  "2024-01-01T00:00:00Z",
		"file": map[string]any{
			"name":             "image.iso",
			"mimeType":         "application/octet-stream",
			"templatingEngine": false,
			"path":             "/image.iso",
		},
	}
}

func TestTemplateAssetList_HappyPath(t *testing.T) {
	page1 := []map[string]any{
		templateAssetFixture(1),
		templateAssetFixture(2),
	}
	srv := testutils.MultiPageServer("/api/v2/template-assets", []any{page1})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := TemplateAssetList(ctx, nil, nil, nil); err != nil {
		t.Fatalf("TemplateAssetList() unexpected error: %v", err)
	}
}

func TestTemplateAssetList_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/template-assets": testutils.ErrorHandler(500, "server error"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := TemplateAssetList(ctx, nil, nil, nil); err == nil {
		t.Fatal("TemplateAssetList() expected error, got nil")
	}
}

func TestTemplateAssetList_Empty(t *testing.T) {
	srv := testutils.MultiPageServer("/api/v2/template-assets", []any{[]map[string]any{}})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := TemplateAssetList(ctx, nil, nil, nil); err != nil {
		t.Fatalf("TemplateAssetList() unexpected error on empty: %v", err)
	}
}

func TestTemplateAssetList_Pagination(t *testing.T) {
	makeItems := func(start, count int) []map[string]any {
		items := make([]map[string]any, count)
		for i := range items {
			items[i] = templateAssetFixture(start + i)
		}
		return items
	}

	page1 := makeItems(1, 100)
	page2 := makeItems(101, 100)
	page3 := makeItems(201, 5)

	srv := testutils.MultiPageServer("/api/v2/template-assets", []any{page1, page2, page3})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := TemplateAssetList(ctx, nil, nil, nil); err != nil {
		t.Fatalf("TemplateAssetList() pagination error: %v", err)
	}
}

func TestTemplateAssetList_WithFilters(t *testing.T) {
	page1 := []map[string]any{
		templateAssetFixture(3),
	}
	srv := testutils.MultiPageServer("/api/v2/template-assets", []any{page1})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	// Filters are passed as query params; server ignores them in test but the call should succeed
	if err := TemplateAssetList(ctx, []string{"5"}, []string{"logo"}, []string{"image/png"}); err != nil {
		t.Fatalf("TemplateAssetList() with filters unexpected error: %v", err)
	}
}

func TestTemplateAssetGet_HappyPath(t *testing.T) {
	asset := templateAssetFixture(4)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/template-assets/4": testutils.JSONHandler(200, asset),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := TemplateAssetGet(ctx, "4"); err != nil {
		t.Fatalf("TemplateAssetGet() unexpected error: %v", err)
	}
}

func TestTemplateAssetGet_NotFound(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/template-assets/99": testutils.ErrorHandler(404, "not found"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := TemplateAssetGet(ctx, "99"); err == nil {
		t.Fatal("TemplateAssetGet() expected error for not found, got nil")
	}
}

func TestTemplateAssetGet_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := TemplateAssetGet(ctx, "not-a-number"); err == nil {
		t.Fatal("TemplateAssetGet() expected error for invalid ID, got nil")
	}
}

func TestTemplateAssetCreate_HappyPath(t *testing.T) {
	created := templateAssetFixture(30)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/template-assets": testutils.JSONHandler(201, created),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	config := []byte(`{"templateId":10,"usage":"logo","file":{"name":"logo.png","mimeType":"image/png","templatingEngine":false,"path":"/logo.png"}}`)
	if err := TemplateAssetCreate(ctx, config); err != nil {
		t.Fatalf("TemplateAssetCreate() unexpected error: %v", err)
	}
}

func TestTemplateAssetCreate_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/template-assets": testutils.ErrorHandler(400, "bad request"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	config := []byte(`{"templateId":10,"usage":"logo"}`)
	if err := TemplateAssetCreate(ctx, config); err == nil {
		t.Fatal("TemplateAssetCreate() expected error, got nil")
	}
}

func TestTemplateAssetDelete_HappyPath(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/template-assets/9": testutils.RawHandler(204, ""),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := TemplateAssetDelete(ctx, "9"); err != nil {
		t.Fatalf("TemplateAssetDelete() unexpected error: %v", err)
	}
}

func TestTemplateAssetDelete_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/template-assets/9": testutils.ErrorHandler(500, "server error"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := TemplateAssetDelete(ctx, "9"); err == nil {
		t.Fatal("TemplateAssetDelete() expected error, got nil")
	}
}
