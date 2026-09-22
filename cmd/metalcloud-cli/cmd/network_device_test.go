package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// networkDeviceItem satisfies sdk.NetworkDevice required fields.
var networkDeviceItem = map[string]interface{}{
	"id":                            "1",
	"revision":                      1.0,
	"status":                        "active",
	"siteId":                        1.0,
	"identifierString":              "sw-01",
	"description":                   "leaf switch",
	"chassisIdentifier":             "chassis-01",
	"country":                       "US",
	"city":                          "Chicago",
	"datacenterMeta":                "rack-a",
	"datacenterRoom":                "room-1",
	"datacenterRack":                "rack-01",
	"rackPositionUpperUnit":         1.0,
	"rackPositionLowerUnit":         1.0,
	"managementAddress":             "10.0.0.100",
	"managementAddressPrefixLength": 24.0,
	"managementAddressGateway":      "10.0.0.1",
	"managementPort":                22.0,
	// float32 fields (not bool) in the SDK model
	"syslogEnabled":                     0.0,
	"username":                          "admin",
	"managementMacAddress":              "00:11:22:33:44:55",
	"serialNumber":                      "SN-SW-001",
	"driver":                            "sonic_enterprise",
	"position":                          "leaf",
	"orderIndex":                        0.0,
	"tags":                              []interface{}{},
	"readyForInitialConfiguration":      0.0,
	"bootstrapReadinessCheckInProgress": 0.0,
	"subnetOobId":                       0.0,
	"subnetOobIndex":                    0.0,
	// bool fields
	"requiresOsInstall":                false,
	"bootstrapExpectedPartnerHostname": "",
	"loopbackAddressIpv6":              "",
	"asn":                              0.0,
	"vtepAddressIpv6":                  "",
	"mlagSystemMac":                    "",
	"mlagDomainId":                     0.0,
	"quarantineVlan":                   0.0,
	// map fields (not bool)
	"variablesMaterializedForOSAssets":      map[string]interface{}{},
	"secretsMaterializedForOSAssets":        map[string]interface{}{},
	"bootstrapReadinessCheckResult":         map[string]interface{}{},
	"isGateway":                             false,
	"vendorId":                              1.0,
	"applyIdentifierAsHostnameOnNextDeploy": false,
	"snmpServiceEnabled":                    false,
	"snmpMonitoringEnabled":                 false,
	"driftDetectionSyncStatus":              "synced",
	"backupEnabled":                         false,
	"driftDetectionEnabled":                 false,
	"portCount":                             32,
	"tagsMap":                               map[string]interface{}{},
}

func newNetworkDeviceTestServer() *httptest.Server {
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		// network-device list — FetchAllPages on GET /api/v2/network-devices
		mux.HandleFunc("/api/v2/network-devices", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(networkDeviceItem))
		})
		// network-device get — GetNetworkDevice GET /api/v2/network-devices/{id}
		mux.HandleFunc("/api/v2/network-devices/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(networkDeviceItem)
		})
	}))
}

// --- network-device list ---

func TestNetworkDeviceList_HappyPath(t *testing.T) {
	srv := newNetworkDeviceTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "network-device", "list")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "sw-01") {
		t.Fatalf("expected sw-01 in output, got: %s", out)
	}
}

func TestNetworkDeviceList_Alias(t *testing.T) {
	srv := newNetworkDeviceTestServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-device", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestNetworkDeviceList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "network-device", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- network-device get ---

func TestNetworkDeviceGet_HappyPath(t *testing.T) {
	srv := newNetworkDeviceTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "network-device", "get", "1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "sw-01") {
		t.Fatalf("expected sw-01 in output, got: %s", out)
	}
}

func TestNetworkDeviceGet_NoArgs(t *testing.T) {
	srv := newNetworkDeviceTestServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-device", "get"); err == nil {
		t.Fatal("expected error when no args given to network-device get")
	}
}

// --- network-device create ---

func TestNetworkDeviceCreate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/network-devices", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(networkDeviceItem)
		})
	}))
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "nd-create-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"driver":"sonic_enterprise","position":"leaf","username":"admin","managementPassword":"pass"}`)
	f.Close()

	if _, execErr := runCLI(t, srv, "network-device", "create", "--config-source", f.Name()); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

// --- network-device update ---

func TestNetworkDeviceUpdate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/network-devices/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(networkDeviceItem)
		})
	}))
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "nd-update-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"description":"updated"}`)
	f.Close()

	if _, execErr := runCLI(t, srv, "network-device", "update", "1", "--config-source", f.Name()); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

// --- network-device delete ---

func TestNetworkDeviceDelete(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/network-devices/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(networkDeviceItem)
		})
	}))
	defer srv.Close()

	if _, execErr := runCLI(t, srv, "network-device", "delete", "1"); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

