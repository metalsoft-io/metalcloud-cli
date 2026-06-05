package quota_profile

import (
	"context"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var quotaProfilePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Name": {
			Title: "Name",
			Order: 2,
		},
		"Description": {
			Title: "Description",
			Order: 3,
		},
		"Limits": {
			Hidden: true,
		},
	},
}

func QuotaProfileList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing quota profiles")

	client := api.GetApiClient(ctx)

	quotaProfileList, httpRes, err := client.SecurityAPI.GetQuotaProfiles(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(quotaProfileList, &quotaProfilePrintConfig)
}

func QuotaProfileGet(ctx context.Context, profileId string) error {
	logger.Get().Info().Msgf("Get quota profile '%s'", profileId)

	client := api.GetApiClient(ctx)

	quotaProfile, httpRes, err := client.SecurityAPI.GetQuotaProfile(ctx, profileId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(quotaProfile, &quotaProfilePrintConfig)
}

func QuotaProfileCreate(ctx context.Context, config []byte) error {
	logger.Get().Info().Msgf("Creating quota profile")

	var createConfig sdk.CreateQuotaProfile
	err := utils.UnmarshalContent(config, &createConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	quotaProfile, httpRes, err := client.SecurityAPI.CreateQuotaProfile(ctx).CreateQuotaProfile(createConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(quotaProfile, &quotaProfilePrintConfig)
}

func QuotaProfileUpdate(ctx context.Context, profileId string, config []byte) error {
	logger.Get().Info().Msgf("Updating quota profile '%s'", profileId)

	var updateConfig sdk.EditQuotaProfile
	err := utils.UnmarshalContent(config, &updateConfig)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	quotaProfile, httpRes, err := client.SecurityAPI.UpdateQuotaProfile(ctx, profileId).EditQuotaProfile(updateConfig).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(quotaProfile, &quotaProfilePrintConfig)
}

func QuotaProfileDelete(ctx context.Context, profileId string) error {
	logger.Get().Info().Msgf("Deleting quota profile '%s'", profileId)

	client := api.GetApiClient(ctx)

	httpRes, err := client.SecurityAPI.DeleteQuotaProfile(ctx, profileId).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Quota profile '%s' deleted successfully", profileId)
	return nil
}

func QuotaProfileConfigExample(ctx context.Context) error {
	quotaProfileConfig := sdk.CreateQuotaProfile{
		Id:          sdk.PtrString("example-quota-profile"),
		Name:        "Example Quota Profile",
		Description: sdk.PtrString("Example quota profile for testing"),
		Limits: &sdk.PatchQuotaProfileLimitsDto{
			InfrastructureServerGroupMaxCount: sdk.PtrFloat32(10),
			InfrastructureDriveMaxCount:       sdk.PtrFloat32(20),
			ServerGroupInstancesMaxCount:      sdk.PtrFloat32(5),
			UserSshKeysCountMax:               sdk.PtrFloat32(10),
			AllowedServerTypes:                []string{},
			AllowedSites:                      []string{},
		},
	}

	return formatter.PrintResult(quotaProfileConfig, nil)
}
