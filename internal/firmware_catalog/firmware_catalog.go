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
	// First unmarshal into a more generic map to extract ID if needed
	var rawConfig map[string]interface{}
	if err := json.Unmarshal(config, &rawConfig); err != nil {
		return err
	}

	// If firmwareCatalogId is empty, try to get it from the config
	if firmwareCatalogId == "" {
		if id, exists := rawConfig["id"]; exists && id != nil {
			firmwareCatalogId = fmt.Sprintf("%v", id)
			logger.Get().Info().Msgf("Using firmware catalog ID '%s' from config", firmwareCatalogId)
		} else {
			return fmt.Errorf("firmware catalog ID is required either as parameter or in the config file")
		}
	} else {
		if id, exists := rawConfig["id"]; exists && id != nil {
			if id != firmwareCatalogId {
				return fmt.Errorf("firmware catalog ID '%s' in config does not match the provided ID '%s'", id, firmwareCatalogId)
			}
		}
	}

	logger.Get().Info().Msgf("Updating firmware catalog '%s'", firmwareCatalogId)

	firmwareCatalogIdNumeric, err := getFirmwareCatalogId(firmwareCatalogId)
	if err != nil {
		return err
	}

	// Try to unmarshal into the proper update structure
	var firmwareCatalogConfig sdk.UpdateFirmwareCatalog
	err = json.Unmarshal(config, &firmwareCatalogConfig)

	// If unmarshaling into UpdateFirmwareCatalog fails, try with FirmwareCatalog
	if err != nil || len(firmwareCatalogConfig.AdditionalProperties) > 0 {
		var fullCatalog sdk.FirmwareCatalog
		if err := json.Unmarshal(config, &fullCatalog); err != nil {
			// If both fail, return an error
			return fmt.Errorf("firmware catalog config does not match the expected format")
		}

		// Copy relevant fields from FirmwareCatalog to UpdateFirmwareCatalog
		firmwareCatalogConfig = sdk.UpdateFirmwareCatalog{}

		firmwareCatalogConfig.Name = fullCatalog.Name
		firmwareCatalogConfig.Description = fullCatalog.Description
		firmwareCatalogConfig.Vendor = sdk.FirmwareVendorType(fullCatalog.Vendor)
		firmwareCatalogConfig.VendorId = fullCatalog.VendorId
		firmwareCatalogConfig.VendorUrl = fullCatalog.VendorUrl
		firmwareCatalogConfig.UpdateType = sdk.CatalogUpdateType(fullCatalog.UpdateType)
		firmwareCatalogConfig.VendorConfiguration = fullCatalog.VendorConfiguration
		firmwareCatalogConfig.MetalsoftServerTypesSupported = fullCatalog.MetalsoftServerTypesSupported
		firmwareCatalogConfig.VendorServerTypesSupported = fullCatalog.VendorServerTypesSupported

		logger.Get().Info().Msg("Converted FirmwareCatalog to UpdateFirmwareCatalog format")
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
		Name:        "DELL Enterprise catalog",
		Description: sdk.PtrString("The catalog contains the latest BIOS, firmware, drivers, and certain applications for both Microsoft Windows and Linux operating systems."),
		Vendor:      "dell",
		VendorId:    sdk.PtrString("48912fae-2b46-4b4c-bafe-2c709c7b0ad2"),
		VendorUrl:   sdk.PtrString("https://downloads.dell.com/catalog/Catalog.gz"),
		UpdateType:  "online",
		VendorConfiguration: map[string]interface{}{
			"update_frequency": "daily",
		},
		MetalsoftServerTypesSupported: []string{"M.4.8.2", "M.4.16.2"},
		VendorServerTypesSupported:    []string{"R740", "R640"},
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
