package container_instance

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/internal/infrastructure"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var containerInstancePrintConfig = formatter.PrintConfig{
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
		"GroupId": {
			Title: "Group",
			Order: 4,
		},
		"TypeId": {
			Title: "Type",
			Order: 5,
		},
		"ContainerId": {
			Title: "Container",
			Order: 6,
		},
		"CpuCores": {
			Title: "Cores",
			Order: 7,
		},
		"RamGB": {
			Title: "RAM GB",
			Order: 8,
		},
		"DiskSizeGB": {
			Title: "Disk GB",
			Order: 9,
		},
		"ServiceStatus": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       10,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       11,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       12,
		},
	},
}

var containerInstanceConfigPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Revision": {
			Title: "Revision",
			Order: 1,
		},
		"Label": {
			MaxWidth: 30,
			Order:    2,
		},
		"TypeId": {
			Title: "Type",
			Order: 3,
		},
		"ContainerId": {
			Title: "Container",
			Order: 4,
		},
		"CpuCores": {
			Title: "Cores",
			Order: 5,
		},
		"RamGB": {
			Title: "RAM GB",
			Order: 6,
		},
		"DiskSizeGB": {
			Title: "Disk GB",
			Order: 7,
		},
		"DeployType": {
			Title: "Deploy Type",
			Order: 8,
		},
		"DeployStatus": {
			Title:       "Deploy Status",
			Transformer: formatter.FormatStatusValue,
			Order:       9,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       10,
		},
	},
}

var containerInstanceCredentialsPrintConfig = formatter.PrintConfig{
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
}

// PrintConfig exposes the container instance table layout so other packages
// that list container instances render them identically.
func PrintConfig() *formatter.PrintConfig {
	return &containerInstancePrintConfig
}

func ContainerInstanceList(ctx context.Context, infrastructureIdOrLabel string) error {
	logger.Get().Info().Msgf("List all container instances for infrastructure '%s'", infrastructureIdOrLabel)

	infraId, err := resolveInfrastructureId(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.ContainerInstanceAPI.
		GetInfrastructureContainerInstances(ctx, infraId).
		SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &containerInstancePrintConfig)
}

