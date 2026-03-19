package network_device_configuration_template

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
	ValidDeviceTemplateActions      = []string{"add-global-config", "remove-global-config", "add-neighbor", "remove-neighbor"}
	ValidDeviceTemplateNetworkTypes = []string{"underlay", "overlay"}
	ValidNetworkDeviceDrivers       = []string{"cisco_aci51", "nvidia_ufm", "nexus9000", "cumulus42", "arista_eos", "dell_s4048", "hp5800", "hp5900", "hp5950", "dummy", "junos", "os_10", "sonic_enterprise", "vmware_vds", "cumulus_linux", "brocade", "nvidia_dpu", "dell_s4000", "dell_s6010", "junos18"}
	ValidNetworkDevicePositions     = []string{"all", "tor", "north", "spine", "leaf", "other"}
	ValidBgpNumberings              = []string{"numbered", "unnumbered"}
	ValidBgpLinkConfigurations      = []string{"disabled", "active", "passive"}
)

type ConfigFieldGuideEntry struct {
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
		"BgpNumbering": {
			Title: "BGP Numbering",
			Order: 6,
		},
		"BgpLinkConfiguration": {
			Title: "BGP Link Configuration",
			Order: 7,
		},
		"LibraryLabel": {
			Title: "Library Label",
			Order: 8,
		},
	},
}

func NetworkDeviceConfigurationTemplateList(ctx context.Context, filterId []string, filterLibraryLabel []string) error {
	logger.Get().Info().Msgf("Listing all network device configuration templates")

	client := api.GetApiClient(ctx)

	request := client.NetworkDeviceBGPConfigurationTemplateAPI.GetNetworkDeviceBGPConfigurationTemplates(ctx)

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
	if formatter.IsNativeFormat() {
		type templateExample struct {
			Action                      string `json:"action" yaml:"action"`
			NetworkType                 string `json:"networkType" yaml:"networkType"`
			NetworkDeviceDriver         string `json:"networkDeviceDriver" yaml:"networkDeviceDriver"`
			NetworkDevicePosition       string `json:"networkDevicePosition" yaml:"networkDevicePosition"`
			RemoteNetworkDevicePosition string `json:"remoteNetworkDevicePosition" yaml:"remoteNetworkDevicePosition"`
			BgpNumbering                string `json:"bgpNumbering" yaml:"bgpNumbering"`
			BgpLinkConfiguration        string `json:"bgpLinkConfiguration" yaml:"bgpLinkConfiguration"`
			ExecutionType               string `json:"executionType" yaml:"executionType"`
			LibraryLabel                string `json:"libraryLabel" yaml:"libraryLabel"`
			Preparation                 string `json:"preparation,omitempty" yaml:"preparation,omitempty"`
			Configuration               string `json:"configuration" yaml:"configuration"`
		}
		return formatter.PrintResult(templateExample{
			Action:                      strings.Join(ValidDeviceTemplateActions, "|"),
			NetworkType:                 strings.Join(ValidDeviceTemplateNetworkTypes, "|"),
			NetworkDeviceDriver:         strings.Join(ValidNetworkDeviceDrivers, "|"),
			NetworkDevicePosition:       strings.Join(ValidNetworkDevicePositions, "|"),
			RemoteNetworkDevicePosition: strings.Join(ValidNetworkDevicePositions, "|"),
			BgpNumbering:                strings.Join(ValidBgpNumberings, "|"),
			BgpLinkConfiguration:        strings.Join(ValidBgpLinkConfigurations, "|"),
			ExecutionType:               "cli",
			LibraryLabel:                "<string>",
			Configuration:               "<base64 encoded commands>",
		}, nil)
	}

	entries := []ConfigFieldGuideEntry{
		{Field: "action", Required: "yes", AcceptedValues: strings.Join(ValidDeviceTemplateActions, ", "), Example: ValidDeviceTemplateActions[0]},
		{Field: "networkType", Required: "yes", AcceptedValues: strings.Join(ValidDeviceTemplateNetworkTypes, ", "), Example: ValidDeviceTemplateNetworkTypes[0]},
		{Field: "networkDeviceDriver", Required: "yes", AcceptedValues: strings.Join(ValidNetworkDeviceDrivers, ", "), Example: "junos"},
		{Field: "networkDevicePosition", Required: "yes", AcceptedValues: strings.Join(ValidNetworkDevicePositions, ", "), Example: ValidNetworkDevicePositions[0]},
		{Field: "remoteNetworkDevicePosition", Required: "yes", AcceptedValues: strings.Join(ValidNetworkDevicePositions, ", "), Example: ValidNetworkDevicePositions[0]},
		{Field: "bgpNumbering", Required: "yes", AcceptedValues: strings.Join(ValidBgpNumberings, ", "), Example: ValidBgpNumberings[0]},
		{Field: "bgpLinkConfiguration", Required: "yes", AcceptedValues: strings.Join(ValidBgpLinkConfigurations, ", "), Example: ValidBgpLinkConfigurations[1]},
		{Field: "executionType", Required: "yes", AcceptedValues: "cli", Example: "cli"},
		{Field: "libraryLabel", Required: "yes", AcceptedValues: "any string", Example: "my-template"},
		{Field: "preparation", Required: "no", AcceptedValues: "base64 encoded commands", Example: ""},
		{Field: "configuration", Required: "yes", AcceptedValues: "base64 encoded commands", Example: ""},
	}
	return formatter.PrintResult(entries, &configFieldGuidePrintConfig)
}

func NetworkDeviceConfigurationTemplateGet(ctx context.Context, networkDeviceConfigurationTemplateId string) error {
	logger.Get().Info().Msgf("Get network device configuration template %s details", networkDeviceConfigurationTemplateId)

	networkDeviceConfigurationTemplateIdNumeric, err := getNetworkDeviceConfigurationTemplateId(networkDeviceConfigurationTemplateId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkDeviceConfigurationTemplate, httpRes, err := client.NetworkDeviceBGPConfigurationTemplateAPI.
		GetNetworkDeviceBGPConfigurationTemplate(ctx, networkDeviceConfigurationTemplateIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}
	return formatter.PrintResult(networkDeviceConfigurationTemplate, &NetworkDeviceConfigurationTemplatePrintConfig)
}

func NetworkDeviceConfigurationTemplateCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating network device configuration template")

	var networkDeviceConfigurationTemplateConfig sdk.CreateNetworkDeviceBGPConfigurationTemplate
	err := utils.UnmarshalContent(config, &networkDeviceConfigurationTemplateConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkDeviceConfigurationTemplateInfo, httpRes, err := client.NetworkDeviceBGPConfigurationTemplateAPI.
		CreateNetworkDeviceBGPConfigurationTemplate(ctx).
		CreateNetworkDeviceBGPConfigurationTemplate(networkDeviceConfigurationTemplateConfig).
		Execute()
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

	var networkDeviceConfigurationTemplateConfig sdk.UpdateNetworkDeviceBGPConfigurationTemplate
	err = utils.UnmarshalContent(config, &networkDeviceConfigurationTemplateConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkDeviceConfigurationTemplateInfo, httpRes, err := client.NetworkDeviceBGPConfigurationTemplateAPI.
		UpdateNetworkDeviceBGPConfigurationTemplate(ctx, networkDeviceConfigurationTemplateIdNumeric).
		UpdateNetworkDeviceBGPConfigurationTemplate(networkDeviceConfigurationTemplateConfig).
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

	httpRes, err := client.NetworkDeviceBGPConfigurationTemplateAPI.
		DeleteNetworkDeviceBGPConfigurationTemplate(ctx, networkDeviceConfigurationTemplateIdNumeric).
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
