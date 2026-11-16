package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

// ServerFirmwareInventoryEntry is the typed view of one Redfish
// SoftwareInventory entry returned by GetServerFirmwareInventory. The SDK
// returns []map[string]interface{}; we round-trip into this struct so the
// formatter can drive a table from PrintConfig.
type ServerFirmwareInventoryEntry struct {
	Id                           string                         `json:"Id"`
	Name                         string                         `json:"Name"`
	Description                  string                         `json:"Description,omitempty"`
	Version                      string                         `json:"Version,omitempty"`
	Updateable                   bool                           `json:"Updateable"`
	Status                       *ServerFirmwareInventoryStatus `json:"Status,omitempty"`
	MatchedComponents            []map[string]interface{}       `json:"MatchedComponents,omitempty"`
	MatchedComponentsDeviceClass string                         `json:"MatchedComponentsDeviceClass,omitempty"`
	Oem                          map[string]interface{}         `json:"Oem,omitempty"`
}

type ServerFirmwareInventoryStatus struct {
	Health string `json:"Health,omitempty"`
	State  string `json:"State,omitempty"`
}

var serverFirmwareInventoryPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id":          {Title: "ID", Order: 1},
		"Name":        {Title: "Name", MaxWidth: 40, Order: 2},
		"Description": {Title: "Description", MaxWidth: 32, Order: 3},
		"Version":     {Title: "Version", Order: 4},
		"Updateable": {
			Title: "Updateable",
			Order: 5,
			Transformer: func(v interface{}) string {
				if b, ok := v.(bool); ok && b {
					return "Yes"
				}
				return "No"
			},
		},
		"Status.Health": {Title: "Health", Order: 6},
		"Status.State":  {Title: "State", Order: 7},
		"MatchedComponents": {
			Title: "Matches",
			Order: 8,
			Transformer: func(v interface{}) string {
				if arr, ok := v.([]interface{}); ok {
					return fmt.Sprintf("%d", len(arr))
				}
				return "0"
			},
		},
	},
}

var serverComponentPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
			Order: 1,
		},
		"ExternalId": {
			Title: "External ID",
			Order: 2,
		},
		"Name": {
			Title: "Name",
			Order: 3,
		},
		"Type": {
			Title: "Type",
			Order: 4,
		},
		"FirmwareVersion": {
			Title: "Current Version",
			Order: 5,
		},
		"FirmwareTargetVersion": {
			Title: "Target Version",
			Order: 6,
		},
		"FirmwareUpdateable": {
			Title: "Updateable",
			Order: 7,
			Transformer: func(value interface{}) string {
				if val, ok := value.(float32); ok {
					if val == 1 {
						return "Yes"
					}
					return "No"
				}
				return fmt.Sprintf("%v", value)
			},
		},
		"FirmwareStatus": {
			Title: "Status",
			Order: 8,
		},
		"FirmwareUpdateTimestamp": {
			Title: "Last Update",
			Order: 9,
		},
		"FirmwareScheduledTimestamp": {
			Title: "Scheduled",
			Order: 10,
		},
		"ComponentClass": {
			Title: "Class",
			Order: 11,
		},
		"ComponentVendor": {
			Title: "Vendor",
			Order: 12,
		},
		"ComponentModel": {
			Title:    "Model",
			Order:    13,
			MaxWidth: 30,
		},
		"BaseboardId": {
			Title: "Baseboard",
			Order: 14,
		},
	},
}

// ComponentsListFlags filters `server firmware components`; every filter maps to a `filter.*` query parameter.
type ComponentsListFlags struct {
	FilterComponentClass []string
	FilterBaseboardId    []string
}

// ServerFirmwareComponentsList lists all firmware components for a server
func ServerFirmwareComponentsList(ctx context.Context, serverId string) error {
	return ServerFirmwareComponentsListFiltered(ctx, serverId, ComponentsListFlags{})
}

