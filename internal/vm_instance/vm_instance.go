package vm_instance

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

var vmInstancePrintConfig = formatter.PrintConfig{
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
		"TypeId": {
			Title: "Type ID",
			Order: 6,
		},
		"DiskSizeGB": {
			Title: "Disk GB",
			Order: 7,
		},
		"RamGB": {
			Title: "RAM GB",
			Order: 8,
		},
		"CpuCores": {
			Title: "CPU Cores",
			Order: 9,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       10,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       11,
		},
	},
}

var vmInstanceConfigPrintConfig = formatter.PrintConfig{
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
			Title: "Type ID",
			Order: 3,
		},
		"DiskSizeGB": {
			Title: "Disk GB",
			Order: 4,
		},
		"RamGB": {
			Title: "RAM GB",
			Order: 5,
		},
		"CpuCores": {
			Title: "CPU Cores",
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
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       9,
		},
	},
}

var vmInstanceCredentialsPrintConfig = formatter.PrintConfig{
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

func VMInstanceGet(ctx context.Context, infrastructureIdOrLabel string, vmInstanceId string) error {
	logger.Get().Info().Msgf("Get VM instance '%s' in infrastructure '%s'", vmInstanceId, infrastructureIdOrLabel)

	infraId, instanceId, err := resolveIds(ctx, infrastructureIdOrLabel, vmInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmInstanceInfo, httpRes, err := client.VMInstanceAPI.
		GetInfrastructureVMInstance(ctx, infraId, instanceId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmInstanceInfo, &vmInstancePrintConfig)
}

func VMInstanceList(ctx context.Context, infrastructureIdOrLabel string) error {
	logger.Get().Info().Msgf("List all VM instances for infrastructure '%s'", infrastructureIdOrLabel)

	infraId, err := resolveInfrastructureId(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.VMInstanceAPI.
		GetInfrastructureVMInstances(ctx, infraId).
		SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &vmInstancePrintConfig)
}

func VMInstanceConfigExample(ctx context.Context) error {
	example := sdk.CreateVMInstance{
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

// VMInstanceCreate creates a VM instance from a complete configuration payload
// (JSON or YAML).
func VMInstanceCreate(ctx context.Context, infrastructureIdOrLabel string, config []byte) error {
	var create sdk.CreateVMInstance
	if err := utils.UnmarshalContent(config, &create); err != nil {
		return err
	}

	return vmInstanceCreate(ctx, infrastructureIdOrLabel, create)
}

// VMInstanceCreateFromFlags creates a VM instance from individual command line
// flags.
func VMInstanceCreateFromFlags(ctx context.Context, infrastructureIdOrLabel string, vmTypeId string, groupId string, diskSizeGB string, tags []string) error {
	typeIdNumerical, err := utils.GetInt64FromString(vmTypeId)
	if err != nil {
		return err
	}

	groupIdNumerical, err := utils.GetInt64FromString(groupId)
	if err != nil {
		return err
	}

	create := sdk.CreateVMInstance{
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

	return vmInstanceCreate(ctx, infrastructureIdOrLabel, create)
}

func vmInstanceCreate(ctx context.Context, infrastructureIdOrLabel string, create sdk.CreateVMInstance) error {
	logger.Get().Info().Msgf("Create VM instance in infrastructure '%s'", infrastructureIdOrLabel)

	infraId, err := resolveInfrastructureId(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmInstanceInfo, httpRes, err := client.VMInstanceAPI.
		CreateVMInstance(ctx, infraId).
		CreateVMInstance(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmInstanceInfo, &vmInstancePrintConfig)
}

func VMInstanceDelete(ctx context.Context, infrastructureIdOrLabel string, vmInstanceId string) error {
	logger.Get().Info().Msgf("Delete VM instance '%s' from infrastructure '%s'", vmInstanceId, infrastructureIdOrLabel)

	infraId, instanceId, revision, err := resolveIdsAndRevision(ctx, infrastructureIdOrLabel, vmInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.VMInstanceAPI.
		DeleteVMInstance(ctx, infraId, instanceId).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM instance '%s' deleted", vmInstanceId)
	return nil
}

func VMInstanceGetConfig(ctx context.Context, infrastructureIdOrLabel string, vmInstanceId string) error {
	logger.Get().Info().Msgf("Get configuration of VM instance '%s' in infrastructure '%s'", vmInstanceId, infrastructureIdOrLabel)

	infraId, instanceId, err := resolveIds(ctx, infrastructureIdOrLabel, vmInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmInstanceConfig, httpRes, err := client.VMInstanceAPI.
		GetVMInstanceConfigInfo(ctx, infraId, instanceId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmInstanceConfig, &vmInstanceConfigPrintConfig)
}

func VMInstanceUpdateConfig(ctx context.Context, infrastructureIdOrLabel string, vmInstanceId string, config []byte) error {
	logger.Get().Info().Msgf("Update configuration of VM instance '%s' in infrastructure '%s'", vmInstanceId, infrastructureIdOrLabel)

	var update sdk.UpdateVMInstance
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	infraId, instanceId, err := resolveIds(ctx, infrastructureIdOrLabel, vmInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// The config sub-resource is guarded by its own revision, not by the VM
	// instance revision.
	currentConfig, httpRes, err := client.VMInstanceAPI.
		GetVMInstanceConfigInfo(ctx, infraId, instanceId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	updatedConfig, httpRes, err := client.VMInstanceAPI.
		UpdateVMInstanceConfig(ctx, infraId, instanceId).
		UpdateVMInstance(update).
		IfMatch(strconv.FormatInt(currentConfig.Revision, 10)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(updatedConfig, nil)
}

func VMInstanceUpdateMeta(ctx context.Context, infrastructureIdOrLabel string, vmInstanceId string, config []byte) error {
	logger.Get().Info().Msgf("Update metadata of VM instance '%s' in infrastructure '%s'", vmInstanceId, infrastructureIdOrLabel)

	var update sdk.UpdateVMInstanceMeta
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	infraId, instanceId, err := resolveIds(ctx, infrastructureIdOrLabel, vmInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmInstanceInfo, httpRes, err := client.VMInstanceAPI.
		PatchVMInstanceMeta(ctx, infraId, instanceId).
		UpdateVMInstanceMeta(update).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmInstanceInfo, &vmInstancePrintConfig)
}

func VMInstanceApplyType(ctx context.Context, infrastructureIdOrLabel string, vmInstanceId string, vmTypeId string) error {
	logger.Get().Info().Msgf("Apply VM type '%s' on VM instance '%s' in infrastructure '%s'", vmTypeId, vmInstanceId, infrastructureIdOrLabel)

	vmTypeIdNumerical, err := utils.GetInt64FromString(vmTypeId)
	if err != nil {
		return err
	}

	infraId, instanceId, revision, err := resolveIdsAndRevision(ctx, infrastructureIdOrLabel, vmInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmInstanceInfo, httpRes, err := client.VMInstanceAPI.
		ApplyVMTypeOnVMInstance(ctx, infraId, instanceId, vmTypeIdNumerical).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmInstanceInfo, &vmInstancePrintConfig)
}

// VMInstanceVariables returns the variables available to extensions and
// templates for a VM instance.
//
// The body is read raw instead of through the typed SDK call: the generated
// sdk.VmInstanceContextVariables requires a non-null 'server' block and types
// 'network' as an array, while the API answers with "server": null and a
// 'network' object. The typed call fails on live data with
// "no value given for required property serverId" (and, once a server is
// attached, with "json: cannot unmarshal object into Go struct field
// _VmInstanceContextVariables.network of type []sdk.InstanceNetworkVariables").
func VMInstanceVariables(ctx context.Context, infrastructureIdOrLabel string, vmInstanceId string, usage string) error {
	logger.Get().Info().Msgf("Get variables of VM instance '%s' in infrastructure '%s'", vmInstanceId, infrastructureIdOrLabel)

	infraId, instanceId, err := resolveIds(ctx, infrastructureIdOrLabel, vmInstanceId)
	if err != nil {
		return err
	}

	values := url.Values{}
	if err := addUsageValue(values, usage); err != nil {
		return err
	}

	body, err := api.RawJSONRequest(ctx, http.MethodGet,
		vmInstanceContextPath(infraId, instanceId, "variables", values), nil, nil)
	if err != nil {
		return err
	}

	return printContextObject(body)
}

// VMInstanceOSInstallationData returns the values the OS installer of a VM
// instance is rendered with. It is read raw for the same reason as
// VMInstanceVariables.
//
// The removeEmpty query parameter is declared by the SDK but is not routed by
// every API build: QA01 answers 404 "Cannot GET
// /vm-instances/{id}/os-installation-data&removeEmpty=1" for it.
func VMInstanceOSInstallationData(ctx context.Context, infrastructureIdOrLabel string, vmInstanceId string, usage string, removeEmpty bool) error {
	logger.Get().Info().Msgf("Get OS installation data of VM instance '%s' in infrastructure '%s'", vmInstanceId, infrastructureIdOrLabel)

	infraId, instanceId, err := resolveIds(ctx, infrastructureIdOrLabel, vmInstanceId)
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

	body, err := api.RawJSONRequest(ctx, http.MethodGet,
		vmInstanceContextPath(infraId, instanceId, "os-installation-data", values), nil, nil)
	if err != nil {
		return err
	}

	return printContextObject(body)
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

func vmInstanceContextPath(infraId int64, instanceId int64, resource string, values url.Values) string {
	path := fmt.Sprintf("/api/v2/infrastructures/%d/vm-instances/%d/%s", infraId, instanceId, resource)
	if encoded := values.Encode(); encoded != "" {
		path += "?" + encoded
	}
	return path
}

func VMInstancePowerControl(ctx context.Context, infrastructureIdOrLabel string, vmInstanceId string, action string) error {
	logger.Get().Info().Msgf("Performing '%s' action on VM instance '%s' in infrastructure '%s'", action, vmInstanceId, infrastructureIdOrLabel)

	infraId, instanceId, err := resolveIds(ctx, infrastructureIdOrLabel, vmInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)
	var httpRes *http.Response

	switch action {
	case "start":
		httpRes, err = client.VMInstanceAPI.StartVMInstance(ctx, infraId, instanceId).Execute()
	case "shutdown":
		httpRes, err = client.VMInstanceAPI.ShutdownVMInstance(ctx, infraId, instanceId).Execute()
	case "reboot":
		httpRes, err = client.VMInstanceAPI.RebootVMInstance(ctx, infraId, instanceId).Execute()
	default:
		err := fmt.Errorf("unsupported power action: '%s'. Use 'start', 'shutdown' or 'reboot'", action)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM instance '%s' power action '%s' successful", vmInstanceId, action)
	return nil
}

func VMInstanceGetPowerStatus(ctx context.Context, infrastructureIdOrLabel string, vmInstanceId string) error {
	logger.Get().Info().Msgf("Get power status of VM instance '%s' in infrastructure '%s'", vmInstanceId, infrastructureIdOrLabel)

	infraId, instanceId, err := resolveIds(ctx, infrastructureIdOrLabel, vmInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	powerStatus, httpRes, err := client.VMInstanceAPI.
		GetVMInstancePowerStatus(ctx, infraId, instanceId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(powerStatus, nil)
}

func VMInstanceGetCredentials(ctx context.Context, infrastructureIdOrLabel string, vmInstanceId string) error {
	logger.Get().Info().Msgf("Get credentials of VM instance '%s' in infrastructure '%s'", vmInstanceId, infrastructureIdOrLabel)

	infraId, instanceId, err := resolveIds(ctx, infrastructureIdOrLabel, vmInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	credentials, httpRes, err := client.VMInstanceAPI.
		GetVMInstanceCredentials(ctx, infraId, instanceId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, &vmInstanceCredentialsPrintConfig)
}

func resolveInfrastructureId(ctx context.Context, infrastructureIdOrLabel string) (int64, error) {
	infra, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return 0, err
	}

	return int64(infra.Id), nil
}

func resolveIds(ctx context.Context, infrastructureIdOrLabel string, vmInstanceId string) (int64, int64, error) {
	instanceId, err := GetVMInstanceId(vmInstanceId)
	if err != nil {
		return 0, 0, err
	}

	infraId, err := resolveInfrastructureId(ctx, infrastructureIdOrLabel)
	if err != nil {
		return 0, 0, err
	}

	return infraId, instanceId, nil
}

func resolveIdsAndRevision(ctx context.Context, infrastructureIdOrLabel string, vmInstanceId string) (int64, int64, string, error) {
	infraId, instanceId, err := resolveIds(ctx, infrastructureIdOrLabel, vmInstanceId)
	if err != nil {
		return 0, 0, "", err
	}

	client := api.GetApiClient(ctx)

	vmInstanceInfo, httpRes, err := client.VMInstanceAPI.
		GetInfrastructureVMInstance(ctx, infraId, instanceId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, 0, "", err
	}

	return infraId, instanceId, strconv.FormatInt(vmInstanceInfo.Revision, 10), nil
}

func GetVMInstanceId(vmInstanceId string) (int64, error) {
	vmInstanceIdNumerical, err := utils.GetInt64FromString(vmInstanceId)
	if err != nil {
		err := fmt.Errorf("invalid VM instance ID: '%s'", vmInstanceId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return vmInstanceIdNumerical, nil
}
