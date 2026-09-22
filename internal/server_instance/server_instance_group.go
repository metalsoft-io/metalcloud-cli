package server_instance

import (
	"context"
	"fmt"
	"net/http"
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
			Title: "ID",
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
			Title: "ID",
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

var networkEndpointGroupPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Name": {
			MaxWidth: 40,
			Order:    2,
		},
		"SiteId": {
			Title: "Site ID",
			Order: 3,
		},
		"Revision": {
			Title: "Revision",
			Order: 4,
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

var serverInstanceGroupDriveGroupPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
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
		"DriveSizeMbDefault": {
			Title: "Default Size (MB)",
			Order: 4,
		},
		"StorageType": {
			Title: "Storage Type",
			Order: 5,
		},
		"ServiceStatus": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       6,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
	},
}

var serverInstanceGroupInterfacePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Label": {
			MaxWidth: 30,
			Order:    2,
		},
		"GroupId": {
			Title: "Group ID",
			Order: 3,
		},
		"Index": {
			Title: "Index",
			Order: 4,
		},
		"NetworkId": {
			Title: "Network ID",
			Order: 5,
		},
		"ServiceStatus": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       6,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
	},
}

var logicalNetworkACLPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Sequence": {
			Title: "Seq",
			Order: 2,
		},
		"RuleType": {
			Title: "Type",
			Order: 3,
		},
		"Direction": {
			Title: "Direction",
			Order: 4,
		},
		"ForwardingAction": {
			Title: "Action",
			Order: 5,
		},
		"NetworkProtocol": {
			Title: "Protocol",
			Order: 6,
		},
		"SourceAddress": {
			Title:    "Source",
			MaxWidth: 30,
			Order:    7,
		},
		"DestinationAddress": {
			Title:    "Destination",
			MaxWidth: 30,
			Order:    8,
		},
		"SourcePort": {
			Title: "Src Port",
			Order: 9,
		},
		"DestinationPort": {
			Title: "Dst Port",
			Order: 10,
		},
		"EnforcementPoint": {
			Title: "Enforcement",
			Order: 11,
		},
		"LogicalNetworkId": {
			Title: "Logical Network",
			Order: 12,
		},
	},
}

// logicalNetworkACLRecord mirrors sdk.LogicalNetworkACL without the strict
// required-property and enum validation of the generated model. The
// '.../security/rules' collection is typed in the SDK as returning a single
// LogicalNetworkACL, so the listing is parsed from the raw body instead.
type logicalNetworkACLRecord struct {
	Id                 string `json:"id"`
	Sequence           int32  `json:"sequence"`
	RuleType           string `json:"ruleType"`
	Direction          string `json:"direction"`
	ForwardingAction   string `json:"forwardingAction"`
	NetworkProtocol    string `json:"networkProtocol"`
	SourceAddress      string `json:"sourceAddress"`
	DestinationAddress string `json:"destinationAddress"`
	SourcePort         string `json:"sourcePort"`
	DestinationPort    string `json:"destinationPort"`
	SourceMac          string `json:"sourceMac"`
	DestinationMac     string `json:"destinationMac"`
	EnforcementPoint   string `json:"enforcementPoint"`
	EndpointGroupId    string `json:"endpointGroupId"`
	LogicalNetworkId   string `json:"logicalNetworkId"`
}

func ServerInstanceGroupList(ctx context.Context, infrastructureIdOrLabel string) error {
	logger.Get().Info().Msgf("List all server instance groups for infrastructure %s", infrastructureIdOrLabel)

	infra, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.ServerInstanceGroupAPI.GetInfrastructureServerInstanceGroups(ctx, int64(infra.Id)).SortBy([]string{"id:ASC"})

	groups, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(groups, meta, len(groups), &serverInstanceGroupPrintConfig)
}

