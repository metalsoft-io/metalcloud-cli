package route_domain

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// routeDomainDisplay is a lenient view of a route domain for output. The SDK's
// typed decode of GetRouteDomains/GetRouteDomain fails ("failed to unmarshal
// VlanAllocationStrategy as AutoVlanAllocationStrategy: no value given for
// required property granularityLevel") because the nested VLAN/VNI allocation
// strategy oneOf treats granularityLevel as required. We only display these
// fields, so parse the raw body directly and bypass the broken decode.
type routeDomainDisplay struct {
	Id       int64  `json:"id"`
	Label    string `json:"label"`
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	Revision int64  `json:"revision"`
}

var RouteDomainPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Label": {
			Title: "Label",
			Order: 2,
		},
		"Name": {
			Title: "Name",
			Order: 3,
		},
		"Kind": {
			Title: "Kind",
			Order: 4,
		},
		"Revision": {
			Title: "Revision",
			Order: 5,
		},
	},
}

func RouteDomainList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all route domains")

	rawItems, meta, err := utils.FetchAllPagesRaw(func(page float32) (*http.Response, error) {
		return api.DoJSONRequest(ctx, http.MethodGet,
			fmt.Sprintf("/api/v2/route-domains?page=%.0f&limit=100&sortBy=id:ASC", page), nil)
	})
	if err != nil {
		return err
	}

	records, err := utils.UnmarshalRawItems[routeDomainDisplay](rawItems)
	if err != nil {
		return fmt.Errorf("failed to parse route domains: %w", err)
	}

	return utils.PrintAllRaw(rawItems, records, meta, len(records), &RouteDomainPrintConfig)
}

func RouteDomainGet(ctx context.Context, routeDomainId string) error {
	logger.Get().Info().Msgf("Get route domain %s details", routeDomainId)

	routeDomainIdNumeric, err := getRouteDomainId(routeDomainId)
	if err != nil {
		return err
	}

	httpRes, err := api.DoJSONRequest(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/route-domains/%d", routeDomainIdNumeric), nil)
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}
	defer httpRes.Body.Close()
	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return fmt.Errorf("reading route domain: %w", err)
	}
	var routeDomain routeDomainDisplay
	if err := json.Unmarshal(body, &routeDomain); err != nil {
		return fmt.Errorf("parsing route domain: %w", err)
	}

	return formatter.PrintResult(routeDomain, &RouteDomainPrintConfig)
}

func RouteDomainCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating route domain")

	var routeDomainConfig sdk.CreateRouteDomain
	if err := utils.UnmarshalContent(config, &routeDomainConfig); err != nil {
		return err
	}

	// The API validates the three allocation-strategy fields as arrays and
	// rejects the request ("... must be an array") when they are absent: a nil
	// vrfAllocationStrategies slice serializes as null (it is a required field),
	// and the nil l3V{lan,ni}AllocationStrategies slices are omitted entirely.
	// Normalize nil slices to empty ones so they always marshal as [].
	if routeDomainConfig.VrfAllocationStrategies == nil {
		routeDomainConfig.VrfAllocationStrategies = []sdk.CreateVrfAllocationStrategy{}
	}
	if routeDomainConfig.L3VlanAllocationStrategies == nil {
		routeDomainConfig.L3VlanAllocationStrategies = []sdk.CreateVlanAllocationStrategy{}
	}
	if routeDomainConfig.L3VniAllocationStrategies == nil {
		routeDomainConfig.L3VniAllocationStrategies = []sdk.CreateVniAllocationStrategy{}
	}

	client := api.GetApiClient(ctx)

	routeDomainInfo, httpRes, err := client.RouteDomainAPI.CreateRouteDomain(ctx).CreateRouteDomain(routeDomainConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(routeDomainInfo, &RouteDomainPrintConfig)
}

