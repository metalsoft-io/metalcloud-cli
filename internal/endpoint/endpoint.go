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

func GetEndpointId(endpointId string) (int32, error) {
	endpointIdNumeric, err := strconv.ParseFloat(endpointId, 32)
	if err != nil {
		err := fmt.Errorf("invalid endpoint ID: '%s'", endpointId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return int32(endpointIdNumeric), nil
}
