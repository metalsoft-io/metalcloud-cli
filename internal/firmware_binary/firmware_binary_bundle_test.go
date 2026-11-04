package firmware_binary

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/spf13/viper"
)

func bundleFixture() map[string]any {
	bundle := binaryFixture(4796)
	bundle["name"] = "NVIDIA HGX H100 8-GPU 25.06.1"
	bundle["isBundle"] = true
	bundle["componentClass"] = nil
	bundle["sha256"] = strings.Repeat("a", 64)
	bundle["bundleSteps"] = []map[string]any{
		{"id": "hmc", "componentClass": "hgx_hmc", "targetPattern": "^HGX_FW_BMC_0$", "packageFile": "main.fwpkg",
			"packageVersion": "25.06.1", "sha256": strings.Repeat("a", 64), "activation": "hgx_reset", "dependsOnStep": nil},
		{"id": "gpu", "componentClass": "gpu", "targetPattern": "^HGX_FW_GPU_SXM_\\d+$", "packageFile": "main.fwpkg",
			"packageVersion": "25.06.1", "verifyVersion": "96.00.A5.00.01", "sha256": strings.Repeat("a", 64),
			"activation": "chassis_power_cycle", "dependsOnStep": "hmc"},
	}
	bundle["prerequisites"] = []map[string]any{
		{"platform": "NVIDIA DGXH100", "componentClass": "bios", "minimumVersion": "01.02.05"},
	}
	return bundle
}

func withTextFormat(t *testing.T) {
	t.Helper()
	viper.Set(formatter.ConfigFormat, "text")
	t.Cleanup(func() { viper.Set(formatter.ConfigFormat, "json") })
}

func TestFirmwareBinaryGet_ListsBundleStepsAndPrerequisites(t *testing.T) {
	withTextFormat(t)
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/binary/4796": testutils.JSONHandler(http.StatusOK, bundleFixture()),
	})
	defer srv.Close()

	out := testutils.CaptureStdout(t, func() {
		if err := FirmwareBinaryGet(testutils.SetupTestContext(srv.URL), "4796"); err != nil {
			t.Fatalf("FirmwareBinaryGet() unexpected error: %v", err)
		}
	})

	for _, expected := range []string{"Bundle steps (2):", "chassis_power_cycle", "96.00.A5.00.01", "Prerequisites (1):", "01.02.05"} {
		if !strings.Contains(out, expected) {
			t.Errorf("expected %q in the output:\n%s", expected, out)
		}
	}
}

func TestFirmwareBinaryGet_JsonHasNoExtraTables(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/binary/4796": testutils.JSONHandler(http.StatusOK, bundleFixture()),
	})
	defer srv.Close()

	out := testutils.CaptureStdout(t, func() {
		if err := FirmwareBinaryGet(testutils.SetupTestContext(srv.URL), "4796"); err != nil {
			t.Fatal(err)
		}
	})

	var decoded map[string]any
	if err := json.Unmarshal([]byte(out), &decoded); err != nil {
		t.Fatalf("expected one JSON object, got %q: %v", out, err)
	}
	if steps, _ := decoded["bundleSteps"].([]any); len(steps) != 2 {
		t.Errorf("bundleSteps = %v", decoded["bundleSteps"])
	}
}

func TestFirmwareBinaryUpdate_ChangesOnlyTheGivenFields(t *testing.T) {
	var sent map[string]any
	var method string
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/binary/4796": func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPut {
				method = r.Method
				body, _ := io.ReadAll(r.Body)
				_ = json.Unmarshal(body, &sent)
			}
			testutils.JSONHandler(http.StatusOK, bundleFixture())(w, r)
		},
	})
	defer srv.Close()

	config := `{"updateSeverity": "critical", "prerequisites": null}`
	if err := FirmwareBinaryUpdate(testutils.SetupTestContext(srv.URL), "4796", []byte(config)); err != nil {
		t.Fatalf("FirmwareBinaryUpdate() unexpected error: %v", err)
	}

	if method != http.MethodPut {
		t.Fatalf("expected a PUT, got %q", method)
	}
	if sent["updateSeverity"] != "critical" {
		t.Errorf("updateSeverity = %v", sent["updateSeverity"])
	}
	if steps, _ := sent["bundleSteps"].([]any); len(steps) != 2 {
		t.Errorf("the steps are kept, got %v", sent["bundleSteps"])
	}
	if _, present := sent["prerequisites"]; present {
		t.Errorf("a field set to null is cleared, got %v", sent["prerequisites"])
	}
	for _, readOnly := range []string{"id", "links"} {
		if _, present := sent[readOnly]; present {
			t.Errorf("the read-only %q is not sent", readOnly)
		}
	}
	if sent["catalogId"] != float64(10) || sent["name"] != "NVIDIA HGX H100 8-GPU 25.06.1" {
		t.Errorf("the other fields are kept, got %v", sent)
	}
}

func TestFirmwareBinaryUpdate_RejectsUnknownFields(t *testing.T) {
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/binary/4796": testutils.JSONHandler(http.StatusOK, bundleFixture()),
	})
	defer srv.Close()

	err := FirmwareBinaryUpdate(testutils.SetupTestContext(srv.URL), "4796", []byte(`{"severity": "critical"}`))
	if err == nil || !strings.Contains(err.Error(), `"severity"`) {
		t.Fatalf("expected the unknown field to be named, got %v", err)
	}
}

func TestFirmwareBinaryList_SendsBundleAndClassFilters(t *testing.T) {
	rec, url, closeFn := newRecorder(t, 1)
	defer closeFn()

	err := FirmwareBinaryList(testutils.SetupTestContext(url), ListFlags{
		FilterIsBundle:       []string{"true"},
		FilterComponentClass: []string{"gpu"},
	})
	if err != nil {
		t.Fatal(err)
	}

	query := rec.queries[0]
	if query.Get("filter.isBundle") != "true" || query.Get("filter.componentClass") != "gpu" {
		t.Errorf("filters = %v", query)
	}
}

func TestFirmwareBundleConfigExample(t *testing.T) {
	out := testutils.CaptureStdout(t, func() {
		if err := FirmwareBundleConfigExample(nil); err != nil {
			t.Fatal(err)
		}
	})
	for _, expected := range []string{`"isBundle":true`, `"bundleSteps":[`, `"verifyVersion":"96.00.A5.00.01"`, `"dependsOnStep":null`, `"prerequisites":[`} {
		if !strings.Contains(out, expected) {
			t.Errorf("expected %s in %s", expected, out)
		}
	}
}
