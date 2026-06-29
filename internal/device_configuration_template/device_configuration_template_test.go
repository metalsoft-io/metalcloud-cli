package device_configuration_template

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

func decodeItems(jsonStrings []string) []any {
	out := make([]any, len(jsonStrings))
	for i, s := range jsonStrings {
		var v any
		if err := json.Unmarshal([]byte(s), &v); err != nil {
			panic(fmt.Sprintf("decodeItems: invalid JSON at index %d: %v", i, err))
		}
		out[i] = v
	}
	return out
}

func TestMain(m *testing.M) {
	testutils.SetupTestFormat()
	m.Run()
}

const (
	basePath    = "/api/v2/device-configuration-templates"
	configPath  = basePath + "/config"
	profilePath = basePath + "/profile"
)

const tmplItem = `{
	"id":1,
	"label":"my-device-template",
	"name":"My device template",
	"description":"Example device configuration template",
	"deviceDriver":"junos",
	"executionType":"cli",
	"templateContent":"hostname {{ hostname }}",
	"tags":["example"],
	"revision":"1",
	"createdTimestamp":"2024-01-01T00:00:00Z",
	"updatedTimestamp":"2024-01-01T00:00:00Z"
}`

const profileItem = `{
	"id":"1",
	"deviceConfigurationTemplateId":1,
	"networkDeviceId":100,
	"lifecycleStage":"configuration",
	"isEnabled":true,
	"priority":100,
	"applyMode":"once",
	"revision":"1",
	"createdTimestamp":"2024-01-01T00:00:00Z",
	"updatedTimestamp":"2024-01-01T00:00:00Z"
}`

const renderedItem = `{"rendered":"hostname switch-01","templateContent":"hostname {{ hostname }}"}`

func listHandler(statusCode int, items []string, currentPage, totalPages int) http.HandlerFunc {
	data := "[" + strings.Join(items, ",") + "]"
	body := fmt.Sprintf(`{"data":%s,"meta":{"currentPage":%d,"totalPages":%d,"itemsPerPage":100}}`,
		data, currentPage, totalPages)
	return testutils.RawHandler(statusCode, body)
}

