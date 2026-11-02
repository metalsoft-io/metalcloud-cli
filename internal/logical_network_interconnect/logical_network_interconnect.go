package logical_network_interconnect

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

var logicalNetworkInterconnectPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"Label": {
			Order:    2,
			MaxWidth: 30,
		},
		"Name": {
			Order:    3,
			MaxWidth: 30,
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
		"Revision": {
			Order: 8,
		},
		"CreatedAt": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       9,
		},
	},
}

var logicalNetworkInterconnectLinkPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"LogicalNetworkInterconnectId": {
			Title: "Interconnect",
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

// LogicalNetworkInterconnectListFilters carries the optional list filters.
type LogicalNetworkInterconnectListFilters struct {
	Id                   []string
	Label                []string
	Name                 []string
	Kind                 []string
	Status               []string
	FabricInterconnectId []string
}

func InterconnectList(ctx context.Context, filters LogicalNetworkInterconnectListFilters) error {
	logger.Get().Info().Msgf("Listing logical network interconnects")

	client := api.GetApiClient(ctx)

	request := client.LogicalNetworkInterconnectAPI.GetLogicalNetworkInterconnects(ctx).SortBy([]string{"id:ASC"})
	if len(filters.Id) > 0 {
		request = request.FilterId(utils.ProcessFilterStringSlice(filters.Id))
	}
	if len(filters.Label) > 0 {
		request = request.FilterLabel(utils.ProcessFilterStringSlice(filters.Label))
	}
	if len(filters.Name) > 0 {
		request = request.FilterName(utils.ProcessFilterStringSlice(filters.Name))
	}
	if len(filters.Kind) > 0 {
		request = request.FilterKind(utils.ProcessFilterStringSlice(filters.Kind))
	}
	if len(filters.Status) > 0 {
		request = request.FilterStatus(utils.ProcessFilterStringSlice(filters.Status))
	}
	if len(filters.FabricInterconnectId) > 0 {
		request = request.FilterFabricInterconnectId(utils.ProcessFilterStringSlice(filters.FabricInterconnectId))
	}

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &logicalNetworkInterconnectPrintConfig)
}

func InterconnectGet(ctx context.Context, interconnectIdOrLabel string) error {
	logger.Get().Info().Msgf("Get logical network interconnect '%s'", interconnectIdOrLabel)

	interconnect, err := GetInterconnectByIdOrLabel(ctx, interconnectIdOrLabel)
	if err != nil {
		return err
	}

	return formatter.PrintResult(interconnect, &logicalNetworkInterconnectPrintConfig)
}

func InterconnectConfigExample(ctx context.Context) error {
	kind := sdk.LOGICALNETWORKINTERCONNECTKIND_DCI_EVPN
	example := sdk.CreateLogicalNetworkInterconnect{
		Label:                "dc1-dc2-logical-interconnect",
		Name:                 "DC1 to DC2 logical interconnect",
		Kind:                 &kind,
		FabricInterconnectId: 1,
		Annotations:          &map[string]string{"environment": "production"},
	}

	if formatter.IsNativeFormat() {
		return formatter.PrintResult(example, nil)
	}
	return formatter.PrintYamlResult(example)
}

func InterconnectCreate(ctx context.Context, create sdk.CreateLogicalNetworkInterconnect) error {
	logger.Get().Info().Msgf("Creating logical network interconnect '%s'", create.Label)

	client := api.GetApiClient(ctx)

	interconnect, httpRes, err := client.LogicalNetworkInterconnectAPI.
		CreateLogicalNetworkInterconnect(ctx).
		CreateLogicalNetworkInterconnect(create).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(interconnect, &logicalNetworkInterconnectPrintConfig)
}

