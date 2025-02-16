package response_inspector

import (
	"fmt"
	"net/http"

	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
)

func InspectResponse(httpRes *http.Response, err error) error {
	if err != nil {
		if httpRes != nil {
			err := fmt.Errorf("%s - %s", httpRes.Status, httpRes.Body)
			logger.Get().Error().Err(err).Msg("")
			return err
		}
		logger.Get().Error().Err(err).Msg("")
		return err
	}
	if httpRes.StatusCode != 200 {
		err := fmt.Errorf("%s - %s", httpRes.Status, httpRes.Body)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	return nil
}
