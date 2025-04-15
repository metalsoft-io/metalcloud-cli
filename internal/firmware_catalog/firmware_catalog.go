package firmware_catalog

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

var firmwareCatalogPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Name": {
			MaxWidth: 30,
			Order:    2,
		},
		"Vendor": {
			MaxWidth: 20,
			Order:    3,
		},
		"UpdateType": {
			Title: "Update Type",
			Order: 4,
		},
		"Description": {
			MaxWidth: 40,
			Order:    5,
		},
		"VendorId": {
			Title: "Vendor ID",
			Order: 6,
		},
		"VendorUrl": {
			Title: "Vendor URL",
			Order: 7,
		},
		"VendorReleaseTimestamp": {
			Title:       "Vendor Release",
			Transformer: formatter.FormatDateTimeValue,
			Order:       8,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       9,
		},
	},
}

func FirmwareCatalogList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all firmware catalogs")

	client := api.GetApiClient(ctx)

	firmwareCatalogList, httpRes, err := client.FirmwareCatalogAPI.GetFirmwareCatalogs(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(firmwareCatalogList, &firmwareCatalogPrintConfig)
}

func FirmwareCatalogGet(ctx context.Context, firmwareCatalogId string) error {
	logger.Get().Info().Msgf("Get firmware catalog '%s' details", firmwareCatalogId)

	firmwareCatalogIdNumeric, err := getFirmwareCatalogId(firmwareCatalogId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	firmwareCatalog, httpRes, err := client.FirmwareCatalogAPI.GetFirmwareCatalog(ctx, firmwareCatalogIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(firmwareCatalog, &firmwareCatalogPrintConfig)
}

func FirmwareCatalogCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating firmware catalog")

	var firmwareCatalogConfig sdk.CreateFirmwareCatalog
	err := json.Unmarshal(config, &firmwareCatalogConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	firmwareCatalog, httpRes, err := client.FirmwareCatalogAPI.
		CreateFirmwareCatalogs(ctx).
		CreateFirmwareCatalog(firmwareCatalogConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(firmwareCatalog, &firmwareCatalogPrintConfig)
}

func FirmwareCatalogUpdate(ctx context.Context, firmwareCatalogId string, config []byte) error {
	logger.Get().Info().Msgf("Updating firmware catalog '%s'", firmwareCatalogId)

	firmwareCatalogIdNumeric, err := getFirmwareCatalogId(firmwareCatalogId)
	if err != nil {
		return err
	}

	var firmwareCatalogConfig sdk.UpdateFirmwareCatalog
	err = json.Unmarshal(config, &firmwareCatalogConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	firmwareCatalog, httpRes, err := client.FirmwareCatalogAPI.
		UpdateFirmwareCatalog(ctx, firmwareCatalogIdNumeric).
		UpdateFirmwareCatalog(firmwareCatalogConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(firmwareCatalog, &firmwareCatalogPrintConfig)
}

func FirmwareCatalogDelete(ctx context.Context, firmwareCatalogId string) error {
	logger.Get().Info().Msgf("Deleting firmware catalog '%s'", firmwareCatalogId)

	firmwareCatalogIdNumeric, err := getFirmwareCatalogId(firmwareCatalogId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.FirmwareCatalogAPI.
		DeleteFirmwareCatalog(ctx, firmwareCatalogIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Firmware catalog '%s' deleted", firmwareCatalogId)
	return nil
}

func FirmwareCatalogConfigExample(ctx context.Context) error {
	// Example create firmware catalog configuration
	firmwareCatalogConfiguration := sdk.CreateFirmwareCatalog{
		Name:        "example-firmware-catalog",
		Description: sdk.PtrString("Example firmware catalog for Dell servers"),
		Vendor:      "dell",
		VendorId:    sdk.PtrString("R740"),
		VendorUrl:   sdk.PtrString("https://dell.com/support/firmware/R740"),
		UpdateType:  "online",
		VendorConfiguration: map[string]interface{}{
			"credentials": map[string]interface{}{
				"username": "api_user",
				"api_key":  "API_KEY_HERE",
			},
			"update_frequency": "daily",
		},
		MetalsoftServerTypesSupported: []string{"dell_r740", "dell_r640"},
		VendorServerTypesSupported:    []string{"poweredge_r740", "poweredge_r640"},
	}

	return formatter.PrintResult(firmwareCatalogConfiguration, nil)
}

func getFirmwareCatalogId(firmwareCatalogId string) (float32, error) {
	firmwareCatalogIdNumeric, err := strconv.ParseFloat(firmwareCatalogId, 32)
	if err != nil {
		err := fmt.Errorf("invalid firmware catalog ID: '%s'", firmwareCatalogId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(firmwareCatalogIdNumeric), nil
}