// ServerFirmwareComponentsListFiltered lists the firmware components of a server, filtered by class or baseboard
func ServerFirmwareComponentsListFiltered(ctx context.Context, serverId string, flags ComponentsListFlags) error {
	logger.Get().Info().Msgf("Listing firmware components for server '%s'", serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	request := client.ServerFirmwareAPI.GetServerComponents(ctx, serverIdNumeric).SortBy([]string{"id:ASC"})
	if len(flags.FilterComponentClass) > 0 {
		request = request.FilterComponentClass(utils.ProcessFilterStringSlice(flags.FilterComponentClass))
	}
	if len(flags.FilterBaseboardId) > 0 {
		request = request.FilterBaseboardId(utils.ProcessFilterStringSlice(flags.FilterBaseboardId))
	}

	components, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(components, meta, len(components), &serverComponentPrintConfig)
}

// ServerFirmwareComponentGet retrieves information about a specific component
func ServerFirmwareComponentGet(ctx context.Context, serverId string, componentId string) error {
	logger.Get().Info().Msgf("Getting firmware component '%s' for server '%s'", componentId, serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	componentIdNumeric, err := GetServerId(componentId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	component, httpRes, err := client.ServerFirmwareAPI.GetServerComponentInfo(ctx, serverIdNumeric, componentIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(component, &serverComponentPrintConfig)
}

// ServerFirmwareComponentUpdate updates a server firmware component
func ServerFirmwareComponentUpdate(ctx context.Context, serverId string, componentId string, config []byte) error {
	logger.Get().Info().Msgf("Updating firmware component '%s' for server '%s'", componentId, serverId)

	var updateConfig sdk.UpdateServerComponent
	err := utils.UnmarshalContent(config, &updateConfig)
	if err != nil {
		return err
	}

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	componentIdNumeric, err := GetServerId(componentId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// A target from a firmware package delivered to a target other components share (an HGX tray) is set, or
	// reset, on every component the package updates: the targets before the change show which
	var targetsBefore map[int64]string
	if updateConfig.FirmwareTargetVersion != nil || configSetsTargetVersion(config) {
		targetsBefore, err = serverComponentTargets(ctx, serverIdNumeric)
		if err != nil {
			return err
		}
	}

	component, httpRes, err := client.ServerFirmwareAPI.UpdateServerComponent(ctx, serverIdNumeric, componentIdNumeric).UpdateServerComponent(updateConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Firmware component '%s' updated for server '%s'", componentId, serverId)
	if err := formatter.PrintResult(component, &serverComponentPrintConfig); err != nil {
		return err
	}

	if targetsBefore == nil || !formatter.IsTextFormat() {
		return nil
	}
	changed, err := changedServerComponentTargets(ctx, serverIdNumeric, componentIdNumeric, targetsBefore)
	if err != nil {
		return err
	}
	if len(changed) > 0 {
		fmt.Printf("\nThe same firmware package updates these components, so their target changed too (%d):\n", len(changed))
		return formatter.PrintResult(changed, &packageTargetPrintConfig)
	}
	return nil
}

// configSetsTargetVersion is whether the update sets or resets (null) firmwareTargetVersion.
func configSetsTargetVersion(config []byte) bool {
	var fields map[string]any
	if utils.UnmarshalContent(config, &fields) != nil {
		return false
	}
	_, present := fields["firmwareTargetVersion"]
	return present
}

// PackageTargetChange is a component whose target version changed with another component's.
type PackageTargetChange struct {
	ComponentId    int64
	Component      string
	PreviousTarget string
	TargetVersion  string
}

var packageTargetPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ComponentId":    {Title: "Component ID", Order: 1},
		"Component":      {Title: "Component", Order: 2, MaxWidth: 40},
		"PreviousTarget": {Title: "Previous Target", Order: 3},
		"TargetVersion":  {Title: "Target Version", Order: 4},
	},
}

func serverComponentTargets(ctx context.Context, serverId int64) (map[int64]string, error) {
	client := api.GetApiClient(ctx)
	components, _, err := utils.FetchAllPages(client.ServerFirmwareAPI.GetServerComponents(ctx, serverId).SortBy([]string{"id:ASC"}))
	if err != nil {
		return nil, err
	}
	targets := map[int64]string{}
	for _, component := range components {
		targets[component.Id] = derefString(component.FirmwareTargetVersion)
	}
	return targets, nil
}

// changedServerComponentTargets lists the other components whose target changed since targetsBefore.
func changedServerComponentTargets(ctx context.Context, serverId int64, componentId int64, targetsBefore map[int64]string) ([]PackageTargetChange, error) {
	client := api.GetApiClient(ctx)
	components, _, err := utils.FetchAllPages(client.ServerFirmwareAPI.GetServerComponents(ctx, serverId).SortBy([]string{"id:ASC"}))
	if err != nil {
		return nil, err
	}
	changed := []PackageTargetChange{}
	for _, component := range components {
		target := derefString(component.FirmwareTargetVersion)
		if component.Id == componentId || target == targetsBefore[component.Id] {
			continue
		}
		changed = append(changed, PackageTargetChange{
			ComponentId:    component.Id,
			Component:      derefString(component.Name),
			PreviousTarget: targetsBefore[component.Id],
			TargetVersion:  target,
		})
	}
	return changed, nil
}

// ServerFirmwareUpdateInfo updates firmware information for a server
func ServerFirmwareUpdateInfo(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Updating firmware information for server '%s'", serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerFirmwareAPI.UpdateServerFirmwareInfo(ctx, serverIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Firmware information updated for server '%s'", serverId)
	return nil
}

// ServerFirmwareInventory retrieves firmware inventory from redfish
func ServerFirmwareInventory(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Retrieving firmware inventory for server '%s'", serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	inventory, httpRes, err := client.ServerFirmwareAPI.GetServerFirmwareInventory(ctx, serverIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	raw, err := json.Marshal(inventory)
	if err != nil {
		return fmt.Errorf("failed to re-encode firmware inventory: %w", err)
	}
	entries := []ServerFirmwareInventoryEntry{}
	if err := json.Unmarshal(raw, &entries); err != nil {
		return fmt.Errorf("failed to decode firmware inventory: %w", err)
	}

	return formatter.PrintResult(entries, &serverFirmwareInventoryPrintConfig)
}

// ServerFirmwareUpgrade upgrades firmware for all components on a server
func ServerFirmwareUpgrade(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Upgrading firmware for server '%s'", serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	jobInfo, httpRes, err := client.ServerFirmwareAPI.UpgradeFirmwareOfServer(ctx, serverIdNumeric).Execute()
	if err := inspectFirmwareJobResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Firmware upgrade initiated for server '%s'", serverId)
	// The components the safety checks left out are part of the JSON/YAML output; the text output lists them
	if formatter.IsNativeFormat() {
		return formatter.PrintResult(jobInfo, nil)
	}
	if err := formatter.PrintResult(jobInfo, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"JobId": {
				Title: "Job ID",
				Order: 1,
			},
			"JobGroupId": {
				Title: "Job Group ID",
				Order: 2,
			},
			"Status": {
				Title: "Status",
				Order: 3,
			},
		},
	}); err != nil {
		return err
	}
	if len(jobInfo.RefusedComponents) > 0 {
		fmt.Printf("\nLeft out by the safety checks (%d):\n", len(jobInfo.RefusedComponents))
		return formatter.PrintResult(jobInfo.RefusedComponents, &firmwareRefusalPrintConfig)
	}
	return nil
}

// ServerFirmwareComponentUpgrade upgrades firmware for a specific component
func ServerFirmwareComponentUpgrade(ctx context.Context, serverId string, componentId string, config []byte) error {
	logger.Get().Info().Msgf("Upgrading firmware component '%s' for server '%s'", componentId, serverId)

	var upgradeConfig sdk.FirmwareUpgrade
	err := utils.UnmarshalContent(config, &upgradeConfig)
	if err != nil {
		return err
	}

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	componentIdNumeric, err := GetServerId(componentId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	jobInfo, httpRes, err := client.ServerFirmwareAPI.UpgradeFirmwareOfServerComponent(ctx, serverIdNumeric, componentIdNumeric).FirmwareUpgrade(upgradeConfig).Execute()
	if err := inspectFirmwareJobResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Firmware upgrade initiated for component '%s' on server '%s'", componentId, serverId)
	return formatter.PrintResult(jobInfo, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"JobId": {
				Title: "Job ID",
				Order: 1,
			},
			"JobGroupId": {
				Title: "Job Group ID",
				Order: 2,
			},
			"Status": {
				Title: "Status",
				Order: 3,
			},
		},
	})
}

