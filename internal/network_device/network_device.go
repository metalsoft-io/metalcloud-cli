package network_device

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

var NetworkDevicePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
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

	networkDeviceList, httpRes, err := request.SortBy([]string{"id:ASC"}).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkDeviceList, &NetworkDevicePrintConfig)
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
		SiteId:           sdk.PtrInt32(1),
		Driver:           sdk.NETWORKDEVICEDRIVER_SONIC_ENTERPRISE,
		IdentifierString: sdk.PtrString("example"),
		SerialNumber:     sdk.PtrString("1234567890"),
		ChassisRackId:    sdk.PtrInt32(1),
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

	result, httpRes, err := client.NetworkDeviceAPI.
		DeleteNetworkDevice(ctx, networkDeviceIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network device %s deleting in progress. Job ID %d", networkDeviceId, result.JobId)
	return nil
}

func NetworkDeviceArchive(ctx context.Context, networkDeviceId string) error {
	logger.Get().Info().Msgf("Archiving network device %s", networkDeviceId)

	networkDeviceIdNumeric, revision, err := getNetworkDeviceIdAndRevision(ctx, networkDeviceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	result, httpRes, err := client.NetworkDeviceAPI.
		ArchiveNetworkDevice(ctx, networkDeviceIdNumeric).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network device %s archiving in progress. Job ID %d", networkDeviceId, result.JobId)
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

	portsInfo, err := GetNetworkDevicePorts(ctx, networkDeviceIdNumeric)
	if err != nil {
		return err
	}

	return formatter.PrintResult(portsInfo, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"PortName": {
				Title: "Name",
				Order: 1,
			},
			"Enabled": {
				Title: "Enabled",
				Order: 2,
			},
			"Active": {
				Title: "Active",
				Order: 3,
			},
			"LinkSpeed": {
				Title: "Speed",
				Order: 4,
			},
			"LinkDuplex": {
				Title: "Duplex",
				Order: 5,
			},
			"UtilizationIn": {
				Title: "Utilization In",
				Order: 6,
			},
			"UtilizationOut": {
				Title: "Utilization Out",
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

	httpRes, err := client.NetworkDeviceAPI.
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
				Title: "#",
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

func NetworkDeviceExampleDefaults(ctx context.Context) error {
	networkDeviceDefaults := sdk.CreateNetworkDeviceDefaults{
		DatacenterName:            "site1",
		SerialNumber:              sdk.PtrString("1234"),
		ManagementMacAddress:      "AA:BB:CC:DD:EE:FF",
		Position:                  sdk.PtrString("leaf"),
		IdentifierString:          sdk.PtrString("1234"),
		Asn:                       sdk.PtrInt64(65000),
		CustomVariables:           map[string]interface{}{"var1": "value1", "var2": "value2"},
		MlagDomainId:              sdk.PtrInt32(1),
		LoopbackAddressIpv4:       sdk.PtrString("1.2.3.4"),
		LoopbackAddressIpv6:       sdk.PtrString("0:0:0:0:0:0:0:1"),
		VtepAddressIpv4:           sdk.PtrString("1.2.3.4"),
		VtepAddressIpv6:           sdk.PtrString("0:0:0:0:0:0:0:1"),
		OrderIndex:                sdk.PtrInt32(1),
		OsTemplateId:              sdk.PtrInt32(10),
		MlagPartnerHostname:       sdk.PtrString("partner.example.com"),
		IsPartOfMlagPair:          sdk.PtrBool(true),
		MlagSystemMac:             sdk.PtrString("AA:BB:CC:DD:EE:FF"),
		MlagPeerLinkPortChannelId: sdk.PtrInt32(1),
		MlagPartnerVlanId:         sdk.PtrInt32(100),
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

func GetNetworkDeviceId(networkDeviceId string) (float32, error) {
	networkDeviceIdNumeric, err := strconv.ParseFloat(networkDeviceId, 32)
	if err != nil {
		err := fmt.Errorf("invalid network device ID: '%s'", networkDeviceId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(networkDeviceIdNumeric), nil
}

func getNetworkDeviceIdAndRevision(ctx context.Context, networkDeviceId string) (float32, string, error) {
	networkDeviceIdNumeric, err := GetNetworkDeviceId(networkDeviceId)
	if err != nil {
		return 0, "", err
	}

	client := api.GetApiClient(ctx)

	networkDevice, httpRes, err := client.NetworkDeviceAPI.GetNetworkDevice(ctx, float32(networkDeviceIdNumeric)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}

	return float32(networkDeviceIdNumeric), strconv.Itoa(int(networkDevice.Revision)), nil
}
