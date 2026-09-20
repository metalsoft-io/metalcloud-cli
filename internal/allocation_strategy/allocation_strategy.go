// Package allocation_strategy implements the allocation-strategy sub-resources
// shared by logical networks, logical network profiles, route domains and
// point-to-point links. Every parent exposes the same five operations per
// strategy family (list, get, create, replace, delete) under
// <parent path>/<family segment>[/<strategy id>], so a single generic
// implementation driven by a Parent description serves all of them.
//
// Requests and responses are raw JSON: the SDK models these strategies as
// polymorphic oneOf unions (auto/manual/unnumbered) whose strict decoders
// reject valid API payloads, and the create bodies are passed through from
// the user's --config-source untouched.
package allocation_strategy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// Family is one allocation strategy collection (e.g. VLAN) under a parent.
type Family struct {
	// Name is the value the user passes on the command line.
	Name string
	// PathSegment is appended to the parent's base path, e.g.
	// "vlan/vlan-allocation-strategies".
	PathSegment string
	// Paginated is true when the list endpoint returns a {data, meta}
	// envelope with page/limit parameters, false when it returns a bare array.
	Paginated bool
	// Examples are the create bodies shown by config-example, one per kind.
	Examples []map[string]any
}

// Parent describes a resource that owns allocation strategy collections.
type Parent struct {
	// Name is used in log and error messages, e.g. "logical network".
	Name string
	// ItemPath returns the API path of one parent, e.g.
	// "/api/v2/logical-networks/12".
	ItemPath func(id int64) string
	// ConfigSuffix is inserted between the parent path and the family segment
	// ("/config" for resources whose strategies live under their config).
	ConfigSuffix string
	// Families lists the strategy collections this parent supports.
	Families []Family
}

// FamilyNames returns the family names in declaration order, for help text.
func (p Parent) FamilyNames() []string {
	names := make([]string, 0, len(p.Families))
	for _, f := range p.Families {
		names = append(names, f.Name)
	}
	return names
}

func (p Parent) family(name string) (Family, error) {
	for _, f := range p.Families {
		if f.Name == name {
			return f, nil
		}
	}
	err := fmt.Errorf("unknown %s allocation strategy family '%s'; valid values are: %s",
		p.Name, name, strings.Join(p.FamilyNames(), ", "))
	logger.Get().Error().Err(err).Msg("")
	return Family{}, err
}

// resolve parses the parent reference as a numeric ID and fetches the current
// revision for If-Match. Strategies have no revision of their own: when they
// live under the parent's config, the API guards them with the config
// object's revision (which differs from the entity's), otherwise with the
// parent's own revision.
func (p Parent) resolve(ctx context.Context, parentRef string) (int64, string, error) {
	id, err := utils.GetInt64FromString(parentRef)
	if err != nil {
		err = fmt.Errorf("invalid %s ID: '%s'", p.Name, parentRef)
		logger.Get().Error().Err(err).Msg("")
		return 0, "", err
	}

	revision, err := fetchRevision(ctx, p.ItemPath(id)+p.ConfigSuffix)
	if err != nil {
		return 0, "", err
	}
	if revision == "" && p.ConfigSuffix != "" {
		if revision, err = fetchRevision(ctx, p.ItemPath(id)); err != nil {
			return 0, "", err
		}
	}

	return id, revision, nil
}

func fetchRevision(ctx context.Context, path string) (string, error) {
	body, err := api.RawJSONRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return "", err
	}
	object, err := utils.DecodeRawObject(body)
	if err != nil {
		return "", err
	}
	return utils.RevisionFromRaw(object), nil
}

func (p Parent) collectionPath(id int64, f Family) string {
	return p.ItemPath(id) + p.ConfigSuffix + "/" + f.PathSegment
}

func (p Parent) itemPath(id int64, f Family, strategyId string) (string, error) {
	strategyIdNumeric, err := utils.GetInt64FromString(strategyId)
	if err != nil {
		err = fmt.Errorf("invalid allocation strategy ID: '%s'", strategyId)
		logger.Get().Error().Err(err).Msg("")
		return "", err
	}
	return fmt.Sprintf("%s/%d", p.collectionPath(id, f), strategyIdNumeric), nil
}

