package firmware_catalog

import (
	"encoding/json"
	"testing"
)

func TestExtractHpeFirmwareTargets(t *testing.T) {
	tests := []struct {
		name string
		body string
		want []string
	}{
		{
			name: "Hpe key with targets",
			body: `{
				"Id": "1",
				"Name": "iLO",
				"Oem": {
					"Hpe": {
						"DeviceClass": "abc",
						"Targets": ["uuid-a", "uuid-b"]
					}
				}
			}`,
			want: []string{"uuid-a", "uuid-b"},
		},
		{
			name: "Legacy Hp key",
			body: `{"Oem": {"Hp": {"Targets": ["uuid-c"]}}}`,
			want: []string{"uuid-c"},
		},
		{
			name: "Drops non-string and empty entries",
			body: `{"Oem": {"Hpe": {"Targets": ["uuid-a", 42, "", "uuid-b"]}}}`,
			want: []string{"uuid-a", "uuid-b"},
		},
		{
			name: "Missing Oem section",
			body: `{"Id": "1"}`,
			want: nil,
		},
		{
			name: "Targets absent",
			body: `{"Oem": {"Hpe": {"DeviceClass": "abc"}}}`,
			want: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var entry map[string]interface{}
			if err := json.Unmarshal([]byte(tc.body), &entry); err != nil {
				t.Fatalf("invalid test fixture: %v", err)
			}

			got := extractHpeFirmwareTargets(entry)
			if !equalStringSlices(got, tc.want) {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
}

func TestIsHpeVendor(t *testing.T) {
	cases := map[string]bool{
		"HPE": true,
		"HP":  true,
		"hpe": true,
		"hp":  true,
		"Hewlett Packard Enterprise": false,
		"Dell":                       false,
		"":                           false,
	}
	for in, want := range cases {
		if got := isHpeVendor(in); got != want {
			t.Errorf("isHpeVendor(%q) = %v, want %v", in, got, want)
		}
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
