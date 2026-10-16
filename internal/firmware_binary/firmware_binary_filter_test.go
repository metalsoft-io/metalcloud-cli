package firmware_binary

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/metalsoft-io/metalcloud-cli/internal/testutils"
)

// queryRecorder serves a single page of binaries and records every query string.
type queryRecorder struct {
	mu      sync.Mutex
	queries []url.Values
	total   int
}

func (q *queryRecorder) handler(w http.ResponseWriter, r *http.Request) {
	q.mu.Lock()
	q.queries = append(q.queries, r.URL.Query())
	q.mu.Unlock()

	items := make([]map[string]any, 0)
	for i := 1; i <= q.total && i <= 3; i++ {
		items = append(items, binaryFixture(i))
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(testutils.PaginatedResponse(items, 1, 1))
}

func newRecorder(t *testing.T, total int) (*queryRecorder, string, func()) {
	t.Helper()

	rec := &queryRecorder{total: total}
	srv := testutils.NewTestServer(map[string]http.HandlerFunc{
		"/api/v2/firmware/binary": rec.handler,
	})

	return rec, srv.URL, srv.Close
}

func TestFirmwareBinaryList_SendsServerSideFilters(t *testing.T) {
	rec, url, closeFn := newRecorder(t, 2)
	defer closeFn()

	err := FirmwareBinaryList(testutils.SetupTestContext(url), ListFlags{
		FilterCatalogId:      []string{"8813"},
		FilterUpdateSeverity: []string{"critical"},
		FilterRebootRequired: []string{"false"},
		FilterPackageVersion: []string{"2.15.0"},
		FilterExternalId:     []string{"abc", "def"},
	})
	if err != nil {
		t.Fatalf("FirmwareBinaryList() unexpected error: %v", err)
	}

	if len(rec.queries) != 1 {
		t.Fatalf("expected a single request, got %d", len(rec.queries))
	}

	expect := map[string][]string{
		"filter.catalogId":      {"8813"},
		"filter.updateSeverity": {"critical"},
		"filter.rebootRequired": {"false"},
		"filter.packageVersion": {"2.15.0"},
		"filter.externalId":     {"$eq:abc", "$or:$eq:def"},
	}
	for key, want := range expect {
		if got := rec.queries[0][key]; strings.Join(got, "|") != strings.Join(want, "|") {
			t.Errorf("%s: expected %v, got %v", key, want, got)
		}
	}
}

func TestFirmwareBinaryList_NoFilterParamsByDefault(t *testing.T) {
	rec, url, closeFn := newRecorder(t, 2)
	defer closeFn()

	if err := FirmwareBinaryList(testutils.SetupTestContext(url), ListFlags{}); err != nil {
		t.Fatalf("FirmwareBinaryList() unexpected error: %v", err)
	}

	for key := range rec.queries[0] {
		if strings.HasPrefix(key, "filter.") {
			t.Errorf("unexpected filter parameter %q without any filter flag", key)
		}
	}
}

func TestFirmwareBinaryList_PageWindow(t *testing.T) {
	rec, url, closeFn := newRecorder(t, 3)
	defer closeFn()

	if err := FirmwareBinaryList(testutils.SetupTestContext(url), ListFlags{Page: 2, Limit: 50}); err != nil {
		t.Fatalf("FirmwareBinaryList() unexpected error: %v", err)
	}

	if len(rec.queries) != 1 {
		t.Fatalf("expected a single request for page 2 of 50, got %d", len(rec.queries))
	}
	if rec.queries[0].Get("page") != "2" || rec.queries[0].Get("limit") != "50" {
		t.Errorf("expected page=2 limit=50, got page=%s limit=%s", rec.queries[0].Get("page"), rec.queries[0].Get("limit"))
	}
}

func TestFirmwareBinaryList_LimitOnly(t *testing.T) {
	rec, url, closeFn := newRecorder(t, 3)
	defer closeFn()

	if err := FirmwareBinaryList(testutils.SetupTestContext(url), ListFlags{Limit: 10}); err != nil {
		t.Fatalf("FirmwareBinaryList() unexpected error: %v", err)
	}

	if len(rec.queries) != 1 {
		t.Fatalf("expected a single request for --limit 10, got %d", len(rec.queries))
	}
	if rec.queries[0].Get("page") != "1" || rec.queries[0].Get("limit") != "10" {
		t.Errorf("expected page=1 limit=10, got page=%s limit=%s", rec.queries[0].Get("page"), rec.queries[0].Get("limit"))
	}
}

func TestFirmwareBinaryConfigExample_HasNoVendorField(t *testing.T) {
	out := testutils.CaptureStdout(t, func() {
		if err := FirmwareBinaryConfigExample(testutils.SetupTestContext("http://localhost")); err != nil {
			t.Fatalf("FirmwareBinaryConfigExample() unexpected error: %v", err)
		}
	})

	var example map[string]any
	if err := json.Unmarshal([]byte(out), &example); err != nil {
		t.Fatalf("config example is not valid JSON: %v\n%s", err, out)
	}

	// The API no longer accepts "vendor" on create.
	if _, ok := example["vendor"]; ok {
		t.Errorf("config example must not contain the vendor field: %s", out)
	}
	if example["vendorDownloadUrl"] == nil || example["name"] == nil {
		t.Errorf("config example lost its other fields: %s", out)
	}
}
