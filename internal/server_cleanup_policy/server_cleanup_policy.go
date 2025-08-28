package server_cleanup_policy

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/spf13/cobra"

	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var cleanupPolicyPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Label": {
			Title:    "Label",
			MaxWidth: 30,
			Order:    2,
		},
		"CleanupDrivesForOobEnabledServer": {
			Title:       "Cleanup Drives",
			Transformer: formatter.FormatBooleanValue,
			Order:       3,
		},
		"RecreateRaid": {
			Title:       "Recreate RAID",
			Transformer: formatter.FormatBooleanValue,
			Order:       4,
		},
		"DisableEmbeddedNics": {
			Title:       "Disable Emb. NICs",
			Transformer: formatter.FormatBooleanValue,
			Order:       5,
		},
		"RaidOneDrive": {
			Title: "RAID 1 Drive",
			Order: 6,
		},
		"RaidTwoDrives": {
			Title: "RAID 2 Drives",
			Order: 7,
		},
		"RaidEvenNumberMoreThanTwoDrives": {
			Title: "RAID Even Drives",
			Order: 8,
		},
		"RaidOddNumberMoreThanOneDrive": {
			Title: "RAID Odd Drives",
			Order: 9,
		},
		"SkipRaidActions": {
			Title: "Skip RAID Actions",
			Order: 10,
		},
	},
}

func CleanupPolicyList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all server cleanup policies")

	client := api.GetApiClient(ctx)

	cleanupPoliciesList, httpRes, err := client.ServerCleanupPolicyAPI.GetServerCleanupPolicies(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(cleanupPoliciesList, &cleanupPolicyPrintConfig)
}

