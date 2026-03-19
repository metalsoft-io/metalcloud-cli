package network_device_link_aggregation_configuration_template

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var (
	ValidLinkAggregationTemplateActions = []string{"create", "delete", "add-member", "remove-member"}
	ValidAggregationTypes               = []string{"lag", "mlag", "mlag-peer-link"}
	ValidNetworkDeviceDrivers           = []string{"cisco_aci51", "nvidia_ufm", "nexus9000", "cumulus42", "arista_eos", "dell_s4048", "hp5800", "hp5900", "hp5950", "dummy", "junos", "os_10", "sonic_enterprise", "vmware_vds", "cumulus_linux", "brocade", "nvidia_dpu", "dell_s4000", "dell_s6010", "junos18"}
)

type configFieldGuideEntry struct {
	Field          string `json:"field" yaml:"field"`
	Required       string `json:"required" yaml:"required"`
	AcceptedValues string `json:"acceptedValues" yaml:"acceptedValues"`
	Example        string `json:"example" yaml:"example"`
}

var configFieldGuidePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Field": {
			Title: "Field",
			Order: 1,
		},
		"Required": {
			Title: "Required",
			Order: 2,
		},
		"AcceptedValues": {
			Title:    "Accepted Values",
			MaxWidth: 55,
			Order:    3,
		},
		"Example": {
			Title: "Example",
			Order: 4,
		},
	},
}

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
	if formatter.IsNativeFormat() {
		type templateExample struct {
			Action              string `json:"action" yaml:"action"`
			AggregationType     string `json:"aggregationType" yaml:"aggregationType"`
			NetworkDeviceDriver string `json:"networkDeviceDriver" yaml:"networkDeviceDriver"`
			ExecutionType       string `json:"executionType" yaml:"executionType"`
			LibraryLabel        string `json:"libraryLabel" yaml:"libraryLabel"`
			Preparation         string `json:"preparation,omitempty" yaml:"preparation,omitempty"`
			Configuration       string `json:"configuration" yaml:"configuration"`
		}
		return formatter.PrintResult(templateExample{
			Action:              strings.Join(ValidLinkAggregationTemplateActions, "|"),
			AggregationType:     strings.Join(ValidAggregationTypes, "|"),
			NetworkDeviceDriver: strings.Join(ValidNetworkDeviceDrivers, "|"),
			ExecutionType:       "cli",
			LibraryLabel:        "<string>",
			Configuration:       "<base64 encoded commands>",
		}, nil)
	}

	entries := []configFieldGuideEntry{
		{Field: "action", Required: "yes", AcceptedValues: strings.Join(ValidLinkAggregationTemplateActions, ", "), Example: ValidLinkAggregationTemplateActions[0]},
		{Field: "aggregationType", Required: "yes", AcceptedValues: strings.Join(ValidAggregationTypes, ", "), Example: ValidAggregationTypes[0]},
		{Field: "networkDeviceDriver", Required: "yes", AcceptedValues: strings.Join(ValidNetworkDeviceDrivers, ", "), Example: "junos"},
		{Field: "executionType", Required: "yes", AcceptedValues: "cli", Example: "cli"},
		{Field: "libraryLabel", Required: "yes", AcceptedValues: "any string", Example: "my-template"},
		{Field: "preparation", Required: "no", AcceptedValues: "base64 encoded commands", Example: ""},
		{Field: "configuration", Required: "yes", AcceptedValues: "base64 encoded commands", Example: ""},
	}
	return formatter.PrintResult(entries, &configFieldGuidePrintConfig)
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
