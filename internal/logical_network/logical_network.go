package logical_network

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/internal/fabric"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var logicalNetworkPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Label": {
			MaxWidth: 30,
			Order:    2,
		},
		"Name": {
			MaxWidth: 30,
			Order:    3,
		},
		"Description": {
			MaxWidth: 30,
			Order:    4,
		},
		"LogicalNetworkType": {
			Title: "Type",
			Order: 5,
		},
		"FabricId": {
			Title: "Fabric ID",
			Order: 6,
		},
		"InfrastructureId": {
			Title: "Infra ID",
			Order: 7,
		},
	},
}

func LogicalNetworkList(ctx context.Context, fabricIdOrLabel string) error {
	logger.Get().Info().Msgf("Listing logical networks")

	client := api.GetApiClient(ctx)

	request := client.LogicalNetworksAPI.GetAllLogicalNetworks(ctx)

	if fabricIdOrLabel != "" {
		fabric, err := fabric.GetFabricByIdOrLabel(ctx, fabricIdOrLabel)
		if err != nil {
			return err
		}

		request = request.FilterFabricId([]string{fabric.Id})
	}

	logicalNetworkList, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(logicalNetworkList, &logicalNetworkPrintConfig)
}

func LogicalNetworkGet(ctx context.Context, logicalNetworkId string) error {
	logger.Get().Info().Msgf("Get logical network '%s' details", logicalNetworkId)

	logicalNetworkIdNumeric, err := getLogicalNetworkId(logicalNetworkId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	logicalNetwork, httpRes, err := client.LogicalNetworksAPI.GetLogicalNetworkById(ctx, logicalNetworkIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(logicalNetwork, &logicalNetworkPrintConfig)
}

func LogicalNetworkCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating logical network")

	var logicalNetworkConfig sdk.CreateLogicalNetwork
	err := json.Unmarshal(config, &logicalNetworkConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	logicalNetwork, httpRes, err := client.LogicalNetworksAPI.
		CreateLogicalNetwork(ctx).
		CreateLogicalNetwork(logicalNetworkConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(logicalNetwork, &logicalNetworkPrintConfig)
}

func LogicalNetworkDelete(ctx context.Context, logicalNetworkId string) error {
	logger.Get().Info().Msgf("Deleting logical network '%s'", logicalNetworkId)

	logicalNetworkIdNumeric, err := getLogicalNetworkId(logicalNetworkId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.LogicalNetworksAPI.
		DeleteLogicalNetwork(ctx, logicalNetworkIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Logical network '%s' deleted", logicalNetworkId)
	return nil
}

func LogicalNetworkUpdate(ctx context.Context, logicalNetworkId string, config []byte) error {
	logger.Get().Info().Msgf("Updating logical network '%s'", logicalNetworkId)

	logicalNetworkIdNumeric, err := getLogicalNetworkId(logicalNetworkId)
	if err != nil {
		return err
	}

	var logicalNetworkUpdate sdk.UpdateLogicalNetwork
	err = json.Unmarshal(config, &logicalNetworkUpdate)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	logicalNetwork, httpRes, err := client.LogicalNetworksAPI.
		UpdateLogicalNetwork(ctx, logicalNetworkIdNumeric).
		UpdateLogicalNetwork(logicalNetworkUpdate).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(logicalNetwork, &logicalNetworkPrintConfig)
}

func LogicalNetworkConfigExample(ctx context.Context) error {
	logicalNetworkConfiguration := sdk.CreateLogicalNetwork{
		Label:              sdk.PtrString("example-logical-network"),
		Name:               sdk.PtrString("Example Logical Network"),
		Description:        sdk.PtrString("Example logical network description"),
		FabricId:           1,
		InfrastructureId:   sdk.PtrFloat32(1),
		LogicalNetworkType: "vlan",
		Annotations: map[string]interface{}{
			"example-key": "example-value",
		},
	}

	return formatter.PrintResult(logicalNetworkConfiguration, nil)
}

func getLogicalNetworkId(logicalNetworkId string) (float32, error) {
	logicalNetworkIdNumeric, err := strconv.ParseFloat(logicalNetworkId, 32)
	if err != nil {
		err := fmt.Errorf("invalid logical network ID: '%s'", logicalNetworkId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(logicalNetworkIdNumeric), nil
}
