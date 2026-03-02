package user

import (
	"context"
	"fmt"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var apiKeyPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ApiKey": {
			Title: "API Key",
			Order: 1,
		},
	},
}

var twoFASecretPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Secret2FA": {
			Title: "Secret",
			Order: 1,
		},
		"QrCode": {
			Title: "QR Code",
			Order: 2,
		},
	},
}

func GetApiKey(ctx context.Context) error {
	logger.Get().Info().Msgf("Getting API key for current user")

	client := api.GetApiClient(ctx)

	apiKey, httpRes, err := client.UserAPI.GetUserApiKey(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Println(apiKey.ApiKey)
	return nil
}

func RegenerateApiKey(ctx context.Context) error {
	logger.Get().Info().Msgf("Regenerating API key for current user")

	revision, err := getCurrentUserRevision(ctx)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UserAPI.RegenerateUserApiKey(ctx).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Println("API key has been regenerated. Your previous key is no longer valid.")
	return formatter.PrintResult(userInfo, &userPrintConfig)
}

func Enable2FA(ctx context.Context, token string) error {
	logger.Get().Info().Msgf("Enabling 2FA for current user")

	revision, err := getCurrentUserRevision(ctx)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.UserAPI.EnableUser2FA(ctx).
		TwoFactorAuthenticationToken(sdk.TwoFactorAuthenticationToken{
			Token: token,
		}).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Println("Two-factor authentication enabled successfully.")
	return nil
}

func Disable2FA(ctx context.Context) error {
	logger.Get().Info().Msgf("Disabling 2FA for current user")

	revision, err := getCurrentUserRevision(ctx)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.UserAPI.DisableUser2FA(ctx).IfMatch(revision).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	fmt.Println("Two-factor authentication disabled.")
	return nil
}

func GenerateUser2FASecret(ctx context.Context) error {
	logger.Get().Info().Msgf("Generating 2FA secret for current user")

	client := api.GetApiClient(ctx)

	secret, httpRes, err := client.UserAPI.GenerateUser2FASecret(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(secret, &twoFASecretPrintConfig)
}

func getCurrentUserRevision(ctx context.Context) (string, error) {
	userIdStr := api.GetUserId(ctx)
	userIdNumeric, err := getUserId(userIdStr)
	if err != nil {
		return "", err
	}

	client := api.GetApiClient(ctx)

	userInfo, httpRes, err := client.UsersAPI.GetUser(ctx, userIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", int(userInfo.Revision)), nil
}
