package fabric

import (
	"context"
	"fmt"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var bgpSessionPrintConfig = formatter.PrintConfig{
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
		// The API populates either the scalar id or only the embedded object,
		// depending on the endpoint; fall back to the embedded object's id.
		"NetworkFabricLinkId|NetworkFabricLink.Id": {
			Title: "Link",
			Order: 4,
		},
		"NetworkFabricLinkAggregationId|LinkAggregation.Id": {
			Title: "Link Aggregation",
			Order: 5,
		},
		"BgpNumbering": {
			Title: "Numbering",
			Order: 6,
		},
		"BgpLinkConfiguration": {
			Title: "Link Config",
			Order: 7,
		},
		"Status": {
			Transformer: formatter.FormatStatusValue,
			Order:       8,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       9,
		},
	},
}

// BgpSessionList lists the BGP sessions of a fabric.
func BgpSessionList(ctx context.Context, fabricIdOrLabel string) error {
	logger.Get().Info().Msgf("Listing BGP sessions of fabric '%s'", fabricIdOrLabel)

	fabricId, err := ResolveFabricNumericId(ctx, fabricIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.NetworkFabricAPI.GetNetworkFabricBgpSessions(ctx, fabricId).SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &bgpSessionPrintConfig)
}

// BgpSessionGet shows one BGP session of a fabric.
func BgpSessionGet(ctx context.Context, fabricIdOrLabel string, bgpSessionId string) error {
	logger.Get().Info().Msgf("Get BGP session '%s' of fabric '%s'", bgpSessionId, fabricIdOrLabel)

	fabricId, sessionId, err := resolveFabricAndChildId(ctx, fabricIdOrLabel, bgpSessionId, "BGP session")
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	bgpSession, httpRes, err := client.NetworkFabricAPI.GetNetworkFabricBGPSession(ctx, fabricId, sessionId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(bgpSession, &bgpSessionPrintConfig)
}

// BgpSessionCreate creates a BGP session on a fabric from a JSON/YAML config.
func BgpSessionCreate(ctx context.Context, fabricIdOrLabel string, config []byte) error {
	logger.Get().Info().Msgf("Creating BGP session on fabric '%s'", fabricIdOrLabel)

	fabricId, err := ResolveFabricNumericId(ctx, fabricIdOrLabel)
	if err != nil {
		return err
	}

	var createSession sdk.CreateNetworkFabricBGPSession
	if err := utils.UnmarshalContent(config, &createSession); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	bgpSession, httpRes, err := client.NetworkFabricAPI.
		CreateNetworkFabricBgpSession(ctx, fabricId).
		CreateNetworkFabricBGPSession(createSession).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(bgpSession, &bgpSessionPrintConfig)
}

// BgpSessionUpdate updates a BGP session from a JSON/YAML config.
func BgpSessionUpdate(ctx context.Context, fabricIdOrLabel string, bgpSessionId string, config []byte) error {
	logger.Get().Info().Msgf("Updating BGP session '%s' of fabric '%s'", bgpSessionId, fabricIdOrLabel)

	fabricId, sessionId, err := resolveFabricAndChildId(ctx, fabricIdOrLabel, bgpSessionId, "BGP session")
	if err != nil {
		return err
	}

	var updateSession sdk.UpdateNetworkFabricBGPSession
	if err := utils.UnmarshalContent(config, &updateSession); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	bgpSession, httpRes, err := client.NetworkFabricAPI.
		UpdateNetworkFabricBgpSession(ctx, fabricId, sessionId).
		UpdateNetworkFabricBGPSession(updateSession).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(bgpSession, &bgpSessionPrintConfig)
}

// BgpSessionDelete removes a BGP session from a fabric.
func BgpSessionDelete(ctx context.Context, fabricIdOrLabel string, bgpSessionId string) error {
	logger.Get().Info().Msgf("Deleting BGP session '%s' of fabric '%s'", bgpSessionId, fabricIdOrLabel)

	fabricId, sessionId, err := resolveFabricAndChildId(ctx, fabricIdOrLabel, bgpSessionId, "BGP session")
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkFabricAPI.DeleteNetworkFabricBgpSession(ctx, fabricId, sessionId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("BGP session '%s' deleted from fabric '%s'", bgpSessionId, fabricIdOrLabel)
	return nil
}

// BgpSessionConfigExample prints an example BGP session create body.
func BgpSessionConfigExample(ctx context.Context) error {
	// bgpNumbering is one of inherited, numbered, unnumbered; bgpLinkConfiguration
	// one of disabled, active, passive. Either linkId or linkAggregationId is set.
	example := sdk.CreateNetworkFabricBGPSession{
		BgpNumbering:         "numbered",
		BgpLinkConfiguration: "active",
		LinkId:               sdk.PtrInt64(1),
		CustomVariables:      map[string]interface{}{"example_variable": "example_value"},
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(example, nil)
	}
	return formatter.PrintYamlResult(example)
}

// resolveFabricAndChildId resolves a fabric reference (ID or name) and a numeric
// sub-resource ID, so the fabric sub-resource commands accept fabric labels too.
func resolveFabricAndChildId(ctx context.Context, fabricIdOrLabel string, childId string, childName string) (int64, int64, error) {
	fabricId, err := ResolveFabricNumericId(ctx, fabricIdOrLabel)
	if err != nil {
		return 0, 0, err
	}

	childIdNumeric, err := utils.GetInt64FromString(childId)
	if err != nil {
		err = fmt.Errorf("invalid %s ID: '%s'", childName, childId)
		logger.Get().Error().Err(err).Msg("")
		return 0, 0, err
	}

	return fabricId, childIdNumeric, nil
}
