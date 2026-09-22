package vm_instance

import (
	"context"
	"encoding/json"
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

var vmInstanceGroupPrintConfig = formatter.PrintConfig{
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
		"InstanceCount": {
			Title: "Count",
			Order: 4,
		},
		"ServiceStatus": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       5,
		},
		"DiskSizeGB": {
			Title: "Disk GB",
			Order: 6,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       8,
		},
	},
}

var vmInstanceGroupConfigPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Revision": {
			Title: "Revision",
			Order: 1,
		},
		"Label": {
			MaxWidth: 30,
			Order:    2,
		},
		"InstanceCount": {
			Title: "Count",
			Order: 3,
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
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
	},
}

var vmInstanceGroupInterfacePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Config.Label": {
			Title:    "Label",
			MaxWidth: 30,
			Order:    2,
		},
		"Index": {
			Title: "Index",
			Order: 3,
		},
		"GroupId": {
			Title: "Group ID",
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
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
	},
}

// vmInstanceGroupInterfaceSummary is the table projection of a group interface.
//
// sdk.VMInstanceGroupInterface cannot decode live responses: it requires a
// top-level `label`, but the API only carries the label inside `config`, so the
// strict decoder fails with "no value given for required property label". Both
// interface endpoints therefore read the body raw.
type vmInstanceGroupInterfaceSummary struct {
	Id               int64                                 `json:"id"`
	Index            float32                               `json:"index"`
	GroupId          int64                                 `json:"groupId"`
	NetworkId        *int64                                `json:"networkId"`
	ServiceStatus    string                                `json:"serviceStatus"`
	CreatedTimestamp string                                `json:"createdTimestamp"`
	Config           vmInstanceGroupInterfaceConfigSummary `json:"config"`
}

type vmInstanceGroupInterfaceConfigSummary struct {
	Label string `json:"label"`
}

var vmInstanceGroupNetworkEndpointGroupPrintConfig = formatter.PrintConfig{
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

var vmInstanceGroupNetworkConnectionPrintConfig = formatter.PrintConfig{
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
		"ProvidesDefaultRoute": {
			Title: "Default Route",
			Order: 5,
		},
		"DisableAutoIpAllocation": {
			Title: "No Auto IP",
			Order: 6,
		},
	},
}

func VMInstanceGroupGet(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string) error {
	logger.Get().Info().Msgf("Get VM instance group '%s' in infrastructure '%s'", vmInstanceGroupId, infrastructureIdOrLabel)

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, vmInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmInstanceGroup, httpRes, err := client.VMInstanceGroupAPI.
		GetInfrastructureVMInstanceGroup(ctx, infraId, groupId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmInstanceGroup, &vmInstanceGroupPrintConfig)
}

func VMInstanceGroupList(ctx context.Context, infrastructureIdOrLabel string) error {
	logger.Get().Info().Msgf("List all VM instance groups for infrastructure '%s'", infrastructureIdOrLabel)

	infraId, err := resolveInfrastructureId(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.VMInstanceGroupAPI.
		GetInfrastructureVMInstanceGroups(ctx, infraId).
		SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &vmInstanceGroupPrintConfig)
}

func VMInstanceGroupCreate(ctx context.Context, infrastructureIdOrLabel string, vmTypeId string, diskSizeGB string, instanceCount string, osTemplateId string) error {
	logger.Get().Info().Msgf("Create new VM instance group in infrastructure '%s'", infrastructureIdOrLabel)

	vmTypeIdNumerical, err := utils.GetInt64FromString(vmTypeId)
	if err != nil {
		return err
	}

	diskSizeGBNumerical, err := utils.GetFloat32FromString(diskSizeGB)
	if err != nil {
		return err
	}

	instanceCountNumerical, err := utils.GetFloat32FromString(instanceCount)
	if err != nil {
		return err
	}

	payload := sdk.CreateVMInstanceGroup{
		TypeId:        vmTypeIdNumerical,
		DiskSizeGB:    diskSizeGBNumerical,
		InstanceCount: &instanceCountNumerical,
	}

	if osTemplateId != "" {
		payload.OsTemplateId, err = utils.GetInt64FromString(osTemplateId)
		if err != nil {
			return err
		}
	}

	infraId, err := resolveInfrastructureId(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmInstanceGroupInfo, httpRes, err := client.VMInstanceGroupAPI.
		CreateVMInstanceGroup(ctx, infraId).
		CreateVMInstanceGroup(payload).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmInstanceGroupInfo, &vmInstanceGroupPrintConfig)
}

