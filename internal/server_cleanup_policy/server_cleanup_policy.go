package server_cleanup_policy

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
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
			Title: "Cleanup Drives",
			Order: 3,
		},
		"RecreateRaid": {
			Title: "Recrate RAID",
			Order: 4,
		},
		"DisableEmbeddedNics": {
			Title: "Disable Emb. NICs",
			Order: 5,
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
