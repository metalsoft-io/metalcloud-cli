package network_device

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var NetworkDevicePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"IdentifierString": {
			Title:    "Identifier",
			MaxWidth: 40,
			Order:    2,
		},
		"SiteId": {
			Title: "Site",
			Order: 3,
		},
		"ManagementAddress": {
			Title: "Address",
			Order: 4,
		},
		"ManagementMacAddress": {
			Title: "MAC",
			Order: 5,
		},
		"SerialNumber": {
			Title: "Serial",
			Order: 6,
		},
		"Driver": {
			Order: 7,
		},
		"Status": {
			Order:       8,
			Transformer: formatter.FormatStatusValue,
		},
	},
}

func NetworkDeviceList(ctx context.Context, filterStatus []string) error {
	logger.Get().Info().Msgf("Listing all network devices")

	client := api.GetApiClient(ctx)

	request := client.NetworkDeviceAPI.GetNetworkDevices(ctx)
	if len(filterStatus) > 0 {
		request = request.FilterStatus(utils.ProcessFilterStringSlice(filterStatus))
	}
	request = request.SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &NetworkDevicePrintConfig)
}

func NetworkDeviceGet(ctx context.Context, networkDeviceId string) error {
	logger.Get().Info().Msgf("Get network device %s details", networkDeviceId)

	networkDevice, err := GetNetworkDeviceById(ctx, networkDeviceId)
	if err != nil {
		return err
	}

	return formatter.PrintResult(networkDevice, &NetworkDevicePrintConfig)
}

func NetworkDeviceConfigExample(ctx context.Context) error {
	networkDeviceConfiguration := sdk.CreateNetworkDevice{
		SiteId:           sdk.PtrInt64(1),
		Driver:           sdk.NETWORKDEVICEDRIVER_SONIC_ENTERPRISE,
		IdentifierString: sdk.PtrString("example"),
		SerialNumber:     sdk.PtrString("1234567890"),
		ChassisRackId:    sdk.PtrInt64(1),
		Position:         "leaf",
		IsGateway:        sdk.PtrBool(false),
		IsStorageSwitch:  sdk.PtrBool(false),
		IsBorderDevice:   sdk.PtrBool(false),
	}

	networkDeviceConfiguration.ManagementAddress.Set(sdk.PtrString("1.1.1.1"))
	networkDeviceConfiguration.ManagementPort.Set(sdk.PtrInt32(22))
	networkDeviceConfiguration.Username.Set(sdk.PtrString("admin"))
	networkDeviceConfiguration.ManagementPassword = "password"

	networkDeviceConfiguration.SyslogEnabled.Set(sdk.PtrBool(true))

	networkDeviceConfiguration.ManagementAddressGateway.Set(sdk.PtrString("1.1.1.1"))
	networkDeviceConfiguration.ManagementAddressMask.Set(sdk.PtrString("255.255.255.0"))
	networkDeviceConfiguration.ManagementMAC.Set(sdk.PtrString("AA:BB:CC:DD:EE:FF"))
	networkDeviceConfiguration.ChassisIdentifier.Set(sdk.PtrString("example"))
	networkDeviceConfiguration.LoopbackAddress.Set(sdk.PtrString("127.0.0.1"))
	networkDeviceConfiguration.VtepAddress.Set(nil)
	networkDeviceConfiguration.Asn.Set(sdk.PtrInt64(65000))

	networkDeviceConfiguration.AuthenticationOptions = []sdk.NetworkDeviceAuthOption{
		{Kind: "tacacs", DeviceAuthProviderId: sdk.PtrInt64(1)},
		{Kind: "local"},
	}

	return formatter.PrintResult(networkDeviceConfiguration, nil)
}

func NetworkDeviceCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating network device")

	var networkDeviceConfig sdk.CreateNetworkDevice
	err := utils.UnmarshalContent(config, &networkDeviceConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkDeviceInfo, httpRes, err := client.NetworkDeviceAPI.CreateNetworkDevice(ctx).CreateNetworkDevice(networkDeviceConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkDeviceInfo, &NetworkDevicePrintConfig)
}

