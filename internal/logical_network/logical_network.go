package logical_network

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/internal/fabric"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

type logicalNetworkRaw struct {
	Id               interface{} `json:"id"`
	Label            *string     `json:"label"`
	Name             *string     `json:"name"`
	Kind             *string     `json:"kind"`
	FabricId         interface{} `json:"fabricId"`
	InfrastructureId interface{} `json:"infrastructureId"`
}

var logicalNetworkPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title:       "#",
			Transformer: formatter.FormatIdValue,
			Order:       1,
		},
		"Label": {
			MaxWidth: 30,
			Order:    2,
		},
		"Name": {
			MaxWidth: 30,
			Order:    3,
		},
		"Kind": {
			Title: "Kind",
			Order: 4,
		},
		"FabricId": {
			Title:       "Fabric ID",
			Transformer: formatter.FormatIdValue,
			Order:       5,
		},
		"InfrastructureId": {
			Title:       "Infra ID",
			Transformer: formatter.FormatIdValue,
			Order:       6,
		},
	},
}

type ListFlags struct {
	FilterId               []string
	FilterLabel            []string
	FilterFabricId         []string
	FilterInfrastructureId []string
	FilterKind             []string
	SortBy                 []string
	Page                   int
	Limit                  int
}

func LogicalNetworkList(ctx context.Context, fabricIdOrLabel string, flags ListFlags) error {
	logger.Get().Info().Msgf("Listing logical networks with filters: %+v", flags)

	client := api.GetApiClient(ctx)

	request := client.LogicalNetworkAPI.GetLogicalNetworks(ctx)

	if len(flags.FilterId) > 0 {
		request = request.FilterId(utils.ProcessFilterStringSlice(flags.FilterId))
	}
	if len(flags.FilterLabel) > 0 {
		request = request.FilterLabel(utils.ProcessFilterStringSlice(flags.FilterLabel))
	}
	if len(flags.FilterFabricId) > 0 {
		request = request.FilterFabricId(utils.ProcessFilterStringSlice(flags.FilterFabricId))
	}
	if len(flags.FilterInfrastructureId) > 0 {
		if flags.FilterInfrastructureId[0] == "null" {
			flags.FilterInfrastructureId[0] = "$null"
		}
		request = request.FilterInfrastructureId(utils.ProcessFilterStringSlice(flags.FilterInfrastructureId))
	}
	if len(flags.FilterKind) > 0 {
		request = request.FilterKind(utils.ProcessFilterStringSlice(flags.FilterKind))
	}
	if len(flags.SortBy) > 0 {
		request = request.SortBy(flags.SortBy)
	}

	if fabricIdOrLabel != "" {
		fabric, err := fabric.GetFabricByIdOrLabel(ctx, fabricIdOrLabel)
		if err != nil {
			return err
		}

		request = request.FilterFabricId([]string{fabric.Id})
	}

	printRaw := func(rawItems []json.RawMessage, meta sdk.PaginatedResponseMeta) error {
		records, err := utils.UnmarshalRawItems[logicalNetworkRaw](rawItems)
		if err != nil {
			return fmt.Errorf("failed to parse logical networks: %w", err)
		}
		return utils.PrintAllRaw(rawItems, records, meta, len(records), &logicalNetworkPrintConfig)
	}

	switch {
	case flags.Page > 0:
		rawItems, meta, err := utils.FetchPageWindowRaw(func(p, l float32) (*http.Response, error) {
			_, httpRes, _ := request.Page(p).Limit(l).Execute()
			return httpRes, nil
		}, flags.Page, flags.Limit)
		if err != nil {
			return err
		}
		return printRaw(rawItems, meta)
	case flags.Limit > 0:
		rawItems, meta, err := utils.FetchUpToRaw(func(p, l float32) (*http.Response, error) {
			_, httpRes, _ := request.Page(p).Limit(l).Execute()
			return httpRes, nil
		}, flags.Limit)
		if err != nil {
			return err
		}
		return printRaw(rawItems, meta)
	default:
		rawItems, meta, err := utils.FetchAllPagesRaw(func(page float32) (*http.Response, error) {
			_, httpRes, _ := request.Page(page).Limit(100).Execute()
			return httpRes, nil
		})
		if err != nil {
			return err
		}
		return printRaw(rawItems, meta)
	}
}