// strategyPrintConfig renders the common strategy fields; family-specific
// fields (vlanId, subnetId, prefixLength, ...) are folded into Details.
var strategyPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Kind": {
			Order: 2,
		},
		"Scope.Kind": {
			Title: "Scope",
			Order: 3,
		},
		"Scope.ResourceId": {
			Title: "Scope Resource",
			Order: 4,
		},
		"Details": {
			Order:    5,
			MaxWidth: 80,
		},
	},
}

// commonKeys are rendered as dedicated columns or are noise; everything else
// is family-specific and goes into the Details column.
var commonKeys = map[string]bool{
	"id": true, "kind": true, "scope": true, "createdAt": true, "updatedAt": true,
	"createdTimestamp": true, "updatedTimestamp": true, "links": true,
}

// strategyDisplay is the typed row rendered in table formats. The formatter
// tabulates struct slices, not map slices, so the raw object is projected here.
type strategyDisplay struct {
	Id    any
	Kind  string
	Scope struct {
		Kind       string
		ResourceId any
	}
	Details string
}

// toDisplayRecord decodes a raw strategy into a table row, folding the
// family-specific fields into a compact "key=value" Details string.
func toDisplayRecord(raw json.RawMessage) (strategyDisplay, error) {
	var row strategyDisplay

	record, err := utils.DecodeRawObject(raw)
	if err != nil {
		return row, err
	}

	row.Id = record["id"]
	row.Kind, _ = record["kind"].(string)
	if scope, ok := record["scope"].(map[string]any); ok {
		row.Scope.Kind, _ = scope["kind"].(string)
		row.Scope.ResourceId = scope["resourceId"]
	}

	keys := make([]string, 0, len(record))
	for key := range record {
		if !commonKeys[key] {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		value, _ := json.Marshal(record[key])
		parts = append(parts, fmt.Sprintf("%s=%s", key, value))
	}
	row.Details = strings.Join(parts, " ")

	return row, nil
}

func printStrategies(rawItems []json.RawMessage, meta sdk.PaginatedResponseMeta) error {
	records := make([]strategyDisplay, 0, len(rawItems))
	for _, raw := range rawItems {
		record, err := toDisplayRecord(raw)
		if err != nil {
			return err
		}
		records = append(records, record)
	}
	return utils.PrintAllRaw(rawItems, records, meta, len(records), &strategyPrintConfig)
}

// printStrategy prints one strategy: the raw object in machine formats (lossless),
// the projected row in table formats.
func printStrategy(body []byte) error {
	if formatter.IsNativeFormat() {
		return utils.PrintRawObject(body, nil)
	}
	record, err := toDisplayRecord(body)
	if err != nil {
		return err
	}
	return formatter.PrintResult(record, &strategyPrintConfig)
}

// List prints all strategies of one family under the parent.
func List(ctx context.Context, p Parent, parentRef string, familyName string) error {
	logger.Get().Info().Msgf("Listing %s allocation strategies of %s '%s'", familyName, p.Name, parentRef)

	f, err := p.family(familyName)
	if err != nil {
		return err
	}
	id, _, err := p.resolve(ctx, parentRef)
	if err != nil {
		return err
	}

	path := p.collectionPath(id, f)

	if f.Paginated {
		rawItems, meta, err := utils.FetchAllPagesRaw(func(page float32) (*http.Response, error) {
			return api.DoJSONRequest(ctx, http.MethodGet, fmt.Sprintf("%s?page=%.0f&limit=100", path, page), nil)
		})
		if err != nil {
			return err
		}
		return printStrategies(rawItems, meta)
	}

	body, err := api.RawJSONRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return err
	}
	rawItems, err := utils.DecodeRawItems(body)
	if err != nil {
		return err
	}
	return printStrategies(rawItems, sdk.PaginatedResponseMeta{})
}