func NetworkDeviceCreateBulk(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating network devices in bulk")

	var devicesConfig []sdk.CreateNetworkDevice
	err := utils.UnmarshalContent(config, &devicesConfig)
	if err != nil {
		return err
	}

	if len(devicesConfig) == 0 {
		return fmt.Errorf("no network devices found in configuration")
	}

	client := api.GetApiClient(ctx)

	// Track results for reporting
	results := make([]interface{}, 0)
	errors := make([]error, 0)

	logger.Get().Info().Msgf("Creating %d network devices", len(devicesConfig))

	// Process each network device
	for i, deviceConfig := range devicesConfig {
		label := networkDeviceConfigLabel(deviceConfig, i)

		networkDeviceInfo, httpRes, err := client.NetworkDeviceAPI.CreateNetworkDevice(ctx).CreateNetworkDevice(deviceConfig).Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			logger.Get().Error().Msgf("Failed to create network device %d: %s", i+1, err)
			errors = append(errors, fmt.Errorf("network device %d (%s): %s", i+1, label, err))
			continue
		}

		results = append(results, networkDeviceInfo)
		logger.Get().Info().Msgf("Created network device %d: %s", i+1, label)
	}

	// Print summary
	logger.Get().Info().Msgf("Bulk network device creation complete: %d created, %d failed", len(results), len(errors))

	// Print any errors that occurred
	errorsText := ""
	if len(errors) > 0 {
		logger.Get().Error().Msgf("Errors encountered during bulk creation:")
		for _, err := range errors {
			logger.Get().Error().Msgf("  - %s", err)
			errorsText += fmt.Sprintf("\n  - %s", err)
		}
	}

	// Print the successfully created network devices
	if len(results) > 0 {
		err = formatter.PrintResult(results, &NetworkDevicePrintConfig)
	}

	if len(errors) > 0 || err != nil {
		if err != nil {
			errorsText += fmt.Sprintf("\n  - %s", err)
		}
		return fmt.Errorf("bulk network device creation completed with errors: %s", errorsText)
	}

	return nil
}

// networkDeviceConfigLabel returns a human-friendly identifier for a device
// config, used only in bulk-creation log lines and error messages.
func networkDeviceConfigLabel(device sdk.CreateNetworkDevice, index int) string {
	if device.IdentifierString != nil && *device.IdentifierString != "" {
		return *device.IdentifierString
	}
	if device.ManagementAddress.IsSet() && device.ManagementAddress.Get() != nil && *device.ManagementAddress.Get() != "" {
		return *device.ManagementAddress.Get()
	}
	if device.SerialNumber != nil && *device.SerialNumber != "" {
		return *device.SerialNumber
	}
	return fmt.Sprintf("device[%d]", index)
}

func NetworkDeviceUpdate(ctx context.Context, networkDeviceId string, config []byte) error {
	logger.Get().Info().Msgf("Updating network device")

	networkDeviceIdNumeric, revision, err := getNetworkDeviceIdAndRevision(ctx, networkDeviceId)
	if err != nil {
		return err
	}

	var networkDeviceConfig sdk.UpdateNetworkDevice
	err = utils.UnmarshalContent(config, &networkDeviceConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkDeviceInfo, httpRes, err := client.NetworkDeviceAPI.
		UpdateNetworkDevice(ctx, networkDeviceIdNumeric).
		UpdateNetworkDevice(networkDeviceConfig).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkDeviceInfo, &NetworkDevicePrintConfig)
}

