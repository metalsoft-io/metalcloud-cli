package utils

import (
	"errors"
	"net/http"
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
	records, err := FetchAllPages(fakeRequest{pages: [][]int{{1, 2, 3}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 3 {
		t.Fatalf("expected 3 records, got %d", len(records))
	}
}

func TestFetchAllPagesMultiPage(t *testing.T) {
	records, err := FetchAllPages(fakeRequest{pages: [][]int{{1, 2}, {3, 4}, {5}}})
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
	_, err := FetchAllPages(fakeRequest{pages: [][]int{{1, 2}, {3, 4}}, failOn: 2})
	if err == nil {
		t.Fatal("expected error")
	}
}
