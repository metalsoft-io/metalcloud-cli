package network_device_configuration_template

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

var NetworkDeviceConfigurationTemplatePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"NetworkType": {
			Title:    "Network Type",
			MaxWidth: 40,
			Order:    2,
		},
		"NetworkDeviceDriver": {
			Title: "Network Device Driver",
			Order: 3,
		},
		"NetworkDevicePosition": {
			Title: "Network Device Position",
			Order: 4,
		},
		"RemoteNetworkDevicePosition": {
			Title: "Remote Network Device Position",
			Order: 5,
		},
		"MlagPair": {
			Title: "Mlag Pair",
			Order: 6,
		},
		"BgpNumbering": {
			Title: "BGP Numbering",
			Order: 7,
		},
		"BgpLinkConfiguration": {
			Title: "BGP Link Configuration",
			Order: 8,
		},
		"LibraryLabel": {
			Title: "Library Label",
			Order: 9,
		},
	},
}

func NetworkDeviceConfigurationTemplateList(ctx context.Context, filterId []string, filterLibraryLabel []string) error {
	logger.Get().Info().Msgf("Listing all network device configuration templates")

	client := api.GetApiClient(ctx)

	request := client.NetworkDeviceConfigurationTemplateAPI.GetNetworkDeviceConfigurationTemplates(ctx)

	if len(filterId) > 0 {
		request = request.FilterId(utils.ProcessFilterStringSlice(filterId))
	}
	if len(filterLibraryLabel) > 0 {
		request = request.FilterLibraryLabel(filterLibraryLabel)
	}

	networkDeviceConfigurationTemplateList, httpRes, err := request.SortBy([]string{"id:ASC"}).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkDeviceConfigurationTemplateList, &NetworkDeviceConfigurationTemplatePrintConfig)
}

func NetworkDeviceConfigurationTemplateConfigExample(ctx context.Context) error {
	networkDeviceConfiguration := sdk.CreateNetworkDeviceConfigurationTemplate{
		NetworkType:                 "underlay",
		NetworkDeviceDriver:         "junos",
		NetworkDevicePosition:       "all",
		RemoteNetworkDevicePosition: "all",
		MlagPair:                    1,
		BgpNumbering:                "numbered",
		BgpLinkConfiguration:        "active",
		ExecutionType:               "cli",
		LibraryLabel:                "string",
		Preparation:                 sdk.PtrString("string"),
		Configuration:               "string",
	}

	return formatter.PrintResult(networkDeviceConfiguration, nil)
}

func NetworkDeviceConfigurationTemplateGet(ctx context.Context, networkDeviceConfigurationTemplateId string) error {
	logger.Get().Info().Msgf("Get network device configuration template %s details", networkDeviceConfigurationTemplateId)

	networkDeviceConfigurationTemplateIdNumeric, err := getNetworkDeviceConfigurationTemplateId(networkDeviceConfigurationTemplateId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkDeviceConfigurationTemplate, httpRes, err := client.NetworkDeviceConfigurationTemplateAPI.GetNetworkDeviceConfigurationTemplate(ctx, networkDeviceConfigurationTemplateIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}
	return formatter.PrintResult(networkDeviceConfigurationTemplate, &NetworkDeviceConfigurationTemplatePrintConfig)
}

func NetworkDeviceConfigurationTemplateCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating network device configuration template")

	var networkDeviceConfigurationTemplateConfig sdk.CreateNetworkDeviceConfigurationTemplate
	err := utils.UnmarshalContent(config, &networkDeviceConfigurationTemplateConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkDeviceConfigurationTemplateInfo, httpRes, err := client.NetworkDeviceConfigurationTemplateAPI.CreateNetworkDeviceConfigurationTemplate(ctx).CreateNetworkDeviceConfigurationTemplate(networkDeviceConfigurationTemplateConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkDeviceConfigurationTemplateInfo, &NetworkDeviceConfigurationTemplatePrintConfig)
}

func NetworkDeviceConfigurationTemplateUpdate(ctx context.Context, networkDeviceConfigurationTemplateId string, config []byte) error {
	logger.Get().Info().Msgf("Updating network device configuration template %s", networkDeviceConfigurationTemplateId)

	networkDeviceConfigurationTemplateIdNumeric, err := getNetworkDeviceConfigurationTemplateId(networkDeviceConfigurationTemplateId)
	if err != nil {
		return err
	}

	var networkDeviceConfigurationTemplateConfig sdk.UpdateNetworkDeviceConfigurationTemplate
	err = utils.UnmarshalContent(config, &networkDeviceConfigurationTemplateConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkDeviceConfigurationTemplateInfo, httpRes, err := client.NetworkDeviceConfigurationTemplateAPI.
		UpdateNetworkDeviceConfigurationTemplate(ctx, networkDeviceConfigurationTemplateIdNumeric).
		UpdateNetworkDeviceConfigurationTemplate(networkDeviceConfigurationTemplateConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkDeviceConfigurationTemplateInfo, &NetworkDeviceConfigurationTemplatePrintConfig)
}

func NetworkDeviceConfigurationTemplateDelete(ctx context.Context, networkDeviceConfigurationTemplateId string) error {
	logger.Get().Info().Msgf("Deleting network device configuration template %s", networkDeviceConfigurationTemplateId)

	networkDeviceConfigurationTemplateIdNumeric, err := getNetworkDeviceConfigurationTemplateId(networkDeviceConfigurationTemplateId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceConfigurationTemplateAPI.
		DeleteNetworkDeviceConfigurationTemplate(ctx, networkDeviceConfigurationTemplateIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network device configuration template %s deleted", networkDeviceConfigurationTemplateId)
	return nil
}

func getNetworkDeviceConfigurationTemplateId(networkDeviceConfigurationTemplateId string) (float32, error) {
	networkDeviceConfigurationTemplateIdNumeric, err := strconv.ParseFloat(networkDeviceConfigurationTemplateId, 32)
	if err != nil {
		err := fmt.Errorf("invalid network device configuration template ID: '%s'", networkDeviceConfigurationTemplateId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(networkDeviceConfigurationTemplateIdNumeric), nil
}