// Get prints one strategy.
func Get(ctx context.Context, p Parent, parentRef string, familyName string, strategyId string) error {
	logger.Get().Info().Msgf("Get %s allocation strategy '%s' of %s '%s'", familyName, strategyId, p.Name, parentRef)

	f, err := p.family(familyName)
	if err != nil {
		return err
	}
	id, _, err := p.resolve(ctx, parentRef)
	if err != nil {
		return err
	}
	path, err := p.itemPath(id, f, strategyId)
	if err != nil {
		return err
	}

	body, err := api.RawJSONRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return err
	}
	return printStrategy(body)
}

// Create adds a strategy from a JSON/YAML config body and prints the result.
func Create(ctx context.Context, p Parent, parentRef string, familyName string, config []byte) error {
	logger.Get().Info().Msgf("Creating %s allocation strategy on %s '%s'", familyName, p.Name, parentRef)

	f, err := p.family(familyName)
	if err != nil {
		return err
	}
	id, revision, err := p.resolve(ctx, parentRef)
	if err != nil {
		return err
	}
	body, err := configToJSON(config)
	if err != nil {
		return err
	}

	result, err := api.RawJSONRequest(ctx, http.MethodPost, p.collectionPath(id, f), body, api.IfMatchHeader(revision))
	if err != nil {
		return err
	}
	return printStrategy(result)
}

// Replace overwrites a strategy from a JSON/YAML config body and prints the result.
func Replace(ctx context.Context, p Parent, parentRef string, familyName string, strategyId string, config []byte) error {
	logger.Get().Info().Msgf("Replacing %s allocation strategy '%s' of %s '%s'", familyName, strategyId, p.Name, parentRef)

	f, err := p.family(familyName)
	if err != nil {
		return err
	}
	id, revision, err := p.resolve(ctx, parentRef)
	if err != nil {
		return err
	}
	path, err := p.itemPath(id, f, strategyId)
	if err != nil {
		return err
	}
	body, err := configToJSON(config)
	if err != nil {
		return err
	}

	result, err := api.RawJSONRequest(ctx, http.MethodPut, path, body, api.IfMatchHeader(revision))
	if err != nil {
		return err
	}
	return printStrategy(result)
}

// Delete removes a strategy.
func Delete(ctx context.Context, p Parent, parentRef string, familyName string, strategyId string) error {
	logger.Get().Info().Msgf("Deleting %s allocation strategy '%s' of %s '%s'", familyName, strategyId, p.Name, parentRef)

	f, err := p.family(familyName)
	if err != nil {
		return err
	}
	id, revision, err := p.resolve(ctx, parentRef)
	if err != nil {
		return err
	}
	path, err := p.itemPath(id, f, strategyId)
	if err != nil {
		return err
	}

	if _, err := api.RawJSONRequest(ctx, http.MethodDelete, path, nil, api.IfMatchHeader(revision)); err != nil {
		return err
	}

	logger.Get().Info().Msgf("%s allocation strategy '%s' deleted from %s '%s'", familyName, strategyId, p.Name, parentRef)
	return nil
}

// ConfigExample prints the example create bodies of one family (one per kind).
func ConfigExample(p Parent, familyName string) error {
	f, err := p.family(familyName)
	if err != nil {
		return err
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(f.Examples, nil)
	}
	return formatter.PrintYamlResult(f.Examples)
}

// configToJSON accepts JSON or YAML user input and returns canonical JSON for
// the wire. Going through a generic map keeps the body a pass-through.
func configToJSON(config []byte) ([]byte, error) {
	var object map[string]any
	if err := utils.UnmarshalContent(config, &object); err != nil {
		return nil, err
	}
	return json.Marshal(object)
}

// scopeExample is the scope block shared by every strategy kind. A global
// scope must omit resourceId (the API rejects resourceId 0).
func scopeExample(kind sdk.ResourceScopeKind, resourceId int64) map[string]any {
	scope := map[string]any{"kind": string(kind)}
	if kind != sdk.RESOURCESCOPEKIND_GLOBAL {
		scope["resourceId"] = resourceId
	}
	return scope
}

