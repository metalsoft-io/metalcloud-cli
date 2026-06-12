package firmware_baseline

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

type firmwareBaselineRaw struct {
	Id               interface{} `json:"id"`
	Name             *string     `json:"name"`
	Description      *string     `json:"description"`
	Catalog          interface{} `json:"catalog"`
	CreatedTimestamp interface{} `json:"createdTimestamp"`
}

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
		"Description": {
			MaxWidth: 40,
			Order:    3,
		},
		"Catalog": {
			MaxWidth: 30,
			Order:    4,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       5,
		},
	},
}

func FirmwareBaselineList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all firmware baselines")

	client := api.GetApiClient(ctx)

	rawItems, meta, err := utils.FetchAllPagesRaw(func(p float32) (*http.Response, error) {
		_, httpRes, _ := client.FirmwareBaselineAPI.GetFirmwareBaselines(ctx).SortBy([]string{"id:ASC"}).Page(p).Limit(100).Execute()
		return httpRes, nil
	})
	if err != nil {
		return err
	}
	records, err := utils.UnmarshalRawItems[firmwareBaselineRaw](rawItems)
	if err != nil {
		return fmt.Errorf("failed to parse firmware baselines: %w", err)
	}

	return utils.PrintAllRaw(rawItems, records, meta, len(records), &firmwareBaselinePrintConfig)
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
	err := utils.UnmarshalContent(config, &firmwareBaselineConfig)
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
	err = utils.UnmarshalContent(config, &firmwareBaselineConfig)
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
		MinimumVersions: []sdk.FirmwareMinimumVersion{
			sdk.FirmwareMinimumVersion{ComponentName: "BMC", MinimumVersion: "1.0"},
		},
		Catalog: []string{"catalog-1", "catalog-2"},
	}

	return formatter.PrintResult(firmwareBaselineConfiguration, nil)
}

func FirmwareBaselineSearch(ctx context.Context, searchCriteria []byte) error {
	logger.Get().Info().Msgf("Searching firmware baselines")

	var searchParams sdk.SearchFirmwareBinary
	err := utils.UnmarshalContent(searchCriteria, &searchParams)
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
		Vendor:      sdk.SERVERFIRMWARECATALOGVENDOR_DELL,
		BaselineIds: []string{"baseline-1"},
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
