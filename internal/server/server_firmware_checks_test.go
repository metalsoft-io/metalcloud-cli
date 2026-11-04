package server

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/spf13/viper"
)

func withFirmwareTextFormat(t *testing.T) {
	t.Helper()
	viper.Set(formatter.ConfigFormat, "text")
	t.Cleanup(func() { viper.Set(formatter.ConfigFormat, "json") })
}

func assertContainsAll(t *testing.T, out string, expected ...string) {
	t.Helper()
	for _, text := range expected {
		if !strings.Contains(out, text) {
			t.Errorf("expected %q in the output:\n%s", text, out)
		}
	}
}

const refusedUpgradeJSON = `{
	"statusCode": 400,
	"message": "The firmware upgrade was refused by its safety checks",
	"reasonCode": "all_components_refused",
	"refusals": [
		{"componentId": 2320, "componentName": "HGX_FW_GPU_SXM_3", "reasonCode": "component_not_healthy", "message": "Health is Critical"},
		{"reasonCode": "concurrent_firmware_job", "message": "Another firmware job is running"}
	]
}`

func TestServerFirmwareUpgrade_RefusalIsReadable(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/7/firmware/actions/upgrade": testutils.RawHandler(http.StatusBadRequest, refusedUpgradeJSON),
	})
	defer ts.Close()

	err := ServerFirmwareUpgrade(testutils.SetupTestContext(ts.URL), "7")
	if err == nil {
		t.Fatal("expected the refusal to be an error")
	}
	assertContainsAll(t, err.Error(),
		"refused by its safety checks",
		"HGX_FW_GPU_SXM_3 (2320): component_not_healthy: Health is Critical",
		"  concurrent_firmware_job: Another firmware job is running",
	)
	if strings.Contains(err.Error(), `"statusCode"`) {
		t.Errorf("expected no raw body, got %s", err)
	}
}

func TestServerFirmwareUpgrade_OtherErrorsKeepTheRawBody(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/7/firmware/actions/upgrade": testutils.RawHandler(http.StatusBadRequest, `{"statusCode":400,"message":"Server is not managed"}`),
	})
	defer ts.Close()

	err := ServerFirmwareUpgrade(testutils.SetupTestContext(ts.URL), "7")
	if err == nil || !strings.Contains(err.Error(), "Server is not managed") {
		t.Fatalf("expected the response body in the error, got %v", err)
	}
}

