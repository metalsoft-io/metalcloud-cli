package server_registration_profile

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"

	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var registrationProfilePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "#",
			Order: 1,
		},
		"Name": {
			Title:    "Name",
			MaxWidth: 30,
			Order:    2,
		},
		"Settings.RegisterCredentials": {
			Title: "Register Credentials",
			Order: 3,
		},
		"Settings.MinimumNumberOfConnectedInterfaces": {
			Title: "Min Number Connected Interfaces",
			Order: 4,
		},
		"Settings.AlwaysDiscoverInterfacesWithBDK": {
			Title:       "Always Use BDK",
			Transformer: formatter.FormatBooleanValue,
			Order:       5,
		},
		"Settings.EnableTpm": {
			Title:       "Enable TPM",
			Transformer: formatter.FormatBooleanValue,
			Order:       6,
		},
		"Settings.EnableIntelTxt": {
			Title:       "Enable Intel Txt",
			Transformer: formatter.FormatBooleanValue,
			Order:       7,
		},
		"Settings.EnableSyslogMonitoring": {
			Title:       "Enable Syslog",
			Transformer: formatter.FormatBooleanValue,
			Order:       8,
		},
		"Settings.DisableTpmAfterRegistration": {
			Title:       "Disable TPM After Reg",
			Transformer: formatter.FormatBooleanValue,
			Order:       9,
		},
		"Settings.DefaultVirtualMediaProtocol": {
			Title: "Default Virt Media Protocol",
			Order: 10,
		},
		"Settings.ResetRaidControllers": {
			Title:       "Reset RAID Controllers",
			Transformer: formatter.FormatBooleanValue,
			Order:       11,
		},
		"Settings.CleanupDrives": {
			Title:       "Cleanup Drives",
			Transformer: formatter.FormatBooleanValue,
			Order:       12,
		},
		"Settings.RecreateRaid": {
			Title:       "Recreate RAID",
			Transformer: formatter.FormatBooleanValue,
			Order:       13,
		},
		"Settings.DisableEmbeddedNics": {
			Title:       "Disable Embedded NICs",
			Transformer: formatter.FormatBooleanValue,
			Order:       14,
		},
		"Settings.RaidOneDrive": {
			Title: "RAID One Drive",
			Order: 15,
		},
		"Settings.RaidTwoDrives": {
			Title: "RAID Two Drives",
			Order: 16,
		},
		"Settings.RaidEvenNumberMoreThanTwoDrives": {
			Title: "RAID Even 4+ Drives",
			Order: 17,
		},
		"Settings.RaidOddNumberMoreThanOneDrive": {
			Title: "RAID Odd 3+ Drives",
			Order: 18,
		},
	},
}

func RegistrationProfileList(ctx context.Context) error {
	logger.Get().Info().Msgf("Listing all server cleanup policies")

	client := api.GetApiClient(ctx)

	registrationProfilesList, httpRes, err := client.ServerRegistrationProfileAPI.GetServerRegistrationProfiles(ctx).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(registrationProfilesList, &registrationProfilePrintConfig)
}

func RegistrationProfileGet(ctx context.Context, registrationProfileId string) error {
	logger.Get().Info().Msgf("Get server registration profile '%s'", registrationProfileId)

	registrationProfileIdNumber, err := strconv.ParseFloat(registrationProfileId, 32)
	if err != nil {
		err := fmt.Errorf("invalid server registration profile ID: '%s'", registrationProfileId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	registrationProfile, httpRes, err := client.ServerRegistrationProfileAPI.GetServerRegistrationProfileInfo(ctx, float32(registrationProfileIdNumber)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(registrationProfile, &registrationProfilePrintConfig)
}

func RegistrationProfileCreate(ctx context.Context, name string, settings sdk.ServerRegistrationProfileSettings) error {
	logger.Get().Info().Msgf("Creating server registration profile")

	client := api.GetApiClient(ctx)

	createRequest := sdk.ServerRegistrationProfileCreate{
		Name:     name,
		Settings: settings,
	}

	registrationProfile, httpRes, err := client.ServerRegistrationProfileAPI.
		CreateServerRegistrationProfile(ctx).
		ServerRegistrationProfileCreate(createRequest).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(registrationProfile, &registrationProfilePrintConfig)
}

func RegistrationProfileUpdate(ctx context.Context, registrationProfileId string, name string, settings sdk.ServerRegistrationProfileUpdateSettings) error {
	logger.Get().Info().Msgf("Updating server registration profile '%s'", registrationProfileId)

	registrationProfileIdNumber, revision, err := getRegistrationProfileIdAndRevision(ctx, registrationProfileId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	updateRequest := sdk.ServerRegistrationProfileUpdate{}

	if len(name) > 0 {
		updateRequest.Name = &name
	}

	updateRequest.Settings = &settings

	registrationProfile, httpRes, err := client.ServerRegistrationProfileAPI.
		UpdateServerRegistrationProfile(ctx, float32(registrationProfileIdNumber)).
		ServerRegistrationProfileUpdate(updateRequest).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	return formatter.PrintResult(registrationProfile, &registrationProfilePrintConfig)
}

func RegistrationProfileDelete(ctx context.Context, registrationProfileId string) error {
	logger.Get().Info().Msgf("Deleting server registration profile '%s'", registrationProfileId)

	registrationProfileIdNumber, revision, err := getRegistrationProfileIdAndRevision(ctx, registrationProfileId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	httpRes, err := client.ServerRegistrationProfileAPI.
		DeleteServerRegistrationProfile(ctx, float32(registrationProfileIdNumber)).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Server registration profile '%s' deleted successfully", registrationProfileId)
	return nil
}

func getRegistrationProfileIdAndRevision(ctx context.Context, registrationProfileId string) (float32, string, error) {
	registrationProfileIdNumeric, err := strconv.ParseFloat(registrationProfileId, 32)
	if err != nil {
		err := fmt.Errorf("invalid registration profile ID: '%s'", registrationProfileId)
		logger.Get().Error().Err(err).Msg("")
		return 0, "", err
	}

	client := api.GetApiClient(ctx)

	registrationProfile, httpRes, err := client.ServerRegistrationProfileAPI.GetServerRegistrationProfileInfo(ctx, float32(registrationProfileIdNumeric)).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}

	return float32(registrationProfileIdNumeric), registrationProfile.Revision, nil
}
