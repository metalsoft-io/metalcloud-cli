package network_device

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var networkDeviceDefaultSecretsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"SiteId": {
			Title: "Site ID",
			Order: 2,
		},
		"MacAddressOrSerialNumber": {
			Title:    "MAC/Serial",
			MaxWidth: 30,
			Order:    3,
		},
		"SecretName": {
			Title: "Secret Name",
			Order: 4,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       5,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
	},
}

func NetworkDeviceDefaultSecretsList(ctx context.Context, page int, limit int) error {
	logger.Get().Info().Msgf("Listing network device default secrets")

	client := api.GetApiClient(ctx)

	req := client.NetworkDeviceDefaultSecretsAPI.GetNetworkDevicesDefaultSecrets(ctx)

	if page > 0 {
		req = req.Page(float32(page))
	}

	if limit > 0 {
		req = req.Limit(float32(limit))
	}

	secretsList, httpRes, err := req.SortBy([]string{"id:ASC"}).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(secretsList.Data, &networkDeviceDefaultSecretsPrintConfig)
}

func NetworkDeviceDefaultSecretsGet(ctx context.Context, secretsId string) error {
	logger.Get().Info().Msgf("Get network device default secrets '%s'", secretsId)

	secretsIdNumeric, err := parseSecretsId(secretsId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	secrets, httpRes, err := client.NetworkDeviceDefaultSecretsAPI.GetNetworkDeviceDefaultSecretsInfo(ctx, secretsIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(secrets, &networkDeviceDefaultSecretsPrintConfig)
}

func NetworkDeviceDefaultSecretsGetCredentials(ctx context.Context, secretsId string) error {
	logger.Get().Info().Msgf("Get network device default secrets credentials for '%s'", secretsId)

	secretsIdNumeric, err := parseSecretsId(secretsId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	credentials, httpRes, err := client.NetworkDeviceDefaultSecretsAPI.GetNetworkDeviceDefaultSecretsCredentials(ctx, secretsIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, nil)
}

func NetworkDeviceDefaultSecretsCreate(ctx context.Context, siteId float32, macAddressOrSerialNumber string, secretName string, secretValue string) error {
	logger.Get().Info().Msgf("Creating network device default secrets")

	client := api.GetApiClient(ctx)

	createSecrets := sdk.NewCreateNetworkDeviceDefaultSecrets(siteId, macAddressOrSerialNumber, secretName, secretValue)

	secrets, httpRes, err := client.NetworkDeviceDefaultSecretsAPI.CreateNetworkDeviceDefaultSecrets(ctx).
		CreateNetworkDeviceDefaultSecrets(*createSecrets).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(secrets, &networkDeviceDefaultSecretsPrintConfig)
}

func NetworkDeviceDefaultSecretsUpdate(ctx context.Context, secretsId string, secretValue string) error {
	logger.Get().Info().Msgf("Updating network device default secrets '%s'", secretsId)

	secretsIdNumeric, err := parseSecretsId(secretsId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	updateSecrets := sdk.NewUpdateNetworkDeviceDefaultSecrets()
	updateSecrets.SetSecretValue(secretValue)

	secrets, httpRes, err := client.NetworkDeviceDefaultSecretsAPI.
		UpdateNetworkDeviceDefaultSecrets(ctx, secretsIdNumeric).
		UpdateNetworkDeviceDefaultSecrets(*updateSecrets).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(secrets, &networkDeviceDefaultSecretsPrintConfig)
}

func NetworkDeviceDefaultSecretsDelete(ctx context.Context, secretsId string) error {
	logger.Get().Info().Msgf("Deleting network device default secrets '%s'", secretsId)

	secretsIdNumeric, err := parseSecretsId(secretsId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceDefaultSecretsAPI.DeleteNetworkDeviceDefaultSecrets(ctx, secretsIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network device default secrets with ID %s deleted successfully", secretsId)
	return nil
}

func parseSecretsId(secretsId string) (float32, error) {
	secretsIdNumeric, err := strconv.ParseFloat(secretsId, 32)
	if err != nil {
		err := fmt.Errorf("invalid network device default secrets ID: '%s'", secretsId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(secretsIdNumeric), nil
}
