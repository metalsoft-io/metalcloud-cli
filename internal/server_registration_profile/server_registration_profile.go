package server_registration_profile

import (
	"context"
	"fmt"
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/pkg/api"
	"github.com/metalsoft-io/metalcloud-cli/pkg/formatter"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
	"github.com/metalsoft-io/metalcloud-cli/pkg/response_inspector"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"

	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
)

var registrationProfilePrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"Id": {
			Title: "ID",
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

	request := client.ServerRegistrationProfileAPI.GetServerRegistrationProfiles(ctx).SortBy([]string{"id:ASC"})

	records, meta, err := utils.FetchAllPages(request)
	if err != nil {
		return err
	}

	return utils.PrintAll(records, meta, len(records), &registrationProfilePrintConfig)
}

func RegistrationProfileGet(ctx context.Context, registrationProfileId string) error {
	logger.Get().Info().Msgf("Get server registration profile '%s'", registrationProfileId)

	registrationProfileIdNumber, err := strconv.ParseInt(registrationProfileId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid server registration profile ID: '%s'", registrationProfileId)
		logger.Get().Error().Err(err).Msg("")
		return err
	}

	client := api.GetApiClient(ctx)

	registrationProfile, httpRes, err := client.ServerRegistrationProfileAPI.GetServerRegistrationProfileInfo(ctx, registrationProfileIdNumber).Execute()
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
		UpdateServerRegistrationProfile(ctx, registrationProfileIdNumber).
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
		DeleteServerRegistrationProfile(ctx, registrationProfileIdNumber).
		IfMatch(revision).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	logger.Get().Info().Msgf("Server registration profile '%s' deleted successfully", registrationProfileId)
	return nil
}

func getRegistrationProfileIdAndRevision(ctx context.Context, registrationProfileId string) (int64, string, error) {
	registrationProfileIdNumeric, err := strconv.ParseInt(registrationProfileId, 10, 64)
	if err != nil {
		err := fmt.Errorf("invalid registration profile ID: '%s'", registrationProfileId)
		logger.Get().Error().Err(err).Msg("")
		return 0, "", err
	}

	client := api.GetApiClient(ctx)

	registrationProfile, httpRes, err := client.ServerRegistrationProfileAPI.GetServerRegistrationProfileInfo(ctx, registrationProfileIdNumeric).Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return 0, "", err
	}

	return registrationProfileIdNumeric, registrationProfile.Revision, nil
}

// registrationProfileSettingsPrintConfig renders a bare settings object, as
// returned by the system-defaults endpoint (the profile print config addresses
// the same fields through the nested "Settings." prefix).
var registrationProfileSettingsPrintConfig = formatter.PrintConfig{
	FieldsConfig: map[string]formatter.RecordFieldConfig{
		"RegisterCredentials": {
			Title: "Register Credentials",
			Order: 1,
		},
		"MinimumNumberOfConnectedInterfaces": {
			Title: "Min Number Connected Interfaces",
			Order: 2,
		},
		"AlwaysDiscoverInterfacesWithBDK": {
			Title:       "Always Use BDK",
			Transformer: formatter.FormatBooleanValue,
			Order:       3,
		},
		"EnableTpm": {
			Title:       "Enable TPM",
			Transformer: formatter.FormatBooleanValue,
			Order:       4,
		},
		"EnableIntelTxt": {
			Title:       "Enable Intel Txt",
			Transformer: formatter.FormatBooleanValue,
			Order:       5,
		},
		"EnableSyslogMonitoring": {
			Title:       "Enable Syslog",
			Transformer: formatter.FormatBooleanValue,
			Order:       6,
		},
		"DisableTpmAfterRegistration": {
			Title:       "Disable TPM After Reg",
			Transformer: formatter.FormatBooleanValue,
			Order:       7,
		},
		"DefaultVirtualMediaProtocol": {
			Title: "Default Virt Media Protocol",
			Order: 8,
		},
		"ResetRaidControllers": {
			Title:       "Reset RAID Controllers",
			Transformer: formatter.FormatBooleanValue,
			Order:       9,
		},
		"CleanupDrives": {
			Title:       "Cleanup Drives",
			Transformer: formatter.FormatBooleanValue,
			Order:       10,
		},
		"RecreateRaid": {
			Title:       "Recreate RAID",
			Transformer: formatter.FormatBooleanValue,
			Order:       11,
		},
		"DisableEmbeddedNics": {
			Title:       "Disable Embedded NICs",
			Transformer: formatter.FormatBooleanValue,
			Order:       12,
		},
		"RaidOneDrive": {
			Title: "RAID One Drive",
			Order: 13,
		},
		"RaidTwoDrives": {
			Title: "RAID Two Drives",
			Order: 14,
		},
		"RaidEvenNumberMoreThanTwoDrives": {
			Title: "RAID Even 4+ Drives",
			Order: 15,
		},
		"RaidOddNumberMoreThanOneDrive": {
			Title: "RAID Odd 3+ Drives",
			Order: 16,
		},
		"DpuMode": {
			Title: "DPU Mode",
			Order: 17,
		},
	},
}

// RegistrationProfileSearch resolves the registration profile that applies to a
// site, falling back to the system-wide default when the site has none. The
// site ID is mandatory: the SDK rejects the call without it.
func RegistrationProfileSearch(ctx context.Context, siteId int) error {
	logger.Get().Info().Msgf("Searching for the server registration profile of site %d", siteId)

	client := api.GetApiClient(ctx)

	registrationProfile, httpRes, err := client.ServerRegistrationProfileAPI.
		SearchServerRegistrationProfileInfo(ctx).
		SiteId(float32(siteId)).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if registrationProfile == nil {
		logger.Get().Info().Msgf("No server registration profile matched the search")
		return nil
	}

	return formatter.PrintResult(registrationProfile, &registrationProfilePrintConfig)
}

// RegistrationProfileForServer shows the registration profile a server was, or
// would be, registered with.
func RegistrationProfileForServer(ctx context.Context, serverId string) error {
	logger.Get().Info().Msgf("Get the server registration profile of server '%s'", serverId)

	serverIdNumeric, err := utils.GetInt64FromString(serverId)
	if err != nil {
		return err
	}

	client := api.GetApiClient(ctx)

	registrationProfile, httpRes, err := client.ServerRegistrationProfileAPI.
		GetServerRegistrationProfileInfoForServer(ctx, serverIdNumeric).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if registrationProfile == nil {
		logger.Get().Info().Msgf("Server '%s' has no server registration profile", serverId)
		return nil
	}

	return formatter.PrintResult(registrationProfile, &registrationProfilePrintConfig)
}

// RegistrationProfileSystemDefaults shows the built-in registration settings
// used when no profile applies.
func RegistrationProfileSystemDefaults(ctx context.Context) error {
	logger.Get().Info().Msgf("Get the system default server registration settings")

	client := api.GetApiClient(ctx)

	settings, httpRes, err := client.ServerRegistrationProfileAPI.
		GetServerRegistrationProfileSystemDefaults(ctx).
		Execute()
	if err := response_inspector.InspectResponse(httpRes, err); err != nil {
		return err
	}

	if settings == nil {
		logger.Get().Info().Msgf("No system default server registration settings returned")
		return nil
	}

	return formatter.PrintResult(settings, &registrationProfileSettingsPrintConfig)
}
