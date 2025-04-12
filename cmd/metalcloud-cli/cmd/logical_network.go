package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/logical_network"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	logicalNetworkFlags = struct {
		configSource string
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
			return logical_network.LogicalNetworkList(cmd.Context(), fabricIdOrLabel)
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

	logicalNetworkCreateCmd = &cobra.Command{
		Use:          "create",
		Aliases:      []string{"new"},
		Short:        "Create a new logical network.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_PROFILES_READ},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(logicalNetworkFlags.configSource)
			if err != nil {
				return err
			}

			return logical_network.LogicalNetworkCreate(cmd.Context(), config)
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
)

func init() {
	rootCmd.AddCommand(logicalNetworkCmd)

	logicalNetworkCmd.AddCommand(logicalNetworkListCmd)
	logicalNetworkCmd.AddCommand(logicalNetworkGetCmd)

	logicalNetworkCmd.AddCommand(logicalNetworkCreateCmd)
	logicalNetworkCreateCmd.Flags().StringVar(&logicalNetworkFlags.configSource, "config-source", "", "Source of the new logical network configuration. Can be 'pipe' or path to a JSON file.")
	logicalNetworkCreateCmd.MarkFlagsOneRequired("config-source")

	logicalNetworkCmd.AddCommand(logicalNetworkDeleteCmd)

	logicalNetworkCmd.AddCommand(logicalNetworkUpdateCmd)
	logicalNetworkUpdateCmd.Flags().StringVar(&logicalNetworkFlags.configSource, "config-source", "", "Source of the logical network updates. Can be 'pipe' or path to a JSON file.")
	logicalNetworkUpdateCmd.MarkFlagsOneRequired("config-source")
}
