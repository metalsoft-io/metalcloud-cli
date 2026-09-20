package fabric

import (
	"context"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var linkAggregationPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"NetworkFabricId": {
			Title: "Fabric",
			Order: 2,
		},
		"Name": {
			Order:    3,
			MaxWidth: 40,
		},
		"Type": {
			Order: 4,
		},
		"MlagDomainIdentifier": {
			Title: "MLAG Domain",
			Order: 5,
		},
		"Status": {
			Transformer: formatter.FormatStatusValue,
			Order:       6,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
	},
}

// LinkAggregationList lists the link aggregations of a fabric.
func LinkAggregationList(ctx context.Context, fabricIdOrLabel string) error {
	logger.Get().Info().Msgf("Listing link aggregations of fabric '%s'", fabricIdOrLabel)

	fabricId, err := ResolveFabricNumericId(ctx, fabricIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.NetworkFabricAPI.GetNetworkFabricLinkAggregations(ctx, fabricId).SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &linkAggregationPrintConfig)
}

// LinkAggregationGet shows one link aggregation of a fabric.
func LinkAggregationGet(ctx context.Context, fabricIdOrLabel string, linkAggregationId string) error {
	logger.Get().Info().Msgf("Get link aggregation '%s' of fabric '%s'", linkAggregationId, fabricIdOrLabel)

	fabricId, aggregationId, err := resolveFabricAndChildId(ctx, fabricIdOrLabel, linkAggregationId, "link aggregation")
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	linkAggregation, httpRes, err := client.NetworkFabricAPI.GetNetworkFabricLinkAggregation(ctx, fabricId, aggregationId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(linkAggregation, &linkAggregationPrintConfig)
}

// LinkAggregationCreate creates a link aggregation on a fabric from a JSON/YAML config.
func LinkAggregationCreate(ctx context.Context, fabricIdOrLabel string, config []byte) error {
	logger.Get().Info().Msgf("Creating link aggregation on fabric '%s'", fabricIdOrLabel)

	fabricId, err := ResolveFabricNumericId(ctx, fabricIdOrLabel)
	if err != nil {
		return err
	}

	var createAggregation sdk.CreateNetworkFabricLinkAggregation
	if err := utils.UnmarshalContent(config, &createAggregation); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	linkAggregation, httpRes, err := client.NetworkFabricAPI.
		CreateNetworkFabricLinkAggregation(ctx, fabricId).
		CreateNetworkFabricLinkAggregation(createAggregation).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(linkAggregation, &linkAggregationPrintConfig)
}

// LinkAggregationUpdate updates a link aggregation from a JSON/YAML config.
func LinkAggregationUpdate(ctx context.Context, fabricIdOrLabel string, linkAggregationId string, config []byte) error {
	logger.Get().Info().Msgf("Updating link aggregation '%s' of fabric '%s'", linkAggregationId, fabricIdOrLabel)

	fabricId, aggregationId, err := resolveFabricAndChildId(ctx, fabricIdOrLabel, linkAggregationId, "link aggregation")
	if err != nil {
		return err
	}

	var updateAggregation sdk.UpdateNetworkFabricLinkAggregation
	if err := utils.UnmarshalContent(config, &updateAggregation); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	linkAggregation, httpRes, err := client.NetworkFabricAPI.
		UpdateNetworkFabricLinkAggregation(ctx, fabricId, aggregationId).
		UpdateNetworkFabricLinkAggregation(updateAggregation).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(linkAggregation, &linkAggregationPrintConfig)
}

// LinkAggregationDelete removes a link aggregation from a fabric.
func LinkAggregationDelete(ctx context.Context, fabricIdOrLabel string, linkAggregationId string) error {
	logger.Get().Info().Msgf("Deleting link aggregation '%s' of fabric '%s'", linkAggregationId, fabricIdOrLabel)

	fabricId, aggregationId, err := resolveFabricAndChildId(ctx, fabricIdOrLabel, linkAggregationId, "link aggregation")
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkFabricAPI.DeleteNetworkFabricLinkAggregation(ctx, fabricId, aggregationId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Link aggregation '%s' deleted from fabric '%s'", linkAggregationId, fabricIdOrLabel)
	return nil
}

// LinkAggregationConfigExample prints an example link aggregation create body.
func LinkAggregationConfigExample(ctx context.Context) error {
	// type is one of lag, mlag, mlag-peer-link; mlagDomainIdentifier applies to
	// mlag-peer-link only. linkIds are network fabric link IDs (see 'fabric get-links').
	example := sdk.CreateNetworkFabricLinkAggregation{
		Type:            "lag",
		LinkIds:         []int64{1, 2},
		CustomVariables: map[string]interface{}{"example_variable": "example_value"},
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(example, nil)
	}
	return formatter.PrintYamlResult(example)
}
