package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// endpointItem satisfies sdk.Endpoint required fields.
var endpointItem = map[string]interface{}{
	"id":               "1",
	"revision":         "1",
	"siteId":           1,
	"name":             "ep-01",
	"label":            "ep-01",
	"createdTimestamp": "2024-01-01T00:00:00Z",
	"updatedTimestamp": "2024-01-01T00:00:00Z",
}

func newEndpointTestServer() *httptest.Server {
	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		// endpoint list — FetchAllPages / GetEndpoints
		mux.HandleFunc("/api/v2/endpoints", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(paginatedList(endpointItem))
		})
		// endpoint get — GetEndpointById
		mux.HandleFunc("/api/v2/endpoints/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(endpointItem)
		})
	}))
}

// --- endpoint list ---

func TestEndpointList_HappyPath(t *testing.T) {
	srv := newEndpointTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "endpoint", "list")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "ep-01") {
		t.Fatalf("expected ep-01 in output, got: %s", out)
	}
}

func TestEndpointList_Alias(t *testing.T) {
	srv := newEndpointTestServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "endpoint", "ls"); err != nil {
		t.Fatalf("alias ls: expected no error, got: %v", err)
	}
}

func TestEndpointList_NoEndpoint(t *testing.T) {
	if _, err := runCLI(t, nil, "endpoint", "list"); err == nil {
		t.Fatal("expected error when endpoint is empty")
	}
}

// --- endpoint get ---

func TestEndpointGet_HappyPath(t *testing.T) {
	srv := newEndpointTestServer()
	defer srv.Close()

	out, err := runCLI(t, srv, "endpoint", "get", "1")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "ep-01") {
		t.Fatalf("expected ep-01 in output, got: %s", out)
	}
}

func TestEndpointGet_NoArgs(t *testing.T) {
	srv := newEndpointTestServer()
	defer srv.Close()

	if _, err := runCLI(t, srv, "endpoint", "get"); err == nil {
		t.Fatal("expected error when no args given to endpoint get")
	}
}

// --- endpoint create ---

func TestEndpointCreate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/endpoints", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(endpointItem)
		})
	}))
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "ep-create-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"siteId":1,"name":"ep-new","label":"ep-new"}`)
	f.Close()

	if _, execErr := runCLI(t, srv, "endpoint", "create", "--config-source", f.Name()); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

// --- endpoint update ---

func TestEndpointUpdate(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/endpoints/1", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(endpointItem)
		})
	}))
	defer srv.Close()

	f, err := os.CreateTemp(t.TempDir(), "ep-update-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	_, _ = f.WriteString(`{"name":"ep-updated"}`)
	f.Close()

	if _, execErr := runCLI(t, srv, "endpoint", "update", "1", "--config-source", f.Name()); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

// --- endpoint delete ---

func TestEndpointDelete(t *testing.T) {
	srv := httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/endpoints/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(endpointItem)
		})
	}))
	defer srv.Close()

	if _, execErr := runCLI(t, srv, "endpoint", "delete", "1"); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}
}

// --- endpoint interface get/add/update/remove and device interfaces ---

// endpointInterfaceItem satisfies sdk.EndpointInterface required fields.
var endpointInterfaceItem = map[string]interface{}{
	"id": 3, "revision": "2", "networkDeviceId": 10,
	"networkDeviceInterfaceId": 20, "networkDeviceInterfaceName": "swp1",
	"macAddress":       "AA:BB:CC:DD:EE:FF",
	"createdTimestamp": "2024-01-01T00:00:00Z", "updatedTimestamp": "2024-01-01T00:00:00Z",
}

// newEndpointInterfaceServer records the If-Match header and body of the last
// mutating request on the interface sub-resource.
func newEndpointInterfaceServer(lastIfMatch *string, lastBody *string) *httptest.Server {
	record := func(r *http.Request) {
		*lastIfMatch = r.Header.Get("If-Match")
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		raw, _ := json.Marshal(body)
		*lastBody = string(raw)
	}

	return httptest.NewServer(newMux(allPerms, func(mux *http.ServeMux) {
		mux.HandleFunc("/api/v2/endpoints/1/interfaces", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				record(r)
				jsonResponse(w, http.StatusCreated, endpointInterfaceItem)
				return
			}
			jsonResponse(w, http.StatusOK, paginatedList(endpointInterfaceItem))
		})
		mux.HandleFunc("/api/v2/endpoints/1/interfaces/3", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodDelete:
				*lastIfMatch = r.Header.Get("If-Match")
				w.WriteHeader(http.StatusNoContent)
			case http.MethodPatch:
				record(r)
				jsonResponse(w, http.StatusOK, endpointInterfaceItem)
			default:
				jsonResponse(w, http.StatusOK, endpointInterfaceItem)
			}
		})
		mux.HandleFunc("/api/v2/network-devices/10", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, networkDeviceFixtureJSON())
		})
		mux.HandleFunc("/api/v2/endpoints/network-devices/10/interfaces", func(w http.ResponseWriter, r *http.Request) {
			jsonResponse(w, http.StatusOK, map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"networkDeviceId": 10, "networkDeviceInterfaceId": 20,
						"networkDeviceInterfaceName": "swp1", "endpointId": 1, "endpointInterfaceId": 3,
					},
				},
			})
		})
	}))
}

