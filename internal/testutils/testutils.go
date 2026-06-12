// Package testutils provides shared test helpers for metalcloud-cli unit tests.
package testutils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/viper"
)

// SetupTestFormat configures viper so formatter.PrintResult uses JSON output.
// Call this from TestMain in each test package.
func SetupTestFormat() {
	viper.Set(formatter.ConfigFormat, "json")
}

// SetupTestContext creates a context with an SDK API client pointed at serverURL.
func SetupTestContext(serverURL string) context.Context {
	cfg := sdk.NewConfiguration()
	cfg.Servers = []sdk.ServerConfiguration{{URL: serverURL}}
	client := sdk.NewAPIClient(cfg)
	ctx := context.WithValue(context.Background(), api.ApiClientContextKey, client)
	ctx = context.WithValue(ctx, sdk.ContextAccessToken, "test-api-key")
	return ctx
}

// NewTestServer creates an httptest.Server that routes requests by URL path.
func NewTestServer(routes map[string]http.HandlerFunc) *httptest.Server {
	mux := http.NewServeMux()
	for path, handler := range routes {
		mux.HandleFunc(path, handler)
	}
	return httptest.NewServer(mux)
}

// JSONHandler returns an http.HandlerFunc that writes statusCode and body as JSON.
func JSONHandler(statusCode int, body any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if err := json.NewEncoder(w).Encode(body); err != nil {
			panic(fmt.Sprintf("testutils.JSONHandler: encode failed: %v", err))
		}
	}
}

// RawHandler returns an http.HandlerFunc that writes a raw JSON string.
func RawHandler(statusCode int, rawJSON string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(rawJSON))
	}
}

// PaginatedMeta builds the standard pagination metadata envelope.
func PaginatedMeta(currentPage, totalPages, itemsPerPage int32) map[string]any {
	return map[string]any{
		"currentPage":  currentPage,
		"totalPages":   totalPages,
		"itemsPerPage": itemsPerPage,
	}
}

// PaginatedResponse wraps a data slice in the standard paginated envelope.
func PaginatedResponse(data any, currentPage, totalPages int32) map[string]any {
	return map[string]any{
		"data": data,
		"meta": PaginatedMeta(currentPage, totalPages, 100),
	}
}

// MultiPageServer creates an httptest.Server that serves sequential pages of items.
// Each call to the path returns the next page until exhausted, then repeats the last.
func MultiPageServer(path string, pages []any) *httptest.Server {
	var counter atomic.Int32
	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		idx := int(counter.Add(1)) - 1
		if idx >= len(pages) {
			idx = len(pages) - 1
		}
		total := int32(len(pages))
		current := int32(idx + 1)
		resp := PaginatedResponse(pages[idx], current, total)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})
	return httptest.NewServer(mux)
}

// ErrorHandler returns an http.HandlerFunc that returns an API error response.
func ErrorHandler(statusCode int, message string) http.HandlerFunc {
	return RawHandler(statusCode, fmt.Sprintf(`{"message":%q,"statusCode":%d}`, message, statusCode))
}
