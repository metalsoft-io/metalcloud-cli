package network_device_default_secrets

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"

	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

type networkDeviceDefaultSecretsRaw struct {
	Id                       interface{} `json:"id"`
	SiteId                   interface{} `json:"siteId"`
	MacAddressOrSerialNumber *string     `json:"macAddressOrSerialNumber"`
	SecretName               *string     `json:"secretName"`
	CreatedTimestamp         *string     `json:"createdTimestamp"`
	UpdatedTimestamp         *string     `json:"updatedTimestamp"`
}

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
			Title: "MAC/Serial",
			Order: 3,
		},
		"SecretName": {
			Title: "Secret Name",
			Order: 4,
		},
		"CreatedTimestamp": {
			Title: "Created",
			Order: 5,
		},
		"UpdatedTimestamp": {
			Title: "Updated",
			Order: 6,
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

	req = req.SortBy([]string{"id:ASC"})

	if page > 0 {
		rawItems, meta, err := utils.FetchPageWindowRaw(func(p, l float32) (*http.Response, error) {
			_, httpRes, _ := req.Page(p).Limit(l).Execute()
			return httpRes, nil
		}, page, limit)
		if err != nil {
			return err
		}
		records, err := utils.UnmarshalRawItems[networkDeviceDefaultSecretsRaw](rawItems)
		if err != nil {
			return fmt.Errorf("failed to parse network device default secrets: %w", err)
		}
		return utils.PrintAllRaw(rawItems, records, meta, len(records), &networkDeviceDefaultSecretsPrintConfig)
	}

	if limit > 0 {
		rawItems, meta, err := utils.FetchUpToRaw(func(p, l float32) (*http.Response, error) {
			_, httpRes, _ := req.Page(p).Limit(l).Execute()
			return httpRes, nil
		}, limit)
		if err != nil {
			return err
		}
		records, err := utils.UnmarshalRawItems[networkDeviceDefaultSecretsRaw](rawItems)
		if err != nil {
			return fmt.Errorf("failed to parse network device default secrets: %w", err)
		}
		return utils.PrintAllRaw(rawItems, records, meta, len(records), &networkDeviceDefaultSecretsPrintConfig)
	}

	rawItems, meta, err := utils.FetchAllPagesRaw(func(p float32) (*http.Response, error) {
		_, httpRes, _ := req.Page(p).Limit(100).Execute()
		return httpRes, nil
	})
	if err != nil {
		return err
	}
	records, err := utils.UnmarshalRawItems[networkDeviceDefaultSecretsRaw](rawItems)
	if err != nil {
		return fmt.Errorf("failed to parse network device default secrets: %w", err)
	}
	return utils.PrintAllRaw(rawItems, records, meta, len(records), &networkDeviceDefaultSecretsPrintConfig)
}

func NetworkDeviceDefaultSecretsGet(ctx context.Context, id string) error {
	logger.Get().Info().Msgf("Getting network device default secrets '%s'", id)

	idNumeric, err := parseId(id)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	secret, httpRes, err := client.NetworkDeviceDefaultSecretsAPI.GetNetworkDeviceDefaultSecretsInfo(ctx, idNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(secret, &networkDeviceDefaultSecretsPrintConfig)
}

func NetworkDeviceDefaultSecretsGetCredentials(ctx context.Context, id string) error {
	logger.Get().Info().Msgf("Getting network device default secrets credentials for '%s'", id)

	idNumeric, err := parseId(id)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	credentials, httpRes, err := client.NetworkDeviceDefaultSecretsAPI.GetNetworkDeviceDefaultSecretsCredentials(ctx, idNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, nil)
}

func NetworkDeviceDefaultSecretsCreate(ctx context.Context, siteId float32, macOrSerial string, secretName string, secretValue string) error {
	logger.Get().Info().Msgf("Creating network device default secrets for '%s'", macOrSerial)

	client := api.GetApiClient(ctx)

	body := sdk.NewCreateNetworkDeviceDefaultSecrets(siteId, macOrSerial, secretName, secretValue)

	secret, httpRes, err := client.NetworkDeviceDefaultSecretsAPI.
		CreateNetworkDeviceDefaultSecrets(ctx).
		CreateNetworkDeviceDefaultSecrets(*body).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(secret, &networkDeviceDefaultSecretsPrintConfig)
}

func NetworkDeviceDefaultSecretsUpdate(ctx context.Context, id string, secretValue string) error {
	logger.Get().Info().Msgf("Updating network device default secrets '%s'", id)

	idNumeric, err := parseId(id)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	body := sdk.UpdateNetworkDeviceDefaultSecrets{
		SecretValue: &secretValue,
	}

	secret, httpRes, err := client.NetworkDeviceDefaultSecretsAPI.
		UpdateNetworkDeviceDefaultSecrets(ctx, idNumeric).
		UpdateNetworkDeviceDefaultSecrets(body).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(secret, &networkDeviceDefaultSecretsPrintConfig)
}

func NetworkDeviceDefaultSecretsDelete(ctx context.Context, id string) error {
	logger.Get().Info().Msgf("Deleting network device default secrets '%s'", id)

	idNumeric, err := parseId(id)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceDefaultSecretsAPI.DeleteNetworkDeviceDefaultSecrets(ctx, idNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network device default secrets '%s' deleted successfully", id)
	return nil
}

func parseId(id string) (float32, error) {
	idNumeric, err := strconv.ParseFloat(id, 32)
	if err != nil {
		err := fmt.Errorf("invalid network device default secrets ID: '%s'", id)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}
	return float32(idNumeric), nil
}