// networkDeviceFixtureJSON returns the minimal network device satisfying the
// strict SDK model, decoded from the shared testutils fixture shape.
func networkDeviceFixtureJSON() map[string]interface{} {
	return map[string]interface{}{
		"id": "10", "revision": 1, "status": "active", "vendorId": 1, "siteId": 1,
		"identifierString": "leaf-01", "applyIdentifierAsHostnameOnNextDeploy": false,
		"description": "", "chassisIdentifier": "", "country": "", "city": "",
		"datacenterMeta": "", "datacenterRoom": "", "datacenterRack": "",
		"rackPositionUpperUnit": 1, "rackPositionLowerUnit": 1,
		"managementAddress": "10.0.0.100", "managementAddressPrefixLength": 24,
		"managementAddressGateway": "10.0.0.1", "managementPort": 22,
		"syslogEnabled": 0, "snmpServiceEnabled": false, "snmpMonitoringEnabled": false,
		"username": "admin", "managementMacAddress": "00:11:22:33:44:55", "serialNumber": "SN-1",
		"driver": "sonic_enterprise", "position": "leaf", "backupEnabled": false,
		"driftDetectionEnabled": false, "driftDetectionSyncStatus": "not_supported",
		"orderIndex": 0, "tags": []interface{}{}, "tagsMap": map[string]interface{}{},
		"readyForInitialConfiguration": 0, "bootstrapReadinessCheckInProgress": 0,
		"subnetOobId": 0, "subnetOobIndex": 0, "requiresOsInstall": false,
		"bootstrapExpectedPartnerHostname": "", "loopbackAddressIpv6": "", "asn": 0,
		"vtepAddressIpv6": "", "mlagSystemMac": "", "mlagDomainId": 0, "quarantineVlan": 0,
		"variablesMaterializedForOSAssets": map[string]interface{}{},
		"secretsMaterializedForOSAssets":   map[string]interface{}{},
		"bootstrapReadinessCheckResult":    map[string]interface{}{},
		"isGateway":                        false, "portCount": 0,
	}
}

func TestEndpointInterfaceGet(t *testing.T) {
	var ifMatch, body string
	srv := newEndpointInterfaceServer(&ifMatch, &body)
	defer srv.Close()

	out, err := runCLI(t, srv, "endpoint", "interface", "get", "1", "3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "swp1") {
		t.Errorf("expected the switch port in output, got: %s", out)
	}

	if _, err := runCLI(t, srv, "endpoint", "interface", "get", "1"); err == nil {
		t.Error("expected an error when the interface id is missing")
	}
}

func TestEndpointInterfaceAdd_Flags(t *testing.T) {
	var ifMatch, body string
	srv := newEndpointInterfaceServer(&ifMatch, &body)
	defer srv.Close()

	_, err := runCLI(t, srv, "endpoint", "interface", "add", "1",
		"--network-device-interface-id", "20", "--mac-address", "AA:BB:CC:DD:EE:FF")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{`"networkDeviceInterfaceId":20`, `"macAddress":"AA:BB:CC:DD:EE:FF"`} {
		if !strings.Contains(body, want) {
			t.Errorf("add body missing %s: %s", want, body)
		}
	}
}

func TestEndpointInterfaceAdd_RequiresPortOrConfigSource(t *testing.T) {
	var ifMatch, body string
	srv := newEndpointInterfaceServer(&ifMatch, &body)
	defer srv.Close()

	if _, err := runCLI(t, srv, "endpoint", "interface", "add", "1"); err == nil {
		t.Fatal("expected an error when neither --network-device-interface-id nor --config-source is given")
	}
}

func TestEndpointInterfaceUpdate_SendsIfMatch(t *testing.T) {
	var ifMatch, body string
	srv := newEndpointInterfaceServer(&ifMatch, &body)
	defer srv.Close()

	if _, err := runCLI(t, srv, "endpoint", "interface", "update", "1", "3", "--mac-address", "11:22:33:44:55:66"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ifMatch != "2" {
		t.Errorf("expected If-Match 2 (the interface revision), got %q", ifMatch)
	}
	if !strings.Contains(body, `"macAddress":"11:22:33:44:55:66"`) {
		t.Errorf("update body missing the MAC address: %s", body)
	}
}

func TestEndpointInterfaceRemove_SendsIfMatch(t *testing.T) {
	var ifMatch, body string
	srv := newEndpointInterfaceServer(&ifMatch, &body)
	defer srv.Close()

	if _, err := runCLI(t, srv, "endpoint", "interface", "remove", "1", "3"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ifMatch != "2" {
		t.Errorf("expected If-Match 2 (the interface revision), got %q", ifMatch)
	}
}

func TestEndpointGetNetworkDeviceInterfaces(t *testing.T) {
	var ifMatch, body string
	srv := newEndpointInterfaceServer(&ifMatch, &body)
	defer srv.Close()

	out, err := runCLI(t, srv, "endpoint", "get-network-device-interfaces", "10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "swp1") {
		t.Errorf("expected the switch port in output, got: %s", out)
	}
}