func NetworkDeviceDelete(ctx context.Context, networkDeviceId string) error {
	logger.Get().Info().Msgf("Deleting network device %s", networkDeviceId)

	networkDeviceIdNumeric, err := GetNetworkDeviceId(networkDeviceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceAPI.
		DeleteNetworkDevice(ctx, networkDeviceIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network device %s deleting in progress.", networkDeviceId)
	return nil
}

func NetworkDeviceArchive(ctx context.Context, networkDeviceId string) error {
	logger.Get().Info().Msgf("Archiving network device %s", networkDeviceId)

	networkDeviceIdNumeric, revision, err := getNetworkDeviceIdAndRevision(ctx, networkDeviceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceAPI.
		ArchiveNetworkDevice(ctx, networkDeviceIdNumeric).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network device %s archiving in progress.", networkDeviceId)
	return nil
}

// DiscoveryTargets lists the discovery types a network device supports. Passing
// all of them (the default) performs a full discovery.
var DiscoveryTargets = []string{"hardware", "software", "ports"}

// NetworkDeviceDiscover initiates discovery for a network device. targets selects
// which discovery types to run ("hardware", "software", "ports"); an empty slice
// runs all of them. Discovered data is persisted to the device inventory.
func NetworkDeviceDiscover(ctx context.Context, networkDeviceId string, targets []string) error {
	networkDeviceIdNumeric, err := GetNetworkDeviceId(networkDeviceId)
	if err != nil {
		return err
	}

	if len(targets) == 0 {
		targets = DiscoveryTargets
	}
	for _, target := range targets {
		if !slices.Contains(DiscoveryTargets, target) {
			return fmt.Errorf("invalid discovery target '%s'; valid targets are: %s",
				target, strings.Join(DiscoveryTargets, ", "))
		}
	}

	logger.Get().Info().Msgf("Discovering network device %s (%s)", networkDeviceId, strings.Join(targets, ", "))

	discoveryQuery := sdk.DiscoveryQuery{
		Discover:    targets,
		PersistData: true,
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceAPI.
		DiscoverNetworkDevice(ctx, networkDeviceIdNumeric).
		DiscoveryQuery(discoveryQuery).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network device %s discovery initiated", networkDeviceId)
	return nil
}

func NetworkDeviceGetCredentials(ctx context.Context, networkDeviceId string) error {
	logger.Get().Info().Msgf("Getting network device %s credentials", networkDeviceId)

	networkDeviceIdNumeric, err := GetNetworkDeviceId(networkDeviceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, httpRes, err := client.NetworkDeviceAPI.
		GetNetworkDeviceCredentials(ctx, networkDeviceIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	// Parse response body to display credentials
	credentialsMap, err := response_inspector.ParseResponseBody(httpRes)
	if err != nil {
		return err
	}

	return formatter.PrintResult(credentialsMap, nil)
}

func NetworkDeviceGetPorts(ctx context.Context, networkDeviceId string) error {
	logger.Get().Info().Msgf("Getting network device %s ports", networkDeviceId)

	networkDeviceIdNumeric, err := GetNetworkDeviceId(networkDeviceId)
	if err != nil {
		return err
	}

	portsInfo, err := GetNetworkDevicePorts(ctx, float32(networkDeviceIdNumeric))
	if err != nil {
		return err
	}

	return formatter.PrintResult(portsInfo, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"InterfaceId": {
				Title: "ID",
				Order: 1,
			},
			"InterfaceName": {
				Title: "Name",
				Order: 2,
			},
			"InterfaceDescription": {
				Title:    "Description",
				MaxWidth: 30,
				Order:    3,
			},
			"Kind": {
				Order: 4,
			},
			"MacAddress": {
				Title: "MAC",
				Order: 5,
			},
			"LagIdentifier": {
				Title: "LAG",
				Order: 6,
			},
			"Tags": {
				Order: 7,
			},
		},
	})
}

func NetworkDeviceSetPortStatus(ctx context.Context, networkDeviceId string, portId string, action string) error {
	logger.Get().Info().Msgf("Setting port status for network device %s port %s to %s", networkDeviceId, portId, action)

	networkDeviceIdNumeric, err := GetNetworkDeviceId(networkDeviceId)
	if err != nil {
		return err
	}

	if action != "up" && action != "down" {
		return fmt.Errorf("invalid port action: '%s'. Valid actions are: up, down", action)
	}

	client := api.GetApiClient(ctx)

	portStatus := sdk.NetworkDevicePortStatus{
		Ports:  []string{portId},
		Status: action == "up",
	}

	httpRes, err := client.NetworkDeviceAPI.
		SetNetworkDevicePortStatus(ctx, networkDeviceIdNumeric).
		NetworkDevicePortStatus(portStatus).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Port %s status for network device %s set to %s", portId, networkDeviceId, action)
	return nil
}

var networkDevicePortConfigPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Enabled": {
			Title: "Enabled",
			Order: 1,
		},
		"Description": {
			Title: "Description",
			Order: 2,
		},
		"Mtu": {
			Title: "MTU",
			Order: 3,
		},
		"Revision": {
			Title: "Revision",
			Order: 4,
		},
	},
}

var networkDevicePortIpPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"InterfaceId": {
			Title: "Interface",
			Order: 2,
		},
		"Address": {
			Title: "Address",
			Order: 3,
		},
		"PrefixLength": {
			Title: "Prefix",
			Order: 4,
		},
		"ServiceStatus": {
			Title: "Status",
			Order: 5,
		},
	},
}

