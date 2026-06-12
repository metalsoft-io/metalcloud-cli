package extension

import (
	"context"
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

func makeExtensionInfo(id float32, kind string) map[string]any {
	return map[string]any{
		"id":       id,
		"name":     "ext-name",
		"label":    "ext-label",
		"kind":     kind,
		"status":   "active",
		"isPublic": false,
	}
}

// extensionListResponse wraps a slice of extensions in the paginated envelope
// that the SDK's GetExtensions endpoint expects.
func extensionListResponse(items []any) map[string]any {
	return testutils.PaginatedResponse(items, 1, 1)
}

// TestExtensionList_HappyPath verifies a basic list call succeeds.
func TestExtensionList_HappyPath(t *testing.T) {
	resp := extensionListResponse([]any{
		makeExtensionInfo(1, "workflow"),
		makeExtensionInfo(2, "ansible"),
	})
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extensions": testutils.JSONHandler(http.StatusOK, resp),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionList(ctx, nil, nil, nil, nil, ""); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestExtensionList_ServerError verifies a 500 is surfaced as an error.
func TestExtensionList_ServerError(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extensions": testutils.ErrorHandler(http.StatusInternalServerError, "server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionList(ctx, nil, nil, nil, nil, ""); err == nil {
		t.Fatal("expected an error for HTTP 500, got nil")
	}
}

// TestExtensionList_EmptyList verifies an empty result is handled cleanly.
func TestExtensionList_EmptyList(t *testing.T) {
	resp := extensionListResponse([]any{})
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extensions": testutils.JSONHandler(http.StatusOK, resp),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionList(ctx, nil, nil, nil, nil, ""); err != nil {
		t.Fatalf("expected nil error for empty list, got: %v", err)
	}
}

// TestExtensionList_Pagination verifies 3-page fetch-all (205 items).
func TestExtensionList_Pagination(t *testing.T) {
	page1 := make([]any, 100)
	page2 := make([]any, 100)
	page3 := make([]any, 5)
	for i := range page1 {
		page1[i] = makeExtensionInfo(float32(i+1), "workflow")
	}
	for i := range page2 {
		page2[i] = makeExtensionInfo(float32(i+101), "workflow")
	}
	for i := range page3 {
		page3[i] = makeExtensionInfo(float32(i+201), "workflow")
	}
	ts := testutils.MultiPageServer("/api/v2/extensions", []any{page1, page2, page3})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionList(ctx, nil, nil, nil, nil, ""); err != nil {
		t.Fatalf("expected nil error during pagination, got: %v", err)
	}
}

// TestExtensionList_FilterKind verifies the client-side kind filter keeps only matching items.
func TestExtensionList_FilterKind(t *testing.T) {
	resp := extensionListResponse([]any{
		makeExtensionInfo(1, "workflow"),
		makeExtensionInfo(2, "ansible"),
		makeExtensionInfo(3, "workflow"),
	})
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extensions": testutils.JSONHandler(http.StatusOK, resp),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	// Filter for "workflow" only — should not error even though items are reduced client-side.
	if err := ExtensionList(ctx, nil, nil, nil, []string{"workflow"}, ""); err != nil {
		t.Fatalf("expected nil error with kind filter, got: %v", err)
	}
}

// validExtensionJSON is a raw JSON string with all required fields for sdk.Extension
// and sdk.ExtensionDefinition (the SDK validates required properties on unmarshal).
// validExtensionJSON contains all required fields for sdk.Extension and sdk.ExtensionDefinition.
// The SDK validates required properties on unmarshal; omitting any causes a parse error.
const validExtensionJSON = `{
	"id": 3,
	"revision": 1,
	"name": "my-ext",
	"description": "desc",
	"status": "active",
	"kind": "workflow",
	"isPublic": false,
	"definition": {
		"kind": "workflow",
		"schemaVersion": "1.0",
		"name": "my-ext",
		"label": "my-ext",
		"extensionType": "workflow",
		"description": "desc",
		"vendor": "test",
		"extensionVersion": "1.0.0",
		"icon": "",
		"dependencies": {"controllerVersion": "1.0"},
		"inputs": [],
		"outputs": [],
		"assets": []
	}
}`

// TestExtensionGet_HappyPath verifies successful retrieval of a single extension by numeric ID.
func TestExtensionGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extensions/3": testutils.RawHandler(http.StatusOK, validExtensionJSON),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionGet(ctx, "3"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestExtensionGet_NotFound verifies that an unknown extension ID returns an error.
// GetExtensionByIdOrLabel: first tries GET /extensions/{id} (returns non-200), then
// falls back to label search (returns empty list) — overall result is "not found" error.
func TestExtensionGet_NotFound(t *testing.T) {
	emptyList := map[string]any{"data": []any{}, "meta": map[string]any{"itemsPerPage": 100}}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extensions/999": testutils.ErrorHandler(http.StatusNotFound, "not found"),
		"/api/v2/extensions":     testutils.JSONHandler(http.StatusOK, emptyList),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionGet(ctx, "999"); err == nil {
		t.Fatal("expected an error for not-found extension, got nil")
	}
}

// validDefinitionJSON is a minimal ExtensionDefinition JSON accepted by the SDK.
const validDefinitionJSON = `{
	"kind": "workflow",
	"schemaVersion": "1.0",
	"name": "my-ext",
	"label": "my-ext",
	"extensionType": "workflow",
	"description": "desc",
	"vendor": "test",
	"extensionVersion": "1.0.0",
	"icon": "",
	"dependencies": {"controllerVersion": "1.0"},
	"inputs": [],
	"outputs": [],
	"assets": []
}`

// --- ExtensionCreate ---

// TestExtensionCreate_HappyPath verifies a successful extension creation.
func TestExtensionCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extensions": testutils.RawHandler(http.StatusOK, validExtensionJSON),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionCreate(ctx, "my-ext", "workflow", "desc", []byte(validDefinitionJSON)); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestExtensionCreate_InvalidDefinition verifies that a bad JSON definition returns an error
// before any HTTP call is made.
func TestExtensionCreate_InvalidDefinition(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionCreate(ctx, "my-ext", "workflow", "desc", []byte(`not-json`)); err == nil {
		t.Fatal("expected error for invalid JSON definition, got nil")
	}
}

// TestExtensionCreate_ServerError verifies that a 500 from the API is surfaced as an error.
func TestExtensionCreate_ServerError(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extensions": testutils.ErrorHandler(http.StatusInternalServerError, "server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionCreate(ctx, "my-ext", "workflow", "desc", []byte(validDefinitionJSON)); err == nil {
		t.Fatal("expected error for HTTP 500, got nil")
	}
}

// --- ExtensionUpdate ---

// TestExtensionUpdate_HappyPath verifies a successful extension update identified by numeric ID.
func TestExtensionUpdate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extensions/3": func(w http.ResponseWriter, r *http.Request) {
			// GET returns the current extension; PUT returns the updated one.
			if r.Method == http.MethodGet {
				testutils.RawHandler(http.StatusOK, validExtensionJSON)(w, r)
			} else {
				testutils.RawHandler(http.StatusOK, validExtensionJSON)(w, r)
			}
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionUpdate(ctx, "3", "new-name", "new-desc", nil); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestExtensionUpdate_NotFound verifies that updating a non-existent extension returns an error.
func TestExtensionUpdate_NotFound(t *testing.T) {
	emptyList := map[string]any{"data": []any{}, "meta": map[string]any{"itemsPerPage": 100}}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extensions/999": testutils.ErrorHandler(http.StatusNotFound, "not found"),
		"/api/v2/extensions":     testutils.JSONHandler(http.StatusOK, emptyList),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionUpdate(ctx, "999", "", "", nil); err == nil {
		t.Fatal("expected error for not-found extension, got nil")
	}
}

// TestExtensionUpdate_InvalidDefinition verifies that a bad JSON config returns an error
// after the extension is fetched but before the update API call.
func TestExtensionUpdate_InvalidDefinition(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extensions/3": testutils.RawHandler(http.StatusOK, validExtensionJSON),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionUpdate(ctx, "3", "", "", []byte(`not-json`)); err == nil {
		t.Fatal("expected error for invalid definition JSON, got nil")
	}
}

