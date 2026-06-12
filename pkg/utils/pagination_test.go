package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

type fakeList struct {
	data []int
	meta sdk.PaginatedResponseMeta
}

func (l fakeList) GetData() []int                     { return l.data }
func (l fakeList) GetMeta() sdk.PaginatedResponseMeta { return l.meta }

type fakeRequest struct {
	pages  [][]int
	page   float32
	failOn float32
}

func (r fakeRequest) Page(p float32) fakeRequest  { r.page = p; return r }
func (r fakeRequest) Limit(l float32) fakeRequest { return r }

func (r fakeRequest) Execute() (fakeList, *http.Response, error) {
	if r.failOn > 0 && r.page == r.failOn {
		return fakeList{}, &http.Response{StatusCode: 500}, errors.New("boom")
	}
	current := int32(r.page)
	total := int32(len(r.pages))
	items := int32(0)
	for _, p := range r.pages {
		items += int32(len(p))
	}
	return fakeList{
		data: r.pages[int(r.page)-1],
		meta: sdk.PaginatedResponseMeta{
			CurrentPage: &current,
			TotalPages:  &total,
			TotalItems:  &items,
		},
	}, &http.Response{StatusCode: 200}, nil
}

func TestFetchAllPagesSinglePage(t *testing.T) {
	records, _, err := FetchAllPages(fakeRequest{pages: [][]int{{1, 2, 3}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 3 {
		t.Fatalf("expected 3 records, got %d", len(records))
	}
}

func TestFetchAllPagesMultiPage(t *testing.T) {
	records, _, err := FetchAllPages(fakeRequest{pages: [][]int{{1, 2}, {3, 4}, {5}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 5 {
		t.Fatalf("expected 5 records, got %d", len(records))
	}
	for i, v := range []int{1, 2, 3, 4, 5} {
		if records[i] != v {
			t.Fatalf("expected %d at index %d, got %d", v, i, records[i])
		}
	}
}

func TestFetchAllPagesErrorMidStream(t *testing.T) {
	_, _, err := FetchAllPages(fakeRequest{pages: [][]int{{1, 2}, {3, 4}}, failOn: 2})
	if err == nil {
		t.Fatal("expected error")
	}
}

// fakeRequestNoMeta does not set TotalPages so FetchAllPages uses the len<limit fallback.
type fakeRequestNoMeta struct {
	pages    [][]int
	page     float32
	execCalls *int
}

func (r fakeRequestNoMeta) Page(p float32) fakeRequestNoMeta  { r.page = p; return r }
func (r fakeRequestNoMeta) Limit(l float32) fakeRequestNoMeta { return r }

func (r fakeRequestNoMeta) Execute() (fakeList, *http.Response, error) {
	if r.execCalls != nil {
		*r.execCalls++
	}
	idx := int(r.page) - 1
	if idx < 0 || idx >= len(r.pages) {
		return fakeList{}, &http.Response{StatusCode: 200}, nil
	}
	current := int32(r.page)
	return fakeList{
		data: r.pages[idx],
		meta: sdk.PaginatedResponseMeta{CurrentPage: &current},
	}, &http.Response{StatusCode: 200}, nil
}

// TestFetchAllPagesExactly100Items — exactly 100 items on page 1, no TotalPages.
// The fallback branch breaks when len(batch) < defaultPageSize (100), so a full
// page of 100 would loop to page 2. We model that: page 1 has 100 items, page 2 has 0.
// Execute must be called exactly twice.
func TestFetchAllPagesExactly100Items(t *testing.T) {
	hundred := make([]int, 100)
	for i := range hundred {
		hundred[i] = i + 1
	}
	calls := 0
	req := fakeRequestNoMeta{
		pages:     [][]int{hundred, {}},
		execCalls: &calls,
	}
	records, _, err := FetchAllPages(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 100 {
		t.Fatalf("expected 100 records, got %d", len(records))
	}
	if calls != 2 {
		t.Fatalf("expected Execute called exactly 2 times, got %d", calls)
	}
}

// TestFetchAllPagesZeroItems — page 1 returns 0 items; loop must terminate immediately.
func TestFetchAllPagesZeroItems(t *testing.T) {
	calls := 0
	req := fakeRequestNoMeta{
		pages:     [][]int{{}},
		execCalls: &calls,
	}
	records, _, err := FetchAllPages(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 0 {
		t.Fatalf("expected 0 records, got %d", len(records))
	}
	if calls != 1 {
		t.Fatalf("expected Execute called exactly once, got %d", calls)
	}
}

// captureStderr redirects os.Stderr for the duration of f and returns what was written.
func captureStderr(f func()) string {
	r, w, _ := os.Pipe()
	old := os.Stderr
	os.Stderr = w
	f()
	w.Close()
	os.Stderr = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestPaginationSummaryWithTotalItems(t *testing.T) {
	total := int32(200)
	meta := sdk.PaginatedResponseMeta{TotalItems: &total}
	out := captureStderr(func() { PrintPaginationSummary(50, meta) })
	want := fmt.Sprintf("Returned %d out of %d records", 50, 200)
	if !strings.Contains(out, want) {
		t.Errorf("expected %q in stderr output, got %q", want, out)
	}
}

// fakeRawServer simulates a paginated API endpoint returning records 1..total
// as {"id":N} objects, honoring page and limit query params (limit capped at 100).
func fakeRawServer(total int) func(page, limit float32) (*http.Response, error) {
	return func(page, limit float32) (*http.Response, error) {
		l := int(limit)
		if l <= 0 || l > 100 {
			l = 100
		}
		p := int(page)
		if p < 1 {
			p = 1
		}
		start := (p - 1) * l
		items := []string{}
		for i := start; i < start+l && i < total; i++ {
			items = append(items, fmt.Sprintf(`{"id":%d}`, i+1))
		}
		totalPages := (total + l - 1) / l
		body := fmt.Sprintf(`{"data":[%s],"meta":{"totalItems":%d,"totalPages":%d,"currentPage":%d}}`,
			strings.Join(items, ","), total, totalPages, p)
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	}
}

func rawIDs(t *testing.T, items []json.RawMessage) []int {
	t.Helper()
	ids := make([]int, 0, len(items))
	for _, it := range items {
		var v struct {
			Id int `json:"id"`
		}
		if err := json.Unmarshal(it, &v); err != nil {
			t.Fatal(err)
		}
		ids = append(ids, v.Id)
	}
	return ids
}

func TestFetchUpToRawUnderPageSize(t *testing.T) {
	items, _, err := FetchUpToRaw(fakeRawServer(500), 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 5 {
		t.Fatalf("expected 5 records, got %d", len(items))
	}
}

func TestFetchUpToRawOverPageSize(t *testing.T) {
	// limit 150 must span 2 API pages (100 + 50)
	items, _, err := FetchUpToRaw(fakeRawServer(500), 150)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 150 {
		t.Fatalf("expected 150 records, got %d", len(items))
	}
	ids := rawIDs(t, items)
	if ids[0] != 1 || ids[149] != 150 {
		t.Fatalf("expected ids 1..150, got first=%d last=%d", ids[0], ids[149])
	}
}

func TestFetchUpToRawFewerAvailable(t *testing.T) {
	items, _, err := FetchUpToRaw(fakeRawServer(30), 150)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 30 {
		t.Fatalf("expected 30 records, got %d", len(items))
	}
}

func TestFetchPageWindowRawSmallLimit(t *testing.T) {
	// page 2 of size 10 = records 11..20
	items, _, err := FetchPageWindowRaw(fakeRawServer(500), 2, 10)
	if err != nil {
		t.Fatal(err)
	}
	ids := rawIDs(t, items)
	if len(ids) != 10 || ids[0] != 11 || ids[9] != 20 {
		t.Fatalf("expected ids 11..20, got %v", ids)
	}
}

func TestFetchPageWindowRawLimitOver100(t *testing.T) {
	// page 3 of size 120 = records 241..360, spans API pages of 100
	items, _, err := FetchPageWindowRaw(fakeRawServer(500), 3, 120)
	if err != nil {
		t.Fatal(err)
	}
	ids := rawIDs(t, items)
	if len(ids) != 120 {
		t.Fatalf("expected 120 records, got %d", len(ids))
	}
	if ids[0] != 241 || ids[119] != 360 {
		t.Fatalf("expected ids 241..360, got first=%d last=%d", ids[0], ids[119])
	}
}

func TestFetchPageWindowRawPastEnd(t *testing.T) {
	// page 10 of size 100 when only 500 records exist → empty
	items, _, err := FetchPageWindowRaw(fakeRawServer(500), 10, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatalf("expected 0 records past end, got %d", len(items))
	}
}

func TestFetchPageWindowRawPartialLastPage(t *testing.T) {
	// page 2 of size 120 when only 200 exist = records 121..200 (80 records)
	items, _, err := FetchPageWindowRaw(fakeRawServer(200), 2, 120)
	if err != nil {
		t.Fatal(err)
	}
	ids := rawIDs(t, items)
	if len(ids) != 80 {
		t.Fatalf("expected 80 records, got %d", len(ids))
	}
	if ids[0] != 121 || ids[79] != 200 {
		t.Fatalf("expected ids 121..200, got first=%d last=%d", ids[0], ids[79])
	}
}

func TestFetchAllPagesRawMultiPage(t *testing.T) {
	items, meta, err := FetchAllPagesRaw(func(page float32) (*http.Response, error) {
		return fakeRawServer(250)(page, 100)
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 250 {
		t.Fatalf("expected 250 records, got %d", len(items))
	}
	if meta.TotalItems == nil || *meta.TotalItems != 250 {
		t.Fatalf("expected TotalItems=250, got %v", meta.TotalItems)
	}
}

func TestPaginationSummaryWithoutTotalItems(t *testing.T) {
	meta := sdk.PaginatedResponseMeta{} // TotalItems nil
	out := captureStderr(func() { PrintPaginationSummary(42, meta) })
	want := fmt.Sprintf("Returned %d out of %d records", 42, 42)
	if !strings.Contains(out, want) {
		t.Errorf("expected %q in stderr output, got %q", want, out)
	}
}
