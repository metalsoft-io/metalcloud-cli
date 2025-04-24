package logical_network

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/internal/fabric"
	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var logicalNetworkPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"VlanLogicalNetwork": {
			Hidden: true,
			InnerFields: map[string]formatter.RecordFieldConfig{
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
					Title: "Kind",
					Order: 4,
				},
				"FabricId": {
					Title: "Fabric ID",
					Order: 5,
				},
				"InfrastructureId": {
					Title: "Infra ID",
					Order: 6,
				},
			},
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
}

func LogicalNetworkList(ctx context.Context, fabricIdOrLabel string, flags ListFlags) error {
	logger.Get().Info().Msgf("Listing logical networks")

	client := api.GetApiClient(ctx)

	request := client.LogicalNetworkAPI.GetLogicalNetworks(ctx)

	if len(flags.FilterId) > 0 {
		request = request.FilterId(flags.FilterId)
	}
	if len(flags.FilterLabel) > 0 {
		request = request.FilterLabel(flags.FilterLabel)
	}
	if len(flags.FilterFabricId) > 0 {
		request = request.FilterFabricId(flags.FilterFabricId)
	}
	if len(flags.FilterInfrastructureId) > 0 {
		request = request.FilterInfrastructureId(flags.FilterInfrastructureId)
	}
	if len(flags.FilterKind) > 0 {
		request = request.FilterKind(flags.FilterKind)
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

	logicalNetworkList, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(logicalNetworkList, &logicalNetworkPrintConfig)
}

func LogicalNetworkGet(ctx context.Context, logicalNetworkId string) error {
	logger.Get().Info().Msgf("Get logical network '%s' details", logicalNetworkId)

	logicalNetworkIdNumeric, err := getLogicalNetworkId(logicalNetworkId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	logicalNetwork, httpRes, err := client.LogicalNetworkAPI.GetLogicalNetwork(ctx, logicalNetworkIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(logicalNetwork, &logicalNetworkPrintConfig)
}

func LogicalNetworkConfigExample(ctx context.Context, kind string) error {
	logicalNetworkConfiguration := sdk.CreateLogicalNetworkRequest{}

	if kind == string(sdk.LOGICALNETWORKKIND_VLAN) {
		logicalNetworkConfiguration.CreateVlanLogicalNetwork = &sdk.CreateVlanLogicalNetwork{
			Kind:             sdk.LOGICALNETWORKKIND_VLAN,
			Label:            sdk.PtrString("example-logical-network"),
			Name:             sdk.PtrString("Example Logical Network"),
			FabricId:         1,
			InfrastructureId: *sdk.NewNullableInt32(sdk.PtrInt32(1)),
			Vlan: sdk.CreateVlanLogicalNetworkVlanProperties{
				VlanAllocationStrategies: []sdk.CreateVlanAllocationStrategy{
					{
						CreateAutoVlanAllocationStrategy: &sdk.CreateAutoVlanAllocationStrategy{
							Kind: sdk.ALLOCATIONSTRATEGYKIND_AUTO,
							Scope: sdk.CreateResourceScope{
								Kind:       sdk.RESOURCESCOPEKIND_FABRIC,
								ResourceId: 1,
							},
							GranularityLevel: *sdk.NewNullableVlanAllocationGranularityLevel(sdk.VLANALLOCATIONGRANULARITYLEVEL_NETWORK_DEVICE.Ptr()),
						},
					},
				},
			},
			Ipv4: sdk.CreateVlanLogicalNetworkIpv4Properties{
				SubnetAllocationStrategies: []sdk.CreateIpv4SubnetAllocationStrategy{
					{
						CreateAutoIpv4SubnetAllocationStrategy: &sdk.CreateAutoIpv4SubnetAllocationStrategy{
							Kind: sdk.ALLOCATIONSTRATEGYKIND_AUTO,
							Scope: sdk.CreateResourceScope{
								Kind:       sdk.RESOURCESCOPEKIND_FABRIC,
								ResourceId: 1,
							},
							PrefixLength:  24,
							SubnetPoolIds: []int32{2, 3},
						},
					},
				},
			},
			Ipv6: &sdk.CreateVlanLogicalNetworkIpv6Properties{
				SubnetAllocationStrategies: []sdk.CreateIpv6SubnetAllocationStrategy{
					{
						CreateAutoIpv6SubnetAllocationStrategy: &sdk.CreateAutoIpv6SubnetAllocationStrategy{
							Kind: sdk.ALLOCATIONSTRATEGYKIND_AUTO,
							Scope: sdk.CreateResourceScope{
								Kind:       sdk.RESOURCESCOPEKIND_FABRIC,
								ResourceId: 1,
							},
							PrefixLength:  64,
							SubnetPoolIds: []int32{2, 3},
						},
					},
				},
			},
			RouteDomainId: *sdk.NewNullableInt32(sdk.PtrInt32(1)),
			Annotations: &map[string]string{
				"example": "example",
			},
		}
	} else if kind == string(sdk.LOGICALNETWORKKIND_VXLAN) {
		logicalNetworkConfiguration.CreateVxlanLogicalNetwork = &sdk.CreateVxlanLogicalNetwork{
			Kind:             sdk.LOGICALNETWORKKIND_VXLAN,
			Label:            sdk.PtrString("example-logical-network"),
			Name:             sdk.PtrString("Example Logical Network"),
			FabricId:         1,
			InfrastructureId: *sdk.NewNullableInt32(sdk.PtrInt32(1)),
			Vlan: sdk.CreateVxlanLogicalNetworkVlanProperties{
				VlanAllocationStrategies: []sdk.CreateVlanAllocationStrategy{
					{
						CreateAutoVlanAllocationStrategy: &sdk.CreateAutoVlanAllocationStrategy{
							Kind: sdk.ALLOCATIONSTRATEGYKIND_AUTO,
							Scope: sdk.CreateResourceScope{
								Kind:       sdk.RESOURCESCOPEKIND_FABRIC,
								ResourceId: 1,
							},
							GranularityLevel: *sdk.NewNullableVlanAllocationGranularityLevel(sdk.VLANALLOCATIONGRANULARITYLEVEL_NETWORK_DEVICE.Ptr()),
						},
					},
				},
			},
			Vxlan: sdk.CreateVxlanLogicalNetworkVxlanProperties{
				VniAllocationStrategies: []sdk.CreateVniAllocationStrategy{
					{
						CreateAutoVniAllocationStrategy: &sdk.CreateAutoVniAllocationStrategy{
							Kind: sdk.ALLOCATIONSTRATEGYKIND_AUTO,
							Scope: sdk.CreateResourceScope{
								Kind:       sdk.RESOURCESCOPEKIND_FABRIC,
								ResourceId: 1,
							},
						},
					},
				},
			},
			Ipv4: sdk.CreateVxlanLogicalNetworkIpv4Properties{
				SubnetAllocationStrategies: []sdk.CreateIpv4SubnetAllocationStrategy{
					{
						CreateAutoIpv4SubnetAllocationStrategy: &sdk.CreateAutoIpv4SubnetAllocationStrategy{
							Kind: sdk.ALLOCATIONSTRATEGYKIND_AUTO,
							Scope: sdk.CreateResourceScope{
								Kind:       sdk.RESOURCESCOPEKIND_FABRIC,
								ResourceId: 1,
							},
							PrefixLength:  24,
							SubnetPoolIds: []int32{2, 3},
						},
					},
				},
			},
			Ipv6: &sdk.CreateVxlanLogicalNetworkIpv6Properties{
				SubnetAllocationStrategies: []sdk.CreateIpv6SubnetAllocationStrategy{
					{
						CreateAutoIpv6SubnetAllocationStrategy: &sdk.CreateAutoIpv6SubnetAllocationStrategy{
							Kind: sdk.ALLOCATIONSTRATEGYKIND_AUTO,
							Scope: sdk.CreateResourceScope{
								Kind:       sdk.RESOURCESCOPEKIND_FABRIC,
								ResourceId: 1,
							},
							PrefixLength:  64,
							SubnetPoolIds: []int32{2, 3},
						},
					},
				},
			},
			RouteDomainId: *sdk.NewNullableInt32(sdk.PtrInt32(1)),
			Annotations: &map[string]string{
				"example": "example",
			},
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

	var logicalNetworkConfig sdk.CreateLogicalNetworkRequest
	if kind == string(sdk.LOGICALNETWORKKIND_VLAN) {
		vlanLogicalNetwork := sdk.CreateVlanLogicalNetwork{}
		err := json.Unmarshal(config, &vlanLogicalNetwork)
		if err != nil {
			return err
		}

		logicalNetworkConfig.CreateVlanLogicalNetwork = &vlanLogicalNetwork
	} else if kind == string(sdk.LOGICALNETWORKKIND_VXLAN) {
		vxlanLogicalNetwork := sdk.CreateVxlanLogicalNetwork{}
		err := json.Unmarshal(config, &vxlanLogicalNetwork)
		if err != nil {
			return err
		}

		logicalNetworkConfig.CreateVxlanLogicalNetwork = &vxlanLogicalNetwork
	} else {
		err := fmt.Errorf("unsupported logical network kind '%s'", kind)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	logicalNetwork, httpRes, err := client.LogicalNetworkAPI.
		CreateLogicalNetwork(ctx).
		CreateLogicalNetworkRequest(logicalNetworkConfig).
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
	err = json.Unmarshal(config, &logicalNetworkUpdate)
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

	httpRes, err := client.LogicalNetworkAPI.
		DeleteLogicalNetwork(ctx, logicalNetworkIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Logical network '%s' deleted", logicalNetworkId)
	return nil
}

func getLogicalNetworkId(logicalNetworkId string) (float32, error) {
	logicalNetworkIdNumeric, err := strconv.ParseFloat(logicalNetworkId, 32)
	if err != nil {
		err := fmt.Errorf("invalid logical network ID: '%s'", logicalNetworkId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(logicalNetworkIdNumeric), nil
}
