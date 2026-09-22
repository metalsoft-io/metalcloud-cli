package server_instance

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/internal/infrastructure"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// ServerInstanceFilters carries the query filters shared by the global and the
// infrastructure scoped server instance listings.
type ServerInstanceFilters struct {
	InfrastructureId   []string
	GroupId            []string
	ServerId           []string
	ServiceStatus      []string
	ConfigServerId     []string
	ConfigDeployStatus []string
	ConfigDeployType   []string
}

var serverInstancePrintConfig = formatter.PrintConfig{
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
		"GroupId": {
			Title: "Group ID",
			Order: 4,
		},
		"ServiceStatus": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       5,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
	},
}

var serverInstanceConfigPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Label": {
			Title: "Label",
			Order: 1,
		},
		"GroupId": {
			Title: "Group ID",
			Order: 2,
		},
		"ServerTypeId": {
			Title: "Server Type",
			Order: 3,
		},
		"ServerId": {
			Title: "Server ID",
			Order: 4,
		},
		"OsTemplateId": {
			Title: "OS Template",
			Order: 5,
		},
		"Hostname": {
			Title: "Hostname",
			Order: 6,
		},
		"DeployType": {
			Title: "Deploy Type",
			Order: 7,
		},
		"DeployStatus": {
			Title:       "Deploy Status",
			Transformer: formatter.FormatStatusValue,
			Order:       8,
		},
		"Revision": {
			Title: "Revision",
			Order: 9,
		},
	},
}

var serverInstanceDrivePrintConfig = formatter.PrintConfig{
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
			Title: "Drive Group",
			Order: 3,
		},
		"SizeMb": {
			Title: "Size (MB)",
			Order: 4,
		},
		"StorageType": {
			Title: "Storage Type",
			Order: 5,
		},
		"StoragePoolId": {
			Title: "Storage Pool",
			Order: 6,
		},
		"ServiceStatus": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       7,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       8,
		},
	},
}

var serverInstanceInterfacePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Label": {
			MaxWidth: 30,
			Order:    2,
		},
		"InstanceId": {
			Title: "Instance ID",
			Order: 3,
		},
		"Index": {
			Title: "Index",
			Order: 4,
		},
		"CapacityMbps": {
			Title: "Capacity (Mbps)",
			Order: 5,
		},
		"NetworkId": {
			Title: "Network ID",
			Order: 6,
		},
		"ServiceStatus": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       7,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       8,
		},
	},
}

var serverInstanceInterfaceConfigPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Label": {
			MaxWidth: 30,
			Order:    1,
		},
		"InstanceId": {
			Title: "Instance ID",
			Order: 2,
		},
		"Index": {
			Title: "Index",
			Order: 3,
		},
		"CapacityMbps": {
			Title: "Capacity (Mbps)",
			Order: 4,
		},
		"NetworkId": {
			Title: "Network ID",
			Order: 5,
		},
		"DeployType": {
			Title: "Deploy Type",
			Order: 6,
		},
		"DeployStatus": {
			Title:       "Deploy Status",
			Transformer: formatter.FormatStatusValue,
			Order:       7,
		},
		"Revision": {
			Title: "Revision",
			Order: 8,
		},
	},
}

var serverInstanceStatisticsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ServerStatus": {
			Title: "Server Status",
			Order: 1,
		},
		"Site": {
			Title: "Site",
			Order: 2,
		},
	},
}

