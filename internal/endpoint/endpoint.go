package endpoint

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

type EndpointInterfaceDetails struct {
	Id                         *float32
	MacAddress                 *string
	NetworkDeviceId            *float32
	NetworkDeviceInterfaceId   *float32
	NetworkDeviceName          *string
	NetworkDeviceInterfaceName *string
}

var endpointPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
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
			Title: "#",
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

	request := client.EndpointAPI.GetEndpoints(ctx)

	if len(filterSite) > 0 {
		request = request.FilterSiteId(utils.ProcessFilterStringSlice(filterSite))
	}

	if len(filterExternalId) > 0 {
		request = request.FilterExternalId(utils.ProcessFilterStringSlice(filterExternalId))
	}

	endpointList, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(endpointList, &endpointPrintConfig)
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

	endpointInterfaces, httpRes, err := client.EndpointAPI.GetEndpointInterfaces(ctx, endpointIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	endpointInterfacesList := make([]EndpointInterfaceDetails, 0, len(endpointInterfaces.Data))
	for _, iface := range endpointInterfaces.Data {
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

func GetEndpointId(endpointId string) (int32, error) {
	endpointIdNumeric, err := strconv.ParseFloat(endpointId, 32)
	if err != nil {
		err := fmt.Errorf("invalid endpoint ID: '%s'", endpointId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return int32(endpointIdNumeric), nil
}
