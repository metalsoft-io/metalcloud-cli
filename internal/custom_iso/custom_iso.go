package custom_iso

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/internal/server"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var customIsoPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Name": {
			MaxWidth: 30,
			Order:    2,
		},
		"DisplayName": {
			Title:    "Display Name",
			MaxWidth: 30,
			Order:    3,
		},
		"AccessUrl": {
			Title: "Access URL",
			Order: 4,
		},
		"IsPublic": {
			Title: "Public",
			Order: 5,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
	},
}

func CustomIsoList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all custom ISOs")

	client := api.GetApiClient(ctx)

	customIsoList, httpRes, err := client.CustomIsoAPI.GetCustomIsos(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(customIsoList, &customIsoPrintConfig)
}

func CustomIsoGet(ctx context.Context, customIsoId string) error {
	logger.Get().Info().Msgf("Get custom ISO '%s' details", customIsoId)

	customIsoIdNumeric, err := getCustomIsoId(customIsoId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	customIso, httpRes, err := client.CustomIsoAPI.GetCustomIso(ctx, customIsoIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(customIso, &customIsoPrintConfig)
}

func CustomIsoCreate(ctx context.Context, customIsoConfig sdk.CreateCustomIso) error {
	logger.Get().Info().Msgf("Creating custom ISO")

	client := api.GetApiClient(ctx)

	customIso, httpRes, err := client.CustomIsoAPI.
		CreateCustomIso(ctx).
		CreateCustomIso(customIsoConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(customIso, &customIsoPrintConfig)
}

func CustomIsoUpdate(ctx context.Context, customIsoId string, config []byte) error {
	logger.Get().Info().Msgf("Updating custom ISO '%s'", customIsoId)

	customIsoIdNumeric, err := getCustomIsoId(customIsoId)
	if err != nil {
		return err
	}

	var customIsoConfig sdk.UpdateCustomIso
	err = utils.UnmarshalContent(config, &customIsoConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	customIso, httpRes, err := client.CustomIsoAPI.
		UpdateCustomIso(ctx, customIsoIdNumeric).
		UpdateCustomIso(customIsoConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(customIso, &customIsoPrintConfig)
}

func CustomIsoDelete(ctx context.Context, customIsoId string) error {
	logger.Get().Info().Msgf("Deleting custom ISO '%s'", customIsoId)

	customIsoIdNumeric, err := getCustomIsoId(customIsoId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.CustomIsoAPI.
		DeleteCustomIso(ctx, customIsoIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Custom ISO '%s' deleted", customIsoId)
	return nil
}

func CustomIsoMakePublic(ctx context.Context, customIsoId string) error {
	logger.Get().Info().Msgf("Making custom ISO '%s' public", customIsoId)

	customIsoIdNumeric, err := getCustomIsoId(customIsoId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	customIso, httpRes, err := client.CustomIsoAPI.
		MakeCustomIsoPublic(ctx, customIsoIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Custom ISO '%s' is now public", customIsoId)
	return formatter.PrintResult(customIso, &customIsoPrintConfig)
}

func CustomIsoBootServer(ctx context.Context, customIsoId string, serverId string) error {
	logger.Get().Info().Msgf("Booting server '%s' with custom ISO '%s'", serverId, customIsoId)

	customIsoIdNumeric, err := getCustomIsoId(customIsoId)
	if err != nil {
		return err
	}

	serverIdNumeric, err := server.GetServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	jobInfo, httpRes, err := client.CustomIsoAPI.
		BootCustomIsoIntoServer(ctx, customIsoIdNumeric, serverIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Server '%s' boot with custom ISO '%s' initiated", serverId, customIsoId)
	return formatter.PrintResult(jobInfo, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"JobId": {
				Title: "Job ID",
				Order: 1,
			},
			"JobGroupId": {
				Title: "Job Group ID",
				Order: 2,
			},
		},
	})
}

func CustomIsoConfigExample(ctx context.Context) error {
	// Example create custom ISO configuration
	customIsoConfiguration := sdk.CreateCustomIso{
		Label:     "example-iso",
		Name:      sdk.PtrString("Example Custom ISO"),
		AccessUrl: "http://example.com/isos/example.iso",
	}

	return formatter.PrintResult(customIsoConfiguration, nil)
}

func getCustomIsoId(customIsoId string) (float32, error) {
	customIsoIdNumeric, err := strconv.ParseFloat(customIsoId, 32)
	if err != nil {
		err := fmt.Errorf("invalid custom ISO ID: '%s'", customIsoId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(customIsoIdNumeric), nil
}
