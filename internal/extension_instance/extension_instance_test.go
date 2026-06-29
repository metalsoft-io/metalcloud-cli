package extension_instance

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

// infrastructureSearchResponse returns the envelope that GetInfrastructures (search) returns.
// All required fields for sdk.Infrastructure must be present to avoid SDK validation errors.
func infrastructureSearchResponse(infraId float32, label string) map[string]any {
	return map[string]any{
		"data": []any{
			map[string]any{
				"id":               infraId,
				"label":            label,
				"revision":         float32(1),
				"serviceStatus":    "active",
				"datacenterName":   "dc1",
				"siteId":           1,
				"createdTimestamp": "2024-01-01T00:00:00Z",
				"updatedTimestamp": "2024-01-01T00:00:00Z",
				"designIsLocked":   0,
				"config":           map[string]any{},
			},
		},
		"meta": map[string]any{"itemsPerPage": 100},
	}
}

func makeExtensionInstance(id float32) map[string]any {
	return map[string]any{
		"id":                  id,
		"revision":            float32(1),
		"label":               "ext-instance",
		"automaticManagement": float32(1),
		"infrastructureId":    float32(10),
		"infrastructure": map[string]any{
			"id":    float32(10),
			"label": "infra-1",
		},
		"extensionId":      float32(5),
		"serviceStatus":    "active",
		"updatedTimestamp": "2024-01-01T00:00:00Z",
		"inputVariables":   []any{},
		"outputVariables":  []any{},
		"config": map[string]any{
			"revision":            float32(1),
			"label":               "ext-instance",
			"automaticManagement": float32(1),
			"deployType":          "none",
			"deployStatus":        "idle",
			"updatedTimestamp":    "2024-01-01T00:00:00Z",
		},
	}
}

// TestExtensionInstanceList_HappyPath verifies listing instances for an infrastructure.
func TestExtensionInstanceList_HappyPath(t *testing.T) {
	instances := testutils.PaginatedResponse([]any{makeExtensionInstance(1), makeExtensionInstance(2)}, 1, 1)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                   testutils.JSONHandler(http.StatusOK, infrastructureSearchResponse(10, "infra-1")),
		"/api/v2/infrastructures/10/extension-instances": testutils.JSONHandler(http.StatusOK, instances),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionInstanceList(ctx, "infra-1", nil, nil, nil, nil); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestExtensionInstanceList_InfraNotFound verifies infrastructure lookup failure is surfaced.
func TestExtensionInstanceList_InfraNotFound(t *testing.T) {
	emptyList := map[string]any{"data": []any{}, "meta": map[string]any{"itemsPerPage": 100}}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures": testutils.JSONHandler(http.StatusOK, emptyList),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionInstanceList(ctx, "no-such-infra", nil, nil, nil, nil); err == nil {
		t.Fatal("expected an error for missing infrastructure, got nil")
	}
}

// TestExtensionInstanceList_ServerError verifies a 500 from the instances endpoint is surfaced.
func TestExtensionInstanceList_ServerError(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                   testutils.JSONHandler(http.StatusOK, infrastructureSearchResponse(10, "infra-1")),
		"/api/v2/infrastructures/10/extension-instances": testutils.ErrorHandler(http.StatusInternalServerError, "server error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionInstanceList(ctx, "infra-1", nil, nil, nil, nil); err == nil {
		t.Fatal("expected an error for HTTP 500, got nil")
	}
}

// TestExtensionInstanceGet_HappyPath verifies successful retrieval of a single instance.
func TestExtensionInstanceGet_HappyPath(t *testing.T) {
	inst := makeExtensionInstance(42)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extension-instances/42": testutils.JSONHandler(http.StatusOK, inst),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionInstanceGet(ctx, "42"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestExtensionInstanceGet_NotFound verifies a 404 is surfaced as an error.
func TestExtensionInstanceGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extension-instances/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionInstanceGet(ctx, "99"); err == nil {
		t.Fatal("expected an error for HTTP 404, got nil")
	}
}

// TestExtensionInstanceGet_InvalidId verifies a non-numeric ID is rejected immediately.
func TestExtensionInstanceGet_InvalidId(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := ExtensionInstanceGet(ctx, "bad-id"); err == nil {
		t.Fatal("expected an error for invalid ID, got nil")
	}
}

// TestExtensionInstanceCreate_HappyPath verifies successful creation of an instance.
func TestExtensionInstanceCreate_HappyPath(t *testing.T) {
	inst := makeExtensionInstance(1)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/infrastructures":                   testutils.JSONHandler(http.StatusOK, infrastructureSearchResponse(10, "infra-1")),
		"/api/v2/infrastructures/10/extension-instances": testutils.JSONHandler(http.StatusCreated, inst),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	payload := sdk.CreateExtensionInstance{
		ExtensionId: sdk.PtrInt64(5),
		Label:       sdk.PtrString("new-instance"),
	}
	if err := ExtensionInstanceCreate(ctx, "infra-1", payload); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

// TestExtensionInstanceDelete_HappyPath verifies successful deletion of an instance.
func TestExtensionInstanceDelete_HappyPath(t *testing.T) {
	inst := makeExtensionInstance(7)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/extension-instances/7": func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				testutils.JSONHandler(http.StatusOK, inst)(w, r)
			case http.MethodDelete:
				w.WriteHeader(http.StatusNoContent)
			default:
				http.NotFound(w, r)
			}
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := ExtensionInstanceDelete(ctx, "7"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}