func ContainerInstanceGet(ctx context.Context, infrastructureIdOrLabel string, containerInstanceId string) error {
	logger.Get().Info().Msgf("Get container instance '%s' in infrastructure '%s'", containerInstanceId, infrastructureIdOrLabel)

	infraId, instanceId, err := resolveIds(ctx, infrastructureIdOrLabel, containerInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	containerInstance, httpRes, err := client.ContainerInstanceAPI.
		GetInfrastructureContainerInstance(ctx, infraId, instanceId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(containerInstance, &containerInstancePrintConfig)
}

func ContainerInstanceConfigExample(ctx context.Context) error {
	example := sdk.CreateContainerInstance{
		TypeId:     1,
		GroupId:    1,
		DiskSizeGB: sdk.PtrFloat32(40),
		Tags:       []string{"tag1", "tag2"},
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(example, nil)
	}

	return formatter.PrintYamlResult(example)
}

// ContainerInstanceCreate creates a container instance from a complete
// configuration payload (JSON or YAML).
func ContainerInstanceCreate(ctx context.Context, infrastructureIdOrLabel string, config []byte) error {
	var create sdk.CreateContainerInstance
	if err := utils.UnmarshalContent(config, &create); err != nil {
		return err
	}

	return containerInstanceCreate(ctx, infrastructureIdOrLabel, create)
}

// ContainerInstanceCreateFromFlags creates a container instance from individual
// command line flags.
func ContainerInstanceCreateFromFlags(ctx context.Context, infrastructureIdOrLabel string, containerTypeId string, groupId string, diskSizeGB string, tags []string) error {
	typeIdNumerical, err := utils.GetInt64FromString(containerTypeId)
	if err != nil {
		return err
	}

	groupIdNumerical, err := utils.GetInt64FromString(groupId)
	if err != nil {
		return err
	}

	create := sdk.CreateContainerInstance{
		TypeId:  typeIdNumerical,
		GroupId: groupIdNumerical,
		Tags:    tags,
	}

	if diskSizeGB != "" {
		diskSize, err := utils.GetFloat32FromString(diskSizeGB)
		if err != nil {
			return err
		}
		create.DiskSizeGB = &diskSize
	}

	return containerInstanceCreate(ctx, infrastructureIdOrLabel, create)
}

func containerInstanceCreate(ctx context.Context, infrastructureIdOrLabel string, create sdk.CreateContainerInstance) error {
	logger.Get().Info().Msgf("Create container instance in infrastructure '%s'", infrastructureIdOrLabel)

	infraId, err := resolveInfrastructureId(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	containerInstance, httpRes, err := client.ContainerInstanceAPI.
		CreateContainerInstance(ctx, infraId).
		CreateContainerInstance(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(containerInstance, &containerInstancePrintConfig)
}

func ContainerInstanceDelete(ctx context.Context, infrastructureIdOrLabel string, containerInstanceId string) error {
	logger.Get().Info().Msgf("Delete container instance '%s' from infrastructure '%s'", containerInstanceId, infrastructureIdOrLabel)

	infraId, instanceId, revision, err := resolveIdsAndRevision(ctx, infrastructureIdOrLabel, containerInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ContainerInstanceAPI.
		DeleteContainerInstance(ctx, infraId, instanceId).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Container instance '%s' deleted", containerInstanceId)
	return nil
}

func ContainerInstanceGetConfig(ctx context.Context, infrastructureIdOrLabel string, containerInstanceId string) error {
	logger.Get().Info().Msgf("Get configuration of container instance '%s' in infrastructure '%s'", containerInstanceId, infrastructureIdOrLabel)

	infraId, instanceId, err := resolveIds(ctx, infrastructureIdOrLabel, containerInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	config, httpRes, err := client.ContainerInstanceAPI.
		GetContainerInstanceConfigInfo(ctx, infraId, instanceId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(config, &containerInstanceConfigPrintConfig)
}

func ContainerInstanceUpdateConfig(ctx context.Context, infrastructureIdOrLabel string, containerInstanceId string, config []byte) error {
	logger.Get().Info().Msgf("Update configuration of container instance '%s' in infrastructure '%s'", containerInstanceId, infrastructureIdOrLabel)

	var update sdk.UpdateContainerInstance
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	infraId, instanceId, err := resolveIds(ctx, infrastructureIdOrLabel, containerInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// The config sub-resource is guarded by its own revision, not by the
	// container instance revision.
	currentConfig, httpRes, err := client.ContainerInstanceAPI.
		GetContainerInstanceConfigInfo(ctx, infraId, instanceId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	updatedConfig, httpRes, err := client.ContainerInstanceAPI.
		UpdateContainerInstanceConfig(ctx, infraId, instanceId).
		UpdateContainerInstance(update).
		IfMatch(strconv.FormatInt(currentConfig.Revision, 10)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(updatedConfig, nil)
}

func ContainerInstanceUpdateMeta(ctx context.Context, infrastructureIdOrLabel string, containerInstanceId string, config []byte) error {
	logger.Get().Info().Msgf("Update metadata of container instance '%s' in infrastructure '%s'", containerInstanceId, infrastructureIdOrLabel)

	var update sdk.UpdateContainerInstanceMeta
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	infraId, instanceId, err := resolveIds(ctx, infrastructureIdOrLabel, containerInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	containerInstance, httpRes, err := client.ContainerInstanceAPI.
		PatchContainerInstanceMeta(ctx, infraId, instanceId).
		UpdateContainerInstanceMeta(update).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(containerInstance, &containerInstancePrintConfig)
}

func ContainerInstanceCredentials(ctx context.Context, infrastructureIdOrLabel string, containerInstanceId string) error {
	logger.Get().Info().Msgf("Get credentials of container instance '%s' in infrastructure '%s'", containerInstanceId, infrastructureIdOrLabel)

	infraId, instanceId, err := resolveIds(ctx, infrastructureIdOrLabel, containerInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	credentials, httpRes, err := client.ContainerInstanceAPI.
		GetContainerInstanceCredentials(ctx, infraId, instanceId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, &containerInstanceCredentialsPrintConfig)
}

func ContainerInstanceVariables(ctx context.Context, infrastructureIdOrLabel string, containerInstanceId string, usage string) error {
	logger.Get().Info().Msgf("Get variables of container instance '%s' in infrastructure '%s'", containerInstanceId, infrastructureIdOrLabel)

	infraId, instanceId, err := resolveIds(ctx, infrastructureIdOrLabel, containerInstanceId)
	if err != nil {
		return err
	}

	values := url.Values{}
	if err := addUsageValue(values, usage); err != nil {
		return err
	}

	// sdk.ContainerInstanceContextVariables requires a non-null `server` object,
	// but container instances answer with "server": null, so the strict decoder
	// fails with "no value given for required property serverId". Read the body
	// raw instead.
	body, err := api.RawJSONRequest(ctx, http.MethodGet,
		variablesPath(infraId, instanceId, "variables", values), nil, nil)
	if err != nil {
		return err
	}

	return utils.PrintRawObject(body, nil)
}

func ContainerInstanceOSInstallationData(ctx context.Context, infrastructureIdOrLabel string, containerInstanceId string, usage string, removeEmpty bool) error {
	logger.Get().Info().Msgf("Get OS installation data of container instance '%s' in infrastructure '%s'", containerInstanceId, infrastructureIdOrLabel)

	infraId, instanceId, err := resolveIds(ctx, infrastructureIdOrLabel, containerInstanceId)
	if err != nil {
		return err
	}

	values := url.Values{}
	if err := addUsageValue(values, usage); err != nil {
		return err
	}
	if removeEmpty {
		values.Set("removeEmpty", "1")
	}

	// Same schema drift as the variables endpoint: `server` comes back null for
	// container instances while the SDK model requires it.
	body, err := api.RawJSONRequest(ctx, http.MethodGet,
		variablesPath(infraId, instanceId, "os-installation-data", values), nil, nil)
	if err != nil {
		return err
	}

	return utils.PrintRawObject(body, nil)
}

// addUsageValue validates usage against the SDK enum and adds it to values.
func addUsageValue(values url.Values, usage string) error {
	if usage == "" {
		return nil
	}

	usageType, err := sdk.NewVariableUsageTypeFromValue(usage)
	if err != nil {
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	values.Set("usage", string(*usageType))
	return nil
}

func variablesPath(infraId int64, instanceId int64, resource string, values url.Values) string {
	path := fmt.Sprintf("/api/v2/infrastructures/%d/container-instances/%d/%s", infraId, instanceId, resource)
	if encoded := values.Encode(); encoded != "" {
		path += "?" + encoded
	}
	return path
}

func ContainerInstancePowerStatus(ctx context.Context, infrastructureIdOrLabel string, containerInstanceId string) error {
	logger.Get().Info().Msgf("Get power status of container instance '%s' in infrastructure '%s'", containerInstanceId, infrastructureIdOrLabel)

	infraId, instanceId, err := resolveIds(ctx, infrastructureIdOrLabel, containerInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	powerStatus, httpRes, err := client.ContainerInstanceAPI.
		GetContainerInstancePowerStatus(ctx, infraId, instanceId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(powerStatus, nil)
}

// ContainerInstancePowerControl runs one of the 'start', 'shutdown' or 'reboot'
// actions against a container instance. These endpoints answer 204 with no body.
func ContainerInstancePowerControl(ctx context.Context, infrastructureIdOrLabel string, containerInstanceId string, action string) error {
	logger.Get().Info().Msgf("Performing '%s' action on container instance '%s' in infrastructure '%s'", action, containerInstanceId, infrastructureIdOrLabel)

	infraId, instanceId, err := resolveIds(ctx, infrastructureIdOrLabel, containerInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	var httpRes *http.Response

	switch action {
	case "start":
		httpRes, err = client.ContainerInstanceAPI.StartContainerInstance(ctx, infraId, instanceId).Execute()
	case "shutdown":
		httpRes, err = client.ContainerInstanceAPI.ShutdownContainerInstance(ctx, infraId, instanceId).Execute()
	case "reboot":
		httpRes, err = client.ContainerInstanceAPI.RebootContainerInstance(ctx, infraId, instanceId).Execute()
	default:
		err := fmt.Errorf("unsupported power action: '%s'. Use 'start', 'shutdown' or 'reboot'", action)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Container instance '%s' power action '%s' successful", containerInstanceId, action)
	return nil
}

func ContainerInstanceApplyType(ctx context.Context, infrastructureIdOrLabel string, containerInstanceId string, containerTypeId string) error {
	logger.Get().Info().Msgf("Apply container type '%s' on container instance '%s' in infrastructure '%s'", containerTypeId, containerInstanceId, infrastructureIdOrLabel)

	containerTypeIdNumerical, err := utils.GetInt64FromString(containerTypeId)
	if err != nil {
		return err
	}

	infraId, instanceId, revision, err := resolveIdsAndRevision(ctx, infrastructureIdOrLabel, containerInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	containerInstance, httpRes, err := client.ContainerInstanceAPI.
		ApplyContainerTypeOnContainerInstance(ctx, infraId, instanceId, containerTypeIdNumerical).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(containerInstance, &containerInstancePrintConfig)
}

func resolveInfrastructureId(ctx context.Context, infrastructureIdOrLabel string) (int64, error) {
	infra, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return 0, err
	}

	return int64(infra.Id), nil
}

func resolveIds(ctx context.Context, infrastructureIdOrLabel string, containerInstanceId string) (int64, int64, error) {
	instanceId, err := GetContainerInstanceId(containerInstanceId)
	if err != nil {
		return 0, 0, err
	}

	infraId, err := resolveInfrastructureId(ctx, infrastructureIdOrLabel)
	if err != nil {
		return 0, 0, err
	}

	return infraId, instanceId, nil
}

func resolveIdsAndRevision(ctx context.Context, infrastructureIdOrLabel string, containerInstanceId string) (int64, int64, string, error) {
	infraId, instanceId, err := resolveIds(ctx, infrastructureIdOrLabel, containerInstanceId)
	if err != nil {
		return 0, 0, "", err
	}

	client := api.GetApiClient(ctx)

	containerInstance, httpRes, err := client.ContainerInstanceAPI.
		GetInfrastructureContainerInstance(ctx, infraId, instanceId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, 0, "", err
	}

	return infraId, instanceId, strconv.FormatInt(containerInstance.Revision, 10), nil
}

func GetContainerInstanceId(containerInstanceId string) (int64, error) {
	containerInstanceIdNumerical, err := utils.GetInt64FromString(containerInstanceId)
	if err != nil {
		err := fmt.Errorf("invalid container instance ID: '%s'", containerInstanceId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return containerInstanceIdNumerical, nil
}
