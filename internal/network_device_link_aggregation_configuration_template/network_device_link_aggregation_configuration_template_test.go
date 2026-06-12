package network_device_link_aggregation_configuration_template

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

const lagTmplItem = `{
	"id":1,
	"action":"create",
	"aggregationType":"lag",
	"networkDeviceDriver":"sonic_enterprise",
	"executionType":"cli",
	"libraryLabel":"my-lag-template",
	"configuration":"Y21k",
	"createdTimestamp":"2024-01-01T00:00:00Z",
	"updatedTimestamp":"2024-01-01T00:00:00Z"
}`

func lagTmplListHandler(statusCode int, items []string, currentPage, totalPages int) http.HandlerFunc {
	data := "[" + strings.Join(items, ",") + "]"
	body := fmt.Sprintf(`{"data":%s,"meta":{"currentPage":%d,"totalPages":%d,"itemsPerPage":100}}`,
		data, currentPage, totalPages)
	return testutils.RawHandler(statusCode, body)
}

// --- List ---

func TestNetworkDeviceLinkAggregationConfigurationTemplateList_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-link-aggregation-configuration-templates": lagTmplListHandler(http.StatusOK, []string{lagTmplItem, lagTmplItem}, 1, 1),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceLinkAggregationConfigurationTemplateList(ctx, nil, nil); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceLinkAggregationConfigurationTemplateList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-link-aggregation-configuration-templates": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceLinkAggregationConfigurationTemplateList(ctx, nil, nil); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

func TestNetworkDeviceLinkAggregationConfigurationTemplateList_Empty(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-link-aggregation-configuration-templates": lagTmplListHandler(http.StatusOK, []string{}, 1, 1),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceLinkAggregationConfigurationTemplateList(ctx, nil, nil); err != nil {
		t.Errorf("expected nil error for empty list, got: %v", err)
	}
}

func TestNetworkDeviceLinkAggregationConfigurationTemplateList_Pagination(t *testing.T) {
	makeItems := func(n int) []string {
		items := make([]string, n)
		for i := range items {
			items[i] = fmt.Sprintf(`{
				"id":%d,"action":"create","aggregationType":"lag",
				"networkDeviceDriver":"sonic_enterprise","executionType":"cli",
				"libraryLabel":"label-%d","configuration":"Y21k",
				"createdTimestamp":"2024-01-01T00:00:00Z","updatedTimestamp":"2024-01-01T00:00:00Z"
			}`, i+1, i+1)
		}
		return items
	}

	ts := testutils.MultiPageServer("/api/v2/network-device-link-aggregation-configuration-templates", []any{
		decodeItems(makeItems(100)),
		decodeItems(makeItems(100)),
		decodeItems(makeItems(5)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceLinkAggregationConfigurationTemplateList(ctx, nil, nil); err != nil {
		t.Errorf("expected nil error across 3 pages, got: %v", err)
	}
}

// --- Get ---

func TestNetworkDeviceLinkAggregationConfigurationTemplateGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-link-aggregation-configuration-templates/1": testutils.RawHandler(http.StatusOK, lagTmplItem),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceLinkAggregationConfigurationTemplateGet(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceLinkAggregationConfigurationTemplateGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-link-aggregation-configuration-templates/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceLinkAggregationConfigurationTemplateGet(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestNetworkDeviceLinkAggregationConfigurationTemplateGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceLinkAggregationConfigurationTemplateGet(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

// --- Create ---

func TestNetworkDeviceLinkAggregationConfigurationTemplateCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-link-aggregation-configuration-templates": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, lagTmplItem)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{
		"action":"create","aggregationType":"lag",
		"networkDeviceDriver":"sonic_enterprise","executionType":"cli",
		"libraryLabel":"my-lag-template","configuration":"Y21k"
	}`)
	if err := NetworkDeviceLinkAggregationConfigurationTemplateCreate(ctx, config); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceLinkAggregationConfigurationTemplateCreate_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-link-aggregation-configuration-templates": testutils.ErrorHandler(http.StatusBadRequest, "validation error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceLinkAggregationConfigurationTemplateCreate(ctx, []byte(`{}`)); err == nil {
		t.Error("expected error for 400, got nil")
	}
}

// --- Update ---

func TestNetworkDeviceLinkAggregationConfigurationTemplateUpdate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-link-aggregation-configuration-templates/1": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPatch {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, lagTmplItem)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceLinkAggregationConfigurationTemplateUpdate(ctx, "1", []byte(`{"libraryLabel":"updated"}`)); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceLinkAggregationConfigurationTemplateUpdate_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-link-aggregation-configuration-templates/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceLinkAggregationConfigurationTemplateUpdate(ctx, "99", []byte(`{}`)); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

// --- Delete ---

func TestNetworkDeviceLinkAggregationConfigurationTemplateDelete_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-link-aggregation-configuration-templates/1": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceLinkAggregationConfigurationTemplateDelete(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceLinkAggregationConfigurationTemplateDelete_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-device-link-aggregation-configuration-templates/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceLinkAggregationConfigurationTemplateDelete(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestNetworkDeviceLinkAggregationConfigurationTemplateDelete_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceLinkAggregationConfigurationTemplateDelete(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}
