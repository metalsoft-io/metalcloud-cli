// Package endpoint_instance implements the business logic behind the
// 'endpoint-instance' and 'endpoint-instance-group' command groups.
//
// Endpoint instances are created inside an infrastructure
// (POST /infrastructures/{id}/endpoint-instances) but are addressed globally
// afterwards (/endpoint-instances/{id}), which is why the create helpers take
// an infrastructure id or label while every other helper takes the instance id.
package endpoint_instance

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/internal/infrastructure"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// EndpointInstanceFilters carries the query filters shared by the global and
// the infrastructure scoped endpoint instance listings.
type EndpointInstanceFilters struct {
	InfrastructureId   []string
	GroupId            []string
	EndpointId         []string
	ServiceStatus      []string
	ConfigEndpointId   []string
	ConfigDeployStatus []string
	ConfigDeployType   []string
}

var endpointInstancePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Label": {
			Order:    2,
			MaxWidth: 30,
		},
		"InfrastructureId": {
			Title: "Infra ID",
			Order: 3,
		},
		"GroupId": {
			Title: "Group ID",
			Order: 4,
		},
		"EndpointId": {
			Title: "Endpoint ID",
			Order: 5,
		},
		"Hostname": {
			Title:    "Hostname",
			Order:    6,
			MaxWidth: 30,
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

var endpointInstanceConfigPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Label": {
			Order:    1,
			MaxWidth: 30,
		},
		"GroupId": {
			Title: "Group ID",
			Order: 2,
		},
		"EndpointId": {
			Title: "Endpoint ID",
			Order: 3,
		},
		"Hostname": {
			Title:    "Hostname",
			Order:    4,
			MaxWidth: 30,
		},
		"DeployType": {
			Title: "Deploy Type",
			Order: 5,
		},
		"DeployStatus": {
			Title:       "Deploy Status",
			Transformer: formatter.FormatStatusValue,
			Order:       6,
		},
		"Revision": {
			Title: "Revision",
			Order: 7,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       8,
		},
	},
}

// PrintConfig exposes the endpoint instance table layout so that other
// packages listing endpoint instances render them identically.
func PrintConfig() *formatter.PrintConfig {
	return &endpointInstancePrintConfig
}

// EndpointInstanceList lists endpoint instances. When infrastructureIdOrLabel
// is empty the global /endpoint-instances collection is used, otherwise the
// listing is scoped to that infrastructure.
func EndpointInstanceList(ctx context.Context, infrastructureIdOrLabel string, filters EndpointInstanceFilters) error {
	client := api.GetApiClient(ctx)

	if infrastructureIdOrLabel != "" {
		logger.Get().Info().Msgf("Listing endpoint instances of infrastructure '%s'", infrastructureIdOrLabel)

		infra, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
		if err != nil {
			return err
		}

		request := client.EndpointInstanceAPI.GetInfrastructureEndpointInstances(ctx, int64(infra.Id)).SortBy([]string{"id:ASC"})
		if len(filters.GroupId) > 0 {
			request = request.FilterGroupId(utils.ProcessFilterStringSlice(filters.GroupId))
		}
		if len(filters.EndpointId) > 0 {
			request = request.FilterEndpointId(utils.ProcessFilterStringSlice(filters.EndpointId))
		}
		if len(filters.ServiceStatus) > 0 {
			request = request.FilterServiceStatus(utils.ProcessFilterStringSlice(filters.ServiceStatus))
		}
		if len(filters.ConfigEndpointId) > 0 {
			request = request.FilterConfigEndpointId(utils.ProcessFilterStringSlice(filters.ConfigEndpointId))
		}
		if len(filters.ConfigDeployStatus) > 0 {
			request = request.FilterConfigDeployStatus(utils.ProcessFilterStringSlice(filters.ConfigDeployStatus))
		}
		if len(filters.ConfigDeployType) > 0 {
			request = request.FilterConfigDeployType(utils.ProcessFilterStringSlice(filters.ConfigDeployType))
		}

		records, meta, err := utils.FetchAllPages(request)
		if err != nil {
			return err
		}

		return utils.PrintAll(records, meta, len(records), &endpointInstancePrintConfig)
	}

	logger.Get().Info().Msg("Listing endpoint instances")

	request := client.EndpointInstanceAPI.GetEndpointInstances(ctx).SortBy([]string{"id:ASC"})
	if len(filters.InfrastructureId) > 0 {
		request = request.FilterInfrastructureId(utils.ProcessFilterStringSlice(filters.InfrastructureId))
	}
	if len(filters.GroupId) > 0 {
		request = request.FilterGroupId(utils.ProcessFilterStringSlice(filters.GroupId))
	}
	if len(filters.EndpointId) > 0 {
		request = request.FilterEndpointId(utils.ProcessFilterStringSlice(filters.EndpointId))
	}
	if len(filters.ServiceStatus) > 0 {
		request = request.FilterServiceStatus(utils.ProcessFilterStringSlice(filters.ServiceStatus))
	}
	if len(filters.ConfigEndpointId) > 0 {
		request = request.FilterConfigEndpointId(utils.ProcessFilterStringSlice(filters.ConfigEndpointId))
	}
	if len(filters.ConfigDeployStatus) > 0 {
		request = request.FilterConfigDeployStatus(utils.ProcessFilterStringSlice(filters.ConfigDeployStatus))
	}
	if len(filters.ConfigDeployType) > 0 {
		request = request.FilterConfigDeployType(utils.ProcessFilterStringSlice(filters.ConfigDeployType))
	}

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &endpointInstancePrintConfig)
}

