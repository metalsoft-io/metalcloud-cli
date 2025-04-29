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
		Short:   "Logical network profile management",
		Long:    `Logical network profile management commands.`,
	}

	logicalNetworkProfileListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all logical network profiles.",
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
		Use:          "get logical_network_profile_id",
		Aliases:      []string{"show"},
		Short:        "Get logical network profile details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network_profile.LogicalNetworkProfileGet(cmd.Context(), args[0])
		},
	}

	logicalNetworkProfileConfigExampleCmd = &cobra.Command{
		Use:          "config-example kind",
		Aliases:      []string{"example"},
		Short:        "Get example of logical network profile configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network_profile.LogicalNetworkProfileConfigExample(cmd.Context(), args[0])
		},
	}

	logicalNetworkProfileCreateCmd = &cobra.Command{
		Use:          "create kind",
		Aliases:      []string{"new"},
		Short:        "Create a new logical network profile.",
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
		Use:          "update logical_network_profile_id",
		Aliases:      []string{"edit"},
		Short:        "Update a logical network profile.",
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
		Use:          "delete logical_network_profile_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a logical network profile.",
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