// Shared family definitions. Logical networks and profiles use the six
// network families; route domains use the L3 families; point-to-point links
// use the point-to-point address families.
var (
	Ipv4Subnet = Family{
		Name: "ipv4", PathSegment: "ipv4/subnet-allocation-strategies", Paginated: true,
		Examples: []map[string]any{
			{"kind": "auto", "scope": scopeExample(sdk.RESOURCESCOPEKIND_FABRIC, 1), "subnetPoolIds": []int64{1}, "prefixLength": 24, "gatewayPlacement": "first"},
			{"kind": "manual", "scope": scopeExample(sdk.RESOURCESCOPEKIND_FABRIC, 1), "subnetId": 1, "gatewayPlacement": "first"},
		},
	}
	Ipv6Subnet = Family{
		Name: "ipv6", PathSegment: "ipv6/subnet-allocation-strategies", Paginated: true,
		Examples: []map[string]any{
			{"kind": "auto", "scope": scopeExample(sdk.RESOURCESCOPEKIND_FABRIC, 1), "subnetPoolIds": []int64{1}, "prefixLength": 64, "gatewayPlacement": "first"},
			{"kind": "manual", "scope": scopeExample(sdk.RESOURCESCOPEKIND_FABRIC, 1), "subnetId": 1, "gatewayPlacement": "first"},
		},
	}
	Vlan = Family{
		Name: "vlan", PathSegment: "vlan/vlan-allocation-strategies", Paginated: true,
		Examples: []map[string]any{
			{"kind": "auto", "scope": scopeExample(sdk.RESOURCESCOPEKIND_FABRIC, 1), "granularityLevel": string(sdk.VLANALLOCATIONGRANULARITYLEVEL_FABRIC)},
			{"kind": "manual", "scope": scopeExample(sdk.RESOURCESCOPEKIND_FABRIC, 1), "granularityLevel": string(sdk.VLANALLOCATIONGRANULARITYLEVEL_FABRIC), "vlanId": 100},
		},
	}
	Vni = Family{
		Name: "vni", PathSegment: "vxlan/vni-allocation-strategies", Paginated: true,
		Examples: []map[string]any{
			{"kind": "auto", "scope": scopeExample(sdk.RESOURCESCOPEKIND_FABRIC, 1)},
			{"kind": "manual", "scope": scopeExample(sdk.RESOURCESCOPEKIND_FABRIC, 1), "vni": 10100},
		},
	}
	Pkey = Family{
		Name: "pkey", PathSegment: "pkey/pkey-allocation-strategies", Paginated: true,
		Examples: []map[string]any{
			{"kind": "auto", "scope": scopeExample(sdk.RESOURCESCOPEKIND_FABRIC, 1)},
			{"kind": "manual", "scope": scopeExample(sdk.RESOURCESCOPEKIND_FABRIC, 1), "pkey": 100},
		},
	}
	Zone = Family{
		Name: "zone", PathSegment: "zone/zone-allocation-strategies", Paginated: true,
		Examples: []map[string]any{
			{"kind": "auto", "scope": scopeExample(sdk.RESOURCESCOPEKIND_FABRIC, 1)},
			{"kind": "manual", "scope": scopeExample(sdk.RESOURCESCOPEKIND_FABRIC, 1), "zoneName": "zone-a"},
		},
	}
	L3Vlan = Family{
		Name: "l3-vlan", PathSegment: "l3-vlan-allocation-strategies", Paginated: true,
		Examples: []map[string]any{
			{"kind": "auto", "scope": scopeExample(sdk.RESOURCESCOPEKIND_FABRIC, 1), "granularityLevel": string(sdk.VLANALLOCATIONGRANULARITYLEVEL_FABRIC)},
			{"kind": "manual", "scope": scopeExample(sdk.RESOURCESCOPEKIND_FABRIC, 1), "granularityLevel": string(sdk.VLANALLOCATIONGRANULARITYLEVEL_FABRIC), "vlanId": 400},
		},
	}
	L3Vni = Family{
		Name: "l3-vni", PathSegment: "l3-vni-allocation-strategies", Paginated: true,
		Examples: []map[string]any{
			{"kind": "auto", "scope": scopeExample(sdk.RESOURCESCOPEKIND_FABRIC, 1)},
			{"kind": "manual", "scope": scopeExample(sdk.RESOURCESCOPEKIND_FABRIC, 1), "vni": 9100},
		},
	}
	Vrf = Family{
		Name: "vrf", PathSegment: "vrf-allocation-strategies", Paginated: true,
		Examples: []map[string]any{
			{"kind": "manual", "scope": scopeExample(sdk.RESOURCESCOPEKIND_FABRIC, 1), "name": "vrf-tenant-a"},
		},
	}
	Ipv4PointToPoint = Family{
		Name: "ipv4", PathSegment: "ipv4/subnet-allocation-strategies", Paginated: false,
		Examples: []map[string]any{
			{"kind": "auto", "scope": scopeExample(sdk.RESOURCESCOPEKIND_GLOBAL, 0), "subnetPoolIds": []int64{1}, "interfaceABinding": string(sdk.POINTTOPOINTINTERFACEBINDING_AUTO)},
			{"kind": "manual", "scope": scopeExample(sdk.RESOURCESCOPEKIND_GLOBAL, 0), "subnetId": 1, "interfaceABinding": string(sdk.POINTTOPOINTINTERFACEBINDING_A_FIRST)},
			{"kind": "unnumbered", "scope": scopeExample(sdk.RESOURCESCOPEKIND_GLOBAL, 0)},
		},
	}
	Ipv6PointToPoint = Family{
		Name: "ipv6", PathSegment: "ipv6/subnet-allocation-strategies", Paginated: false,
		Examples: []map[string]any{
			{"kind": "auto", "scope": scopeExample(sdk.RESOURCESCOPEKIND_GLOBAL, 0), "subnetPoolIds": []int64{1}, "interfaceABinding": string(sdk.POINTTOPOINTINTERFACEBINDING_AUTO)},
			{"kind": "manual", "scope": scopeExample(sdk.RESOURCESCOPEKIND_GLOBAL, 0), "subnetId": 1, "interfaceABinding": string(sdk.POINTTOPOINTINTERFACEBINDING_A_FIRST)},
			{"kind": "unnumbered", "scope": scopeExample(sdk.RESOURCESCOPEKIND_GLOBAL, 0)},
		},
	}
)

