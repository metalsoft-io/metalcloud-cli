package endpoint_instance

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/internal/infrastructure"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// EndpointInstanceGroupFilters carries the query filters of the
// infrastructure scoped endpoint instance group listing.
type EndpointInstanceGroupFilters struct {
	ExtensionInstanceId []string
	ServiceStatus       []string
	ConfigDeployStatus  []string
	ConfigDeployType    []string
}

var endpointInstanceGroupPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Label": {
			Order:    2,
			MaxWidth: 30,
		},
		"InfrastructureId": {
			Title: "Infra ID",
			Order: 3,
		},
		"EndpointGroupName": {
			Title:    "Endpoint Group",
			Order:    4,
			MaxWidth: 30,
		},
		"ExtensionInstanceId": {
			Title: "Extension Instance",
			Order: 5,
		},
		"ResourcePoolId": {
			Title: "Resource Pool",
			Order: 6,
		},
		"ServiceStatus": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       7,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       8,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       9,
		},
	},
}

var endpointInstanceGroupConfigPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Label": {
			Order:    1,
			MaxWidth: 30,
		},
		"EndpointGroupName": {
			Title:    "Endpoint Group",
			Order:    2,
			MaxWidth: 30,
		},
		"Hostname": {
			Title:    "Hostname",
			Order:    3,
			MaxWidth: 30,
		},
		"DeployType": {
			Title: "Deploy Type",
			Order: 4,
		},
		"DeployStatus": {
			Title:       "Deploy Status",
			Transformer: formatter.FormatStatusValue,
			Order:       5,
		},
		"Revision": {
			Title: "Revision",
			Order: 6,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
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
			Order:    2,
			MaxWidth: 40,
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

var networkConnectionPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"AccessMode": {
			Title: "Access Mode",
			Order: 2,
		},
		"Tagged": {
			Title: "Tagged",
			Order: 3,
		},
		"Mtu": {
			Title: "MTU",
			Order: 4,
		},
		"InterfaceCount": {
			Title: "Interfaces",
			Order: 5,
		},
		"ProvidesDefaultRoute": {
			Title: "Default Route",
			Order: 6,
		},
		"DisableAutoIpAllocation": {
			Title: "No Auto IP",
			Order: 7,
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
			Order:    7,
			MaxWidth: 30,
		},
		"DestinationAddress": {
			Title:    "Destination",
			Order:    8,
			MaxWidth: 30,
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

// GroupPrintConfig exposes the endpoint instance group table layout.
func GroupPrintConfig() *formatter.PrintConfig {
	return &endpointInstanceGroupPrintConfig
}

func EndpointInstanceGroupList(ctx context.Context, infrastructureIdOrLabel string, filters EndpointInstanceGroupFilters) error {
	logger.Get().Info().Msgf("Listing endpoint instance groups of infrastructure '%s'", infrastructureIdOrLabel)

	infra, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.EndpointInstanceGroupAPI.GetInfrastructureEndpointInstanceGroups(ctx, int64(infra.Id)).SortBy([]string{"id:ASC"})
	if len(filters.ExtensionInstanceId) > 0 {
		request = request.FilterExtensionInstanceId(utils.ProcessFilterStringSlice(filters.ExtensionInstanceId))
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

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &endpointInstanceGroupPrintConfig)
}

func EndpointInstanceGroupGet(ctx context.Context, endpointInstanceGroupId string) error {
	logger.Get().Info().Msgf("Get endpoint instance group '%s'", endpointInstanceGroupId)

	group, err := getEndpointInstanceGroup(ctx, endpointInstanceGroupId)
	if err != nil {
		return err
	}

	return formatter.PrintResult(group, &endpointInstanceGroupPrintConfig)
}

func EndpointInstanceGroupConfigExample(ctx context.Context) error {
	example := sdk.EndpointInstanceGroupCreate{
		Label:               sdk.PtrString("my-endpoint-instance-group"),
		EndpointGroupName:   sdk.PtrString("my-endpoint-group"),
		ExtensionInstanceId: sdk.PtrInt64(1),
		Hostname:            sdk.PtrString("my-endpoint-group"),
		ResourcePoolId:      sdk.PtrInt64(1),
		Tags:                []string{"example"},
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(example, nil)
	}
	return formatter.PrintYamlResult(example)
}

func EndpointInstanceGroupCreate(ctx context.Context, infrastructureIdOrLabel string, create sdk.EndpointInstanceGroupCreate) error {
	logger.Get().Info().Msgf("Creating endpoint instance group in infrastructure '%s'", infrastructureIdOrLabel)

	infra, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	group, httpRes, err := client.EndpointInstanceGroupAPI.
		CreateEndpointInstanceGroup(ctx, int64(infra.Id)).
		EndpointInstanceGroupCreate(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(group, &endpointInstanceGroupPrintConfig)
}

func EndpointInstanceGroupDelete(ctx context.Context, endpointInstanceGroupId string) error {
	logger.Get().Info().Msgf("Deleting endpoint instance group '%s'", endpointInstanceGroupId)

	group, err := getEndpointInstanceGroup(ctx, endpointInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.EndpointInstanceGroupAPI.
		DeleteEndpointInstanceGroup(ctx, group.Id).
		IfMatch(strconv.FormatInt(group.Revision, 10)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Endpoint instance group '%s' deleted", endpointInstanceGroupId)
	return nil
}

func EndpointInstanceGroupConfigGet(ctx context.Context, endpointInstanceGroupId string) error {
	logger.Get().Info().Msgf("Get endpoint instance group '%s' configuration", endpointInstanceGroupId)

	config, _, err := getEndpointInstanceGroupConfig(ctx, endpointInstanceGroupId)
	if err != nil {
		return err
	}

	return formatter.PrintResult(config, &endpointInstanceGroupConfigPrintConfig)
}

func EndpointInstanceGroupConfigUpdate(ctx context.Context, endpointInstanceGroupId string, config []byte) error {
	logger.Get().Info().Msgf("Updating endpoint instance group '%s' configuration", endpointInstanceGroupId)

	var update sdk.EndpointInstanceGroupUpdate
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	currentConfig, groupId, err := getEndpointInstanceGroupConfig(ctx, endpointInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	updatedConfig, httpRes, err := client.EndpointInstanceGroupAPI.
		UpdateEndpointInstanceGroupConfig(ctx, groupId).
		IfMatch(strconv.FormatInt(currentConfig.Revision, 10)).
		EndpointInstanceGroupUpdate(update).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(updatedConfig, &endpointInstanceGroupConfigPrintConfig)
}

func EndpointInstanceGroupMetaUpdate(ctx context.Context, endpointInstanceGroupId string, config []byte) error {
	logger.Get().Info().Msgf("Updating endpoint instance group '%s' metadata", endpointInstanceGroupId)

	var meta sdk.GenericMeta
	if err := utils.UnmarshalContent(config, &meta); err != nil {
		return err
	}

	groupId, err := GetEndpointInstanceGroupId(endpointInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.EndpointInstanceGroupAPI.
		UpdateEndpointInstanceGroupMeta(ctx, groupId).
		GenericMeta(meta).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Endpoint instance group '%s' metadata updated", endpointInstanceGroupId)
	return nil
}

func EndpointInstanceGroupInstances(ctx context.Context, endpointInstanceGroupId string) error {
	logger.Get().Info().Msgf("Listing endpoint instances of group '%s'", endpointInstanceGroupId)

	groupId, err := GetEndpointInstanceGroupId(endpointInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.EndpointInstanceGroupAPI.GetEndpointInstanceGroupEndpointInstances(ctx, groupId).SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &endpointInstancePrintConfig)
}

func EndpointInstanceGroupNetworkList(ctx context.Context, endpointInstanceGroupId string) error {
	logger.Get().Info().Msgf("Get network configuration of endpoint instance group '%s'", endpointInstanceGroupId)

	groupId, err := GetEndpointInstanceGroupId(endpointInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkConfiguration, httpRes, err := client.EndpointInstanceGroupAPI.GetEndpointInstanceGroupNetworkConfiguration(ctx, groupId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkConfiguration, &networkEndpointGroupPrintConfig)
}

// EndpointInstanceGroupNetworkReplace creates or replaces the network
// configuration of an endpoint instance group (PUT /config/networking).
//
// The generated SDK request exposes no body setter for this operation, so a
// caller supplied configuration is sent through a raw request carrying the
// If-Match of the group configuration; without a configuration the typed SDK
// call is used, which sends no body at all.
func EndpointInstanceGroupNetworkReplace(ctx context.Context, endpointInstanceGroupId string, config []byte) error {
	logger.Get().Info().Msgf("Replacing network configuration of endpoint instance group '%s'", endpointInstanceGroupId)

	if len(config) == 0 {
		groupId, err := GetEndpointInstanceGroupId(endpointInstanceGroupId)
		if err != nil {
			return err
		}

		client := api.GetApiClient(ctx)

		networkConfiguration, httpRes, err := client.EndpointInstanceGroupAPI.UpdateEndpointInstanceGroupNetworkConfiguration(ctx, groupId).Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			return err
		}

		return formatter.PrintResult(networkConfiguration, &networkEndpointGroupPrintConfig)
	}

	currentConfig, groupId, err := getEndpointInstanceGroupConfig(ctx, endpointInstanceGroupId)
	if err != nil {
		return err
	}

	body, err := api.RawJSONRequest(ctx, http.MethodPut,
		fmt.Sprintf("/api/v2/endpoint-instance-groups/%d/config/networking", groupId),
		config,
		api.IfMatchHeader(strconv.FormatInt(currentConfig.Revision, 10)))
	if err != nil {
		return err
	}

	if len(body) == 0 {
		logger.Get().Info().Msgf("Network configuration of endpoint instance group '%s' replaced", endpointInstanceGroupId)
		return nil
	}

	return utils.PrintRawObject(body, &networkEndpointGroupPrintConfig)
}

func EndpointInstanceGroupNetworkConnections(ctx context.Context, endpointInstanceGroupId string) error {
	logger.Get().Info().Msgf("Listing network connections of endpoint instance group '%s'", endpointInstanceGroupId)

	groupId, err := GetEndpointInstanceGroupId(endpointInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	connections, httpRes, err := client.EndpointInstanceGroupAPI.GetEndpointInstanceGroupNetworkConfigurationConnections(ctx, groupId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(connections, &networkConnectionPrintConfig)
}

func EndpointInstanceGroupNetworkGet(ctx context.Context, endpointInstanceGroupId string, connectionId string) error {
	logger.Get().Info().Msgf("Get network connection '%s' of endpoint instance group '%s'", connectionId, endpointInstanceGroupId)

	groupId, err := GetEndpointInstanceGroupId(endpointInstanceGroupId)
	if err != nil {
		return err
	}

	connectionIdNumerical, err := utils.GetInt64FromString(connectionId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	connection, httpRes, err := client.EndpointInstanceGroupAPI.
		GetEndpointInstanceGroupNetworkConfigurationConnectionById(ctx, groupId, connectionIdNumerical).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(connection, &networkConnectionPrintConfig)
}

func EndpointInstanceGroupNetworkConnectionConfigExample(ctx context.Context) error {
	example := sdk.CreateEndpointInstanceGroupNetworkConnection{
		LogicalNetworkId:     "1",
		Tagged:               true,
		AccessMode:           sdk.NETWORKENDPOINTGROUPALLOWEDACCESSMODE_L2,
		Mtu:                  sdk.PtrInt32(1500),
		ProvidesDefaultRoute: sdk.PtrBool(false),
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(example, nil)
	}
	return formatter.PrintYamlResult(example)
}

func EndpointInstanceGroupNetworkConnect(ctx context.Context, endpointInstanceGroupId string, create sdk.CreateEndpointInstanceGroupNetworkConnection) error {
	logger.Get().Info().Msgf("Connecting endpoint instance group '%s' to logical network '%s'", endpointInstanceGroupId, create.LogicalNetworkId)

	groupId, err := GetEndpointInstanceGroupId(endpointInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	connection, httpRes, err := client.EndpointInstanceGroupAPI.
		CreateEndpointInstanceGroupNetworkConfigurationConnection(ctx, groupId).
		CreateEndpointInstanceGroupNetworkConnection(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(connection, &networkConnectionPrintConfig)
}

func EndpointInstanceGroupNetworkUpdate(ctx context.Context, endpointInstanceGroupId string, connectionId string, update sdk.UpdateNetworkEndpointGroupLogicalNetwork) error {
	logger.Get().Info().Msgf("Updating network connection '%s' of endpoint instance group '%s'", connectionId, endpointInstanceGroupId)

	groupId, err := GetEndpointInstanceGroupId(endpointInstanceGroupId)
	if err != nil {
		return err
	}

	// The SDK types this path parameter as float32 while every other
	// connection operation types it as int64.
	connectionIdNumerical, err := utils.GetFloat32FromString(connectionId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	connection, httpRes, err := client.EndpointInstanceGroupAPI.
		UpdateEndpointInstanceGroupNetworkConfigurationConnection(ctx, groupId, connectionIdNumerical).
		UpdateNetworkEndpointGroupLogicalNetwork(update).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(connection, &networkConnectionPrintConfig)
}

func EndpointInstanceGroupNetworkDisconnect(ctx context.Context, endpointInstanceGroupId string, connectionId string) error {
	logger.Get().Info().Msgf("Disconnecting network connection '%s' from endpoint instance group '%s'", connectionId, endpointInstanceGroupId)

	groupId, err := GetEndpointInstanceGroupId(endpointInstanceGroupId)
	if err != nil {
		return err
	}

	connectionIdNumerical, err := utils.GetInt64FromString(connectionId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.EndpointInstanceGroupAPI.
		DeleteEndpointInstanceGroupNetworkConfigurationConnection(ctx, groupId, connectionIdNumerical).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network connection '%s' disconnected", connectionId)
	return nil
}

// EndpointInstanceGroupACLList lists the security rules of a network
// connection. The SDK types this collection endpoint as returning a single
// LogicalNetworkACL, so the body is parsed raw and rendered as a list; a
// single object response is still rendered correctly.
func EndpointInstanceGroupACLList(ctx context.Context, endpointInstanceGroupId string, connectionId string) error {
	logger.Get().Info().Msgf("Listing security rules of network connection '%s' of endpoint instance group '%s'", connectionId, endpointInstanceGroupId)

	groupId, connectionIdNumerical, err := getGroupAndConnectionId(endpointInstanceGroupId, connectionId)
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

func EndpointInstanceGroupACLGet(ctx context.Context, endpointInstanceGroupId string, connectionId string, ruleId string) error {
	logger.Get().Info().Msgf("Get security rule '%s' of network connection '%s'", ruleId, connectionId)

	groupId, connectionIdNumerical, err := getGroupAndConnectionId(endpointInstanceGroupId, connectionId)
	if err != nil {
		return err
	}

	ruleIdNumerical, err := utils.GetInt64FromString(ruleId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	rule, httpRes, err := client.EndpointInstanceGroupAPI.
		GetEndpointInstanceGroupLogicalNetworkACLById(ctx, groupId, connectionIdNumerical, ruleIdNumerical).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(rule, &logicalNetworkACLPrintConfig)
}

func EndpointInstanceGroupACLConfigExample(ctx context.Context) error {
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

func EndpointInstanceGroupACLAdd(ctx context.Context, endpointInstanceGroupId string, connectionId string, create sdk.CreateLogicalNetworkACL) error {
	logger.Get().Info().Msgf("Adding security rule to network connection '%s' of endpoint instance group '%s'", connectionId, endpointInstanceGroupId)

	groupId, connectionIdNumerical, err := getGroupAndConnectionId(endpointInstanceGroupId, connectionId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	rule, httpRes, err := client.EndpointInstanceGroupAPI.
		CreateEndpointInstanceGroupLogicalNetworkACL(ctx, groupId, connectionIdNumerical).
		CreateLogicalNetworkACL(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(rule, &logicalNetworkACLPrintConfig)
}

func EndpointInstanceGroupACLUpdate(ctx context.Context, endpointInstanceGroupId string, connectionId string, ruleId string, config []byte) error {
	logger.Get().Info().Msgf("Updating security rule '%s' of network connection '%s'", ruleId, connectionId)

	var update sdk.UpdateLogicalNetworkACL
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	groupId, connectionIdNumerical, err := getGroupAndConnectionId(endpointInstanceGroupId, connectionId)
	if err != nil {
		return err
	}

	ruleIdNumerical, err := utils.GetInt64FromString(ruleId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	rule, httpRes, err := client.EndpointInstanceGroupAPI.
		UpdateEndpointInstanceGroupLogicalNetworkACL(ctx, groupId, connectionIdNumerical, ruleIdNumerical).
		UpdateLogicalNetworkACL(update).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(rule, &logicalNetworkACLPrintConfig)
}

func EndpointInstanceGroupACLRemove(ctx context.Context, endpointInstanceGroupId string, connectionId string, ruleId string) error {
	logger.Get().Info().Msgf("Removing security rule '%s' from network connection '%s'", ruleId, connectionId)

	groupId, connectionIdNumerical, err := getGroupAndConnectionId(endpointInstanceGroupId, connectionId)
	if err != nil {
		return err
	}

	ruleIdNumerical, err := utils.GetInt64FromString(ruleId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.EndpointInstanceGroupAPI.
		DeleteEndpointInstanceGroupLogicalNetworkACL(ctx, groupId, connectionIdNumerical, ruleIdNumerical).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Security rule '%s' removed", ruleId)
	return nil
}

func aclCollectionPath(groupId int64, connectionId int64) string {
	return fmt.Sprintf("/api/v2/endpoint-instance-groups/%d/config/networking/connections/%d/security/rules", groupId, connectionId)
}

func getGroupAndConnectionId(endpointInstanceGroupId string, connectionId string) (int64, int64, error) {
	groupId, err := GetEndpointInstanceGroupId(endpointInstanceGroupId)
	if err != nil {
		return 0, 0, err
	}

	connectionIdNumerical, err := utils.GetInt64FromString(connectionId)
	if err != nil {
		return 0, 0, err
	}

	return groupId, connectionIdNumerical, nil
}

func getEndpointInstanceGroup(ctx context.Context, endpointInstanceGroupId string) (*sdk.EndpointInstanceGroup, error) {
	groupId, err := GetEndpointInstanceGroupId(endpointInstanceGroupId)
	if err != nil {
		return nil, err
	}

	client := api.GetApiClient(ctx)

	group, httpRes, err := client.EndpointInstanceGroupAPI.GetEndpointInstanceGroup(ctx, groupId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	return group, nil
}

// getEndpointInstanceGroupConfig returns the group configuration together with
// the numeric group id. The configuration carries its own revision, which
// guards every write below the /config path.
func getEndpointInstanceGroupConfig(ctx context.Context, endpointInstanceGroupId string) (*sdk.EndpointInstanceGroupConfiguration, int64, error) {
	groupId, err := GetEndpointInstanceGroupId(endpointInstanceGroupId)
	if err != nil {
		return nil, 0, err
	}

	client := api.GetApiClient(ctx)

	config, httpRes, err := client.EndpointInstanceGroupAPI.GetEndpointInstanceGroupConfig(ctx, groupId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, 0, err
	}

	return config, groupId, nil
}

func GetEndpointInstanceGroupId(endpointInstanceGroupId string) (int64, error) {
	groupId, err := strconv.ParseInt(endpointInstanceGroupId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid endpoint instance group ID: '%s'", endpointInstanceGroupId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return groupId, nil
}
