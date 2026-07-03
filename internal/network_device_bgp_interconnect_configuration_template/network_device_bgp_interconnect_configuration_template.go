package network_device_bgp_interconnect_configuration_template

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
	ValidExecutionTypes       = []string{"cli", "json_patch"}
	ValidNetworkDeviceDrivers = []string{"cisco_aci51", "nvidia_ufm", "nexus9000", "cumulus42", "arista_eos", "dell_s4048", "hp5800", "hp5900", "hp5950", "dummy", "junos", "os_10", "sonic_enterprise", "vmware_vds", "cumulus_linux", "brocade", "nvidia_dpu", "dell_s4000", "dell_s6010", "junos18"}
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

var NetworkDeviceBGPInterconnectConfigurationTemplatePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Label": {
			Title: "Label",
			Order: 2,
		},
		"Name": {
			Title: "Name",
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
		"Revision": {
			Title: "Revision",
			Order: 6,
		},
	},
}

func NetworkDeviceBGPInterconnectConfigurationTemplateList(ctx context.Context, filterId []string, filterNetworkDeviceDriver []string) error {
	logger.Get().Info().Msgf("Listing all network device BGP interconnect configuration templates")

	client := api.GetApiClient(ctx)

	request := client.NetworkDeviceBGPInterconnectConfigurationTemplateAPI.GetNetworkDeviceBGPInterconnectConfigurationTemplates(ctx)

	if len(filterId) > 0 {
		request = request.FilterId(utils.ProcessFilterStringSlice(filterId))
	}
	if len(filterNetworkDeviceDriver) > 0 {
		request = request.FilterNetworkDeviceDriver(filterNetworkDeviceDriver)
	}

	request = request.SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &NetworkDeviceBGPInterconnectConfigurationTemplatePrintConfig)
}

func NetworkDeviceBGPInterconnectConfigurationTemplateConfigExample(ctx context.Context) error {
	if formatter.IsNativeFormat() {
		type templateExample struct {
			Label               string `json:"label" yaml:"label"`
			Name                string `json:"name" yaml:"name"`
			NetworkDeviceDriver string `json:"networkDeviceDriver" yaml:"networkDeviceDriver"`
			ExecutionType       string `json:"executionType" yaml:"executionType"`
			AddGlobalConfig     string `json:"addGlobalConfig,omitempty" yaml:"addGlobalConfig,omitempty"`
			RemoveGlobalConfig  string `json:"removeGlobalConfig,omitempty" yaml:"removeGlobalConfig,omitempty"`
			AddNeighbor         string `json:"addNeighbor,omitempty" yaml:"addNeighbor,omitempty"`
			RemoveNeighbor      string `json:"removeNeighbor,omitempty" yaml:"removeNeighbor,omitempty"`
		}
		return formatter.PrintResult(templateExample{
			Label:               "<string>",
			Name:                "<string>",
			NetworkDeviceDriver: strings.Join(ValidNetworkDeviceDrivers, "|"),
			ExecutionType:       strings.Join(ValidExecutionTypes, "|"),
			AddGlobalConfig:     "<base64 encoded commands>",
			RemoveGlobalConfig:  "<base64 encoded commands>",
			AddNeighbor:         "<base64 encoded commands>",
			RemoveNeighbor:      "<base64 encoded commands>",
		}, nil)
	}

	entries := []configFieldGuideEntry{
		{Field: "label", Required: "yes", AcceptedValues: "any string", Example: "my-interconnect-template"},
		{Field: "name", Required: "yes", AcceptedValues: "any string", Example: "My Interconnect Template"},
		{Field: "networkDeviceDriver", Required: "yes", AcceptedValues: strings.Join(ValidNetworkDeviceDrivers, ", "), Example: "junos"},
		{Field: "executionType", Required: "yes", AcceptedValues: strings.Join(ValidExecutionTypes, ", "), Example: ValidExecutionTypes[0]},
		{Field: "addGlobalConfig", Required: "no", AcceptedValues: "base64 encoded commands", Example: ""},
		{Field: "removeGlobalConfig", Required: "no", AcceptedValues: "base64 encoded commands", Example: ""},
		{Field: "addNeighbor", Required: "no", AcceptedValues: "base64 encoded commands", Example: ""},
		{Field: "removeNeighbor", Required: "no", AcceptedValues: "base64 encoded commands", Example: ""},
	}
	return formatter.PrintResult(entries, &configFieldGuidePrintConfig)
}

