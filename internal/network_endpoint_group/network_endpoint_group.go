package network_endpoint_group

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

var networkEndpointGroupPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Name": {
			Order:    2,
			MaxWidth: 40,
		},
		"SiteId": {
			Title: "Site",
			Order: 3,
		},
		"Revision": {
			Order: 4,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       5,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
	},
}

var networkEndpointGroupLogicalNetworkPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"LogicalNetworkId": {
			Title: "Logical Network",
			Order: 1,
		},
		"NetworkEndpointGroupId": {
			Title: "Endpoint Group",
			Order: 2,
		},
		"AccessMode": {
			Title: "Access Mode",
			Order: 3,
		},
		"Tagged": {
			Order: 4,
		},
		"Mtu": {
			Title: "MTU",
			Order: 5,
		},
		"ProvidesDefaultRoute": {
			Title: "Default Route",
			Order: 6,
		},
		"DisableAutoIpAllocation": {
			Title: "No Auto IP",
			Order: 7,
		},
	},
}

// NetworkEndpointGroupListFilters carries the optional list filters.
type NetworkEndpointGroupListFilters struct {
	Id     []string
	Name   []string
	SiteId []string
}

func NetworkEndpointGroupList(ctx context.Context, filters NetworkEndpointGroupListFilters) error {
	logger.Get().Info().Msgf("Listing network endpoint groups")

	client := api.GetApiClient(ctx)

	request := client.NetworkEndpointGroupAPI.GetNetworkEndpointGroups(ctx).SortBy([]string{"id:ASC"})
	if len(filters.Id) > 0 {
		request = request.FilterId(utils.ProcessFilterStringSlice(filters.Id))
	}
	if len(filters.Name) > 0 {
		request = request.FilterName(utils.ProcessFilterStringSlice(filters.Name))
	}
	if len(filters.SiteId) > 0 {
		request = request.FilterSiteId(utils.ProcessFilterStringSlice(filters.SiteId))
	}

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &networkEndpointGroupPrintConfig)
}

func NetworkEndpointGroupGet(ctx context.Context, groupIdOrName string) error {
	logger.Get().Info().Msgf("Get network endpoint group '%s'", groupIdOrName)

	group, err := GetNetworkEndpointGroupByIdOrLabel(ctx, groupIdOrName)
	if err != nil {
		return err
	}

	return formatter.PrintResult(group, &networkEndpointGroupPrintConfig)
}

func NetworkEndpointGroupConfigExample(ctx context.Context) error {
	example := sdk.CreateNetworkEndpointGroup{
		Name:   "dc1-endpoint-group",
		SiteId: sdk.PtrInt64(1),
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(example, nil)
	}
	return formatter.PrintYamlResult(example)
}