// methodHandler dispatches to per-method handlers, returning 405 for any method
// without an entry. Used for endpoints that the CLI hits with more than one verb
// (e.g. update/delete first GET the resource to read its revision).
func methodHandler(handlers map[string]http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h, ok := handlers[r.Method]; ok {
			h(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// --- Template List ---

func TestDeviceConfigurationTemplateList_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		configPath: listHandler(http.StatusOK, []string{tmplItem, tmplItem}, 1, 1),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateList(ctx, nil, nil, nil); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestDeviceConfigurationTemplateList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		configPath: testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateList(ctx, nil, nil, nil); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

func TestDeviceConfigurationTemplateList_Empty(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		configPath: listHandler(http.StatusOK, []string{}, 1, 1),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateList(ctx, nil, nil, nil); err != nil {
		t.Errorf("expected nil error for empty list, got: %v", err)
	}
}

func TestDeviceConfigurationTemplateList_Pagination(t *testing.T) {
	makeItems := func(n int) []string {
		items := make([]string, n)
		for i := range items {
			items[i] = fmt.Sprintf(`{
				"id":%d,"label":"label-%d","deviceDriver":"junos","executionType":"cli",
				"revision":"1","createdTimestamp":"2024-01-01T00:00:00Z",
				"updatedTimestamp":"2024-01-01T00:00:00Z"
			}`, i+1, i+1)
		}
		return items
	}

	ts := testutils.MultiPageServer(configPath, []any{
		decodeItems(makeItems(100)),
		decodeItems(makeItems(100)),
		decodeItems(makeItems(5)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateList(ctx, nil, nil, nil); err != nil {
		t.Errorf("expected nil error across 3 pages, got: %v", err)
	}
}

// --- Template Get ---

func TestDeviceConfigurationTemplateGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		configPath + "/1": testutils.RawHandler(http.StatusOK, tmplItem),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateGet(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestDeviceConfigurationTemplateGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		configPath + "/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateGet(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestDeviceConfigurationTemplateGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateGet(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

// --- Template Create ---

func TestDeviceConfigurationTemplateCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		configPath: methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: testutils.RawHandler(http.StatusCreated, tmplItem),
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"label":"my-device-template","deviceDriver":"junos","executionType":"cli","templateContent":"x"}`)
	if err := DeviceConfigurationTemplateCreate(ctx, config); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestDeviceConfigurationTemplateCreate_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		configPath: testutils.ErrorHandler(http.StatusBadRequest, "validation error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateCreate(ctx, []byte(`{}`)); err == nil {
		t.Error("expected error for 400, got nil")
	}
}

// --- Template Update ---

func TestDeviceConfigurationTemplateUpdate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		configPath + "/1": methodHandler(map[string]http.HandlerFunc{
			http.MethodGet:   testutils.RawHandler(http.StatusOK, tmplItem),
			http.MethodPatch: testutils.RawHandler(http.StatusOK, tmplItem),
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateUpdate(ctx, "1", []byte(`{"label":"updated"}`)); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestDeviceConfigurationTemplateUpdate_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		configPath + "/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateUpdate(ctx, "99", []byte(`{}`)); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

// --- Template Delete ---

func TestDeviceConfigurationTemplateDelete_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		configPath + "/1": methodHandler(map[string]http.HandlerFunc{
			http.MethodGet:    testutils.RawHandler(http.StatusOK, tmplItem),
			http.MethodDelete: testutils.RawHandler(http.StatusNoContent, ""),
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateDelete(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestDeviceConfigurationTemplateDelete_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		configPath + "/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateDelete(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestDeviceConfigurationTemplateDelete_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateDelete(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

// --- Template Render ---

func TestDeviceConfigurationTemplateRender_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		configPath + "/actions/render": methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: testutils.RawHandler(http.StatusOK, renderedItem),
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"templateContent":"hostname {{ hostname }}","variables":{"hostname":"switch-01"}}`)
	if err := DeviceConfigurationTemplateRender(ctx, config); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestDeviceConfigurationTemplateRenderSaved_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		configPath + "/1/actions/render": methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: testutils.RawHandler(http.StatusOK, renderedItem),
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateRenderSaved(ctx, "1", []byte(`{"variables":{"hostname":"switch-01"}}`)); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

// --- ConfigExample (offline) ---

func TestDeviceConfigurationTemplateConfigExample(t *testing.T) {
	if err := DeviceConfigurationTemplateConfigExample(nil); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestDeviceConfigurationTemplateProfileConfigExample(t *testing.T) {
	if err := DeviceConfigurationTemplateProfileConfigExample(nil); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

// --- Profile List ---

func TestDeviceConfigurationTemplateProfileList_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		profilePath: listHandler(http.StatusOK, []string{profileItem, profileItem}, 1, 1),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateProfileList(ctx, nil, nil, nil, nil); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestDeviceConfigurationTemplateProfileList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		profilePath: testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateProfileList(ctx, nil, nil, nil, nil); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

// --- Profile Get ---

func TestDeviceConfigurationTemplateProfileGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		profilePath + "/1": testutils.RawHandler(http.StatusOK, profileItem),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateProfileGet(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestDeviceConfigurationTemplateProfileGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateProfileGet(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

// --- Profile Create ---

func TestDeviceConfigurationTemplateProfileCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		profilePath: methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: testutils.RawHandler(http.StatusCreated, profileItem),
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"deviceConfigurationTemplateId":1,"networkDeviceId":100}`)
	if err := DeviceConfigurationTemplateProfileCreate(ctx, config); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

// --- Profile Update ---

func TestDeviceConfigurationTemplateProfileUpdate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		profilePath + "/1": methodHandler(map[string]http.HandlerFunc{
			http.MethodGet:   testutils.RawHandler(http.StatusOK, profileItem),
			http.MethodPatch: testutils.RawHandler(http.StatusOK, profileItem),
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateProfileUpdate(ctx, "1", []byte(`{"isEnabled":false}`)); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

// --- Profile Delete ---

func TestDeviceConfigurationTemplateProfileDelete_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		profilePath + "/1": methodHandler(map[string]http.HandlerFunc{
			http.MethodGet:    testutils.RawHandler(http.StatusOK, profileItem),
			http.MethodDelete: testutils.RawHandler(http.StatusNoContent, ""),
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateProfileDelete(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

// --- Profile Render ---

func TestDeviceConfigurationTemplateProfileRender_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		profilePath + "/1/actions/render": methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: testutils.RawHandler(http.StatusOK, renderedItem),
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateProfileRender(ctx, "1", []byte(`{"networkDeviceId":100}`)); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

// --- Profile Find / Render Applicable ---

func TestDeviceConfigurationTemplateProfileFindApplicable_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		profilePath + "/actions/find-applicable": methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: testutils.RawHandler(http.StatusOK, `{"items":[],"alreadyApplied":[]}`),
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateProfileFindApplicable(ctx, []byte(`{"networkDeviceId":100}`)); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestDeviceConfigurationTemplateProfileRenderApplicable_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		profilePath + "/actions/render-applicable": methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: testutils.RawHandler(http.StatusOK, `{"items":[],"alreadyAppliedItems":[],"joined":""}`),
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := DeviceConfigurationTemplateProfileRenderApplicable(ctx, []byte(`{"networkDeviceId":100}`)); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

// --- Profile Bulk Assign ---

func TestDeviceConfigurationTemplateProfileBulkAssign_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		profilePath + "/bulk": methodHandler(map[string]http.HandlerFunc{
			http.MethodPost: testutils.RawHandler(http.StatusOK, `{"created":[],"skipped":[],"targetDeviceCount":0}`),
		}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"deviceConfigurationTemplateId":1,"networkDeviceIds":[100,101]}`)
	if err := DeviceConfigurationTemplateProfileBulkAssign(ctx, config); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}