func TestNetworkDeviceList_Formats(t *testing.T) {
	srv := newNetworkDeviceTestServer()
	defer srv.Close()
	for _, format := range []string{"json", "csv", "yaml", "text", "md"} {
		t.Run(format, func(t *testing.T) {
			out, err := runCLIFormat(t, srv, format, "network-device", "list")
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

// ---------------------------------------------------------------------------
// Network device sub-resource commands
// ---------------------------------------------------------------------------

// ndCmdBreakoutItem satisfies sdk.NetworkEquipmentBreakout.
var ndCmdBreakoutItem = map[string]interface{}{
	"id": 5, "portName": "Ethernet0", "serviceStatus": "active", "pendingDelete": false,
	"revision":       2,
	"breakoutGroups": []interface{}{map[string]interface{}{"numberOfInterfaces": 4, "speed": "25G"}},
	"config": map[string]interface{}{
		"revision":       9,
		"breakoutGroups": []interface{}{map[string]interface{}{"numberOfInterfaces": 2, "speed": "50G"}},
	},
}

// ndCmdVendorItem satisfies sdk.NetworkDeviceVendors.
var ndCmdVendorItem = map[string]interface{}{
	"id": 3, "revision": 4, "kind": "sonic_enterprise",
	"oidGroups":             []interface{}{},
	"healthCheckRules":      map[string]interface{}{"rules": []interface{}{}},
	"optionalFilesToBackup": map[string]interface{}{"files": map[string]interface{}{"default": []interface{}{"/etc/config"}}},
}

// ndCmdDriftItem satisfies sdk.NetworkDeviceDriftHistory.
var ndCmdDriftItem = map[string]interface{}{
	"id": 9, "networkDeviceId": 1, "snapshotId": "abc123",
	"configurationDrift": "- a\n+ b", "createdTimestamp": "2024-01-01T00:00:00Z",
}

// newNetworkDeviceSubResourceServer serves the sub-resource endpoints used by
// the command-level tests below.
func newNetworkDeviceSubResourceServer() *httptest.Server {
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/network-devices/1", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, networkDeviceItem)
		})
		mux.HandleFunc("/api/v2/network-devices/1/breakouts", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, []interface{}{ndCmdBreakoutItem})
		})
		mux.HandleFunc("/api/v2/network-devices/1/breakouts/5", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, ndCmdBreakoutItem)
		})
		mux.HandleFunc("/api/v2/network-devices/1/drift-history", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(ndCmdDriftItem))
		})
		mux.HandleFunc("/api/v2/network-devices/1/secrets", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{"names": []interface{}{"enable_password"}})
		})
		mux.HandleFunc("/api/v2/network-devices/1/actions/snmp-monitoring-subscribe", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
		mux.HandleFunc("/api/v2/network-devices/statistics", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{"deviceCount": 2, "portCount": 64})
		})
		mux.HandleFunc("/api/v2/network-devices/vendors", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, paginatedList(ndCmdVendorItem))
		})
		mux.HandleFunc("/api/v2/network-devices/drivers/capabilities", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, []interface{}{
				map[string]interface{}{"driver": "sonic_enterprise", "activationActions": []interface{}{"start-registration"}},
			})
		})
	}))
}

func TestNetworkDeviceSubResourceReadCommands(t *testing.T) {
	srv := newNetworkDeviceSubResourceServer()
	defer srv.Close()

	cases := []struct {
		name     string
		args     []string
		contains string
	}{
		{"breakout list", []string{"network-device", "breakout", "list", "1"}, "Ethernet0"},
		{"breakout get", []string{"network-device", "breakout", "get", "1", "5"}, "Ethernet0"},
		{"drift list", []string{"network-device", "drift", "list", "1"}, "abc123"},
		{"secret list", []string{"network-device", "secret", "list", "1"}, "enable_password"},
		{"statistics", []string{"network-device", "statistics"}, "64"},
		{"vendor list", []string{"network-device", "vendor", "list"}, "sonic_enterprise"},
		{"driver capabilities", []string{"network-device", "driver-capabilities"}, "start-registration"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := runCLI(t, srv, tc.args...)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if !strings.Contains(out, tc.contains) {
				t.Errorf("expected %q in output, got: %s", tc.contains, out)
			}
		})
	}
}

func TestNetworkDeviceSnmpMonitoringEnableCmd(t *testing.T) {
	srv := newNetworkDeviceSubResourceServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "network-device", "snmp-monitoring", "enable", "1"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestNetworkDeviceSubResourceArgValidation(t *testing.T) {
	srv := newNetworkDeviceSubResourceServer()
	defer srv.Close()

	cases := [][]string{
		{"network-device", "breakout", "get", "1"},
		{"network-device", "port", "ip", "list", "1", "42"},
		{"network-device", "drift", "acknowledge", "1"},
		{"network-device", "secret", "set", "1"},
		{"network-device", "replace", "1"},
	}

	for _, args := range cases {
		t.Run(strings.Join(args[1:], " "), func(t *testing.T) {
			if _, err := runCLI(t, srv, args...); err == nil {
				t.Error("expected an argument validation error, got nil")
			}
		})
	}
}

func TestNetworkDeviceSubResourceRequiredFlags(t *testing.T) {
	srv := newNetworkDeviceSubResourceServer()
	defer srv.Close()

	cases := [][]string{
		{"network-device", "breakout", "create", "1"},
		{"network-device", "port", "add", "1"},
		{"network-device", "re-provision", "1"},
		{"network-device", "run-extension", "1"},
		{"network-device", "vendor", "update", "3"},
		{"network-device", "set-health-monitoring-filter"},
	}

	for _, args := range cases {
		t.Run(strings.Join(args[1:], " "), func(t *testing.T) {
			if _, err := runCLI(t, srv, args...); err == nil {
				t.Error("expected a required flag error, got nil")
			}
		})
	}
}

func TestNetworkDeviceConfigExampleCmds(t *testing.T) {
	srv := newNetworkDeviceSubResourceServer()
	defer srv.Close()

	cases := [][]string{
		{"network-device", "port", "config-example"},
		{"network-device", "port", "ip", "config-example"},
		{"network-device", "breakout", "config-example"},
		{"network-device", "breakout", "update-config-example"},
		{"network-device", "vendor", "config-example"},
		{"network-device", "return-to-planned-config-example"},
		{"network-device", "run-extension-config-example"},
	}

	for _, args := range cases {
		t.Run(strings.Join(args[1:], " "), func(t *testing.T) {
			out, err := runCLI(t, srv, args...)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if out == "" {
				t.Error("expected example output, got empty string")
			}
		})
	}
}
