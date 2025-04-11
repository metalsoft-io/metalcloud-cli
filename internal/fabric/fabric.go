package fabric

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/metalsoft-io/metalcloud-cli/internal/network_device"
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
		"SiteId": {
			Title: "Site",
			Order: 4,
		},
		"Status": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       5,
		},
		"FabricConfiguration": {
			Hidden: true,
			InnerFields: map[string]formatter.RecordFieldConfig{
				"EthernetFabric": {
					Hidden: true,
					InnerFields: map[string]formatter.RecordFieldConfig{
						"FabricType": {
							Title: "Type",
							Order: 6,
						},
						"DefaultVlan": {
							Title: "Default VLAN",
							Order: 7,
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

	fabricInfo, err := GetFabricByIdOrLabel(ctx, fabricId)
	if err != nil {
		return err
	}

	return formatter.PrintResult(fabricInfo, &fabricPrintConfig)
}

func FabricConfigExample(ctx context.Context, fabricType string) error {
	var fabricConfiguration sdk.NetworkFabricFabricConfiguration
	switch fabricType {
	case "ethernet":
		ethernetConfig := sdk.EthernetFabric{
			FabricType:                       sdk.FABRICTYPE_ETHERNET,
			DefaultNetworkProfileId:          sdk.PtrInt32(101),
			GnmiMonitoringEnabled:            sdk.PtrBool(false),
			SyslogMonitoringEnabled:          sdk.PtrBool(true),
			ZeroTouchEnabled:                 sdk.PtrBool(false),
			AsnRanges:                        []string{"65000-65010"},
			DefaultVlan:                      sdk.PtrInt32(10),
			ExtraInternalIPsPerSubnet:        sdk.PtrInt32(2),
			LagRanges:                        []string{"100-200", "300-400"},
			LeafSwitchesHaveMlagPairs:        sdk.PtrBool(false),
			MlagRanges:                       []string{"30-40", "50-60"},
			NumberOfSpinesNextToLeafSwitches: sdk.PtrInt32(2),
			PreventVlanCleanup:               []string{"1000-1100"},
			PreventCleanupFromUplinks:        sdk.PtrBool(true),
			ReservedVlans:                    []string{"2000-2100", "2200-2300"},
			VlanRanges:                       []string{"3000-3100", "2000-2100"},
			VniPrefix:                        sdk.PtrInt32(5000),
			VrfVlanRanges:                    []string{"400-450", "460-470"},
		}

		fabricConfiguration = sdk.NetworkFabricFabricConfiguration{
			EthernetFabric: &ethernetConfig,
		}
	case "fibre_channel":
		fcConfig := sdk.FibreChannelFabric{
			FabricType:               sdk.FABRICTYPE_FIBRE_CHANNEL,
			DefaultNetworkProfileId:  sdk.PtrInt32(101),
			GnmiMonitoringEnabled:    sdk.PtrBool(false),
			SyslogMonitoringEnabled:  sdk.PtrBool(true),
			ZeroTouchEnabled:         sdk.PtrBool(false),
			VsanId:                   sdk.PtrInt32(1),
			TopologyType:             sdk.FABRICTOPOLOGYTYPE_MESH,
			Mtu:                      sdk.PtrFloat32(1200),
			ZoningConfiguration:      map[string]interface{}{"zone1": []string{"wwn1", "wwn2"}},
			InteropMode:              sdk.PtrString("full"),
			QosConfiguration:         map[string]interface{}{"qos1": "low"},
			TrunkingConfiguration:    map[string]interface{}{"trunk1": []string{"wwn1", "wwn2"}},
			PortChannelConfiguration: map[string]interface{}{"port1": []string{"wwn1", "wwn2"}},
		}

		fabricConfiguration = sdk.NetworkFabricFabricConfiguration{
			FibreChannelFabric: &fcConfig,
		}
	default:
		err := fmt.Errorf("invalid fabric type: '%s'", fabricType)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	return formatter.PrintResult(fabricConfiguration, nil)
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

	fabricInfo, err := GetFabricByIdOrLabel(ctx, fabricId)
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

func FabricActivate(ctx context.Context, fabricId string) error {
	logger.Get().Info().Msgf("Activate fabric '%s'", fabricId)

	fabricIdNumeric, err := utils.GetFloat32FromString(fabricId)
	if err != nil {
		return err
	}

	fabricInfo, err := GetFabricByIdOrLabel(ctx, fabricId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	fabricInfo, httpRes, err := client.NetworkFabricAPI.
		ActivateNetworkFabric(ctx, int32(fabricIdNumeric)).
		IfMatch(fabricInfo.Revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(fabricInfo, &fabricPrintConfig)
}

func FabricDevicesGet(ctx context.Context, fabricId string) error {
	logger.Get().Info().Msgf("Get fabric '%s' devices", fabricId)

	fabricIdNumeric, err := utils.GetFloat32FromString(fabricId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	devicesList, httpRes, err := client.NetworkFabricAPI.GetFabricNetworkDevices(ctx, int32(fabricIdNumeric)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(devicesList.Data, &network_device.NetworkDevicePrintConfig)
}

func FabricDevicesAdd(ctx context.Context, fabricId string, deviceIds []string) error {
	logger.Get().Info().Msgf("Adding devices '%v' to fabric '%s'", deviceIds, fabricId)

	fabricInfo, err := GetFabricByIdOrLabel(ctx, fabricId)
	if err != nil {
		return err
	}

	fabricIdNumeric, err := utils.GetFloat32FromString(fabricInfo.Id)
	if err != nil {
		return err
	}

	deviceIdsNumeric := make([]float32, 0)
	for _, deviceId := range deviceIds {
		device, err := network_device.GetNetworkDeviceById(ctx, deviceId)
		if err != nil {
			return err
		}

		if *fabricInfo.SiteId != device.SiteId {
			err := fmt.Errorf("device '%s' is not in the same site as fabric '%s'", deviceId, fabricId)
			logger.Get().Error().Err(err).Msg("")
			return err
		}

		deviceIdNumeric, err := utils.GetFloat32FromString(device.Id)
		if err != nil {
			return err
		}

		deviceIdsNumeric = append(deviceIdsNumeric, deviceIdNumeric)
	}

	client := api.GetApiClient(ctx)

	_, httpRes, err := client.NetworkFabricAPI.AddNetworkDevicesToFabric(ctx, int32(fabricIdNumeric)).
		NetworkDevicesToFabric(sdk.NetworkDevicesToFabric{NetworkDeviceIds: deviceIdsNumeric}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return nil
}

func FabricDevicesRemove(ctx context.Context, fabricId string, deviceId string) error {
	logger.Get().Info().Msgf("Removing device '%s' from fabric '%s'", deviceId, fabricId)

	fabricInfo, err := GetFabricByIdOrLabel(ctx, fabricId)
	if err != nil {
		return err
	}

	fabricIdNumeric, err := utils.GetFloat32FromString(fabricInfo.Id)
	if err != nil {
		return err
	}

	device, err := network_device.GetNetworkDeviceById(ctx, deviceId)
	if err != nil {
		return err
	}

	deviceIdNumeric, err := utils.GetFloat32FromString(device.Id)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, httpRes, err := client.NetworkFabricAPI.RemoveNetworkDeviceFromFabric(ctx, int32(fabricIdNumeric), int32(deviceIdNumeric)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return nil
}

func GetFabricByIdOrLabel(ctx context.Context, fabricIdOrLabel string) (*sdk.NetworkFabric, error) {
	client := api.GetApiClient(ctx)

	fabricIdNumber, err := utils.GetFloat32FromString(fabricIdOrLabel)
	if err == nil {
		fabricInfo, httpRes, err := client.NetworkFabricAPI.GetNetworkFabricById(ctx, fabricIdNumber).Execute()
		if err = response_inspector.InspectResponse(httpRes, err); err == nil {
			return fabricInfo, nil
		}
	}

	fabrics, httpRes, err := client.NetworkFabricAPI.
		GetNetworkFabrics(ctx).
		FilterName([]string{fabricIdOrLabel}).
		Execute()
	if err = response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	if len(fabrics.Data) == 0 {
		err := fmt.Errorf("fabric '%s' not found", fabricIdOrLabel)
		logger.Get().Error().Err(err).Msg("")
		return nil, err
	}

	return &fabrics.Data[0], nil
}