func InterconnectUpdate(ctx context.Context, interconnectIdOrLabel string, config []byte) error {
	logger.Get().Info().Msgf("Updating logical network interconnect '%s'", interconnectIdOrLabel)

	var update sdk.UpdateLogicalNetworkInterconnectDto
	if err := utils.UnmarshalContent(config, &update); err != nil {
		return err
	}

	interconnectId, revision, err := resolveInterconnectIdAndRevision(ctx, interconnectIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	interconnect, httpRes, err := client.LogicalNetworkInterconnectAPI.
		UpdateLogicalNetworkInterconnect(ctx, interconnectId).
		UpdateLogicalNetworkInterconnectDto(update).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(interconnect, &logicalNetworkInterconnectPrintConfig)
}

func InterconnectDelete(ctx context.Context, interconnectIdOrLabel string) error {
	logger.Get().Info().Msgf("Deleting logical network interconnect '%s'", interconnectIdOrLabel)

	interconnectId, revision, err := resolveInterconnectIdAndRevision(ctx, interconnectIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.LogicalNetworkInterconnectAPI.
		DeleteLogicalNetworkInterconnect(ctx, interconnectId).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Logical network interconnect '%s' deleted", interconnectIdOrLabel)
	return nil
}

func InterconnectLinkList(ctx context.Context, interconnectIdOrLabel string, filterLogicalNetworkId []string, filterStatus []string) error {
	logger.Get().Info().Msgf("Listing links of logical network interconnect '%s'", interconnectIdOrLabel)

	interconnectId, err := ResolveInterconnectNumericId(ctx, interconnectIdOrLabel)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.LogicalNetworkInterconnectAPI.
		GetLogicalNetworkInterconnectLinks(ctx, interconnectId).
		SortBy([]string{"id:ASC"})
	if len(filterLogicalNetworkId) > 0 {
		request = request.FilterLogicalNetworkId(utils.ProcessFilterStringSlice(filterLogicalNetworkId))
	}
	if len(filterStatus) > 0 {
		request = request.FilterStatus(utils.ProcessFilterStringSlice(filterStatus))
	}

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &logicalNetworkInterconnectLinkPrintConfig)
}

func InterconnectLinkGet(ctx context.Context, interconnectIdOrLabel string, linkId string) error {
	logger.Get().Info().Msgf("Get link '%s' of logical network interconnect '%s'", linkId, interconnectIdOrLabel)

	interconnectId, err := ResolveInterconnectNumericId(ctx, interconnectIdOrLabel)
	if err != nil {
		return err
	}

	linkIdNumeric, err := parseId(linkId, "logical network interconnect link")
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	link, httpRes, err := client.LogicalNetworkInterconnectAPI.
		GetLogicalNetworkInterconnectLinkById(ctx, interconnectId, linkIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(link, &logicalNetworkInterconnectLinkPrintConfig)
}

func InterconnectLinkAdd(ctx context.Context, interconnectIdOrLabel string, logicalNetworkId string) error {
	logger.Get().Info().Msgf("Adding logical network '%s' to logical network interconnect '%s'", logicalNetworkId, interconnectIdOrLabel)

	interconnectId, err := ResolveInterconnectNumericId(ctx, interconnectIdOrLabel)
	if err != nil {
		return err
	}

	logicalNetworkIdNumeric, err := parseId(logicalNetworkId, "logical network")
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	link, httpRes, err := client.LogicalNetworkInterconnectAPI.
		AddLogicalNetworkToLogicalNetworkInterconnect(ctx, interconnectId).
		AddLogicalNetworkToInterconnect(sdk.AddLogicalNetworkToInterconnect{
			LogicalNetworkId: logicalNetworkIdNumeric,
		}).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(link, &logicalNetworkInterconnectLinkPrintConfig)
}

func InterconnectLinkRemove(ctx context.Context, interconnectIdOrLabel string, linkId string) error {
	logger.Get().Info().Msgf("Removing link '%s' from logical network interconnect '%s'", linkId, interconnectIdOrLabel)

	interconnectId, err := ResolveInterconnectNumericId(ctx, interconnectIdOrLabel)
	if err != nil {
		return err
	}

	linkIdNumeric, err := parseId(linkId, "logical network interconnect link")
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.LogicalNetworkInterconnectAPI.
		RemoveLogicalNetworkFromLogicalNetworkInterconnect(ctx, interconnectId, linkIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Link '%s' removed from logical network interconnect '%s'", linkId, interconnectIdOrLabel)
	return nil
}

// GetInterconnectByIdOrLabel resolves a logical network interconnect by
// numeric ID first and falls back to an exact label match.
func GetInterconnectByIdOrLabel(ctx context.Context, interconnectIdOrLabel string) (*sdk.LogicalNetworkInterconnect, error) {
	client := api.GetApiClient(ctx)

	if idNumeric, err := utils.GetInt64FromString(interconnectIdOrLabel); err == nil {
		interconnect, httpRes, err := client.LogicalNetworkInterconnectAPI.GetLogicalNetworkInterconnectById(ctx, idNumeric).Execute()
		if err = response_inspector.InspectResponse(httpRes, err); err == nil {
			return interconnect, nil
		}
	}

	list, httpRes, err := client.LogicalNetworkInterconnectAPI.
		GetLogicalNetworkInterconnects(ctx).
		FilterLabel([]string{interconnectIdOrLabel}).
		Execute()
	if err = response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}

	for i := range list.Data {
		if list.Data[i].Label == interconnectIdOrLabel {
			return &list.Data[i], nil
		}
	}

	err = fmt.Errorf("logical network interconnect '%s' not found", interconnectIdOrLabel)
	logger.Get().Error().Err(err).Msg("")
	return nil, err
}

// ResolveInterconnectNumericId resolves a logical network interconnect ID or
// label to its numeric ID.
func ResolveInterconnectNumericId(ctx context.Context, interconnectIdOrLabel string) (int64, error) {
	interconnectId, _, err := resolveInterconnectIdAndRevision(ctx, interconnectIdOrLabel)
	return interconnectId, err
}

func resolveInterconnectIdAndRevision(ctx context.Context, interconnectIdOrLabel string) (int64, string, error) {
	interconnect, err := GetInterconnectByIdOrLabel(ctx, interconnectIdOrLabel)
	if err != nil {
		return 0, "", err
	}

	interconnectId, err := utils.GetInt64FromString(interconnect.Id)
	if err != nil {
		return 0, "", fmt.Errorf("invalid logical network interconnect ID %q: %w", interconnect.Id, err)
	}

	return interconnectId, interconnect.Revision, nil
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
