package fabric

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/internal/fabric_switch_config"
	"github.com/metalsoft-io/metalcloud-cli/internal/network_device"
	"github.com/metalsoft-io/metalcloud-cli/internal/site"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var fabricPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Name": {
			Order: 2,
		},
		"Description": {
			Order: 3,
		},
		"SiteId": {
			Title: "Site",
			Order: 4,
		},
		"Status": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       5,
		},
		"FabricConfiguration.EthernetFabric.FabricType|FabricConfiguration.InfinibandFabric.FabricType|FabricConfiguration.FibreChannelFabric.FabricType": {
			Title: "Type",
			Order: 6,
		},
	},
}

func FabricList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all fabrics")

	client := api.GetApiClient(ctx)

	request := client.NetworkFabricAPI.GetNetworkFabrics(ctx).SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &fabricPrintConfig)
}

func FabricGet(ctx context.Context, fabricId string) error {
	logger.Get().Info().Msgf("Get fabric '%s'", fabricId)

	fabricInfo, err := GetFabricByIdOrLabel(ctx, fabricId)
	if err != nil {
		return err
	}

	return formatter.PrintResult(fabricInfo, &fabricPrintConfig)
}

func FabricConfigExample(ctx context.Context, fabricType string) error {
	var fabricConfiguration sdk.NetworkFabricFabricConfiguration
	switch fabricType {
	case "ethernet":
		ethernetConfig := sdk.EthernetFabric{
			FabricType:                 sdk.FABRICTYPE_ETHERNET,
			SyslogMonitoringEnabled:    sdk.PtrBool(true),
			GnmiMonitoringEnabled:      sdk.PtrBool(false),
			ServerOnlyOperationEnabled: sdk.PtrBool(false),
			LagRanges:                  []string{"100-200", "300-400"},
			MlagRanges:                 []string{"30-40", "50-60"},
			VlanRanges:                 []string{"3000-3100", "2000-2100"},
			ReservedVlans:              []string{"2000-2100", "2200-2300"},
			PreventVlanCleanup:         []string{"1000-1100"},
			VniPrefix:                  sdk.PtrInt32(5000),
			L3VniPrefix:                sdk.PtrInt32(9000),
			VrfVlanRanges:              []string{"400-450", "460-470"},
		}

		fabricConfiguration = sdk.NetworkFabricFabricConfiguration{
			EthernetFabric: &ethernetConfig,
		}
	case "infiniband":
		infinibandConfig := sdk.InfinibandFabric{
			FabricType:                 sdk.FABRICTYPE_INFINIBAND,
			SyslogMonitoringEnabled:    sdk.PtrBool(true),
			GnmiMonitoringEnabled:      sdk.PtrBool(false),
			ServerOnlyOperationEnabled: sdk.PtrBool(false),
			PkeyRanges:                 []string{"200-2100"},
			PreventPKeyCleanup:         []string{"1000-1100"},
			ReservedPkeys:              []string{"1-100"},
		}

		fabricConfiguration = sdk.NetworkFabricFabricConfiguration{
			InfinibandFabric: &infinibandConfig,
		}
	case "fibre_channel":
		fcConfig := sdk.FibreChannelFabric{
			FabricType:                 sdk.FABRICTYPE_FIBRE_CHANNEL,
			SyslogMonitoringEnabled:    sdk.PtrBool(true),
			GnmiMonitoringEnabled:      sdk.PtrBool(false),
			ServerOnlyOperationEnabled: sdk.PtrBool(false),
		}

		fabricConfiguration = sdk.NetworkFabricFabricConfiguration{
			FibreChannelFabric: &fcConfig,
		}
	default:
		err := fmt.Errorf("invalid fabric type: '%s'", fabricType)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	return formatter.PrintResult(fabricConfiguration, nil)
}

func FabricCreate(ctx context.Context, siteIdOrLabel string, fabricName string, fabricType string, description string, config []byte) error {
	logger.Get().Info().Msgf("Create fabric '%s'", fabricName)

	site, err := site.GetSiteByIdOrLabel(ctx, siteIdOrLabel)
	if err != nil {
		return err
	}

	var fabricConfiguration sdk.NetworkFabricFabricConfiguration
	switch fabricType {
	case "ethernet":
		ethernetConfig := sdk.EthernetFabric{}
		err := utils.UnmarshalContent(config, &ethernetConfig)
		if err != nil {
			return err
		}

		fabricConfiguration = sdk.NetworkFabricFabricConfiguration{
			EthernetFabric: &ethernetConfig,
		}
	case "infiniband":
		infinibandConfig := sdk.InfinibandFabric{}
		err := utils.UnmarshalContent(config, &infinibandConfig)
		if err != nil {
			return err
		}

		fabricConfiguration = sdk.NetworkFabricFabricConfiguration{
			InfinibandFabric: &infinibandConfig,
		}
	case "fibre_channel":
		fcConfig := sdk.FibreChannelFabric{}
		err := utils.UnmarshalContent(config, &fcConfig)
		if err != nil {
			return err
		}

		fabricConfiguration = sdk.NetworkFabricFabricConfiguration{
			FibreChannelFabric: &fcConfig,
		}
	default:
		err := fmt.Errorf("invalid fabric type: '%s'", fabricType)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	createFabric := sdk.CreateNetworkFabric{
		Name:                fabricName,
		Description:         sdk.PtrString(description),
		SiteId:              sdk.PtrInt64(site.Id),
		FabricConfiguration: fabricConfiguration,
	}

	client := api.GetApiClient(ctx)

	fabricInfo, httpRes, err := client.NetworkFabricAPI.CreateNetworkFabric(ctx).CreateNetworkFabric(createFabric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(fabricInfo, &fabricPrintConfig)
}

func FabricUpdate(ctx context.Context, fabricId string, fabricName string, description string, config []byte) error {
	logger.Get().Info().Msgf("Update fabric '%s'", fabricId)

	fabricInfo, err := GetFabricByIdOrLabel(ctx, fabricId)
	if err != nil {
		return err
	}
	fabricIdNumber, err := utils.GetInt64FromString(fabricInfo.Id)
	if err != nil {
		return err
	}

	var fabricConfiguration sdk.NetworkFabricFabricConfiguration
	if fabricInfo.FabricConfiguration.EthernetFabric != nil {
		ethernetConfig := sdk.EthernetFabric{}
		err := utils.UnmarshalContent(config, &ethernetConfig)
		if err != nil {
			return err
		}

		fabricConfiguration = sdk.NetworkFabricFabricConfiguration{
			EthernetFabric: &ethernetConfig,
		}
	} else if fabricInfo.FabricConfiguration.InfinibandFabric != nil {
		infinibandConfig := sdk.InfinibandFabric{}
		err := utils.UnmarshalContent(config, &infinibandConfig)
		if err != nil {
			return err
		}

		fabricConfiguration = sdk.NetworkFabricFabricConfiguration{
			InfinibandFabric: &infinibandConfig,
		}
	} else if fabricInfo.FabricConfiguration.FibreChannelFabric != nil {
		fcConfig := sdk.FibreChannelFabric{}
		err := utils.UnmarshalContent(config, &fcConfig)
		if err != nil {
			return err
		}

		fabricConfiguration = sdk.NetworkFabricFabricConfiguration{
			FibreChannelFabric: &fcConfig,
		}
	} else {
		err := fmt.Errorf("invalid fabric type")
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	updateFabric := sdk.UpdateNetworkFabric{
		Name:                sdk.PtrString(fabricName),
		Description:         sdk.PtrString(description),
		SiteId:              fabricInfo.SiteId,
		FabricConfiguration: fabricConfiguration,
	}

	client := api.GetApiClient(ctx)

	fabricInfoUpdated, httpRes, err := client.NetworkFabricAPI.UpdateNetworkFabric(ctx, fabricIdNumber).
		UpdateNetworkFabric(updateFabric).
		IfMatch(fabricInfo.Revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(fabricInfoUpdated, &fabricPrintConfig)
}

func FabricActivate(ctx context.Context, fabricId string) error {
	logger.Get().Info().Msgf("Activate fabric '%s'", fabricId)

	fabricIdNumeric, err := utils.GetInt64FromString(fabricId)
	if err != nil {
		return err
	}

	fabricInfo, err := GetFabricByIdOrLabel(ctx, fabricId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	fabricInfo, httpRes, err := client.NetworkFabricAPI.
		ActivateNetworkFabric(ctx, fabricIdNumeric).
		IfMatch(fabricInfo.Revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(fabricInfo, &fabricPrintConfig)
}

func FabricDeploy(ctx context.Context, fabricId string) error {
	logger.Get().Info().Msgf("Deploy fabric '%s'", fabricId)

	fabricIdNumeric, err := utils.GetInt64FromString(fabricId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	jobInfo, httpRes, err := client.NetworkFabricAPI.
		DeployNetworkFabric(ctx, fabricIdNumeric).
		NetworkFabricDeployOptions(*sdk.NewNetworkFabricDeployOptions(false)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(jobInfo, nil)
}

func FabricDevicesGet(ctx context.Context, fabricId string) error {
	logger.Get().Info().Msgf("Get fabric '%s' devices", fabricId)

	fabricIdNumeric, err := utils.GetInt64FromString(fabricId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	devicesList, httpRes, err := client.NetworkFabricAPI.GetFabricNetworkDevices(ctx, fabricIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(devicesList.Data, &network_device.NetworkDevicePrintConfig)
}

func FabricDevicesAdd(ctx context.Context, fabricId string, deviceIds []string) error {
	logger.Get().Info().Msgf("Adding devices '%v' to fabric '%s'", deviceIds, fabricId)

	fabricInfo, err := GetFabricByIdOrLabel(ctx, fabricId)
	if err != nil {
		return err
	}

	fabricIdNumeric, err := utils.GetInt64FromString(fabricInfo.Id)
	if err != nil {
		return err
	}

	deviceIdsNumeric := make([]int64, 0)
	for _, deviceId := range deviceIds {
		device, err := network_device.GetNetworkDeviceById(ctx, deviceId)
		if err != nil {
			return err
		}

		if *fabricInfo.SiteId != device.SiteId {
			err := fmt.Errorf("device '%s' is not in the same site as fabric '%s'", deviceId, fabricId)
			logger.Get().Error().Err(err).Msg("")
			return err
		}

		deviceIdNumeric, err := utils.GetInt64FromString(device.Id)
		if err != nil {
			return err
		}

		deviceIdsNumeric = append(deviceIdsNumeric, deviceIdNumeric)
	}

	client := api.GetApiClient(ctx)

	_, httpRes, err := client.NetworkFabricAPI.AddNetworkDevicesToFabric(ctx, fabricIdNumeric).
		NetworkDevicesToFabric(sdk.NetworkDevicesToFabric{NetworkDeviceIds: deviceIdsNumeric}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return nil
}

func FabricDevicesRemove(ctx context.Context, fabricId string, deviceId string) error {
	logger.Get().Info().Msgf("Removing device '%s' from fabric '%s'", deviceId, fabricId)

	fabricInfo, err := GetFabricByIdOrLabel(ctx, fabricId)
	if err != nil {
		return err
	}

	fabricIdNumeric, err := utils.GetInt64FromString(fabricInfo.Id)
	if err != nil {
		return err
	}

	device, err := network_device.GetNetworkDeviceById(ctx, deviceId)
	if err != nil {
		return err
	}

	deviceIdNumeric, err := utils.GetInt64FromString(device.Id)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, httpRes, err := client.NetworkFabricAPI.RemoveNetworkDeviceFromFabric(ctx, fabricIdNumeric, deviceIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return nil
}

// FabricConfigureSwitches applies a declarative fabric-switch configuration
// (hostnames, ASNs, loopbacks, port enable/descriptions, point-to-point links)
// to every device in a fabric. It is the CLI port of the standalone
// configure_switches.py script: idempotent, with a --dry-run preview.
func FabricConfigureSwitches(ctx context.Context, fabricIdOrLabel string, config []byte, dryRun bool) error {
	fabricInfo, err := GetFabricByIdOrLabel(ctx, fabricIdOrLabel)
	if err != nil {
		return err
	}
	fabricId, err := strconv.ParseInt(fabricInfo.Id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid fabric ID %q: %w", fabricInfo.Id, err)
	}

	cfg, err := fabric_switch_config.LoadConfig(config)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)
	switchClient := fabric_switch_config.NewSDKClient(ctx, client)

	if dryRun {
		logger.Get().Info().Msgf("Dry run: computing the plan for fabric %d without writing", fabricId)
	}

	result, err := fabric_switch_config.Configure(switchClient, cfg, fabricId, dryRun)
	if err != nil {
		return err
	}

	for _, w := range result.Warnings {
		logger.Get().Warn().Msg(w)
	}

	keys := make([]string, 0, len(result.Counters))
	for k := range result.Counters {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	summary := ""
	for _, k := range keys {
		if summary != "" {
			summary += ", "
		}
		summary += fmt.Sprintf("%s=%d", k, result.Counters[k])
	}
	if summary == "" {
		summary = "nothing to do"
	}
	suffix := ""
	if dryRun {
		suffix = " (dry-run, no changes made)"
	}
	logger.Get().Info().Msgf("Summary: %s, failures=%d%s", summary, result.Failures, suffix)

	if result.Failures > 0 {
		return fmt.Errorf("fabric switch configuration completed with %d failure(s)", result.Failures)
	}
	return nil
}

// FabricConfigureSwitchesExample prints a commented, ready-to-edit example of
// the fabric switch configuration accepted by FabricConfigureSwitches. The
// output is valid YAML and can be piped straight into the configure-switches
// command's --config-source.
func FabricConfigureSwitchesExample(ctx context.Context) error {
	fmt.Print(fabric_switch_config.ExampleConfigYAML())
	return nil
}

func GetFabricByIdOrLabel(ctx context.Context, fabricIdOrLabel string) (*sdk.NetworkFabric, error) {
	client := api.GetApiClient(ctx)

	fabricIdNumber, err := utils.GetInt64FromString(fabricIdOrLabel)
	if err == nil {
		fabricInfo, httpRes, err := client.NetworkFabricAPI.GetNetworkFabricById(ctx, fabricIdNumber).Execute()
		if err = response_inspector.InspectResponse(httpRes, err); err == nil {
			return fabricInfo, nil
		}
	}

	fabrics, httpRes, err := client.NetworkFabricAPI.
		GetNetworkFabrics(ctx).
		FilterName([]string{fabricIdOrLabel}).
		Execute()
	if err = response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	if len(fabrics.Data) == 0 {
		err := fmt.Errorf("fabric '%s' not found", fabricIdOrLabel)
		logger.Get().Error().Err(err).Msg("")
		return nil, err
	}

	return &fabrics.Data[0], nil
}
