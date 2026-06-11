package utils

import (
	"fmt"
	"net/http"
	"os"

	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// PaginatedList is implemented by all SDK paginated list models.
type PaginatedList[T any] interface {
	GetData() []T
	GetMeta() sdk.PaginatedResponseMeta
}

// PaginatedRequest is implemented by all SDK fluent list requests that support pagination.
type PaginatedRequest[Req any, L any] interface {
	Page(float32) Req
	Limit(float32) Req
	Execute() (L, *http.Response, error)
}

const defaultPageSize = 100

// FetchAllPages loops a fluent SDK list request through all pages and returns the combined records
// and the last page's meta (for summary printing). Callers should print the summary AFTER rendering
// output using PrintPaginationSummary.
func FetchAllPages[Req PaginatedRequest[Req, L], L PaginatedList[T], T any](req Req) ([]T, sdk.PaginatedResponseMeta, error) {
	records := make([]T, 0)
	req = req.Limit(defaultPageSize)

	var meta sdk.PaginatedResponseMeta
	for page := float32(1); ; page++ {
		result, httpRes, err := req.Page(page).Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			return nil, meta, err
		}

		batch := result.GetData()
		records = append(records, batch...)
		meta = result.GetMeta()

		// Prefer TotalPages when available; fall back to "fewer than limit = last page".
		if meta.TotalPages != nil {
			if *meta.TotalPages <= int32(page) {
				break
			}
		} else if len(batch) < defaultPageSize {
			break
		}
	}

	return records, meta, nil
}

// PrintAll renders results via formatter then prints the pagination summary to stderr.
// Use this instead of formatter.PrintResult directly after FetchAllPages.
func PrintAll(result interface{}, meta sdk.PaginatedResponseMeta, count int, printConfig *formatter.PrintConfig) error {
	if err := formatter.PrintResult(result, printConfig); err != nil {
		return err
	}
	PrintPaginationSummary(count, meta)
	return nil
}

// PrintPaginationSummary prints a "Returned X out of Y records" line to stderr.
// Pass total=-1 when the server-side total is unknown.
func PrintPaginationSummary(returned int, meta sdk.PaginatedResponseMeta) {
	if meta.TotalItems != nil {
		fmt.Fprintf(os.Stderr, "Returned %d out of %d records\n", returned, *meta.TotalItems)
	} else {
		// TotalItems not returned by API; use returned count as total (fetch-all case).
		fmt.Fprintf(os.Stderr, "Returned %d out of %d records\n", returned, returned)
	}
}