// Parents owning allocation strategies.
var (
	LogicalNetwork = Parent{
		Name:         "logical network",
		ItemPath:     func(id int64) string { return fmt.Sprintf("/api/v2/logical-networks/%d", id) },
		ConfigSuffix: "/config",
		Families:     []Family{Ipv4Subnet, Ipv6Subnet, Vlan, Vni, Pkey, Zone},
	}
	LogicalNetworkProfile = Parent{
		Name:         "logical network profile",
		ItemPath:     func(id int64) string { return fmt.Sprintf("/api/v2/logical-network-profiles/%d", id) },
		ConfigSuffix: "",
		Families:     []Family{Ipv4Subnet, Ipv6Subnet, Vlan, Vni, Pkey, Zone},
	}
	RouteDomain = Parent{
		Name:         "route domain",
		ItemPath:     func(id int64) string { return fmt.Sprintf("/api/v2/route-domains/%d", id) },
		ConfigSuffix: "/config",
		Families:     []Family{L3Vlan, L3Vni, Vrf},
	}
	PointToPointLink = Parent{
		Name:         "point-to-point link",
		ItemPath:     func(id int64) string { return fmt.Sprintf("/api/v2/point-to-point-links/%d", id) },
		ConfigSuffix: "/config",
		Families:     []Family{Ipv4PointToPoint, Ipv6PointToPoint},
	}
)
