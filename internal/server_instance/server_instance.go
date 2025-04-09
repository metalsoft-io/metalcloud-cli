package server_instance

import (
	"context"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
)

var serverInstancePrintConfig = formatter.PrintConfig{
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

func ServerInstanceGet(ctx context.Context, serverInstanceId string) error {
	logger.Get().Info().Msgf("Get server instance details for %s", serverInstanceId)

	serverInstanceIdNumerical, err := utils.GetFloat32FromString(serverInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	serverInstanceInfo, httpRes, err := client.ServerInstanceAPI.GetServerInstance(ctx, int32(serverInstanceIdNumerical)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(serverInstanceInfo, &serverInstancePrintConfig)
}
