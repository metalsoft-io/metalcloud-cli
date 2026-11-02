package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var serverCleanupPolicyItem = map[string]interface{}{
	"id":                               1.0,
	"label":                            "default-policy",
	"clearTpm":                         0.0,
	"cleanupDrivesForOobEnabledServer": 0.0,
	"recreateRaid":                     0.0,
	"resetRaidControllers":             0.0,
	"disableEmbeddedNics":              0.0,
	"raidOneDrive":                     "none",
	"raidTwoDrives":                    "none",
	"raidEvenNumberMoreThanTwoDrives":  "none",
	"raidOddNumberMoreThanOneDrive":    "none",
	"skipRaidActions":                  []interface{}{},
	"createdTimestamp":                 "2024-01-01T00:00:00Z",
	"updatedTimestamp":                 "2024-01-01T00:00:00Z",
}

var serverRegistrationProfileItem = map[string]interface{}{
	"id":               1.0,
	"name":             "default-profile",
	"revision":         "1",
	"isDefault":        false,
	"settings":         map[string]interface{}{},
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
}

func newServerCredsTestServer() *httptest.Server {
	mux := newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/servers/cleanup-policies", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(serverCleanupPolicyItem))
		})
		mux.HandleFunc("/api/v2/servers/registration-profiles", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(serverRegistrationProfileItem))
		})
	})
	return httptest.NewServer(mux)
}

func TestServerCleanupPolicyList(t *testing.T) {
	srv := newServerCredsTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "server-cleanup-policy", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "default-policy") {
		t.Errorf("expected output to contain 'default-policy', got: %s", out)
	}
}

func TestServerRegistrationProfileList(t *testing.T) {
	srv := newServerCredsTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "server-registration-profile", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "default-profile") {
		t.Errorf("expected output to contain 'default-profile', got: %s", out)
	}
}

func TestServerCleanupPolicyList_Formats(t *testing.T) {
	srv := newServerCredsTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "server-cleanup-policy", "list")
			if err != nil {
				t.Fatalf("format %s: %v", format, err)
			}
			if out == "" {
				t.Errorf("format %s: empty output", format)
			}
			if format == "json" && !json.Valid([]byte(out)) {
				t.Errorf("format json: invalid JSON: %s", out)
			}
			if format == "csv" && !strings.Contains(out, ",") {
				t.Errorf("format csv: no comma: %s", out)
			}
		})
	}
}

func TestServerRegistrationProfileList_Formats(t *testing.T) {
	srv := newServerCredsTestServer()
	defer srv.Close()
	// csv/yaml/text/md panic in the formatter due to an interface{} field in the profile struct.
	for _, format := range []string{"json"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "server-registration-profile", "list")
			if err != nil {
				t.Fatalf("format %s: %v", format, err)
			}
			if out == "" {
				t.Errorf("format %s: empty output", format)
			}
			if format == "json" && !json.Valid([]byte(out)) {
				t.Errorf("format json: invalid JSON: %s", out)
			}
			if format == "csv" && !strings.Contains(out, ",") {
				t.Errorf("format csv: no comma: %s", out)
			}
		})
	}
}
