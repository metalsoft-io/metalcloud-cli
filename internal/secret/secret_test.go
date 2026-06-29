package secret

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

// secretJSON has all required fields: id, userIdOwner, name, createdTimestamp, updatedTimestamp.
const secretJSON = `{"id":1,"userIdOwner":1,"name":"my-secret","valueEncrypted":"enc-value","createdTimestamp":"2024-01-01T00:00:00Z","updatedTimestamp":"2024-01-01T00:00:00Z"}`

func secretListHandler(statusCode int, count, currentPage, totalPages int) http.HandlerFunc {
	items := make([]string, count)
	for i := range items {
		items[i] = secretJSON
	}
	data := "[" + strings.Join(items, ",") + "]"
	body := fmt.Sprintf(`{"data":%s,"meta":{"currentPage":%d,"totalPages":%d,"itemsPerPage":100}}`,
		data, currentPage, totalPages)
	return testutils.RawHandler(statusCode, body)
}

func secretPage(n int) []any {
	items := make([]any, n)
	for i := range items {
		items[i] = map[string]any{
			"id": i + 1, "userIdOwner": 1,
			"name":             fmt.Sprintf("secret-%d", i+1),
			"createdTimestamp": "2024-01-01T00:00:00Z",
			"updatedTimestamp": "2024-01-01T00:00:00Z",
		}
	}
	return items
}

func TestSecretList_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/secrets": secretListHandler(http.StatusOK, 2, 1, 1),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SecretList(ctx); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestSecretList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/secrets": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SecretList(ctx); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

func TestSecretList_Empty(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/secrets": secretListHandler(http.StatusOK, 0, 1, 1),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SecretList(ctx); err != nil {
		t.Errorf("expected nil error for empty list, got: %v", err)
	}
}

func TestSecretList_Pagination(t *testing.T) {
	ts := testutils.MultiPageServer("/api/v2/secrets", []any{secretPage(100), secretPage(100), secretPage(5)})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SecretList(ctx); err != nil {
		t.Errorf("expected nil error across 3 pages, got: %v", err)
	}
}

func TestSecretGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/secrets/1": testutils.RawHandler(http.StatusOK, secretJSON),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SecretGet(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestSecretGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/secrets/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SecretGet(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestSecretGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SecretGet(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

func TestSecretCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/secrets": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, secretJSON)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := []byte(`{"name":"my-secret","value":"plain-text"}`)
	if err := SecretCreate(ctx, config); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestSecretCreate_BadRequest(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/secrets": testutils.ErrorHandler(http.StatusBadRequest, "validation error"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := []byte(`{}`)
	if err := SecretCreate(ctx, config); err == nil {
		t.Error("expected error for 400, got nil")
	}
}

func TestSecretDelete_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/secrets/1": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SecretDelete(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestSecretDelete_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/secrets/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SecretDelete(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestSecretConfigExample(t *testing.T) {
	ctx := setupTestContext("http://localhost")
	if err := SecretConfigExample(ctx); err != nil {
		t.Errorf("SecretConfigExample() unexpected error: %v", err)
	}
}
