package point_to_point_link

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// linkConfigPrintConfig renders the staged config object of a link. The nested
// ipv4/ipv6 blocks are printed as compact JSON so they fit in a table cell.
var linkConfigPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"DeployType": {
			Title: "Deploy Type",
			Order: 2,
		},
		"DeployStatus": {
			Title:       "Deploy Status",
			Transformer: formatter.FormatStatusValue,
			Order:       3,
		},
		"Mtu": {
			Title: "MTU",
			Order: 4,
		},
		"Ipv4": {
			Title:       "IPv4",
			Order:       5,
			MaxWidth:    60,
			Transformer: formatCompactJSONValue,
		},
		"Ipv6": {
			Title:       "IPv6",
			Order:       6,
			MaxWidth:    60,
			Transformer: formatCompactJSONValue,
		},
		"Revision": {
			Title: "Revision",
			Order: 7,
		},
	},
}

var staticRoutePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"DestinationPrefix": {
			Title: "Destination Prefix",
			Order: 2,
		},
	},
}

// formatCompactJSONValue renders a nested value as compact JSON for table output.
func formatCompactJSONValue(value interface{}) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprint(value)
	}
	return string(encoded)
}

// PointToPointLinkUpdate updates a link's own properties (label, name,
// description, annotations) from a JSON/YAML config. The link's revision is
// sent as If-Match.
func PointToPointLinkUpdate(ctx context.Context, linkId string, config []byte) error {
	logger.Get().Info().Msgf("Updating point-to-point link %s", linkId)

	linkIdNumeric, revision, err := getPointToPointLinkIdAndRevision(ctx, linkId)
	if err != nil {
		return err
	}

	// Validate the user's body against the SDK model, then forward it as raw
	// JSON: the typed response is the same oneOf that fails to decode once the
	// link carries an ipv4 strategy.
	var update sdk.UpdatePointToPointLink
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	body, err := json.Marshal(update)
	if err != nil {
		return err
	}

	result, err := api.RawJSONRequest(ctx, http.MethodPatch,
		fmt.Sprintf("/api/v2/point-to-point-links/%d", int64(linkIdNumeric)), body, api.IfMatchHeader(revision))
	if err != nil {
		return err
	}

	var link p2pLinkDisplay
	if err := json.Unmarshal(result, &link); err != nil {
		return fmt.Errorf("parsing updated point-to-point link: %w", err)
	}

	return formatter.PrintResult(link, &PointToPointLinkPrintConfig)
}

// PointToPointLinkConfigExampleUpdate prints an example link update body.
func PointToPointLinkConfigExampleUpdate(ctx context.Context) error {
	example := sdk.UpdatePointToPointLink{
		Label:       sdk.PtrString("leaf01-swp33s0-to-spine01-swp1s0"),
		Name:        sdk.PtrString("leaf01 to spine01"),
		Description: sdk.PtrString("uplink"),
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(example, nil)
	}
	return formatter.PrintYamlResult(example)
}

// PointToPointLinkConfigGet shows the staged configuration object of a link.
func PointToPointLinkConfigGet(ctx context.Context, linkId string) error {
	logger.Get().Info().Msgf("Get configuration of point-to-point link %s", linkId)

	linkIdNumeric, err := getPointToPointLinkId(linkId)
	if err != nil {
		return err
	}

	body, err := api.RawJSONRequest(ctx, http.MethodGet, linkConfigPath(linkIdNumeric), nil, nil)
	if err != nil {
		return err
	}

	// The config embeds the subnet allocation strategies, whose typed model is
	// the oneOf that fails strict decode - parse the raw object instead.
	return utils.PrintRawObject(body, &linkConfigPrintConfig)
}

// PointToPointLinkConfigUpdate updates the staged configuration of a link (e.g.
// its MTU). The config object carries its own revision, which is what the API
// guards this endpoint with - not the link's revision.
func PointToPointLinkConfigUpdate(ctx context.Context, linkId string, config []byte) error {
	logger.Get().Info().Msgf("Updating configuration of point-to-point link %s", linkId)

	linkIdNumeric, revision, err := getLinkConfigIdAndRevision(ctx, linkId)
	if err != nil {
		return err
	}

	var update sdk.UpdatePointToPointLinkConfig
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	body, err := json.Marshal(update)
	if err != nil {
		return err
	}

	result, err := api.RawJSONRequest(ctx, http.MethodPatch, linkConfigPath(linkIdNumeric), body, api.IfMatchHeader(revision))
	if err != nil {
		return err
	}

	return utils.PrintRawObject(result, &linkConfigPrintConfig)
}

// PointToPointLinkConfigExampleConfig prints an example config update body.
func PointToPointLinkConfigExampleConfig(ctx context.Context) error {
	example := map[string]interface{}{"mtu": 9216}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(example, nil)
	}
	return formatter.PrintYamlResult(example)
}

