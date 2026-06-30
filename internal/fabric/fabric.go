package fabric

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/internal/fabric_switch_config"
	"github.com/metalsoft-io/metalcloud-cli/internal/fabric_template_config"
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
	var fabricConfiguration interface{}
	switch fabricType {
	case "ethernet":
		fabricConfiguration = sdk.EthernetFabric{
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

	case "infiniband":
		fabricConfiguration = sdk.InfinibandFabric{
			FabricType:                 sdk.FABRICTYPE_INFINIBAND,
			SyslogMonitoringEnabled:    sdk.PtrBool(true),
			GnmiMonitoringEnabled:      sdk.PtrBool(false),
			ServerOnlyOperationEnabled: sdk.PtrBool(false),
			PkeyRanges:                 []string{"200-2100"},
			PreventPKeyCleanup:         []string{"1000-1100"},
			ReservedPkeys:              []string{"1-100"},
		}

	case "fibre_channel":
		fabricConfiguration = sdk.FibreChannelFabric{
			FabricType:                 sdk.FABRICTYPE_FIBRE_CHANNEL,
			SyslogMonitoringEnabled:    sdk.PtrBool(true),
			GnmiMonitoringEnabled:      sdk.PtrBool(false),
			ServerOnlyOperationEnabled: sdk.PtrBool(false),
		}

	default:
		err := fmt.Errorf("invalid fabric type: '%s'", fabricType)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(fabricConfiguration, nil)
	} else {
		return formatter.PrintYamlResult(fabricConfiguration)
	}
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

// switchImportConfig is the YAML/JSON shape consumed by FabricImportDevices:
// optional defaults deep-merged into each switch (per-switch keys win).
type switchImportConfig struct {
	Defaults map[string]interface{}   `json:"defaults" yaml:"defaults"`
	Switches []map[string]interface{} `json:"switches" yaml:"switches"`
}

// importRequiredFields are the per-switch fields required to create a device in
// practice (siteId is derived from the fabric, not supplied by the user).
var importRequiredFields = []string{
	"driver", "position", "username", "managementPassword",
	"managementAddress", "managementPort", "identifierString",
}

// FabricImportDevices bulk-imports switches into MetalSoft and attaches them to
// the fabric, idempotently. Switches that already exist (matched by management
// address, identifier or serial) are not recreated, and only devices not yet
// attached are attached. The switches' site is derived from the fabric.
// --dry-run reports the plan without writing.
func FabricImportDevices(ctx context.Context, fabricId string, config []byte, dryRun bool) error {
	logger.Get().Info().Msgf("Importing switches into fabric '%s'", fabricId)

	var importConfig switchImportConfig
	if err := utils.UnmarshalContent(config, &importConfig); err != nil {
		return err
	}
	if len(importConfig.Switches) == 0 {
		return fmt.Errorf("config must contain a non-empty 'switches' list")
	}

	fabricInfo, err := GetFabricByIdOrLabel(ctx, fabricId)
	if err != nil {
		return err
	}
	if fabricInfo.SiteId == nil {
		return fmt.Errorf("fabric '%s' (%q) has no siteId; cannot derive the site for the switches", fabricId, fabricInfo.Name)
	}
	siteId := *fabricInfo.SiteId
	fabricIdNumeric, err := utils.GetInt64FromString(fabricInfo.Id)
	if err != nil {
		return err
	}
	logger.Get().Info().Msgf("Target fabric: %q (id=%d, siteId=%d)", fabricInfo.Name, fabricIdNumeric, siteId)

	client := api.GetApiClient(ctx)

	// Existing devices in the site (idempotency) + fabric membership.
	existing, _, err := utils.FetchAllPages(client.NetworkDeviceAPI.GetNetworkDevices(ctx).FilterSiteId([]string{strconv.FormatInt(siteId, 10)}))
	if err != nil {
		return err
	}
	byMgmt := map[string]int64{}
	byIdent := map[string]int64{}
	bySerial := map[string]int64{}
	for i := range existing {
		dev := &existing[i]
		id, convErr := strconv.ParseInt(dev.Id, 10, 64)
		if convErr != nil {
			continue
		}
		if dev.ManagementAddress != "" {
			byMgmt[dev.ManagementAddress] = id
		}
		if dev.IdentifierString != "" {
			byIdent[dev.IdentifierString] = id
		}
		if dev.SerialNumber != "" {
			bySerial[dev.SerialNumber] = id
		}
	}

	fabricDevices, _, err := utils.FetchAllPages(client.NetworkFabricAPI.GetFabricNetworkDevices(ctx, fabricIdNumeric))
	if err != nil {
		return err
	}
	attached := map[int64]bool{}
	for i := range fabricDevices {
		if id, convErr := strconv.ParseInt(fabricDevices[i].Id, 10, 64); convErr == nil {
			attached[id] = true
		}
	}

	defaults := importConfig.Defaults
	created, skipped, failed, alreadyAttached := 0, 0, 0, 0
	toAttach := map[int64]bool{}

	for index, entry := range importConfig.Switches {
		merged := deepMergeStringMaps(defaults, entry)
		coerceTagsMap(merged)
		delete(merged, "siteId") // derived from the fabric, never user-supplied
		label := switchImportLabel(merged, index)

		if errs := validateSwitchMap(merged); len(errs) > 0 {
			failed++
			logger.Get().Error().Msgf("[%s] invalid: %s", label, strings.Join(errs, "; "))
			continue
		}

		if id, ok := findExistingDeviceId(merged, byMgmt, byIdent, bySerial); ok {
			skipped++
			if attached[id] {
				alreadyAttached++
				logger.Get().Info().Msgf("[%s] exists (id=%d) and already attached; skipping", label, id)
			} else {
				toAttach[id] = true
				logger.Get().Info().Msgf("[%s] exists (id=%d); will attach to fabric", label, id)
			}
			continue
		}

		if dryRun {
			created++
			logger.Get().Info().Msgf("[%s] would create network device", label)
			continue
		}

		createDevice, convErr := buildCreateNetworkDevice(merged, siteId)
		if convErr != nil {
			failed++
			logger.Get().Error().Msgf("[%s] invalid: %s", label, convErr.Error())
			continue
		}
		device, httpRes, createErr := client.NetworkDeviceAPI.CreateNetworkDevice(ctx).CreateNetworkDevice(createDevice).Execute()
		if err := response_inspector.InspectResponse(httpRes, createErr); err != nil {
			failed++
			logger.Get().Error().Msgf("[%s] create failed: %s", label, err.Error())
			continue
		}
		created++
		logger.Get().Info().Msgf("[%s] created network device id=%s", label, device.Id)
		if id, convErr := strconv.ParseInt(device.Id, 10, 64); convErr == nil {
			toAttach[id] = true
		}
	}

	// Attach phase.
	attachIds := make([]int64, 0, len(toAttach))
	for id := range toAttach {
		attachIds = append(attachIds, id)
	}
	sort.Slice(attachIds, func(i, j int) bool { return attachIds[i] < attachIds[j] })

	if len(attachIds) > 0 {
		if dryRun {
			logger.Get().Info().Msgf("Would attach %d device(s) to fabric %d: %v", len(attachIds), fabricIdNumeric, attachIds)
		} else {
			_, httpRes, attachErr := client.NetworkFabricAPI.
				AddNetworkDevicesToFabric(ctx, fabricIdNumeric).
				NetworkDevicesToFabric(sdk.NetworkDevicesToFabric{NetworkDeviceIds: attachIds}).
				Execute()
			if err := response_inspector.InspectResponse(httpRes, attachErr); err != nil {
				failed += len(attachIds)
				logger.Get().Error().Msgf("attach to fabric %d failed: %s", fabricIdNumeric, err.Error())
			} else {
				logger.Get().Info().Msgf("Attached %d device(s) to fabric %d", len(attachIds), fabricIdNumeric)
			}
		}
	} else {
		logger.Get().Info().Msgf("No devices need attaching.")
	}

	verb := "created"
	if dryRun {
		verb = "would create"
	}
	suffix := ""
	if dryRun {
		suffix = " (dry-run, no changes made)"
	}
	logger.Get().Info().Msgf(
		"Summary: %s=%d, skipped(existing)=%d, queued-for-attach=%d, already-attached=%d, failed=%d%s",
		verb, created, skipped, len(attachIds), alreadyAttached, failed, suffix)

	if failed > 0 {
		return fmt.Errorf("switch import completed with %d failure(s)", failed)
	}
	return nil
}

// deepMergeStringMaps returns a deep merge of base and override; override wins,
// nested maps merge recursively.
func deepMergeStringMaps(base, override map[string]interface{}) map[string]interface{} {
	out := map[string]interface{}{}
	for k, v := range base {
		out[k] = v
	}
	for k, v := range override {
		if existing, ok := out[k]; ok {
			if em, ok1 := existing.(map[string]interface{}); ok1 {
				if vm, ok2 := v.(map[string]interface{}); ok2 {
					out[k] = deepMergeStringMaps(em, vm)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

// coerceTagsMap stringifies scalar tagsMap values in place (the API wants string
// values; YAML parses `rack: 42` as int and `managed: true` as bool).
func coerceTagsMap(sw map[string]interface{}) {
	tags, ok := sw["tagsMap"].(map[string]interface{})
	if !ok {
		return
	}
	for key, value := range tags {
		switch v := value.(type) {
		case bool:
			if v {
				tags[key] = "true"
			} else {
				tags[key] = "false"
			}
		case string:
			// leave as-is
		case int, int64, float64:
			tags[key] = fmt.Sprintf("%v", v)
		}
	}
}

func switchImportLabel(sw map[string]interface{}, index int) string {
	for _, key := range []string{"identifierString", "managementAddress", "serialNumber"} {
		if s, ok := sw[key].(string); ok && s != "" {
			return s
		}
	}
	return fmt.Sprintf("switches[%d]", index)
}

func validateSwitchMap(sw map[string]interface{}) []string {
	var errs []string
	for _, field := range importRequiredFields {
		v, ok := sw[field]
		if !ok || v == nil || v == "" {
			errs = append(errs, fmt.Sprintf("missing required field '%s'", field))
		}
	}
	if tags, present := sw["tagsMap"]; present {
		tm, ok := tags.(map[string]interface{})
		if !ok {
			errs = append(errs, "'tagsMap' must be a mapping of string keys to string values")
		} else {
			for key, value := range tm {
				if _, isStr := value.(string); !isStr {
					errs = append(errs, fmt.Sprintf("'tagsMap.%s' must be a string value", key))
				}
			}
		}
	}
	return errs
}

func findExistingDeviceId(sw map[string]interface{}, byMgmt, byIdent, bySerial map[string]int64) (int64, bool) {
	if s, ok := sw["managementAddress"].(string); ok && s != "" {
		if id, found := byMgmt[s]; found {
			return id, true
		}
	}
	if s, ok := sw["identifierString"].(string); ok && s != "" {
		if id, found := byIdent[s]; found {
			return id, true
		}
	}
	if s, ok := sw["serialNumber"].(string); ok && s != "" {
		if id, found := bySerial[s]; found {
			return id, true
		}
	}
	return 0, false
}

// buildCreateNetworkDevice maps a merged switch entry onto a CreateNetworkDevice,
// injecting the fabric-derived siteId. It goes through JSON so the SDK's field
// tags / nullable types are honored.
func buildCreateNetworkDevice(sw map[string]interface{}, siteId int64) (sdk.CreateNetworkDevice, error) {
	raw, err := json.Marshal(sw)
	if err != nil {
		return sdk.CreateNetworkDevice{}, err
	}
	var createDevice sdk.CreateNetworkDevice
	if err := json.Unmarshal(raw, &createDevice); err != nil {
		return sdk.CreateNetworkDevice{}, err
	}
	createDevice.SiteId = sdk.PtrInt64(siteId)
	return createDevice, nil
}

// FabricRescanLinks re-scans the fabric's links, deriving the physical-link
// records from the LLDP data present once the ports are up (the "Discover links"
// action). It is idempotent: on a fabric whose links already exist it matches
// rows in place. Returns the kicked-off job.
func FabricRescanLinks(ctx context.Context, fabricId string, updateLLDPInfo bool) error {
	logger.Get().Info().Msgf("Rescan links of fabric '%s'", fabricId)

	fabricInfo, err := GetFabricByIdOrLabel(ctx, fabricId)
	if err != nil {
		return err
	}
	fabricIdNumeric, err := utils.GetInt64FromString(fabricInfo.Id)
	if err != nil {
		return err
	}

	options := sdk.NewNetworkFabricLinkRescanOptions()
	options.UpdateLLDPInfo = sdk.PtrBool(updateLLDPInfo)

	client := api.GetApiClient(ctx)

	jobInfo, httpRes, err := client.NetworkFabricAPI.
		RescanNetworkFabricLinks(ctx, fabricIdNumeric).
		NetworkFabricLinkRescanOptions(*options).
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
// to every device in a fabric. It is idempotent, with a --dry-run preview.
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

// FabricConfigureFreeform registers the base freeform device-configuration
// template + one profile per switch.
func FabricConfigureFreeform(ctx context.Context, fabricIdOrLabel string, config []byte, dryRun bool, verify bool) error {
	fabricId, err := resolveFabricNumericId(ctx, fabricIdOrLabel)
	if err != nil {
		return err
	}
	client := fabric_template_config.NewSDKClient(ctx, api.GetApiClient(ctx))
	_, err = fabric_template_config.RunFreeform(client, config, fabricId, dryRun, verify)
	return err
}

// FabricConfigureBgp registers the BGP underlay (+ l3evpn overlay/PFC/VRF)
// templates and per-switch profiles, and reconciles device customVariables.
func FabricConfigureBgp(ctx context.Context, fabricIdOrLabel string, config []byte, dryRun bool, verify bool) error {
	fabricId, err := resolveFabricNumericId(ctx, fabricIdOrLabel)
	if err != nil {
		return err
	}
	client := fabric_template_config.NewSDKClient(ctx, api.GetApiClient(ctx))
	_, err = fabric_template_config.RunBgp(client, config, fabricId, dryRun, verify)
	return err
}

// FabricConfigureFreeformExample prints a ready-to-edit freeform config example.
func FabricConfigureFreeformExample(ctx context.Context) error {
	fmt.Print(fabric_template_config.ExampleFreeformYAML())
	return nil
}

// FabricConfigureBgpExample prints a ready-to-edit bgp config example.
func FabricConfigureBgpExample(ctx context.Context) error {
	fmt.Print(fabric_template_config.ExampleBgpYAML())
	return nil
}

func resolveFabricNumericId(ctx context.Context, fabricIdOrLabel string) (int64, error) {
	fabricInfo, err := GetFabricByIdOrLabel(ctx, fabricIdOrLabel)
	if err != nil {
		return 0, err
	}
	fabricId, err := strconv.ParseInt(fabricInfo.Id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid fabric ID %q: %w", fabricInfo.Id, err)
	}
	return fabricId, nil
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
