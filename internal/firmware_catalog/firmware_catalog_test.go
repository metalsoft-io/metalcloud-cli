package firmware_catalog

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func catalogFixture(id int) map[string]any {
	return map[string]any{
		"id":               id,
		"name":             "Test Catalog",
		"vendor":           "dell",
		"updateType":       "online",
		"createdTimestamp": "2024-01-01T00:00:00Z",
		"links":            []any{},
	}
}

func TestFirmwareCatalogList_HappyPath(t *testing.T) {
	page1 := []map[string]any{
		catalogFixture(1),
		catalogFixture(2),
	}
	srv := testutils.MultiPageServer("/api/v2/firmware/catalog", []any{page1})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareCatalogList(ctx); err != nil {
		t.Fatalf("FirmwareCatalogList() unexpected error: %v", err)
	}
}

func TestFirmwareCatalogList_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/catalog": testutils.ErrorHandler(500, "internal server error"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareCatalogList(ctx); err == nil {
		t.Fatal("FirmwareCatalogList() expected error, got nil")
	}
}

func TestFirmwareCatalogList_Empty(t *testing.T) {
	srv := testutils.MultiPageServer("/api/v2/firmware/catalog", []any{[]map[string]any{}})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareCatalogList(ctx); err != nil {
		t.Fatalf("FirmwareCatalogList() unexpected error on empty: %v", err)
	}
}

func TestFirmwareCatalogList_Pagination(t *testing.T) {
	makeItems := func(start, count int) []map[string]any {
		items := make([]map[string]any, count)
		for i := range items {
			items[i] = catalogFixture(start + i)
		}
		return items
	}

	page1 := makeItems(1, 100)
	page2 := makeItems(101, 100)
	page3 := makeItems(201, 5)

	srv := testutils.MultiPageServer("/api/v2/firmware/catalog", []any{page1, page2, page3})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareCatalogList(ctx); err != nil {
		t.Fatalf("FirmwareCatalogList() pagination error: %v", err)
	}
}

func TestFirmwareCatalogGet_HappyPath(t *testing.T) {
	catalog := catalogFixture(42)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/catalog/42": testutils.JSONHandler(200, catalog),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareCatalogGet(ctx, "42"); err != nil {
		t.Fatalf("FirmwareCatalogGet() unexpected error: %v", err)
	}
}

func TestFirmwareCatalogGet_NotFound(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/catalog/99": testutils.ErrorHandler(404, "not found"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareCatalogGet(ctx, "99"); err == nil {
		t.Fatal("FirmwareCatalogGet() expected error for not found, got nil")
	}
}

func TestFirmwareCatalogGet_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := FirmwareCatalogGet(ctx, "not-a-number"); err == nil {
		t.Fatal("FirmwareCatalogGet() expected error for invalid ID, got nil")
	}
}

func TestFirmwareCatalogDelete_HappyPath(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/catalog/7": testutils.RawHandler(204, ""),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareCatalogDelete(ctx, "7"); err != nil {
		t.Fatalf("FirmwareCatalogDelete() unexpected error: %v", err)
	}
}

func TestFirmwareCatalogDelete_Error(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/catalog/7": testutils.ErrorHandler(500, "server error"),
	})
	defer srv.Close()

	ctx := testutils.SetupTestContext(srv.URL)
	if err := FirmwareCatalogDelete(ctx, "7"); err == nil {
		t.Fatal("FirmwareCatalogDelete() expected error, got nil")
	}
}
