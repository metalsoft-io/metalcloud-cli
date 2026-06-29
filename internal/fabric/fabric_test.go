package fabric

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	viper.Set("format", "json")
	m.Run()
}

func setupTestContext(serverURL string) context.Context {
	cfg := sdk.NewConfiguration()
	cfg.Servers = []sdk.ServerConfiguration{{URL: serverURL}}
	client := sdk.NewAPIClient(cfg)
	ctx := context.WithValue(context.Background(), api.ApiClientContextKey, client)
	ctx = context.WithValue(ctx, sdk.ContextAccessToken, "test-api-key")
	return ctx
}

// fabricJSON: fabricConfiguration is the ethernet fabric itself (discriminator on top-level fabricType).
const fabricJSON = `{"id":"1","name":"fabric-1","revision":"1","createdTimestamp":"2024-01-01T00:00:00Z","updatedTimestamp":"2024-01-01T00:00:00Z","fabricConfiguration":{"fabricType":"ethernet"}}`

func fabricListHandler(statusCode int, count, currentPage, totalPages int) http.HandlerFunc {
	items := make([]string, count)
	for i := range items {
		items[i] = fabricJSON
	}
	data := "[" + strings.Join(items, ",") + "]"
	body := fmt.Sprintf(`{"data":%s,"meta":{"currentPage":%d,"totalPages":%d,"itemsPerPage":100}}`,
		data, currentPage, totalPages)
	return testutils.RawHandler(statusCode, body)
}

func fabricPage(n int) []any {
	items := make([]any, n)
	for i := range items {
		items[i] = map[string]any{
			"id": fmt.Sprintf("%d", i+1), "name": fmt.Sprintf("fabric-%d", i+1),
			"revision": "1", "createdTimestamp": "2024-01-01T00:00:00Z",
			"updatedTimestamp": "2024-01-01T00:00:00Z",
			"fabricConfiguration": map[string]any{"fabricType": "ethernet"},
		}
	}
	return items
}

func TestFabricList_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-fabrics": fabricListHandler(http.StatusOK, 2, 1, 1),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricList(ctx); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestFabricList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-fabrics": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricList(ctx); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

func TestFabricList_Empty(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-fabrics": fabricListHandler(http.StatusOK, 0, 1, 1),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricList(ctx); err != nil {
		t.Errorf("expected nil error for empty list, got: %v", err)
	}
}

func TestFabricList_Pagination(t *testing.T) {
	ts := testutils.MultiPageServer("/api/v2/network-fabrics", []any{fabricPage(100), fabricPage(100), fabricPage(5)})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricList(ctx); err != nil {
		t.Errorf("expected nil error across 3 pages, got: %v", err)
	}
}

func TestFabricGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-fabrics/1": testutils.RawHandler(http.StatusOK, fabricJSON),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricGet(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestFabricGet_NotFound(t *testing.T) {
	// FabricGet uses GetFabricByIdOrLabel: first tries numeric GET, then filter-by-name.
	// Return 404 for the numeric GET and empty list for the filter fallback.
	notFoundBody := `{"message":"not found","statusCode":404}`
	emptyListBody := `{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-fabrics/99": testutils.RawHandler(http.StatusNotFound, notFoundBody),
		"/api/v2/network-fabrics":    testutils.RawHandler(http.StatusOK, emptyListBody),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricGet(ctx, "99"); err == nil {
		t.Error("expected error for not-found fabric, got nil")
	}
}

func TestFabricGet_InvalidId_ByLabel_NotFound(t *testing.T) {
	emptyListBody := `{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-fabrics": testutils.RawHandler(http.StatusOK, emptyListBody),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricGet(ctx, "no-such-label"); err == nil {
		t.Error("expected error for missing label, got nil")
	}
}

func TestFabricConfigExample_Ethernet(t *testing.T) {
	ctx := setupTestContext("http://localhost")
	if err := FabricConfigExample(ctx, "ethernet"); err != nil {
		t.Errorf("FabricConfigExample(ethernet) unexpected error: %v", err)
	}
}

func TestFabricConfigExample_InvalidType(t *testing.T) {
	ctx := setupTestContext("http://localhost")
	if err := FabricConfigExample(ctx, "invalid-type"); err == nil {
		t.Error("expected error for invalid fabric type, got nil")
	}
}
