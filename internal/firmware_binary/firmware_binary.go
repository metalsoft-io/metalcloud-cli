package firmware_binary

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

var firmwareBinaryPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Name": {
			MaxWidth: 30,
			Order:    2,
		},
		"CatalogId": {
			Title: "Catalog ID",
			Order: 3,
		},
		"PackageVersion": {
			Title: "Version",
			Order: 4,
		},
		"UpdateSeverity": {
			Title: "Severity",
			Order: 5,
		},
		"RebootRequired": {
			Title: "Reboot Required",
			Order: 6,
		},
		"VendorReleaseTimestamp": {
			Title:       "Vendor Release",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
		"PackageId": {
			Title: "Package ID",
			Order: 8,
		},
		"ExternalId": {
			Title: "External ID",
			Order: 9,
		},
		"VendorDownloadUrl": {
			Title: "Download URL",
			Order: 10,
		},
		"VendorInfoUrl": {
			Title: "Info URL",
			Order: 11,
		},
	},
}

func FirmwareBinaryList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all firmware binaries")

	client := api.GetApiClient(ctx)

	firmwareBinaryList, httpRes, err := client.FirmwareBinaryAPI.GetFirmwareBinaries(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(firmwareBinaryList, &firmwareBinaryPrintConfig)
}

func FirmwareBinaryGet(ctx context.Context, firmwareBinaryId string) error {
	logger.Get().Info().Msgf("Get firmware binary '%s' details", firmwareBinaryId)

	firmwareBinaryIdNumeric, err := getFirmwareBinaryId(firmwareBinaryId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	firmwareBinary, httpRes, err := client.FirmwareBinaryAPI.GetFirmwareBinary(ctx, firmwareBinaryIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(firmwareBinary, &firmwareBinaryPrintConfig)
}

func FirmwareBinaryCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating firmware binary")

	var firmwareBinaryConfig sdk.CreateFirmwareBinary
	err := json.Unmarshal(config, &firmwareBinaryConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	firmwareBinary, httpRes, err := client.FirmwareBinaryAPI.
		CreateFirmwareBinary(ctx).
		CreateFirmwareBinary(firmwareBinaryConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(firmwareBinary, &firmwareBinaryPrintConfig)
}

func FirmwareBinaryDelete(ctx context.Context, firmwareBinaryId string) error {
	logger.Get().Info().Msgf("Deleting firmware binary '%s'", firmwareBinaryId)

	firmwareBinaryIdNumeric, err := getFirmwareBinaryId(firmwareBinaryId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.FirmwareBinaryAPI.
		DeleteFirmwareBinary(ctx, firmwareBinaryIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Firmware binary '%s' deleted", firmwareBinaryId)
	return nil
}

func FirmwareBinaryConfigExample(ctx context.Context) error {
	// Example create firmware binary configuration
	firmwareBinaryConfiguration := sdk.CreateFirmwareBinary{
		Name:                   "BIOS-R740-2.15.0",
		CatalogId:              1,
		VendorDownloadUrl:      "https://dell.com/downloads/firmware/R740/BIOS-2.15.0.bin",
		VendorInfoUrl:          sdk.PtrString("https://dell.com/support/firmware/R740/BIOS/2.15.0"),
		ExternalId:             sdk.PtrString("DELL-R740-BIOS-2.15.0"),
		PackageId:              sdk.PtrString("BIOS"),
		PackageVersion:         sdk.PtrString("2.15.0"),
		RebootRequired:         true,
		UpdateSeverity:         "recommended",
		VendorReleaseTimestamp: sdk.PtrString("2024-04-01T12:00:00Z"),
		VendorSupportedDevices: []map[string]interface{}{
			{
				"model": "PowerEdge R740",
				"type":  "server",
			},
		},
		VendorSupportedSystems: []map[string]interface{}{
			{
				"os":      "any",
				"version": "any",
			},
		},
		Vendor: map[string]interface{}{
			"name":    "Dell Inc.",
			"contact": "support@dell.com",
		},
	}

	return formatter.PrintResult(firmwareBinaryConfiguration, nil)
}

func getFirmwareBinaryId(firmwareBinaryId string) (float32, error) {
	firmwareBinaryIdNumeric, err := strconv.ParseFloat(firmwareBinaryId, 32)
	if err != nil {
		err := fmt.Errorf("invalid firmware binary ID: '%s'", firmwareBinaryId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(firmwareBinaryIdNumeric), nil
}
