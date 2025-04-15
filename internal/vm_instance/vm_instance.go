package vm_instance

import (
	"context"
	"fmt"
	"net/http"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
)

var vmInstancePrintConfig = formatter.PrintConfig{
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

func VMInstanceGet(ctx context.Context, infrastructureId string, vmInstanceId string) error {
	logger.Get().Info().Msgf("Get VM instance details for %s in infrastructure %s", vmInstanceId, infrastructureId)

	infraIdNumerical, err := utils.GetFloat32FromString(infrastructureId)
	if err != nil {
		return err
	}

	vmInstanceIdNumerical, err := utils.GetFloat32FromString(vmInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmInstanceInfo, httpRes, err := client.VMInstanceAPI.GetInfrastructureVMInstance(
		ctx, infraIdNumerical, vmInstanceIdNumerical).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmInstanceInfo, &vmInstancePrintConfig)
}

func VMInstanceList(ctx context.Context, infrastructureId string) error {
	logger.Get().Info().Msgf("List all VM instances for infrastructure %s", infrastructureId)

	infraIdNumerical, err := utils.GetFloat32FromString(infrastructureId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmInstancesList, httpRes, err := client.VMInstanceAPI.GetInfrastructureVMInstances(
		ctx, infraIdNumerical).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmInstancesList.Data, &vmInstancePrintConfig)
}

func VMInstanceGetConfig(ctx context.Context, infrastructureId string, vmInstanceId string) error {
	logger.Get().Info().Msgf("Get VM instance configuration for %s in infrastructure %s", vmInstanceId, infrastructureId)

	infraIdNumerical, err := utils.GetFloat32FromString(infrastructureId)
	if err != nil {
		return err
	}

	vmInstanceIdNumerical, err := utils.GetFloat32FromString(vmInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmInstanceConfig, httpRes, err := client.VMInstanceAPI.GetVMInstanceConfigInfo(
		ctx, infraIdNumerical, vmInstanceIdNumerical).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmInstanceConfig, nil)
}

func VMInstancePowerControl(ctx context.Context, infrastructureId string, vmInstanceId string, action string) error {
	logger.Get().Info().Msgf("Performing %s action on VM instance %s in infrastructure %s", action, vmInstanceId, infrastructureId)

	infraIdNumerical, err := utils.GetFloat32FromString(infrastructureId)
	if err != nil {
		return err
	}

	vmInstanceIdNumerical, err := utils.GetFloat32FromString(vmInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)
	var httpRes *http.Response

	switch action {
	case "start":
		httpRes, err = client.VMInstanceAPI.StartVMInstance(
			ctx, infraIdNumerical, vmInstanceIdNumerical).Execute()
	case "shutdown":
		httpRes, err = client.VMInstanceAPI.ShutdownVMInstance(
			ctx, infraIdNumerical, vmInstanceIdNumerical).Execute()
	case "reboot":
		httpRes, err = client.VMInstanceAPI.RebootVMInstance(
			ctx, infraIdNumerical, vmInstanceIdNumerical).Execute()
	default:
		return fmt.Errorf("unsupported power action: %s. Use 'start', 'shutdown', or 'reboot'", action)
	}

	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM instance %s power action '%s' successful", vmInstanceId, action)
	return nil
}

func VMInstanceGetPowerStatus(ctx context.Context, infrastructureId string, vmInstanceId string) error {
	logger.Get().Info().Msgf("Get VM instance power status for %s in infrastructure %s", vmInstanceId, infrastructureId)

	infraIdNumerical, err := utils.GetFloat32FromString(infrastructureId)
	if err != nil {
		return err
	}

	vmInstanceIdNumerical, err := utils.GetFloat32FromString(vmInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	powerStatus, httpRes, err := client.VMInstanceAPI.GetVMInstancePowerStatus(
		ctx, infraIdNumerical, vmInstanceIdNumerical).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM instance %s power status: %s", vmInstanceId, powerStatus)
	return nil
}