// StaticRouteList lists the staged static routes of one address family.
func StaticRouteList(ctx context.Context, linkId string, family string) error {
	logger.Get().Info().Msgf("Listing %s static routes of point-to-point link %s", family, linkId)

	path, err := staticRoutesPath(linkId, family)
	if err != nil {
		return err
	}

	body, err := api.RawJSONRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return err
	}

	rawItems, err := utils.DecodeRawItems(body)
	if err != nil {
		return err
	}

	routes, err := utils.UnmarshalRawItems[sdk.PointToPointStaticRoute](rawItems)
	if err != nil {
		return fmt.Errorf("parsing static routes: %w", err)
	}

	return utils.PrintAllRaw(rawItems, routes, sdk.PaginatedResponseMeta{}, len(routes), &staticRoutePrintConfig)
}

// StaticRouteGet shows one staged static route.
func StaticRouteGet(ctx context.Context, linkId string, family string, routeId string) error {
	logger.Get().Info().Msgf("Get %s static route %s of point-to-point link %s", family, routeId, linkId)

	path, err := staticRoutePath(linkId, family, routeId)
	if err != nil {
		return err
	}

	body, err := api.RawJSONRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return err
	}

	return utils.PrintRawObject(body, &staticRoutePrintConfig)
}

// StaticRouteAdd stages a static route on one address family of a link. The
// config object's revision guards the write.
func StaticRouteAdd(ctx context.Context, linkId string, family string, destinationPrefix string) error {
	logger.Get().Info().Msgf("Adding %s static route '%s' to point-to-point link %s", family, destinationPrefix, linkId)

	path, err := staticRoutesPath(linkId, family)
	if err != nil {
		return err
	}

	_, revision, err := getLinkConfigIdAndRevision(ctx, linkId)
	if err != nil {
		return err
	}

	body, err := json.Marshal(sdk.CreatePointToPointStaticRoute{DestinationPrefix: destinationPrefix})
	if err != nil {
		return err
	}

	result, err := api.RawJSONRequest(ctx, http.MethodPost, path, body, api.IfMatchHeader(revision))
	if err != nil {
		return err
	}

	return utils.PrintRawObject(result, &staticRoutePrintConfig)
}

// StaticRouteRemove removes one staged static route.
func StaticRouteRemove(ctx context.Context, linkId string, family string, routeId string) error {
	logger.Get().Info().Msgf("Removing %s static route %s from point-to-point link %s", family, routeId, linkId)

	path, err := staticRoutePath(linkId, family, routeId)
	if err != nil {
		return err
	}

	_, revision, err := getLinkConfigIdAndRevision(ctx, linkId)
	if err != nil {
		return err
	}

	if _, err := api.RawJSONRequest(ctx, http.MethodDelete, path, nil, api.IfMatchHeader(revision)); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Static route %s removed from point-to-point link %s", routeId, linkId)
	return nil
}

// StaticRouteFamilies are the address families accepted by the static route commands.
var StaticRouteFamilies = []string{"ipv4", "ipv6"}

func validateFamily(family string) error {
	for _, valid := range StaticRouteFamilies {
		if family == valid {
			return nil
		}
	}
	err := fmt.Errorf("invalid address family '%s'; valid values are: %s", family, strings.Join(StaticRouteFamilies, ", "))
	logger.Get().Error().Err(err).Msg("")
	return err
}

func linkConfigPath(linkIdNumeric float32) string {
	return fmt.Sprintf("/api/v2/point-to-point-links/%d/config", int64(linkIdNumeric))
}

func staticRoutesPath(linkId string, family string) (string, error) {
	if err := validateFamily(family); err != nil {
		return "", err
	}

	linkIdNumeric, err := getPointToPointLinkId(linkId)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/static-routes", linkConfigPath(linkIdNumeric), family), nil
}

func staticRoutePath(linkId string, family string, routeId string) (string, error) {
	collection, err := staticRoutesPath(linkId, family)
	if err != nil {
		return "", err
	}

	routeIdNumeric, err := utils.GetInt64FromString(routeId)
	if err != nil {
		err = fmt.Errorf("invalid static route ID: '%s'", routeId)
		logger.Get().Error().Err(err).Msg("")
		return "", err
	}

	return fmt.Sprintf("%s/%d", collection, routeIdNumeric), nil
}

// getLinkConfigIdAndRevision returns the link's numeric ID and the revision of
// its config object. Config sub-resources (static routes, allocation
// strategies) are guarded by the config's revision, which differs from the
// link's own.
func getLinkConfigIdAndRevision(ctx context.Context, linkId string) (float32, string, error) {
	linkIdNumeric, err := getPointToPointLinkId(linkId)
	if err != nil {
		return 0, "", err
	}

	body, err := api.RawJSONRequest(ctx, http.MethodGet, linkConfigPath(linkIdNumeric), nil, nil)
	if err != nil {
		return 0, "", err
	}

	object, err := utils.DecodeRawObject(body)
	if err != nil {
		return 0, "", err
	}

	revision := utils.RevisionFromRaw(object)
	if revision == "" {
		// Older deployments return the config without a revision; fall back to
		// the link's own revision so the write is still guarded.
		_, revision, err = getPointToPointLinkIdAndRevision(ctx, linkId)
		if err != nil {
			return 0, "", err
		}
	}

	return linkIdNumeric, revision, nil
}
