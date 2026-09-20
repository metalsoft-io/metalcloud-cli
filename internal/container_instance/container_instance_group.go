package container_instance

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

var containerInstanceGroupPrintConfig = formatter.PrintConfig{
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
			Title: "Infra",
			Order: 3,
		},
		"InstanceCount": {
			Title: "Count",
			Order: 4,
		},
		"DiskSizeGB": {
			Title: "Disk GB",
			Order: 5,
		},
		"VmPoolId": {
			Title: "Pool",
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

var containerInstanceGroupConfigPrintConfig = formatter.PrintConfig{
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

var containerInstanceGroupInterfacePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
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
			Title: "Group",
			Order: 4,
		},
		"NetworkId": {
			Title: "Network",
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

// containerInstanceGroupInterfaceSummary is the table projection of a group
// interface.
//
// sdk.ContainerInstanceGroupInterface cannot decode live responses: it requires
// a top-level `label`, but the API only carries the label inside `config`, so
// the strict decoder fails with "no value given for required property label".
// Both interface endpoints therefore read the body raw.
type containerInstanceGroupInterfaceSummary struct {
	Id               int64                                        `json:"id"`
	Index            float32                                      `json:"index"`
	GroupId          int64                                        `json:"groupId"`
	NetworkId        *int64                                       `json:"networkId"`
	ServiceStatus    string                                       `json:"serviceStatus"`
	CreatedTimestamp string                                       `json:"createdTimestamp"`
	Config           containerInstanceGroupInterfaceConfigSummary `json:"config"`
}

type containerInstanceGroupInterfaceConfigSummary struct {
	Label string `json:"label"`
}

var containerInstanceGroupNetworkEndpointGroupPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Name": {
			MaxWidth: 40,
			Order:    2,
		},
		"SiteId": {
			Title: "Site",
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

var containerInstanceGroupNetworkConnectionPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
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

func ContainerInstanceGroupList(ctx context.Context, infrastructureIdOrLabel string) error {
	logger.Get().Info().Msgf("List all container instance groups for infrastructure '%s'", infrastructureIdOrLabel)

	infraId, err := resolveInfrastructureId(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.ContainerInstanceGroupAPI.
		GetInfrastructureContainerInstanceGroups(ctx, infraId).
		SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &containerInstanceGroupPrintConfig)
}

func ContainerInstanceGroupGet(ctx context.Context, infrastructureIdOrLabel string, containerInstanceGroupId string) error {
	logger.Get().Info().Msgf("Get container instance group '%s' in infrastructure '%s'", containerInstanceGroupId, infrastructureIdOrLabel)

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, containerInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	group, httpRes, err := client.ContainerInstanceGroupAPI.
		GetInfrastructureContainerInstanceGroup(ctx, infraId, groupId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(group, &containerInstanceGroupPrintConfig)
}

func ContainerInstanceGroupConfigExample(ctx context.Context) error {
	example := sdk.CreateContainerInstanceGroup{
		InstanceCount: sdk.PtrFloat32(1),
		DiskSizeGB:    40,
		TypeId:        1,
		OsTemplateId:  1,
		VmPoolId:      1,
		Tags:          []string{"tag1", "tag2"},
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(example, nil)
	}

	return formatter.PrintYamlResult(example)
}

// ContainerInstanceGroupCreate creates a container instance group from a
// complete configuration payload (JSON or YAML).
func ContainerInstanceGroupCreate(ctx context.Context, infrastructureIdOrLabel string, config []byte) error {
	var create sdk.CreateContainerInstanceGroup
	if err := utils.UnmarshalContent(config, &create); err != nil {
		return err
	}

	return containerInstanceGroupCreate(ctx, infrastructureIdOrLabel, create)
}

// ContainerInstanceGroupCreateFromFlags creates a container instance group from
// individual command line flags.
func ContainerInstanceGroupCreateFromFlags(ctx context.Context, infrastructureIdOrLabel string, containerTypeId string, osTemplateId string, vmPoolId string, diskSizeGB string, instanceCount string, tags []string) error {
	typeIdNumerical, err := utils.GetInt64FromString(containerTypeId)
	if err != nil {
		return err
	}

	osTemplateIdNumerical, err := utils.GetInt64FromString(osTemplateId)
	if err != nil {
		return err
	}

	vmPoolIdNumerical, err := utils.GetInt64FromString(vmPoolId)
	if err != nil {
		return err
	}

	diskSize, err := utils.GetFloat32FromString(diskSizeGB)
	if err != nil {
		return err
	}

	create := sdk.CreateContainerInstanceGroup{
		TypeId:       typeIdNumerical,
		OsTemplateId: osTemplateIdNumerical,
		VmPoolId:     vmPoolIdNumerical,
		DiskSizeGB:   diskSize,
		Tags:         tags,
	}

	if instanceCount != "" {
		count, err := utils.GetFloat32FromString(instanceCount)
		if err != nil {
			return err
		}
		create.InstanceCount = &count
	}

	return containerInstanceGroupCreate(ctx, infrastructureIdOrLabel, create)
}

func containerInstanceGroupCreate(ctx context.Context, infrastructureIdOrLabel string, create sdk.CreateContainerInstanceGroup) error {
	logger.Get().Info().Msgf("Create container instance group in infrastructure '%s'", infrastructureIdOrLabel)

	infraId, err := resolveInfrastructureId(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	group, httpRes, err := client.ContainerInstanceGroupAPI.
		CreateContainerInstanceGroup(ctx, infraId).
		CreateContainerInstanceGroup(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(group, &containerInstanceGroupPrintConfig)
}

func ContainerInstanceGroupDelete(ctx context.Context, infrastructureIdOrLabel string, containerInstanceGroupId string) error {
	logger.Get().Info().Msgf("Delete container instance group '%s' from infrastructure '%s'", containerInstanceGroupId, infrastructureIdOrLabel)

	infraId, groupId, revision, err := resolveGroupIdsAndRevision(ctx, infrastructureIdOrLabel, containerInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ContainerInstanceGroupAPI.
		DeleteContainerInstanceGroup(ctx, infraId, groupId).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Container instance group '%s' deleted", containerInstanceGroupId)
	return nil
}

func ContainerInstanceGroupGetConfig(ctx context.Context, infrastructureIdOrLabel string, containerInstanceGroupId string) error {
	logger.Get().Info().Msgf("Get configuration of container instance group '%s' in infrastructure '%s'", containerInstanceGroupId, infrastructureIdOrLabel)

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, containerInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	config, httpRes, err := client.ContainerInstanceGroupAPI.
		GetContainerInstanceGroupConfigInfo(ctx, infraId, groupId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(config, &containerInstanceGroupConfigPrintConfig)
}

func ContainerInstanceGroupUpdateConfig(ctx context.Context, infrastructureIdOrLabel string, containerInstanceGroupId string, config []byte) error {
	logger.Get().Info().Msgf("Update configuration of container instance group '%s' in infrastructure '%s'", containerInstanceGroupId, infrastructureIdOrLabel)

	var update sdk.UpdateContainerInstanceGroup
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, containerInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// The config sub-resource is guarded by its own revision, not by the
	// container instance group revision.
	currentConfig, httpRes, err := client.ContainerInstanceGroupAPI.
		GetContainerInstanceGroupConfigInfo(ctx, infraId, groupId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	updatedConfig, httpRes, err := client.ContainerInstanceGroupAPI.
		UpdateContainerInstanceGroupConfig(ctx, infraId, groupId).
		UpdateContainerInstanceGroup(update).
		IfMatch(strconv.FormatInt(currentConfig.Revision, 10)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(updatedConfig, &containerInstanceGroupConfigPrintConfig)
}

func ContainerInstanceGroupUpdateMeta(ctx context.Context, infrastructureIdOrLabel string, containerInstanceGroupId string, config []byte) error {
	logger.Get().Info().Msgf("Update metadata of container instance group '%s' in infrastructure '%s'", containerInstanceGroupId, infrastructureIdOrLabel)

	var update sdk.UpdateContainerInstanceGroupMeta
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, containerInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	group, httpRes, err := client.ContainerInstanceGroupAPI.
		PatchContainerInstanceGroupMeta(ctx, infraId, groupId).
		UpdateContainerInstanceGroupMeta(update).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(group, &containerInstanceGroupPrintConfig)
}

func ContainerInstanceGroupInstances(ctx context.Context, infrastructureIdOrLabel string, containerInstanceGroupId string) error {
	logger.Get().Info().Msgf("List instances of container instance group '%s' in infrastructure '%s'", containerInstanceGroupId, infrastructureIdOrLabel)

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, containerInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.ContainerInstanceGroupAPI.
		GetContainerInstanceGroupContainerInstances(ctx, infraId, groupId).
		SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &containerInstancePrintConfig)
}

func ContainerInstanceGroupInterfaces(ctx context.Context, infrastructureIdOrLabel string, containerInstanceGroupId string) error {
	logger.Get().Info().Msgf("List interfaces of container instance group '%s' in infrastructure '%s'", containerInstanceGroupId, infrastructureIdOrLabel)

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, containerInstanceGroupId)
	if err != nil {
		return err
	}

	rawItems, meta, err := utils.FetchAllPagesRaw(func(page float32) (*http.Response, error) {
		return api.DoJSONRequest(ctx, http.MethodGet,
			fmt.Sprintf("/api/v2/infrastructures/%d/container-instance-groups/%d/interfaces?page=%.0f&limit=100&sortBy=id:ASC", infraId, groupId, page), nil)
	})
	if err != nil {
		return err
	}

	records, err := utils.UnmarshalRawItems[containerInstanceGroupInterfaceSummary](rawItems)
	if err != nil {
		return fmt.Errorf("failed to parse container instance group interfaces: %w", err)
	}

	return utils.PrintAllRaw(rawItems, records, meta, len(records), &containerInstanceGroupInterfacePrintConfig)
}

func ContainerInstanceGroupInterfaceGet(ctx context.Context, infrastructureIdOrLabel string, containerInstanceGroupId string, interfaceId string) error {
	logger.Get().Info().Msgf("Get interface '%s' of container instance group '%s' in infrastructure '%s'", interfaceId, containerInstanceGroupId, infrastructureIdOrLabel)

	interfaceIdNumerical, err := utils.GetInt64FromString(interfaceId)
	if err != nil {
		return err
	}

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, containerInstanceGroupId)
	if err != nil {
		return err
	}

	body, err := api.RawJSONRequest(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/infrastructures/%d/container-instance-groups/%d/interfaces/%d", infraId, groupId, interfaceIdNumerical), nil, nil)
	if err != nil {
		return err
	}

	if formatter.IsNativeFormat() {
		return utils.PrintRawObject(body, &containerInstanceGroupInterfacePrintConfig)
	}

	var record containerInstanceGroupInterfaceSummary
	if err := json.Unmarshal(body, &record); err != nil {
		return fmt.Errorf("failed to parse container instance group interface: %w", err)
	}

	return formatter.PrintResult(record, &containerInstanceGroupInterfacePrintConfig)
}

func ContainerInstanceGroupApplyType(ctx context.Context, infrastructureIdOrLabel string, containerInstanceGroupId string, containerTypeId string) error {
	logger.Get().Info().Msgf("Apply container type '%s' on container instance group '%s' in infrastructure '%s'", containerTypeId, containerInstanceGroupId, infrastructureIdOrLabel)

	containerTypeIdNumerical, err := utils.GetInt64FromString(containerTypeId)
	if err != nil {
		return err
	}

	infraId, groupId, revision, err := resolveGroupIdsAndRevision(ctx, infrastructureIdOrLabel, containerInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	group, httpRes, err := client.ContainerInstanceGroupAPI.
		ApplyContainerTypeOnContainerInstanceGroup(ctx, infraId, groupId, containerTypeIdNumerical).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(group, &containerInstanceGroupPrintConfig)
}

func ContainerInstanceGroupNetworkConfiguration(ctx context.Context, infrastructureIdOrLabel string, containerInstanceGroupId string) error {
	logger.Get().Info().Msgf("Get network configuration of container instance group '%s' in infrastructure '%s'", containerInstanceGroupId, infrastructureIdOrLabel)

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, containerInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkConfiguration, httpRes, err := client.ContainerInstanceGroupAPI.
		GetContainerInstanceGroupNetworkConfiguration(ctx, infraId, groupId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkConfiguration, &containerInstanceGroupNetworkEndpointGroupPrintConfig)
}

func ContainerInstanceGroupNetworkList(ctx context.Context, infrastructureIdOrLabel string, containerInstanceGroupId string) error {
	logger.Get().Info().Msgf("List network connections of container instance group '%s' in infrastructure '%s'", containerInstanceGroupId, infrastructureIdOrLabel)

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, containerInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	connections, httpRes, err := client.ContainerInstanceGroupAPI.
		GetContainerInstanceGroupNetworkConfigurationConnections(ctx, infraId, groupId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(connections.Data, &containerInstanceGroupNetworkConnectionPrintConfig)
}

func ContainerInstanceGroupNetworkGet(ctx context.Context, infrastructureIdOrLabel string, containerInstanceGroupId string, connectionId string) error {
	logger.Get().Info().Msgf("Get network connection '%s' of container instance group '%s' in infrastructure '%s'", connectionId, containerInstanceGroupId, infrastructureIdOrLabel)

	connectionIdNumerical, err := utils.GetInt64FromString(connectionId)
	if err != nil {
		return err
	}

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, containerInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	connection, httpRes, err := client.ContainerInstanceGroupAPI.
		GetContainerInstanceGroupNetworkConfigurationConnectionById(ctx, infraId, groupId, connectionIdNumerical).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(connection, &containerInstanceGroupNetworkConnectionPrintConfig)
}

func ContainerInstanceGroupNetworkConnect(ctx context.Context, infrastructureIdOrLabel string, containerInstanceGroupId string, logicalNetworkId string, accessMode string, tagged string, redundancy string) error {
	logger.Get().Info().Msgf("Connect container instance group '%s' in infrastructure '%s' to logical network '%s'", containerInstanceGroupId, infrastructureIdOrLabel, logicalNetworkId)

	isTagged, err := strconv.ParseBool(tagged)
	if err != nil {
		err := fmt.Errorf("invalid tagged value: '%s'", tagged)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, containerInstanceGroupId)
	if err != nil {
		return err
	}

	payload := sdk.CreateContainerInstanceGroupNetworkConnection{
		LogicalNetworkId: logicalNetworkId,
		AccessMode:       sdk.NetworkEndpointGroupAllowedAccessMode(accessMode),
		Tagged:           isTagged,
	}

	if redundancy != "" {
		payload.Redundancy = *sdk.NewNullableRedundancyConfig(&sdk.RedundancyConfig{
			Mode: sdk.NetworkEndpointGroupRedundancyMode(redundancy),
		})
	}

	client := api.GetApiClient(ctx)

	connection, httpRes, err := client.ContainerInstanceGroupAPI.
		CreateContainerInstanceGroupNetworkConfigurationConnection(ctx, infraId, groupId).
		CreateContainerInstanceGroupNetworkConnection(payload).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(connection, &containerInstanceGroupNetworkConnectionPrintConfig)
}

func ContainerInstanceGroupNetworkUpdate(ctx context.Context, infrastructureIdOrLabel string, containerInstanceGroupId string, connectionId string, accessMode string, tagged string, redundancy string) error {
	logger.Get().Info().Msgf("Update network connection '%s' of container instance group '%s' in infrastructure '%s'", connectionId, containerInstanceGroupId, infrastructureIdOrLabel)

	connectionIdNumerical, err := utils.GetInt64FromString(connectionId)
	if err != nil {
		return err
	}

	payload := sdk.UpdateContainerInstanceGroupNetworkConnection{}

	if accessMode != "" {
		accessModeValue := sdk.NetworkEndpointGroupAllowedAccessMode(accessMode)
		payload.AccessMode = &accessModeValue
	}

	if tagged != "" {
		isTagged, err := strconv.ParseBool(tagged)
		if err != nil {
			err := fmt.Errorf("invalid tagged value: '%s'", tagged)
			logger.Get().Error().Err(err).Msg("")
			return err
		}
		payload.Tagged = &isTagged
	}

	if redundancy != "" {
		payload.Redundancy = *sdk.NewNullableRedundancyConfig(&sdk.RedundancyConfig{
			Mode: sdk.NetworkEndpointGroupRedundancyMode(redundancy),
		})
	}

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, containerInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	connection, httpRes, err := client.ContainerInstanceGroupAPI.
		UpdateContainerInstanceGroupNetworkConfigurationConnection(ctx, infraId, groupId, connectionIdNumerical).
		UpdateContainerInstanceGroupNetworkConnection(payload).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(connection, &containerInstanceGroupNetworkConnectionPrintConfig)
}

func ContainerInstanceGroupNetworkDisconnect(ctx context.Context, infrastructureIdOrLabel string, containerInstanceGroupId string, connectionId string) error {
	logger.Get().Info().Msgf("Disconnect network connection '%s' from container instance group '%s' in infrastructure '%s'", connectionId, containerInstanceGroupId, infrastructureIdOrLabel)

	connectionIdNumerical, err := utils.GetInt64FromString(connectionId)
	if err != nil {
		return err
	}

	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, containerInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ContainerInstanceGroupAPI.
		DeleteContainerInstanceGroupNetworkConfigurationConnection(ctx, infraId, groupId, connectionIdNumerical).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network connection '%s' successfully removed", connectionId)
	return nil
}

func resolveGroupIds(ctx context.Context, infrastructureIdOrLabel string, containerInstanceGroupId string) (int64, int64, error) {
	groupId, err := GetContainerInstanceGroupId(containerInstanceGroupId)
	if err != nil {
		return 0, 0, err
	}

	infraId, err := resolveInfrastructureId(ctx, infrastructureIdOrLabel)
	if err != nil {
		return 0, 0, err
	}

	return infraId, groupId, nil
}

func resolveGroupIdsAndRevision(ctx context.Context, infrastructureIdOrLabel string, containerInstanceGroupId string) (int64, int64, string, error) {
	infraId, groupId, err := resolveGroupIds(ctx, infrastructureIdOrLabel, containerInstanceGroupId)
	if err != nil {
		return 0, 0, "", err
	}

	client := api.GetApiClient(ctx)

	group, httpRes, err := client.ContainerInstanceGroupAPI.
		GetInfrastructureContainerInstanceGroup(ctx, infraId, groupId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, 0, "", err
	}

	return infraId, groupId, strconv.FormatInt(group.Revision, 10), nil
}

func GetContainerInstanceGroupId(containerInstanceGroupId string) (int64, error) {
	containerInstanceGroupIdNumerical, err := utils.GetInt64FromString(containerInstanceGroupId)
	if err != nil {
		err := fmt.Errorf("invalid container instance group ID: '%s'", containerInstanceGroupId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return containerInstanceGroupIdNumerical, nil
}
