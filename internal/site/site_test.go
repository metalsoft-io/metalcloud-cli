package site

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

// siteListResponse returns a paginated site list JSON with the given items.
func siteItem(id int, name string) map[string]interface{} {
	return map[string]interface{}{
		"id":              id,
		"name":            name,
		"slug":            name,
		"isHidden":        false,
		"isInMaintenance": false,
		"revision":        1,
		"location":        map[string]interface{}{"address": ""},
	}
}

func TestSiteList_HappyPath(t *testing.T) {
	items := []interface{}{siteItem(1, "site-one"), siteItem(2, "site-two")}
	ts := testutils.MultiPageServer("/api/v2/sites", []interface{}{items})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := SiteList(ctx); err != nil {
		t.Errorf("SiteList: expected nil error, got: %v", err)
	}
}

func TestSiteList_Empty(t *testing.T) {
	ts := testutils.MultiPageServer("/api/v2/sites", []interface{}{[]interface{}{}})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := SiteList(ctx); err != nil {
		t.Errorf("SiteList empty: expected nil error, got: %v", err)
	}
}

func TestSiteList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := SiteList(ctx); err == nil {
		t.Error("SiteList with 500: expected error, got nil")
	}
}

func TestSiteList_MultiPage(t *testing.T) {
	page1 := make([]interface{}, 100)
	page2 := make([]interface{}, 100)
	page3 := make([]interface{}, 5)
	for i := range page1 {
		page1[i] = siteItem(i+1, "site-p1")
	}
	for i := range page2 {
		page2[i] = siteItem(100+i+1, "site-p2")
	}
	for i := range page3 {
		page3[i] = siteItem(200+i+1, "site-p3")
	}

	ts := testutils.MultiPageServer("/api/v2/sites", []interface{}{page1, page2, page3})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := SiteList(ctx); err != nil {
		t.Errorf("SiteList multi-page: expected nil error, got: %v", err)
	}
}

func TestSiteGet_Success(t *testing.T) {
	item := siteItem(3, "my-site")
	resp := map[string]interface{}{
		"data": []interface{}{item},
		"meta": testutils.PaginatedMeta(1, 1, 100),
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites": testutils.JSONHandler(http.StatusOK, resp),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	// GetSiteByIdOrLabel searches by name match — "my-site" will match.
	if err := SiteGet(ctx, "my-site"); err != nil {
		t.Errorf("SiteGet: expected nil error, got: %v", err)
	}
}

func TestSiteGet_NotFound(t *testing.T) {
	resp := map[string]interface{}{
		"data": []interface{}{},
		"meta": testutils.PaginatedMeta(1, 0, 100),
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites": testutils.JSONHandler(http.StatusOK, resp),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := SiteGet(ctx, "nonexistent"); err == nil {
		t.Error("SiteGet not-found: expected error, got nil")
	}
}

func TestSiteGet_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites": testutils.ErrorHandler(http.StatusInternalServerError, "server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := SiteGet(ctx, "any"); err == nil {
		t.Error("SiteGet with 500: expected error, got nil")
	}
}

func TestSiteCreate_Success(t *testing.T) {
	item := siteItem(10, "new-site")
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites": testutils.JSONHandler(http.StatusOK, item),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := SiteCreate(ctx, "new-site"); err != nil {
		t.Errorf("SiteCreate: expected nil error, got: %v", err)
	}
}

func TestSiteCreate_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites": testutils.ErrorHandler(http.StatusBadRequest, "bad request"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := SiteCreate(ctx, "bad-site"); err == nil {
		t.Error("SiteCreate with 400: expected error, got nil")
	}
}

func TestSiteUpdate_Success(t *testing.T) {
	searchResp := map[string]interface{}{
		"data": []interface{}{siteItem(3, "my-site")},
		"meta": testutils.PaginatedMeta(1, 1, 100),
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites":   testutils.JSONHandler(http.StatusOK, searchResp),
		"/api/v2/sites/3": testutils.JSONHandler(http.StatusOK, siteItem(3, "updated-site")),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := SiteUpdate(ctx, "my-site", "updated-site"); err != nil {
		t.Errorf("SiteUpdate: expected nil error, got: %v", err)
	}
}

func TestSiteUpdate_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := SiteUpdate(ctx, "my-site", "updated-site"); err == nil {
		t.Error("SiteUpdate with 500: expected error, got nil")
	}
}

func TestSiteGetConfig_Success(t *testing.T) {
	searchResp := map[string]interface{}{
		"data": []interface{}{siteItem(3, "my-site")},
		"meta": testutils.PaginatedMeta(1, 1, 100),
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites":          testutils.JSONHandler(http.StatusOK, searchResp),
		"/api/v2/sites/3/config": testutils.JSONHandler(http.StatusOK, map[string]interface{}{"dnsZoneId": 1}),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := SiteGetConfig(ctx, "my-site"); err != nil {
		t.Errorf("SiteGetConfig: expected nil error, got: %v", err)
	}
}

func TestSiteGetConfig_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := SiteGetConfig(ctx, "my-site"); err == nil {
		t.Error("SiteGetConfig with 500: expected error, got nil")
	}
}