func CleanupPolicyGet(ctx context.Context, cleanupPolicyId string) error {
	logger.Get().Info().Msgf("Get server cleanup policy '%s'", cleanupPolicyId)

	cleanupPolicyIdNumber, err := strconv.ParseFloat(cleanupPolicyId, 32)
	if err != nil {
		err := fmt.Errorf("invalid server cleanup policy ID: '%s'", cleanupPolicyId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	cleanupPolicy, httpRes, err := client.ServerCleanupPolicyAPI.GetServerCleanupPolicyInfo(ctx, float32(cleanupPolicyIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(cleanupPolicy, &cleanupPolicyPrintConfig)
}

func CleanupPolicyCreate(ctx context.Context, label string, cleanupDrives, recreateRaid, disableEmbeddedNics bool, raidOneDrive, raidTwoDrives, raidEvenDrives, raidOddDrives, skipRaidActionsStr string) error {
	logger.Get().Info().Msgf("Creating server cleanup policy")

	client := api.GetApiClient(ctx)

	// Parse skip-raid-actions as comma-separated list
	var skipRaidActions []string
	if skipRaidActionsStr != "" {
		skipRaidActions = strings.Split(skipRaidActionsStr, ",")
		for i, action := range skipRaidActions {
			skipRaidActions[i] = strings.TrimSpace(action)
		}
	}

	// Convert boolean values to float32 for API compatibility
	var cleanupDrivesFloat float32
	var recreateRaidFloat float32
	var disableEmbeddedNicsFloat float32

	if cleanupDrives {
		cleanupDrivesFloat = 1
	}
	if recreateRaid {
		recreateRaidFloat = 1
	}
	if disableEmbeddedNics {
		disableEmbeddedNicsFloat = 1
	}

	createRequest := sdk.CreateServerCleanupPolicy{
		Label:                            label,
		CleanupDrivesForOobEnabledServer: cleanupDrivesFloat,
		RecreateRaid:                     recreateRaidFloat,
		DisableEmbeddedNics:              disableEmbeddedNicsFloat,
		RaidOneDrive:                     raidOneDrive,
		RaidTwoDrives:                    raidTwoDrives,
		RaidEvenNumberMoreThanTwoDrives:  raidEvenDrives,
		RaidOddNumberMoreThanOneDrive:    raidOddDrives,
		SkipRaidActions:                  skipRaidActions,
	}

	cleanupPolicy, httpRes, err := client.ServerCleanupPolicyAPI.CreateServerCleanupPolicy(ctx).CreateServerCleanupPolicy(createRequest).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(cleanupPolicy, &cleanupPolicyPrintConfig)
}

func CleanupPolicyUpdate(ctx context.Context, cleanupPolicyId string, label string, cleanupDrives, recreateRaid, disableEmbeddedNics bool, raidOneDrive, raidTwoDrives, raidEvenDrives, raidOddDrives, skipRaidActionsStr string, cmd *cobra.Command) error {
	logger.Get().Info().Msgf("Updating server cleanup policy '%s'", cleanupPolicyId)

	cleanupPolicyIdNumber, err := strconv.ParseFloat(cleanupPolicyId, 32)
	if err != nil {
		err := fmt.Errorf("invalid server cleanup policy ID: '%s'", cleanupPolicyId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	updateRequest := sdk.UpdateServerCleanupPolicy{}

	// Only set fields that were provided
	if cmd.Flags().Changed("label") {
		updateRequest.Label = &label
	}

	if cmd.Flags().Changed("cleanup-drives") {
		var cleanupDrivesFloat float32
		if cleanupDrives {
			cleanupDrivesFloat = 1
		}
		updateRequest.CleanupDrivesForOobEnabledServer = &cleanupDrivesFloat
	}

	if cmd.Flags().Changed("recreate-raid") {
		var recreateRaidFloat float32
		if recreateRaid {
			recreateRaidFloat = 1
		}
		updateRequest.RecreateRaid = &recreateRaidFloat
	}

	if cmd.Flags().Changed("disable-embedded-nics") {
		var disableEmbeddedNicsFloat float32
		if disableEmbeddedNics {
			disableEmbeddedNicsFloat = 1
		}
		updateRequest.DisableEmbeddedNics = &disableEmbeddedNicsFloat
	}

	if cmd.Flags().Changed("raid-one-drive") {
		updateRequest.RaidOneDrive = &raidOneDrive
	}

	if cmd.Flags().Changed("raid-two-drives") {
		updateRequest.RaidTwoDrives = &raidTwoDrives
	}

	if cmd.Flags().Changed("raid-even-drives") {
		updateRequest.RaidEvenNumberMoreThanTwoDrives = &raidEvenDrives
	}

	if cmd.Flags().Changed("raid-odd-drives") {
		updateRequest.RaidOddNumberMoreThanOneDrive = &raidOddDrives
	}

	if cmd.Flags().Changed("skip-raid-actions") {
		var skipRaidActions []string
		if skipRaidActionsStr != "" {
			skipRaidActions = strings.Split(skipRaidActionsStr, ",")
			for i, action := range skipRaidActions {
				skipRaidActions[i] = strings.TrimSpace(action)
			}
		}
		updateRequest.SkipRaidActions = skipRaidActions
	}

	cleanupPolicy, httpRes, err := client.ServerCleanupPolicyAPI.UpdateServerCleanupPolicy(ctx, float32(cleanupPolicyIdNumber)).UpdateServerCleanupPolicy(updateRequest).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(cleanupPolicy, &cleanupPolicyPrintConfig)
}

func CleanupPolicyDelete(ctx context.Context, cleanupPolicyId string) error {
	logger.Get().Info().Msgf("Deleting server cleanup policy '%s'", cleanupPolicyId)

	cleanupPolicyIdNumber, err := strconv.ParseFloat(cleanupPolicyId, 32)
	if err != nil {
		err := fmt.Errorf("invalid server cleanup policy ID: '%s'", cleanupPolicyId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerCleanupPolicyAPI.DeleteServerCleanupPolicy(ctx, float32(cleanupPolicyIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Server cleanup policy '%s' deleted successfully", cleanupPolicyId)
	return nil
}
