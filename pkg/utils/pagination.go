package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/viper"
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

// FetchUpTo fetches up to n records across pages using the typed SDK request.
// Use for --limit-only paths (no explicit --page). The page size is held constant
// (capped at 100) so page numbers stay well-defined across requests.
func FetchUpTo[Req PaginatedRequest[Req, L], L PaginatedList[T], T any](req Req, n int) ([]T, sdk.PaginatedResponseMeta, error) {
	records := make([]T, 0)
	var meta sdk.PaginatedResponseMeta

	pageSize := n
	if pageSize > defaultPageSize {
		pageSize = defaultPageSize
	}
	req = req.Limit(float32(pageSize))

	for page := 1; len(records) < n; page++ {
		result, httpRes, err := req.Page(float32(page)).Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			return nil, meta, err
		}

		batch := result.GetData()
		records = append(records, batch...)
		meta = result.GetMeta()

		if len(batch) < pageSize {
			break
		}
	}

	if len(records) > n {
		records = records[:n]
	}
	return records, meta, nil
}

// FetchPageWindow fetches the records belonging to user-requested page `page` of size `limit`,
// even when limit exceeds the API's per-page cap (100). It computes the record window
// [(page-1)*limit, page*limit) and collects it across as many API pages as required.
func FetchPageWindow[Req PaginatedRequest[Req, L], L PaginatedList[T], T any](req Req, page, limit int) ([]T, sdk.PaginatedResponseMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = defaultPageSize
	}

	start := (page - 1) * limit // 0-based index of first wanted record
	end := page * limit         // exclusive

	apiPageSize := limit
	if apiPageSize > defaultPageSize {
		apiPageSize = defaultPageSize
	}
	req = req.Limit(float32(apiPageSize))

	window := make([]T, 0, limit)
	var meta sdk.PaginatedResponseMeta

	for apiPage := start/apiPageSize + 1; len(window) < limit; apiPage++ {
		result, httpRes, err := req.Page(float32(apiPage)).Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			return nil, meta, err
		}

		batch := result.GetData()
		meta = result.GetMeta()

		// Index of first record on this API page within the overall record stream.
		pageStart := (apiPage - 1) * apiPageSize
		for i, item := range batch {
			abs := pageStart + i
			if abs >= start && abs < end {
				window = append(window, item)
			}
		}

		if len(batch) < apiPageSize {
			break // last page reached
		}
		if pageStart+len(batch) >= end {
			break // window complete
		}
	}

	return window, meta, nil
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

// PrintAllRaw renders full-fidelity raw API objects for json/yaml output, and the
// typed records (matching printConfig fields) for table formats. This keeps machine
// formats lossless while table formats use the safe raw structs.
func PrintAllRaw(rawItems []json.RawMessage, records interface{}, meta sdk.PaginatedResponseMeta, count int, printConfig *formatter.PrintConfig) error {
	format := strings.ToLower(viper.GetString(formatter.ConfigFormat))
	if format == "json" || format == "yaml" {
		combined := make([]interface{}, 0, len(rawItems))
		for _, item := range rawItems {
			var v interface{}
			if err := json.Unmarshal(item, &v); err != nil {
				return fmt.Errorf("failed to decode raw item: %w", err)
			}
			combined = append(combined, v)
		}
		if err := formatter.PrintResult(combined, printConfig); err != nil {
			return err
		}
		PrintPaginationSummary(count, meta)
		return nil
	}
	return PrintAll(records, meta, count, printConfig)
}

// rawPaginatedEnvelope is the generic JSON envelope for all paginated SDK endpoints.
type rawPaginatedEnvelope struct {
	Data json.RawMessage `json:"data"`
	Meta struct {
		TotalPages  *int32 `json:"totalPages"`
		CurrentPage *int32 `json:"currentPage"`
		TotalItems  *int32 `json:"totalItems"`
	} `json:"meta"`
}

// FetchAllPagesRaw fetches all pages without relying on SDK type unmarshalling.
// Use when the SDK struct's UnmarshalJSON rejects valid API responses due to schema drift.
// fetch receives the 1-based page number and must return the raw *http.Response from Execute().
func FetchAllPagesRaw(fetch func(page float32) (*http.Response, error)) ([]json.RawMessage, sdk.PaginatedResponseMeta, error) {
	var all []json.RawMessage
	var lastMeta sdk.PaginatedResponseMeta

	for page := float32(1); ; page++ {
		httpRes, fetchErr := fetch(page)
		if httpRes != nil && httpRes.StatusCode >= 400 {
			return nil, lastMeta, response_inspector.InspectResponse(httpRes, fetchErr)
		}
		if httpRes == nil {
			return nil, lastMeta, fetchErr
		}

		body, err := io.ReadAll(httpRes.Body)
		if err != nil {
			return nil, lastMeta, fmt.Errorf("failed to read response body: %w", err)
		}

		var envelope rawPaginatedEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			return nil, lastMeta, fmt.Errorf("failed to parse paginated response: %w", err)
		}

		var items []json.RawMessage
		if err := json.Unmarshal(envelope.Data, &items); err != nil {
			return nil, lastMeta, fmt.Errorf("failed to parse response data: %w", err)
		}
		all = append(all, items...)

		if envelope.Meta.TotalItems != nil {
			lastMeta.TotalItems = envelope.Meta.TotalItems
		}
		if envelope.Meta.TotalPages != nil {
			lastMeta.TotalPages = envelope.Meta.TotalPages
			if *envelope.Meta.TotalPages <= int32(page) {
				break
			}
		} else if len(items) < defaultPageSize {
			break
		}
	}

	return all, lastMeta, nil
}