func LogicalNetworkGet(ctx context.Context, logicalNetworkId string) error {
	logger.Get().Info().Msgf("Get logical network '%s' details", logicalNetworkId)

	logicalNetworkIdNumeric, err := getLogicalNetworkId(logicalNetworkId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// Raw-body parse: see logicalNetworkRaw — the nested allocation-strategy
	// oneOf mismatch breaks typed decoding.
	_, httpRes, sdkErr := client.LogicalNetworkAPI.GetLogicalNetwork(ctx, logicalNetworkIdNumeric).Execute()
	if httpRes != nil && httpRes.StatusCode >= 400 {
		return response_inspector.InspectResponse(httpRes, sdkErr)
	}
	if httpRes == nil {
		return sdkErr
	}

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var logicalNetwork logicalNetworkRaw
	if err := json.Unmarshal(body, &logicalNetwork); err != nil {
		return fmt.Errorf("failed to parse logical network: %w", err)
	}

	return formatter.PrintResult(logicalNetwork, &logicalNetworkPrintConfig)
}

func LogicalNetworkConfigExample(ctx context.Context, kind string) error {
	logicalNetworkConfiguration := sdk.CreateLogicalNetwork{}

	logicalNetworkConfiguration.Label = sdk.PtrString("example-logical-network")
	logicalNetworkConfiguration.Name = sdk.PtrString("Example Logical Network")
	logicalNetworkConfiguration.FabricId = 1
	logicalNetworkConfiguration.InfrastructureId = *sdk.NewNullableInt64(sdk.PtrInt64(1))

	if kind == string(sdk.LOGICALNETWORKKIND_VLAN) {
		logicalNetworkConfiguration.Kind = sdk.LOGICALNETWORKKIND_VLAN
		logicalNetworkConfiguration.Vlan = &sdk.CreateLogicalNetworkVlanProperties{
			VlanAllocationStrategies: []sdk.CreateVlanAllocationStrategy{
				{
					CreateAutoVlanAllocationStrategy: &sdk.CreateAutoVlanAllocationStrategy{
						Kind: "auto",
						Scope: sdk.CreateResourceScope{
							Kind:       sdk.RESOURCESCOPEKIND_FABRIC,
							ResourceId: *sdk.NewNullableInt64(sdk.PtrInt64(1)),
						},
						GranularityLevel: *sdk.NewNullableVlanAllocationGranularityLevel(sdk.VLANALLOCATIONGRANULARITYLEVEL_NETWORK_DEVICE.Ptr()),
					},
				},
			},
		}
		logicalNetworkConfiguration.Ipv4 = &sdk.CreateLogicalNetworkIpv4Properties{
			SubnetAllocationStrategies: []sdk.CreateIpv4SubnetAllocationStrategy{
				{
					CreateAutoIpv4SubnetAllocationStrategy: &sdk.CreateAutoIpv4SubnetAllocationStrategy{
						Kind: "auto",
						Scope: sdk.CreateResourceScope{
							Kind:       sdk.RESOURCESCOPEKIND_FABRIC,
							ResourceId: *sdk.NewNullableInt64(sdk.PtrInt64(1)),
						},
						PrefixLength:  24,
						SubnetPoolIds: []int64{2, 3},
					},
				},
			},
		}
		logicalNetworkConfiguration.Ipv6 = &sdk.CreateLogicalNetworkIpv6Properties{
			SubnetAllocationStrategies: []sdk.CreateIpv6SubnetAllocationStrategy{
				{
					CreateAutoIpv6SubnetAllocationStrategy: &sdk.CreateAutoIpv6SubnetAllocationStrategy{
						Kind: "auto",
						Scope: sdk.CreateResourceScope{
							Kind:       sdk.RESOURCESCOPEKIND_FABRIC,
							ResourceId: *sdk.NewNullableInt64(sdk.PtrInt64(1)),
						},
						PrefixLength:  64,
						SubnetPoolIds: []int64{2, 3},
					},
				},
			},
		}
		logicalNetworkConfiguration.RouteDomainId = *sdk.NewNullableInt64(sdk.PtrInt64(1))
		logicalNetworkConfiguration.Annotations = &map[string]string{
			"example": "example",
		}
	} else if kind == string(sdk.LOGICALNETWORKKIND_VXLAN) {
		logicalNetworkConfiguration.Kind = sdk.LOGICALNETWORKKIND_VXLAN
		logicalNetworkConfiguration.Vlan = &sdk.CreateLogicalNetworkVlanProperties{
			VlanAllocationStrategies: []sdk.CreateVlanAllocationStrategy{
				{
					CreateAutoVlanAllocationStrategy: &sdk.CreateAutoVlanAllocationStrategy{
						Kind: "auto",
						Scope: sdk.CreateResourceScope{
							Kind:       sdk.RESOURCESCOPEKIND_FABRIC,
							ResourceId: *sdk.NewNullableInt64(sdk.PtrInt64(1)),
						},
						GranularityLevel: *sdk.NewNullableVlanAllocationGranularityLevel(sdk.VLANALLOCATIONGRANULARITYLEVEL_NETWORK_DEVICE.Ptr()),
					},
				},
			},
		}
		logicalNetworkConfiguration.Vxlan = &sdk.CreateLogicalNetworkVxlanProperties{
			VniAllocationStrategies: []sdk.CreateVniAllocationStrategy{
				{
					CreateAutoVniAllocationStrategy: &sdk.CreateAutoVniAllocationStrategy{
						Kind: "auto",
						Scope: sdk.CreateResourceScope{
							Kind:       sdk.RESOURCESCOPEKIND_FABRIC,
							ResourceId: *sdk.NewNullableInt64(sdk.PtrInt64(1)),
						},
					},
				},
			},
		}
		logicalNetworkConfiguration.Ipv4 = &sdk.CreateLogicalNetworkIpv4Properties{
			SubnetAllocationStrategies: []sdk.CreateIpv4SubnetAllocationStrategy{
				{
					CreateAutoIpv4SubnetAllocationStrategy: &sdk.CreateAutoIpv4SubnetAllocationStrategy{
						Kind: "auto",
						Scope: sdk.CreateResourceScope{
							Kind:       sdk.RESOURCESCOPEKIND_FABRIC,
							ResourceId: *sdk.NewNullableInt64(sdk.PtrInt64(1)),
						},
						PrefixLength:  24,
						SubnetPoolIds: []int64{2, 3},
					},
				},
			},
		}
		logicalNetworkConfiguration.Ipv6 = &sdk.CreateLogicalNetworkIpv6Properties{
			SubnetAllocationStrategies: []sdk.CreateIpv6SubnetAllocationStrategy{
				{
					CreateAutoIpv6SubnetAllocationStrategy: &sdk.CreateAutoIpv6SubnetAllocationStrategy{
						Kind: "auto",
						Scope: sdk.CreateResourceScope{
							Kind:       sdk.RESOURCESCOPEKIND_FABRIC,
							ResourceId: *sdk.NewNullableInt64(sdk.PtrInt64(1)),
						},
						PrefixLength:  64,
						SubnetPoolIds: []int64{2, 3},
					},
				},
			},
		}
		logicalNetworkConfiguration.RouteDomainId = *sdk.NewNullableInt64(sdk.PtrInt64(1))
		logicalNetworkConfiguration.Annotations = &map[string]string{
			"example": "example",
		}
	} else {
		err := fmt.Errorf("unsupported logical network kind '%s'", kind)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	return formatter.PrintResult(logicalNetworkConfiguration, nil)
}

func LogicalNetworkCreate(ctx context.Context, kind string, config []byte) error {
	logger.Get().Info().Msgf("Creating logical network")

	var logicalNetworkConfig sdk.CreateLogicalNetwork
	err := utils.UnmarshalContent(config, &logicalNetworkConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	logicalNetwork, httpRes, err := client.LogicalNetworkAPI.
		CreateLogicalNetwork(ctx).
		CreateLogicalNetwork(logicalNetworkConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(logicalNetwork, &logicalNetworkPrintConfig)
}

func LogicalNetworkUpdate(ctx context.Context, logicalNetworkId string, config []byte) error {
	logger.Get().Info().Msgf("Updating logical network '%s'", logicalNetworkId)

	logicalNetworkIdNumeric, err := getLogicalNetworkId(logicalNetworkId)
	if err != nil {
		return err
	}

	var logicalNetworkUpdate sdk.UpdateLogicalNetwork
	err = utils.UnmarshalContent(config, &logicalNetworkUpdate)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	logicalNetwork, httpRes, err := client.LogicalNetworkAPI.
		UpdateLogicalNetwork(ctx, logicalNetworkIdNumeric).
		UpdateLogicalNetwork(logicalNetworkUpdate).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(logicalNetwork, &logicalNetworkPrintConfig)
}

func LogicalNetworkDelete(ctx context.Context, logicalNetworkId string) error {
	logger.Get().Info().Msgf("Deleting logical network '%s'", logicalNetworkId)

	logicalNetworkIdNumeric, err := getLogicalNetworkId(logicalNetworkId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	ln, httpRes, err := client.LogicalNetworkAPI.GetLogicalNetwork(ctx, logicalNetworkIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}
	revision := strconv.Itoa(int(ln.Revision))

	httpRes, err = client.LogicalNetworkAPI.
		DeleteLogicalNetwork(ctx, logicalNetworkIdNumeric).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Logical network '%s' deleted", logicalNetworkId)
	return nil
}

func getLogicalNetworkId(logicalNetworkId string) (int64, error) {
	logicalNetworkIdNumeric, err := strconv.ParseInt(logicalNetworkId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid logical network ID: '%s'", logicalNetworkId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return logicalNetworkIdNumeric, nil
}

// ---------------------------------------------------------------------------
// Logical network config and attached resources
// ---------------------------------------------------------------------------

// logicalNetworkConfigPrintConfig renders the scalar fields of a logical
// network config. The nested allocation-strategy collections are managed by
// the 'logical-network allocation-strategy' commands and are only rendered in
// the json/yaml output.
var logicalNetworkConfigPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Kind": {
			Order: 2,
		},
		"DeployType": {
			Title: "Deploy Type",
			Order: 3,
		},
		"DeployStatus": {
			Title:       "Deploy Status",
			Transformer: formatter.FormatStatusValue,
			Order:       4,
		},
		"Mtu": {
			Title: "MTU",
			Order: 5,
		},
		"Revision": {
			Order: 6,
		},
		"UpdatedAt": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
	},
}

var externalConnectionPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Label": {
			MaxWidth: 30,
			Order:    2,
		},
		"Name": {
			MaxWidth: 30,
			Order:    3,
		},
		"FabricId": {
			Title: "Fabric ID",
			Order: 4,
		},
		"CreatedAt": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       5,
		},
	},
}

var externalConnectionLogicalNetworkPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"ExternalConnectionId": {
			Title: "External Connection",
			Order: 2,
		},
		"LogicalNetworkId": {
			Title: "Logical Network",
			Order: 3,
		},
		"Status": {
			Transformer: formatter.FormatStatusValue,
			Order:       4,
		},
		"CreatedAt": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       5,
		},
	},
}

var logicalNetworkInterconnectPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Label": {
			MaxWidth: 30,
			Order:    2,
		},
		"Name": {
			MaxWidth: 30,
			Order:    3,
		},
		"Kind": {
			Order: 4,
		},
		"FabricInterconnectId": {
			Title: "Fabric Interconnect",
			Order: 5,
		},
		"TransportId": {
			Title: "Transport",
			Order: 6,
		},
		"Status": {
			Transformer: formatter.FormatStatusValue,
			Order:       7,
		},
		"CreatedAt": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       8,
		},
	},
}

func LogicalNetworkConfigGet(ctx context.Context, logicalNetworkId string) error {
	logger.Get().Info().Msgf("Get logical network '%s' config", logicalNetworkId)

	config, _, err := getLogicalNetworkConfig(ctx, logicalNetworkId)
	if err != nil {
		return err
	}

	return formatter.PrintResult(config, &logicalNetworkConfigPrintConfig)
}