// revisionMismatchRe extracts the server-expected revision out of an optimistic
// locking 409 body ("... found 7 ..."), used to retry the staged port-IP POST.
var revisionMismatchRe = regexp.MustCompile(`found (\d+)`)

// NetworkDeviceUpdatePortConfig patches the staged config (enabled flag and/or
// description) of a single network device port. The port is addressed by its
// numeric interface id. The current config revision is read first and sent as
// If-Match for optimistic concurrency.
func NetworkDeviceUpdatePortConfig(ctx context.Context, networkDeviceId string, portId string, enabled *bool, description *string) error {
	logger.Get().Info().Msgf("Updating port %s config of network device %s", portId, networkDeviceId)

	networkDeviceIdNumeric, err := GetNetworkDeviceId(networkDeviceId)
	if err != nil {
		return err
	}

	portIdNumeric, err := getNetworkDevicePortId(portId)
	if err != nil {
		return err
	}

	if enabled == nil && description == nil {
		return fmt.Errorf("nothing to update: specify --enabled and/or --description")
	}

	client := api.GetApiClient(ctx)

	port, httpRes, err := client.NetworkDeviceAPI.
		GetNetworkDevicePort(ctx, networkDeviceIdNumeric, portIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	configUpdate := sdk.UpdateNetworkEquipmentInterfaceConfig{}
	if enabled != nil {
		configUpdate.Enabled = *sdk.NewNullableBool(enabled)
	}
	if description != nil {
		configUpdate.Description = *sdk.NewNullableString(description)
	}

	revision := strconv.FormatInt(port.Config.Revision, 10)

	configInfo, httpRes, err := client.NetworkDeviceAPI.
		UpdateNetworkDevicePortConfig(ctx, networkDeviceIdNumeric, portIdNumeric).
		UpdateNetworkEquipmentInterfaceConfig(configUpdate).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(configInfo, &networkDevicePortConfigPrintConfig)
}

// NetworkDeviceAddPortIp stages a new IP address on a network device port
// (addressed by its numeric interface id), e.g. a /32 loopback address.
//
// The POST is guarded by optimistic locking, but the interface entity revision
// is not served directly: the single-port GET exposes a zero-based
// config.revision while the lock checks a one-based counter. We therefore send
// config.revision + 1, and if the server still rejects with a 409 naming a
// different current revision, retry once with the value it expects.
func NetworkDeviceAddPortIp(ctx context.Context, networkDeviceId string, portId string, family string, address string, prefixLength int32) error {
	logger.Get().Info().Msgf("Adding %s/%d to port %s of network device %s", address, prefixLength, portId, networkDeviceId)

	networkDeviceIdNumeric, err := GetNetworkDeviceId(networkDeviceId)
	if err != nil {
		return err
	}

	portIdNumeric, err := getNetworkDevicePortId(portId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	port, httpRes, err := client.NetworkDeviceAPI.
		GetNetworkDevicePort(ctx, networkDeviceIdNumeric, portIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	payload := sdk.AddNetworkEquipmentInterfaceIp{
		Address:      address,
		PrefixLength: prefixLength,
	}

	revision := strconv.FormatInt(port.Config.Revision+1, 10)

	ipInfo, httpRes, err := client.NetworkDeviceAPI.
		AddNetworkDevicePortIp(ctx, networkDeviceIdNumeric, portIdNumeric, family).
		AddNetworkEquipmentInterfaceIp(payload).
		IfMatch(revision).
		Execute()

	// Optimistic-lock retry: the server reports the revision it actually expects.
	if err != nil && httpRes != nil && httpRes.StatusCode == 409 {
		if expected := expectedRevisionFromError(err); expected != "" {
			ipInfo, httpRes, err = client.NetworkDeviceAPI.
				AddNetworkDevicePortIp(ctx, networkDeviceIdNumeric, portIdNumeric, family).
				AddNetworkEquipmentInterfaceIp(payload).
				IfMatch(expected).
				Execute()
		}
	}
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(ipInfo, &networkDevicePortIpPrintConfig)
}

func getNetworkDevicePortId(portId string) (int64, error) {
	portIdNumeric, err := strconv.ParseInt(portId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid port (interface) ID: '%s'", portId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return portIdNumeric, nil
}

func expectedRevisionFromError(err error) string {
	var apiErr sdk.GenericOpenAPIError
	if errors.As(err, &apiErr) {
		if m := revisionMismatchRe.FindSubmatch(apiErr.Body()); m != nil {
			return string(m[1])
		}
	}
	return ""
}

func NetworkDeviceReset(ctx context.Context, networkDeviceId string) error {
	logger.Get().Info().Msgf("Resetting network device %s", networkDeviceId)

	networkDeviceIdNumeric, err := GetNetworkDeviceId(networkDeviceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceAPI.
		ResetNetworkDevice(ctx, networkDeviceIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network device %s reset initiated", networkDeviceId)
	return nil
}

func NetworkDeviceSetFailed(ctx context.Context, networkDeviceId string) error {
	logger.Get().Info().Msgf("Changing network device %s status to failed", networkDeviceId)

	networkDeviceIdNumeric, eTag, err := getNetworkDeviceIdAndRevision(ctx, networkDeviceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, httpRes, err := client.NetworkDeviceAPI.
		SetNetworkDeviceAsFailed(ctx, networkDeviceIdNumeric).
		IfMatch(eTag).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network device %s status changed to failed", networkDeviceId)
	return nil
}

func NetworkDeviceEnableSyslog(ctx context.Context, networkDeviceId string) error {
	logger.Get().Info().Msgf("Enabling syslog for network device %s", networkDeviceId)

	networkDeviceIdNumeric, err := GetNetworkDeviceId(networkDeviceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, httpRes, err := client.NetworkDeviceAPI.
		EnableNetworkDeviceSyslog(ctx, networkDeviceIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Syslog enabled for network device %s", networkDeviceId)
	return nil
}

func GetNetworkDeviceById(ctx context.Context, networkDeviceId string) (*sdk.NetworkDevice, error) {
	networkDeviceIdNumeric, err := GetNetworkDeviceId(networkDeviceId)
	if err != nil {
		return nil, err
	}

	client := api.GetApiClient(ctx)

	networkDevice, httpRes, err := client.NetworkDeviceAPI.GetNetworkDevice(ctx, networkDeviceIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	return networkDevice, nil
}

func GetNetworkDeviceByName(ctx context.Context, siteName string, networkDeviceName string) (*sdk.NetworkDevice, error) {
	client := api.GetApiClient(ctx)

	request := client.NetworkDeviceAPI.GetNetworkDevices(ctx)

	request = request.FilterDatacenterName([]string{siteName})
	request = request.FilterIdentifierString([]string{networkDeviceName})

	networkDevice, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	return &networkDevice.Data[0], nil
}

// GetNetworkDeviceByIdOrLabel resolves a network device by its numeric id or by
// its identifierString (the switch hostname/label). A purely numeric ref is
// treated as an id; anything else is matched exactly against identifierString.
func GetNetworkDeviceByIdOrLabel(ctx context.Context, deviceRef string) (*sdk.NetworkDevice, error) {
	return getNetworkDeviceByIdOrLabel(ctx, deviceRef, nil)
}

// GetNetworkDeviceByIdOrLabelInSite behaves like GetNetworkDeviceByIdOrLabel but
// scopes a label lookup to a single site. Switch labels (e.g. "leaf-su00-r0")
// are reused across sites, so scoping to the caller's site disambiguates them.
// A numeric ref still resolves directly by id, independent of the site.
func GetNetworkDeviceByIdOrLabelInSite(ctx context.Context, deviceRef string, siteId int64) (*sdk.NetworkDevice, error) {
	return getNetworkDeviceByIdOrLabel(ctx, deviceRef, &siteId)
}

func getNetworkDeviceByIdOrLabel(ctx context.Context, deviceRef string, siteId *int64) (*sdk.NetworkDevice, error) {
	client := api.GetApiClient(ctx)

	if idNumeric, err := strconv.ParseInt(deviceRef, 10, 64); err == nil {
		networkDevice, httpRes, err := client.NetworkDeviceAPI.GetNetworkDevice(ctx, idNumeric).Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			return nil, err
		}
		return networkDevice, nil
	}

	request := client.NetworkDeviceAPI.GetNetworkDevices(ctx).
		FilterIdentifierString([]string{deviceRef})
	if siteId != nil {
		request = request.FilterSiteId([]string{strconv.FormatInt(*siteId, 10)})
	}

	list, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	matches := make([]sdk.NetworkDevice, 0, len(list.Data))
	for _, nd := range list.Data {
		if nd.IdentifierString == deviceRef {
			matches = append(matches, nd)
		}
	}

	scope := ""
	if siteId != nil {
		scope = fmt.Sprintf(" in site %d", *siteId)
	}

	switch len(matches) {
	case 0:
		return nil, fmt.Errorf("no network device found with label '%s'%s", deviceRef, scope)
	case 1:
		return &matches[0], nil
	default:
		return nil, fmt.Errorf("multiple network devices found with label '%s'%s; use the numeric id instead", deviceRef, scope)
	}
}

func GetNetworkDevicePorts(ctx context.Context, networkDeviceId float32) ([]sdk.NetworkDeviceInterface, error) {
	client := api.GetApiClient(ctx)

	portsInfo, httpRes, err := client.NetworkDeviceAPI.
		GetNetworkDevicePorts(ctx, networkDeviceId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	return portsInfo.Data, nil
}

func GetNetworkDeviceId(networkDeviceId string) (int64, error) {
	networkDeviceIdNumeric, err := strconv.ParseInt(networkDeviceId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid network device ID: '%s'", networkDeviceId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return networkDeviceIdNumeric, nil
}

func getNetworkDeviceIdAndRevision(ctx context.Context, networkDeviceId string) (int64, string, error) {
	networkDeviceIdNumeric, err := GetNetworkDeviceId(networkDeviceId)
	if err != nil {
		return 0, "", err
	}

	client := api.GetApiClient(ctx)

	networkDevice, httpRes, err := client.NetworkDeviceAPI.GetNetworkDevice(ctx, networkDeviceIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}

	return networkDeviceIdNumeric, strconv.Itoa(int(networkDevice.Revision)), nil
}

var networkDeviceDefaultSecretsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"SiteId": {
			Title: "Site ID",
			Order: 2,
		},
		"MacAddressOrSerialNumber": {
			Title:    "MAC/Serial",
			MaxWidth: 30,
			Order:    3,
		},
		"SecretName": {
			Title: "Secret Name",
			Order: 4,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       5,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
	},
}
