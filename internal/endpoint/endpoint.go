package endpoint

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/internal/network_device"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

type EndpointInterfaceDetails struct {
	Id                         *int64
	MacAddress                 *string
	NetworkDeviceId            *int64
	NetworkDeviceInterfaceId   *int64
	NetworkDeviceName          *string
	NetworkDeviceInterfaceName *string
}

var endpointPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"SiteId": {
			Title: "Site",
			Order: 2,
		},
		"Name": {
			Title: "Name",
			Order: 3,
		},
		"Label": {
			Title: "Label",
			Order: 4,
		},
		"ExternalId": {
			Title: "External Id",
			Order: 5,
		},
	},
}

var endpointInterfacePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"MacAddress": {
			Title: "MAC Address",
			Order: 2,
		},
		"NetworkDeviceName": {
			Title: "Switch",
			Order: 3,
		},
		"NetworkDeviceInterfaceName": {
			Title: "Switch Port",
			Order: 4,
		},
	},
}

func EndpointList(ctx context.Context, filterSite []string, filterExternalId []string) error {
	logger.Get().Info().Msgf("Listing all endpoints")

	client := api.GetApiClient(ctx)

	request := client.EndpointAPI.GetEndpoints(ctx).SortBy([]string{"id:ASC"})

	if len(filterSite) > 0 {
		request = request.FilterSiteId(utils.ProcessFilterStringSlice(filterSite))
	}

	if len(filterExternalId) > 0 {
		request = request.FilterExternalId(utils.ProcessFilterStringSlice(filterExternalId))
	}

	endpoints, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(endpoints, meta, len(endpoints), &endpointPrintConfig)
}