func LogicalNetworkConfigUpdate(ctx context.Context, logicalNetworkId string, config []byte) error {
	logger.Get().Info().Msgf("Updating logical network '%s' config", logicalNetworkId)

	var settings sdk.UpdateLogicalNetworkConfigGlobalSettings
	if err := utils.UnmarshalContent(config, &settings); err != nil {
		return err
	}

	logicalNetworkIdNumeric, revision, err := getLogicalNetworkConfigRevision(ctx, logicalNetworkId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	updated, httpRes, err := client.LogicalNetworkAPI.
		UpdateLogicalNetworkConfig(ctx, logicalNetworkIdNumeric).
		UpdateLogicalNetworkConfigGlobalSettings(settings).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(updated, &logicalNetworkConfigPrintConfig)
}

// LogicalNetworkApplyProfiles applies a logical network profile to the logical
// network's config, replacing the config's allocation strategies with the
// profile's.
func LogicalNetworkApplyProfiles(ctx context.Context, logicalNetworkId string, profileId string) error {
	logger.Get().Info().Msgf("Applying profile '%s' to logical network '%s' config", profileId, logicalNetworkId)

	profileIdNumeric, err := utils.GetInt64FromString(profileId)
	if err != nil {
		err = fmt.Errorf("invalid logical network profile ID: '%s'", profileId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	logicalNetworkIdNumeric, revision, err := getLogicalNetworkConfigRevision(ctx, logicalNetworkId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// The body is always sent: an optional body that is never set serializes
	// as a literal "null", which the API rejects with 400.
	body := sdk.ApplyProfilesToLogicalNetworkConfig{}
	body.SetLogicalNetworkProfileId(profileIdNumeric)

	config, httpRes, err := client.LogicalNetworkAPI.
		ApplyProfilesToLogicalNetworkConfig(ctx, logicalNetworkIdNumeric).
		ApplyProfilesToLogicalNetworkConfig(body).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(config, &logicalNetworkConfigPrintConfig)
}

func LogicalNetworkCreateFromProfile(ctx context.Context, create sdk.CreateLogicalNetworkFromProfile) error {
	logger.Get().Info().Msgf("Creating logical network from profile %d", create.LogicalNetworkProfileId)

	client := api.GetApiClient(ctx)

	logicalNetwork, httpRes, err := client.LogicalNetworkAPI.
		CreateLogicalNetworkFromProfile(ctx).
		CreateLogicalNetworkFromProfile(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(logicalNetwork, &logicalNetworkPrintConfig)
}

func LogicalNetworkExternalConnections(ctx context.Context, logicalNetworkId string) error {
	logger.Get().Info().Msgf("Listing external connections attached to logical network '%s'", logicalNetworkId)

	logicalNetworkIdNumeric, err := getLogicalNetworkId(logicalNetworkId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.LogicalNetworkAPI.
		GetLogicalNetworkAttachedExternalConnections(ctx, logicalNetworkIdNumeric).
		SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &externalConnectionPrintConfig)
}

func LogicalNetworkExternalConnectionLogicalNetworks(ctx context.Context, logicalNetworkId string) error {
	logger.Get().Info().Msgf("Listing external connection logical networks of logical network '%s'", logicalNetworkId)

	logicalNetworkIdNumeric, err := getLogicalNetworkId(logicalNetworkId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.LogicalNetworkAPI.
		GetLogicalNetworkAttachedExternalConnectionLogicalNetworks(ctx, logicalNetworkIdNumeric).
		SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &externalConnectionLogicalNetworkPrintConfig)
}

func LogicalNetworkInterconnects(ctx context.Context, logicalNetworkId string) error {
	logger.Get().Info().Msgf("Listing logical network interconnects of logical network '%s'", logicalNetworkId)

	logicalNetworkIdNumeric, err := getLogicalNetworkId(logicalNetworkId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.LogicalNetworkAPI.
		GetLogicalNetworkAttachedLogicalNetworkInterconnects(ctx, logicalNetworkIdNumeric).
		SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &logicalNetworkInterconnectPrintConfig)
}

func LogicalNetworkDetachExternalConnection(ctx context.Context, logicalNetworkId string, externalConnectionId string) error {
	logger.Get().Info().Msgf("Detaching external connection '%s' from logical network '%s'", externalConnectionId, logicalNetworkId)

	logicalNetworkIdNumeric, err := getLogicalNetworkId(logicalNetworkId)
	if err != nil {
		return err
	}

	externalConnectionIdNumeric, err := utils.GetInt64FromString(externalConnectionId)
	if err != nil {
		err = fmt.Errorf("invalid external connection ID: '%s'", externalConnectionId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.LogicalNetworkAPI.
		DetachExternalConnectionLogicalNetwork(ctx, logicalNetworkIdNumeric, externalConnectionIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("External connection '%s' detached from logical network '%s'", externalConnectionId, logicalNetworkId)
	return nil
}

// getLogicalNetworkConfig fetches a logical network's config object and its
// numeric id. The config carries its own revision, which differs from the
// logical network's and is the entity tag the config endpoints expect in
// If-Match (without it they answer 428, with the entity's revision 409).
func getLogicalNetworkConfig(ctx context.Context, logicalNetworkId string) (*sdk.LogicalNetworkConfig, int64, error) {
	logicalNetworkIdNumeric, err := getLogicalNetworkId(logicalNetworkId)
	if err != nil {
		return nil, 0, err
	}

	client := api.GetApiClient(ctx)

	config, httpRes, err := client.LogicalNetworkAPI.GetLogicalNetworkConfig(ctx, logicalNetworkIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, 0, err
	}

	return config, logicalNetworkIdNumeric, nil
}

func getLogicalNetworkConfigRevision(ctx context.Context, logicalNetworkId string) (int64, string, error) {
	config, logicalNetworkIdNumeric, err := getLogicalNetworkConfig(ctx, logicalNetworkId)
	if err != nil {
		return 0, "", err
	}

	return logicalNetworkIdNumeric, strconv.FormatInt(config.Revision, 10), nil
}
