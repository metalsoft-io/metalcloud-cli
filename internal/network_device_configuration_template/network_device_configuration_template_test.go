package network_device_configuration_template

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

const bgpTmplItem = `{
	"id":1,
	"action":"add-neighbor",
	"networkType":"underlay",
	"networkDeviceDriver":"sonic_enterprise",
	"networkDevicePosition":"leaf",
	"remoteNetworkDevicePosition":"spine",
	"bgpNumbering":"numbered",
	"bgpLinkConfiguration":"active",
	"executionType":"cli",
	"libraryLabel":"my-template",
	"configuration":"Y21k",
	"createdTimestamp":"2024-01-01T00:00:00Z",
	"updatedTimestamp":"2024-01-01T00:00:00Z"
}`

func bgpTmplListHandler(statusCode int, items []string, currentPage, totalPages int) http.HandlerFunc {
	data := "[" + strings.Join(items, ",") + "]"
	body := fmt.Sprintf(`{"data":%s,"meta":{"currentPage":%d,"totalPages":%d,"itemsPerPage":100}}`,
		data, currentPage, totalPages)
	return testutils.RawHandler(statusCode, body)
}

// --- List ---

func TestNetworkDeviceConfigurationTemplateList_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates": bgpTmplListHandler(http.StatusOK, []string{bgpTmplItem, bgpTmplItem}, 1, 1),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateList(ctx, nil, nil); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceConfigurationTemplateList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateList(ctx, nil, nil); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

func TestNetworkDeviceConfigurationTemplateList_Empty(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates": bgpTmplListHandler(http.StatusOK, []string{}, 1, 1),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateList(ctx, nil, nil); err != nil {
		t.Errorf("expected nil error for empty list, got: %v", err)
	}
}

func TestNetworkDeviceConfigurationTemplateList_Pagination(t *testing.T) {
	makeItems := func(n int) []string {
		items := make([]string, n)
		for i := range items {
			items[i] = fmt.Sprintf(`{
				"id":%d,"action":"add-neighbor","networkType":"underlay",
				"networkDeviceDriver":"sonic_enterprise","networkDevicePosition":"leaf",
				"remoteNetworkDevicePosition":"spine","bgpNumbering":"numbered",
				"bgpLinkConfiguration":"active","executionType":"cli",
				"libraryLabel":"label-%d","configuration":"Y21k",
				"createdTimestamp":"2024-01-01T00:00:00Z","updatedTimestamp":"2024-01-01T00:00:00Z"
			}`, i+1, i+1)
		}
		return items
	}

	ts := testutils.MultiPageServer("/api/v2/network-device-bgp-configuration-templates", []any{
		decodeItems(makeItems(100)),
		decodeItems(makeItems(100)),
		decodeItems(makeItems(5)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateList(ctx, nil, nil); err != nil {
		t.Errorf("expected nil error across 3 pages, got: %v", err)
	}
}

// --- Get ---

func TestNetworkDeviceConfigurationTemplateGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates/1": testutils.RawHandler(http.StatusOK, bgpTmplItem),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateGet(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceConfigurationTemplateGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateGet(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestNetworkDeviceConfigurationTemplateGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateGet(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

// --- Create ---

func TestNetworkDeviceConfigurationTemplateCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, bgpTmplItem)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{
		"action":"add-neighbor","networkType":"underlay",
		"networkDeviceDriver":"sonic_enterprise","networkDevicePosition":"leaf",
		"remoteNetworkDevicePosition":"spine","bgpNumbering":"numbered",
		"bgpLinkConfiguration":"active","executionType":"cli",
		"libraryLabel":"my-template","configuration":"Y21k"
	}`)
	if err := NetworkDeviceConfigurationTemplateCreate(ctx, config); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceConfigurationTemplateCreate_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates": testutils.ErrorHandler(http.StatusBadRequest, "validation error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateCreate(ctx, []byte(`{}`)); err == nil {
		t.Error("expected error for 400, got nil")
	}
}

// --- Update ---

func TestNetworkDeviceConfigurationTemplateUpdate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates/1": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPatch {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, bgpTmplItem)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateUpdate(ctx, "1", []byte(`{"libraryLabel":"updated"}`)); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceConfigurationTemplateUpdate_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateUpdate(ctx, "99", []byte(`{}`)); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

// --- Delete ---

func TestNetworkDeviceConfigurationTemplateDelete_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates/1": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateDelete(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceConfigurationTemplateDelete_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-bgp-configuration-templates/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateDelete(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestNetworkDeviceConfigurationTemplateDelete_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceConfigurationTemplateDelete(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}
