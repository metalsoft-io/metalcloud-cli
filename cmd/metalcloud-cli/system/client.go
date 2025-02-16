package system

import (
	"context"
	"fmt"

	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

func GetApiClient(ctx context.Context) (*sdk.APIClient, error) {
	client := ctx.Value(ApiClientContextKey)
	if client == nil {
		err := fmt.Errorf("SDK client not found in context")
		logger.Get().Error().Err(err).Msg("")
		return nil, err
	}

	apiClient, ok := client.(*sdk.APIClient)
	if !ok {
		err := fmt.Errorf("invalid SDK client in context")
		logger.Get().Error().Err(err).Msg("")
		return nil, err
	}

	return apiClient, nil
}
