package point_to_point_link

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var PointToPointLinkPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Label": {
			Title: "Label",
			Order: 2,
		},
		"Description": {
			Title: "Description",
			Order: 3,
		},
		"RoutingActivation": {
			Title: "Routing Activation",
			Order: 4,
		},
		"ServiceStatus": {
			Title: "Status",
			Order: 5,
		},
		"Revision": {
			Title: "Revision",
			Order: 6,
		},
	},
}

// p2pLinkDisplay is a lenient view of a point-to-point link for output. The
// SDK's typed decode of GetPointToPointLinks fails ("data matches more than one
// schema in oneOf(...201Response)") whenever a link carries an ipv4 strategy -
// the oneOf's "unnumbered" variant ({kind, scope}) also matches auto/manual
// strategies. We parse the raw body for the displayed fields instead. JSON null
// (e.g. an unset description) decodes into the zero value, which is fine here.
type p2pLinkDisplay struct {
	Id                int64  `json:"id"`
	Label             string `json:"label"`
	Description       string `json:"description"`
	RoutingActivation string `json:"routingActivation"`
	ServiceStatus     string `json:"serviceStatus"`
	Revision          int64  `json:"revision"`
}

// PointToPointLinkList lists point-to-point links, optionally filtered by a
// referenced interface id or route domain id.
func PointToPointLinkList(ctx context.Context, interfaceId string, routeDomainId string) error {
	logger.Get().Info().Msgf("Listing point-to-point links")

	path := "/api/v2/point-to-point-links"
	query := url.Values{}
	if interfaceId != "" {
		query.Set("interfaceId", interfaceId)
	}
	if routeDomainId != "" {
		query.Set("routeDomainId", routeDomainId)
	}
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	httpRes, err := api.DoJSONRequest(ctx, http.MethodGet, path, nil)
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}
	defer httpRes.Body.Close()
	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return fmt.Errorf("reading point-to-point links: %w", err)
	}
	var links []p2pLinkDisplay
	if err := json.Unmarshal(body, &links); err != nil {
		return fmt.Errorf("parsing point-to-point links: %w", err)
	}

	return formatter.PrintResult(links, &PointToPointLinkPrintConfig)
}

// PointToPointLinkGet shows a single point-to-point link.
func PointToPointLinkGet(ctx context.Context, linkId string) error {
	logger.Get().Info().Msgf("Get point-to-point link %s details", linkId)

	linkIdNumeric, err := getPointToPointLinkId(linkId)
	if err != nil {
		return err
	}

	httpRes, err := api.DoJSONRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v2/point-to-point-links/%d", int64(linkIdNumeric)), nil)
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}
	defer httpRes.Body.Close()
	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return fmt.Errorf("reading point-to-point link: %w", err)
	}
	var link p2pLinkDisplay
	if err := json.Unmarshal(body, &link); err != nil {
		return fmt.Errorf("parsing point-to-point link: %w", err)
	}

	return formatter.PrintResult(link, &PointToPointLinkPrintConfig)
}

// PointToPointLinkCreate creates a point-to-point link from a raw configuration
// body. The body may stage one or more IPv4/IPv6 subnet allocation strategies
// inline (ipv4.subnetAllocationStrategies), so a fully-connected link plus its
// manual /31 strategy can be created in a single call.
func PointToPointLinkCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating point-to-point link")

	// Parse into a generic map and forward as raw JSON rather than round-tripping
	// through sdk.CreatePointToPointLink: the SDK re-serializes a manual strategy's
	// scope with resourceId:0, which the API rejects for a global scope. Sending
	// the user's body verbatim preserves an omitted (global) resourceId.
	var linkConfig map[string]interface{}
	if err := utils.UnmarshalContent(config, &linkConfig); err != nil {
		return err
	}

	body, err := json.Marshal(linkConfig)
	if err != nil {
		return err
	}

	httpRes, err := api.DoJSONRequest(ctx, http.MethodPost, "/api/v2/point-to-point-links", body)
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}
	defer httpRes.Body.Close()
	respBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return fmt.Errorf("reading created point-to-point link: %w", err)
	}
	var link p2pLinkDisplay
	if err := json.Unmarshal(respBody, &link); err != nil {
		return fmt.Errorf("parsing created point-to-point link: %w", err)
	}

	return formatter.PrintResult(link, &PointToPointLinkPrintConfig)
}

