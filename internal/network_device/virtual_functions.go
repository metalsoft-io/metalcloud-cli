package network_device

import (
	"context"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
)

var networkDeviceVirtualFunctionPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"NetworkDeviceId": {
			Title: "Device",
			Order: 2,
		},
		"InterfaceId": {
			Title: "Interface",
			Order: 3,
		},
		"Name": {
			Order:    4,
			MaxWidth: 40,
		},
		"Index": {
			Order: 5,
		},
		"LogicalNetworkId": {
			Title: "Logical Network",
			Order: 6,
		},
		"Status": {
			Transformer: formatter.FormatStatusValue,
			Order:       7,
		},
	},
}

// NetworkDeviceVirtualFunctionList lists all virtual functions of a network
// device, across all of its interfaces.
func NetworkDeviceVirtualFunctionList(ctx context.Context, networkDeviceRef string) error {
	logger.Get().Info().Msgf("Listing virtual functions of network device '%s'", networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.NetworkDeviceAPI.
		GetNetworkDeviceVirtualFunctions(ctx, float32(networkDeviceIdNumeric)).
		SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &networkDeviceVirtualFunctionPrintConfig)
}

func NetworkDeviceVirtualFunctionGet(ctx context.Context, networkDeviceRef string, virtualFunctionId string) error {
	logger.Get().Info().Msgf("Getting virtual function %s of network device '%s'", virtualFunctionId, networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	virtualFunctionIdNumeric, err := utils.GetInt64FromString(virtualFunctionId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	virtualFunction, httpRes, err := client.NetworkDeviceAPI.
		GetNetworkDeviceVirtualFunction(ctx, float32(networkDeviceIdNumeric), float32(virtualFunctionIdNumeric)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(virtualFunction, &networkDeviceVirtualFunctionPrintConfig)
}

// NetworkDevicePortVirtualFunctionList lists the virtual functions of a single
// network device port.
func NetworkDevicePortVirtualFunctionList(ctx context.Context, networkDeviceRef string, portId string) error {
	logger.Get().Info().Msgf("Listing virtual functions of port %s of network device '%s'", portId, networkDeviceRef)

	networkDeviceIdNumeric, portIdNumeric, err := resolvePortIds(ctx, networkDeviceRef, portId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.NetworkDeviceAPI.
		GetNetworkDevicePortVirtualFunctions(ctx, float32(networkDeviceIdNumeric), float32(portIdNumeric)).
		SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &networkDeviceVirtualFunctionPrintConfig)
}

func NetworkDevicePortVirtualFunctionGet(ctx context.Context, networkDeviceRef string, portId string, virtualFunctionId string) error {
	logger.Get().Info().Msgf("Getting virtual function %s of port %s of network device '%s'", virtualFunctionId, portId, networkDeviceRef)

	networkDeviceIdNumeric, portIdNumeric, err := resolvePortIds(ctx, networkDeviceRef, portId)
	if err != nil {
		return err
	}

	virtualFunctionIdNumeric, err := utils.GetInt64FromString(virtualFunctionId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	virtualFunction, httpRes, err := client.NetworkDeviceAPI.
		GetNetworkDevicePortVirtualFunction(ctx, float32(networkDeviceIdNumeric), float32(portIdNumeric), float32(virtualFunctionIdNumeric)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(virtualFunction, &networkDeviceVirtualFunctionPrintConfig)
}
