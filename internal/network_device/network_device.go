package network_device

import (
	"context"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
)

var networkDevicePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"IdentifierString": {
			Title: "Identifier",
			Order: 2,
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

func NetworkDeviceList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all network devices")

	client := api.GetApiClient(ctx)

	networkDeviceList, httpRes, err := client.NetworkDeviceAPI.GetNetworkDevices(ctx).SortBy([]string{"id:ASC"}).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkDeviceList, &networkDevicePrintConfig)
}

func NetworkDeviceGet(ctx context.Context, networkDeviceId string) error {
	logger.Get().Info().Msgf("Get network device %s details", networkDeviceId)

	networkDeviceIdNumeric, err := utils.GetFloat32FromString(networkDeviceId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	networkDevice, httpRes, err := client.NetworkDeviceAPI.GetNetworkDevice(ctx, networkDeviceIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(networkDevice, &networkDevicePrintConfig)
}