// ServerInstanceList lists server instances. When infrastructureIdOrLabel is
// empty the global /server-instances collection is used, otherwise the listing
// is scoped to that infrastructure, which is resolved by ID or by label.
func ServerInstanceList(ctx context.Context, infrastructureIdOrLabel string, filters ServerInstanceFilters) error {
	client := api.GetApiClient(ctx)

	if infrastructureIdOrLabel != "" {
		logger.Get().Info().Msgf("Listing server instances for infrastructure '%s'", infrastructureIdOrLabel)

		infra, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
		if err != nil {
			return err
		}

		request := client.ServerInstanceAPI.GetInfrastructureServerInstances(ctx, int64(infra.Id)).SortBy([]string{"id:ASC"})
		if len(filters.GroupId) > 0 {
			request = request.FilterGroupId(utils.ProcessFilterStringSlice(filters.GroupId))
		}
		if len(filters.ServerId) > 0 {
			request = request.FilterServerId(utils.ProcessFilterStringSlice(filters.ServerId))
		}
		if len(filters.ServiceStatus) > 0 {
			request = request.FilterServiceStatus(utils.ProcessFilterStringSlice(filters.ServiceStatus))
		}
		if len(filters.ConfigServerId) > 0 {
			request = request.FilterConfigServerId(utils.ProcessFilterStringSlice(filters.ConfigServerId))
		}
		if len(filters.ConfigDeployStatus) > 0 {
			request = request.FilterConfigDeployStatus(utils.ProcessFilterStringSlice(filters.ConfigDeployStatus))
		}
		if len(filters.ConfigDeployType) > 0 {
			request = request.FilterConfigDeployType(utils.ProcessFilterStringSlice(filters.ConfigDeployType))
		}

		instances, meta, err := utils.FetchAllPages(request)
		if err != nil {
			return err
		}

		return utils.PrintAll(instances, meta, len(instances), &serverInstancePrintConfig)
	}

	logger.Get().Info().Msg("Listing server instances")

	request := client.ServerInstanceAPI.GetServerInstances(ctx).SortBy([]string{"id:ASC"})
	if len(filters.InfrastructureId) > 0 {
		request = request.FilterInfrastructureId(utils.ProcessFilterStringSlice(filters.InfrastructureId))
	}
	if len(filters.GroupId) > 0 {
		request = request.FilterGroupId(utils.ProcessFilterStringSlice(filters.GroupId))
	}
	if len(filters.ServerId) > 0 {
		request = request.FilterServerId(utils.ProcessFilterStringSlice(filters.ServerId))
	}
	if len(filters.ServiceStatus) > 0 {
		request = request.FilterServiceStatus(utils.ProcessFilterStringSlice(filters.ServiceStatus))
	}
	if len(filters.ConfigServerId) > 0 {
		request = request.FilterConfigServerId(utils.ProcessFilterStringSlice(filters.ConfigServerId))
	}
	if len(filters.ConfigDeployStatus) > 0 {
		request = request.FilterConfigDeployStatus(utils.ProcessFilterStringSlice(filters.ConfigDeployStatus))
	}
	if len(filters.ConfigDeployType) > 0 {
		request = request.FilterConfigDeployType(utils.ProcessFilterStringSlice(filters.ConfigDeployType))
	}

	instances, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(instances, meta, len(instances), &serverInstancePrintConfig)
}

