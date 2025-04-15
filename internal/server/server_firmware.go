package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var serverComponentPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
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
	},
}

// ServerFirmwareComponentsList lists all firmware components for a server
func ServerFirmwareComponentsList(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Listing firmware components for server '%s'", serverId)

	serverIdNumeric, err := GetServerId(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	componentsList, httpRes, err := client.ServerFirmwareAPI.GetServerComponents(ctx, serverIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(componentsList, &serverComponentPrintConfig)
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
	err := json.Unmarshal(config, &updateConfig)
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

	component, httpRes, err := client.ServerFirmwareAPI.UpdateServerComponent(ctx, serverIdNumeric, componentIdNumeric).UpdateServerComponent(updateConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Firmware component '%s' updated for server '%s'", componentId, serverId)
	return formatter.PrintResult(component, &serverComponentPrintConfig)
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

	return formatter.PrintResult(inventory, nil)
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
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Firmware upgrade initiated for server '%s'", serverId)
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

// ServerFirmwareComponentUpgrade upgrades firmware for a specific component
func ServerFirmwareComponentUpgrade(ctx context.Context, serverId string, componentId string, config []byte) error {
	logger.Get().Info().Msgf("Upgrading firmware component '%s' for server '%s'", componentId, serverId)

	var upgradeConfig sdk.FirmwareUpgrade
	err := json.Unmarshal(config, &upgradeConfig)
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
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
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
	err := json.Unmarshal(config, &scheduleConfig)
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
	logger.Get().Info().Msgf("Generating firmware upgrade audit")

	var auditConfig sdk.GenerateFirmwareUpgradeAudit
	err := json.Unmarshal(config, &auditConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	audit, httpRes, err := client.ServerFirmwareAPI.GenerateServersFirmwareUpgradeAudit(ctx).GenerateFirmwareUpgradeAudit(auditConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(audit, nil)
}
