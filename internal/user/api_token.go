package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var apiTokenPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title:       "ID",
			Transformer: formatter.FormatIdValue,
			Order:       1,
		},
		"Name": {
			MaxWidth: 30,
			Order:    2,
		},
		"Label": {
			MaxWidth: 30,
			Order:    3,
		},
		"Scope": {
			Order: 4,
		},
		"Status": {
			Transformer: formatter.FormatStatusValue,
			Order:       5,
		},
		"IsLegacy": {
			Title: "Legacy",
			Order: 6,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
		"ExpiresTimestamp": {
			Title:       "Expires",
			Transformer: formatter.FormatDateTimeValue,
			Order:       8,
		},
	},
}

var apiTokenCreatedPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title:       "ID",
			Transformer: formatter.FormatIdValue,
			Order:       1,
		},
		"Label": {
			MaxWidth: 30,
			Order:    2,
		},
		"ExpiresTimestamp": {
			Title:       "Expires",
			Transformer: formatter.FormatDateTimeValue,
			Order:       3,
		},
		"Token": {
			Order: 4,
		},
	},
}

var apiTokenSettingsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"DefaultValidityDuration": {
			Title:       "Default Validity (s)",
			Transformer: formatter.FormatIntegerValue,
			Order:       1,
		},
		"MaxPerUser": {
			Title:       "Max Per User",
			Transformer: formatter.FormatIntegerValue,
			Order:       2,
		},
		"MaxExpiresTimestamp": {
			Title:       "Max Expires",
			Transformer: formatter.FormatDateTimeValue,
			Order:       3,
		},
	},
}

// ApiTokenCreateOptions carries the flag-based input of 'user api-token create'.
type ApiTokenCreateOptions struct {
	Name        string
	Label       string
	Description string
	Scope       string
	Expires     string
	Permissions []string
}

func ApiTokenList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing API tokens of the current user")

	client := api.GetApiClient(ctx)

	tokens, httpRes, err := client.UserAPI.ListUserApiTokens(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(tokens, &apiTokenPrintConfig)
}

func ApiTokenGet(ctx context.Context, apiTokenId string) error {
	logger.Get().Info().Msgf("Getting API token '%s'", apiTokenId)

	apiTokenIdNumeric, err := getApiTokenId(apiTokenId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	token, httpRes, err := client.UserAPI.GetUserApiToken(ctx, apiTokenIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(token, &apiTokenPrintConfig)
}

func ApiTokenSettings(ctx context.Context) error {
	logger.Get().Info().Msgf("Getting API token settings")

	client := api.GetApiClient(ctx)

	settings, httpRes, err := client.UserAPI.GetUserApiTokenSettings(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(settings, &apiTokenSettingsPrintConfig)
}

// ApiTokenCreate creates an API token either from a JSON/YAML config (when
// config is non-nil) or from the flag-based options.
func ApiTokenCreate(ctx context.Context, options ApiTokenCreateOptions, config []byte) error {
	logger.Get().Info().Msgf("Creating API token")

	var createApiToken sdk.CreateApiToken

	if config != nil {
		if err := utils.UnmarshalContent(config, &createApiToken); err != nil {
			return err
		}
	} else {
		expires, err := resolveApiTokenExpiry(ctx, options.Expires)
		if err != nil {
			return err
		}

		createApiToken = sdk.CreateApiToken{
			Name:             options.Name,
			ExpiresTimestamp: expires,
			Permissions: sdk.ApiTokenPermissionDocument{
				Permissions: options.Permissions,
			},
		}
		if options.Label != "" {
			createApiToken.Label = sdk.PtrString(options.Label)
		}
		if options.Description != "" {
			createApiToken.Description = sdk.PtrString(options.Description)
		}
		if options.Scope != "" {
			scope, err := sdk.NewApiTokenScopeFromValue(options.Scope)
			if err != nil {
				return err
			}
			createApiToken.Scope = scope
		}
	}

	// The API rejects null arrays in the permission document.
	if createApiToken.Permissions.Permissions == nil {
		createApiToken.Permissions.Permissions = []string{}
	}
	if createApiToken.Permissions.Policies == nil {
		createApiToken.Permissions.Policies = []sdk.PolicyEntry{}
	}

	client := api.GetApiClient(ctx)

	token, httpRes, err := client.UserAPI.CreateUserApiToken(ctx).CreateApiToken(createApiToken).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return printApiTokenCreated(token)
}

func ApiTokenRegenerate(ctx context.Context, apiTokenId string, expires string) error {
	logger.Get().Info().Msgf("Regenerating API token '%s'", apiTokenId)

	apiTokenIdNumeric, err := getApiTokenId(apiTokenId)
	if err != nil {
		return err
	}

	expiresTimestamp, err := resolveApiTokenExpiry(ctx, expires)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	token, httpRes, err := client.UserAPI.
		RegenerateUserApiToken(ctx, apiTokenIdNumeric).
		RegenerateApiToken(sdk.RegenerateApiToken{ExpiresTimestamp: expiresTimestamp}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return printApiTokenCreated(token)
}

func ApiTokenDelete(ctx context.Context, apiTokenId string) error {
	logger.Get().Info().Msgf("Deleting API token '%s'", apiTokenId)

	apiTokenIdNumeric, err := getApiTokenId(apiTokenId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.UserAPI.DeleteUserApiToken(ctx, apiTokenIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("API token '%s' deleted", apiTokenId)
	return nil
}

func printApiTokenCreated(token *sdk.ApiTokenCreated) error {
	if formatter.IsTextFormat() {
		fmt.Println("Store the token value now: it will not be shown again.")
	}

	return formatter.PrintResult(token, &apiTokenCreatedPrintConfig)
}

// resolveApiTokenExpiry parses an RFC 3339 timestamp or a YYYY-MM-DD date
// (midnight UTC). When empty, the controller's default token validity is used.
func resolveApiTokenExpiry(ctx context.Context, expires string) (time.Time, error) {
	expires = strings.TrimSpace(expires)

	if expires == "" {
		client := api.GetApiClient(ctx)

		settings, httpRes, err := client.UserAPI.GetUserApiTokenSettings(ctx).Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			return time.Time{}, err
		}

		return time.Now().UTC().Add(time.Duration(settings.DefaultValidityDuration) * time.Second).Truncate(time.Second), nil
	}

	if t, err := time.Parse(time.RFC3339, expires); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.DateOnly, expires); err == nil {
		return t, nil
	}

	err := fmt.Errorf("invalid expiry '%s': use an RFC 3339 timestamp (2027-01-01T00:00:00Z) or a date (2027-01-01)", expires)
	logger.Get().Error().Err(err).Msg("")
	return time.Time{}, err
}

func getApiTokenId(apiTokenId string) (int64, error) {
	apiTokenIdNumeric, err := utils.GetInt64FromString(apiTokenId)
	if err != nil {
		err := fmt.Errorf("invalid API token ID: '%s'", apiTokenId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return apiTokenIdNumeric, nil
}
