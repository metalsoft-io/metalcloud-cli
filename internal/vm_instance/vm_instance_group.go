package vm_instance

import (
	"context"
	"encoding/json"
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

func VMInstanceGroupGet(ctx context.Context, infrastructureId string, vmInstanceGroupId string) error {
	logger.Get().Info().Msgf("Get VM instance group details for %s in infrastructure %s", vmInstanceGroupId, infrastructureId)

	infraIdNumerical, err := utils.GetFloat32FromString(infrastructureId)
	if err != nil {
		return err
	}

	vmInstanceGroupIdNumerical, err := utils.GetFloat32FromString(vmInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmInstanceGroup, httpRes, err := client.VMInstanceGroupAPI.GetInfrastructureVMInstanceGroup(
		ctx, infraIdNumerical, vmInstanceGroupIdNumerical).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmInstanceGroup, &vmInstanceGroupPrintConfig)
}

func VMInstanceGroupList(ctx context.Context, infrastructureId string) error {
	logger.Get().Info().Msgf("List all VM instance groups for infrastructure %s", infrastructureId)

	infraIdNumerical, err := utils.GetFloat32FromString(infrastructureId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmInstanceGroupsList, httpRes, err := client.VMInstanceGroupAPI.GetInfrastructureVMInstanceGroups(
		ctx, infraIdNumerical).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmInstanceGroupsList.Data, &vmInstanceGroupPrintConfig)
}

func VMInstanceGroupCreate(ctx context.Context, infrastructureId string, vmTypeId string, diskSizeGB string, instanceCount string, osTemplateId string) error {
	logger.Get().Info().Msgf("Create new VM instance group in infrastructure %s", infrastructureId)

	infraIdNumerical, err := utils.GetFloat32FromString(infrastructureId)
	if err != nil {
		return err
	}

	vmTypeIdNumerical, err := utils.GetFloat32FromString(vmTypeId)
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
		payload.OsTemplateId, err = utils.GetFloat32FromString(osTemplateId)
		if err != nil {
			return err
		}
	}

	client := api.GetApiClient(ctx)

	vmInstanceGroupInfo, httpRes, err := client.VMInstanceGroupAPI.CreateVMInstanceGroup(
		ctx, infraIdNumerical).CreateVMInstanceGroup(payload).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmInstanceGroupInfo, &vmInstanceGroupPrintConfig)
}

func VMInstanceGroupUpdate(ctx context.Context, infrastructureId string, vmInstanceGroupId string, label string, customVariables []byte) error {
	logger.Get().Info().Msgf("Update VM instance group %s in infrastructure %s", vmInstanceGroupId, infrastructureId)

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

	infraIdNumerical, vmInstanceGroupIdNumerical, revision, err := getVmInstanceGroupIdAndRevision(ctx, infrastructureId, vmInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmInstanceGroup, httpRes, err := client.VMInstanceGroupAPI.
		UpdateVMInstanceGroupConfig(ctx, infraIdNumerical, vmInstanceGroupIdNumerical).
		UpdateVMInstanceGroup(payload).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmInstanceGroup, &vmInstanceGroupPrintConfig)
}

func VMInstanceGroupDelete(ctx context.Context, infrastructureId string, vmInstanceGroupId string) error {
	logger.Get().Info().Msgf("Delete VM instance group %s from infrastructure %s", vmInstanceGroupId, infrastructureId)

	infraIdNumerical, vmInstanceGroupIdNumerical, revision, err := getVmInstanceGroupIdAndRevision(ctx, infrastructureId, vmInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.VMInstanceGroupAPI.
		DeleteVMInstanceGroup(ctx, infraIdNumerical, vmInstanceGroupIdNumerical).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("VM instance group %s successfully deleted", vmInstanceGroupId)
	return nil
}

func VMInstanceGroupInstances(ctx context.Context, infrastructureId string, vmInstanceGroupId string) error {
	logger.Get().Info().Msgf("List instances of VM instance group %s in infrastructure %s", vmInstanceGroupId, infrastructureId)

	infraIdNumerical, err := utils.GetFloat32FromString(infrastructureId)
	if err != nil {
		return err
	}

	vmInstanceGroupIdNumerical, err := utils.GetFloat32FromString(vmInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	vmInstancesList, httpRes, err := client.VMInstanceGroupAPI.GetVMInstanceGroupVMInstances(
		ctx, infraIdNumerical, vmInstanceGroupIdNumerical).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(vmInstancesList.Data, &vmInstancePrintConfig)
}

func getVmInstanceGroupIdAndRevision(ctx context.Context, infrastructureId string, vmInstanceGroupId string) (float32, float32, string, error) {
	infraIdNumerical, err := utils.GetFloat32FromString(infrastructureId)
	if err != nil {
		return 0, 0, "", err
	}

	vmInstanceGroupIdNumerical, err := utils.GetFloat32FromString(vmInstanceGroupId)
	if err != nil {
		return 0, 0, "", err
	}

	client := api.GetApiClient(ctx)

	vmGroup, httpRes, err := client.VMInstanceGroupAPI.
		GetInfrastructureVMInstanceGroup(ctx, infraIdNumerical, vmInstanceGroupIdNumerical).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, 0, "", err
	}

	return infraIdNumerical, vmInstanceGroupIdNumerical, strconv.Itoa(int(vmGroup.Revision)), nil
}
