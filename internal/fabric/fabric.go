package fabric

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

var fabricPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
		},
		"Name":        {},
		"Description": {},
		"FabricConfiguration.FabricType": {
			Title: "Type",
		},
	},
}

func FabricList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all fabrics")

	client := api.GetApiClient(ctx)

	fabricList, httpRes, err := client.NetworkFabricAPI.GetNetworkFabrics(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(fabricList, &fabricPrintConfig)
}

func FabricGet(ctx context.Context, fabricId string) error {
	logger.Get().Info().Msgf("Get fabric '%s'", fabricId)

	fabricIdNumber, err := strconv.ParseFloat(fabricId, 32)
	if err != nil {
		err := fmt.Errorf("invalid fabric ID: '%s'", fabricId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	fabricInfo, httpRes, err := client.NetworkFabricAPI.GetNetworkFabricById(ctx, float32(fabricIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(fabricInfo, &fabricPrintConfig)
}

func FabricCreate(ctx context.Context, fabricName string, fabricDescription string, fabricType string) error {
	logger.Get().Info().Msgf("Create fabric '%s'", fabricName)

	var fabricConfiguration sdk.NetworkFabricFabricConfiguration
	switch fabricType {
	case "ethernet":
		fabricConfiguration = sdk.NetworkFabricFabricConfiguration{
			EthernetFabric: &sdk.EthernetFabric{
				FabricType: "ethernet",
			},
		}
	case "fibre_channel":
		fabricConfiguration = sdk.NetworkFabricFabricConfiguration{
			FibreChannelFabric: &sdk.FibreChannelFabric{
				FabricType: "fibre_channel",
			},
		}
	default:
		err := fmt.Errorf("invalid fabric type: '%s'", fabricType)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	createFabric := sdk.CreateNetworkFabric{
		Name:                fabricName,
		Description:         &fabricDescription,
		FabricConfiguration: fabricConfiguration,
	}

	client := api.GetApiClient(ctx)

	fabricInfo, httpRes, err := client.NetworkFabricAPI.CreateNetworkFabric(ctx).CreateNetworkFabric(createFabric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(fabricInfo, &fabricPrintConfig)
}

func FabricUpdate(ctx context.Context, fabricId string) error {
	logger.Get().Info().Msgf("Update fabric '%s'", fabricId)

	fabricIdNumber, err := strconv.ParseFloat(fabricId, 32)
	if err != nil {
		err := fmt.Errorf("invalid fabric ID: '%s'", fabricId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	fabricInfo, httpRes, err := client.NetworkFabricAPI.GetNetworkFabricById(ctx, float32(fabricIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(fabricInfo, &fabricPrintConfig)
}
