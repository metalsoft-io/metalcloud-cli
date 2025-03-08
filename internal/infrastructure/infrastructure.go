package infrastructure

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

var infrastructurePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
		},
		"Label": {},
		"ServiceStatus": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
		},
		"UserIdOwner": {
			Title: "Owner",
		},
		"SiteId": {
			Title: "Site",
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
		},
	},
}

func InfrastructureList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all infrastructures")

	client := api.GetApiClient(ctx)

	infrastructureList, httpRes, err := client.InfrastructureAPI.GetInfrastructures(ctx).SortBy([]string{"id:ASC"}).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(infrastructureList, &infrastructurePrintConfig)
}

func InfrastructureGet(ctx context.Context, infrastructureId string) error {
	logger.Get().Info().Msgf("Get infrastructure '%s'", infrastructureId)

	infrastructureIdNumber, err := strconv.ParseFloat(infrastructureId, 32)
	if err != nil {
		err := fmt.Errorf("invalid infrastructure ID: '%s'", infrastructureId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	infrastructureInfo, httpRes, err := client.InfrastructureAPI.GetInfrastructure(ctx, float32(infrastructureIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(infrastructureInfo, &infrastructurePrintConfig)
}

func InfrastructureCreate(ctx context.Context, infrastructureLabel string, infrastructureDescription string, siteId string) error {
	logger.Get().Info().Msgf("Create infrastructure '%s'", infrastructureLabel)

	siteIdNumber, err := strconv.ParseFloat(siteId, 32)
	if err != nil {
		err := fmt.Errorf("invalid site ID: '%s'", siteId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	createInfrastructure := sdk.InfrastructureCreate{
		Label:       infrastructureLabel,
		Description: &infrastructureDescription,
		SiteId:      float32(siteIdNumber),
	}

	client := api.GetApiClient(ctx)

	infrastructureInfo, httpRes, err := client.InfrastructureAPI.CreateInfrastructure(ctx).InfrastructureCreate(createInfrastructure).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(infrastructureInfo, &infrastructurePrintConfig)
}

func InfrastructureUpdate(ctx context.Context, infrastructureId string) error {
	logger.Get().Info().Msgf("Update infrastructure '%s'", infrastructureId)

	infrastructureIdNumber, err := strconv.ParseFloat(infrastructureId, 32)
	if err != nil {
		err := fmt.Errorf("invalid infrastructure ID: '%s'", infrastructureId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	infrastructureInfo, httpRes, err := client.InfrastructureAPI.UpdateInfrastructureConfiguration(ctx, float32(infrastructureIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(infrastructureInfo, &infrastructurePrintConfig)
}
