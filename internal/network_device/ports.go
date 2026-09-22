package network_device

import (
	"context"
	"net/http"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var networkDevicePortPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"InterfaceId": {
			Title: "ID",
			Order: 1,
		},
		"InterfaceName": {
			Title: "Name",
			Order: 2,
		},
		"Kind": {
			Order: 3,
		},
		"InterfaceDescription": {
			Title:    "Description",
			MaxWidth: 30,
			Order:    4,
		},
		"Mtu": {
			Title: "MTU",
			Order: 5,
		},
		"Enabled": {
			Order: 6,
		},
		"ServiceStatus": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       7,
		},
		"PendingDelete": {
			Title: "Pending Delete",
			Order: 8,
		},
	},
}

var networkDeviceLivePortPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"PortName": {
			Title: "Name",
			Order: 1,
		},
		"Enabled": {
			Order: 2,
		},
		"Active": {
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
			Title: "Util In",
			Order: 6,
		},
		"UtilizationOut": {
			Title: "Util Out",
			Order: 7,
		},
	},
}

// NetworkDevicePortCreate creates a new logical interface (e.g. a loopback or a
// sub-interface) on a network device.
func NetworkDevicePortCreate(ctx context.Context, networkDeviceRef string, config []byte) error {
	logger.Get().Info().Msgf("Creating port on network device '%s'", networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	var portConfig sdk.CreateNetworkEquipmentInterface
	if err := utils.UnmarshalContent(config, &portConfig); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	portInfo, httpRes, err := client.NetworkDeviceAPI.
		CreateNetworkDevicePort(ctx, networkDeviceIdNumeric).
		CreateNetworkEquipmentInterface(portConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(portInfo, &networkDevicePortPrintConfig)
}

func NetworkDevicePortConfigExample(ctx context.Context) error {
	portConfig := sdk.CreateNetworkEquipmentInterface{
		Kind:        "loopback",
		Name:        "Loopback1",
		Description: sdk.PtrString("management loopback"),
		Mtu:         sdk.PtrInt32(9216),
		Enabled:     sdk.PtrBool(true),
	}

	return formatter.PrintResult(portConfig, nil)
}

// NetworkDevicePortDelete removes a logical interface from a network device.
func NetworkDevicePortDelete(ctx context.Context, networkDeviceRef string, portId string) error {
	logger.Get().Info().Msgf("Deleting port %s of network device '%s'", portId, networkDeviceRef)

	networkDeviceIdNumeric, portIdNumeric, err := resolvePortIds(ctx, networkDeviceRef, portId)
	if err != nil {
		return err
	}

	revision, err := getNetworkDevicePortRevision(ctx, networkDeviceIdNumeric, portIdNumeric)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, err = executeWithRevisionRetry(revision, func(revision string) (any, *http.Response, error) {
		httpRes, err := client.NetworkDeviceAPI.
			DeleteNetworkDevicePort(ctx, networkDeviceIdNumeric, portIdNumeric).
			IfMatch(revision).
			Execute()
		return nil, httpRes, err
	})
	if err != nil {
		return err
	}

	logger.Get().Info().Msgf("Port %s of network device '%s' deleted", portId, networkDeviceRef)
	return nil
}

// NetworkDevicePortGetConfig shows the staged (desired) configuration of a port.
func NetworkDevicePortGetConfig(ctx context.Context, networkDeviceRef string, portId string) error {
	logger.Get().Info().Msgf("Getting config of port %s of network device '%s'", portId, networkDeviceRef)

	networkDeviceIdNumeric, portIdNumeric, err := resolvePortIds(ctx, networkDeviceRef, portId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	configInfo, httpRes, err := client.NetworkDeviceAPI.
		GetNetworkDevicePortConfig(ctx, networkDeviceIdNumeric, portIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(configInfo, &networkDevicePortConfigPrintConfig)
}

// NetworkDevicePortsLiveStatus queries the device itself for the operational
// state of its physical ports (link state, speed, utilization). This is the
// POST .../actions/ports action and differs from 'get-ports', which lists the
// interface inventory stored by MetalSoft.
func NetworkDevicePortsLiveStatus(ctx context.Context, networkDeviceRef string) error {
	logger.Get().Info().Msgf("Getting live port status of network device '%s'", networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	portsInfo, httpRes, err := client.NetworkDeviceAPI.
		GetPorts(ctx, float32(networkDeviceIdNumeric)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if portsInfo == nil {
		return errEmptyBody("port status")
	}

	return formatter.PrintResult(portsInfo.Ports, &networkDeviceLivePortPrintConfig)
}

// NetworkDevicePortIpList lists the IP addresses of one address family staged
// on a network device port.
func NetworkDevicePortIpList(ctx context.Context, networkDeviceRef string, portId string, family string) error {
	logger.Get().Info().Msgf("Listing %s addresses of port %s of network device '%s'", family, portId, networkDeviceRef)

	networkDeviceIdNumeric, portIdNumeric, err := resolvePortIds(ctx, networkDeviceRef, portId)
	if err != nil {
		return err
	}

	if err := validateAddressFamily(family); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	ips, httpRes, err := client.NetworkDeviceAPI.
		ListNetworkDevicePortIps(ctx, networkDeviceIdNumeric, portIdNumeric, family).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(ips, &networkDevicePortIpPrintConfig)
}

// NetworkDevicePortIpGet shows a single IP address of a network device port.
func NetworkDevicePortIpGet(ctx context.Context, networkDeviceRef string, portId string, family string, ipId string) error {
	logger.Get().Info().Msgf("Getting %s address %s of port %s of network device '%s'", family, ipId, portId, networkDeviceRef)

	networkDeviceIdNumeric, portIdNumeric, err := resolvePortIds(ctx, networkDeviceRef, portId)
	if err != nil {
		return err
	}

	if err := validateAddressFamily(family); err != nil {
		return err
	}

	ipIdNumeric, err := utils.GetInt64FromString(ipId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	ipInfo, httpRes, err := client.NetworkDeviceAPI.
		GetNetworkDevicePortIp(ctx, networkDeviceIdNumeric, portIdNumeric, family, ipIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(ipInfo, &networkDevicePortIpPrintConfig)
}

// NetworkDevicePortIpRemove removes one IP address from a network device port.
func NetworkDevicePortIpRemove(ctx context.Context, networkDeviceRef string, portId string, family string, ipId string) error {
	logger.Get().Info().Msgf("Removing %s address %s from port %s of network device '%s'", family, ipId, portId, networkDeviceRef)

	networkDeviceIdNumeric, portIdNumeric, err := resolvePortIds(ctx, networkDeviceRef, portId)
	if err != nil {
		return err
	}

	if err := validateAddressFamily(family); err != nil {
		return err
	}

	ipIdNumeric, err := utils.GetInt64FromString(ipId)
	if err != nil {
		return err
	}

	revision, err := getNetworkDevicePortRevision(ctx, networkDeviceIdNumeric, portIdNumeric)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	_, err = executeWithRevisionRetry(revision, func(revision string) (any, *http.Response, error) {
		httpRes, err := client.NetworkDeviceAPI.
			RemoveNetworkDevicePortIp(ctx, networkDeviceIdNumeric, portIdNumeric, family, ipIdNumeric).
			IfMatch(revision).
			Execute()
		return nil, httpRes, err
	})
	if err != nil {
		return err
	}

	logger.Get().Info().Msgf("Address %s removed from port %s of network device '%s'", ipId, portId, networkDeviceRef)
	return nil
}

// NetworkDevicePortIpReplace replaces the complete IP address set of one
// address family on a network device port with the supplied list.
func NetworkDevicePortIpReplace(ctx context.Context, networkDeviceRef string, portId string, family string, config []byte) error {
	logger.Get().Info().Msgf("Replacing %s addresses of port %s of network device '%s'", family, portId, networkDeviceRef)

	networkDeviceIdNumeric, portIdNumeric, err := resolvePortIds(ctx, networkDeviceRef, portId)
	if err != nil {
		return err
	}

	if err := validateAddressFamily(family); err != nil {
		return err
	}

	var ipsConfig sdk.ReplaceNetworkEquipmentInterfaceIps
	if err := utils.UnmarshalContent(config, &ipsConfig); err != nil {
		return err
	}

	revision, err := getNetworkDevicePortRevision(ctx, networkDeviceIdNumeric, portIdNumeric)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	ips, err := executeWithRevisionRetry(revision, func(revision string) ([]sdk.NetworkEquipmentInterfaceIp, *http.Response, error) {
		return client.NetworkDeviceAPI.
			ReplaceNetworkDevicePortIps(ctx, networkDeviceIdNumeric, portIdNumeric, family).
			ReplaceNetworkEquipmentInterfaceIps(ipsConfig).
			IfMatch(revision).
			Execute()
	})
	if err != nil {
		return err
	}

	return formatter.PrintResult(ips, &networkDevicePortIpPrintConfig)
}

func NetworkDevicePortIpReplaceConfigExample(ctx context.Context) error {
	ipsConfig := sdk.ReplaceNetworkEquipmentInterfaceIps{
		Ips: []sdk.AddNetworkEquipmentInterfaceIp{
			{Address: "10.0.0.1", PrefixLength: 32},
			{Address: "10.0.0.2", PrefixLength: 32},
		},
	}

	return formatter.PrintResult(ipsConfig, nil)
}

// resolvePortIds resolves both the device reference and the numeric port
// (interface) id used by the port sub-resources.
func resolvePortIds(ctx context.Context, networkDeviceRef string, portId string) (int64, int64, error) {
	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return 0, 0, err
	}

	portIdNumeric, err := getNetworkDevicePortId(portId)
	if err != nil {
		return 0, 0, err
	}

	return networkDeviceIdNumeric, portIdNumeric, nil
}

// getNetworkDevicePortRevision derives the optimistic-lock revision the port
// sub-resources expect. The interface entity revision is not served directly:
// the single-port GET exposes a zero-based config.revision while the lock
// checks a one-based counter, so config.revision + 1 is sent. See
// executeWithRevisionRetry for the correction path when that guess is stale.
func getNetworkDevicePortRevision(ctx context.Context, networkDeviceId int64, portId int64) (string, error) {
	client := api.GetApiClient(ctx)

	port, httpRes, err := client.NetworkDeviceAPI.
		GetNetworkDevicePort(ctx, networkDeviceId, portId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return "", err
	}

	if port == nil {
		return "", errEmptyBody("port")
	}

	return strconv.FormatInt(port.Config.Revision+1, 10), nil
}

// executeWithRevisionRetry runs an optimistic-locking request and, when the
// server answers 409 naming the revision it actually expects, retries once
// with that value.
func executeWithRevisionRetry[T any](revision string, execute func(revision string) (T, *http.Response, error)) (T, error) {
	result, httpRes, err := execute(revision)

	if err != nil && httpRes != nil && httpRes.StatusCode == http.StatusConflict {
		if expected := expectedRevisionFromError(err); expected != "" && expected != revision {
			result, httpRes, err = execute(expected)
		}
	}

	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		var empty T
		return empty, err
	}

	return result, nil
}
