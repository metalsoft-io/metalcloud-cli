package logical_network_profile

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var logicalNetworkProfilePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"LogicalNetworkProfile": {
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
			},
		},
	},
}

type ListFlags struct {
	FilterId       []string
	FilterLabel    []string
	FilterKind     []string
	FilterName     []string
	FilterFabricId []string
	SortBy         []string
}

func LogicalNetworkProfileList(ctx context.Context, flags ListFlags) error {
	logger.Get().Info().Msg("Listing logical network profiles")

	client := api.GetApiClient(ctx)
	request := client.LogicalNetworkProfileAPI.GetLogicalNetworkProfiles(ctx)

	if len(flags.FilterId) > 0 {
		request = request.FilterId(utils.ProcessFilterStringSlice(flags.FilterId))
	}
	if len(flags.FilterLabel) > 0 {
		request = request.FilterLabel(utils.ProcessFilterStringSlice(flags.FilterLabel))
	}
	if len(flags.FilterKind) > 0 {
		request = request.FilterKind(utils.ProcessFilterStringSlice(flags.FilterKind))
	}
	if len(flags.FilterName) > 0 {
		request = request.FilterName(utils.ProcessFilterStringSlice(flags.FilterName))
	}
	if len(flags.FilterFabricId) > 0 {
		request = request.FilterFabricId(utils.ProcessFilterStringSlice(flags.FilterFabricId))
	}
	if len(flags.SortBy) > 0 {
		request = request.SortBy(flags.SortBy)
	}

	profiles, httpRes, err := request.Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(profiles, &logicalNetworkProfilePrintConfig)
}

func LogicalNetworkProfileGet(ctx context.Context, profileId string) error {
	logger.Get().Info().Msgf("Get logical network profile '%s' details", profileId)

	id, err := getLogicalNetworkProfileId(profileId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)
	profile, httpRes, err := client.LogicalNetworkProfileAPI.GetLogicalNetworkProfile(ctx, id).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(profile, &logicalNetworkProfilePrintConfig)
}

func LogicalNetworkProfileConfigExample(ctx context.Context, kind string) error {
	example := sdk.CreateLogicalNetworkProfileRequest{}

	switch kind {
	case string(sdk.LOGICALNETWORKKIND_VLAN):
		example.CreateVlanLogicalNetworkProfile = &sdk.CreateVlanLogicalNetworkProfile{
			Kind:     sdk.LOGICALNETWORKKIND_VLAN,
			Label:    sdk.PtrString("example-vlan"),
			Name:     sdk.PtrString("example-vlan"),
			FabricId: 1,
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
							SubnetPoolIds: []int32{1},
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
							SubnetPoolIds: []int32{1},
						},
					},
				},
			},
			RouteDomainId: *sdk.NewNullableInt32(sdk.PtrInt32(1)),
		}
	case string(sdk.LOGICALNETWORKKIND_VXLAN):
		example.CreateVxlanLogicalNetworkProfile = &sdk.CreateVxlanLogicalNetworkProfile{
			Kind:     sdk.LOGICALNETWORKKIND_VXLAN,
			Label:    sdk.PtrString("example-vxlan"),
			Name:     sdk.PtrString("example-vxlan"),
			FabricId: 1,
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
							SubnetPoolIds: []int32{1},
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
							SubnetPoolIds: []int32{1},
						},
					},
				},
			},
			RouteDomainId: *sdk.NewNullableInt32(sdk.PtrInt32(1)),
		}
	default:
		return fmt.Errorf("invalid logical network profile kind: '%s'", kind)
	}

	return formatter.PrintResult(example, nil)
}

func LogicalNetworkProfileCreate(ctx context.Context, config []byte, kind string) error {
	logger.Get().Info().Msg("Creating logical network profile")

	var req sdk.CreateLogicalNetworkProfileRequest

	switch kind {
	case string(sdk.LOGICALNETWORKKIND_VLAN):
		var vlanReq sdk.CreateVlanLogicalNetworkProfile
		if err := json.Unmarshal(config, &vlanReq); err != nil {
			return err
		}
		req.CreateVlanLogicalNetworkProfile = &vlanReq
	case string(sdk.LOGICALNETWORKKIND_VXLAN):
		var vxlanReq sdk.CreateVxlanLogicalNetworkProfile
		if err := json.Unmarshal(config, &vxlanReq); err != nil {
			return err
		}
		req.CreateVxlanLogicalNetworkProfile = &vxlanReq
	default:
		return fmt.Errorf("invalid logical network profile kind: '%s'", kind)
	}

	client := api.GetApiClient(ctx)
	profile, httpRes, err := client.LogicalNetworkProfileAPI.
		CreateLogicalNetworkProfile(ctx).
		CreateLogicalNetworkProfileRequest(req).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(profile, &logicalNetworkProfilePrintConfig)
}

func LogicalNetworkProfileUpdate(ctx context.Context, profileId string, config []byte) error {
	logger.Get().Info().Msgf("Updating logical network profile '%s'", profileId)

	id, err := getLogicalNetworkProfileId(profileId)
	if err != nil {
		return err
	}

	var req sdk.UpdateLogicalNetworkProfile
	if err := json.Unmarshal(config, &req); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)
	profile, httpRes, err := client.LogicalNetworkProfileAPI.
		UpdateLogicalNetworkProfile(ctx, id).
		UpdateLogicalNetworkProfile(req).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(profile, &logicalNetworkProfilePrintConfig)
}

func LogicalNetworkProfileDelete(ctx context.Context, profileId string) error {
	logger.Get().Info().Msgf("Deleting logical network profile '%s'", profileId)

	id, err := getLogicalNetworkProfileId(profileId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)
	httpRes, err := client.LogicalNetworkProfileAPI.
		DeleteLogicalNetworkProfile(ctx, id).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Logical network profile '%s' deleted", profileId)
	return nil
}

func getLogicalNetworkProfileId(profileId string) (float32, error) {
	id, err := strconv.ParseFloat(profileId, 32)
	if err != nil {
		err := fmt.Errorf("invalid logical network profile ID: '%s'", profileId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}
	return float32(id), nil
}
