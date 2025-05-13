package secret

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var secretPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Name": {
			MaxWidth: 30,
			Order:    2,
		},
		"ValueEncrypted": {
			MaxWidth: 50,
			Title:    "Value (encrypted)",
			Order:    3,
			Transformer: func(value interface{}) string {
				if value == nil {
					return ""
				}
				return fmt.Sprint(value)
			},
		},
		"Usage": {
			Title: "Usage Type",
			Order: 4,
		},
		"UserIdOwner": {
			Title: "Owner ID",
			Order: 5,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
	},
}

func SecretList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all secrets")

	client := api.GetApiClient(ctx)

	secretList, httpRes, err := client.SecretsAPI.GetSecrets(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(secretList, &secretPrintConfig)
}

func SecretGet(ctx context.Context, secretId string) error {
	logger.Get().Info().Msgf("Get secret '%s' details", secretId)

	secretIdNumeric, err := getSecretId(secretId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	secret, httpRes, err := client.SecretsAPI.GetSecret(ctx, secretIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(secret, &secretPrintConfig)
}

func SecretCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating secret")

	var secretConfig sdk.CreateSecret
	err := json.Unmarshal(config, &secretConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	secret, httpRes, err := client.SecretsAPI.
		CreateSecret(ctx).
		CreateSecret(secretConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(secret, &secretPrintConfig)
}

func SecretUpdate(ctx context.Context, secretId string, config []byte) error {
	logger.Get().Info().Msgf("Updating secret '%s'", secretId)

	secretIdNumeric, err := getSecretId(secretId)
	if err != nil {
		return err
	}

	var secretConfig sdk.UpdateSecret
	err = json.Unmarshal(config, &secretConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	secret, httpRes, err := client.SecretsAPI.
		UpdateSecret(ctx, secretIdNumeric).
		UpdateSecret(secretConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(secret, &secretPrintConfig)
}

func SecretDelete(ctx context.Context, secretId string) error {
	logger.Get().Info().Msgf("Deleting secret '%s'", secretId)

	secretIdNumeric, err := getSecretId(secretId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.SecretsAPI.
		DeleteSecret(ctx, secretIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Secret '%s' deleted", secretId)
	return nil
}

func SecretConfigExample(ctx context.Context) error {
	// Example create secret configuration
	secretConfiguration := sdk.CreateSecret{
		Name:  "example-secret",
		Value: "my-secret-value",
		Usage: (*sdk.VariableUsageType)(sdk.PtrString("credential")),
	}

	return formatter.PrintResult(secretConfiguration, nil)
}

func getSecretId(secretId string) (float32, error) {
	secretIdNumeric, err := strconv.ParseFloat(secretId, 32)
	if err != nil {
		err := fmt.Errorf("invalid secret ID: '%s'", secretId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(secretIdNumeric), nil
}
