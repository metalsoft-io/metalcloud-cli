package dns_zone

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const dnsRecordSetSingle = `{
	"id": 456,
	"name": "www",
	"type": "A",
	"records": ["10.0.0.1"],
	"ttl": 300,
	"status": "active",
	"zoneName": "example.com",
	"siteId": 1,
	"infrastructureId": 1,
	"zoneId": 123,
	"revision": 1,
	"createdBy": 1,
	"createdAt": "2024-01-01T00:00:00Z",
	"links": []
}`

func dnsExtrasServer(t *testing.T, routes map[string]http.HandlerFunc) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	for path, handler := range routes {
		mux.HandleFunc(path, handler)
	}
	return httptest.NewServer(mux)
}

func dnsExtrasJSON(status int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}
}

func TestDNSZoneNameservers(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		ts := dnsExtrasServer(t, map[string]http.HandlerFunc{
			"/api/v2/dns-zones/123/nameservers": dnsExtrasJSON(http.StatusOK, `["ns1.example.com","ns2.example.com"]`),
		})
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		if err := DNSZoneNameservers(ctx, "123"); err != nil {
			t.Errorf("DNSZoneNameservers: expected nil error, got: %v", err)
		}
	})

	t.Run("Empty", func(t *testing.T) {
		ts := dnsExtrasServer(t, map[string]http.HandlerFunc{
			"/api/v2/dns-zones/123/nameservers": dnsExtrasJSON(http.StatusOK, `[]`),
		})
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		if err := DNSZoneNameservers(ctx, "123"); err != nil {
			t.Errorf("DNSZoneNameservers empty: expected nil error, got: %v", err)
		}
	})

	t.Run("InvalidId", func(t *testing.T) {
		ctx := setupTestContext("http://127.0.0.1:1")
		if err := DNSZoneNameservers(ctx, "not-a-number"); err == nil {
			t.Error("expected error for invalid zone ID, got nil")
		}
	})

	t.Run("HttpError", func(t *testing.T) {
		ts := dnsExtrasServer(t, map[string]http.HandlerFunc{
			"/api/v2/dns-zones/123/nameservers": dnsExtrasJSON(http.StatusNotFound, `{"message":"not found"}`),
		})
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		if err := DNSZoneNameservers(ctx, "123"); err == nil {
			t.Error("expected error for HTTP 404, got nil")
		}
	})
}

func TestDNSZoneRecord(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		ts := dnsExtrasServer(t, map[string]http.HandlerFunc{
			"/api/v2/dns-zones/123/recordsets/456": dnsExtrasJSON(http.StatusOK, dnsRecordSetSingle),
		})
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		if err := DNSZoneRecord(ctx, "123", "456"); err != nil {
			t.Errorf("DNSZoneRecord: expected nil error, got: %v", err)
		}
	})

	t.Run("InvalidRecordSetId", func(t *testing.T) {
		ctx := setupTestContext("http://127.0.0.1:1")
		if err := DNSZoneRecord(ctx, "123", "not-a-number"); err == nil {
			t.Error("expected error for invalid record set ID, got nil")
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		ts := dnsExtrasServer(t, map[string]http.HandlerFunc{
			"/api/v2/dns-zones/123/recordsets/999": dnsExtrasJSON(http.StatusNotFound, `{"message":"not found"}`),
		})
		defer ts.Close()

		ctx := setupTestContext(ts.URL)
		if err := DNSZoneRecord(ctx, "123", "999"); err == nil {
			t.Error("expected error for HTTP 404, got nil")
		}
	})
}

func TestDNSZoneRecords_GlobalList(t *testing.T) {
	// An empty zone ID lists every record set through the global endpoint.
	ts := dnsExtrasServer(t, map[string]http.HandlerFunc{
		"/api/v2/dns-recordsets": dnsExtrasJSON(http.StatusOK,
			`{"data":[`+dnsRecordSetSingle+`],"meta":{"currentPage":1,"totalPages":1,"itemsPerPage":100}}`),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := DNSZoneRecords(ctx, ""); err != nil {
		t.Errorf("DNSZoneRecords global: expected nil error, got: %v", err)
	}
}

func TestDNSZoneRecords_GlobalListError(t *testing.T) {
	ts := dnsExtrasServer(t, map[string]http.HandlerFunc{
		"/api/v2/dns-recordsets": dnsExtrasJSON(http.StatusInternalServerError, `{"message":"internal error"}`),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := DNSZoneRecords(ctx, ""); err == nil {
		t.Error("expected error for HTTP 500, got nil")
	}
}
