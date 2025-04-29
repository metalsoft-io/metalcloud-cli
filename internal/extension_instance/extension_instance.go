package extension_instance

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/internal/infrastructure"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var extensionInstancePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Label": {
			MaxWidth: 30,
			Order:    2,
		},
		"ExtensionId": {
			Title: "Extension ID",
			Order: 3,
		},
		"InfrastructureId": {
			Title: "Infra ID",
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

func ExtensionInstanceList(ctx context.Context, infrastructureIdOrLabel string) error {
	logger.Get().Info().Msgf("List all extension instances for infrastructure %s", infrastructureIdOrLabel)

	infra, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	instances, httpRes, err := client.ExtensionInstanceAPI.GetExtensionInstances(ctx, float32(infra.Id)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(instances.Data, &extensionInstancePrintConfig)
}

func ExtensionInstanceGet(ctx context.Context, extensionInstanceId string) error {
	logger.Get().Info().Msgf("Get extension instance details for %s", extensionInstanceId)

	id, err := GetExtensionInstanceId(extensionInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	instance, httpRes, err := client.ExtensionInstanceAPI.GetExtensionInstance(ctx, id).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(instance, &extensionInstancePrintConfig)
}

func ExtensionInstanceCreate(ctx context.Context, infrastructureIdOrLabel string, config []byte) error {
	logger.Get().Info().Msgf("Create new extension instance in infrastructure %s", infrastructureIdOrLabel)

	var payload sdk.CreateExtensionInstance
	if err := json.Unmarshal(config, &payload); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	infra, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	instance, httpRes, err := client.ExtensionInstanceAPI.CreateExtensionInstance(ctx, float32(infra.Id)).CreateExtensionInstance(payload).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(instance, &extensionInstancePrintConfig)
}

func ExtensionInstanceUpdate(ctx context.Context, extensionInstanceId string, config []byte) error {
	logger.Get().Info().Msgf("Update extension instance %s", extensionInstanceId)

	id, err := GetExtensionInstanceId(extensionInstanceId)
	if err != nil {
		return err
	}

	var payload sdk.UpdateExtensionInstance
	if err := json.Unmarshal(config, &payload); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	client := api.GetApiClient(ctx)

	instance, httpRes, err := client.ExtensionInstanceAPI.GetExtensionInstance(ctx, id).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	instance, httpRes, err = client.ExtensionInstanceAPI.UpdateExtensionInstance(ctx, id).
		IfMatch(strconv.Itoa(int(instance.Revision))).
		UpdateExtensionInstance(payload).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(instance, &extensionInstancePrintConfig)
}

func ExtensionInstanceDelete(ctx context.Context, extensionInstanceId string) error {
	logger.Get().Info().Msgf("Delete extension instance %s", extensionInstanceId)

	id, err := GetExtensionInstanceId(extensionInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	instance, httpRes, err := client.ExtensionInstanceAPI.GetExtensionInstance(ctx, id).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	httpRes, err = client.ExtensionInstanceAPI.DeleteExtensionInstance(ctx, id).
		IfMatch(strconv.Itoa(int(instance.Revision))).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Extension instance %s successfully deleted", extensionInstanceId)
	return nil
}

func GetExtensionInstanceId(extensionInstanceId string) (float32, error) {
	id, err := strconv.ParseFloat(extensionInstanceId, 32)
	if err != nil {
		err := fmt.Errorf("invalid extension instance ID: '%s'", extensionInstanceId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}
	return float32(id), nil
}
