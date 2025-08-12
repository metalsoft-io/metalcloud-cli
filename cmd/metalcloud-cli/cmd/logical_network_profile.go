package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/logical_network_profile"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	logicalNetworkProfileFlags = struct {
		configSource   string
		filterId       []string
		filterLabel    []string
		filterKind     []string
		filterName     []string
		filterFabricId []string
		sortBy         []string
		kind           string
	}{}

	logicalNetworkProfileCmd = &cobra.Command{
		Use:     "logical-network-profile [command]",
		Aliases: []string{"lnp", "network-profile"},
		Short:   "Manage logical network profiles for network configuration templates",
		Long: `Manage logical network profiles which define network configuration templates
that can be applied to infrastructure deployments. These profiles contain
network settings, routing rules, and connectivity configurations.

Available commands:
  list          List all logical network profiles with filtering options
  get           Get detailed information about a specific profile
  create        Create a new logical network profile from configuration
  update        Update an existing logical network profile
  delete        Delete a logical network profile
  config-example Get example configuration for a specific profile kind

Examples:
  metalcloud-cli logical-network-profile list
  metalcloud-cli lnp get 12345
  metalcloud-cli network-profile create --config-source profile.json`,
	}

	logicalNetworkProfileListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List logical network profiles with optional filtering",
		Long: `List all logical network profiles in the system with optional filtering capabilities.

This command displays a tabular view of logical network profiles including their
ID, name, kind, label, fabric ID, and other metadata. Use the filter flags to
narrow down results based on specific criteria.

Flags:
  --filter-id          Filter profiles by one or more profile IDs
  --filter-label       Filter profiles by label pattern matching
  --filter-kind        Filter profiles by profile kind (e.g., 'cisco', 'juniper')
  --filter-name        Filter profiles by name pattern matching
  --filter-fabric-id   Filter profiles by associated fabric ID
  --sort-by            Sort results by specified fields with direction

Examples:
  # List all logical network profiles
  metalcloud-cli logical-network-profile list

  # Filter by profile kind
  metalcloud-cli lnp list --filter-kind cisco

  # Filter by multiple criteria
  metalcloud-cli network-profile list --filter-kind cisco --filter-label "prod"

  # Sort by name in descending order
  metalcloud-cli lnp ls --sort-by name:DESC

  # Filter by fabric ID and sort by ID
  metalcloud-cli lnp list --filter-fabric-id 100 --sort-by id:ASC`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network_profile.LogicalNetworkProfileList(cmd.Context(), logical_network_profile.ListFlags{
				FilterId:       logicalNetworkProfileFlags.filterId,
				FilterLabel:    logicalNetworkProfileFlags.filterLabel,
				FilterKind:     logicalNetworkProfileFlags.filterKind,
				FilterName:     logicalNetworkProfileFlags.filterName,
				FilterFabricId: logicalNetworkProfileFlags.filterFabricId,
				SortBy:         logicalNetworkProfileFlags.sortBy,
			})
		},
	}

	logicalNetworkProfileGetCmd = &cobra.Command{
		Use:     "get logical_network_profile_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a logical network profile",
		Long: `Get comprehensive details about a specific logical network profile.

This command displays detailed information about a logical network profile including
its configuration, metadata, associated fabric information, and deployment settings.
The profile ID can be obtained from the list command.

Required Arguments:
  logical_network_profile_id    The unique identifier of the profile to retrieve

Examples:
  # Get details for profile ID 12345
  metalcloud-cli logical-network-profile get 12345

  # Get profile details using short alias
  metalcloud-cli lnp show 12345

  # Get profile details for piping to other commands
  metalcloud-cli network-profile get 12345 --output json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network_profile.LogicalNetworkProfileGet(cmd.Context(), args[0])
		},
	}

	logicalNetworkProfileConfigExampleCmd = &cobra.Command{
		Use:     "config-example kind",
		Aliases: []string{"example"},
		Short:   "Get example configuration for a specific profile kind",
		Long: `Get an example configuration template for creating logical network profiles.

This command displays a sample configuration in JSON format that can be used
as a starting point for creating new logical network profiles. The configuration
includes all required and optional fields with example values.

Required Arguments:
  kind    The profile kind to generate an example for (e.g., 'cisco', 'juniper', 'arista')

Examples:
  # Get example configuration for Cisco profiles
  metalcloud-cli logical-network-profile config-example cisco

  # Get example and save to file
  metalcloud-cli lnp example juniper > profile-template.json

  # Get example using alias
  metalcloud-cli network-profile example arista`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network_profile.LogicalNetworkProfileConfigExample(cmd.Context(), args[0])
		},
	}

	logicalNetworkProfileCreateCmd = &cobra.Command{
		Use:     "create kind",
		Aliases: []string{"new"},
		Short:   "Create a new logical network profile from configuration",
		Long: `Create a new logical network profile from a JSON configuration.

This command creates a new logical network profile based on the specified kind
and configuration provided via file or stdin. The configuration must match the
schema for the specified profile kind.

Required Arguments:
  kind    The type of profile to create (e.g., 'cisco', 'juniper', 'arista')

Required Flags:
  --config-source    Source of configuration data (required)
                     - 'pipe' to read from stdin
                     - path to JSON file containing profile configuration

Examples:
  # Create profile from JSON file
  metalcloud-cli logical-network-profile create cisco --config-source profile.json

  # Create profile from stdin
  cat profile.json | metalcloud-cli lnp create juniper --config-source pipe

  # Create profile using alias
  metalcloud-cli network-profile new arista --config-source ./configs/arista-profile.json

  # Get example configuration first, then create
  metalcloud-cli lnp example cisco > cisco-profile.json
  # Edit cisco-profile.json with your settings
  metalcloud-cli lnp create cisco --config-source cisco-profile.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			kind := args[0]
			config, err := utils.ReadConfigFromPipeOrFile(logicalNetworkProfileFlags.configSource)
			if err != nil {
				return err
			}
			return logical_network_profile.LogicalNetworkProfileCreate(cmd.Context(), config, kind)
		},
	}

	logicalNetworkProfileUpdateCmd = &cobra.Command{
		Use:     "update logical_network_profile_id",
		Aliases: []string{"edit"},
		Short:   "Update an existing logical network profile",
		Long: `Update an existing logical network profile with new configuration.

This command modifies an existing logical network profile by applying the provided
configuration changes. The profile ID must be valid and the configuration must
match the schema for the profile's kind.

Required Arguments:
  logical_network_profile_id    The unique identifier of the profile to update

Required Flags:
  --config-source    Source of configuration data (required)
                     - 'pipe' to read from stdin
                     - path to JSON file containing updated profile configuration

Examples:
  # Update profile from JSON file
  metalcloud-cli logical-network-profile update 12345 --config-source updated-profile.json

  # Update profile from stdin
  cat updated-config.json | metalcloud-cli lnp update 12345 --config-source pipe

  # Update profile using alias
  metalcloud-cli network-profile edit 12345 --config-source ./configs/new-settings.json

  # Get current profile, modify, then update
  metalcloud-cli lnp get 12345 --output json > current-profile.json
  # Edit current-profile.json with your changes
  metalcloud-cli lnp update 12345 --config-source current-profile.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(logicalNetworkProfileFlags.configSource)
			if err != nil {
				return err
			}
			return logical_network_profile.LogicalNetworkProfileUpdate(cmd.Context(), args[0], config)
		},
	}

	logicalNetworkProfileDeleteCmd = &cobra.Command{
		Use:     "delete logical_network_profile_id",
		Aliases: []string{"rm"},
		Short:   "Delete a logical network profile",
		Long: `Delete a logical network profile from the system.

This command permanently removes a logical network profile and all its associated
configuration. The profile must not be in use by any active deployments before
it can be deleted.

Required Arguments:
  logical_network_profile_id    The unique identifier of the profile to delete

Examples:
  # Delete profile by ID
  metalcloud-cli logical-network-profile delete 12345

  # Delete profile using short alias
  metalcloud-cli lnp rm 12345

  # Delete profile using alias
  metalcloud-cli network-profile delete 12345

Warning: This operation is irreversible. Ensure the profile is not in use
by any active infrastructure deployments before deletion.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network_profile.LogicalNetworkProfileDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(logicalNetworkProfileCmd)

	logicalNetworkProfileCmd.AddCommand(logicalNetworkProfileListCmd)
	logicalNetworkProfileListCmd.Flags().StringSliceVar(&logicalNetworkProfileFlags.filterId, "filter-id", nil, "Filter by profile ID.")
	logicalNetworkProfileListCmd.Flags().StringSliceVar(&logicalNetworkProfileFlags.filterLabel, "filter-label", nil, "Filter by profile label.")
	logicalNetworkProfileListCmd.Flags().StringSliceVar(&logicalNetworkProfileFlags.filterKind, "filter-kind", nil, "Filter by profile kind.")
	logicalNetworkProfileListCmd.Flags().StringSliceVar(&logicalNetworkProfileFlags.filterName, "filter-name", nil, "Filter by profile name.")
	logicalNetworkProfileListCmd.Flags().StringSliceVar(&logicalNetworkProfileFlags.filterFabricId, "filter-fabric-id", nil, "Filter by fabric ID.")
	logicalNetworkProfileListCmd.Flags().StringSliceVar(&logicalNetworkProfileFlags.sortBy, "sort-by", nil, "Sort by fields (e.g., id:ASC, name:DESC).")

	logicalNetworkProfileCmd.AddCommand(logicalNetworkProfileGetCmd)

	logicalNetworkProfileCmd.AddCommand(logicalNetworkProfileConfigExampleCmd)

	logicalNetworkProfileCmd.AddCommand(logicalNetworkProfileCreateCmd)
	logicalNetworkProfileCreateCmd.Flags().StringVar(&logicalNetworkProfileFlags.configSource, "config-source", "", "Source of the new logical network profile configuration. Can be 'pipe' or path to a JSON file.")
	logicalNetworkProfileCreateCmd.MarkFlagsOneRequired("config-source")

	logicalNetworkProfileCmd.AddCommand(logicalNetworkProfileUpdateCmd)
	logicalNetworkProfileUpdateCmd.Flags().StringVar(&logicalNetworkProfileFlags.configSource, "config-source", "", "Source of the logical network profile updates. Can be 'pipe' or path to a JSON file.")
	logicalNetworkProfileUpdateCmd.MarkFlagsOneRequired("config-source")

	logicalNetworkProfileCmd.AddCommand(logicalNetworkProfileDeleteCmd)
}