// PointToPointLinkDelete deletes a point-to-point link.
func PointToPointLinkDelete(ctx context.Context, linkId string) error {
	logger.Get().Info().Msgf("Deleting point-to-point link %s", linkId)

	linkIdNumeric, revision, err := getPointToPointLinkIdAndRevision(ctx, linkId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.PointToPointLinkAPI.
		DeletePointToPointLink(ctx, linkIdNumeric).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Point-to-point link %s deleted", linkId)
	return nil
}

// PointToPointLinkAddIpv4Strategy attaches a manual IPv4 subnet allocation
// strategy to an existing link. The link's current revision is sent as If-Match.
func PointToPointLinkAddIpv4Strategy(ctx context.Context, linkId string, subnetId int64, binding string, scopeKind string, scopeResourceId int64) error {
	logger.Get().Info().Msgf("Adding manual IPv4 allocation strategy (subnet %d) to point-to-point link %s", subnetId, linkId)

	linkIdNumeric, revision, err := getPointToPointLinkIdAndRevision(ctx, linkId)
	if err != nil {
		return err
	}

	if _, err := sdk.NewPointToPointInterfaceBindingFromValue(binding); err != nil {
		return fmt.Errorf("invalid interface binding '%s': %w", binding, err)
	}

	if _, err := sdk.NewResourceScopeKindFromValue(scopeKind); err != nil {
		return fmt.Errorf("invalid scope kind '%s': %w", scopeKind, err)
	}

	// Built as a raw body: the SDK's CreateResourceScope always serializes
	// resourceId, but a global scope must omit it (the API rejects resourceId:0
	// with "scope.resourceId must not be less than 1"). Only non-global scopes
	// carry a resourceId.
	scope := map[string]any{"kind": scopeKind}
	if scopeKind != string(sdk.RESOURCESCOPEKIND_GLOBAL) {
		scope["resourceId"] = scopeResourceId
	}
	strategy := map[string]any{
		"kind":              string(sdk.POINTTOPOINTALLOCATIONSTRATEGYKIND_MANUAL),
		"subnetId":          subnetId,
		"scope":             scope,
		"interfaceABinding": binding,
	}

	body, err := json.Marshal(strategy)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/api/v2/point-to-point-links/%d/config/ipv4/subnet-allocation-strategies", int64(linkIdNumeric))
	httpRes, err := api.DoJSONRequestWithHeaders(ctx, http.MethodPost, path, body,
		map[string]string{"If-Match": revision})
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}
	defer httpRes.Body.Close()

	body, err = io.ReadAll(httpRes.Body)
	if err != nil {
		return fmt.Errorf("reading created strategy: %w", err)
	}
	// Decode into a generic map: the typed 201 response is a oneOf that fails
	// strict decode (same schema-drift bug as the list/get paths).
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("parsing created strategy: %w", err)
	}
	return formatter.PrintResult(result, nil)
}

// PointToPointLinkConfigExample prints an example create body, including a staged
// manual IPv4 /31 allocation strategy on interface A.
func PointToPointLinkConfigExample(ctx context.Context) error {
	// Built as a map so the manual strategy's global scope prints as
	// {"kind":"global"} with no resourceId - the shape the API expects (a global
	// scope must omit resourceId, and the create path forwards the body verbatim).
	example := map[string]interface{}{
		"description":       "leaf-swp33s0 to spine-swp1s0",
		"routingActivation": string(sdk.POINTTOPOINTLINKROUTINGACTIVATION_DEFAULT),
		"mtu":               9216,
		"interfaceA": map[string]interface{}{
			"type":        string(sdk.POINTTOPOINTINTERFACETYPE_NETWORK_EQUIPMENT_INTERFACE),
			"interfaceId": 1001,
		},
		"interfaceB": map[string]interface{}{
			"type":        string(sdk.POINTTOPOINTINTERFACETYPE_NETWORK_EQUIPMENT_INTERFACE),
			"interfaceId": 2002,
		},
		"ipv4": map[string]interface{}{
			"subnetAllocationStrategies": []map[string]interface{}{
				{
					"kind":              string(sdk.POINTTOPOINTALLOCATIONSTRATEGYKIND_MANUAL),
					"scope":             map[string]interface{}{"kind": string(sdk.RESOURCESCOPEKIND_GLOBAL)},
					"subnetId":          12345,
					"interfaceABinding": string(sdk.POINTTOPOINTINTERFACEBINDING_A_FIRST),
				},
			},
		},
	}

	return formatter.PrintResult(example, nil)
}

func getPointToPointLinkId(linkId string) (float32, error) {
	linkIdNumeric, err := strconv.ParseInt(linkId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid point-to-point link ID: '%s'", linkId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(linkIdNumeric), nil
}

func getPointToPointLinkIdAndRevision(ctx context.Context, linkId string) (float32, string, error) {
	linkIdNumeric, err := getPointToPointLinkId(linkId)
	if err != nil {
		return 0, "", err
	}

	// Fetch the revision via a raw GET: the typed decode trips the same
	// oneOf("...201Response") bug as the listing once the link has an ipv4 strategy.
	httpRes, err := api.DoJSONRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v2/point-to-point-links/%d", int64(linkIdNumeric)), nil)
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}
	defer httpRes.Body.Close()
	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return 0, "", fmt.Errorf("reading point-to-point link: %w", err)
	}
	var link p2pLinkDisplay
	if err := json.Unmarshal(body, &link); err != nil {
		return 0, "", fmt.Errorf("parsing point-to-point link: %w", err)
	}

	return linkIdNumeric, strconv.FormatInt(link.Revision, 10), nil
}
