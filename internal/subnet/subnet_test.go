package subnet

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	// formatter.PrintResult requires a valid format; use "json" so it serialises
	// to stdout without needing a TTY or table renderer.
	viper.Set("format", "json")
	m.Run()
}

func setupTestContext(serverURL string) context.Context {
	cfg := sdk.NewConfiguration()
	cfg.Servers = []sdk.ServerConfiguration{{URL: serverURL}}
	client := sdk.NewAPIClient(cfg)
	ctx := context.WithValue(context.Background(), api.ApiClientContextKey, client)
	ctx = context.WithValue(ctx, sdk.ContextAccessToken, "test-api-key")
	return ctx
}

const subnetIpsResponse = `{
	"data": [
		{
			"id": 1,
			"name": "ip-1",
			"address": "10.0.0.1",
			"ipVersion": "ipv4",
			"subnetId": 1,
			"annotations": {},
			"createdAt": "2024-01-01T00:00:00Z",
			"updatedAt": "2024-01-01T00:00:00Z",
			"revision": 1,
			"tags": {}
		}
	],
	"meta": {"itemsPerPage": 10}
}`

const subnetIpRangesResponse = `{
	"data": [
		{
			"id": 1,
			"name": "range-1",
			"startAddress": "10.0.0.1",
			"endAddress": "10.0.0.254",
			"ipVersion": "ipv4",
			"subnetId": 1,
			"annotations": {},
			"createdAt": "2024-01-01T00:00:00Z",
			"updatedAt": "2024-01-01T00:00:00Z",
			"revision": 1,
			"tags": {}
		}
	],
	"meta": {"itemsPerPage": 10}
}`

func newSubnetMockServer(ipsBody, ipRangesBody string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		path := r.URL.Path
		switch {
		case strings.Contains(path, "/ip-ranges"):
			fmt.Fprint(w, ipRangesBody)
		case strings.Contains(path, "/ips"):
			fmt.Fprint(w, ipsBody)
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestSubnetIps(t *testing.T) {
	ts := newSubnetMockServer(subnetIpsResponse, subnetIpRangesResponse)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SubnetIps(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestSubnetIpRanges(t *testing.T) {
	ts := newSubnetMockServer(subnetIpsResponse, subnetIpRangesResponse)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SubnetIpRanges(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// capacity, remove-ip and remove-ip-range
// ---------------------------------------------------------------------------

// subnetItem satisfies the required properties of sdk.Subnet. Its revision (4)
// is the entity tag the remove commands must echo back in If-Match.
var subnetItem = map[string]any{
	"id": 1, "label": "subnet-1", "name": "subnet-1", "revision": 4,
	"ipVersion": "ipv4", "networkAddress": "10.0.0.0", "prefixLength": 24,
	"netmask": "255.255.255.0", "isPool": false,
	"parentSubnetId": nil, "defaultGatewayAddress": "10.0.0.1",
	"allocationDenylist": []any{}, "childOverlapAllowRules": []any{},
	"annotations": map[string]any{}, "tags": map[string]any{},
	"createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z",
}

var subnetCapacityItem = map[string]any{
	"subnetId": 1, "prefix": "10.0.0.0/24", "isPool": false,
	"totalIpCount": "254", "freeIpCount": "200", "usedIpCount": "54",
	"freePrefixes": nil, "freePrefixCounts": nil,
	"requestedPrefixLength": nil, "allocatablePrefixCount": nil,
}

type subnetRecordedRequest struct {
	Method string
	Path   string
	Header http.Header
}

// newSubnetCapacityServer serves the subnet, its capacity and the two delete
// sub-resources, recording every request for assertions.
func newSubnetCapacityServer(reqs *[]subnetRecordedRequest) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*reqs = append(*reqs, subnetRecordedRequest{Method: r.Method, Path: r.URL.Path, Header: r.Header.Clone()})
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v2/subnets/1/capacity":
			_ = json.NewEncoder(w).Encode(subnetCapacityItem)
		case "/api/v2/subnets/1/ips/7", "/api/v2/subnets/1/ip-ranges/9":
			w.WriteHeader(http.StatusNoContent)
		case "/api/v2/subnets/1":
			_ = json.NewEncoder(w).Encode(subnetItem)
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestSubnetCapacity(t *testing.T) {
	var reqs []subnetRecordedRequest
	ts := newSubnetCapacityServer(&reqs)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SubnetCapacity(ctx, "1", 0); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	last := reqs[len(reqs)-1]
	if last.Method != http.MethodGet || last.Path != "/api/v2/subnets/1/capacity" {
		t.Fatalf("unexpected request %s %s", last.Method, last.Path)
	}
}

func TestSubnetCapacity_InvalidId(t *testing.T) {
	ctx := setupTestContext("http://localhost")
	if err := SubnetCapacity(ctx, "not-a-number", 0); err == nil {
		t.Fatal("expected an error for an invalid subnet ID, got nil")
	}
}

func TestSubnetIpRemove_SendsIfMatch(t *testing.T) {
	var reqs []subnetRecordedRequest
	ts := newSubnetCapacityServer(&reqs)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SubnetIpRemove(ctx, "1", "7"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	last := reqs[len(reqs)-1]
	if last.Method != http.MethodDelete || last.Path != "/api/v2/subnets/1/ips/7" {
		t.Fatalf("unexpected request %s %s", last.Method, last.Path)
	}
	if got := last.Header.Get("If-Match"); got != "4" {
		t.Errorf("expected If-Match 4 (the subnet revision), got %q", got)
	}
}

func TestSubnetIpRangeRemove_SendsIfMatch(t *testing.T) {
	var reqs []subnetRecordedRequest
	ts := newSubnetCapacityServer(&reqs)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SubnetIpRangeRemove(ctx, "1", "9"); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	last := reqs[len(reqs)-1]
	if last.Method != http.MethodDelete || last.Path != "/api/v2/subnets/1/ip-ranges/9" {
		t.Fatalf("unexpected request %s %s", last.Method, last.Path)
	}
	if got := last.Header.Get("If-Match"); got != "4" {
		t.Errorf("expected If-Match 4 (the subnet revision), got %q", got)
	}
}

func TestSubnetIpRemove_InvalidIpId(t *testing.T) {
	var reqs []subnetRecordedRequest
	ts := newSubnetCapacityServer(&reqs)
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := SubnetIpRemove(ctx, "1", "not-a-number"); err == nil {
		t.Fatal("expected an error for an invalid IP ID, got nil")
	}
	if err := SubnetIpRangeRemove(ctx, "1", "not-a-number"); err == nil {
		t.Fatal("expected an error for an invalid IP range ID, got nil")
	}
}