// ServerFirmwareScheduleUpgrade schedules a firmware upgrade for a server
func ServerFirmwareScheduleUpgrade(ctx context.Context, serverId string, config []byte) error {
	logger.Get().Info().Msgf("Scheduling firmware upgrade for server '%s'", serverId)

	var scheduleConfig sdk.ScheduleFirmwareUpgrade
	err := utils.UnmarshalContent(config, &scheduleConfig)
	if err != nil {
		return err
	}

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerFirmwareAPI.ScheduleServerFirmwareUpgrade(ctx, serverIdNumeric).ScheduleFirmwareUpgrade(scheduleConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Firmware upgrade scheduled for server '%s'", serverId)
	return nil
}

// ServerFirmwareFetchVersions fetches available firmware versions for a server
func ServerFirmwareFetchVersions(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Fetching available firmware versions for server '%s'", serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerFirmwareAPI.FetchAndUpdateServerFirmwareAvailableVersions(ctx, serverIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Available firmware versions fetched for server '%s'", serverId)
	return nil
}

// ServerFirmwareGenerateAudit generates a firmware upgrade audit for servers
func ServerFirmwareGenerateAudit(ctx context.Context, config []byte) error {
	return ServerFirmwareGenerateAuditWithPreflight(ctx, config, false)
}

