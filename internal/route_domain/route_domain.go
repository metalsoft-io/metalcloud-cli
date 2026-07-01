package route_domain

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

	client := api.GetApiClient(ctx)

	request := client.RouteDomainAPI.GetRouteDomains(ctx).SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &RouteDomainPrintConfig)
}

func RouteDomainGet(ctx context.Context, routeDomainId string) error {
	logger.Get().Info().Msgf("Get route domain %s details", routeDomainId)

	routeDomainIdNumeric, err := getRouteDomainId(routeDomainId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	routeDomain, httpRes, err := client.RouteDomainAPI.GetRouteDomain(ctx, routeDomainIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(routeDomain, &RouteDomainPrintConfig)
}

func RouteDomainCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating route domain")

	var routeDomainConfig sdk.CreateRouteDomain
	if err := utils.UnmarshalContent(config, &routeDomainConfig); err != nil {
		return err
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
				Scope: sdk.CreateResourceScope{Kind: sdk.RESOURCESCOPEKIND_FABRIC, ResourceId: 1},
				Name:  "tenant1",
			}),
		},
		L3VniAllocationStrategies: []sdk.CreateVniAllocationStrategy{
			sdk.CreateAutoVniAllocationStrategyAsCreateVniAllocationStrategy(&sdk.CreateAutoVniAllocationStrategy{
				Kind:             sdk.ALLOCATIONSTRATEGYKIND_AUTO,
				Scope:            sdk.CreateResourceScope{Kind: sdk.RESOURCESCOPEKIND_FABRIC, ResourceId: 1},
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

	client := api.GetApiClient(ctx)

	routeDomain, httpRes, err := client.RouteDomainAPI.GetRouteDomain(ctx, routeDomainIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}

	return routeDomainIdNumeric, strconv.FormatInt(routeDomain.Revision, 10), nil
}
