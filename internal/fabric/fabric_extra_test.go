package fabric

import (
	"net/http"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// ---- JSON fixtures ----

// siteItem: minimal Site with all SDK-required fields (id, revision, slug, name).
const siteItem = `{"id":10,"revision":1,"slug":"site-1","name":"site-1"}`
const siteListBody = `{"data":[` + siteItem + `],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`

// ndItem2: minimal NetworkDevice for fabric tests (siteId=10 to match siteItem).
const ndItem2 = `{
	"id":"5","revision":1,"status":"active","siteId":10,
	"identifierString":"switch-a","description":"","chassisIdentifier":"",
	"country":"","city":"","datacenterMeta":"","datacenterRoom":"","datacenterRack":"",
	"rackPositionUpperUnit":0,"rackPositionLowerUnit":0,
	"managementAddress":"10.0.0.1","managementAddressPrefixLength":24,
	"managementAddressGateway":"10.0.0.254","managementPort":22,
	"syslogEnabled":0,"username":"admin","managementMacAddress":"AA:BB:CC:DD:EE:01",
	"serialNumber":"SN001","driver":"sonic_enterprise","position":"leaf",
	"orderIndex":1,"tags":[],"readyForInitialConfiguration":0,
	"bootstrapReadinessCheckInProgress":0,"subnetOobId":0,"subnetOobIndex":0,
	"requiresOsInstall":false,"bootstrapExpectedPartnerHostname":"",
	"loopbackAddressIpv6":"","asn":65000,"vtepAddressIpv6":"",
	"mlagSystemMac":"","mlagDomainId":0,"quarantineVlan":0,
	"variablesMaterializedForOSAssets":{},"secretsMaterializedForOSAssets":{},
	"bootstrapReadinessCheckResult":{},"isGateway":false
}`
const ndListBody = `{"data":[` + ndItem2 + `],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`

// fabricLinkItem: minimal NetworkFabricLink with all SDK-required fields.
const fabricLinkItem = `{"id":1,"networkFabricId":1,"linkType":"fabric","status":"active","revision":"1","createdTimestamp":"2024-01-01T00:00:00Z","updatedTimestamp":"2024-01-01T00:00:00Z"}`
const fabricLinkListBody = `{"data":[` + fabricLinkItem + `],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`

// jobInfoItem: minimal JobInfo (only jobGroupId is required).
const jobInfoItem = `{"jobGroupId":1}`

// fabricRoutes returns routes that serve fabricJSON at id=1 for both GET-by-id and list.
func fabricRoutes() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/api/v2/network-fabrics/1": testutils.RawHandler(http.StatusOK, fabricJSON),
		"/api/v2/network-fabrics":   testutils.RawHandler(http.StatusOK, `{"data":[`+fabricJSON+`],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
	}
}

// ---- GetFabricByIdOrLabel ----

func TestGetFabricByIdOrLabel_ById(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-fabrics/1": testutils.RawHandler(http.StatusOK, fabricJSON),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	fabric, err := GetFabricByIdOrLabel(ctx, "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fabric == nil {
		t.Fatal("expected fabric, got nil")
	}
}

func TestGetFabricByIdOrLabel_ByLabel(t *testing.T) {
	listBody := `{"data":[` + fabricJSON + `],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-fabrics": testutils.RawHandler(http.StatusOK, listBody),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	fabric, err := GetFabricByIdOrLabel(ctx, "fabric-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fabric == nil {
		t.Fatal("expected fabric, got nil")
	}
}

// ---- FabricCreate ----

func TestFabricCreate_HappyPath_Ethernet(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/sites":           testutils.RawHandler(http.StatusOK, siteListBody),
		"/api/v2/network-fabrics": testutils.RawHandler(http.StatusOK, fabricJSON),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := []byte(`{"fabricType":"ethernet"}`)
	if err := FabricCreate(ctx, "site-1", "test-fabric", "ethernet", "desc", config); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestFabricCreate_HappyPath_Infiniband(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/sites":           testutils.RawHandler(http.StatusOK, siteListBody),
		"/api/v2/network-fabrics": testutils.RawHandler(http.StatusOK, fabricJSON),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := []byte(`{"fabricType":"infiniband"}`)
	if err := FabricCreate(ctx, "site-1", "test-fabric", "infiniband", "desc", config); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestFabricCreate_HappyPath_FibreChannel(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/sites":           testutils.RawHandler(http.StatusOK, siteListBody),
		"/api/v2/network-fabrics": testutils.RawHandler(http.StatusOK, fabricJSON),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := []byte(`{"fabricType":"fibre_channel"}`)
	if err := FabricCreate(ctx, "site-1", "test-fabric", "fibre_channel", "desc", config); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestFabricCreate_InvalidType(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/sites": testutils.RawHandler(http.StatusOK, siteListBody),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricCreate(ctx, "site-1", "test-fabric", "bogus", "desc", []byte(`{}`)); err == nil {
		t.Error("expected error for invalid fabric type, got nil")
	}
}

func TestFabricCreate_SiteNotFound(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/sites": testutils.RawHandler(http.StatusOK, `{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricCreate(ctx, "no-such-site", "test-fabric", "ethernet", "desc", []byte(`{}`)); err == nil {
		t.Error("expected error for missing site, got nil")
	}
}

func TestFabricCreate_APIError(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/sites":           testutils.RawHandler(http.StatusOK, siteListBody),
		"/api/v2/network-fabrics": testutils.ErrorHandler(http.StatusInternalServerError, "server error"),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := []byte(`{"fabricType":"ethernet"}`)
	if err := FabricCreate(ctx, "site-1", "test-fabric", "ethernet", "desc", config); err == nil {
		t.Error("expected error for API failure, got nil")
	}
}

// ---- FabricUpdate ----
// fabricJSON has fabricConfiguration={"fabricType":"ethernet"} which the SDK
// unmarshals as EthernetFabric (discriminator-based oneOf).

func TestFabricUpdate_HappyPath_Ethernet(t *testing.T) {
	var callCount int
	routes := map[string]http.HandlerFunc{
		"/api/v2/network-fabrics/1": func(w http.ResponseWriter, r *http.Request) {
			callCount++
			testutils.RawHandler(http.StatusOK, fabricJSON)(w, r)
		},
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := []byte(`{"fabricType":"ethernet"}`)
	if err := FabricUpdate(ctx, "1", "new-name", "new-desc", config); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestFabricUpdate_FabricNotFound(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/network-fabrics/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
		"/api/v2/network-fabrics":    testutils.RawHandler(http.StatusOK, `{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricUpdate(ctx, "99", "name", "desc", []byte(`{}`)); err == nil {
		t.Error("expected error for missing fabric, got nil")
	}
}

func TestFabricUpdate_APIError(t *testing.T) {
	var callCount int
	handler := func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if r.Method == http.MethodGet {
			testutils.RawHandler(http.StatusOK, fabricJSON)(w, r)
			return
		}
		testutils.ErrorHandler(http.StatusInternalServerError, "server error")(w, r)
	}
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/network-fabrics/1": handler,
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := []byte(`{"fabricType":"ethernet"}`)
	if err := FabricUpdate(ctx, "1", "name", "desc", config); err == nil {
		t.Error("expected error for PATCH/PUT failure, got nil")
	}
}

// ---- FabricActivate ----

func TestFabricActivate_HappyPath(t *testing.T) {
	routes := fabricRoutes()
	routes["/api/v2/network-fabrics/1/actions/activate"] = testutils.RawHandler(http.StatusOK, fabricJSON)
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricActivate(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestFabricActivate_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(nil)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricActivate(ctx, "not-a-number"); err == nil {
		t.Error("expected error for non-numeric id, got nil")
	}
}

func TestFabricActivate_APIError(t *testing.T) {
	routes := fabricRoutes()
	routes["/api/v2/network-fabrics/1/actions/activate"] = testutils.ErrorHandler(http.StatusInternalServerError, "server error")
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricActivate(ctx, "1"); err == nil {
		t.Error("expected error for API failure, got nil")
	}
}

// ---- FabricDeploy ----

func TestFabricDeploy_HappyPath(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/network-fabrics/1/actions/deploy": testutils.RawHandler(http.StatusOK, jobInfoItem),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricDeploy(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestFabricDeploy_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(nil)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricDeploy(ctx, "not-a-number"); err == nil {
		t.Error("expected error for non-numeric id, got nil")
	}
}

func TestFabricDeploy_APIError(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/network-fabrics/1/actions/deploy": testutils.ErrorHandler(http.StatusInternalServerError, "server error"),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricDeploy(ctx, "1"); err == nil {
		t.Error("expected error for API failure, got nil")
	}
}

// ---- FabricDevicesGet ----

func TestFabricDevicesGet_HappyPath(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/network-fabrics/1/network-devices": testutils.RawHandler(http.StatusOK, ndListBody),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricDevicesGet(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestFabricDevicesGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(nil)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricDevicesGet(ctx, "not-a-number"); err == nil {
		t.Error("expected error for non-numeric id, got nil")
	}
}

func TestFabricDevicesGet_APIError(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/network-fabrics/1/network-devices": testutils.ErrorHandler(http.StatusInternalServerError, "server error"),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricDevicesGet(ctx, "1"); err == nil {
		t.Error("expected error for API failure, got nil")
	}
}

// ---- FabricDevicesAdd ----

// fabricWithSite: fabricJSON extended with siteId:10 to match ndItem2.
const fabricWithSite = `{"id":"1","name":"fabric-1","revision":"1","siteId":10,"createdTimestamp":"2024-01-01T00:00:00Z","updatedTimestamp":"2024-01-01T00:00:00Z","fabricConfiguration":{"fabricType":"ethernet"}}`

func TestFabricDevicesAdd_HappyPath(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/network-fabrics/1": testutils.RawHandler(http.StatusOK, fabricWithSite),
		"/api/v2/network-fabrics":   testutils.RawHandler(http.StatusOK, `{"data":[`+fabricWithSite+`],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
		"/api/v2/network-devices/5": testutils.RawHandler(http.StatusOK, ndItem2),
		"/api/v2/network-devices":   testutils.RawHandler(http.StatusOK, ndListBody),
		// POST to /api/v2/network-fabrics/1/network-devices returns the updated fabric
		"/api/v2/network-fabrics/1/network-devices": testutils.RawHandler(http.StatusOK, fabricWithSite),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricDevicesAdd(ctx, "1", []string{"5"}); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestFabricDevicesAdd_FabricNotFound(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/network-fabrics/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
		"/api/v2/network-fabrics":    testutils.RawHandler(http.StatusOK, `{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricDevicesAdd(ctx, "99", []string{"5"}); err == nil {
		t.Error("expected error for missing fabric, got nil")
	}
}

func TestFabricDevicesAdd_DeviceNotFound(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/network-fabrics/1": testutils.RawHandler(http.StatusOK, fabricWithSite),
		"/api/v2/network-fabrics":   testutils.RawHandler(http.StatusOK, `{"data":[`+fabricWithSite+`],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
		"/api/v2/network-devices/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
		"/api/v2/network-devices":    testutils.RawHandler(http.StatusOK, `{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricDevicesAdd(ctx, "1", []string{"99"}); err == nil {
		t.Error("expected error for missing device, got nil")
	}
}

// ---- FabricDevicesRemove ----

func TestFabricDevicesRemove_HappyPath(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/network-fabrics/1":             testutils.RawHandler(http.StatusOK, fabricWithSite),
		"/api/v2/network-fabrics":               testutils.RawHandler(http.StatusOK, `{"data":[`+fabricWithSite+`],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
		"/api/v2/network-devices/5":             testutils.RawHandler(http.StatusOK, ndItem2),
		"/api/v2/network-devices":               testutils.RawHandler(http.StatusOK, ndListBody),
		"/api/v2/network-fabrics/1/network-devices/5": testutils.RawHandler(http.StatusOK, fabricWithSite),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricDevicesRemove(ctx, "1", "5"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestFabricDevicesRemove_FabricNotFound(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/network-fabrics/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
		"/api/v2/network-fabrics":    testutils.RawHandler(http.StatusOK, `{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricDevicesRemove(ctx, "99", "5"); err == nil {
		t.Error("expected error for missing fabric, got nil")
	}
}

// ---- FabricLinksGet ----

func TestFabricLinksGet_HappyPath(t *testing.T) {
	routes := fabricRoutes()
	routes["/api/v2/network-fabrics/1/links"] = testutils.RawHandler(http.StatusOK, fabricLinkListBody)
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricLinksGet(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestFabricLinksGet_FabricNotFound(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/network-fabrics/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
		"/api/v2/network-fabrics":    testutils.RawHandler(http.StatusOK, `{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricLinksGet(ctx, "99"); err == nil {
		t.Error("expected error for missing fabric, got nil")
	}
}

func TestFabricLinksGet_APIError(t *testing.T) {
	routes := fabricRoutes()
	routes["/api/v2/network-fabrics/1/links"] = testutils.ErrorHandler(http.StatusInternalServerError, "server error")
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricLinksGet(ctx, "1"); err == nil {
		t.Error("expected error for API failure, got nil")
	}
}

// ---- FabricLinkAdd ----

func TestFabricLinkAdd_HappyPath(t *testing.T) {
	routes := fabricRoutes()
	routes["/api/v2/network-fabrics/1/links"] = testutils.RawHandler(http.StatusOK, fabricLinkItem)
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	ifA := float32(10)
	ifB := float32(20)
	createLink := sdk.CreateNetworkFabricLink{
		LinkType:                  "fabric",
		NetworkDeviceAInterfaceId: &ifA,
		NetworkDeviceBInterfaceId: &ifB,
	}
	if err := FabricLinkAdd(ctx, "1", createLink); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestFabricLinkAdd_FabricNotFound(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/network-fabrics/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
		"/api/v2/network-fabrics":    testutils.RawHandler(http.StatusOK, `{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricLinkAdd(ctx, "99", sdk.CreateNetworkFabricLink{}); err == nil {
		t.Error("expected error for missing fabric, got nil")
	}
}

func TestFabricLinkAdd_APIError(t *testing.T) {
	routes := fabricRoutes()
	routes["/api/v2/network-fabrics/1/links"] = testutils.ErrorHandler(http.StatusInternalServerError, "server error")
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricLinkAdd(ctx, "1", sdk.CreateNetworkFabricLink{}); err == nil {
		t.Error("expected error for API failure, got nil")
	}
}

// ---- FabricLinkRemove ----

func TestFabricLinkRemove_HappyPath(t *testing.T) {
	routes := fabricRoutes()
	routes["/api/v2/network-fabrics/1/links/7"] = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricLinkRemove(ctx, "1", "7"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestFabricLinkRemove_FabricNotFound(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v2/network-fabrics/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
		"/api/v2/network-fabrics":    testutils.RawHandler(http.StatusOK, `{"data":[],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
	}
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricLinkRemove(ctx, "99", "7"); err == nil {
		t.Error("expected error for missing fabric, got nil")
	}
}

func TestFabricLinkRemove_InvalidLinkId(t *testing.T) {
	routes := fabricRoutes()
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricLinkRemove(ctx, "1", "not-a-number"); err == nil {
		t.Error("expected error for non-numeric link id, got nil")
	}
}

func TestFabricLinkRemove_APIError(t *testing.T) {
	routes := fabricRoutes()
	routes["/api/v2/network-fabrics/1/links/7"] = testutils.ErrorHandler(http.StatusInternalServerError, "server error")
	ts := testutils.NewTestServer(routes)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := FabricLinkRemove(ctx, "1", "7"); err == nil {
		t.Error("expected error for API failure, got nil")
	}
}
