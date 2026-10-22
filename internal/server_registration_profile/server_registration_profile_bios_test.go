package server_registration_profile

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/spf13/viper"
)

const registrationProfileWithComplianceItem = `{
	"id": 7, "name": "compliance-profile", "revision": "1",
	"isDefault": false,
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
	"settings": {
		"registerCredentials": "default",
		"minimumNumberOfConnectedInterfaces": 1,
		"defaultVirtualMediaProtocol": "ipmi",
		"firmwareBaselineId": 1902,
		"biosProfile": [
			{"attributeName": "ProcCStates", "attributeValue": "Disabled", "serverVendor": "Dell", "serverModel": "R640"},
			{"attributeName": "SriovGlobalEnable", "attributeValue": "Enabled", "serverVendor": "HPE"}
		]
	},
	"links": []
}`

// renderCSV runs fn with CSV output and returns the first row keyed by column title.
func renderCSV(t *testing.T, fn func()) map[string]string {
	t.Helper()

	viper.Set(formatter.ConfigFormat, "csv")
	defer viper.Set(formatter.ConfigFormat, "json")

	out := testutils.CaptureStdout(t, fn)

	records, err := csv.NewReader(strings.NewReader(out)).ReadAll()
	if err != nil {
		t.Fatalf("output is not valid CSV: %v\n%s", err, out)
	}
	if len(records) < 2 {
		t.Fatalf("expected a header and a row, got %q", out)
	}

	row := map[string]string{}
	for i, name := range records[0] {
		row[name] = records[1][i]
	}
	return row
}

func assertComplianceRow(t *testing.T, row map[string]string) {
	t.Helper()

	if got := row["Firmware Baseline"]; got != "1902" {
		t.Errorf("Firmware Baseline: expected 1902, got %q", got)
	}

	bios := row["BIOS Profile"]
	for _, want := range []string{"ProcCStates=Disabled (Dell R640)", "SriovGlobalEnable=Enabled (HPE)"} {
		if !strings.Contains(bios, want) {
			t.Errorf("BIOS Profile: expected %q in %q", want, bios)
		}
	}
}

func TestRegistrationProfileGet_ShowsFirmwareBaselineAndBiosProfile(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/registration-profiles/7": testutils.RawHandler(http.StatusOK, registrationProfileWithComplianceItem),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	assertComplianceRow(t, renderCSV(t, func() {
		if err := RegistrationProfileGet(ctx, "7"); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	}))
}

func TestRegistrationProfileList_ShowsFirmwareBaselineAndBiosProfile(t *testing.T) {
	list := `{"data":[` + registrationProfileWithComplianceItem + `],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/registration-profiles": testutils.RawHandler(http.StatusOK, list),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	assertComplianceRow(t, renderCSV(t, func() {
		if err := RegistrationProfileList(ctx); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	}))
}

func TestRegistrationProfileGet_WithoutComplianceFields(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/registration-profiles/1": testutils.RawHandler(http.StatusOK, registrationProfileItem),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	row := renderCSV(t, func() {
		if err := RegistrationProfileGet(ctx, "1"); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	for _, column := range []string{"Firmware Baseline", "BIOS Profile"} {
		got, ok := row[column]
		if !ok {
			t.Errorf("missing column %q", column)
		} else if got != "" {
			t.Errorf("%s: expected empty for a profile without it, got %q", column, got)
		}
	}
}

func TestRegistrationProfileGet_JsonOutputUnchanged(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/registration-profiles/7": testutils.RawHandler(http.StatusOK, registrationProfileWithComplianceItem),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	out := testutils.CaptureStdout(t, func() {
		if err := RegistrationProfileGet(ctx, "7"); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	var decoded map[string]any
	if err := json.Unmarshal([]byte(out), &decoded); err != nil {
		t.Fatalf("json output is not valid: %v\n%s", err, out)
	}
	if _, leaked := decoded["BiosAttributes"]; leaked {
		t.Error("table-only BiosAttributes field leaked into json output")
	}

	settings, _ := decoded["settings"].(map[string]any)
	if settings["firmwareBaselineId"] != float64(1902) {
		t.Errorf("settings.firmwareBaselineId lost: %v", settings["firmwareBaselineId"])
	}
	if bios, _ := settings["biosProfile"].([]any); len(bios) != 2 {
		t.Errorf("settings.biosProfile should keep both entries, got %v", settings["biosProfile"])
	}
}

func TestRegistrationProfileSystemDefaults_ShowsFirmwareBaselineAndBiosProfile(t *testing.T) {
	const settingsJSON = `{
		"registerCredentials": "user",
		"firmwareBaselineId": 5,
		"biosProfile": [{"attributeName": "BootMode", "attributeValue": "Uefi", "serverVendor": "Dell"}]
	}`
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/servers/registration-profiles/system-defaults": testutils.RawHandler(http.StatusOK, settingsJSON),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)

	row := renderCSV(t, func() {
		if err := RegistrationProfileSystemDefaults(ctx); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	if row["Firmware Baseline"] != "5" {
		t.Errorf("Firmware Baseline: expected 5, got %q", row["Firmware Baseline"])
	}
	if row["BIOS Profile"] != "BootMode=Uefi (Dell)" {
		t.Errorf("BIOS Profile: got %q", row["BIOS Profile"])
	}
}
