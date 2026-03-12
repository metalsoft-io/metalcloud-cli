package network_device_link_aggregation_configuration_template

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var NetworkDeviceLinkAggregationConfigurationTemplatePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Action": {
			Title: "Action",
			Order: 2,
		},
		"AggregationType": {
			Title: "Aggregation Type",
			Order: 3,
		},
		"NetworkDeviceDriver": {
			Title: "Network Device Driver",
			Order: 4,
		},
		"ExecutionType": {
			Title: "Execution Type",
			Order: 5,
		},
		"LibraryLabel": {
			Title: "Library Label",
			Order: 6,
		},
	},
}

func NetworkDeviceLinkAggregationConfigurationTemplateList(ctx context.Context, filterId []string, filterLibraryLabel []string) error {
	logger.Get().Info().Msgf("Listing all network device link aggregation configuration templates")

	client := api.GetApiClient(ctx)

	request := client.NetworkDeviceLinkAggregationConfigurationTemplateAPI.GetNetworkDeviceLinkAggregationConfigurationTemplates(ctx)

	if len(filterId) > 0 {
		request = request.FilterId(utils.ProcessFilterStringSlice(filterId))
	}
	if len(filterLibraryLabel) > 0 {
		request = request.FilterLibraryLabel(filterLibraryLabel)
	}

	networkDeviceLinkAggregationConfigurationTemplateList, httpRes, err := request.SortBy([]string{"id:ASC"}).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkDeviceLinkAggregationConfigurationTemplateList, &NetworkDeviceLinkAggregationConfigurationTemplatePrintConfig)
}

func NetworkDeviceLinkAggregationConfigurationTemplateConfigExample(ctx context.Context) error {
	networkDeviceLinkAggregationConfiguration := sdk.CreateNetworkDeviceLinkAggregationConfigurationTemplate{
		Action:              "create",
		AggregationType:     "lacp",
		NetworkDeviceDriver: "junos",
		ExecutionType:       "cli",
		LibraryLabel:        "string",
		Preparation:         sdk.PtrString("string"),
		Configuration:       "string",
	}

	return formatter.PrintResult(networkDeviceLinkAggregationConfiguration, nil)
}

func NetworkDeviceLinkAggregationConfigurationTemplateGet(ctx context.Context, networkDeviceLinkAggregationConfigurationTemplateId string) error {
	logger.Get().Info().Msgf("Get network device link aggregation configuration template %s details", networkDeviceLinkAggregationConfigurationTemplateId)

	networkDeviceLinkAggregationConfigurationTemplateIdNumeric, err := getNetworkDeviceLinkAggregationConfigurationTemplateId(networkDeviceLinkAggregationConfigurationTemplateId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkDeviceLinkAggregationConfigurationTemplate, httpRes, err := client.NetworkDeviceLinkAggregationConfigurationTemplateAPI.
		GetNetworkDeviceLinkAggregationConfigurationTemplate(ctx, networkDeviceLinkAggregationConfigurationTemplateIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}
	return formatter.PrintResult(networkDeviceLinkAggregationConfigurationTemplate, &NetworkDeviceLinkAggregationConfigurationTemplatePrintConfig)
}

func NetworkDeviceLinkAggregationConfigurationTemplateCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating network device link aggregation configuration template")

	var networkDeviceLinkAggregationConfigurationTemplateConfig sdk.CreateNetworkDeviceLinkAggregationConfigurationTemplate
	err := utils.UnmarshalContent(config, &networkDeviceLinkAggregationConfigurationTemplateConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkDeviceLinkAggregationConfigurationTemplateInfo, httpRes, err := client.NetworkDeviceLinkAggregationConfigurationTemplateAPI.
		CreateNetworkDeviceLinkAggregationConfigurationTemplate(ctx).
		CreateNetworkDeviceLinkAggregationConfigurationTemplate(networkDeviceLinkAggregationConfigurationTemplateConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkDeviceLinkAggregationConfigurationTemplateInfo, &NetworkDeviceLinkAggregationConfigurationTemplatePrintConfig)
}

func NetworkDeviceLinkAggregationConfigurationTemplateUpdate(ctx context.Context, networkDeviceLinkAggregationConfigurationTemplateId string, config []byte) error {
	logger.Get().Info().Msgf("Updating network device link aggregation configuration template %s", networkDeviceLinkAggregationConfigurationTemplateId)

	networkDeviceLinkAggregationConfigurationTemplateIdNumeric, err := getNetworkDeviceLinkAggregationConfigurationTemplateId(networkDeviceLinkAggregationConfigurationTemplateId)
	if err != nil {
		return err
	}

	var networkDeviceLinkAggregationConfigurationTemplateConfig sdk.UpdateNetworkDeviceLinkAggregationConfigurationTemplate
	err = utils.UnmarshalContent(config, &networkDeviceLinkAggregationConfigurationTemplateConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkDeviceLinkAggregationConfigurationTemplateInfo, httpRes, err := client.NetworkDeviceLinkAggregationConfigurationTemplateAPI.
		UpdateNetworkDeviceLinkAggregationConfigurationTemplate(ctx, networkDeviceLinkAggregationConfigurationTemplateIdNumeric).
		UpdateNetworkDeviceLinkAggregationConfigurationTemplate(networkDeviceLinkAggregationConfigurationTemplateConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkDeviceLinkAggregationConfigurationTemplateInfo, &NetworkDeviceLinkAggregationConfigurationTemplatePrintConfig)
}

func NetworkDeviceLinkAggregationConfigurationTemplateDelete(ctx context.Context, networkDeviceLinkAggregationConfigurationTemplateId string) error {
	logger.Get().Info().Msgf("Deleting network device link aggregation configuration template %s", networkDeviceLinkAggregationConfigurationTemplateId)

	networkDeviceLinkAggregationConfigurationTemplateIdNumeric, err := getNetworkDeviceLinkAggregationConfigurationTemplateId(networkDeviceLinkAggregationConfigurationTemplateId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceLinkAggregationConfigurationTemplateAPI.
		DeleteNetworkDeviceLinkAggregationConfigurationTemplate(ctx, networkDeviceLinkAggregationConfigurationTemplateIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network device link aggregation configuration template %s deleted", networkDeviceLinkAggregationConfigurationTemplateId)
	return nil
}

func getNetworkDeviceLinkAggregationConfigurationTemplateId(networkDeviceLinkAggregationConfigurationTemplateId string) (float32, error) {
	networkDeviceLinkAggregationConfigurationTemplateIdNumeric, err := strconv.ParseFloat(networkDeviceLinkAggregationConfigurationTemplateId, 32)
	if err != nil {
		err := fmt.Errorf("invalid network device link aggregation configuration template ID: '%s'", networkDeviceLinkAggregationConfigurationTemplateId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(networkDeviceLinkAggregationConfigurationTemplateIdNumeric), nil
}