func NetworkEndpointGroupCreate(ctx context.Context, create sdk.CreateNetworkEndpointGroup) error {
	logger.Get().Info().Msgf("Creating network endpoint group '%s'", create.Name)

	client := api.GetApiClient(ctx)

	group, httpRes, err := client.NetworkEndpointGroupAPI.
		CreateNetworkEndpointGroup(ctx).
		CreateNetworkEndpointGroup(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(group, &networkEndpointGroupPrintConfig)
}

func NetworkEndpointGroupUpdate(ctx context.Context, groupIdOrName string, config []byte) error {
	logger.Get().Info().Msgf("Updating network endpoint group '%s'", groupIdOrName)

	var update sdk.UpdateNetworkEndpointGroup
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	groupId, revision, err := resolveGroupIdAndRevision(ctx, groupIdOrName)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	group, httpRes, err := client.NetworkEndpointGroupAPI.
		UpdateNetworkEndpointGroup(ctx, groupId).
		UpdateNetworkEndpointGroup(update).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(group, &networkEndpointGroupPrintConfig)
}

func NetworkEndpointGroupDelete(ctx context.Context, groupIdOrName string) error {
	logger.Get().Info().Msgf("Deleting network endpoint group '%s'", groupIdOrName)

	groupId, err := ResolveNetworkEndpointGroupNumericId(ctx, groupIdOrName)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkEndpointGroupAPI.DeleteNetworkEndpointGroup(ctx, groupId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Network endpoint group '%s' deleted", groupIdOrName)
	return nil
}

// NetworkEndpointGroupLogicalNetworkList lists the logical networks of a
// network endpoint group. The endpoint is not paginated: it returns the whole
// list in a single {"data": [...]} envelope.
func NetworkEndpointGroupLogicalNetworkList(ctx context.Context, groupIdOrName string) error {
	logger.Get().Info().Msgf("Listing logical networks of network endpoint group '%s'", groupIdOrName)

	groupId, err := ResolveNetworkEndpointGroupNumericId(ctx, groupIdOrName)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	logicalNetworks, httpRes, err := client.NetworkEndpointGroupAPI.
		GetNetworkEndpointGroupLogicalNetworks(ctx, groupId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(logicalNetworks, &networkEndpointGroupLogicalNetworkPrintConfig)
}

func NetworkEndpointGroupLogicalNetworkGet(ctx context.Context, groupIdOrName string, logicalNetworkId string) error {
	logger.Get().Info().Msgf("Get logical network '%s' of network endpoint group '%s'", logicalNetworkId, groupIdOrName)

	groupId, err := ResolveNetworkEndpointGroupNumericId(ctx, groupIdOrName)
	if err != nil {
		return err
	}

	logicalNetworkIdNumeric, err := parseId(logicalNetworkId, "logical network")
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	logicalNetwork, httpRes, err := client.NetworkEndpointGroupAPI.
		GetNetworkEndpointGroupLogicalNetwork(ctx, groupId, logicalNetworkIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(logicalNetwork, &networkEndpointGroupLogicalNetworkPrintConfig)
}

// NetworkEndpointGroupLogicalNetworkAdd attaches a logical network to a
// network endpoint group. The endpoint returns no body.
func NetworkEndpointGroupLogicalNetworkAdd(ctx context.Context, groupIdOrName string, create sdk.CreateNetworkEndpointGroupLogicalNetwork) error {
	logger.Get().Info().Msgf("Adding logical network '%s' to network endpoint group '%s'", create.LogicalNetworkId, groupIdOrName)

	groupId, err := ResolveNetworkEndpointGroupNumericId(ctx, groupIdOrName)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkEndpointGroupAPI.
		AddLogicalNetworksToNetworkEndpointGroup(ctx, groupId).
		CreateNetworkEndpointGroupLogicalNetwork(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Logical network '%s' added to network endpoint group '%s'", create.LogicalNetworkId, groupIdOrName)
	return nil
}

// NetworkEndpointGroupLogicalNetworkUpdate updates the settings of one logical
// network of a network endpoint group. The attachment carries no revision of
// its own, so the group's revision is used as the If-Match entity tag.
func NetworkEndpointGroupLogicalNetworkUpdate(ctx context.Context, groupIdOrName string, logicalNetworkId string, update sdk.UpdateNetworkEndpointGroupLogicalNetwork) error {
	logger.Get().Info().Msgf("Updating logical network '%s' of network endpoint group '%s'", logicalNetworkId, groupIdOrName)

	groupId, revision, err := resolveGroupIdAndRevision(ctx, groupIdOrName)
	if err != nil {
		return err
	}

	logicalNetworkIdNumeric, err := parseId(logicalNetworkId, "logical network")
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	logicalNetwork, httpRes, err := client.NetworkEndpointGroupAPI.
		UpdateNetworkEndpointGroupLogicalNetwork(ctx, groupId, logicalNetworkIdNumeric).
		UpdateNetworkEndpointGroupLogicalNetwork(update).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(logicalNetwork, &networkEndpointGroupLogicalNetworkPrintConfig)
}

func NetworkEndpointGroupLogicalNetworkRemove(ctx context.Context, groupIdOrName string, logicalNetworkId string) error {
	logger.Get().Info().Msgf("Removing logical network '%s' from network endpoint group '%s'", logicalNetworkId, groupIdOrName)

	groupId, err := ResolveNetworkEndpointGroupNumericId(ctx, groupIdOrName)
	if err != nil {
		return err
	}

	logicalNetworkIdNumeric, err := parseId(logicalNetworkId, "logical network")
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkEndpointGroupAPI.
		RemoveLogicalNetworkFromNetworkEndpointGroup(ctx, groupId, logicalNetworkIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Logical network '%s' removed from network endpoint group '%s'", logicalNetworkId, groupIdOrName)
	return nil
}

// GetNetworkEndpointGroupByIdOrLabel resolves a network endpoint group by
// numeric ID first and falls back to an exact name match.
func GetNetworkEndpointGroupByIdOrLabel(ctx context.Context, groupIdOrName string) (*sdk.NetworkEndpointGroup, error) {
	client := api.GetApiClient(ctx)

	if idNumeric, err := utils.GetInt64FromString(groupIdOrName); err == nil {
		group, httpRes, err := client.NetworkEndpointGroupAPI.GetNetworkEndpointGroupById(ctx, idNumeric).Execute()
		if err = response_inspector.InspectResponse(httpRes, err); err == nil {
			return group, nil
		}
	}

	list, httpRes, err := client.NetworkEndpointGroupAPI.
		GetNetworkEndpointGroups(ctx).
		FilterName([]string{groupIdOrName}).
		Execute()
	if err = response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	for i := range list.Data {
		if list.Data[i].Name == groupIdOrName {
			return &list.Data[i], nil
		}
	}

	err = fmt.Errorf("network endpoint group '%s' not found", groupIdOrName)
	logger.Get().Error().Err(err).Msg("")
	return nil, err
}

// ResolveNetworkEndpointGroupNumericId resolves a network endpoint group ID or
// name to its numeric ID.
func ResolveNetworkEndpointGroupNumericId(ctx context.Context, groupIdOrName string) (int64, error) {
	groupId, _, err := resolveGroupIdAndRevision(ctx, groupIdOrName)
	return groupId, err
}

func resolveGroupIdAndRevision(ctx context.Context, groupIdOrName string) (int64, string, error) {
	group, err := GetNetworkEndpointGroupByIdOrLabel(ctx, groupIdOrName)
	if err != nil {
		return 0, "", err
	}

	groupId, err := utils.GetInt64FromString(group.Id)
	if err != nil {
		return 0, "", fmt.Errorf("invalid network endpoint group ID %q: %w", group.Id, err)
	}

	return groupId, group.Revision, nil
}

func parseId(value string, what string) (int64, error) {
	id, err := utils.GetInt64FromString(value)
	if err != nil {
		err = fmt.Errorf("invalid %s ID: '%s'", what, value)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}
	return id, nil
}
