package network_device

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// decodeItems converts a slice of JSON strings into a slice of decoded objects
// suitable for use with testutils.MultiPageServer.
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

// minimal NetworkDevice JSON — only required fields
const ndItem = `{
	"id":"1","revision":1,"status":"active","vendorId":1,"siteId":1,
	"identifierString":"sw-01","applyIdentifierAsHostnameOnNextDeploy":false,
	"description":"","chassisIdentifier":"",
	"country":"","city":"","datacenterMeta":"","datacenterRoom":"","datacenterRack":"",
	"rackPositionUpperUnit":0,"rackPositionLowerUnit":0,
	"managementAddress":"10.0.0.1","managementAddressPrefixLength":24,
	"managementAddressGateway":"10.0.0.254","managementPort":22,
	"syslogEnabled":0,"snmpServiceEnabled":false,"snmpMonitoringEnabled":false,
	"username":"admin","managementMacAddress":"AA:BB:CC:DD:EE:01",
	"serialNumber":"SN001","driver":"sonic_enterprise","position":"leaf",
	"driftDetectionSyncStatus":"",
	"orderIndex":1,"tags":[],"tagsMap":{},"readyForInitialConfiguration":0,
	"bootstrapReadinessCheckInProgress":0,"subnetOobId":0,"subnetOobIndex":0,
	"requiresOsInstall":false,"bootstrapExpectedPartnerHostname":"",
	"loopbackAddressIpv6":"","asn":65000,"vtepAddressIpv6":"",
	"mlagSystemMac":"","mlagDomainId":0,"quarantineVlan":0,
	"variablesMaterializedForOSAssets":{},"secretsMaterializedForOSAssets":{},
	"bootstrapReadinessCheckResult":{},"isGateway":false
}`

func ndListHandler(statusCode int, items []string, currentPage, totalPages int) http.HandlerFunc {
	data := "[" + strings.Join(items, ",") + "]"
	body := fmt.Sprintf(`{"data":%s,"meta":{"currentPage":%d,"totalPages":%d,"itemsPerPage":100}}`,
		data, currentPage, totalPages)
	return testutils.RawHandler(statusCode, body)
}

// --- NetworkDeviceList ---

func TestNetworkDeviceList_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices": ndListHandler(http.StatusOK, []string{ndItem, ndItem}, 1, 1),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceList(ctx, nil); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceList(ctx, nil); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

func TestNetworkDeviceList_Empty(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices": ndListHandler(http.StatusOK, []string{}, 1, 1),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceList(ctx, nil); err != nil {
		t.Errorf("expected nil error for empty list, got: %v", err)
	}
}

