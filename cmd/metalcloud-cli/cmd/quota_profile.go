package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/quota_profile"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

// Quota Profile commands
var (
	quotaProfileFlags = struct {
		configSource string
	}{}

	quotaProfileCmd = &cobra.Command{
		Use:     "quota-profile [command]",
		Aliases: []string{"qp", "quota-profiles"},
		Short:   "Quota Profile management",
		Long: `Quota Profile management commands.

This command group provides comprehensive quota profile management capabilities
including creation, retrieval, updating, and deletion of quota profiles. Quota
profiles define resource limits that can be applied to users.

Available commands:
  - Basic operations: list, get, create, update, delete
  - Configuration: config-example

Use "metalcloud-cli quota-profile [command] --help" for detailed information about each command.
`,
	}

	quotaProfileListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List quota profiles",
		Long: `List all quota profiles in the MetalSoft infrastructure.

This command displays information about all quota profiles including their IDs,
names, and descriptions.

Examples:
  # List all quota profiles
  metalcloud-cli quota-profile list
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_QUOTA_PROFILES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return quota_profile.QuotaProfileList(cmd.Context())
		},
	}

	quotaProfileGetCmd = &cobra.Command{
		Use:     "get profile_id",
		Aliases: []string{"show"},
		Short:   "Get detailed quota profile information",
		Long: `Get detailed information for a specific quota profile.

This command retrieves comprehensive information about a quota profile including
its configuration, limits, and other metadata.

Required Arguments:
  profile_id            The ID of the quota profile to retrieve information for

Examples:
  # Get quota profile information
  metalcloud-cli quota-profile get example-quota-profile

  # Get quota profile information using alias
  metalcloud-cli quota-profile show example-quota-profile
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_QUOTA_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return quota_profile.QuotaProfileGet(cmd.Context(), args[0])
		},
	}

	quotaProfileCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new quota profile",
		Long: `Create a new quota profile in MetalSoft.

The quota profile configuration must be provided using the --config-source flag.
The configuration source can be a path to a JSON file or 'pipe' to read from
standard input.

Required Flags:
  --config-source       Source of the quota profile configuration. Can be 'pipe' or path to a JSON file

Examples:
  # Create using JSON configuration file
  metalcloud-cli quota-profile create --config-source ./quota-profile.json

  # Create using piped JSON configuration
  cat quota-profile.json | metalcloud-cli quota-profile create --config-source pipe
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_QUOTA_PROFILES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(quotaProfileFlags.configSource)
			if err != nil {
				return err
			}
			return quota_profile.QuotaProfileCreate(cmd.Context(), config)
		},
	}

	quotaProfileUpdateCmd = &cobra.Command{
		Use:     "update profile_id",
		Aliases: []string{"edit"},
		Short:   "Update quota profile information",
		Long: `Update quota profile information.

This command updates quota profile configuration using a JSON configuration file
or piped JSON data. The configuration must be provided via the --config-source flag.

Required Arguments:
  profile_id            The ID of the quota profile to update

Required Flags:
  --config-source       Source of the quota profile update configuration. Can be 'pipe' or path to a JSON file

Examples:
  # Update quota profile using JSON configuration file
  metalcloud-cli quota-profile update example-quota-profile --config-source ./quota-profile-update.json

  # Update quota profile using piped JSON configuration
  echo '{"description": "Updated description"}' | metalcloud-cli quota-profile update example-quota-profile --config-source pipe
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_QUOTA_PROFILES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(quotaProfileFlags.configSource)
			if err != nil {
				return err
			}
			return quota_profile.QuotaProfileUpdate(cmd.Context(), args[0], config)
		},
	}

	quotaProfileDeleteCmd = &cobra.Command{
		Use:     "delete profile_id",
		Aliases: []string{"rm"},
		Short:   "Delete a quota profile",
		Long: `Delete a quota profile from MetalSoft infrastructure.

This command permanently deletes a quota profile. This action cannot be undone,
so use with caution.

Required Arguments:
  profile_id            The ID of the quota profile to delete

Examples:
  # Delete a quota profile
  metalcloud-cli quota-profile delete example-quota-profile

  # Delete a quota profile using alias
  metalcloud-cli quota-profile rm example-quota-profile
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_QUOTA_PROFILES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return quota_profile.QuotaProfileDelete(cmd.Context(), args[0])
		},
	}

	quotaProfileConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Show quota profile configuration example",
		Long: `Show an example quota profile configuration.

This command outputs an example quota profile configuration that can be used as a
starting point for creating or updating quota profiles.

Examples:
  # Show the configuration example
  metalcloud-cli quota-profile config-example
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_QUOTA_PROFILES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return quota_profile.QuotaProfileConfigExample(cmd.Context())
		},
	}
)

func init() {
	rootCmd.AddCommand(quotaProfileCmd)

	quotaProfileCmd.AddCommand(quotaProfileListCmd)

	quotaProfileCmd.AddCommand(quotaProfileGetCmd)

	quotaProfileCmd.AddCommand(quotaProfileCreateCmd)
	quotaProfileCreateCmd.Flags().StringVar(&quotaProfileFlags.configSource, "config-source", "", "Source of the new quota profile configuration. Can be 'pipe' or path to a JSON file.")
	quotaProfileCreateCmd.MarkFlagsOneRequired("config-source")

	quotaProfileCmd.AddCommand(quotaProfileUpdateCmd)
	quotaProfileUpdateCmd.Flags().StringVar(&quotaProfileFlags.configSource, "config-source", "", "Source of the quota profile update configuration. Can be 'pipe' or path to a JSON file.")
	quotaProfileUpdateCmd.MarkFlagsOneRequired("config-source")

	quotaProfileCmd.AddCommand(quotaProfileDeleteCmd)

	quotaProfileCmd.AddCommand(quotaProfileConfigExampleCmd)
}