func EndpointGet(ctx context.Context, endpointId string) error {
	logger.Get().Info().Msgf("Get endpoint '%s'", endpointId)

	endpointIdNumeric, err := GetEndpointId(endpointId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	endpointInfo, httpRes, err := client.EndpointAPI.GetEndpointById(ctx, endpointIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(endpointInfo, &endpointPrintConfig)
}

func EndpointCreate(ctx context.Context, endpointConfig sdk.CreateEndpoint) error {
	logger.Get().Info().Msgf("Creating new endpoint")

	client := api.GetApiClient(ctx)

	endpointInfo, httpRes, err := client.EndpointAPI.CreateEndpoint(ctx).CreateEndpoint(endpointConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if endpointConfig.EndpointInterfaces != nil && len(endpointConfig.EndpointInterfaces) > 0 {

	}

	return formatter.PrintResult(endpointInfo, &endpointPrintConfig)
}

// EndpointInterfaceInput describes one endpoint interface in a create/bulk-create
// config. The switch port may be given directly by its numeric
// networkDeviceInterfaceId, or looked up by network device (id or label) plus
// interface name (e.g. "swp9s0").
type EndpointInterfaceInput struct {
	NetworkDeviceInterfaceId *int64  `json:"networkDeviceInterfaceId,omitempty" yaml:"networkDeviceInterfaceId,omitempty"`
	NetworkDevice            *string `json:"networkDevice,omitempty" yaml:"networkDevice,omitempty"`
	Interface                *string `json:"interface,omitempty" yaml:"interface,omitempty"`
	MacAddress               *string `json:"macAddress,omitempty" yaml:"macAddress,omitempty"`
}

// EndpointInput is a lenient view of an endpoint create body that additionally
// allows endpoint interfaces to reference their switch port by network
// device/interface label instead of a numeric id.
type EndpointInput struct {
	ExternalId         *string                  `json:"externalId,omitempty" yaml:"externalId,omitempty"`
	SiteId             int64                    `json:"siteId" yaml:"siteId"`
	Name               string                   `json:"name" yaml:"name"`
	Label              string                   `json:"label" yaml:"label"`
	EndpointInterfaces []EndpointInterfaceInput `json:"endpointInterfaces,omitempty" yaml:"endpointInterfaces,omitempty"`
}

// EndpointCreateBulk creates multiple endpoints in a single call from a config
// holding a list of endpoint definitions. Each endpoint interface may be given
// by numeric networkDeviceInterfaceId, or by network device + interface label
// (resolved here to the numeric id before the bulk call).
func EndpointCreateBulk(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating endpoints in bulk")

	var inputs []EndpointInput
	if err := utils.UnmarshalContent(config, &inputs); err != nil {
		return err
	}
	if len(inputs) == 0 {
		return fmt.Errorf("no endpoints found in configuration")
	}

	resolver := newInterfaceResolver()
	endpoints := make([]sdk.CreateEndpoint, 0, len(inputs))
	for _, input := range inputs {
		endpoint, err := input.toCreateEndpoint(ctx, resolver)
		if err != nil {
			return err
		}
		endpoints = append(endpoints, endpoint)
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.EndpointAPI.
		BulkCreateEndpoints(ctx).
		BulkCreateEndpoints(sdk.BulkCreateEndpoints{Endpoints: endpoints}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Created %d endpoint(s) in bulk", len(endpoints))
	return nil
}

// toCreateEndpoint converts a lenient EndpointInput into the SDK create body,
// resolving each interface's network device/interface label to a numeric id.
func (input EndpointInput) toCreateEndpoint(ctx context.Context, resolver *interfaceResolver) (sdk.CreateEndpoint, error) {
	endpoint := sdk.CreateEndpoint{
		ExternalId: input.ExternalId,
		SiteId:     input.SiteId,
		Name:       input.Name,
		Label:      input.Label,
	}

	for i, ifaceInput := range input.EndpointInterfaces {
		interfaceId, err := resolver.resolve(ctx, input.SiteId, ifaceInput)
		if err != nil {
			return sdk.CreateEndpoint{}, fmt.Errorf("endpoint '%s' interface #%d: %w", input.Label, i+1, err)
		}

		endpoint.EndpointInterfaces = append(endpoint.EndpointInterfaces, sdk.CreateEndpointInterface{
			NetworkDeviceInterfaceId: interfaceId,
			MacAddress:               ifaceInput.MacAddress,
		})
	}

	return endpoint, nil
}

// interfaceResolver resolves (network device label, interface name) pairs to
// numeric interface ids, caching each device's port inventory so a device is
// fetched at most once across a bulk request. Devices are cached per site
// because switch labels are reused across sites.
type interfaceResolver struct {
	// "siteId/deviceRef" -> (interfaceName -> interfaceId)
	ports map[string]map[string]int64
}

func newInterfaceResolver() *interfaceResolver {
	return &interfaceResolver{ports: map[string]map[string]int64{}}
}

func (r *interfaceResolver) resolve(ctx context.Context, siteId int64, input EndpointInterfaceInput) (int64, error) {
	// A numeric id, if given, wins and needs no lookup.
	if input.NetworkDeviceInterfaceId != nil {
		return *input.NetworkDeviceInterfaceId, nil
	}

	if input.NetworkDevice == nil || *input.NetworkDevice == "" || input.Interface == nil || *input.Interface == "" {
		return 0, fmt.Errorf("interface must specify either networkDeviceInterfaceId, or both networkDevice and interface")
	}

	deviceRef := *input.NetworkDevice
	interfaceName := *input.Interface

	cacheKey := fmt.Sprintf("%d/%s", siteId, deviceRef)
	ports, ok := r.ports[cacheKey]
	if !ok {
		logger.Get().Debug().Msgf("Resolving ports for network device '%s' in site %d", deviceRef, siteId)

		device, err := network_device.GetNetworkDeviceByIdOrLabelInSite(ctx, deviceRef, siteId)
		if err != nil {
			return 0, err
		}

		deviceIdNumeric, err := strconv.ParseFloat(device.Id, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid network device id '%s': %w", device.Id, err)
		}

		interfaces, err := network_device.GetNetworkDevicePorts(ctx, float32(deviceIdNumeric))
		if err != nil {
			return 0, err
		}

		ports = make(map[string]int64, len(interfaces))
		for _, iface := range interfaces {
			ports[iface.InterfaceName] = iface.InterfaceId
		}
		r.ports[cacheKey] = ports

		logger.Get().Debug().Msgf("Network device '%s' (id %s) in site %d has %d ports", deviceRef, device.Id, siteId, len(ports))
	}

	interfaceId, ok := ports[interfaceName]
	if !ok {
		return 0, fmt.Errorf("interface '%s' not found on network device '%s'", interfaceName, deviceRef)
	}

	logger.Get().Debug().Msgf("Resolved network device '%s' interface '%s' to id %d", deviceRef, interfaceName, interfaceId)
	return interfaceId, nil
}

func EndpointUpdate(ctx context.Context, endpointId string, endpointUpdates sdk.UpdateEndpoint) error {
	logger.Get().Info().Msgf("Updating endpoint '%s'", endpointId)

	endpointIdNumeric, err := GetEndpointId(endpointId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	endpointInfo, httpRes, err := client.EndpointAPI.GetEndpointById(ctx, endpointIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	endpointInfo, httpRes, err = client.EndpointAPI.
		UpdateEndpoint(ctx, endpointIdNumeric).
		UpdateEndpoint(endpointUpdates).
		IfMatch(endpointInfo.Revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(endpointInfo, &endpointPrintConfig)
}

func EndpointDelete(ctx context.Context, endpointId string) error {
	logger.Get().Info().Msgf("Deleting endpoint '%s'", endpointId)

	endpointIdNumeric, err := GetEndpointId(endpointId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	endpointInfo, httpRes, err := client.EndpointAPI.GetEndpointById(ctx, endpointIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	httpRes, err = client.EndpointAPI.
		DeleteEndpoint(ctx, endpointIdNumeric).
		IfMatch(endpointInfo.Revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Endpoint '%s' deleted successfully", endpointId)

	return nil
}

func EndpointInterfaceList(ctx context.Context, endpointId string) error {
	logger.Get().Info().Msgf("Listing interfaces for endpoint '%s'", endpointId)

	endpointIdNumeric, err := GetEndpointId(endpointId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.EndpointAPI.GetEndpointInterfaces(ctx, endpointIdNumeric).SortBy([]string{"id:ASC"})

	interfaces, _, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	endpointInterfacesList := make([]EndpointInterfaceDetails, 0, len(interfaces))
	for _, iface := range interfaces {
		endpointInterface := EndpointInterfaceDetails{
			Id:                         &iface.Id,
			MacAddress:                 iface.MacAddress,
			NetworkDeviceId:            &iface.NetworkDeviceId,
			NetworkDeviceInterfaceId:   &iface.NetworkDeviceInterfaceId,
			NetworkDeviceInterfaceName: &iface.NetworkDeviceInterfaceName,
		}

		endpointInterfacesList = append(endpointInterfacesList, endpointInterface)
	}

	return formatter.PrintResult(endpointInterfacesList, &endpointInterfacePrintConfig)
}

func GetEndpointId(endpointId string) (int64, error) {
	endpointIdNumeric, err := strconv.ParseInt(endpointId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid endpoint ID: '%s'", endpointId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return endpointIdNumeric, nil
}
