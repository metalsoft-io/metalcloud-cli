package response_inspector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
)

func InspectResponse(httpRes *http.Response, err error) error {
	if httpRes != nil && httpRes.StatusCode >= 400 {
		err := fmt.Errorf("%s - %s", httpRes.Status, errorBody(httpRes))
		logger.Get().Error().Err(err).Msg("")
		return err
	}
	if err != nil {
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	return nil
}

// errorBody returns the response body text for an error message. SDK calls
// leave a replayable buffer in Body, but raw HTTP calls leave the live network
// body, which prints as a struct dump; read it and put a replayable copy back.
func errorBody(httpRes *http.Response) string {
	if httpRes.Body == nil {
		return ""
	}
	if stringer, ok := httpRes.Body.(fmt.Stringer); ok {
		return stringer.String()
	}
	bodyBytes, readErr := io.ReadAll(httpRes.Body)
	_ = httpRes.Body.Close()
	httpRes.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	if readErr != nil {
		return fmt.Sprintf("<unreadable body: %v>", readErr)
	}
	return string(bodyBytes)
}

func ParseResponseBody(httpRes *http.Response) (map[string]interface{}, error) {
	if httpRes == nil {
		return nil, fmt.Errorf("http response is nil")
	}

	bodyBytes, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	defer httpRes.Body.Close()

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return result, nil
}