func TestNetworkDeviceList_Pagination(t *testing.T) {
	makeItems := func(n int) []string {
		items := make([]string, n)
		for i := range items {
			items[i] = fmt.Sprintf(`{
				"id":"%d","revision":1,"status":"active","vendorId":1,"siteId":1,
				"identifierString":"sw-%d","applyIdentifierAsHostnameOnNextDeploy":false,
				"description":"","chassisIdentifier":"",
				"country":"","city":"","datacenterMeta":"","datacenterRoom":"","datacenterRack":"",
				"rackPositionUpperUnit":0,"rackPositionLowerUnit":0,
				"managementAddress":"10.0.0.1","managementAddressPrefixLength":24,
				"managementAddressGateway":"10.0.0.254","managementPort":22,
				"syslogEnabled":0,"snmpServiceEnabled":false,"snmpMonitoringEnabled":false,
				"username":"admin","managementMacAddress":"AA:BB:CC:DD:EE:01",
				"serialNumber":"SN%d","driver":"sonic_enterprise","position":"leaf",
				"driftDetectionSyncStatus":"",
				"orderIndex":1,"tags":[],"tagsMap":{},"readyForInitialConfiguration":0,
				"bootstrapReadinessCheckInProgress":0,"subnetOobId":0,"subnetOobIndex":0,
				"requiresOsInstall":false,"bootstrapExpectedPartnerHostname":"",
				"loopbackAddressIpv6":"","asn":65000,"vtepAddressIpv6":"",
				"mlagSystemMac":"","mlagDomainId":0,"quarantineVlan":0,
				"variablesMaterializedForOSAssets":{},"secretsMaterializedForOSAssets":{},
				"bootstrapReadinessCheckResult":{},"isGateway":false
			}`, i+1, i+1, i+1)
		}
		return items
	}

	ts := testutils.MultiPageServer("/api/v2/network-devices", []any{
		decodeItems(makeItems(100)),
		decodeItems(makeItems(100)),
		decodeItems(makeItems(5)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceList(ctx, nil); err != nil {
		t.Errorf("expected nil error across 3 pages, got: %v", err)
	}
}

// --- NetworkDeviceGet ---

func TestNetworkDeviceGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/1": testutils.RawHandler(http.StatusOK, ndItem),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceGet(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceGet(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestNetworkDeviceGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceGet(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

// --- NetworkDeviceCreate ---

func TestNetworkDeviceCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, ndItem)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{
		"siteId":1,"driver":"sonic_enterprise","identifierString":"sw-01",
		"serialNumber":"SN001","position":"leaf","isGateway":false,
		"isStorageSwitch":false,"isBorderDevice":false,
		"managementPassword":"password","username":"admin"
	}`)
	if err := NetworkDeviceCreate(ctx, config); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceCreate_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices": testutils.ErrorHandler(http.StatusBadRequest, "validation error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceCreate(ctx, []byte(`{}`)); err == nil {
		t.Error("expected error for 400, got nil")
	}
}

// --- NetworkDeviceUpdate ---

func TestNetworkDeviceUpdate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/1": func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, ndItem)
			case http.MethodPatch:
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, ndItem)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceUpdate(ctx, "1", []byte(`{"identifierString":"sw-updated"}`)); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceUpdate_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceUpdate(ctx, "99", []byte(`{}`)); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

// --- NetworkDeviceDelete ---

func TestNetworkDeviceDelete_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/1": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDelete(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceDelete_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDelete(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestNetworkDeviceDelete_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDelete(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

// --- NetworkDeviceDefaultSecretsList (in network_device package) ---

const ndSecretsItem = `{
	"id":1,"siteId":1,"macAddressOrSerialNumber":"AA:BB:CC:DD:EE:01",
	"secretName":"admin-password","createdTimestamp":"2024-01-01T00:00:00Z",
	"updatedTimestamp":"2024-01-01T00:00:00Z"
}`

func ndSecretsListHandler(statusCode int, items []string, currentPage, totalPages int) http.HandlerFunc {
	data := "[" + strings.Join(items, ",") + "]"
	body := fmt.Sprintf(`{"data":%s,"meta":{"currentPage":%d,"totalPages":%d,"itemsPerPage":100}}`,
		data, currentPage, totalPages)
	return testutils.RawHandler(statusCode, body)
}

func TestNetworkDeviceDefaultSecretsList_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets": ndSecretsListHandler(http.StatusOK, []string{ndSecretsItem, ndSecretsItem}, 1, 1),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsList(ctx, 0, 0); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceDefaultSecretsList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsList(ctx, 0, 0); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

func TestNetworkDeviceDefaultSecretsList_Empty(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets": ndSecretsListHandler(http.StatusOK, []string{}, 1, 1),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsList(ctx, 0, 0); err != nil {
		t.Errorf("expected nil error for empty list, got: %v", err)
	}
}

func TestNetworkDeviceDefaultSecretsList_Pagination(t *testing.T) {
	makeItems := func(n int) []string {
		items := make([]string, n)
		for i := range items {
			items[i] = fmt.Sprintf(`{
				"id":%d,"siteId":1,"macAddressOrSerialNumber":"AA:BB:CC:DD:EE:%02d",
				"secretName":"secret-%d","createdTimestamp":"2024-01-01T00:00:00Z",
				"updatedTimestamp":"2024-01-01T00:00:00Z"
			}`, i+1, i+1, i+1)
		}
		return items
	}

	ts := testutils.MultiPageServer("/api/v2/network-devices/default-secrets", []any{
		decodeItems(makeItems(100)),
		decodeItems(makeItems(100)),
		decodeItems(makeItems(5)),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsList(ctx, 0, 0); err != nil {
		t.Errorf("expected nil error across 3 pages, got: %v", err)
	}
}

func TestNetworkDeviceDefaultSecretsList_SinglePage(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets": ndSecretsListHandler(http.StatusOK, []string{ndSecretsItem}, 1, 2),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	// page=1, limit=5: uses single-page path (no FetchAllPages)
	if err := NetworkDeviceDefaultSecretsList(ctx, 1, 5); err != nil {
		t.Errorf("expected nil error for page/limit path, got: %v", err)
	}
}

// --- NetworkDeviceArchive ---

func TestNetworkDeviceArchive_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/1": func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, ndItem)
			case http.MethodPost:
				w.WriteHeader(http.StatusNoContent)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		},
		"/api/v2/network-devices/1/actions/archive": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceArchive(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceArchive_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceArchive(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestNetworkDeviceArchive_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceArchive(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

// --- NetworkDeviceDiscover ---

func TestNetworkDeviceDiscover_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDiscover(ctx, "not-a-number", nil); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

func TestNetworkDeviceDiscover_InvalidTarget(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDiscover(ctx, "1", []string{"bogus"}); err == nil {
		t.Error("expected error for invalid discovery target, got nil")
	}
}

// TestNetworkDeviceDiscover_DefaultAllTargets verifies that with no targets the
// full discovery payload (all three types, persistData=true) is sent.
func TestNetworkDeviceDiscover_DefaultAllTargets(t *testing.T) {
	var body sdk.DiscoveryQuery
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/1/actions/discover": func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewDecoder(r.Body).Decode(&body)
			w.WriteHeader(http.StatusOK)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDiscover(ctx, "1", nil); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if !slices.Equal(body.Discover, DiscoveryTargets) {
		t.Errorf("expected discover=%v, got %v", DiscoveryTargets, body.Discover)
	}
	if !body.PersistData {
		t.Error("expected persistData=true")
	}
}

// TestNetworkDeviceDiscover_SpecificTarget verifies a single requested target is
// forwarded verbatim.
func TestNetworkDeviceDiscover_SpecificTarget(t *testing.T) {
	var body sdk.DiscoveryQuery
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/1/actions/discover": func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewDecoder(r.Body).Decode(&body)
			w.WriteHeader(http.StatusOK)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDiscover(ctx, "1", []string{"ports"}); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if !slices.Equal(body.Discover, []string{"ports"}) {
		t.Errorf("expected discover=[ports], got %v", body.Discover)
	}
}

// --- NetworkDeviceGetCredentials ---

func TestNetworkDeviceGetCredentials_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/1/credentials": testutils.RawHandler(http.StatusOK, `{"username":"admin","password":"secret","host":"10.0.0.1","port":22,"datacenter":"site1","driver":"sonic_enterprise","hostname":"sw-01"}`),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceGetCredentials(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceGetCredentials_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/99/credentials": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceGetCredentials(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestNetworkDeviceGetCredentials_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceGetCredentials(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

// --- NetworkDeviceGetPorts ---

const ndPortItem = `{
	"interfaceId":1,"networkDeviceId":1,"interfaceName":"Ethernet1","kind":"port",
	"dirtyBit":0,"cachedUpdatedTimestamp":"2024-01-01T00:00:00Z",
	"tags":{},"config":{"revision":1},"ipv4":{"addresses":[]},"ipv6":{"addresses":[]},
	"portName":"Ethernet1","enabled":true,"active":true,
	"linkSpeed":10000,"linkDuplex":"full","utilizationIn":0,"utilizationOut":0
}`

func ndPortsListHandler(statusCode int, items []string) http.HandlerFunc {
	data := "[" + strings.Join(items, ",") + "]"
	body := fmt.Sprintf(`{"data":%s,"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`, data)
	return testutils.RawHandler(statusCode, body)
}

func TestNetworkDeviceGetPorts_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/1/ports": ndPortsListHandler(http.StatusOK, []string{ndPortItem}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceGetPorts(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceGetPorts_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/1/ports": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceGetPorts(ctx, "1"); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

func TestNetworkDeviceGetPorts_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceGetPorts(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

// --- NetworkDeviceReset ---

func TestNetworkDeviceReset_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/1/actions/reset": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceReset(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceReset_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/1/actions/reset": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceReset(ctx, "1"); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

func TestNetworkDeviceReset_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceReset(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

// --- NetworkDeviceSetFailed ---

func TestNetworkDeviceSetFailed_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/1": testutils.RawHandler(http.StatusOK, ndItem),
		"/api/v2/network-devices/1/actions/set-as-failed": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, ndItem)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceSetFailed(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceSetFailed_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceSetFailed(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestNetworkDeviceSetFailed_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceSetFailed(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

// --- NetworkDeviceEnableSyslog ---

func TestNetworkDeviceEnableSyslog_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/1/actions/syslog-subscribe": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceEnableSyslog(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceEnableSyslog_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/1/actions/syslog-subscribe": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceEnableSyslog(ctx, "1"); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

func TestNetworkDeviceEnableSyslog_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceEnableSyslog(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

// --- NetworkDeviceSetPortStatus ---

func TestNetworkDeviceSetPortStatus_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/1/actions/set-port-status": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceSetPortStatus(ctx, "1", "Ethernet1", "up"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceSetPortStatus_Down(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/1/actions/set-port-status": func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceSetPortStatus(ctx, "1", "Ethernet1", "down"); err != nil {
		t.Errorf("expected nil error for down action, got: %v", err)
	}
}

func TestNetworkDeviceSetPortStatus_InvalidAction(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceSetPortStatus(ctx, "1", "Ethernet1", "invalid"); err == nil {
		t.Error("expected error for invalid action, got nil")
	}
}

func TestNetworkDeviceSetPortStatus_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceSetPortStatus(ctx, "not-a-number", "Ethernet1", "up"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

func TestNetworkDeviceSetPortStatus_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/1/actions/set-port-status": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceSetPortStatus(ctx, "1", "Ethernet1", "up"); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

// --- NetworkDeviceGetDefaults ---

const ndDefaultsItem = `{
	"id":1,"datacenterName":"site1","serialNumber":"SN001",
	"managementMacAddress":"AA:BB:CC:DD:EE:01","position":"leaf",
	"identifierString":"sw-01","asn":65000
}`

func ndDefaultsListHandler(statusCode int, items []string) http.HandlerFunc {
	data := "[" + strings.Join(items, ",") + "]"
	body := fmt.Sprintf(`{"data":%s,"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`, data)
	return testutils.RawHandler(statusCode, body)
}

func TestNetworkDeviceGetDefaults_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/defaults/1": ndDefaultsListHandler(http.StatusOK, []string{ndDefaultsItem}),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceGetDefaults(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceGetDefaults_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/defaults/1": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceGetDefaults(ctx, "1"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestNetworkDeviceGetDefaults_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceGetDefaults(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid site ID, got nil")
	}
}

// --- NetworkDeviceAddDefaults ---

func TestNetworkDeviceAddDefaults_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/defaults": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"datacenterName":"site1","managementMacAddress":"AA:BB:CC:DD:EE:01"}`)
	if err := NetworkDeviceAddDefaults(ctx, config); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceAddDefaults_MissingMac(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"datacenterName":"site1"}`)
	if err := NetworkDeviceAddDefaults(ctx, config); err == nil {
		t.Error("expected error for missing MAC address, got nil")
	}
}

func TestNetworkDeviceAddDefaults_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/defaults": testutils.ErrorHandler(http.StatusBadRequest, "validation error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	config := []byte(`{"datacenterName":"site1","managementMacAddress":"AA:BB:CC:DD:EE:01"}`)
	if err := NetworkDeviceAddDefaults(ctx, config); err == nil {
		t.Error("expected error for 400, got nil")
	}
}

// --- NetworkDeviceDeleteDefaults ---

func TestNetworkDeviceDeleteDefaults_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/defaults/1/1": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDeleteDefaults(ctx, "1", "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceDeleteDefaults_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/defaults/1/1": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDeleteDefaults(ctx, "1", "1"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestNetworkDeviceDeleteDefaults_InvalidSiteId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDeleteDefaults(ctx, "not-a-number", "1"); err == nil {
		t.Error("expected error for invalid site ID, got nil")
	}
}

func TestNetworkDeviceDeleteDefaults_InvalidDefaultsId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDeleteDefaults(ctx, "1", "not-a-number"); err == nil {
		t.Error("expected error for invalid defaults ID, got nil")
	}
}

// --- NetworkDeviceDefaultSecretsGet ---

func TestNetworkDeviceDefaultSecretsGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets/1": testutils.RawHandler(http.StatusOK, ndSecretsItem),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsGet(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceDefaultSecretsGet_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsGet(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestNetworkDeviceDefaultSecretsGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsGet(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

// --- NetworkDeviceDefaultSecretsGetCredentials ---

func TestNetworkDeviceDefaultSecretsGetCredentials_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets/1/credentials": testutils.RawHandler(http.StatusOK, `{"secretValue":"password123"}`),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsGetCredentials(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceDefaultSecretsGetCredentials_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets/99/credentials": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsGetCredentials(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestNetworkDeviceDefaultSecretsGetCredentials_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsGetCredentials(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

// --- NetworkDeviceDefaultSecretsCreate ---

func TestNetworkDeviceDefaultSecretsCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, ndSecretsItem)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsCreate(ctx, 1, "AA:BB:CC:DD:EE:01", "admin-password", "secret123"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceDefaultSecretsCreate_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets": testutils.ErrorHandler(http.StatusBadRequest, "validation error"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsCreate(ctx, 1, "AA:BB:CC:DD:EE:01", "admin-password", "secret123"); err == nil {
		t.Error("expected error for 400, got nil")
	}
}

// --- NetworkDeviceDefaultSecretsUpdate ---

func TestNetworkDeviceDefaultSecretsUpdate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets/1": func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPatch:
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, ndSecretsItem)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsUpdate(ctx, "1", "newpassword"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceDefaultSecretsUpdate_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsUpdate(ctx, "99", "newpassword"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestNetworkDeviceDefaultSecretsUpdate_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsUpdate(ctx, "not-a-number", "newpassword"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

// --- NetworkDeviceDefaultSecretsDelete ---

func TestNetworkDeviceDefaultSecretsDelete_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets/1": func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsDelete(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkDeviceDefaultSecretsDelete_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-devices/default-secrets/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsDelete(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestNetworkDeviceDefaultSecretsDelete_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := testutils.SetupTestContext(ts.URL)
	if err := NetworkDeviceDefaultSecretsDelete(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

func TestNetworkDeviceConfigExample(t *testing.T) {
	ctx := testutils.SetupTestContext("http://localhost")
	if err := NetworkDeviceConfigExample(ctx); err != nil {
		t.Errorf("NetworkDeviceConfigExample() unexpected error: %v", err)
	}
}
