package server_instance

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

var serverInstanceGroupPrintConfig = formatter.PrintConfig{
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

func ServerInstanceGroupList(ctx context.Context, infrastructureIdOrLabel string) error {
	logger.Get().Info().Msgf("List all server instance groups for infrastructure %s", infrastructureIdOrLabel)

	infra, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	serverInstanceGroupList, httpRes, err := client.ServerInstanceGroupAPI.GetInfrastructureServerInstanceGroups(ctx, int32(infra.Id)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(serverInstanceGroupList, &serverInstanceGroupPrintConfig)
}

func ServerInstanceGroupGet(ctx context.Context, serverInstanceGroupId string) error {
	logger.Get().Info().Msgf("Get server instance group details for %s", serverInstanceGroupId)

	serverInstanceGroupIdNumerical, err := utils.GetFloat32FromString(serverInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	serverInstanceGroup, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroup(ctx, int32(serverInstanceGroupIdNumerical)).Execute()
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
		Label:         &label,
		ServerTypeId:  sdk.PtrInt32(int32(serverType.Id)),
		InstanceCount: sdk.PtrInt32(int32(serverInstanceCountNumerical)),
	}

	if osTemplateId != "" {
		osTemplate, err := os_template.GetOsTemplateByIdOrLabel(ctx, osTemplateId)
		if err != nil {
			return err
		}

		payload.OsTemplateId = sdk.PtrInt32(int32(osTemplate.Id))
	}

	client := api.GetApiClient(ctx)

	serverInstanceGroupInfo, httpRes, err := client.ServerInstanceGroupAPI.CreateServerInstanceGroup(ctx, int32(infra.Id)).ServerInstanceGroupCreate(payload).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(serverInstanceGroupInfo, &serverInstanceGroupPrintConfig)
}

func ServerInstanceGroupUpdate(ctx context.Context, serverInstanceGroupId string, label string, instanceCount int, osTemplateId int) error {
	logger.Get().Info().Msgf("Update server instance group %s", serverInstanceGroupId)

	serverInstanceGroupIdNumerical, err := utils.GetFloat32FromString(serverInstanceGroupId)
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
		payload.OsTemplateId = sdk.PtrInt32(int32(osTemplateId))
	}

	client := api.GetApiClient(ctx)

	serverInstanceGroupConfig, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroupConfig(ctx, int32(serverInstanceGroupIdNumerical)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	serverInstanceGroupConfig, httpRes, err = client.ServerInstanceGroupAPI.UpdateServerInstanceGroupConfig(ctx, int32(serverInstanceGroupIdNumerical)).
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

	serverInstanceGroupIdNumerical, err := utils.GetFloat32FromString(serverInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	serverInstanceGroup, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroup(ctx, int32(serverInstanceGroupIdNumerical)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	httpRes, err = client.ServerInstanceGroupAPI.DeleteServerInstanceGroup(ctx, int32(serverInstanceGroupIdNumerical)).
		IfMatch(strconv.Itoa(int(serverInstanceGroup.Revision))).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return nil
}

func ServerInstanceGroupInstances(ctx context.Context, serverInstanceGroupId string) error {
	logger.Get().Info().Msgf("List instances of server instance group %s", serverInstanceGroupId)

	serverInstanceGroupIdNumerical, err := utils.GetFloat32FromString(serverInstanceGroupId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	serverInstancesList, httpRes, err := client.ServerInstanceGroupAPI.GetServerInstanceGroupServerInstances(ctx, int32(serverInstanceGroupIdNumerical)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(serverInstancesList, &serverInstancePrintConfig)
}
