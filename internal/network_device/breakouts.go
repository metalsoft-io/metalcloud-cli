package network_device

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var networkDeviceBreakoutPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"PortName": {
			Title: "Port",
			Order: 2,
		},
		"AppliedGroups": {
			Title:    "Applied Groups",
			MaxWidth: 30,
			Order:    3,
		},
		"StagedGroups": {
			Title:    "Staged Groups",
			MaxWidth: 30,
			Order:    4,
		},
		"ServiceStatus": {
			Title:       "Status",
			Transformer: formatter.FormatStatusValue,
			Order:       5,
		},
		"PendingDelete": {
			Title: "Pending Delete",
			Order: 6,
		},
		"Revision": {
			Order: 7,
		},
	},
}

var networkDeviceBreakoutConfigPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"StagedGroups": {
			Title:    "Staged Groups",
			MaxWidth: 60,
			Order:    1,
		},
		"Revision": {
			Order: 2,
		},
	},
}

// breakoutRow is the table projection of a breakout: the breakout group lists
// are flattened into readable summaries because the tabular formatter renders
// struct-valued fields as empty cells.
type breakoutRow struct {
	Id            int64
	PortName      string
	AppliedGroups string
	StagedGroups  string
	ServiceStatus string
	PendingDelete bool
	Revision      int64
}

// breakoutConfigRow is the table projection of a breakout config buffer.
type breakoutConfigRow struct {
	StagedGroups string
	Revision     int64
}

// formatBreakoutGroups renders breakout groups as "4x25G, 2x50G". A group
// without a forced speed (the device derives it from the port lane capability)
// renders as "4x".
func formatBreakoutGroups(groups []sdk.BreakoutGroup) string {
	if len(groups) == 0 {
		return ""
	}

	parts := make([]string, 0, len(groups))
	for _, group := range groups {
		part := fmt.Sprintf("%dx", group.NumberOfInterfaces)
		if group.Speed.IsSet() && group.Speed.Get() != nil {
			part += *group.Speed.Get()
		}
		parts = append(parts, part)
	}

	return strings.Join(parts, ", ")
}

func toBreakoutRow(breakout sdk.NetworkEquipmentBreakout) breakoutRow {
	return breakoutRow{
		Id:            breakout.Id,
		PortName:      breakout.PortName,
		AppliedGroups: formatBreakoutGroups(breakout.BreakoutGroups),
		StagedGroups:  formatBreakoutGroups(breakout.Config.BreakoutGroups),
		ServiceStatus: breakout.ServiceStatus,
		PendingDelete: breakout.PendingDelete,
		Revision:      breakout.Revision,
	}
}

func toBreakoutConfigRow(breakoutConfig *sdk.NetworkEquipmentBreakoutConfigDto) breakoutConfigRow {
	return breakoutConfigRow{
		StagedGroups: formatBreakoutGroups(breakoutConfig.BreakoutGroups),
		Revision:     breakoutConfig.Revision,
	}
}

