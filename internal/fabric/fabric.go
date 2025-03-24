package fabric

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/internal/site"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var fabricPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Name": {
			Order: 2,
		},
		"Description": {
			Order: 3,
		},
		"FabricConfiguration": {
			Hidden: true,
			InnerFields: map[string]formatter.RecordFieldConfig{
				"EthernetFabric": {
					Hidden: true,
					InnerFields: map[string]formatter.RecordFieldConfig{
						"FabricType": {
							Title: "Type",
							Order: 4,
						},
						"DefaultVlan": {
							Title: "Default VLAN",
							Order: 5,
						},
					},
				},
			},
		},
	},
}

func FabricList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all fabrics")

	client := api.GetApiClient(ctx)

	fabricList, httpRes, err := client.NetworkFabricAPI.GetNetworkFabrics(ctx).SortBy([]string{"id:ASC"}).Execute()
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

func FabricCreate(ctx context.Context, siteIdOrLabel string, fabricName string, fabricType string, description string, config []byte) error {
	logger.Get().Info().Msgf("Create fabric '%s'", fabricName)

	site, err := site.GetSiteByIdOrLabel(ctx, siteIdOrLabel)
	if err != nil {
		return err
	}

	var fabricConfiguration sdk.NetworkFabricFabricConfiguration
	switch fabricType {
	case "ethernet":
		ethernetConfig := sdk.EthernetFabric{}
		err := json.Unmarshal(config, &ethernetConfig)
		if err != nil {
			return err
		}

		fabricConfiguration = sdk.NetworkFabricFabricConfiguration{
			EthernetFabric: &ethernetConfig,
		}
	case "fibre_channel":
		fcConfig := sdk.FibreChannelFabric{}
		err := json.Unmarshal(config, &fcConfig)
		if err != nil {
			return err
		}

		fabricConfiguration = sdk.NetworkFabricFabricConfiguration{
			FibreChannelFabric: &fcConfig,
		}
	default:
		err := fmt.Errorf("invalid fabric type: '%s'", fabricType)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	createFabric := sdk.CreateNetworkFabric{
		Name:                fabricName,
		Description:         sdk.PtrString(description),
		SiteId:              sdk.PtrFloat32(float32(site.Id)),
		FabricConfiguration: fabricConfiguration,
	}

	client := api.GetApiClient(ctx)

	fabricInfo, httpRes, err := client.NetworkFabricAPI.CreateNetworkFabric(ctx).CreateNetworkFabric(createFabric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(fabricInfo, &fabricPrintConfig)
}

func FabricUpdate(ctx context.Context, fabricId string, fabricName string, description string, config []byte) error {
	logger.Get().Info().Msgf("Update fabric '%s'", fabricId)

	fabricInfo, err := GetFabricById(ctx, fabricId)
	if err != nil {
		return err
	}
	fabricIdNumber, err := utils.GetFloat32FromString(fabricInfo.Id)
	if err != nil {
		return err
	}

	var fabricConfiguration sdk.NetworkFabricFabricConfiguration
	if fabricInfo.FabricConfiguration.EthernetFabric != nil {
		ethernetConfig := sdk.EthernetFabric{}
		err := json.Unmarshal(config, &ethernetConfig)
		if err != nil {
			return err
		}

		fabricConfiguration = sdk.NetworkFabricFabricConfiguration{
			EthernetFabric: &ethernetConfig,
		}
	} else if fabricInfo.FabricConfiguration.FibreChannelFabric != nil {
		fcConfig := sdk.FibreChannelFabric{}
		err := json.Unmarshal(config, &fcConfig)
		if err != nil {
			return err
		}

		fabricConfiguration = sdk.NetworkFabricFabricConfiguration{
			FibreChannelFabric: &fcConfig,
		}
	} else {
		err := fmt.Errorf("invalid fabric type")
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	updateFabric := sdk.UpdateNetworkFabric{
		Name:                sdk.PtrString(fabricName),
		Description:         sdk.PtrString(description),
		SiteId:              fabricInfo.SiteId,
		FabricConfiguration: fabricConfiguration,
	}

	client := api.GetApiClient(ctx)

	fabricInfoUpdated, httpRes, err := client.NetworkFabricAPI.UpdateNetworkFabric(ctx, int32(fabricIdNumber)).
		UpdateNetworkFabric(updateFabric).
		IfMatch(fabricInfo.Revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(fabricInfoUpdated, &fabricPrintConfig)
}

func GetFabricById(ctx context.Context, fabricId string) (*sdk.NetworkFabric, error) {
	client := api.GetApiClient(ctx)

	fabricIdNumber, err := utils.GetFloat32FromString(fabricId)
	if err != nil {
		return nil, err
	}

	fabricInfo, httpRes, err := client.NetworkFabricAPI.GetNetworkFabricById(ctx, fabricIdNumber).Execute()
	if err = response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	return fabricInfo, nil
}
