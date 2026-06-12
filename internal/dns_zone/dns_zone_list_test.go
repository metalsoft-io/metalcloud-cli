package dns_zone

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// dnsZoneJSON has all required fields.
const dnsZoneJSON = `{"id":1,"label":"zone-1","zoneName":"example.com","zoneType":"master","soaEmail":"admin@example.com","soaSerial":1,"ttl":300,"nameServers":["ns1.example.com"],"isDefault":false,"status":"active","revision":1,"createdBy":1,"createdAt":"2024-01-01T00:00:00Z"}`

func dnsZoneListHandler(statusCode int, count, currentPage, totalPages int) http.HandlerFunc {
	items := make([]string, count)
	for i := range items {
		items[i] = dnsZoneJSON
	}
	data := "[" + strings.Join(items, ",") + "]"
	body := fmt.Sprintf(`{"data":%s,"meta":{"currentPage":%d,"totalPages":%d,"itemsPerPage":100}}`,
		data, currentPage, totalPages)
	return testutils.RawHandler(statusCode, body)
}

func dnsZonePage(n int) []any {
	items := make([]any, n)
	for i := range items {
		items[i] = map[string]any{
			"id": i + 1, "label": fmt.Sprintf("zone-%d", i+1),
			"zoneName": fmt.Sprintf("zone-%d.com", i+1), "zoneType": "master",
			"soaEmail": "admin@example.com", "soaSerial": 1, "ttl": 300,
			"nameServers": []string{"ns1.example.com"}, "isDefault": false,
			"status": "active", "revision": 1, "createdBy": 1,
			"createdAt": "2024-01-01T00:00:00Z",
		}
	}
	return items
}

func TestDNSZoneList_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/dns-zones": dnsZoneListHandler(http.StatusOK, 2, 1, 1),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := DNSZoneList(ctx, nil); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestDNSZoneList_Error(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/dns-zones": testutils.ErrorHandler(http.StatusInternalServerError, "internal error"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := DNSZoneList(ctx, nil); err == nil {
		t.Error("expected error for 500, got nil")
	}
}

func TestDNSZoneList_Empty(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/dns-zones": dnsZoneListHandler(http.StatusOK, 0, 1, 1),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := DNSZoneList(ctx, nil); err != nil {
		t.Errorf("expected nil error for empty list, got: %v", err)
	}
}

func TestDNSZoneList_Pagination(t *testing.T) {
	ts := testutils.MultiPageServer("/api/v2/dns-zones", []any{dnsZonePage(100), dnsZonePage(100), dnsZonePage(5)})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := DNSZoneList(ctx, nil); err != nil {
		t.Errorf("expected nil error across 3 pages, got: %v", err)
	}
}

func TestDNSZoneGet_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/dns-zones/1": testutils.RawHandler(http.StatusOK, dnsZoneJSON),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := DNSZoneGet(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestDNSZoneGet_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/dns-zones/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := DNSZoneGet(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestDNSZoneGet_InvalidId(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := DNSZoneGet(ctx, "not-a-number"); err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

func TestDNSZoneCreate_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/dns-zones": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, dnsZoneJSON)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := sdk.CreateDnsZone{
		ZoneName:    "example.com",
		IsDefault:   false,
		NameServers: []string{"ns1.example.com"},
	}
	if err := DNSZoneCreate(ctx, config); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestDNSZoneCreate_BadRequest(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/dns-zones": testutils.ErrorHandler(http.StatusBadRequest, "validation error"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	config := sdk.CreateDnsZone{ZoneName: "bad"}
	if err := DNSZoneCreate(ctx, config); err == nil {
		t.Error("expected error for 400, got nil")
	}
}

func TestDNSZoneDelete_HappyPath(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/dns-zones/1": func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Method == http.MethodGet {
				fmt.Fprint(w, dnsZoneJSON)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		},
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := DNSZoneDelete(ctx, "1"); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestDNSZoneDelete_NotFound(t *testing.T) {
	ts := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/dns-zones/99": testutils.ErrorHandler(http.StatusNotFound, "not found"),
	})
	defer ts.Close()

	ctx := setupTestContext(ts.URL)
	if err := DNSZoneDelete(ctx, "99"); err == nil {
		t.Error("expected error for 404, got nil")
	}
}
