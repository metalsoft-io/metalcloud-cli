package network_device

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
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

func NetworkDeviceDiscover(ctx context.Context, networkDeviceId string) error {
	logger.Get().Info().Msgf("Discovering network device %s", networkDeviceId)

	networkDeviceIdNumeric, err := GetNetworkDeviceId(networkDeviceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceAPI.
		DiscoverNetworkDevice(ctx, networkDeviceIdNumeric).
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

func NetworkDeviceGetDefaults(ctx context.Context, siteId string) error {
	logger.Get().Info().Msgf("Getting network device defaults for site %s", siteId)

	siteIdNumeric, err := utils.GetFloat32FromString(siteId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	defaults, httpRes, err := client.NetworkDeviceAPI.
		GetNetworkDeviceDefaults(ctx, siteIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(defaults, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"Id": {
				Title: "ID",
				Order: 1,
			},
			"DatacenterName": {
				Title: "Site",
				Order: 2,
			},
			"SerialNumber": {
				Title: "Serial",
				Order: 3,
			},
			"ManagementMacAddress": {
				Title: "MAC",
				Order: 4,
			},
			"Position": {
				Title: "Position",
				Order: 5,
			},
			"IdentifierString": {
				Title: "Identifier",
				Order: 6,
			},
			"Asn": {
				Title: "ASN",
				Order: 7,
			},
		},
	})
}

func NetworkDeviceAddDefaults(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Adding network device defaults")

	var networkDeviceDefaults sdk.CreateNetworkDeviceDefaults
	err := utils.UnmarshalContent(config, &networkDeviceDefaults)
	if err != nil {
		return err
	}
	if networkDeviceDefaults.ManagementMacAddress == "" {
		return fmt.Errorf("invalid content - please make sure the correct format is specified")
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceAPI.
		AddNetworkDeviceDefaults(ctx).
		CreateNetworkDeviceDefaults(networkDeviceDefaults).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return nil
}

func NetworkDeviceDeleteDefaults(ctx context.Context, siteId string, defaultsId string) error {
	logger.Get().Info().Msgf("Deleting network device defaults %s for site %s", defaultsId, siteId)

	siteIdNumeric, err := utils.GetFloat32FromString(siteId)
	if err != nil {
		return err
	}

	defaultsIdNumeric, err := utils.GetInt64FromString(defaultsId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceAPI.
		RemoveNetworkDeviceDefaults(ctx, siteIdNumeric, defaultsIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network device defaults %s for site %s deleted successfully", defaultsId, siteId)
	return nil
}

func NetworkDeviceExampleDefaults(ctx context.Context) error {
	networkDeviceDefaults := sdk.CreateNetworkDeviceDefaults{
		DatacenterName:            "site1",
		SerialNumber:              sdk.PtrString("1234"),
		ManagementMacAddress:      "AA:BB:CC:DD:EE:FF",
		Position:                  sdk.PtrString("leaf"),
		IdentifierString:          sdk.PtrString("1234"),
		Asn:                       sdk.PtrInt64(65000),
		CustomVariables:           map[string]interface{}{"var1": "value1", "var2": "value2"},
		MlagDomainId:              sdk.PtrInt64(1),
		LoopbackAddressIpv4:       sdk.PtrString("1.2.3.4"),
		LoopbackAddressIpv6:       sdk.PtrString("0:0:0:0:0:0:0:1"),
		VtepAddressIpv4:           sdk.PtrString("1.2.3.4"),
		VtepAddressIpv6:           sdk.PtrString("0:0:0:0:0:0:0:1"),
		OrderIndex:                sdk.PtrInt32(1),
		OsTemplateId:              sdk.PtrInt64(10),
		MlagPartnerHostname:       sdk.PtrString("partner.example.com"),
		IsPartOfMlagPair:          sdk.PtrBool(true),
		MlagSystemMac:             sdk.PtrString("AA:BB:CC:DD:EE:FF"),
		MlagPeerLinkPortChannelId: sdk.PtrInt64(1),
		MlagPartnerVlanId:         sdk.PtrInt32(100),
		AuthenticationOptions: []sdk.NetworkDeviceAuthOption{
			{Kind: "tacacs", DeviceAuthProviderId: sdk.PtrInt64(1)},
			{Kind: "local"},
		},
	}

	return formatter.PrintResult(networkDeviceDefaults, nil)
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

func NetworkDeviceDefaultSecretsList(ctx context.Context, page int, limit int) error {
	logger.Get().Info().Msgf("Listing network device default secrets")

	client := api.GetApiClient(ctx)

	req := client.NetworkDeviceDefaultSecretsAPI.GetNetworkDevicesDefaultSecrets(ctx)

	req = req.SortBy([]string{"id:ASC"})

	switch {
	case page > 0:
		records, meta, err := utils.FetchPageWindow(req, page, limit)
		if err != nil {
			return err
		}
		return utils.PrintAll(records, meta, len(records), &networkDeviceDefaultSecretsPrintConfig)
	case limit > 0:
		records, meta, err := utils.FetchUpTo(req, limit)
		if err != nil {
			return err
		}
		return utils.PrintAll(records, meta, len(records), &networkDeviceDefaultSecretsPrintConfig)
	default:
		records, meta, err := utils.FetchAllPages(req)
		if err != nil {
			return err
		}
		return utils.PrintAll(records, meta, len(records), &networkDeviceDefaultSecretsPrintConfig)
	}
}

func NetworkDeviceDefaultSecretsGet(ctx context.Context, secretsId string) error {
	logger.Get().Info().Msgf("Get network device default secrets '%s'", secretsId)

	secretsIdNumeric, err := parseDefaultSecretsId(secretsId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	secrets, httpRes, err := client.NetworkDeviceDefaultSecretsAPI.GetNetworkDeviceDefaultSecretsInfo(ctx, secretsIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(secrets, &networkDeviceDefaultSecretsPrintConfig)
}

func NetworkDeviceDefaultSecretsGetCredentials(ctx context.Context, secretsId string) error {
	logger.Get().Info().Msgf("Get network device default secrets credentials for '%s'", secretsId)

	secretsIdNumeric, err := parseDefaultSecretsId(secretsId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	credentials, httpRes, err := client.NetworkDeviceDefaultSecretsAPI.GetNetworkDeviceDefaultSecretsCredentials(ctx, secretsIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(credentials, nil)
}

func NetworkDeviceDefaultSecretsCreate(ctx context.Context, siteId float32, macAddressOrSerialNumber string, secretName string, secretValue string) error {
	logger.Get().Info().Msgf("Creating network device default secrets")

	client := api.GetApiClient(ctx)

	createSecrets := sdk.NewCreateNetworkDeviceDefaultSecrets(int64(siteId), macAddressOrSerialNumber, secretName, secretValue)

	secrets, httpRes, err := client.NetworkDeviceDefaultSecretsAPI.CreateNetworkDeviceDefaultSecrets(ctx).
		CreateNetworkDeviceDefaultSecrets(*createSecrets).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(secrets, &networkDeviceDefaultSecretsPrintConfig)
}

func NetworkDeviceDefaultSecretsUpdate(ctx context.Context, secretsId string, secretValue string) error {
	logger.Get().Info().Msgf("Updating network device default secrets '%s'", secretsId)

	secretsIdNumeric, err := parseDefaultSecretsId(secretsId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	updateSecrets := sdk.NewUpdateNetworkDeviceDefaultSecrets()
	updateSecrets.SetSecretValue(secretValue)

	secrets, httpRes, err := client.NetworkDeviceDefaultSecretsAPI.
		UpdateNetworkDeviceDefaultSecrets(ctx, secretsIdNumeric).
		UpdateNetworkDeviceDefaultSecrets(*updateSecrets).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(secrets, &networkDeviceDefaultSecretsPrintConfig)
}

func NetworkDeviceDefaultSecretsDelete(ctx context.Context, secretsId string) error {
	logger.Get().Info().Msgf("Deleting network device default secrets '%s'", secretsId)

	secretsIdNumeric, err := parseDefaultSecretsId(secretsId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceDefaultSecretsAPI.DeleteNetworkDeviceDefaultSecrets(ctx, secretsIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network device default secrets with ID %s deleted successfully", secretsId)
	return nil
}

func NetworkDeviceDefaultSecretsBatchCreate(ctx context.Context, filePath string) error {
	logger.Get().Info().Msgf("Batch creating network device default secrets from %s", filePath)

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open CSV file: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("unable to read CSV header: %w", err)
	}

	colIndex := map[string]int{}
	for i, col := range header {
		colIndex[strings.TrimSpace(strings.ToLower(col))] = i
	}

	requiredCols := []string{"siteid", "macaddressorserialnumber", "secretname", "secretvalue"}
	for _, col := range requiredCols {
		if _, ok := colIndex[col]; !ok {
			return fmt.Errorf("missing required CSV column: %s (expected columns: siteId,macAddressOrSerialNumber,secretName,secretValue)", col)
		}
	}

	client := api.GetApiClient(ctx)

	var created []sdk.NetworkDeviceDefaultSecrets
	rowNum := 1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV row %d: %w", rowNum+1, err)
		}
		rowNum++

		siteIdStr := strings.TrimSpace(record[colIndex["siteid"]])
		siteIdFloat, err := strconv.ParseFloat(siteIdStr, 32)
		if err != nil {
			return fmt.Errorf("invalid siteId on row %d: '%s'", rowNum, siteIdStr)
		}

		macOrSerial := strings.TrimSpace(record[colIndex["macaddressorserialnumber"]])
		secretName := strings.TrimSpace(record[colIndex["secretname"]])
		secretValue := strings.TrimSpace(record[colIndex["secretvalue"]])

		createSecrets := sdk.NewCreateNetworkDeviceDefaultSecrets(int64(siteIdFloat), macOrSerial, secretName, secretValue)

		secrets, httpRes, err := client.NetworkDeviceDefaultSecretsAPI.CreateNetworkDeviceDefaultSecrets(ctx).
			CreateNetworkDeviceDefaultSecrets(*createSecrets).
			Execute()
		if err := response_inspector.InspectResponse(httpRes, err); err != nil {
			return fmt.Errorf("failed to create secret on row %d: %w", rowNum, err)
		}

		created = append(created, *secrets)
		logger.Get().Info().Msgf("Created secret #%d for %s", int(secrets.Id), macOrSerial)
	}

	if len(created) == 0 {
		return fmt.Errorf("no data rows found in CSV file")
	}

	return formatter.PrintResult(created, &networkDeviceDefaultSecretsPrintConfig)
}

func parseDefaultSecretsId(secretsId string) (int64, error) {
	secretsIdNumeric, err := strconv.ParseInt(secretsId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid network device default secrets ID: '%s'", secretsId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return secretsIdNumeric, nil
}
