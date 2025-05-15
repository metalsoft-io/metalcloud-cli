package server_instance

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/internal/infrastructure"
	"github.com/metalsoft-io/metalcloud-cli/internal/os_template"
	"github.com/metalsoft-io/metalcloud-cli/internal/server_type"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var serverInstanceGroupPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Label": {
			MaxWidth: 30,
			Order:    2,
		},
		"InfrastructureId": {
			Title: "Infra ID",
			Order: 3,
		},
		"ServiceStatus": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       4,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       5,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
	},
}

var serverInstanceGroupConfigPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Label": {
			MaxWidth: 30,
			Order:    1,
		},
		"InstanceCount": {
			Title: "Count",
			Order: 2,
		},
		"VolumeTemplateId": {
			Title: "OS Template Id",
			Order: 3,
		},
	},
}

var networkConnectionPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"NetworkId": {
			Title: "Network ID",
			Order: 2,
		},
		"Network": {
			Title:    "Network",
			Order:    3,
			MaxWidth: 30,
		},
		"SubnetId": {
			Title: "Subnet ID",
			Order: 4,
		},
		"SubnetGateway": {
			Title: "Gateway",
			Order: 5,
		},
		"SubnetPrefixSize": {
			Title: "Prefix Size",
			Order: 6,
		},
	},
}

