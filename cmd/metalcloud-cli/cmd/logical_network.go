package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/logical_network"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	logicalNetworkFlags = struct {
		configSource           string
		filterId               []string
		filterLabel            []string
		filterFabricId         []string
		filterInfrastructureId []string
		filterKind             []string
		sortBy                 []string
	}{}

	logicalNetworkCmd = &cobra.Command{
		Use:     "logical-network [command]",
		Aliases: []string{"ln", "network", "logical_network"},
		Short:   "Logical network management",
		Long:    `Logical network management commands.`,
	}

	logicalNetworkListCmd = &cobra.Command{
		Use:          "list [fabric_id_or_label]",
		Aliases:      []string{"ls"},
		Short:        "List all logical networks, optionally filtered by fabric.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fabricIdOrLabel := ""
			if len(args) > 0 {
				fabricIdOrLabel = args[0]
			}

			return logical_network.LogicalNetworkList(cmd.Context(), fabricIdOrLabel, logical_network.ListFlags{
				FilterId:               logicalNetworkFlags.filterId,
				FilterLabel:            logicalNetworkFlags.filterLabel,
				FilterFabricId:         logicalNetworkFlags.filterFabricId,
				FilterInfrastructureId: logicalNetworkFlags.filterInfrastructureId,
				FilterKind:             logicalNetworkFlags.filterKind,
				SortBy:                 logicalNetworkFlags.sortBy,
			})
		},
	}

	logicalNetworkGetCmd = &cobra.Command{
		Use:          "get logical_network_id",
		Aliases:      []string{"show"},
		Short:        "Get logical network details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network.LogicalNetworkGet(cmd.Context(), args[0])
		},
	}

	logicalNetworkConfigExampleCmd = &cobra.Command{
		Use:          "config-example kind",
		Aliases:      []string{"example"},
		Short:        "Get example of logical network configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network.LogicalNetworkConfigExample(cmd.Context(), args[0])
		},
	}

	logicalNetworkCreateCmd = &cobra.Command{
		Use:          "create kind",
		Aliases:      []string{"new"},
		Short:        "Create a new logical network.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(logicalNetworkFlags.configSource)
			if err != nil {
				return err
			}

			return logical_network.LogicalNetworkCreate(cmd.Context(), args[0], config)
		},
	}

	logicalNetworkUpdateCmd = &cobra.Command{
		Use:          "update logical_network_id",
		Aliases:      []string{"edit"},
		Short:        "Update a logical network.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(logicalNetworkFlags.configSource)
			if err != nil {
				return err
			}

			return logical_network.LogicalNetworkUpdate(cmd.Context(), args[0], config)
		},
	}

	logicalNetworkDeleteCmd = &cobra.Command{
		Use:          "delete logical_network_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a logical network.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logical_network.LogicalNetworkDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(logicalNetworkCmd)

	logicalNetworkCmd.AddCommand(logicalNetworkListCmd)
	logicalNetworkListCmd.Flags().StringSliceVar(&logicalNetworkFlags.filterId, "filter-id", nil, "Filter by logical network ID.")
	logicalNetworkListCmd.Flags().StringSliceVar(&logicalNetworkFlags.filterLabel, "filter-label", nil, "Filter by logical network label.")
	logicalNetworkListCmd.Flags().StringSliceVar(&logicalNetworkFlags.filterFabricId, "filter-fabric-id", nil, "Filter by fabric ID.")
	logicalNetworkListCmd.Flags().StringSliceVar(&logicalNetworkFlags.filterInfrastructureId, "filter-infrastructure-id", nil, "Filter by infrastructure ID.")
	logicalNetworkListCmd.Flags().StringSliceVar(&logicalNetworkFlags.filterKind, "filter-kind", nil, "Filter by logical network kind.")
	logicalNetworkListCmd.Flags().StringSliceVar(&logicalNetworkFlags.sortBy, "sort-by", nil, "Sort by fields (e.g., id:ASC, name:DESC).")

	logicalNetworkCmd.AddCommand(logicalNetworkGetCmd)

	logicalNetworkCmd.AddCommand(logicalNetworkConfigExampleCmd)

	logicalNetworkCmd.AddCommand(logicalNetworkCreateCmd)
	logicalNetworkCreateCmd.Flags().StringVar(&logicalNetworkFlags.configSource, "config-source", "", "Source of the new logical network configuration. Can be 'pipe' or path to a JSON file.")
	logicalNetworkCreateCmd.MarkFlagsOneRequired("config-source")

	logicalNetworkCmd.AddCommand(logicalNetworkUpdateCmd)
	logicalNetworkUpdateCmd.Flags().StringVar(&logicalNetworkFlags.configSource, "config-source", "", "Source of the logical network updates. Can be 'pipe' or path to a JSON file.")
	logicalNetworkUpdateCmd.MarkFlagsOneRequired("config-source")

	logicalNetworkCmd.AddCommand(logicalNetworkDeleteCmd)
}
