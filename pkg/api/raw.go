package api

import (
	"context"
	"fmt"
	"io"

	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
)

// RawJSONRequest issues a request through DoJSONRequestWithHeaders, runs the
// standard response inspection, and returns the response body. It is the
// one-call form of the "raw request, inspect, read, close" sequence that
// commands need when the typed SDK models reject a valid API payload.
// The body is nil for responses without content (e.g. 204).
func RawJSONRequest(ctx context.Context, method, path string, body []byte, headers map[string]string) ([]byte, error) {
	httpRes, err := DoJSONRequestWithHeaders(ctx, method, path, body, headers)
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	responseBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return responseBody, nil
}

// IfMatchHeader returns the headers map carrying an If-Match entity tag, or nil
// when revision is empty so that no header is sent.
func IfMatchHeader(revision string) map[string]string {
	if revision == "" {
		return nil
	}
	return map[string]string{"If-Match": revision}
}