func TestSiteUpdateConfig_Success(t *testing.T) {
	searchResp := map[string]interface{}{
		"data": []interface{}{siteItem(3, "my-site")},
		"meta": testutils.PaginatedMeta(1, 1, 100),
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites":          testutils.JSONHandler(http.StatusOK, searchResp),
		"/api/v2/sites/3":        testutils.JSONHandler(http.StatusOK, siteItem(3, "my-site")),
		"/api/v2/sites/3/config": testutils.JSONHandler(http.StatusOK, map[string]interface{}{"dnsZoneId": 1}),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := SiteUpdateConfig(ctx, "my-site", []byte(`{"dnsZoneId":1}`)); err != nil {
		t.Errorf("SiteUpdateConfig: expected nil error, got: %v", err)
	}
}

func TestSiteUpdateConfig_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := SiteUpdateConfig(ctx, "my-site", []byte(`{"dnsZoneId":1}`)); err == nil {
		t.Error("SiteUpdateConfig with 500: expected error, got nil")
	}
}

func TestSiteDecommission_Success(t *testing.T) {
	searchResp := map[string]interface{}{
		"data": []interface{}{siteItem(3, "my-site")},
		"meta": testutils.PaginatedMeta(1, 1, 100),
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites":                           testutils.JSONHandler(http.StatusOK, searchResp),
		"/api/v2/sites/3/actions/decommission": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := SiteDecommission(ctx, "my-site"); err != nil {
		t.Errorf("SiteDecommission: expected nil error, got: %v", err)
	}
}

func TestSiteDecommission_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := SiteDecommission(ctx, "my-site"); err == nil {
		t.Error("SiteDecommission with 500: expected error, got nil")
	}
}

func TestSiteGetAgents_Success(t *testing.T) {
	searchResp := map[string]interface{}{
		"data": []interface{}{siteItem(3, "my-site")},
		"meta": testutils.PaginatedMeta(1, 1, 100),
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites":               testutils.JSONHandler(http.StatusOK, searchResp),
		"/api/v2/sites/3/controllers": testutils.JSONHandler(http.StatusOK, []interface{}{}),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := SiteGetAgents(ctx, "my-site"); err != nil {
		t.Errorf("SiteGetAgents: expected nil error, got: %v", err)
	}
}

func TestSiteGetAgents_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := SiteGetAgents(ctx, "my-site"); err == nil {
		t.Error("SiteGetAgents with 500: expected error, got nil")
	}
}

func TestSiteOneLiner_Success(t *testing.T) {
	searchResp := map[string]interface{}{
		"data": []interface{}{siteItem(3, "my-site")},
		"meta": testutils.PaginatedMeta(1, 1, 100),
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites": testutils.JSONHandler(http.StatusOK, searchResp),
		"/api/v2/sites/3/controllers/actions/get/one-liner": testutils.JSONHandler(http.StatusOK, map[string]interface{}{"command": "curl http://example.com"}),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := SiteOneLiner(ctx, "my-site", sdk.GenerateSiteControllerOneliner{}); err != nil {
		t.Errorf("SiteOneLiner: expected nil error, got: %v", err)
	}
}

func TestSiteOneLiner_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if err := SiteOneLiner(ctx, "my-site", sdk.GenerateSiteControllerOneliner{}); err == nil {
		t.Error("SiteOneLiner with 500: expected error, got nil")
	}
}

func TestGetSiteByIdOrLabel_ById(t *testing.T) {
	items := []interface{}{siteItem(1, "site-one")}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites": testutils.JSONHandler(http.StatusOK, map[string]interface{}{
			"data": items,
			"meta": map[string]interface{}{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
		}),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	result, err := GetSiteByIdOrLabel(ctx, "1")
	if err != nil {
		t.Fatalf("GetSiteByIdOrLabel(\"1\") unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("GetSiteByIdOrLabel(\"1\") returned nil")
	}
}

func TestGetSiteByIdOrLabel_ByName(t *testing.T) {
	items := []interface{}{siteItem(2, "my-site")}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites": testutils.JSONHandler(http.StatusOK, map[string]interface{}{
			"data": items,
			"meta": map[string]interface{}{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
		}),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	result, err := GetSiteByIdOrLabel(ctx, "my-site")
	if err != nil {
		t.Fatalf("GetSiteByIdOrLabel(\"my-site\") unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("GetSiteByIdOrLabel(\"my-site\") returned nil")
	}
}

func TestGetSiteByIdOrLabel_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/sites": testutils.JSONHandler(http.StatusOK, map[string]interface{}{
			"data": []interface{}{},
			"meta": map[string]interface{}{"currentPage": 1, "totalPages": 1, "itemsPerPage": 100},
		}),
	})
	defer ts.Close()
	ctx := testutils.SetupTestContext(ts.URL)
	if _, err := GetSiteByIdOrLabel(ctx, "no-such"); err == nil {
		t.Error("GetSiteByIdOrLabel: expected error for not-found, got nil")
	}
}