func EndpointInstanceGet(ctx context.Context, endpointInstanceId string) error {
	logger.Get().Info().Msgf("Get endpoint instance '%s'", endpointInstanceId)

	endpointInstance, err := getEndpointInstance(ctx, endpointInstanceId)
	if err != nil {
		return err
	}

	return formatter.PrintResult(endpointInstance, &endpointInstancePrintConfig)
}

func EndpointInstanceConfigExample(ctx context.Context) error {
	example := sdk.EndpointInstanceCreate{
		Label:      sdk.PtrString("my-endpoint-instance"),
		GroupId:    sdk.PtrInt64(1),
		EndpointId: 1,
		Tags:       []string{"example"},
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(example, nil)
	}
	return formatter.PrintYamlResult(example)
}

func EndpointInstanceCreate(ctx context.Context, infrastructureIdOrLabel string, create sdk.EndpointInstanceCreate) error {
	logger.Get().Info().Msgf("Creating endpoint instance in infrastructure '%s'", infrastructureIdOrLabel)

	infra, err := infrastructure.GetInfrastructureByIdOrLabel(ctx, infrastructureIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	endpointInstance, httpRes, err := client.EndpointInstanceAPI.
		CreateEndpointInstance(ctx, int64(infra.Id)).
		EndpointInstanceCreate(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(endpointInstance, &endpointInstancePrintConfig)
}

func EndpointInstanceDelete(ctx context.Context, endpointInstanceId string) error {
	logger.Get().Info().Msgf("Deleting endpoint instance '%s'", endpointInstanceId)

	endpointInstance, err := getEndpointInstance(ctx, endpointInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.EndpointInstanceAPI.
		DeleteEndpointInstance(ctx, endpointInstance.Id).
		IfMatch(strconv.FormatInt(endpointInstance.Revision, 10)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Endpoint instance '%s' deleted", endpointInstanceId)
	return nil
}

func EndpointInstanceConfigGet(ctx context.Context, endpointInstanceId string) error {
	logger.Get().Info().Msgf("Get endpoint instance '%s' configuration", endpointInstanceId)

	config, _, err := getEndpointInstanceConfig(ctx, endpointInstanceId)
	if err != nil {
		return err
	}

	return formatter.PrintResult(config, &endpointInstanceConfigPrintConfig)
}

func EndpointInstanceConfigUpdate(ctx context.Context, endpointInstanceId string, config []byte) error {
	logger.Get().Info().Msgf("Updating endpoint instance '%s' configuration", endpointInstanceId)

	var update sdk.EndpointInstanceUpdate
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	currentConfig, endpointInstanceIdNumerical, err := getEndpointInstanceConfig(ctx, endpointInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	updatedConfig, httpRes, err := client.EndpointInstanceAPI.
		UpdateEndpointInstanceConfig(ctx, endpointInstanceIdNumerical).
		IfMatch(strconv.FormatInt(currentConfig.Revision, 10)).
		EndpointInstanceUpdate(update).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(updatedConfig, &endpointInstanceConfigPrintConfig)
}

func EndpointInstanceMetaUpdate(ctx context.Context, endpointInstanceId string, config []byte) error {
	logger.Get().Info().Msgf("Updating endpoint instance '%s' metadata", endpointInstanceId)

	var meta sdk.GenericMeta
	if err := utils.UnmarshalContent(config, &meta); err != nil {
		return err
	}

	endpointInstanceIdNumerical, err := GetEndpointInstanceId(endpointInstanceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.EndpointInstanceAPI.
		UpdateEndpointInstanceMeta(ctx, endpointInstanceIdNumerical).
		GenericMeta(meta).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Endpoint instance '%s' metadata updated", endpointInstanceId)
	return nil
}

func getEndpointInstance(ctx context.Context, endpointInstanceId string) (*sdk.EndpointInstance, error) {
	endpointInstanceIdNumerical, err := GetEndpointInstanceId(endpointInstanceId)
	if err != nil {
		return nil, err
	}

	client := api.GetApiClient(ctx)

	endpointInstance, httpRes, err := client.EndpointInstanceAPI.GetEndpointInstance(ctx, endpointInstanceIdNumerical).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	return endpointInstance, nil
}

// getEndpointInstanceConfig returns the instance configuration together with
// the numeric instance id. The configuration carries its own revision, which
// guards every write below the /config path.
func getEndpointInstanceConfig(ctx context.Context, endpointInstanceId string) (*sdk.EndpointInstanceConfiguration, int64, error) {
	endpointInstanceIdNumerical, err := GetEndpointInstanceId(endpointInstanceId)
	if err != nil {
		return nil, 0, err
	}

	client := api.GetApiClient(ctx)

	config, httpRes, err := client.EndpointInstanceAPI.GetEndpointInstanceConfig(ctx, endpointInstanceIdNumerical).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, 0, err
	}

	return config, endpointInstanceIdNumerical, nil
}

func GetEndpointInstanceId(endpointInstanceId string) (int64, error) {
	endpointInstanceIdNumerical, err := strconv.ParseInt(endpointInstanceId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid endpoint instance ID: '%s'", endpointInstanceId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return endpointInstanceIdNumerical, nil
}