// FetchUpToRaw fetches up to n records across pages without SDK type unmarshalling.
// Use for --limit-only paths (no explicit --page). fetch receives the page number and
// desired page size (capped at 100 by the caller).
func FetchUpToRaw(fetch func(page, limit float32) (*http.Response, error), n int) ([]json.RawMessage, sdk.PaginatedResponseMeta, error) {
	var all []json.RawMessage
	var lastMeta sdk.PaginatedResponseMeta

	// Page size must stay constant across requests — page N is defined relative to it.
	pageSize := n
	if pageSize > defaultPageSize {
		pageSize = defaultPageSize
	}

	for page := 1; len(all) < n; page++ {
		httpRes, fetchErr := fetch(float32(page), float32(pageSize))
		if httpRes != nil && httpRes.StatusCode >= 400 {
			return nil, lastMeta, response_inspector.InspectResponse(httpRes, fetchErr)
		}
		if httpRes == nil {
			return nil, lastMeta, fetchErr
		}

		body, err := io.ReadAll(httpRes.Body)
		if err != nil {
			return nil, lastMeta, fmt.Errorf("failed to read response body: %w", err)
		}

		var envelope rawPaginatedEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			return nil, lastMeta, fmt.Errorf("failed to parse paginated response: %w", err)
		}

		var items []json.RawMessage
		if err := json.Unmarshal(envelope.Data, &items); err != nil {
			return nil, lastMeta, fmt.Errorf("failed to parse response data: %w", err)
		}
		all = append(all, items...)

		if envelope.Meta.TotalItems != nil {
			lastMeta.TotalItems = envelope.Meta.TotalItems
		}
		if envelope.Meta.TotalPages != nil {
			lastMeta.TotalPages = envelope.Meta.TotalPages
		}
		if len(items) < pageSize {
			break
		}
	}

	if len(all) > n {
		all = all[:n]
	}
	return all, lastMeta, nil
}

// FetchPageWindowRaw fetches the records belonging to user-requested page `page` of size `limit`,
// even when limit exceeds the API's per-page cap (100). It computes the record window
// [(page-1)*limit, page*limit) and collects it across as many API pages as required.
// fetch receives the API page number and API page size.
func FetchPageWindowRaw(fetch func(page, limit float32) (*http.Response, error), page, limit int) ([]json.RawMessage, sdk.PaginatedResponseMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = defaultPageSize
	}

	start := (page - 1) * limit // 0-based index of first wanted record
	end := page * limit         // exclusive

	apiPageSize := limit
	if apiPageSize > defaultPageSize {
		apiPageSize = defaultPageSize
	}

	var window []json.RawMessage
	var lastMeta sdk.PaginatedResponseMeta

	for apiPage := start/apiPageSize + 1; len(window) < limit; apiPage++ {
		httpRes, fetchErr := fetch(float32(apiPage), float32(apiPageSize))
		if httpRes != nil && httpRes.StatusCode >= 400 {
			return nil, lastMeta, response_inspector.InspectResponse(httpRes, fetchErr)
		}
		if httpRes == nil {
			return nil, lastMeta, fetchErr
		}

		body, err := io.ReadAll(httpRes.Body)
		if err != nil {
			return nil, lastMeta, fmt.Errorf("failed to read response body: %w", err)
		}

		var envelope rawPaginatedEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			return nil, lastMeta, fmt.Errorf("failed to parse paginated response: %w", err)
		}

		var items []json.RawMessage
		if err := json.Unmarshal(envelope.Data, &items); err != nil {
			return nil, lastMeta, fmt.Errorf("failed to parse response data: %w", err)
		}

		if envelope.Meta.TotalItems != nil {
			lastMeta.TotalItems = envelope.Meta.TotalItems
		}
		if envelope.Meta.TotalPages != nil {
			lastMeta.TotalPages = envelope.Meta.TotalPages
		}

		// Index of first record on this API page within the overall record stream.
		pageStart := (apiPage - 1) * apiPageSize
		for i, item := range items {
			abs := pageStart + i
			if abs >= start && abs < end {
				window = append(window, item)
			}
		}

		if len(items) < apiPageSize {
			break // last page reached
		}
		if pageStart+len(items) >= end {
			break // window complete
		}
	}

	return window, lastMeta, nil
}

// ParseRawPage reads a single paginated API response from httpRes without SDK type unmarshalling.
// Use for single-page list paths that would otherwise fail due to SDK schema drift.
func ParseRawPage(httpRes *http.Response) ([]json.RawMessage, sdk.PaginatedResponseMeta, error) {
	var meta sdk.PaginatedResponseMeta
	if httpRes == nil {
		return nil, meta, fmt.Errorf("no response from server")
	}
	if httpRes.StatusCode >= 400 {
		return nil, meta, response_inspector.InspectResponse(httpRes, nil)
	}
	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, meta, fmt.Errorf("failed to read response body: %w", err)
	}
	var envelope rawPaginatedEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, meta, fmt.Errorf("failed to parse paginated response: %w", err)
	}
	var items []json.RawMessage
	if err := json.Unmarshal(envelope.Data, &items); err != nil {
		return nil, meta, fmt.Errorf("failed to parse response data: %w", err)
	}
	if envelope.Meta.TotalItems != nil {
		meta.TotalItems = envelope.Meta.TotalItems
	}
	return items, meta, nil
}

// UnmarshalRawItems unmarshals a []json.RawMessage slice (from FetchAllPagesRaw) into a typed slice.
func UnmarshalRawItems[T any](rawItems []json.RawMessage) ([]T, error) {
	result := make([]T, 0, len(rawItems))
	for _, item := range rawItems {
		var r T
		if err := json.Unmarshal(item, &r); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
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
