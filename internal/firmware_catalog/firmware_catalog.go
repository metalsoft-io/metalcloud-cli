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

type FirmwareCatalogCreateOptions struct {
	Name                    string   `json:"name"`
	Description             string   `json:"description,omitempty"`
	Vendor                  string   `json:"vendor"`
	UpdateType              string   `json:"update_type"`
	VendorUrl               string   `json:"vendor_url,omitempty"`
	VendorLocalCatalogPath  string   `json:"vendor_local_catalog_path,omitempty"`
	VendorLocalBinariesPath string   `json:"vendor_local_binaries_path,omitempty"`
	ServerTypesFilter       []string `json:"server_types_filter,omitempty"`
	VendorSystemsFilter     []string `json:"vendor_systems_filter,omitempty"`
	DownloadBinaries        bool     `json:"download_binaries,omitempty"`
	UploadBinaries          bool     `json:"upload_binaries,omitempty"`
	RepoBaseUrl             string   `json:"repo_base_url,omitempty"`
	RepoSshHost             string   `json:"repo_ssh_host,omitempty"`
	RepoSshUser             string   `json:"repo_ssh_user,omitempty"`
	UserPrivateKeyPath      string   `json:"user_private_key_path,omitempty"`
	KnownHostsPath          string   `json:"known_hosts_path,omitempty"`
	IgnoreHostKeyCheck      bool     `json:"ignore_host_key_check,omitempty"`
}

func FirmwareCatalogCreate(ctx context.Context, firmwareCatalogOptions FirmwareCatalogCreateOptions) error {
	logger.Get().Info().Msgf("Creating firmware catalog")

	vendorCatalog, err := NewVendorCatalogFromCreateOptions(firmwareCatalogOptions)
	if err != nil {
		return err
	}

	err = vendorCatalog.ProcessVendorCatalog(ctx)
	if err != nil {
		return err
	}

	err = vendorCatalog.CreateMetalsoftCatalog(ctx)
	if err != nil {
		return err
	}

	return FirmwareCatalogGet(ctx, fmt.Sprintf("%d", int(vendorCatalog.CatalogInfo.Id)))
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

func getFirmwareCatalogId(firmwareCatalogId string) (float32, error) {
	firmwareCatalogIdNumeric, err := strconv.ParseFloat(firmwareCatalogId, 32)
	if err != nil {
		err := fmt.Errorf("invalid firmware catalog ID: '%s'", firmwareCatalogId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(firmwareCatalogIdNumeric), nil
}