func TestServerFirmwareUpgrade_ListsLeftOutComponents(t *testing.T) {
	withFirmwareTextFormat(t)
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/7/firmware/actions/upgrade": testutils.RawHandler(http.StatusAccepted, `{
			"jobId": 991, "jobGroupId": 55,
			"refusedComponents": [{"componentId": 2320, "componentName": "HGX_FW_GPU_SXM_3", "reasonCode": "component_not_healthy", "message": "Health is Critical"}]
		}`),
	})
	defer ts.Close()

	out := testutils.CaptureStdout(t, func() {
		if err := ServerFirmwareUpgrade(testutils.SetupTestContext(ts.URL), "7"); err != nil {
			t.Fatal(err)
		}
	})
	assertContainsAll(t, out, "991", "Left out by the safety checks (1):", "HGX_FW_GPU_SXM_3", "component_not_healthy")
}

func TestServerFirmwareUpgrade_JsonKeepsTheRefusedComponents(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/7/firmware/actions/upgrade": testutils.RawHandler(http.StatusAccepted, `{
			"jobId": 991,
			"refusedComponents": [{"componentId": 2320, "reasonCode": "component_not_healthy", "message": "Health is Critical"}]
		}`),
	})
	defer ts.Close()

	out := testutils.CaptureStdout(t, func() {
		if err := ServerFirmwareUpgrade(testutils.SetupTestContext(ts.URL), "7"); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, `"refusedComponents":[`) || strings.Contains(out, "Left out") {
		t.Errorf("expected the JSON job info only, got %s", out)
	}
}

func serverComponentJSON(id int64, name, class, target string) map[string]any {
	component := map[string]any{
		"id": id, "name": name, "componentClass": class, "componentVendor": "NVIDIA",
		"firmwareUpdateable": 1, "firmwareUpgradeNeedsConfirmation": 0, "firmwareStatus": "idle",
		"firmwareVersion": "96.00.74.00.01",
	}
	if target != "" {
		component["firmwareTargetVersion"] = target
	}
	if class == "gpu" {
		component["baseboardId"] = "HGX_Baseboard_0"
	}
	return component
}

func TestServerFirmwareComponentUpdate_ListsTheOtherTargetsThePackageChanged(t *testing.T) {
	withFirmwareTextFormat(t)
	var mu sync.Mutex
	updated := false
	components := func() []any {
		mu.Lock()
		defer mu.Unlock()
		target := ""
		if updated {
			target = "96.00.A5.00.01"
		}
		return []any{
			serverComponentJSON(2318, "HGX_FW_GPU_SXM_1", "gpu", target),
			serverComponentJSON(2319, "HGX_FW_GPU_SXM_2", "gpu", target),
			serverComponentJSON(2001, "BIOS", "bios", ""),
		}
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/7/firmware/components": func(w http.ResponseWriter, r *http.Request) {
			testutils.JSONHandler(http.StatusOK, testutils.PaginatedResponse(components(), 1, 1))(w, r)
		},
		"/api/v2/servers/7/firmware/components/2318": func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			updated = true
			mu.Unlock()
			testutils.JSONHandler(http.StatusOK, serverComponentJSON(2318, "HGX_FW_GPU_SXM_1", "gpu", "96.00.A5.00.01"))(w, r)
		},
	})
	defer ts.Close()

	out := testutils.CaptureStdout(t, func() {
		err := ServerFirmwareComponentUpdate(testutils.SetupTestContext(ts.URL), "7", "2318", []byte(`{"firmwareTargetVersion":"96.00.A5.00.01"}`))
		if err != nil {
			t.Fatal(err)
		}
	})
	assertContainsAll(t, out, "so their target changed too (1):", "HGX_FW_GPU_SXM_2")
	if strings.Contains(out, "BIOS") {
		t.Errorf("an unchanged component is not listed:\n%s", out)
	}
}

func TestServerFirmwareComponentUpdate_NoTargetNoExtraRequests(t *testing.T) {
	withFirmwareTextFormat(t)
	listed := false
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/7/firmware/components": func(w http.ResponseWriter, r *http.Request) {
			listed = true
			testutils.JSONHandler(http.StatusOK, testutils.PaginatedResponse([]any{}, 1, 1))(w, r)
		},
		"/api/v2/servers/7/firmware/components/2001": testutils.JSONHandler(http.StatusOK, serverComponentJSON(2001, "BIOS", "bios", "")),
	})
	defer ts.Close()

	testutils.CaptureStdout(t, func() {
		if err := ServerFirmwareComponentUpdate(testutils.SetupTestContext(ts.URL), "7", "2001", []byte(`{"firmwareUpdateable":1}`)); err != nil {
			t.Fatal(err)
		}
	})
	if listed {
		t.Error("the components are only compared when the target version changes")
	}
}

func TestServerFirmwareComponentsListFiltered_SendsTheFilters(t *testing.T) {
	withFirmwareTextFormat(t)
	rec := &hwRecorder{}
	var query string
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/7/firmware/components": func(w http.ResponseWriter, r *http.Request) {
			query = r.URL.RawQuery
			testutils.JSONHandler(http.StatusOK, testutils.PaginatedResponse([]any{serverComponentJSON(2318, "HGX_FW_GPU_SXM_1", "gpu", "")}, 1, 1))(w, r)
		},
	})
	defer ts.Close()

	out := testutils.CaptureStdout(t, func() {
		err := ServerFirmwareComponentsListFiltered(testutils.SetupTestContext(ts.URL), "7", ComponentsListFlags{
			FilterComponentClass: []string{"gpu"},
			FilterBaseboardId:    []string{"HGX_Baseboard_0"},
		})
		if err != nil {
			t.Fatal(err)
		}
	})
	assertContainsAll(t, query, "filter.componentClass=gpu", "filter.baseboardId=HGX_Baseboard_0")
	assertContainsAll(t, out, "CLASS", "BASEBOARD", "HGX_Baseboard_0", "NVIDIA")
}

const preflightAuditJSON = `[{
	"serverId": 7,
	"baseline": {"id": 5, "name": "dgx-h100"},
	"components": [
		{"componentId": 2318, "name": "HGX_FW_GPU_SXM_1", "status": "upgrade_available", "currentVersion": "96.00.74.00.01",
			"targetVersion": "96.00.A5.00.01", "upgrades": [], "downgrades": [], "componentClass": "gpu", "baseboardId": "HGX_Baseboard_0",
			"checks": [{"check": "component_health", "status": "passed"}]},
		{"componentId": 2320, "name": "HGX_FW_GPU_SXM_3", "status": "upgrade_available", "currentVersion": "96.00.74.00.01",
			"targetVersion": "96.00.A5.00.01", "upgrades": [], "downgrades": [], "componentClass": "gpu", "baseboardId": "HGX_Baseboard_0",
			"checks": [{"check": "component_health", "status": "refused", "reasonCode": "component_not_healthy", "message": "Health is Critical"}]}
	],
	"baseboards": [{"baseboardId": "HGX_Baseboard_0", "componentIds": [2318, 2320], "componentClasses": ["gpu"]}],
	"preflight": {
		"checkedAt": "2026-11-04T10:00:00Z",
		"jobRefused": false,
		"checks": [{"check": "concurrent_job", "status": "passed"}],
		"refusals": [{"componentId": 2320, "componentName": "HGX_FW_GPU_SXM_3", "reasonCode": "component_not_healthy", "message": "Health is Critical"}],
		"warnings": ["The prerequisites of bundle 4796 could not be read"],
		"plannedJobs": [{
			"componentIds": [2318, 2319], "componentNames": ["HGX_FW_GPU_SXM_1", "HGX_FW_GPU_SXM_2"], "targetVersion": "25.06.1",
			"bundleId": "4796", "bundleName": "NVIDIA HGX H100 8-GPU", "stepIndex": 2, "stepCount": 3, "stepIds": ["gpu"],
			"activations": ["chassis_power_cycle"]
		}]
	}
}]`

func TestServerFirmwareGenerateAuditWithPreflight(t *testing.T) {
	withFirmwareTextFormat(t)
	rec := &hwRecorder{}
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/firmware/actions/generate-upgrade-audit": testutils.RawHandler(http.StatusOK, preflightAuditJSON),
	})
	defer ts.Close()

	out := testutils.CaptureStdout(t, func() {
		if err := ServerFirmwareGenerateAuditWithPreflight(testutils.SetupTestContext(ts.URL), []byte(`{"serverIds":[7]}`), true); err != nil {
			t.Fatal(err)
		}
	})

	var sent map[string]any
	if err := json.Unmarshal([]byte(rec.last().Body), &sent); err != nil || sent["includePreflight"] != true {
		t.Errorf("expected includePreflight to be sent, got %q", rec.last().Body)
	}
	assertContainsAll(t, out,
		"refused: component_health",
		"passed",
		"Server 7 GPU baseboards:",
		"HGX_Baseboard_0",
		"Server 7: an upgrade started now would run 1 job(s)",
		"concurrent_job",
		"component_not_healthy",
		"Not checked (these do not refuse the upgrade):",
		"The prerequisites of bundle 4796 could not be read",
		"Jobs, in the order they run:",
		"NVIDIA HGX H100 8-GPU step 2/3 (gpu)",
		"HGX_FW_GPU_SXM_1..2",
		"chassis_power_cycle",
	)
}

func TestServerFirmwareGenerateAudit_WithoutPreflightSendsNoFlag(t *testing.T) {
	rec := &hwRecorder{}
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/firmware/actions/generate-upgrade-audit": testutils.RawHandler(http.StatusOK, `[]`),
	})
	defer ts.Close()

	testutils.CaptureStdout(t, func() {
		if err := ServerFirmwareGenerateAuditWithPreflight(testutils.SetupTestContext(ts.URL), []byte(`{"serverIds":[7]}`), false); err != nil {
			t.Fatal(err)
		}
	})
	if strings.Contains(rec.last().Body, "includePreflight") {
		t.Errorf("expected no includePreflight, got %q", rec.last().Body)
	}
}

const searchDryRunJSON = `{
	"baseline": {"id": 5, "name": "dgx-h100"},
	"componentBinaries": {},
	"matchReport": {
		"2001": {"componentId": 2001, "componentName": "BIOS", "vendor": "nvidia", "identity": {}, "matched": [], "rejected": [],
			"note": "No binary of the searched catalogs is for this component"},
		"2318": {"componentId": 2318, "componentName": "HGX_FW_GPU_SXM_1", "vendor": "nvidia", "identity": {"class": "gpu"},
			"matched": [{"binaryId": 4796, "catalogId": 10, "name": "NVIDIA HGX H100 8-GPU", "packageVersion": "96.00.A5.00.01",
				"bundleVersion": "25.06.1", "rules": ["supportedSystem", "componentClass", "bundleStep"]}],
			"rejected": [{"binaryId": 4700, "catalogId": 10, "name": "NVIDIA HGX H100 8-GPU", "packageVersion": "96.00.74.00.01",
				"rules": ["supportedSystem"], "reason": "baseboard: the step targets another baseboard"}]}
	}
}`

func TestServerFirmwareSearchDryRun(t *testing.T) {
	withFirmwareTextFormat(t)
	rec := &hwRecorder{}
	ts := hwServer(t, rec, map[string]http.HandlerFunc{
		"/api/v2/servers/7/firmware/actions/search-dry-run": testutils.RawHandler(http.StatusOK, searchDryRunJSON),
	})
	defer ts.Close()

	out := testutils.CaptureStdout(t, func() {
		err := ServerFirmwareSearchDryRun(testutils.SetupTestContext(ts.URL), "7", ServerFirmwareSearchDryRunOptions{
			CatalogIds:         []int64{10, 11},
			IncludeAllVersions: true,
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	var sent map[string]any
	if err := json.Unmarshal([]byte(rec.last().Body), &sent); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(sent["catalogIds"], []any{float64(10), float64(11)}) || sent["includeAllVersions"] != true || sent["baselineId"] != nil {
		t.Errorf("request = %v", sent)
	}
	assertContainsAll(t, out,
		"Baseline: dgx-h100",
		"matched", "4796", "25.06.1", "supportedSystem, componentClass, bundleStep",
		"rejected", "baseboard: the step targets another baseboard",
		"none", "No binary of the searched catalogs",
	)
	if strings.Index(out, "HGX_FW_GPU_SXM_1") > strings.Index(out, "BIOS") {
		t.Errorf("components with a match are listed first:\n%s", out)
	}
}

func TestServerFirmwareSearchDryRun_BaselineOrCatalogs(t *testing.T) {
	err := ServerFirmwareSearchDryRun(testutils.SetupTestContext("http://127.0.0.1:1"), "7", ServerFirmwareSearchDryRunOptions{
		BaselineId: 5,
		CatalogIds: []int64{10},
	})
	if err == nil || !strings.Contains(err.Error(), "not both") {
		t.Fatalf("expected an error, got %v", err)
	}
}

func TestCollapseComponentNames(t *testing.T) {
	got := collapseComponentNames([]string{
		"HGX_FW_GPU_SXM_1", "HGX_FW_GPU_SXM_2", "HGX_FW_GPU_SXM_3", "HGX_FW_GPU_SXM_5", "HGX_FW_BMC_0", "BIOS",
	})
	want := []string{"HGX_FW_GPU_SXM_1..3", "HGX_FW_GPU_SXM_5", "HGX_FW_BMC_0", "BIOS"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("collapseComponentNames() = %v, want %v", got, want)
	}
}
