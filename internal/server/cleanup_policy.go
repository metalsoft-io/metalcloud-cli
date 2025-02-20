package server

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
)

var cleanupPolicyPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
		},
		"Label": {
			Title: "Label",
		},
		"CleanupDrivesForOobEnabledServer": {
			Title: "Cleanup Server Drives",
		},
		"RecreateRaid": {
			Title: "Recrate RAID",
		},
		"DisableEmbeddedNics": {
			Title: "Disable Embedded NICs",
		},
		"RaidOneDrive": {
			Title: "RAID 1 Drive",
		},
		"RaidTwoDrives": {
			Title: "RAID 2 Drives",
		},
		"RaidEvenNumberMoreThanTwoDrives": {
			Title: "RAID Even Drives",
		},
		"RaidOddNumberMoreThanOneDrive": {
			Title: "RAID Odd Drives",
		},
		"SkipRaidActions": {
			Title: "Skip RAID Actions",
		},
	},
}

func CleanupPolicyList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all server cleanup policies")

	client, err := system.GetApiClient(ctx)
	if err != nil {
		return err
	}

	cleanupPoliciesList, httpRes, err := client.ServerCleanupPolicyAPI.GetServerCleanupPolicies(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(cleanupPoliciesList, &cleanupPolicyPrintConfig)
}

func CleanupPolicyGet(ctx context.Context, cleanupPolicyId string) error {
	logger.Get().Info().Msgf("Get server cleanup policy '%s'", cleanupPolicyId)

	client, err := system.GetApiClient(ctx)
	if err != nil {
		return err
	}

	cleanupPolicyIdNumber, err := strconv.ParseFloat(cleanupPolicyId, 32)
	if err != nil {
		err := fmt.Errorf("invalid server cleanup policy ID: '%s'", cleanupPolicyId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	cleanupPolicy, httpRes, err := client.ServerCleanupPolicyAPI.GetServerCleanupPolicyInfo(ctx, float32(cleanupPolicyIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(cleanupPolicy, &cleanupPolicyPrintConfig)
}