// ServerFirmwareGenerateAuditWithPreflight generates the audit and, with includePreflight, runs on each server
// the safety checks an upgrade of the target versions would run now (BMC checks included) and lists the jobs it
// would create. Nothing is changed.
func ServerFirmwareGenerateAuditWithPreflight(ctx context.Context, config []byte, includePreflight bool) error {
	logger.Get().Info().Msgf("Generating firmware upgrade audit")

	var auditConfig sdk.GenerateFirmwareUpgradeAudit
	err := utils.UnmarshalContent(config, &auditConfig)
	if err != nil {
		return err
	}
	if includePreflight {
		auditConfig.IncludePreflight = sdk.PtrBool(true)
	}

	client := api.GetApiClient(ctx)

	audit, httpRes, err := client.ServerFirmwareAPI.GenerateServersFirmwareUpgradeAudit(ctx).GenerateFirmwareUpgradeAudit(auditConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	// JSON and YAML output keep the full audit (baseline, upgrade and downgrade options).
	if formatter.IsNativeFormat() {
		return formatter.PrintResult(audit, nil)
	}

	if err := formatter.PrintResult(flattenServerFirmwareAudit(audit), &serverFirmwareAuditPrintConfig); err != nil {
		return err
	}
	if !formatter.IsTextFormat() {
		return nil
	}
	return printFirmwareAuditDetails(audit)
}

// ServerFirmwareAuditRow is the tabular view of the audit of one server component.
// A server that could not be audited is represented by a single row carrying the error.
type ServerFirmwareAuditRow struct {
	ServerId       int64
	Baseline       string
	ComponentId    string
	Component      string
	Type           string
	CurrentVersion string
	TargetVersion  string
	LatestVersion  string
	Status         string
	BelowMinimum   string
	Upgrades       int
	Downgrades     int
	Checks         string
	Error          string
}

var serverFirmwareAuditPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"ServerId":       {Title: "Server ID", Order: 1},
		"Baseline":       {Title: "Baseline", Order: 2},
		"ComponentId":    {Title: "Component ID", Order: 3},
		"Component":      {Title: "Component", Order: 4, MaxWidth: 40},
		"Type":           {Title: "Type", Order: 5},
		"CurrentVersion": {Title: "Current Version", Order: 6},
		"TargetVersion":  {Title: "Target Version", Order: 7},
		"LatestVersion":  {Title: "Latest Version", Order: 8},
		"Status":         {Title: "Status", Order: 9},
		"BelowMinimum":   {Title: "Below Minimum", Order: 10},
		"Upgrades":       {Title: "Upgrades", Order: 11},
		"Downgrades":     {Title: "Downgrades", Order: 12},
		"Checks":         {Title: "Checks", Order: 13, MaxWidth: 50},
		"Error":          {Title: "Error", Order: 14, MaxWidth: 60},
	},
}