func ServerInstanceGroupGet(ctx context.Context, serverInstanceGroupId string) error {
	logger.Get().Info().Msgf("Get server instance group details for %s", serverInstanceGroupId)

	serverInstanceGroupIdNumerical, err := GetServerInstanceGroupId(serverInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	serverInstanceGroup, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroup(ctx, serverInstanceGroupIdNumerical).Execute()
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
		DefaultServerTypeId: int64(serverType.Id),
		InstanceCount:       sdk.PtrInt32(int32(serverInstanceCountNumerical)),
	}

	if osTemplateId != "" {
		osTemplate, err := os_template.GetOsTemplateByIdOrLabel(ctx, osTemplateId)
		if err != nil {
			return err
		}

		payload.OsTemplateId = sdk.PtrInt64(int64(osTemplate.Id))
	}

	client := api.GetApiClient(ctx)

	serverInstanceGroupInfo, httpRes, err := client.ServerInstanceGroupAPI.CreateServerInstanceGroup(ctx, int64(infra.Id)).ServerInstanceGroupCreate(payload).Execute()
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
		payload.OsTemplateId = sdk.PtrInt64(int64(osTemplateId))
	}

	client := api.GetApiClient(ctx)

	serverInstanceGroupConfig, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroupConfig(ctx, serverInstanceGroupIdNumerical).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	serverInstanceGroupConfig, httpRes, err = client.ServerInstanceGroupAPI.UpdateServerInstanceGroupConfig(ctx, serverInstanceGroupIdNumerical).
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

	serverInstanceGroup, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroup(ctx, serverInstanceGroupIdNumerical).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	httpRes, err = client.ServerInstanceGroupAPI.DeleteServerInstanceGroup(ctx, serverInstanceGroupIdNumerical).
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

	serverInstancesList, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroupServerInstances(ctx, serverInstanceGroupIdNumerical).Execute()
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

	connections, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroupNetworkConfigurationConnections(ctx, serverInstanceGroupIdNumerical).Execute()
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

	networkConnectionIdNumerical, err := utils.GetInt64FromString(networkConnectionId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	connection, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroupNetworkConfigurationConnectionById(ctx, serverInstanceGroupIdNumerical, networkConnectionIdNumerical).Execute()
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
		CreateServerInstanceGroupNetworkConfigurationConnection(ctx, serverInstanceGroupIdNumerical).
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
		UpdateServerInstanceGroupNetworkConfigurationConnection(ctx, serverInstanceGroupIdNumerical, networkConnectionIdNumerical).
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

	networkConnectionIdNumerical, err := utils.GetInt64FromString(networkConnectionId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerInstanceGroupAPI.DeleteServerInstanceGroupNetworkConfigurationConnection(ctx, serverInstanceGroupIdNumerical, networkConnectionIdNumerical).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network connection %s successfully deleted", networkConnectionId)
	return nil
}

func ServerInstanceGroupMetaUpdate(ctx context.Context, serverInstanceGroupId string, config []byte) error {
	logger.Get().Info().Msgf("Updating metadata of server instance group '%s'", serverInstanceGroupId)

	var meta sdk.GenericMeta
	if err := utils.UnmarshalContent(config, &meta); err != nil {
		return err
	}

	groupId, err := GetServerInstanceGroupId(serverInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerInstanceGroupAPI.
		UpdateServerInstanceGroupMeta(ctx, groupId).
		GenericMeta(meta).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Server instance group '%s' metadata updated", serverInstanceGroupId)
	return nil
}

func ServerInstanceGroupDriveGroups(ctx context.Context, serverInstanceGroupId string) error {
	logger.Get().Info().Msgf("Listing drive groups of server instance group '%s'", serverInstanceGroupId)

	groupId, err := GetServerInstanceGroupId(serverInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	driveGroups, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroupDriveGroups(ctx, groupId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(driveGroups, &serverInstanceGroupDriveGroupPrintConfig)
}

func ServerInstanceGroupInterfaces(ctx context.Context, serverInstanceGroupId string, filters ServerInstanceFilters) error {
	logger.Get().Info().Msgf("Listing interfaces of server instance group '%s'", serverInstanceGroupId)

	groupId, err := GetServerInstanceGroupId(serverInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.ServerInstanceGroupAPI.GetServerInstanceGroupInterfaces(ctx, groupId).SortBy([]string{"id:ASC"})
	if len(filters.InfrastructureId) > 0 {
		request = request.FilterInfrastructureId(utils.ProcessFilterStringSlice(filters.InfrastructureId))
	}
	if len(filters.ServiceStatus) > 0 {
		request = request.FilterServiceStatus(utils.ProcessFilterStringSlice(filters.ServiceStatus))
	}
	if len(filters.ConfigDeployStatus) > 0 {
		request = request.FilterConfigDeployStatus(utils.ProcessFilterStringSlice(filters.ConfigDeployStatus))
	}
	if len(filters.ConfigDeployType) > 0 {
		request = request.FilterConfigDeployType(utils.ProcessFilterStringSlice(filters.ConfigDeployType))
	}

	interfaces, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(interfaces, meta, len(interfaces), &serverInstanceGroupInterfacePrintConfig)
}

func ServerInstanceGroupInterfaceGet(ctx context.Context, serverInstanceGroupId string, interfaceId string) error {
	logger.Get().Info().Msgf("Get interface '%s' of server instance group '%s'", interfaceId, serverInstanceGroupId)

	groupId, err := GetServerInstanceGroupId(serverInstanceGroupId)
	if err != nil {
		return err
	}

	interfaceIdNumerical, err := utils.GetInt64FromString(interfaceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	groupInterface, httpRes, err := client.ServerInstanceGroupAPI.
		GetServerInstanceGroupInterface(ctx, groupId, interfaceIdNumerical).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(groupInterface, &serverInstanceGroupInterfacePrintConfig)
}

// ServerInstanceGroupNetworkConfiguration returns the network endpoint group
// that holds the network configuration of a server instance group
// (GET /config/networking). The per-connection listing lives in
// ServerInstanceGroupNetworkList.
func ServerInstanceGroupNetworkConfiguration(ctx context.Context, serverInstanceGroupId string) error {
	logger.Get().Info().Msgf("Get network configuration of server instance group '%s'", serverInstanceGroupId)

	groupId, err := GetServerInstanceGroupId(serverInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkConfiguration, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroupNetworkConfiguration(ctx, groupId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkConfiguration, &networkEndpointGroupPrintConfig)
}

// ServerInstanceGroupNetworkReplace creates or replaces the network
// configuration of a server instance group (PUT /config/networking).
//
// The generated SDK request exposes no body setter for this operation, so a
// caller supplied configuration is sent through a raw request carrying the
// If-Match of the group configuration; without a configuration the typed SDK
// call is used, which sends no body at all.
func ServerInstanceGroupNetworkReplace(ctx context.Context, serverInstanceGroupId string, config []byte) error {
	logger.Get().Info().Msgf("Replacing network configuration of server instance group '%s'", serverInstanceGroupId)

	if len(config) == 0 {
		groupId, err := GetServerInstanceGroupId(serverInstanceGroupId)
		if err != nil {
			return err
		}

		client := api.GetApiClient(ctx)

		networkConfiguration, httpRes, err := client.ServerInstanceGroupAPI.UpdateServerInstanceGroupNetworkConfiguration(ctx, groupId).Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			return err
		}

		return formatter.PrintResult(networkConfiguration, &networkEndpointGroupPrintConfig)
	}

	currentConfig, groupId, err := getServerInstanceGroupConfig(ctx, serverInstanceGroupId)
	if err != nil {
		return err
	}

	body, err := api.RawJSONRequest(ctx, http.MethodPut,
		fmt.Sprintf("/api/v2/server-instance-groups/%d/config/networking", groupId),
		config,
		api.IfMatchHeader(strconv.FormatInt(currentConfig.Revision, 10)))
	if err != nil {
		return err
	}

	if len(body) == 0 {
		logger.Get().Info().Msgf("Network configuration of server instance group '%s' replaced", serverInstanceGroupId)
		return nil
	}

	return utils.PrintRawObject(body, &networkEndpointGroupPrintConfig)
}

// ServerInstanceGroupACLList lists the security rules of a network connection.
// The SDK types this collection endpoint as returning a single
// LogicalNetworkACL, so the body is parsed raw and rendered as a list; a
// single object response is still rendered correctly.
func ServerInstanceGroupACLList(ctx context.Context, serverInstanceGroupId string, connectionId string) error {
	logger.Get().Info().Msgf("Listing security rules of network connection '%s' of server instance group '%s'", connectionId, serverInstanceGroupId)

	groupId, connectionIdNumerical, err := getGroupAndConnectionId(serverInstanceGroupId, connectionId)
	if err != nil {
		return err
	}

	body, err := api.RawJSONRequest(ctx, http.MethodGet, aclCollectionPath(groupId, connectionIdNumerical), nil, nil)
	if err != nil {
		return err
	}

	// A bare object body (no "data" envelope) is a single rule, not a list.
	if object, objectErr := utils.DecodeRawObject(body); objectErr == nil {
		if _, hasData := object["data"]; !hasData {
			return formatter.PrintResult(object, &logicalNetworkACLPrintConfig)
		}
	}

	rawItems, err := utils.DecodeRawItems(body)
	if err != nil {
		return err
	}

	records, err := utils.UnmarshalRawItems[logicalNetworkACLRecord](rawItems)
	if err != nil {
		return fmt.Errorf("failed to parse security rules: %w", err)
	}

	return utils.PrintAllRaw(rawItems, records, sdk.PaginatedResponseMeta{}, len(records), &logicalNetworkACLPrintConfig)
}

func ServerInstanceGroupACLGet(ctx context.Context, serverInstanceGroupId string, connectionId string, ruleId string) error {
	logger.Get().Info().Msgf("Get security rule '%s' of network connection '%s'", ruleId, connectionId)

	groupId, connectionIdNumerical, err := getGroupAndConnectionId(serverInstanceGroupId, connectionId)
	if err != nil {
		return err
	}

	ruleIdNumerical, err := utils.GetInt64FromString(ruleId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	rule, httpRes, err := client.ServerInstanceGroupAPI.
		GetServerInstanceGroupLogicalNetworkACLById(ctx, groupId, connectionIdNumerical, ruleIdNumerical).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(rule, &logicalNetworkACLPrintConfig)
}

func ServerInstanceGroupACLConfigExample(ctx context.Context) error {
	example := sdk.CreateLogicalNetworkACL{
		RuleType:           sdk.ACLTYPE_IPV4,
		Direction:          sdk.ACLDIRECTION_IN,
		Sequence:           10,
		ForwardingAction:   sdk.ACLFORWARDINGACTION_ALLOW,
		NetworkProtocol:    sdk.PtrString("tcp"),
		SourceAddress:      sdk.PtrString("10.0.0.0/24"),
		DestinationAddress: sdk.PtrString("10.0.1.0/24"),
		SourcePort:         sdk.PtrString("any"),
		DestinationPort:    sdk.PtrString("443"),
		EnforcementPoint:   sdk.ACLENFORCEMENTPOINT_SVI,
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(example, nil)
	}
	return formatter.PrintYamlResult(example)
}

func ServerInstanceGroupACLAdd(ctx context.Context, serverInstanceGroupId string, connectionId string, create sdk.CreateLogicalNetworkACL) error {
	logger.Get().Info().Msgf("Adding security rule to network connection '%s' of server instance group '%s'", connectionId, serverInstanceGroupId)

	groupId, connectionIdNumerical, err := getGroupAndConnectionId(serverInstanceGroupId, connectionId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	rule, httpRes, err := client.ServerInstanceGroupAPI.
		CreateServerInstanceGroupLogicalNetworkACL(ctx, groupId, connectionIdNumerical).
		CreateLogicalNetworkACL(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(rule, &logicalNetworkACLPrintConfig)
}

func ServerInstanceGroupACLUpdate(ctx context.Context, serverInstanceGroupId string, connectionId string, ruleId string, config []byte) error {
	logger.Get().Info().Msgf("Updating security rule '%s' of network connection '%s'", ruleId, connectionId)

	var update sdk.UpdateLogicalNetworkACL
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	groupId, connectionIdNumerical, err := getGroupAndConnectionId(serverInstanceGroupId, connectionId)
	if err != nil {
		return err
	}

	ruleIdNumerical, err := utils.GetInt64FromString(ruleId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	rule, httpRes, err := client.ServerInstanceGroupAPI.
		UpdateServerInstanceGroupLogicalNetworkACL(ctx, groupId, connectionIdNumerical, ruleIdNumerical).
		UpdateLogicalNetworkACL(update).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(rule, &logicalNetworkACLPrintConfig)
}

func ServerInstanceGroupACLRemove(ctx context.Context, serverInstanceGroupId string, connectionId string, ruleId string) error {
	logger.Get().Info().Msgf("Removing security rule '%s' from network connection '%s'", ruleId, connectionId)

	groupId, connectionIdNumerical, err := getGroupAndConnectionId(serverInstanceGroupId, connectionId)
	if err != nil {
		return err
	}

	ruleIdNumerical, err := utils.GetInt64FromString(ruleId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerInstanceGroupAPI.
		DeleteServerInstanceGroupLogicalNetworkACL(ctx, groupId, connectionIdNumerical, ruleIdNumerical).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Security rule '%s' removed", ruleId)
	return nil
}

func aclCollectionPath(groupId int64, connectionId int64) string {
	return fmt.Sprintf("/api/v2/server-instance-groups/%d/config/networking/connections/%d/security/rules", groupId, connectionId)
}

func getGroupAndConnectionId(serverInstanceGroupId string, connectionId string) (int64, int64, error) {
	groupId, err := GetServerInstanceGroupId(serverInstanceGroupId)
	if err != nil {
		return 0, 0, err
	}

	connectionIdNumerical, err := utils.GetInt64FromString(connectionId)
	if err != nil {
		return 0, 0, err
	}

	return groupId, connectionIdNumerical, nil
}

// getServerInstanceGroupConfig returns the group configuration together with
// the numeric group id. The configuration carries its own revision, which
// guards every write below the /config path.
func getServerInstanceGroupConfig(ctx context.Context, serverInstanceGroupId string) (*sdk.ServerInstanceGroupConfiguration, int64, error) {
	groupId, err := GetServerInstanceGroupId(serverInstanceGroupId)
	if err != nil {
		return nil, 0, err
	}

	client := api.GetApiClient(ctx)

	config, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroupConfig(ctx, groupId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, 0, err
	}

	return config, groupId, nil
}

func GetServerInstanceGroupId(serverInstanceGroupId string) (int64, error) {
	serverInstanceGroupIdNumeric, err := strconv.ParseInt(serverInstanceGroupId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid server instance group ID: '%s'", serverInstanceGroupId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return serverInstanceGroupIdNumeric, nil
}