func VMInstanceGroupUpdate(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string, label string, customVariables []byte) error {
	logger.Get().Info().Msgf("Update VM instance group '%s' in infrastructure '%s'", vmInstanceGroupId, infrastructureIdOrLabel)

	payload := sdk.UpdateVMInstanceGroup{}

	if label != "" {
		payload.Label = &label
	}

	if customVariables != nil {
		err := json.Unmarshal(customVariables, &payload.CustomVariables)
		if err != nil {
			return err
		}
	}

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, vmInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// The config sub-resource is guarded by its own revision, not by the VM
	// instance group revision.
	currentConfig, httpRes, err := client.VMInstanceGroupAPI.
		GetVMInstanceGroupConfigInfo(ctx, infraId, groupId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	vmInstanceGroupConfig, httpRes, err := client.VMInstanceGroupAPI.
		UpdateVMInstanceGroupConfig(ctx, infraId, groupId).
		UpdateVMInstanceGroup(payload).
		IfMatch(strconv.FormatInt(currentConfig.Revision, 10)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmInstanceGroupConfig, &vmInstanceGroupConfigPrintConfig)
}

func VMInstanceGroupDelete(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string) error {
	logger.Get().Info().Msgf("Delete VM instance group '%s' from infrastructure '%s'", vmInstanceGroupId, infrastructureIdOrLabel)

	infraId, groupId, revision, err := resolveGroupIdsAndRevision(ctx, infrastructureIdOrLabel, vmInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.VMInstanceGroupAPI.
		DeleteVMInstanceGroup(ctx, infraId, groupId).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM instance group '%s' deleted", vmInstanceGroupId)
	return nil
}

func VMInstanceGroupInstances(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string) error {
	logger.Get().Info().Msgf("List instances of VM instance group '%s' in infrastructure '%s'", vmInstanceGroupId, infrastructureIdOrLabel)

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, vmInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.VMInstanceGroupAPI.
		GetVMInstanceGroupVMInstances(ctx, infraId, groupId).
		SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &vmInstancePrintConfig)
}

func VMInstanceGroupGetConfig(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string) error {
	logger.Get().Info().Msgf("Get configuration of VM instance group '%s' in infrastructure '%s'", vmInstanceGroupId, infrastructureIdOrLabel)

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, vmInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	config, httpRes, err := client.VMInstanceGroupAPI.
		GetVMInstanceGroupConfigInfo(ctx, infraId, groupId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(config, &vmInstanceGroupConfigPrintConfig)
}

func VMInstanceGroupUpdateMeta(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string, config []byte) error {
	logger.Get().Info().Msgf("Update metadata of VM instance group '%s' in infrastructure '%s'", vmInstanceGroupId, infrastructureIdOrLabel)

	var update sdk.UpdateVMInstanceGroupMeta
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, vmInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	group, httpRes, err := client.VMInstanceGroupAPI.
		PatchVMInstanceGroupMeta(ctx, infraId, groupId).
		UpdateVMInstanceGroupMeta(update).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(group, &vmInstanceGroupPrintConfig)
}

func VMInstanceGroupApplyType(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string, vmTypeId string) error {
	logger.Get().Info().Msgf("Apply VM type '%s' on VM instance group '%s' in infrastructure '%s'", vmTypeId, vmInstanceGroupId, infrastructureIdOrLabel)

	vmTypeIdNumerical, err := utils.GetInt64FromString(vmTypeId)
	if err != nil {
		return err
	}

	infraId, groupId, revision, err := resolveGroupIdsAndRevision(ctx, infrastructureIdOrLabel, vmInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	group, httpRes, err := client.VMInstanceGroupAPI.
		ApplyVMTypeOnVMInstanceGroup(ctx, infraId, groupId, vmTypeIdNumerical).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(group, &vmInstanceGroupPrintConfig)
}

func VMInstanceGroupInterfaces(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string) error {
	logger.Get().Info().Msgf("List interfaces of VM instance group '%s' in infrastructure '%s'", vmInstanceGroupId, infrastructureIdOrLabel)

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, vmInstanceGroupId)
	if err != nil {
		return err
	}

	rawItems, meta, err := utils.FetchAllPagesRaw(func(page float32) (*http.Response, error) {
		return api.DoJSONRequest(ctx, http.MethodGet,
			fmt.Sprintf("/api/v2/infrastructures/%d/vm-instance-groups/%d/interfaces?page=%.0f&limit=100&sortBy=id:ASC", infraId, groupId, page), nil)
	})
	if err != nil {
		return err
	}

	records, err := utils.UnmarshalRawItems[vmInstanceGroupInterfaceSummary](rawItems)
	if err != nil {
		return fmt.Errorf("failed to parse VM instance group interfaces: %w", err)
	}

	return utils.PrintAllRaw(rawItems, records, meta, len(records), &vmInstanceGroupInterfacePrintConfig)
}

func VMInstanceGroupInterfaceGet(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string, interfaceId string) error {
	logger.Get().Info().Msgf("Get interface '%s' of VM instance group '%s' in infrastructure '%s'", interfaceId, vmInstanceGroupId, infrastructureIdOrLabel)

	interfaceIdNumerical, err := utils.GetInt64FromString(interfaceId)
	if err != nil {
		return err
	}

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, vmInstanceGroupId)
	if err != nil {
		return err
	}

	body, err := api.RawJSONRequest(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/infrastructures/%d/vm-instance-groups/%d/interfaces/%d", infraId, groupId, interfaceIdNumerical), nil, nil)
	if err != nil {
		return err
	}

	if formatter.IsNativeFormat() {
		return utils.PrintRawObject(body, &vmInstanceGroupInterfacePrintConfig)
	}

	var record vmInstanceGroupInterfaceSummary
	if err := json.Unmarshal(body, &record); err != nil {
		return fmt.Errorf("failed to parse VM instance group interface: %w", err)
	}

	return formatter.PrintResult(record, &vmInstanceGroupInterfacePrintConfig)
}

func VMInstanceGroupNetworkConfiguration(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string) error {
	logger.Get().Info().Msgf("Get network configuration of VM instance group '%s' in infrastructure '%s'", vmInstanceGroupId, infrastructureIdOrLabel)

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, vmInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkConfiguration, httpRes, err := client.VMInstanceGroupAPI.
		GetVmInstanceGroupNetworkConfiguration(ctx, infraId, groupId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkConfiguration, &vmInstanceGroupNetworkEndpointGroupPrintConfig)
}

func VMInstanceGroupNetworkConnections(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string) error {
	logger.Get().Info().Msgf("List network connections of VM instance group '%s' in infrastructure '%s'", vmInstanceGroupId, infrastructureIdOrLabel)

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, vmInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	connections, httpRes, err := client.VMInstanceGroupAPI.
		GetVMInstanceGroupNetworkConfigurationConnections(ctx, infraId, groupId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(connections.Data, &vmInstanceGroupNetworkConnectionPrintConfig)
}

func VMInstanceGroupNetworkGet(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string, connectionId string) error {
	logger.Get().Info().Msgf("Get network connection '%s' of VM instance group '%s' in infrastructure '%s'", connectionId, vmInstanceGroupId, infrastructureIdOrLabel)

	connectionIdNumerical, err := utils.GetInt64FromString(connectionId)
	if err != nil {
		return err
	}

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, vmInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	connection, httpRes, err := client.VMInstanceGroupAPI.
		GetVMInstanceGroupNetworkConfigurationConnectionById(ctx, infraId, groupId, connectionIdNumerical).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(connection, &vmInstanceGroupNetworkConnectionPrintConfig)
}

func VMInstanceGroupNetworkConfigExample(ctx context.Context) error {
	example := sdk.CreateVMInstanceGroupNetworkConnection{
		LogicalNetworkId:        "1",
		Tagged:                  true,
		AccessMode:              sdk.NETWORKENDPOINTGROUPALLOWEDACCESSMODE_L2,
		Mtu:                     sdk.PtrInt32(1500),
		ProvidesDefaultRoute:    sdk.PtrBool(false),
		DisableAutoIpAllocation: sdk.PtrBool(false),
		Redundancy: *sdk.NewNullableRedundancyConfig(&sdk.RedundancyConfig{
			Mode: sdk.NETWORKENDPOINTGROUPREDUNDANCYMODE_ACTIVE_BACKUP,
		}),
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(example, nil)
	}

	return formatter.PrintYamlResult(example)
}

// VMInstanceGroupNetworkConnect connects a VM instance group to a logical
// network from a complete configuration payload (JSON or YAML).
func VMInstanceGroupNetworkConnect(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string, config []byte) error {
	var create sdk.CreateVMInstanceGroupNetworkConnection
	if err := utils.UnmarshalContent(config, &create); err != nil {
		return err
	}

	return vmInstanceGroupNetworkConnect(ctx, infrastructureIdOrLabel, vmInstanceGroupId, create)
}

// VMInstanceGroupNetworkConnectFromFlags connects a VM instance group to a
// logical network from individual command line flags.
func VMInstanceGroupNetworkConnectFromFlags(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string, logicalNetworkId string, accessMode string, tagged string, redundancy string, mtu string) error {
	create := sdk.CreateVMInstanceGroupNetworkConnection{
		LogicalNetworkId: logicalNetworkId,
		AccessMode:       sdk.NetworkEndpointGroupAllowedAccessMode(accessMode),
	}

	if tagged != "" {
		isTagged, err := parseTagged(tagged)
		if err != nil {
			return err
		}
		create.Tagged = isTagged
	}

	if redundancy != "" {
		create.Redundancy = *sdk.NewNullableRedundancyConfig(&sdk.RedundancyConfig{
			Mode: sdk.NetworkEndpointGroupRedundancyMode(redundancy),
		})
	}

	if mtu != "" {
		mtuNumerical, err := utils.GetInt64FromString(mtu)
		if err != nil {
			return err
		}
		create.Mtu = sdk.PtrInt32(int32(mtuNumerical))
	}

	return vmInstanceGroupNetworkConnect(ctx, infrastructureIdOrLabel, vmInstanceGroupId, create)
}

func vmInstanceGroupNetworkConnect(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string, create sdk.CreateVMInstanceGroupNetworkConnection) error {
	logger.Get().Info().Msgf("Connect VM instance group '%s' in infrastructure '%s' to logical network '%s'", vmInstanceGroupId, infrastructureIdOrLabel, create.LogicalNetworkId)

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, vmInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	connection, httpRes, err := client.VMInstanceGroupAPI.
		CreateVMInstanceGroupNetworkConfigurationConnection(ctx, infraId, groupId).
		CreateVMInstanceGroupNetworkConnection(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(connection, &vmInstanceGroupNetworkConnectionPrintConfig)
}

// VMInstanceGroupNetworkUpdate updates one network connection of a VM instance
// group from a complete configuration payload (JSON or YAML).
func VMInstanceGroupNetworkUpdate(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string, connectionId string, config []byte) error {
	var update sdk.UpdateVMInstanceGroupNetworkConnection
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	return vmInstanceGroupNetworkUpdate(ctx, infrastructureIdOrLabel, vmInstanceGroupId, connectionId, update)
}

// VMInstanceGroupNetworkUpdateFromFlags updates one network connection of a VM
// instance group from individual command line flags.
func VMInstanceGroupNetworkUpdateFromFlags(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string, connectionId string, accessMode string, tagged string, redundancy string, mtu string) error {
	update := sdk.UpdateVMInstanceGroupNetworkConnection{}

	if accessMode != "" {
		accessModeValue := sdk.NetworkEndpointGroupAllowedAccessMode(accessMode)
		update.AccessMode = &accessModeValue
	}

	if tagged != "" {
		isTagged, err := parseTagged(tagged)
		if err != nil {
			return err
		}
		update.Tagged = &isTagged
	}

	if redundancy != "" {
		update.Redundancy = *sdk.NewNullableRedundancyConfig(&sdk.RedundancyConfig{
			Mode: sdk.NetworkEndpointGroupRedundancyMode(redundancy),
		})
	}

	if mtu != "" {
		mtuNumerical, err := utils.GetInt64FromString(mtu)
		if err != nil {
			return err
		}
		update.Mtu = sdk.PtrInt32(int32(mtuNumerical))
	}

	return vmInstanceGroupNetworkUpdate(ctx, infrastructureIdOrLabel, vmInstanceGroupId, connectionId, update)
}

func vmInstanceGroupNetworkUpdate(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string, connectionId string, update sdk.UpdateVMInstanceGroupNetworkConnection) error {
	logger.Get().Info().Msgf("Update network connection '%s' of VM instance group '%s' in infrastructure '%s'", connectionId, vmInstanceGroupId, infrastructureIdOrLabel)

	connectionIdNumerical, err := utils.GetInt64FromString(connectionId)
	if err != nil {
		return err
	}

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, vmInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	connection, httpRes, err := client.VMInstanceGroupAPI.
		UpdateVMInstanceGroupNetworkConfigurationConnection(ctx, infraId, groupId, connectionIdNumerical).
		UpdateVMInstanceGroupNetworkConnection(update).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(connection, &vmInstanceGroupNetworkConnectionPrintConfig)
}

func VMInstanceGroupNetworkDisconnect(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string, connectionId string) error {
	logger.Get().Info().Msgf("Disconnect network connection '%s' from VM instance group '%s' in infrastructure '%s'", connectionId, vmInstanceGroupId, infrastructureIdOrLabel)

	connectionIdNumerical, err := utils.GetInt64FromString(connectionId)
	if err != nil {
		return err
	}

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, vmInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.VMInstanceGroupAPI.
		DeleteVMInstanceGroupNetworkConfigurationConnection(ctx, infraId, groupId, connectionIdNumerical).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network connection '%s' successfully removed", connectionId)
	return nil
}

func parseTagged(tagged string) (bool, error) {
	isTagged, err := strconv.ParseBool(tagged)
	if err != nil {
		err := fmt.Errorf("invalid tagged value: '%s'", tagged)
		logger.Get().Error().Err(err).Msg("")
		return false, err
	}

	return isTagged, nil
}

func resolveGroupIds(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string) (int64, int64, error) {
	groupId, err := GetVMInstanceGroupId(vmInstanceGroupId)
	if err != nil {
		return 0, 0, err
	}

	infraId, err := resolveInfrastructureId(ctx, infrastructureIdOrLabel)
	if err != nil {
		return 0, 0, err
	}

	return infraId, groupId, nil
}

func resolveGroupIdsAndRevision(ctx context.Context, infrastructureIdOrLabel string, vmInstanceGroupId string) (int64, int64, string, error) {
	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, vmInstanceGroupId)
	if err != nil {
		return 0, 0, "", err
	}

	client := api.GetApiClient(ctx)

	group, httpRes, err := client.VMInstanceGroupAPI.
		GetInfrastructureVMInstanceGroup(ctx, infraId, groupId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, 0, "", err
	}

	return infraId, groupId, strconv.FormatInt(group.Revision, 10), nil
}

func GetVMInstanceGroupId(vmInstanceGroupId string) (int64, error) {
	vmInstanceGroupIdNumerical, err := utils.GetInt64FromString(vmInstanceGroupId)
	if err != nil {
		err := fmt.Errorf("invalid VM instance group ID: '%s'", vmInstanceGroupId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return vmInstanceGroupIdNumerical, nil
}
