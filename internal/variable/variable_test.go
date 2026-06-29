package variable

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

// variableJSON has all required fields: id, userIdOwner, name, value, createdTimestamp, updatedTimestamp.
const variableJSON = `{"id":1,"userIdOwner":1,"name":"my-var","value":{"key":"val"},"createdTimestamp":"2024-01-01T00:00:00Z","updatedTimestamp":"2024-01-01T00:00:00Z"}`

func variableListHandler(statusCode int, count, currentPage, totalPages int) http.HandlerFunc {
	items := make([]string, count)
	for i := range items {
		items[i] = variableJSON
	}
	data := "[" + strings.Join(items, ",") + "]"
	body := fmt.Sprintf(`{"data":%s,"meta":{"currentPage":%d,"totalPages":%d,"itemsPerPage":100}}`,
		data, currentPage, totalPages)
	return testutils.RawHandler(statusCode, body)
}

func variablePage(n int) []any {
	items := make([]any, n)
	for i := range items {
		items[i] = map[string]any{
			"id": i + 1, "userIdOwner": 1,
			"name":             fmt.Sprintf("var-%d", i+1),
			"value":            map[string]any{"k": "v"},
			"createdTimestamp": "2024-01-01T00:00:00Z",
			"updatedTimestamp": "2024-01-01T00:00:00Z",
		}
	}
	return items
}

// variableStringValueJSON carries `value` as a plain string. The SDK types
// Variable.value as map[string]interface{}, so the typed model rejects this;
// raw-body parsing must not.
const variableStringValueJSON = `{"id":2,"userIdOwner":1,"name":"str-var","value":"plain-string-value","createdTimestamp":"2024-01-01T00:00:00Z","updatedTimestamp":"2024-01-01T00:00:00Z"}`

// TestVariableList_StringValue is a regression test for the SDK<->API type
// mismatch where `variable list` failed with
// "json: cannot unmarshal string into Go struct field ...data.value of type map[string]interface {}".
func TestVariableList_StringValue(t *testing.T) {
	body := fmt.Sprintf(`{"data":[%s,%s],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`,
		variableStringValueJSON, variableJSON)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/variables": testutils.RawHandler(http.StatusOK, body),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VariableList(ctx); err != nil {
		t.Errorf("expected nil error for string-valued variable, got: %v", err)
	}
}

// TestVariableGet_StringValue is the single-variable counterpart of the regression test.
func TestVariableGet_StringValue(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/variables/2": testutils.RawHandler(http.StatusOK, variableStringValueJSON),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VariableGet(ctx, "2"); err != nil {
		t.Errorf("expected nil error for string-valued variable, got: %v", err)
	}
}

func TestVariableList_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/variables": variableListHandler(http.StatusOK, 2, 1, 1),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VariableList(ctx); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestVariableList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/variables": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VariableList(ctx); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

func TestVariableList_Empty(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/variables": variableListHandler(http.StatusOK, 0, 1, 1),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VariableList(ctx); err != nil {
		t.Errorf("expected nil error for empty list, got: %v", err)
	}
}

func TestVariableList_Pagination(t *testing.T) {
	ts := testutils.MultiPageServer("/api/v2/variables", []any{variablePage(100), variablePage(100), variablePage(5)})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VariableList(ctx); err != nil {
		t.Errorf("expected nil error across 3 pages, got: %v", err)
	}
}

func TestVariableGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/variables/1": testutils.RawHandler(http.StatusOK, variableJSON),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VariableGet(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestVariableGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/variables/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VariableGet(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestVariableGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VariableGet(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

func TestVariableCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/variables": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, variableJSON)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := []byte(`{"name":"my-var","value":{"key":"val"}}`)
	if err := VariableCreate(ctx, config); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

// TestVariableCreate_StringValueResponse locks in that create tolerates a
// response whose `value` is a plain string (or empty object) — values the
// regenerated SDK's VariableValue oneOf union cannot unmarshal. Create/update
// bypass the union via a raw request + raw response parse.
func TestVariableCreate_StringValueResponse(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/variables": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, variableStringValueJSON)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := []byte(`{"name":"str-var","value":"plain-string-value"}`)
	if err := VariableCreate(ctx, config); err != nil {
		t.Errorf("expected nil error for string-valued create response, got: %v", err)
	}
}

// TestVariableUpdate_EmptyObjectResponse covers the exact oneOf failure mode:
// a response value of `{}` matches zero union branches in the typed SDK.
func TestVariableUpdate_EmptyObjectResponse(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/variables/1": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"id":1,"userIdOwner":1,"name":"v","value":{},"createdTimestamp":"2024-01-01T00:00:00Z","updatedTimestamp":"2024-01-01T00:00:00Z"}`)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := []byte(`{"name":"v","value":{}}`)
	if err := VariableUpdate(ctx, "1", config); err != nil {
		t.Errorf("expected nil error for empty-object update response, got: %v", err)
	}
}

func TestVariableCreate_BadRequest(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/variables": testutils.ErrorHandler(http.StatusBadRequest, "validation error"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := []byte(`{}`)
	if err := VariableCreate(ctx, config); err == nil {
		t.Error("expected error for 400, got nil")
	}
}

func TestVariableDelete_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/variables/1": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VariableDelete(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestVariableDelete_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/variables/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := VariableDelete(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}
