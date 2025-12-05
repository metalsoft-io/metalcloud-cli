package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/server_registration_profile"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
)

// Server Registration Profile commands
var (
	serverRegistrationProfileFlags = struct {
		name         string
		configSource string
	}{}

	serverRegistrationProfileCmd = &cobra.Command{
		Use:     "server-registration-profile [command]",
		Aliases: []string{"srv-rp", "srp"},
		Short:   "Manage server registration profiles",
		Long: `Manage server registration profiles that determine how server registration is performed.

Server registration profiles control how the servers are registered.

Available Commands:
  list      List all server registration profiles
  get       Get details of a specific server registration profile
  create    Create a new server registration profile
  update    Update an existing server registration profile
  delete    Delete a server registration profile

Examples:
  # List all server registration profiles
  metalcloud-cli server-registration-profile list

  # Get details of a specific registration profile
  metalcloud-cli server-registration-profile get srp-123

  # Create a new server registration profile
  metalcloud-cli server-registration-profile create --name "my-profile" --config-source my-profile.json

  # Update an existing registration profile
  metalcloud-cli server-registration-profile update 123 --name "updated-profile"

  # Delete a registration profile
  metalcloud-cli server-registration-profile delete 123

  # Using short aliases
  metalcloud-cli srp list
  metalcloud-cli srv-rp get profile-123`,
	}

	serverRegistrationProfileListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all server registration profiles",
		Long: `List all server registration profiles configured in the system.

This command displays a table of all available server registration profiles with their
key attributes including ID, name, status, and configuration summary.

Output Format:
  By default, output is formatted as a table. Use global flags to change output format:
  --format=json    JSON output
  --format=csv     CSV output
  --format=yaml    YAML output

Examples:
  # List all server registration profiles in table format
  metalcloud-cli server-registration-profile list

  # List policies in JSON format
  metalcloud-cli server-registration-profile list --format=json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_TYPES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_registration_profile.RegistrationProfileList(cmd.Context())
		},
	}

	serverRegistrationProfileGetCmd = &cobra.Command{
		Use:     "get <profile-id>",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific server registration profile",
		Long: `Get detailed information about a specific server registration profile by its ID.

This command retrieves and displays comprehensive information about a server registration
profile including all configuration settings:

Configuration Details:
  - registerCredentials: Password handling mode (user/random)
  - minimumNumberOfConnectedInterfaces: Minimum required connected interfaces
  - alwaysDiscoverInterfacesWithBDK: BDK interface discovery setting
  - enableTpm: TPM enablement
  - enableIntelTxt: Intel TXT enablement
  - enableSyslogMonitoring: Syslog monitoring configuration
  - disableTpmAfterRegistration: Post-registration TPM status
  - biosProfile: Array of BIOS configuration settings
  - defaultVirtualMediaProtocol: Virtual media protocol (HTTPS, etc)
  - resetRaidControllers: RAID controller reset policy
  - cleanupDrives: Drive cleanup policy
  - recreateRaid: RAID recreation policy
  - disableEmbeddedNics: Embedded NIC configuration
  - raidOneDrive: RAID level for single drive (RAID0)
  - raidTwoDrives: RAID level for two drives (RAID1)
  - raidEvenNumberMoreThanTwoDrives: RAID level for even drives 4+ (RAID10)
  - raidOddNumberMoreThanOneDrive: RAID level for odd drives 3+ (RAID5)

Output Format:
  Use global flags to change output format:
  --format=json    JSON output
  --format=csv     CSV output
  --format=yaml    YAML output

Examples:
  # Get profile details by ID
  metalcloud-cli server-registration-profile get 123

  # Get profile in JSON format
  metalcloud-cli server-registration-profile get 123 --format=json

  # Using alias
  metalcloud-cli srp show 123`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_CLEANUP_POLICIES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_registration_profile.RegistrationProfileGet(cmd.Context(), args[0])
		},
	}

	serverRegistrationProfileCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new server registration profile",
		Long: `Create a new server registration profile with specified configuration.

The configuration must be provided via a JSON file or pipe and should include
settings for various aspects of server registration.

Required Flags:
  --name            Name for the server registration profile
  --config-source   Source of configuration (path to JSON file or 'pipe' for stdin)

Configuration Structure (JSON):
{
  "registerCredentials": "user|random",
  "minimumNumberOfConnectedInterfaces": 0,
  "alwaysDiscoverInterfacesWithBDK": true,
  "enableTpm": true,
  "enableIntelTxt": true,
  "enableSyslogMonitoring": true,
  "disableTpmAfterRegistration": false,
  "biosProfile": [
    {
      "key": "BootMode",
      "value": "Uefi"
    }
  ],
  "defaultVirtualMediaProtocol": "HTTPS",
  "resetRaidControllers": true,
  "cleanupDrives": true,
  "recreateRaid": true,
  "disableEmbeddedNics": true,
  "raidOneDrive": "RAID0",
  "raidTwoDrives": "RAID1",
  "raidEvenNumberMoreThanTwoDrives": "RAID10",
  "raidOddNumberMoreThanOneDrive": "RAID5"
}

Configuration Options:
  - registerCredentials: "user" (keep existing) or "random" (generate new)
  - minimumNumberOfConnectedInterfaces: Minimum required connected interfaces (default: 0)
  - alwaysDiscoverInterfacesWithBDK: Always use BDK for discovery (default: true)
  - enableTpm: Enable TPM during registration (default: true)
  - enableIntelTxt: Enable Intel TXT (default: true)
  - enableSyslogMonitoring: Enable syslog monitoring (default: true)
  - disableTpmAfterRegistration: Disable TPM after registration (default: false)
  - biosProfile: Array of BIOS settings (key-value pairs)
  - defaultVirtualMediaProtocol: HTTPS, NFS, CIFS, etc (default: HTTPS)
  - resetRaidControllers: Reset RAID controllers to factory defaults (default: true)
  - cleanupDrives: Clean up drives during registration (default: true)
  - recreateRaid: Recreate RAID configuration (default: true)
  - disableEmbeddedNics: Disable embedded NICs (default: true)
  - raidOneDrive: RAID level for 1 drive (default: RAID0)
  - raidTwoDrives: RAID level for 2 drives (default: RAID1)
  - raidEvenNumberMoreThanTwoDrives: RAID level for even 4+ drives (default: RAID10)
  - raidOddNumberMoreThanOneDrive: RAID level for odd 3+ drives (default: RAID5)

Examples:
  # Create profile from JSON file
  metalcloud-cli server-registration-profile create \
    --name "production-profile" \
    --config-source /path/to/config.json

  # Create profile from stdin
  cat config.json | metalcloud-cli server-registration-profile create \
    --name "staging-profile" \
    --config-source pipe

  # Create minimal profile
  echo '{"registerCredentials": "random"}' | metalcloud-cli srp create \
    --name "minimal-profile" \
    --config-source pipe

  # Create profile with custom RAID settings
  cat > raid-config.json << EOF
{
  "registerCredentials": "user",
  "resetRaidControllers": true,
  "cleanupDrives": true,
  "recreateRaid": true,
  "raidOneDrive": "RAID0",
  "raidTwoDrives": "RAID1",
  "raidEvenNumberMoreThanTwoDrives": "RAID6",
  "raidOddNumberMoreThanOneDrive": "RAID5"
}
EOF
  metalcloud-cli srp create --name "raid-profile" --config-source raid-config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_CLEANUP_POLICIES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			var settings sdk.ServerRegistrationProfileSettings

			if serverRegistrationProfileFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(serverRegistrationProfileFlags.configSource)
				if err != nil {
					return err
				}

				err = utils.UnmarshalContent(config, &settings)
				if err != nil {
					return err
				}
			}

			return server_registration_profile.RegistrationProfileCreate(cmd.Context(), serverRegistrationProfileFlags.name, settings)
		},
	}

	serverRegistrationProfileUpdateCmd = &cobra.Command{
		Use:   "update <profile-id>",
		Short: "Update an existing server registration profile",
		Long: `Update an existing server registration profile by its ID.

You can update the profile name, settings, or both. At least one of --name or
--config-source must be provided.

Required Arguments:
  <profile-id>      ID of the server registration profile to update

Optional Flags:
  --name            New name for the server registration profile
  --config-source   Source of configuration updates (path to JSON file or 'pipe')

Note: At least one of --name or --config-source must be specified.

Configuration Structure (JSON):
{
  "registerCredentials": "user|random",
  "minimumNumberOfConnectedInterfaces": 0,
  "alwaysDiscoverInterfacesWithBDK": true,
  "enableTpm": true,
  "enableIntelTxt": true,
  "enableSyslogMonitoring": true,
  "disableTpmAfterRegistration": false,
  "biosProfile": [
    {
      "key": "BootMode",
      "value": "Uefi"
    }
  ],
  "defaultVirtualMediaProtocol": "HTTPS",
  "resetRaidControllers": true,
  "cleanupDrives": true,
  "recreateRaid": true,
  "disableEmbeddedNics": true,
  "raidOneDrive": "RAID0",
  "raidTwoDrives": "RAID1",
  "raidEvenNumberMoreThanTwoDrives": "RAID10",
  "raidOddNumberMoreThanOneDrive": "RAID5"
}

Update Behavior:
  - Only specified fields in the configuration will be updated
  - Unspecified fields will retain their current values
  - Use --name to update only the profile name
  - Use --config-source to update only the settings
  - Use both to update name and settings simultaneously

Examples:
  # Update profile name only
  metalcloud-cli server-registration-profile update 123 \
    --name "updated-profile-name"

  # Update profile settings from file
  metalcloud-cli server-registration-profile update 123 \
    --config-source /path/to/updates.json

  # Update both name and settings
  metalcloud-cli server-registration-profile update 123 \
    --name "new-name" \
    --config-source updates.json

  # Update from stdin
  cat updates.json | metalcloud-cli srp update 123 --config-source pipe

  # Update specific settings only
  echo '{"enableTpm": false, "enableIntelTxt": false}' | \
    metalcloud-cli srp update 123 --config-source pipe

  # Disable drive cleanup and RAID recreation
  cat > updates.json << EOF
{
  "cleanupDrives": false,
  "recreateRaid": false
}
EOF
  metalcloud-cli srp update 123 --config-source updates.json

  # Change RAID configuration strategy
  echo '{"raidTwoDrives": "RAID0", "raidEvenNumberMoreThanTwoDrives": "RAID5"}' | \
    metalcloud-cli server-registration-profile update 123 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_TYPES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var settings sdk.ServerRegistrationProfileUpdateSettings

			if serverRegistrationProfileFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(serverRegistrationProfileFlags.configSource)
				if err != nil {
					return err
				}

				err = utils.UnmarshalContent(config, &settings)
				if err != nil {
					return err
				}
			}

			return server_registration_profile.RegistrationProfileUpdate(cmd.Context(), args[0], serverRegistrationProfileFlags.name, settings)
		},
	}

	serverRegistrationProfileDeleteCmd = &cobra.Command{
		Use:     "delete <profile-id>",
		Aliases: []string{"rm", "remove"},
		Short:   "Delete a server registration profile",
		Long: `Delete a server registration profile by its ID.

This command permanently removes a server registration profile from the system.
Once deleted, the profile cannot be recovered.

Required Arguments:
  <profile-id>      ID of the server registration profile to delete

Warning:
  - This operation is irreversible
  - Ensure the profile is not currently in use by any servers
  - Servers using this profile may fail to register properly after deletion

Examples:
  # Delete a server registration profile
  metalcloud-cli server-registration-profile delete 123

  # Using aliases
  metalcloud-cli srp rm 123
  metalcloud-cli srv-rp remove 123

  # Delete with confirmation in script
  PROFILE_ID=123
  metalcloud-cli server-registration-profile delete $PROFILE_ID`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_CLEANUP_POLICIES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_registration_profile.RegistrationProfileDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(serverRegistrationProfileCmd)

	serverRegistrationProfileCmd.AddCommand(serverRegistrationProfileListCmd)

	serverRegistrationProfileCmd.AddCommand(serverRegistrationProfileGetCmd)

	serverRegistrationProfileCmd.AddCommand(serverRegistrationProfileCreateCmd)
	serverRegistrationProfileCreateCmd.Flags().StringVar(&serverRegistrationProfileFlags.name, "name", "", "Name for the server registration profile")
	serverRegistrationProfileCreateCmd.Flags().StringVar(&serverRegistrationProfileFlags.configSource, "config-source", "", "Source of the new registration profile configuration. Can be 'pipe' or path to a JSON file.")
	serverRegistrationProfileCreateCmd.MarkFlagRequired("name")
	serverRegistrationProfileCreateCmd.MarkFlagRequired("config-source")

	serverRegistrationProfileCmd.AddCommand(serverRegistrationProfileUpdateCmd)
	serverRegistrationProfileUpdateCmd.Flags().StringVar(&serverRegistrationProfileFlags.name, "name", "", "New name for the server registration profile")
	serverRegistrationProfileUpdateCmd.Flags().StringVar(&serverRegistrationProfileFlags.configSource, "config-source", "", "Source of the registration profile configuration updates. Can be 'pipe' or path to a JSON file.")
	serverRegistrationProfileUpdateCmd.MarkFlagsOneRequired("name", "config-source")

	serverRegistrationProfileCmd.AddCommand(serverRegistrationProfileDeleteCmd)
}
