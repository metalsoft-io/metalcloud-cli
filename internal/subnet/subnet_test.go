package subnet

import (
	"context"
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

const subnetIpsLinksArrayResponse = `{
	"data": [
		{
			"id": 1,
			"name": "ip-1",
			"address": "10.0.0.1",
			"ipVersion": "IPv4",
			"subnetId": 1,
			"links": []
		}
	]
}`

const subnetIpsLinksMapResponse = `{
	"data": [
		{
			"id": 1,
			"name": "ip-1",
			"address": "10.0.0.1",
			"ipVersion": "IPv4",
			"subnetId": 1,
			"links": {"self": "http://example.com/api/v2/subnets/1/ips/1"}
		}
	]
}`

const subnetIpRangesLinksArrayResponse = `{
	"data": [
		{
			"id": 1,
			"name": "range-1",
			"startAddress": "10.0.0.1",
			"endAddress": "10.0.0.254",
			"ipVersion": "IPv4",
			"subnetId": 1,
			"links": []
		}
	]
}`

const subnetIpRangesLinksMapResponse = `{
	"data": [
		{
			"id": 1,
			"name": "range-1",
			"startAddress": "10.0.0.1",
			"endAddress": "10.0.0.254",
			"ipVersion": "IPv4",
			"subnetId": 1,
			"links": {"self": "http://example.com/api/v2/subnets/1/ip-ranges/1"}
		}
	]
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

func TestSubnetIps_LinksAsArray(t *testing.T) {
	t.Run("LinksAsArray", func(t *testing.T) {
		ts := newSubnetMockServer(subnetIpsLinksArrayResponse, subnetIpRangesLinksArrayResponse)
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		if err := SubnetIps(ctx, "1"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})
}

func TestSubnetIps_LinksAsMap(t *testing.T) {
	t.Run("LinksAsMap", func(t *testing.T) {
		ts := newSubnetMockServer(subnetIpsLinksMapResponse, subnetIpRangesLinksMapResponse)
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		if err := SubnetIps(ctx, "1"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})
}

func TestSubnetIpRanges_LinksAsArray(t *testing.T) {
	t.Run("LinksAsArray", func(t *testing.T) {
		ts := newSubnetMockServer(subnetIpsLinksArrayResponse, subnetIpRangesLinksArrayResponse)
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		if err := SubnetIpRanges(ctx, "1"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})
}

func TestSubnetIpRanges_LinksAsMap(t *testing.T) {
	t.Run("LinksAsMap", func(t *testing.T) {
		ts := newSubnetMockServer(subnetIpsLinksMapResponse, subnetIpRangesLinksMapResponse)
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		if err := SubnetIpRanges(ctx, "1"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})
}