func ServerInstanceGroupList(ctx context.Context, infrastructureIdOrLabel string) error {
	logger.Get().Info().Msgf("List all server instance groups for infrastructure %s", infrastructureIdOrLabel)

	infra, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	serverInstanceGroupList, httpRes, err := client.ServerInstanceGroupAPI.GetInfrastructureServerInstanceGroups(ctx, int32(infra.Id)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(serverInstanceGroupList, &serverInstanceGroupPrintConfig)
}

func ServerInstanceGroupGet(ctx context.Context, serverInstanceGroupId string) error {
	logger.Get().Info().Msgf("Get server instance group details for %s", serverInstanceGroupId)

	serverInstanceGroupIdNumerical, err := GetServerInstanceGroupId(serverInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	serverInstanceGroup, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroup(ctx, int32(serverInstanceGroupIdNumerical)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(serverInstanceGroup, &serverInstanceGroupPrintConfig)
}

func ServerInstanceGroupCreate(ctx context.Context, infrastructureIdOrLabel string, label string, serverTypeId string, instanceCount string, osTemplateId string) error {
	logger.Get().Info().Msgf("Create new server instance group in infrastructure %s", infrastructureIdOrLabel)

	serverInstanceCountNumerical, err := utils.GetFloat32FromString(instanceCount)
	if err != nil {
		return err
	}

	infra, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	serverType, err := server_type.GetServerTypeByIdOrLabel(ctx, serverTypeId)
	if err != nil {
		return err
	}

	payload := sdk.ServerInstanceGroupCreate{
		Label:               &label,
		DefaultServerTypeId: int32(serverType.Id),
		InstanceCount:       sdk.PtrInt32(int32(serverInstanceCountNumerical)),
	}

	if osTemplateId != "" {
		osTemplate, err := os_template.GetOsTemplateByIdOrLabel(ctx, osTemplateId)
		if err != nil {
			return err
		}

		payload.OsTemplateId = sdk.PtrInt32(int32(osTemplate.Id))
	}

	client := api.GetApiClient(ctx)

	serverInstanceGroupInfo, httpRes, err := client.ServerInstanceGroupAPI.CreateServerInstanceGroup(ctx, int32(infra.Id)).ServerInstanceGroupCreate(payload).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(serverInstanceGroupInfo, &serverInstanceGroupPrintConfig)
}

func ServerInstanceGroupUpdate(ctx context.Context, serverInstanceGroupId string, label string, instanceCount int, osTemplateId int) error {
	logger.Get().Info().Msgf("Update server instance group %s", serverInstanceGroupId)

	serverInstanceGroupIdNumerical, err := GetServerInstanceGroupId(serverInstanceGroupId)
	if err != nil {
		return err
	}

	payload := sdk.ServerInstanceGroupUpdate{}

	if label != "" {
		payload.Label = &label
	}

	if instanceCount > 0 {
		payload.InstanceCount = sdk.PtrInt32(int32(instanceCount))
	}

	if osTemplateId > 0 {
		payload.OsTemplateId = sdk.PtrInt32(int32(osTemplateId))
	}

	client := api.GetApiClient(ctx)

	serverInstanceGroupConfig, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroupConfig(ctx, int32(serverInstanceGroupIdNumerical)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	serverInstanceGroupConfig, httpRes, err = client.ServerInstanceGroupAPI.UpdateServerInstanceGroupConfig(ctx, int32(serverInstanceGroupIdNumerical)).
		IfMatch(strconv.Itoa(int(serverInstanceGroupConfig.Revision))).
		ServerInstanceGroupUpdate(payload).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(serverInstanceGroupConfig, &serverInstanceGroupConfigPrintConfig)
}

func ServerInstanceGroupDelete(ctx context.Context, serverInstanceGroupId string) error {
	logger.Get().Info().Msgf("Delete server instance group %s", serverInstanceGroupId)

	serverInstanceGroupIdNumerical, err := GetServerInstanceGroupId(serverInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	serverInstanceGroup, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroup(ctx, int32(serverInstanceGroupIdNumerical)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	httpRes, err = client.ServerInstanceGroupAPI.DeleteServerInstanceGroup(ctx, int32(serverInstanceGroupIdNumerical)).
		IfMatch(strconv.Itoa(int(serverInstanceGroup.Revision))).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return nil
}

func ServerInstanceGroupInstances(ctx context.Context, serverInstanceGroupId string) error {
	logger.Get().Info().Msgf("List instances of server instance group %s", serverInstanceGroupId)

	serverInstanceGroupIdNumerical, err := GetServerInstanceGroupId(serverInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	serverInstancesList, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroupServerInstances(ctx, int32(serverInstanceGroupIdNumerical)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(serverInstancesList, &serverInstancePrintConfig)
}

func ServerInstanceGroupNetworkList(ctx context.Context, serverInstanceGroupId string) error {
	logger.Get().Info().Msgf("List network connections for server instance group %s", serverInstanceGroupId)

	serverInstanceGroupIdNumerical, err := GetServerInstanceGroupId(serverInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	connections, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroupNetworkConfigurationConnections(ctx, int32(serverInstanceGroupIdNumerical)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(connections, &networkConnectionPrintConfig)
}

func ServerInstanceGroupNetworkGet(ctx context.Context, serverInstanceGroupId string, networkConnectionId string) error {
	logger.Get().Info().Msgf("Get network connection %s details for server instance group %s", networkConnectionId, serverInstanceGroupId)

	serverInstanceGroupIdNumerical, err := GetServerInstanceGroupId(serverInstanceGroupId)
	if err != nil {
		return err
	}

	networkConnectionIdNumerical, err := utils.GetFloat32FromString(networkConnectionId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	connection, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroupNetworkConfigurationConnectionById(ctx, int32(serverInstanceGroupIdNumerical), int32(networkConnectionIdNumerical)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(connection, &networkConnectionPrintConfig)
}

func ServerInstanceGroupNetworkConnect(ctx context.Context, serverInstanceGroupId string, networkId string, accessMode string, isTagged string, redundancy string) error {
	logger.Get().Info().Msgf("Create network connection for server instance group %s", serverInstanceGroupId)

	serverInstanceGroupIdNumerical, err := GetServerInstanceGroupId(serverInstanceGroupId)
	if err != nil {
		return err
	}

	tagged, err := strconv.ParseBool(isTagged)
	if err != nil {
		return fmt.Errorf("invalid tagged value: %s", isTagged)
	}

	payload := sdk.CreateServerInstanceGroupNetworkConnection{
		LogicalNetworkId: networkId,
		AccessMode:       sdk.NetworkEndpointGroupAllowedAccessMode(accessMode),
		Tagged:           tagged,
	}

	if redundancy != "" {
		payload.Redundancy = *sdk.NewNullableRedundancyConfig(
			&sdk.RedundancyConfig{
				Mode: sdk.NetworkEndpointGroupRedundancyMode(redundancy),
			},
		)
	}

	client := api.GetApiClient(ctx)

	connection, httpRes, err := client.ServerInstanceGroupAPI.
		CreateServerInstanceGroupNetworkConfigurationConnection(ctx, int32(serverInstanceGroupIdNumerical)).
		CreateServerInstanceGroupNetworkConnection(payload).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(connection, &networkConnectionPrintConfig)
}

func ServerInstanceGroupNetworkUpdate(ctx context.Context, serverInstanceGroupId string, networkConnectionId string, accessMode string, isTagged string, redundancy string) error {
	logger.Get().Info().Msgf("Update network connection %s for server instance group %s", networkConnectionId, serverInstanceGroupId)

	serverInstanceGroupIdNumerical, err := GetServerInstanceGroupId(serverInstanceGroupId)
	if err != nil {
		return err
	}

	networkConnectionIdNumerical, err := utils.GetFloat32FromString(networkConnectionId)
	if err != nil {
		return err
	}

	payload := sdk.UpdateNetworkEndpointGroupLogicalNetwork{}

	if accessMode != "" {
		accessModeValue := sdk.NetworkEndpointGroupAllowedAccessMode(accessMode)
		payload.AccessMode = &accessModeValue
	}

	if isTagged != "" {
		tagged, err := strconv.ParseBool(isTagged)
		if err != nil {
			return fmt.Errorf("invalid tagged value: %s", isTagged)
		}

		payload.Tagged = &tagged
	}

	if redundancy != "" {
		payload.Redundancy = *sdk.NewNullableRedundancyConfig(
			&sdk.RedundancyConfig{
				Mode: sdk.NetworkEndpointGroupRedundancyMode(redundancy),
			},
		)
	}

	client := api.GetApiClient(ctx)

	connection, httpRes, err := client.ServerInstanceGroupAPI.
		UpdateServerInstanceGroupNetworkConfigurationConnection(ctx, int32(serverInstanceGroupIdNumerical), networkConnectionIdNumerical).
		UpdateNetworkEndpointGroupLogicalNetwork(payload).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(connection, &networkConnectionPrintConfig)
}

func ServerInstanceGroupNetworkDisconnect(ctx context.Context, serverInstanceGroupId string, networkConnectionId string) error {
	logger.Get().Info().Msgf("Delete network connection %s from server instance group %s", networkConnectionId, serverInstanceGroupId)

	serverInstanceGroupIdNumerical, err := GetServerInstanceGroupId(serverInstanceGroupId)
	if err != nil {
		return err
	}

	networkConnectionIdNumerical, err := utils.GetFloat32FromString(networkConnectionId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerInstanceGroupAPI.DeleteServerInstanceGroupNetworkConfigurationConnection(ctx, int32(serverInstanceGroupIdNumerical), int32(networkConnectionIdNumerical)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network connection %s successfully deleted", networkConnectionId)
	return nil
}

func GetServerInstanceGroupId(serverInstanceGroupId string) (float32, error) {
	serverInstanceGroupIdNumeric, err := strconv.ParseFloat(serverInstanceGroupId, 32)
	if err != nil {
		err := fmt.Errorf("invalid server instance group ID: '%s'", serverInstanceGroupId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(serverInstanceGroupIdNumeric), nil
}
