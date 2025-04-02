package instance

import (
	"context"
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

var instanceGroupPrintConfig = formatter.PrintConfig{
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

var instanceGroupConfigPrintConfig = formatter.PrintConfig{
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

func InstanceGroupList(ctx context.Context, infrastructureIdOrLabel string) error {
	logger.Get().Info().Msgf("List all instance groups for infrastructure %s", infrastructureIdOrLabel)

	infra, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	instanceGroupList, httpRes, err := client.ServerInstanceGroupAPI.GetInfrastructureServerInstanceGroups(ctx, int32(infra.Id)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(instanceGroupList, &instanceGroupPrintConfig)
}

func InstanceGroupGet(ctx context.Context, instanceGroupId string) error {
	logger.Get().Info().Msgf("Get instance group details for %s", instanceGroupId)

	instanceGroupIdNumerical, err := utils.GetFloat32FromString(instanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	instanceGroup, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroup(ctx, int32(instanceGroupIdNumerical)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(instanceGroup, &instanceGroupPrintConfig)
}

func InstanceGroupCreate(ctx context.Context, infrastructureIdOrLabel string, label string, serverTypeId string, instanceCount string, osTemplateId string) error {
	logger.Get().Info().Msgf("Create new instance group in infrastructure %s", infrastructureIdOrLabel)

	instanceCountNumerical, err := utils.GetFloat32FromString(instanceCount)
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
		Label:         &label,
		ServerTypeId:  sdk.PtrInt32(int32(serverType.Id)),
		InstanceCount: sdk.PtrInt32(int32(instanceCountNumerical)),
	}

	if osTemplateId != "" {
		osTemplate, err := os_template.GetOsTemplateByIdOrLabel(ctx, osTemplateId)
		if err != nil {
			return err
		}

		payload.VolumeTemplateId = sdk.PtrInt32(int32(osTemplate.Id))
	}

	client := api.GetApiClient(ctx)

	instanceGroupInfo, httpRes, err := client.ServerInstanceGroupAPI.CreateServerInstanceGroup(ctx, int32(infra.Id)).ServerInstanceGroupCreate(payload).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(instanceGroupInfo, &instanceGroupPrintConfig)
}

func InstanceGroupUpdate(ctx context.Context, instanceGroupId string, label string, instanceCount int, osTemplateId int) error {
	logger.Get().Info().Msgf("Update instance group %s", instanceGroupId)

	instanceGroupIdNumerical, err := utils.GetFloat32FromString(instanceGroupId)
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
		payload.VolumeTemplateId = sdk.PtrInt32(int32(osTemplateId))
	}

	client := api.GetApiClient(ctx)

	instanceGroupConfig, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroupConfig(ctx, int32(instanceGroupIdNumerical)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	instanceGroupConfig, httpRes, err = client.ServerInstanceGroupAPI.UpdateServerInstanceGroupConfig(ctx, int32(instanceGroupIdNumerical)).
		IfMatch(strconv.Itoa(int(instanceGroupConfig.Revision))).
		ServerInstanceGroupUpdate(payload).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(instanceGroupConfig, &instanceGroupConfigPrintConfig)
}

func InstanceGroupDelete(ctx context.Context, instanceGroupId string) error {
	logger.Get().Info().Msgf("Delete instance group %s", instanceGroupId)

	instanceGroupIdNumerical, err := utils.GetFloat32FromString(instanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	instanceGroup, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroup(ctx, int32(instanceGroupIdNumerical)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	httpRes, err = client.ServerInstanceGroupAPI.DeleteServerInstanceGroup(ctx, int32(instanceGroupIdNumerical)).
		IfMatch(strconv.Itoa(int(instanceGroup.Revision))).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return nil
}

func InstanceGroupInstances(ctx context.Context, instanceGroupId string) error {
	logger.Get().Info().Msgf("List instances of instance group %s", instanceGroupId)

	instanceGroupIdNumerical, err := utils.GetFloat32FromString(instanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	instancesList, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroupServerInstances(ctx, int32(instanceGroupIdNumerical)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(instancesList, &instancePrintConfig)
}