func flattenServerFirmwareAudit(audits []sdk.ServerFirmwareAudit) []ServerFirmwareAuditRow {
	rows := []ServerFirmwareAuditRow{}

	for _, audit := range audits {
		baseline := ""
		if audit.Baseline != nil {
			baseline = audit.Baseline.Name
		}

		if audit.Error != nil {
			rows = append(rows, ServerFirmwareAuditRow{
				ServerId: audit.ServerId,
				Baseline: baseline,
				Error:    fmt.Sprintf("%d: %s", int(audit.Error.StatusCode), audit.Error.Message),
			})
			continue
		}

		for _, component := range audit.Components {
			belowMinimum := ""
			if component.BelowMinimumVersion != nil {
				belowMinimum = "No"
				if *component.BelowMinimumVersion {
					belowMinimum = "Yes"
				}
			}

			rows = append(rows, ServerFirmwareAuditRow{
				ServerId:       audit.ServerId,
				Baseline:       baseline,
				ComponentId:    fmt.Sprintf("%d", component.ComponentId),
				Component:      component.Name,
				Type:           derefString(component.Type),
				CurrentVersion: derefString(component.CurrentVersion),
				TargetVersion:  derefString(component.TargetVersion),
				LatestVersion:  derefString(component.LatestVersion),
				Status:         component.Status,
				BelowMinimum:   belowMinimum,
				Upgrades:       len(component.Upgrades),
				Downgrades:     len(component.Downgrades),
				Checks:         summarizeFirmwareChecks(component.Checks),
			})
		}
	}

	return rows
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

var serverFirmwareBatchUpgradePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Successful": {
			Title:    "Successful",
			MaxWidth: 60,
			Order:    1,
		},
		"Failed": {
			Title:    "Failed",
			MaxWidth: 60,
			Order:    2,
		},
	},
}

// ServerFirmwareUpgradeBatch upgrades the firmware of several servers at once.
func ServerFirmwareUpgradeBatch(ctx context.Context, serverIds []string) error {
	logger.Get().Info().Msgf("Upgrading firmware for %d servers", len(serverIds))

	serverIdsNumeric, err := utils.GetInt64SliceFromStrings(serverIds)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	// The body is always sent: skipping the optional setter would serialize a
	// literal null, which the API rejects with 400.
	upgradeConfig := sdk.BatchServerFirmwareUpgrade{
		ServerIds: serverIdsNumeric,
	}

	result, httpRes, err := client.ServerFirmwareAPI.
		UpgradeFirmwareOfServersBatch(ctx).
		BatchServerFirmwareUpgrade(upgradeConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(result, &serverFirmwareBatchUpgradePrintConfig)
}

// ServerFirmwareScheduleUpgradeBatch schedules a firmware upgrade for a batch of
// servers. The full list of server ids is carried in the request body.
func ServerFirmwareScheduleUpgradeBatch(ctx context.Context, serverIds []string, scheduleTimestamp string, confirmationRequired bool) error {
	if len(serverIds) == 0 {
		err := fmt.Errorf("at least one server ID is required")
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	logger.Get().Info().Msgf("Scheduling batch firmware upgrade for %d servers", len(serverIds))

	serverIdsNumeric, err := utils.GetInt64SliceFromStrings(serverIds)
	if err != nil {
		return err
	}

	// The body is always sent: skipping the optional setter would serialize a
	// literal null, which the API rejects with 400.
	scheduleConfig := sdk.BatchScheduleServerFirmwareUpgrade{
		ServerIds:            serverIdsNumeric,
		ConfirmationRequired: sdk.PtrBool(confirmationRequired),
	}

	if scheduleTimestamp != "" {
		scheduleTime, err := time.Parse(time.RFC3339, scheduleTimestamp)
		if err != nil {
			return fmt.Errorf("invalid schedule timestamp '%s' (expected RFC3339, e.g. 2026-01-02T15:04:05Z): %w", scheduleTimestamp, err)
		}
		scheduleConfig.ScheduleUpdateTimestamp = &scheduleTime
	}

	client := api.GetApiClient(ctx)

	result, httpRes, err := client.ServerFirmwareAPI.
		BatchScheduleServerFirmwareUpgrade(ctx).
		BatchScheduleServerFirmwareUpgrade(scheduleConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(result, &formatter.PrintConfig{
		FieldsConfig: map[string]formatter.RecordFieldConfig{
			"Successful": {
				Title: "Scheduled",
				Order: 1,
			},
			"Failed": {
				Title:    "Failed",
				MaxWidth: 60,
				Order:    2,
			},
		},
	})
}
