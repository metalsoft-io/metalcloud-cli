package api

import (
	"context"
	"fmt"

	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
)

type ContextKey string

const (
	ApiClientContextKey ContextKey = "apiClient"
	UserIdContextKey    ContextKey = "userId"
)

func GetApiClient(ctx context.Context) *sdk.APIClient {
	client := ctx.Value(ApiClientContextKey)
	if client == nil {
		err := fmt.Errorf("SDK client not found in context")
		logger.Get().Error().Err(err).Msg("")
		cobra.CheckErr(err)
	}

	apiClient, ok := client.(*sdk.APIClient)
	if !ok {
		err := fmt.Errorf("invalid SDK client in context")
		logger.Get().Error().Err(err).Msg("")
		cobra.CheckErr(err)
	}

	return apiClient
}

func GetApiClientE(ctx context.Context) (*sdk.APIClient, error) {
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

func SetApiClient(ctx context.Context, apiEndpoint string, apiKey string, debug bool) context.Context {
	// Initialize API client using the arguments from the command line or environment variables
	cfg := sdk.NewConfiguration()
	cfg.UserAgent = "metalcloud-cli"
	cfg.Servers = []sdk.ServerConfiguration{
		{
			URL:         apiEndpoint,
			Description: "MetalSoft",
		},
	}

	// Set debug mode
	cfg.Debug = debug

	// Create API client
	apiClient := sdk.NewAPIClient(cfg)

	ctx = context.WithValue(ctx, ApiClientContextKey, apiClient)
	ctx = context.WithValue(ctx, sdk.ContextAccessToken, apiKey)

	return ctx
}

func GetUserId(ctx context.Context) string {
	userId := ctx.Value(UserIdContextKey)
	if userId == nil {
		err := fmt.Errorf("user ID not found in context")
		logger.Get().Error().Err(err).Msg("")
		cobra.CheckErr(err)
	}

	userIdStr, ok := userId.(string)
	if !ok {
		err := fmt.Errorf("invalid user ID in context")
		logger.Get().Error().Err(err).Msg("")
		cobra.CheckErr(err)
	}

	return userIdStr
}

func SetUserId(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, UserIdContextKey, userId)
}
