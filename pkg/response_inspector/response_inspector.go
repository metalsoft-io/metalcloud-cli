package response_inspector

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
)

func InspectResponse(httpRes *http.Response, err error) error {
	if err != nil {
		if httpRes != nil && httpRes.StatusCode >= 400 {
			err := fmt.Errorf("%s - %s", httpRes.Status, httpRes.Body)
			logger.Get().Error().Err(err).Msg("")
			return err
		}
		logger.Get().Error().Err(err).Msg("")
		return err
	}
	if httpRes.StatusCode >= 400 {
		err := fmt.Errorf("%s - %s", httpRes.Status, httpRes.Body)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	return nil
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
