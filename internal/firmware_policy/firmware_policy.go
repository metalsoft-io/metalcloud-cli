package firmware_policy

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var firmwarePolicyPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Label": {
			MaxWidth: 40,
			Order:    2,
		},
		"Status": {
			Order: 3,
		},
		"Action": {
			Order: 4,
		},
		"UserIdOwner": {
			Title: "Owner ID",
			Order: 5,
		},
		"CreatedTimestamp": {
			Title:       "Created",
			Transformer: formatter.FormatDateTimeValue,
			Order:       6,
		},
		"UpdatedTimestamp": {
			Title:       "Updated",
			Transformer: formatter.FormatDateTimeValue,
			Order:       7,
		},
	},
}

var globalFirmwareConfigPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		// Define fields for global firmware configuration
		// Expand this based on the properties of GlobalFirmwareUpgradeConfiguration
	},
}

func FirmwarePolicyList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all firmware policies")

	client := api.GetApiClient(ctx)

	firmwarePolicyList, httpRes, err := client.FirmwarePolicyAPI.GetFirmwarePolicies(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(firmwarePolicyList, &firmwarePolicyPrintConfig)
}

func FirmwarePolicyGet(ctx context.Context, firmwarePolicyId string) error {
	logger.Get().Info().Msgf("Get firmware policy '%s' details", firmwarePolicyId)

	firmwarePolicyIdNumeric, err := getFirmwarePolicyId(firmwarePolicyId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	firmwarePolicy, httpRes, err := client.FirmwarePolicyAPI.GetFirmwarePolicyInfo(ctx, firmwarePolicyIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(firmwarePolicy, &firmwarePolicyPrintConfig)
}

func FirmwarePolicyCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating firmware policy")

	var firmwarePolicyConfig sdk.CreateServerFirmwareUpgradePolicy
	err := json.Unmarshal(config, &firmwarePolicyConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	firmwarePolicy, httpRes, err := client.FirmwarePolicyAPI.
		CreateFirmwarePolicy(ctx).
		CreateServerFirmwareUpgradePolicy(firmwarePolicyConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(firmwarePolicy, &firmwarePolicyPrintConfig)
}

func FirmwarePolicyUpdate(ctx context.Context, firmwarePolicyId string, config []byte) error {
	logger.Get().Info().Msgf("Updating firmware policy '%s'", firmwarePolicyId)

	firmwarePolicyIdNumeric, err := getFirmwarePolicyId(firmwarePolicyId)
	if err != nil {
		return err
	}

	var firmwarePolicyConfig sdk.UpdateServerFirmwareUpgradePolicy
	err = json.Unmarshal(config, &firmwarePolicyConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	firmwarePolicy, httpRes, err := client.FirmwarePolicyAPI.
		UpdateFirmwarePolicy(ctx, firmwarePolicyIdNumeric).
		UpdateServerFirmwareUpgradePolicy(firmwarePolicyConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(firmwarePolicy, &firmwarePolicyPrintConfig)
}

func FirmwarePolicyDelete(ctx context.Context, firmwarePolicyId string) error {
	logger.Get().Info().Msgf("Deleting firmware policy '%s'", firmwarePolicyId)

	firmwarePolicyIdNumeric, err := getFirmwarePolicyId(firmwarePolicyId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.FirmwarePolicyAPI.
		DeleteFirmwarePolicy(ctx, firmwarePolicyIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Firmware policy '%s' deleted", firmwarePolicyId)
	return nil
}

func FirmwarePolicyGenerateAudit(ctx context.Context, firmwarePolicyId string) error {
	logger.Get().Info().Msgf("Generating audit for firmware policy '%s'", firmwarePolicyId)

	firmwarePolicyIdNumeric, err := getFirmwarePolicyId(firmwarePolicyId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	audit, httpRes, err := client.FirmwarePolicyAPI.
		GenerateFirmwarePolicyAudit(ctx, firmwarePolicyIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(audit, nil)
}

func FirmwarePolicyApplyWithGroups(ctx context.Context) error {
	logger.Get().Info().Msgf("Applying all firmware policies linked to server instance groups")

	client := api.GetApiClient(ctx)

	result, httpRes, err := client.FirmwarePolicyAPI.
		ApplyFirmwarePoliciesWithServerInstanceGroups(ctx).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(result, nil)
}

func FirmwarePolicyApplyWithoutGroups(ctx context.Context) error {
	logger.Get().Info().Msgf("Applying all firmware policies not linked to server instance groups")

	client := api.GetApiClient(ctx)

	result, httpRes, err := client.FirmwarePolicyAPI.
		ApplyFirmwarePoliciesWithoutServerInstanceGroups(ctx).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(result, nil)
}

func GetGlobalFirmwareConfiguration(ctx context.Context) error {
	logger.Get().Info().Msgf("Getting global firmware configuration")

	client := api.GetApiClient(ctx)

	config, httpRes, err := client.FirmwarePolicyAPI.
		GetGlobalFirmwareConfiguration(ctx).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(config, &globalFirmwareConfigPrintConfig)
}

func UpdateGlobalFirmwareConfiguration(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Updating global firmware configuration")

	var globalConfig sdk.UpdateGlobalFirmwareUpgradeConfiguration
	err := json.Unmarshal(config, &globalConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	result, httpRes, err := client.FirmwarePolicyAPI.
		UpdateGlobalFirmwareConfiguration(ctx).
		UpdateGlobalFirmwareUpgradeConfiguration(globalConfig).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(result, &globalFirmwareConfigPrintConfig)
}

func FirmwarePolicyConfigExample(ctx context.Context) error {
	// Example create firmware policy configuration
	firmwarePolicyConfiguration := sdk.CreateServerFirmwareUpgradePolicy{
		Label:       "example-firmware-policy",
		UserIdOwner: sdk.PtrFloat32(1),
		Action:      "upgrade",
		Rules: []sdk.ServerFirmwareUpgradePolicyRule{
			{
				Property:  "os",
				Operation: "equals",
				Value:     "ubuntu",
			},
		},
		ServerInstanceGroupIds: []float32{100, 101, 102},
	}

	return formatter.PrintResult(firmwarePolicyConfiguration, nil)
}

func GlobalFirmwareConfigExample(ctx context.Context) error {
	// Example global firmware configuration
	globalConfiguration := sdk.UpdateGlobalFirmwareUpgradeConfiguration{
		Activated:        sdk.PtrBool(true),
		UpgradeStartTime: sdk.PtrString("2023-10-01T00:00:00Z"),
		UpgradeEndTime:   sdk.PtrString("2023-10-01T23:59:59Z"),
	}

	return formatter.PrintResult(globalConfiguration, nil)
}

func getFirmwarePolicyId(firmwarePolicyId string) (float32, error) {
	firmwarePolicyIdNumeric, err := strconv.ParseFloat(firmwarePolicyId, 32)
	if err != nil {
		err := fmt.Errorf("invalid firmware policy ID: '%s'", firmwarePolicyId)
		logger.Get().Error().Err(err).Msg("")
		return 0, err
	}

	return float32(firmwarePolicyIdNumeric), nil
}