func RouteDomainUpdate(ctx context.Context, routeDomainId string, config []byte) error {
	logger.Get().Info().Msgf("Updating route domain %s", routeDomainId)

	routeDomainIdNumeric, revision, err := getRouteDomainIdAndRevision(ctx, routeDomainId)
	if err != nil {
		return err
	}

	var routeDomainConfig sdk.UpdateRouteDomain
	if err := utils.UnmarshalContent(config, &routeDomainConfig); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	routeDomainInfo, httpRes, err := client.RouteDomainAPI.
		UpdateRouteDomain(ctx, routeDomainIdNumeric).
		UpdateRouteDomain(routeDomainConfig).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(routeDomainInfo, &RouteDomainPrintConfig)
}

func RouteDomainDelete(ctx context.Context, routeDomainId string) error {
	logger.Get().Info().Msgf("Deleting route domain %s", routeDomainId)

	routeDomainIdNumeric, revision, err := getRouteDomainIdAndRevision(ctx, routeDomainId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.RouteDomainAPI.
		DeleteRouteDomain(ctx, routeDomainIdNumeric).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Route domain %s deleted", routeDomainId)
	return nil
}

// RouteDomainConfigExample prints an example EVPN-L3VPN route domain create body
// (the tenant VRF for an l3evpn fabric): a manual VRF allocation strategy plus
// an auto L3VNI allocation strategy, both fabric-scoped.
func RouteDomainConfigExample(ctx context.Context) error {
	example := sdk.CreateRouteDomain{
		Label: sdk.PtrString("tenant1"),
		Name:  sdk.PtrString("tenant1"),
		Kind:  sdk.ROUTEDOMAINKIND_EVPN_L3VPN,
		VrfAllocationStrategies: []sdk.CreateVrfAllocationStrategy{
			sdk.CreateManualVrfAllocationStrategyAsCreateVrfAllocationStrategy(&sdk.CreateManualVrfAllocationStrategy{
				Kind:  sdk.ALLOCATIONSTRATEGYKIND_MANUAL,
				Scope: sdk.CreateResourceScope{Kind: sdk.RESOURCESCOPEKIND_FABRIC, ResourceId: *sdk.NewNullableInt64(sdk.PtrInt64(1))},
				Name:  "tenant1",
			}),
		},
		L3VniAllocationStrategies: []sdk.CreateVniAllocationStrategy{
			sdk.CreateAutoVniAllocationStrategyAsCreateVniAllocationStrategy(&sdk.CreateAutoVniAllocationStrategy{
				Kind:             sdk.ALLOCATIONSTRATEGYKIND_AUTO,
				Scope:            sdk.CreateResourceScope{Kind: sdk.RESOURCESCOPEKIND_FABRIC, ResourceId: *sdk.NewNullableInt64(sdk.PtrInt64(1))},
				GranularityLevel: *sdk.NewNullableVniAllocationGranularityLevel(sdk.VNIALLOCATIONGRANULARITYLEVEL_FABRIC.Ptr()),
			}),
		},
		AutoRouteTarget:        sdk.PtrBool(true),
		AutoRouteDistinguisher: sdk.PtrBool(true),
	}

	return formatter.PrintResult(example, nil)
}

func getRouteDomainId(routeDomainId string) (int64, error) {
	routeDomainIdNumeric, err := strconv.ParseInt(routeDomainId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid route domain ID: '%s'", routeDomainId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return routeDomainIdNumeric, nil
}

func getRouteDomainIdAndRevision(ctx context.Context, routeDomainId string) (int64, string, error) {
	routeDomainIdNumeric, err := getRouteDomainId(routeDomainId)
	if err != nil {
		return 0, "", err
	}

	// Raw GET: the typed decode trips the same VlanAllocationStrategy oneOf bug.
	httpRes, err := api.DoJSONRequest(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/route-domains/%d", routeDomainIdNumeric), nil)
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}
	defer httpRes.Body.Close()
	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return 0, "", fmt.Errorf("reading route domain: %w", err)
	}
	var routeDomain routeDomainDisplay
	if err := json.Unmarshal(body, &routeDomain); err != nil {
		return 0, "", fmt.Errorf("parsing route domain: %w", err)
	}

	return routeDomainIdNumeric, strconv.FormatInt(routeDomain.Revision, 10), nil
}