func NetworkDeviceBreakoutList(ctx context.Context, networkDeviceRef string) error {
	logger.Get().Info().Msgf("Listing breakouts of network device '%s'", networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	breakouts, httpRes, err := client.NetworkDeviceAPI.
		ListNetworkDeviceBreakouts(ctx, networkDeviceIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	rows := make([]breakoutRow, 0, len(breakouts))
	for _, breakout := range breakouts {
		rows = append(rows, toBreakoutRow(breakout))
	}

	return printRecords(breakouts, rows, &networkDeviceBreakoutPrintConfig)
}

func NetworkDeviceBreakoutGet(ctx context.Context, networkDeviceRef string, breakoutId string) error {
	logger.Get().Info().Msgf("Getting breakout %s of network device '%s'", breakoutId, networkDeviceRef)

	networkDeviceIdNumeric, breakoutIdNumeric, err := resolveBreakoutIds(ctx, networkDeviceRef, breakoutId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	breakout, httpRes, err := client.NetworkDeviceAPI.
		GetNetworkDeviceBreakout(ctx, networkDeviceIdNumeric, breakoutIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if breakout == nil {
		return errEmptyBody("breakout")
	}

	return printRecords(breakout, toBreakoutRow(*breakout), &networkDeviceBreakoutPrintConfig)
}

func NetworkDeviceBreakoutGetConfig(ctx context.Context, networkDeviceRef string, breakoutId string) error {
	logger.Get().Info().Msgf("Getting config of breakout %s of network device '%s'", breakoutId, networkDeviceRef)

	networkDeviceIdNumeric, breakoutIdNumeric, err := resolveBreakoutIds(ctx, networkDeviceRef, breakoutId)
	if err != nil {
		return err
	}

	breakoutConfig, err := getNetworkDeviceBreakoutConfig(ctx, networkDeviceIdNumeric, breakoutIdNumeric)
	if err != nil {
		return err
	}

	return printRecords(breakoutConfig, toBreakoutConfigRow(breakoutConfig), &networkDeviceBreakoutConfigPrintConfig)
}

func NetworkDeviceBreakoutCreate(ctx context.Context, networkDeviceRef string, config []byte) error {
	logger.Get().Info().Msgf("Creating breakout on network device '%s'", networkDeviceRef)

	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return err
	}

	var breakoutConfig sdk.CreateNetworkEquipmentBreakout
	if err := utils.UnmarshalContent(config, &breakoutConfig); err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	breakout, httpRes, err := client.NetworkDeviceAPI.
		CreateNetworkDeviceBreakout(ctx, networkDeviceIdNumeric).
		CreateNetworkEquipmentBreakout(breakoutConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if breakout == nil {
		return errEmptyBody("breakout")
	}

	return printRecords(breakout, toBreakoutRow(*breakout), &networkDeviceBreakoutPrintConfig)
}

// NetworkDeviceBreakoutUpdateConfig stages a new set of breakout groups on an
// existing breakout. The config sub-resource carries its own optimistic-lock
// revision, independent of the breakout entity revision.
func NetworkDeviceBreakoutUpdateConfig(ctx context.Context, networkDeviceRef string, breakoutId string, config []byte) error {
	logger.Get().Info().Msgf("Updating config of breakout %s of network device '%s'", breakoutId, networkDeviceRef)

	networkDeviceIdNumeric, breakoutIdNumeric, err := resolveBreakoutIds(ctx, networkDeviceRef, breakoutId)
	if err != nil {
		return err
	}

	var updateConfig sdk.UpdateNetworkEquipmentBreakoutConfigDto
	if err := utils.UnmarshalContent(config, &updateConfig); err != nil {
		return err
	}

	currentConfig, err := getNetworkDeviceBreakoutConfig(ctx, networkDeviceIdNumeric, breakoutIdNumeric)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	updatedConfig, httpRes, err := client.NetworkDeviceAPI.
		UpdateNetworkDeviceBreakoutConfig(ctx, networkDeviceIdNumeric, breakoutIdNumeric).
		UpdateNetworkEquipmentBreakoutConfigDto(updateConfig).
		IfMatch(strconv.FormatInt(currentConfig.Revision, 10)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if updatedConfig == nil {
		return errEmptyBody("breakout config")
	}

	return printRecords(updatedConfig, toBreakoutConfigRow(updatedConfig), &networkDeviceBreakoutConfigPrintConfig)
}

func NetworkDeviceBreakoutDelete(ctx context.Context, networkDeviceRef string, breakoutId string) error {
	logger.Get().Info().Msgf("Deleting breakout %s of network device '%s'", breakoutId, networkDeviceRef)

	networkDeviceIdNumeric, breakoutIdNumeric, err := resolveBreakoutIds(ctx, networkDeviceRef, breakoutId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.NetworkDeviceAPI.
		DeleteNetworkDeviceBreakout(ctx, networkDeviceIdNumeric, breakoutIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Breakout %s of network device '%s' deleted", breakoutId, networkDeviceRef)
	return nil
}

func NetworkDeviceBreakoutConfigExample(ctx context.Context) error {
	breakoutConfig := sdk.CreateNetworkEquipmentBreakout{
		PortName: "Ethernet0",
		Groups: []sdk.BreakoutGroup{
			{NumberOfInterfaces: 4, Speed: *sdk.NewNullableString(sdk.PtrString("25G"))},
		},
	}

	return formatter.PrintResult(breakoutConfig, nil)
}

func NetworkDeviceBreakoutUpdateConfigExample(ctx context.Context) error {
	updateConfig := sdk.UpdateNetworkEquipmentBreakoutConfigDto{
		BreakoutGroups: []sdk.BreakoutGroup{
			{NumberOfInterfaces: 4, Speed: *sdk.NewNullableString(sdk.PtrString("25G"))},
		},
	}

	return formatter.PrintResult(updateConfig, nil)
}

func resolveBreakoutIds(ctx context.Context, networkDeviceRef string, breakoutId string) (int64, int64, error) {
	networkDeviceIdNumeric, err := resolveNetworkDeviceId(ctx, networkDeviceRef)
	if err != nil {
		return 0, 0, err
	}

	breakoutIdNumeric, err := utils.GetInt64FromString(breakoutId)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid breakout ID: '%s'", breakoutId)
	}

	return networkDeviceIdNumeric, breakoutIdNumeric, nil
}

func getNetworkDeviceBreakoutConfig(ctx context.Context, networkDeviceId int64, breakoutId int64) (*sdk.NetworkEquipmentBreakoutConfigDto, error) {
	client := api.GetApiClient(ctx)

	breakoutConfig, httpRes, err := client.NetworkDeviceAPI.
		GetNetworkDeviceBreakoutConfig(ctx, networkDeviceId, breakoutId).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return nil, err
	}
	if breakoutConfig == nil {
		return nil, errEmptyBody("breakout config")
	}

	return breakoutConfig, nil
}
