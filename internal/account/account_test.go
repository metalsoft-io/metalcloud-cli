package account

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
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

// accountJSON has all required fields: name, id, revision, limits, config.
const accountJSON = `{"id":1,"name":"acme","revision":1,"limits":{},"config":{"revision":1,"name":"acme"}}`

func accountListHandler(statusCode int, count, currentPage, totalPages int) http.HandlerFunc {
	items := make([]string, count)
	for i := range items {
		items[i] = accountJSON
	}
	data := "[" + strings.Join(items, ",") + "]"
	body := fmt.Sprintf(`{"data":%s,"meta":{"currentPage":%d,"totalPages":%d,"itemsPerPage":100}}`,
		data, currentPage, totalPages)
	return testutils.RawHandler(statusCode, body)
}

func accountPage(n int) []any {
	items := make([]any, n)
	for i := range items {
		items[i] = map[string]any{
			"id": i + 1, "name": fmt.Sprintf("acme-%d", i+1), "revision": 1,
			"limits": map[string]any{},
			"config": map[string]any{"revision": 1, "name": fmt.Sprintf("acme-%d", i+1)},
		}
	}
	return items
}

func TestAccountList_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts": accountListHandler(http.StatusOK, 2, 1, 1),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountList(ctx, false); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestAccountList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountList(ctx, false); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

func TestAccountList_Empty(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts": accountListHandler(http.StatusOK, 0, 1, 1),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountList(ctx, false); err != nil {
		t.Errorf("expected nil error for empty list, got: %v", err)
	}
}

func TestAccountList_Pagination(t *testing.T) {
	ts := testutils.MultiPageServer("/api/v2/accounts", []any{accountPage(100), accountPage(100), accountPage(5)})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountList(ctx, false); err != nil {
		t.Errorf("expected nil error across 3 pages, got: %v", err)
	}
}

func TestAccountGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts/1": testutils.RawHandler(http.StatusOK, accountJSON),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountGet(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestAccountGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountGet(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestAccountGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountGet(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

func TestAccountCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, accountJSON)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := []byte(`{"name":"acme","code":"acme-code"}`)
	if err := AccountCreate(ctx, config); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestAccountCreate_BadRequest(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts": testutils.ErrorHandler(http.StatusBadRequest, "validation error"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := []byte(`{}`)
	if err := AccountCreate(ctx, config); err == nil {
		t.Error("expected error for 400, got nil")
	}
}

func TestAccountArchive_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts/1": testutils.RawHandler(http.StatusOK, accountJSON),
		"/api/v2/accounts/1/actions/archive": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, accountJSON)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountArchive(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestAccountArchive_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/accounts/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := AccountArchive(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}
