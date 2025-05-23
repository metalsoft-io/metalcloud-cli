package firmware_baseline

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

var firmwareBaselinePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Name": {
			MaxWidth: 30,
			Order:    2,
		},
		"Level": {
			Order: 3,
		},
		"LevelFilter": {
			Title:    "Level Filter",
			MaxWidth: 30,
			Order:    4,
		},
		"Description": {
			MaxWidth: 40,
			Order:    5,
		},
		"Catalog": {
			MaxWidth: 30,
			Order:    6,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
	},
}

func FirmwareBaselineList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all firmware baselines")

	client := api.GetApiClient(ctx)

	firmwareBaselineList, httpRes, err := client.FirmwareBaselineAPI.GetFirmwareBaselines(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(firmwareBaselineList, &firmwareBaselinePrintConfig)
}

func FirmwareBaselineGet(ctx context.Context, firmwareBaselineId string) error {
	logger.Get().Info().Msgf("Get firmware baseline '%s' details", firmwareBaselineId)

	firmwareBaselineIdNumeric, err := getFirmwareBaselineId(firmwareBaselineId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	firmwareBaseline, httpRes, err := client.FirmwareBaselineAPI.GetFirmwareBaseline(ctx, firmwareBaselineIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(firmwareBaseline, &firmwareBaselinePrintConfig)
}

func FirmwareBaselineCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating firmware baseline")

	var firmwareBaselineConfig sdk.CreateFirmwareBaseline
	err := json.Unmarshal(config, &firmwareBaselineConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	firmwareBaseline, httpRes, err := client.FirmwareBaselineAPI.
		CreateFirmwareBaseline(ctx).
		CreateFirmwareBaseline(firmwareBaselineConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(firmwareBaseline, &firmwareBaselinePrintConfig)
}

func FirmwareBaselineUpdate(ctx context.Context, firmwareBaselineId string, config []byte) error {
	logger.Get().Info().Msgf("Updating firmware baseline '%s'", firmwareBaselineId)

	firmwareBaselineIdNumeric, err := getFirmwareBaselineId(firmwareBaselineId)
	if err != nil {
		return err
	}

	var firmwareBaselineConfig sdk.UpdateFirmwareBaseline
	err = json.Unmarshal(config, &firmwareBaselineConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	firmwareBaseline, httpRes, err := client.FirmwareBaselineAPI.
		UpdateFirmwareBaseline(ctx, firmwareBaselineIdNumeric).
		UpdateFirmwareBaseline(firmwareBaselineConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(firmwareBaseline, &firmwareBaselinePrintConfig)
}

func FirmwareBaselineDelete(ctx context.Context, firmwareBaselineId string) error {
	logger.Get().Info().Msgf("Deleting firmware baseline '%s'", firmwareBaselineId)

	firmwareBaselineIdNumeric, err := getFirmwareBaselineId(firmwareBaselineId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.FirmwareBaselineAPI.
		DeleteFirmwareBaseline(ctx, firmwareBaselineIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Firmware baseline '%s' deleted", firmwareBaselineId)
	return nil
}

func FirmwareBaselineConfigExample(ctx context.Context) error {
	// Example create firmware baseline configuration
	firmwareBaselineConfiguration := sdk.CreateFirmwareBaseline{
		Name:        "example-firmware-baseline",
		Description: sdk.PtrString("Example firmware baseline for production servers"),
		Level:       "PRODUCTION",
		LevelFilter: []string{"dell_r740", "dell_r640"},
		Catalog:     []string{"catalog-1", "catalog-2"},
	}

	return formatter.PrintResult(firmwareBaselineConfiguration, nil)
}

func FirmwareBaselineSearch(ctx context.Context, searchCriteria []byte) error {
	logger.Get().Info().Msgf("Searching firmware baselines")

	var searchParams sdk.SearchFirmwareBinary
	err := json.Unmarshal(searchCriteria, &searchParams)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	results, httpRes, err := client.FirmwareBaselineSearchAPI.
		SearchFirmwareBaselines(ctx).
		SearchFirmwareBinary(searchParams).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(results, &firmwareBaselinePrintConfig)
}

func FirmwareBaselineSearchExample(ctx context.Context) error {
	// Example search criteria
	searchCriteria := sdk.SearchFirmwareBinary{
		Vendor: sdk.FirmwareVendorType("DELL"),
		BaselineFilter: sdk.BaselineFilter{
			Datacenter: []string{"datacenter-1"},
			ServerType: []string{"dell_r740"},
			OsTemplate: []string{"ubuntu-20.04"},
			BaselineId: []string{"baseline-1"},
		},
		ServerComponentFilter: &sdk.SearchFirmwareBinaryServerComponentFilter{
			DellComponentFilter: &sdk.DellComponentFilter{
				ComponentId: "component-1",
			},
		},
	}

	return formatter.PrintResult(searchCriteria, nil)
}

func getFirmwareBaselineId(firmwareBaselineId string) (float32, error) {
	firmwareBaselineIdNumeric, err := strconv.ParseFloat(firmwareBaselineId, 32)
	if err != nil {
		err := fmt.Errorf("invalid firmware baseline ID: '%s'", firmwareBaselineId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(firmwareBaselineIdNumeric), nil
}
