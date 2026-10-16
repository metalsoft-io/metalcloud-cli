package firmware_binary

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var firmwareBinaryPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
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

// ListFlags holds the server-side filters and pagination options of `firmware-binary list`.
// Every filter maps 1:1 to a `filter.*` query parameter of GET /api/v2/firmware/binary.
type ListFlags struct {
	FilterId                     []string
	FilterCatalogId              []string
	FilterExternalId             []string
	FilterPackageId              []string
	FilterPackageVersion         []string
	FilterUpdateSeverity         []string
	FilterRebootRequired         []string
	FilterCreatedTimestamp       []string
	FilterVendorReleaseTimestamp []string

	Page  int
	Limit int
}

func FirmwareBinaryList(ctx context.Context, flags ListFlags) error {
	logger.Get().Info().Msgf("Listing firmware binaries")

	client := api.GetApiClient(ctx)

	request := client.FirmwareBinaryAPI.GetFirmwareBinaries(ctx).SortBy([]string{"id:ASC"})

	apply := func(values []string, set func([]string) sdk.FirmwareBinaryAPIGetFirmwareBinariesRequest) {
		if len(values) > 0 {
			request = set(utils.ProcessFilterStringSlice(values))
		}
	}

	apply(flags.FilterId, request.FilterId)
	apply(flags.FilterCatalogId, request.FilterCatalogId)
	apply(flags.FilterExternalId, request.FilterExternalId)
	apply(flags.FilterPackageId, request.FilterPackageId)
	apply(flags.FilterPackageVersion, request.FilterPackageVersion)
	apply(flags.FilterUpdateSeverity, request.FilterUpdateSeverity)
	apply(flags.FilterRebootRequired, request.FilterRebootRequired)
	apply(flags.FilterCreatedTimestamp, request.FilterCreatedTimestamp)
	apply(flags.FilterVendorReleaseTimestamp, request.FilterVendorReleaseTimestamp)

	switch {
	case flags.Page > 0:
		records, meta, err := utils.FetchPageWindow(request, flags.Page, flags.Limit)
		if err != nil {
			return err
		}
		return utils.PrintAll(records, meta, len(records), &firmwareBinaryPrintConfig)
	case flags.Limit > 0:
		records, meta, err := utils.FetchUpTo(request, flags.Limit)
		if err != nil {
			return err
		}
		return utils.PrintAll(records, meta, len(records), &firmwareBinaryPrintConfig)
	default:
		records, meta, err := utils.FetchAllPages(request)
		if err != nil {
			return err
		}
		return utils.PrintAll(records, meta, len(records), &firmwareBinaryPrintConfig)
	}
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
	err := utils.UnmarshalContent(config, &firmwareBinaryConfig)
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
	}

	return formatter.PrintResult(firmwareBinaryConfiguration, nil)
}

func getFirmwareBinaryId(firmwareBinaryId string) (int64, error) {
	firmwareBinaryIdNumeric, err := strconv.ParseInt(firmwareBinaryId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid firmware binary ID: '%s'", firmwareBinaryId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return firmwareBinaryIdNumeric, nil
}
