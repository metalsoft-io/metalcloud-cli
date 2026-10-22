package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// pagedServersHandler serves total servers, honouring the page and limit query
// parameters like the real API, and records every query it receives.
func pagedServersHandler(total int, mu *sync.Mutex, queries *[]url.Values) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || !strings.HasSuffix(r.URL.Path, "/servers") {
			http.NotFound(w, r)
			return
		}

		q := r.URL.Query()
		mu.Lock()
		*queries = append(*queries, q)
		mu.Unlock()

		page, _ := strconv.Atoi(q.Get("page"))
		if page < 1 {
			page = 1
		}
		limit, _ := strconv.Atoi(q.Get("limit"))
		if limit < 1 {
			limit = 20
		}

		start := (page - 1) * limit
		end := start + limit
		if end > total {
			end = total
		}

		items := make([]map[string]interface{}, 0)
		for i := start; i < end; i++ {
			items = append(items, map[string]interface{}{
				"serverId": i + 1, "serverUUID": fmt.Sprintf("uuid-%d", i+1), "serialNumber": fmt.Sprintf("SN-%d", i+1),
			})
		}

		totalPages := (total + limit - 1) / limit
		body := buildServersBody(items)
		body = strings.TrimSuffix(body, "}") + fmt.Sprintf(`,"meta":{"totalItems":%d,"totalPages":%d,"currentPage":%d}}`, total, totalPages, page)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, body)
	}
}

func TestServerList_MultiPage(t *testing.T) {
	var mu sync.Mutex
	var queries []url.Values

	ts := httptest.NewServer(pagedServersHandler(205, &mu, &queries))
	defer ts.Close()

	if err := ServerList(setupTestContext(ts.URL), ListFlags{}); err != nil {
		t.Fatalf("ServerList with 205 servers: expected nil error, got: %v", err)
	}

	// 205 servers at 100 per page must be fetched in three requests.
	if len(queries) != 3 {
		t.Fatalf("expected 3 page requests, got %d", len(queries))
	}
	for i, q := range queries {
		if q.Get("page") != strconv.Itoa(i+1) || q.Get("limit") != "100" {
			t.Errorf("request %d: unexpected page/limit %q/%q", i, q.Get("page"), q.Get("limit"))
		}
	}
}

func TestServerList_SendsServerSideFilters(t *testing.T) {
	var mu sync.Mutex
	var queries []url.Values

	ts := httptest.NewServer(pagedServersHandler(3, &mu, &queries))
	defer ts.Close()

	err := ServerList(setupTestContext(ts.URL), ListFlags{
		FilterTags:                       []string{"rack-12"},
		FilterSiteId:                     []string{"2", "3"},
		FilterStatus:                     []string{"active"},
		FilterType:                       []string{"7"},
		FilterFirmwareBaselineCompliance: []string{"$not:$eq:compliant"},
		FilterVendor:                     []string{"Dell"},
		FilterHealthState:                []string{"failed", "warnings_present"},
	})
	if err != nil {
		t.Fatalf("ServerList: unexpected error: %v", err)
	}
	if len(queries) != 1 {
		t.Fatalf("expected 1 request, got %d", len(queries))
	}

	q := queries[0]
	expect := map[string][]string{
		"filter.tags":                             {"rack-12"},
		"filter.siteId":                           {"$eq:2", "$or:$eq:3"},
		"filter.serverStatus":                     {"active"},
		"filter.serverTypeId":                     {"7"},
		"filter.firmwareBaselineComplianceStatus": {"$not:$eq:compliant"},
		"filter.vendor":                           {"Dell"},
		"filter.healthState":                      {"$eq:failed", "$or:$eq:warnings_present"},
	}
	for key, want := range expect {
		got := q[key]
		if strings.Join(got, "|") != strings.Join(want, "|") {
			t.Errorf("%s: expected %v, got %v", key, want, got)
		}
	}
}

func TestServerList_NoFilterParamsByDefault(t *testing.T) {
	var mu sync.Mutex
	var queries []url.Values

	ts := httptest.NewServer(pagedServersHandler(3, &mu, &queries))
	defer ts.Close()

	if err := ServerList(setupTestContext(ts.URL), ListFlags{}); err != nil {
		t.Fatalf("ServerList: unexpected error: %v", err)
	}
	for key := range queries[0] {
		if strings.HasPrefix(key, "filter.") {
			t.Errorf("unexpected filter parameter %q without any filter flag", key)
		}
	}
}

func TestServerList_PageWindow(t *testing.T) {
	var mu sync.Mutex
	var queries []url.Values

	ts := httptest.NewServer(pagedServersHandler(205, &mu, &queries))
	defer ts.Close()

	if err := ServerList(setupTestContext(ts.URL), ListFlags{Page: 2, Limit: 50}); err != nil {
		t.Fatalf("ServerList: unexpected error: %v", err)
	}
	if len(queries) != 1 {
		t.Fatalf("expected a single request for page 2 of 50, got %d", len(queries))
	}
	if queries[0].Get("page") != "2" || queries[0].Get("limit") != "50" {
		t.Errorf("expected page=2 limit=50, got page=%s limit=%s", queries[0].Get("page"), queries[0].Get("limit"))
	}
}

func TestServerList_LimitOnly(t *testing.T) {
	var mu sync.Mutex
	var queries []url.Values

	ts := httptest.NewServer(pagedServersHandler(205, &mu, &queries))
	defer ts.Close()

	if err := ServerList(setupTestContext(ts.URL), ListFlags{Limit: 10}); err != nil {
		t.Fatalf("ServerList: unexpected error: %v", err)
	}
	if len(queries) != 1 {
		t.Fatalf("expected a single request for --limit 10, got %d", len(queries))
	}
	if queries[0].Get("page") != "1" || queries[0].Get("limit") != "10" {
		t.Errorf("expected page=1 limit=10, got page=%s limit=%s", queries[0].Get("page"), queries[0].Get("limit"))
	}
}