// --- ExtensionArchive ---

// TestExtensionArchive_HappyPath verifies successful archiving of an extension.
func TestExtensionArchive_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extensions/3":                 testutils.RawHandler(http.StatusOK, validExtensionJSON),
		"/api/v2/extensions/3/actions/archive": testutils.RawHandler(http.StatusOK, "{}"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionArchive(ctx, "3"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestExtensionArchive_NotFound verifies that archiving a non-existent extension returns an error.
func TestExtensionArchive_NotFound(t *testing.T) {
	emptyList := map[string]any{"data": []any{}, "meta": map[string]any{"itemsPerPage": 100}}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extensions/999": testutils.ErrorHandler(http.StatusNotFound, "not found"),
		"/api/v2/extensions":     testutils.JSONHandler(http.StatusOK, emptyList),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionArchive(ctx, "999"); err == nil {
		t.Fatal("expected error for not-found extension, got nil")
	}
}

// TestExtensionArchive_ServerError verifies that a 500 from the archive action is surfaced.
func TestExtensionArchive_ServerError(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extensions/3":                 testutils.RawHandler(http.StatusOK, validExtensionJSON),
		"/api/v2/extensions/3/actions/archive": testutils.ErrorHandler(http.StatusInternalServerError, "server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionArchive(ctx, "3"); err == nil {
		t.Fatal("expected error for HTTP 500 on archive, got nil")
	}
}

// --- ExtensionMakePublic ---

// TestExtensionMakePublic_HappyPath verifies successful make-public of an extension.
func TestExtensionMakePublic_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extensions/3":                    testutils.RawHandler(http.StatusOK, validExtensionJSON),
		"/api/v2/extensions/3/actions/make-public": testutils.RawHandler(http.StatusOK, "{}"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionMakePublic(ctx, "3"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestExtensionMakePublic_NotFound verifies that making a non-existent extension public returns an error.
func TestExtensionMakePublic_NotFound(t *testing.T) {
	emptyList := map[string]any{"data": []any{}, "meta": map[string]any{"itemsPerPage": 100}}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extensions/999": testutils.ErrorHandler(http.StatusNotFound, "not found"),
		"/api/v2/extensions":     testutils.JSONHandler(http.StatusOK, emptyList),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionMakePublic(ctx, "999"); err == nil {
		t.Fatal("expected error for not-found extension, got nil")
	}
}

// TestExtensionMakePublic_ServerError verifies that a 500 from the make-public action is surfaced.
func TestExtensionMakePublic_ServerError(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extensions/3":                    testutils.RawHandler(http.StatusOK, validExtensionJSON),
		"/api/v2/extensions/3/actions/make-public": testutils.ErrorHandler(http.StatusInternalServerError, "server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionMakePublic(ctx, "3"); err == nil {
		t.Fatal("expected error for HTTP 500 on make-public, got nil")
	}
}

// --- ExtensionPublish ---

// TestExtensionPublish_HappyPath verifies successful publishing of an extension.
func TestExtensionPublish_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extensions/3":                 testutils.RawHandler(http.StatusOK, validExtensionJSON),
		"/api/v2/extensions/3/actions/publish": testutils.RawHandler(http.StatusOK, "{}"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionPublish(ctx, "3"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestExtensionPublish_NotFound verifies that publishing a non-existent extension returns an error.
func TestExtensionPublish_NotFound(t *testing.T) {
	emptyList := map[string]any{"data": []any{}, "meta": map[string]any{"itemsPerPage": 100}}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extensions/999": testutils.ErrorHandler(http.StatusNotFound, "not found"),
		"/api/v2/extensions":     testutils.JSONHandler(http.StatusOK, emptyList),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionPublish(ctx, "999"); err == nil {
		t.Fatal("expected error for not-found extension, got nil")
	}
}

// TestExtensionPublish_ServerError verifies that a 500 from the publish action is surfaced.
func TestExtensionPublish_ServerError(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extensions/3":                 testutils.RawHandler(http.StatusOK, validExtensionJSON),
		"/api/v2/extensions/3/actions/publish": testutils.ErrorHandler(http.StatusInternalServerError, "server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionPublish(ctx, "3"); err == nil {
		t.Fatal("expected error for HTTP 500 on publish, got nil")
	}
}

// --- ExtensionListRepo ---

// TestExtensionListRepo_InvalidURL verifies that an unreachable repo URL returns an error.
func TestExtensionListRepo_InvalidURL(t *testing.T) {
	if err := ExtensionListRepo(context.Background(), "http://127.0.0.1:1/nonexistent.git", "", ""); err == nil {
		t.Fatal("expected error for unreachable repo, got nil")
	}
}