func NetworkDeviceBGPInterconnectConfigurationTemplateGet(ctx context.Context, networkDeviceBGPInterconnectConfigurationTemplateId string) error {
	logger.Get().Info().Msgf("Get network device BGP interconnect configuration template %s details", networkDeviceBGPInterconnectConfigurationTemplateId)

	networkDeviceBGPInterconnectConfigurationTemplateIdNumeric, err := getNetworkDeviceBGPInterconnectConfigurationTemplateId(networkDeviceBGPInterconnectConfigurationTemplateId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkDeviceBGPInterconnectConfigurationTemplate, httpRes, err := client.NetworkDeviceBGPInterconnectConfigurationTemplateAPI.
		GetNetworkDeviceBGPInterconnectConfigurationTemplate(ctx, networkDeviceBGPInterconnectConfigurationTemplateIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}
	return formatter.PrintResult(networkDeviceBGPInterconnectConfigurationTemplate, &NetworkDeviceBGPInterconnectConfigurationTemplatePrintConfig)
}

func NetworkDeviceBGPInterconnectConfigurationTemplateCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating network device BGP interconnect configuration template")

	var networkDeviceBGPInterconnectConfigurationTemplateConfig sdk.CreateNetworkDeviceBGPInterconnectConfigurationTemplate
	err := utils.UnmarshalContent(config, &networkDeviceBGPInterconnectConfigurationTemplateConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkDeviceBGPInterconnectConfigurationTemplateInfo, httpRes, err := client.NetworkDeviceBGPInterconnectConfigurationTemplateAPI.
		CreateNetworkDeviceBGPInterconnectConfigurationTemplate(ctx).
		CreateNetworkDeviceBGPInterconnectConfigurationTemplate(networkDeviceBGPInterconnectConfigurationTemplateConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkDeviceBGPInterconnectConfigurationTemplateInfo, &NetworkDeviceBGPInterconnectConfigurationTemplatePrintConfig)
}

func NetworkDeviceBGPInterconnectConfigurationTemplateUpdate(ctx context.Context, networkDeviceBGPInterconnectConfigurationTemplateId string, config []byte) error {
	logger.Get().Info().Msgf("Updating network device BGP interconnect configuration template %s", networkDeviceBGPInterconnectConfigurationTemplateId)

	networkDeviceBGPInterconnectConfigurationTemplateIdNumeric, err := getNetworkDeviceBGPInterconnectConfigurationTemplateId(networkDeviceBGPInterconnectConfigurationTemplateId)
	if err != nil {
		return err
	}

	var networkDeviceBGPInterconnectConfigurationTemplateConfig sdk.UpdateNetworkDeviceBGPInterconnectConfigurationTemplate
	err = utils.UnmarshalContent(config, &networkDeviceBGPInterconnectConfigurationTemplateConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkDeviceBGPInterconnectConfigurationTemplateInfo, httpRes, err := client.NetworkDeviceBGPInterconnectConfigurationTemplateAPI.
		UpdateNetworkDeviceBGPInterconnectConfigurationTemplate(ctx, networkDeviceBGPInterconnectConfigurationTemplateIdNumeric).
		UpdateNetworkDeviceBGPInterconnectConfigurationTemplate(networkDeviceBGPInterconnectConfigurationTemplateConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkDeviceBGPInterconnectConfigurationTemplateInfo, &NetworkDeviceBGPInterconnectConfigurationTemplatePrintConfig)
}

func NetworkDeviceBGPInterconnectConfigurationTemplateDelete(ctx context.Context, networkDeviceBGPInterconnectConfigurationTemplateId string) error {
	logger.Get().Info().Msgf("Deleting network device BGP interconnect configuration template %s", networkDeviceBGPInterconnectConfigurationTemplateId)

	networkDeviceBGPInterconnectConfigurationTemplateIdNumeric, err := getNetworkDeviceBGPInterconnectConfigurationTemplateId(networkDeviceBGPInterconnectConfigurationTemplateId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceBGPInterconnectConfigurationTemplateAPI.
		DeleteNetworkDeviceBGPInterconnectConfigurationTemplate(ctx, networkDeviceBGPInterconnectConfigurationTemplateIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network device BGP interconnect configuration template %s deleted", networkDeviceBGPInterconnectConfigurationTemplateId)
	return nil
}

func getNetworkDeviceBGPInterconnectConfigurationTemplateId(networkDeviceBGPInterconnectConfigurationTemplateId string) (int64, error) {
	networkDeviceBGPInterconnectConfigurationTemplateIdNumeric, err := strconv.ParseInt(networkDeviceBGPInterconnectConfigurationTemplateId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid network device BGP interconnect configuration template ID: '%s'", networkDeviceBGPInterconnectConfigurationTemplateId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return networkDeviceBGPInterconnectConfigurationTemplateIdNumeric, nil
}