func ServerInstanceGet(ctx context.Context, serverInstanceId string) error {
	logger.Get().Info().Msgf("Get server instance details for %s", serverInstanceId)

	serverInstanceIdNumerical, err := utils.GetInt64FromString(serverInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	serverInstanceInfo, httpRes, err := client.ServerInstanceAPI.GetServerInstance(ctx, serverInstanceIdNumerical).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(serverInstanceInfo, &serverInstancePrintConfig)
}

func ServerInstancePower(ctx context.Context, serverInstanceId string, action string) error {
	logger.Get().Info().Msgf("Setting power for server instance '%s' to '%s'", serverInstanceId, action)

	if err := validatePowerAction(action); err != nil {
		return err
	}

	instanceId, revision, err := getServerInstanceIdAndRevision(ctx, serverInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	powerSet := sdk.ServerInstancePowerSet{
		PowerCommand: action,
	}

	httpRes, err := client.ServerInstanceAPI.
		SetPowerToServerInstance(ctx, instanceId).
		ServerInstancePowerSet(powerSet).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Power command '%s' sent to server instance '%s'", action, serverInstanceId)
	fmt.Printf("Power command '%s' sent to server instance %s\n", action, serverInstanceId)

	return nil
}

func ServerInstancePowerStatus(ctx context.Context, serverInstanceId string) error {
	logger.Get().Info().Msgf("Getting power status for server instance '%s'", serverInstanceId)

	instanceId, err := getServerInstanceId(serverInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	result, httpRes, err := client.ServerInstanceAPI.
		GetPowerFromServerInstance(ctx, instanceId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	status := strings.Trim(strings.TrimSpace(result), "\"")
	if status != "" {
		fmt.Printf("Power status for server instance %s: %s\n", serverInstanceId, status)
	} else {
		fmt.Printf("Power status check initiated for server instance %s\n", serverInstanceId)
	}

	return nil
}

func ServerInstanceCredentials(ctx context.Context, serverInstanceId string) error {
	logger.Get().Info().Msgf("Getting credentials for server instance '%s'", serverInstanceId)

	instanceId, err := getServerInstanceId(serverInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	credentials, httpRes, err := client.ServerInstanceAPI.
		GetServerInstanceCredentials(ctx, instanceId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"Username": {
				Title: "Username",
				Order: 1,
			},
			"InitialPassword": {
				Title: "Password",
				Order: 2,
			},
			"PublicSshKey": {
				Title:    "SSH Public Key",
				MaxWidth: 60,
				Order:    3,
			},
		},
	})
}

func ServerInstanceReinstallOS(ctx context.Context, serverInstanceId string) error {
	logger.Get().Info().Msgf("Reinstalling OS for server instance '%s'", serverInstanceId)

	instanceId, revision, err := getServerInstanceIdAndRevision(ctx, serverInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	reinstall := sdk.ServerInstanceReinstallOS{
		PerformAtNextDeploy: true,
		ReinstallOS:         true,
	}

	httpRes, err := client.ServerInstanceAPI.
		ReinstallServerInstanceOS(ctx, instanceId).
		ServerInstanceReinstallOS(reinstall).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("OS reinstall scheduled for server instance '%s'", serverInstanceId)
	fmt.Printf("OS reinstall scheduled for server instance %s (will take effect at next deploy)\n", serverInstanceId)

	return nil
}

func ServerInstanceConfig(ctx context.Context, serverInstanceId string) error {
	logger.Get().Info().Msgf("Getting configuration for server instance '%s'", serverInstanceId)

	instanceId, err := getServerInstanceId(serverInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	config, httpRes, err := client.ServerInstanceAPI.
		GetServerInstanceConfig(ctx, instanceId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(config, &serverInstanceConfigPrintConfig)
}

// ServerInstanceConfigExample prints a populated sdk.ServerInstanceCreate that
// can be edited and passed back to 'server-instance create --config-source'.
func ServerInstanceConfigExample(ctx context.Context) error {
	example := sdk.ServerInstanceCreate{
		Label:        sdk.PtrString("my-server-instance"),
		GroupId:      sdk.PtrInt64(1),
		ServerTypeId: sdk.PtrInt64(1),
		Hostname:     sdk.PtrString("my-server-instance"),
		OsTemplateId: sdk.PtrInt64(1),
		Tags:         []string{"example"},
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(example, nil)
	}
	return formatter.PrintYamlResult(example)
}

func ServerInstanceCreate(ctx context.Context, infrastructureIdOrLabel string, create sdk.ServerInstanceCreate) error {
	logger.Get().Info().Msgf("Creating server instance in infrastructure '%s'", infrastructureIdOrLabel)

	infra, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	serverInstance, httpRes, err := client.ServerInstanceAPI.
		CreateServerInstance(ctx, int64(infra.Id)).
		ServerInstanceCreate(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(serverInstance, &serverInstancePrintConfig)
}

func ServerInstanceDelete(ctx context.Context, serverInstanceId string) error {
	logger.Get().Info().Msgf("Deleting server instance '%s'", serverInstanceId)

	instanceId, revision, err := getServerInstanceIdAndRevision(ctx, serverInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerInstanceAPI.
		DeleteServerInstance(ctx, instanceId).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Server instance '%s' deleted", serverInstanceId)
	return nil
}

// ServerInstanceConfigUpdate patches the pending configuration of a server
// instance. Writes below /config are guarded by the configuration revision, not
// by the instance revision, so the configuration is fetched first.
func ServerInstanceConfigUpdate(ctx context.Context, serverInstanceId string, config []byte) error {
	logger.Get().Info().Msgf("Updating configuration of server instance '%s'", serverInstanceId)

	var update sdk.ServerInstanceUpdate
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	instanceId, err := getServerInstanceId(serverInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	currentConfig, httpRes, err := client.ServerInstanceAPI.GetServerInstanceConfig(ctx, instanceId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	updatedConfig, httpRes, err := client.ServerInstanceAPI.
		UpdateServerInstanceConfig(ctx, instanceId).
		IfMatch(strconv.FormatInt(currentConfig.Revision, 10)).
		ServerInstanceUpdate(update).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(updatedConfig, &serverInstanceConfigPrintConfig)
}

func ServerInstanceMetaUpdate(ctx context.Context, serverInstanceId string, config []byte) error {
	logger.Get().Info().Msgf("Updating metadata of server instance '%s'", serverInstanceId)

	var meta sdk.GenericMeta
	if err := utils.UnmarshalContent(config, &meta); err != nil {
		return err
	}

	instanceId, err := getServerInstanceId(serverInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerInstanceAPI.
		UpdateServerInstanceMeta(ctx, instanceId).
		GenericMeta(meta).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Server instance '%s' metadata updated", serverInstanceId)
	return nil
}

// ServerInstanceReset resets the deployed server of a server instance
// (POST /server-instances/{id}/actions/reset). This is a different operation
// from 'power reset': the power command goes to /actions/power-set and only
// cycles the power of the server, while this call resets the deployed server
// immediately through the orchestration layer.
func ServerInstanceReset(ctx context.Context, serverInstanceId string) error {
	logger.Get().Info().Msgf("Resetting server instance '%s'", serverInstanceId)

	instanceId, revision, err := getServerInstanceIdAndRevision(ctx, serverInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerInstanceAPI.
		ResetServerInstance(ctx, instanceId).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Server instance '%s' reset", serverInstanceId)
	fmt.Printf("Reset requested for server instance %s\n", serverInstanceId)

	return nil
}

func ServerInstanceDrives(ctx context.Context, serverInstanceId string) error {
	logger.Get().Info().Msgf("Listing drives of server instance '%s'", serverInstanceId)

	instanceId, err := getServerInstanceId(serverInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	drives, httpRes, err := client.ServerInstanceAPI.GetServerInstanceDrives(ctx, instanceId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(drives, &serverInstanceDrivePrintConfig)
}

func ServerInstanceInterfaces(ctx context.Context, serverInstanceId string, filters ServerInstanceFilters) error {
	logger.Get().Info().Msgf("Listing interfaces of server instance '%s'", serverInstanceId)

	instanceId, err := getServerInstanceId(serverInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.ServerInstanceAPI.GetServerInstanceInterfaces(ctx, instanceId).SortBy([]string{"id:ASC"})
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

	return utils.PrintAll(interfaces, meta, len(interfaces), &serverInstanceInterfacePrintConfig)
}

func ServerInstanceInterfaceGet(ctx context.Context, serverInstanceId string, interfaceId string) error {
	logger.Get().Info().Msgf("Get interface '%s' of server instance '%s'", interfaceId, serverInstanceId)

	instanceId, interfaceIdNumerical, err := getServerInstanceAndInterfaceId(serverInstanceId, interfaceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	serverInstanceInterface, httpRes, err := client.ServerInstanceAPI.
		GetServerInstanceInterface(ctx, instanceId, interfaceIdNumerical).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(serverInstanceInterface, &serverInstanceInterfacePrintConfig)
}

// ServerInstanceInterfaceConfigUpdate patches the pending configuration of one
// server instance interface. The If-Match comes from the interface CONFIG
// revision, which is what guards writes below the /config path.
func ServerInstanceInterfaceConfigUpdate(ctx context.Context, serverInstanceId string, interfaceId string, config []byte) error {
	logger.Get().Info().Msgf("Updating configuration of interface '%s' of server instance '%s'", interfaceId, serverInstanceId)

	var update sdk.ServerInstanceInterfaceUpdate
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	instanceId, interfaceIdNumerical, err := getServerInstanceAndInterfaceId(serverInstanceId, interfaceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	currentInterface, httpRes, err := client.ServerInstanceAPI.
		GetServerInstanceInterface(ctx, instanceId, interfaceIdNumerical).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	revision := currentInterface.Revision
	if currentInterface.Config != nil {
		revision = currentInterface.Config.Revision
	}

	updatedConfig, httpRes, err := client.ServerInstanceAPI.
		UpdateServerInstanceInterfaceConfig(ctx, instanceId, interfaceIdNumerical).
		IfMatch(strconv.FormatInt(revision, 10)).
		ServerInstanceInterfaceUpdate(update).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(updatedConfig, &serverInstanceInterfaceConfigPrintConfig)
}

// ServerInstanceOSInstallationData returns the values the OS installer is
// rendered with.
//
// The body is read raw instead of through the typed SDK call: the generated
// ServerInstanceContextOSInstallationData types 'network' as an array and marks
// 'server.serverId' as required, while the API returns 'network' as an object
// and omits the server block for instances with no server deployed. Both make
// the typed call fail with
// "json: cannot unmarshal object into Go struct field
// _ServerInstanceContextOSInstallationData.network of type
// []sdk.InstanceNetworkVariables" and
// "no value given for required property serverId" respectively.
func ServerInstanceOSInstallationData(ctx context.Context, serverInstanceId string, usage string) error {
	logger.Get().Info().Msgf("Get OS installation data of server instance '%s'", serverInstanceId)

	instanceId, err := getServerInstanceId(serverInstanceId)
	if err != nil {
		return err
	}

	body, err := api.RawJSONRequest(ctx, http.MethodGet,
		withUsage(fmt.Sprintf("/api/v2/server-instances/%d/os-installation-data", instanceId), usage),
		nil, nil)
	if err != nil {
		return err
	}

	return printContextObject(body)
}

// ServerInstanceVariables returns the variables available to extensions and
// templates. It is read raw for the same reason as
// ServerInstanceOSInstallationData.
func ServerInstanceVariables(ctx context.Context, serverInstanceId string, usage string) error {
	logger.Get().Info().Msgf("Get variables of server instance '%s'", serverInstanceId)

	instanceId, err := getServerInstanceId(serverInstanceId)
	if err != nil {
		return err
	}

	body, err := api.RawJSONRequest(ctx, http.MethodGet,
		withUsage(fmt.Sprintf("/api/v2/server-instances/%d/variables", instanceId), usage),
		nil, nil)
	if err != nil {
		return err
	}

	return printContextObject(body)
}

// printContextObject renders one of the deeply nested context objects. A table
// of such an object is unreadable, so every non native format falls back to
// YAML.
func printContextObject(body []byte) error {
	object, err := utils.DecodeRawObject(body)
	if err != nil {
		return err
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(object, nil)
	}

	return formatter.PrintYamlResult(object)
}

// withUsage appends the optional usage query parameter to path.
func withUsage(path string, usage string) string {
	if usage == "" {
		return path
	}
	return path + "?usage=" + url.QueryEscape(usage)
}

func ServerInstanceStatistics(ctx context.Context) error {
	logger.Get().Info().Msg("Get server instance statistics")

	client := api.GetApiClient(ctx)

	statistics, httpRes, err := client.ServerInstanceAPI.GetServerInstanceStatistics(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(statistics, &serverInstanceStatisticsPrintConfig)
}

// ServerInstancePowerStatusBatch reads the power status of several server
// instances of one infrastructure in a single call. The endpoint is a POST
// that carries the instance ids in the body, so the body is always sent.
func ServerInstancePowerStatusBatch(ctx context.Context, infrastructureIdOrLabel string, serverInstanceIds []string) error {
	logger.Get().Info().Msgf("Getting power status of server instances %v of infrastructure '%s'", serverInstanceIds, infrastructureIdOrLabel)

	infra, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	if err := validateServerInstanceIds(serverInstanceIds); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	statuses, httpRes, err := client.ServerInstanceAPI.
		GetPowerStatusBatch(ctx, int64(infra.Id)).
		RequestBody(serverInstanceIds).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(statuses, nil)
}

// ServerInstancePowerSetBatch sets the power state of several server instances
// of one infrastructure in a single call.
func ServerInstancePowerSetBatch(ctx context.Context, infrastructureIdOrLabel string, action string, serverInstanceIds []string) error {
	logger.Get().Info().Msgf("Setting power of server instances %v of infrastructure '%s' to '%s'", serverInstanceIds, infrastructureIdOrLabel, action)

	if err := validatePowerAction(action); err != nil {
		return err
	}

	infra, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	if err := validateServerInstanceIds(serverInstanceIds); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerInstanceAPI.
		SetPowerStatusBatch(ctx, int64(infra.Id)).
		InstancesSetPowerState(sdk.InstancesSetPowerState{
			Instances:    serverInstanceIds,
			PowerCommand: action,
		}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Power command '%s' sent to server instances %v", action, serverInstanceIds)
	fmt.Printf("Power command '%s' sent to server instances %s\n", action, strings.Join(serverInstanceIds, ", "))

	return nil
}

// validatePowerAction rejects power commands the API does not accept before
// they reach the wire.
func validatePowerAction(action string) error {
	validActions := map[string]bool{
		"on":    true,
		"off":   true,
		"reset": true,
		"soft":  true,
	}

	if !validActions[action] {
		return fmt.Errorf("invalid power action: '%s'. Valid actions are: on, off, reset, soft", action)
	}

	return nil
}

// validateServerInstanceIds makes sure every batch member is a numeric id; the
// API expects them as strings but rejects anything that is not a number.
func validateServerInstanceIds(serverInstanceIds []string) error {
	if len(serverInstanceIds) == 0 {
		return fmt.Errorf("at least one server instance ID is required")
	}

	for _, serverInstanceId := range serverInstanceIds {
		if _, err := getServerInstanceId(serverInstanceId); err != nil {
			return err
		}
	}

	return nil
}

func getServerInstanceId(serverInstanceId string) (int64, error) {
	id, err := utils.GetInt64FromString(serverInstanceId)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func getServerInstanceAndInterfaceId(serverInstanceId string, interfaceId string) (int64, int64, error) {
	instanceId, err := getServerInstanceId(serverInstanceId)
	if err != nil {
		return 0, 0, err
	}

	interfaceIdNumerical, err := utils.GetInt64FromString(interfaceId)
	if err != nil {
		return 0, 0, err
	}

	return instanceId, interfaceIdNumerical, nil
}

func getServerInstanceIdAndRevision(ctx context.Context, serverInstanceId string) (int64, string, error) {
	instanceId, err := getServerInstanceId(serverInstanceId)
	if err != nil {
		return 0, "", err
	}

	client := api.GetApiClient(ctx)

	instance, httpRes, err := client.ServerInstanceAPI.GetServerInstance(ctx, instanceId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}

	return instanceId, strconv.Itoa(int(instance.Revision)), nil
}
