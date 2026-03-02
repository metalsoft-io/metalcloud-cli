package dns_zone

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

const dnsRecordsLinksArrayResponse = `{
	"data": [
		{
			"id": 1,
			"name": "www",
			"type": "A",
			"records": ["10.0.0.1"],
			"ttl": 300,
			"status": "active",
			"zoneName": "example.com",
			"siteId": 1,
			"infrastructureId": 1,
			"zoneId": 1,
			"revision": 1,
			"createdBy": 1,
			"createdAt": "2024-01-01T00:00:00Z",
			"links": []
		}
	],
	"meta": {"itemsPerPage": 10}
}`

const dnsRecordsLinksMapResponse = `{
	"data": [
		{
			"id": 1,
			"name": "www",
			"type": "A",
			"records": ["10.0.0.1"],
			"ttl": 300,
			"status": "active",
			"zoneName": "example.com",
			"siteId": 1,
			"infrastructureId": 1,
			"zoneId": 1,
			"revision": 1,
			"createdBy": 1,
			"createdAt": "2024-01-01T00:00:00Z",
			"links": [{"rel": "self", "href": "http://example.com/api/v2/dns-zones/1/recordsets/1"}]
		}
	],
	"meta": {"itemsPerPage": 10}
}`

func newDNSZoneMockServer(body string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		path := r.URL.Path
		if strings.Contains(path, "/dns-record-sets") || strings.Contains(path, "/record-sets") || strings.Contains(path, "/recordsets") {
			w.WriteHeader(statusCode)
			fmt.Fprint(w, body)
			return
		}
		w.WriteHeader(statusCode)
		fmt.Fprint(w, body)
	}))
}

func TestDNSZoneRecords_LinksAsArray(t *testing.T) {
	t.Run("LinksAsArray", func(t *testing.T) {
		ts := newDNSZoneMockServer(dnsRecordsLinksArrayResponse, http.StatusOK)
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		if err := DNSZoneRecords(ctx, "1"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})
}

func TestDNSZoneRecords_LinksAsMap(t *testing.T) {
	t.Run("LinksAsMap", func(t *testing.T) {
		ts := newDNSZoneMockServer(dnsRecordsLinksMapResponse, http.StatusOK)
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		if err := DNSZoneRecords(ctx, "1"); err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})
}

func TestDNSZoneRecords_HttpError(t *testing.T) {
	t.Run("HttpError", func(t *testing.T) {
		ts := newDNSZoneMockServer(`{"message":"not found"}`, http.StatusNotFound)
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		if err := DNSZoneRecords(ctx, "1"); err == nil {
			t.Error("expected an error for HTTP 404, got nil")
		}
	})
}
