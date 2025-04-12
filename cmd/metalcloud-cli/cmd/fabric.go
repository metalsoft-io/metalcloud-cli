package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/fabric"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	fabricFlags = struct {
		configSource string
	}{}

	fabricCmd = &cobra.Command{
		Use:     "fabric [command]",
		Aliases: []string{"fc", "fabrics"},
		Short:   "Fabric management",
		Long:    `Fabric management commands.`,
	}

	fabricListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all fabrics.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return fabric.FabricList(cmd.Context())
		},
	}

	fabricGetCmd = &cobra.Command{
		Use:          "get fabric_id",
		Aliases:      []string{"show"},
		Short:        "Get fabric info.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fabric.FabricGet(cmd.Context(), args[0])
		},
	}

	fabricConfigExampleCmd = &cobra.Command{
		Use:          "config-example fabric_type",
		Short:        "Returns example of fabric configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fabric.FabricConfigExample(cmd.Context(), args[0])
		},
	}

	fabricCreateCmd = &cobra.Command{
		Use:          "create site_id_or_label fabric_name fabric_type [fabric_description]",
		Aliases:      []string{"new"},
		Short:        "Create new fabric.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		Args:         cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			description := args[1]
			if len(args) > 3 {
				description = args[3]
			}

			config, err := utils.ReadConfigFromPipeOrFile(fabricFlags.configSource)
			if err != nil {
				return err
			}

			return fabric.FabricCreate(cmd.Context(), args[0], args[1], args[2], description, config)
		},
	}

	fabricUpdateCmd = &cobra.Command{
		Use:          "update fabric_id [fabric_name [fabric_description]]",
		Aliases:      []string{"edit"},
		Short:        "Update fabric configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		Args:         cobra.RangeArgs(1, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := ""
			if len(args) > 1 {
				name = args[1]
			}

			description := ""
			if len(args) > 2 {
				description = args[2]
			}

			config, err := utils.ReadConfigFromPipeOrFile(fabricFlags.configSource)
			if err != nil {
				return err
			}

			return fabric.FabricUpdate(cmd.Context(), args[0], name, description, config)
		},
	}

	fabricActivateCmd = &cobra.Command{
		Use:          "activate fabric_id",
		Aliases:      []string{"start"},
		Short:        "Activate a fabric.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fabric.FabricActivate(cmd.Context(), args[0])
		},
	}

	fabricDevicesGetCmd = &cobra.Command{
		Use:          "get-devices fabric_id",
		Aliases:      []string{"show-devices"},
		Short:        "List fabric devices.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fabric.FabricDevicesGet(cmd.Context(), args[0])
		},
	}

	fabricDevicesAddCmd = &cobra.Command{
		Use:          "add-device fabric_id device_id...",
		Aliases:      []string{"join-device"},
		Short:        "Add network device(s) to a fabric.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		Args:         cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fabric.FabricDevicesAdd(cmd.Context(), args[0], args[1:])
		},
	}

	fabricDevicesRemoveCmd = &cobra.Command{
		Use:          "remove-device fabric_id device_id",
		Aliases:      []string{"delete-device"},
		Short:        "Remove network device from a fabric.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_NETWORK_FABRICS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fabric.FabricDevicesRemove(cmd.Context(), args[0], args[1])
		},
	}
)

func init() {
	rootCmd.AddCommand(fabricCmd)

	fabricCmd.AddCommand(fabricListCmd)

	fabricCmd.AddCommand(fabricGetCmd)

	fabricCmd.AddCommand(fabricConfigExampleCmd)

	fabricCmd.AddCommand(fabricCreateCmd)
	fabricCreateCmd.Flags().StringVar(&fabricFlags.configSource, "config-source", "", "Source of the new fabric configuration. Can be 'pipe' or path to a JSON file.")
	fabricCreateCmd.MarkFlagsOneRequired("config-source")

	fabricCmd.AddCommand(fabricUpdateCmd)
	fabricUpdateCmd.Flags().StringVar(&fabricFlags.configSource, "config-source", "", "Source of the updated fabric configuration. Can be 'pipe' or path to a JSON file.")
	fabricUpdateCmd.MarkFlagsOneRequired("config-source")

	fabricCmd.AddCommand(fabricActivateCmd)

	fabricCmd.AddCommand(fabricDevicesGetCmd)
	fabricCmd.AddCommand(fabricDevicesAddCmd)
	fabricCmd.AddCommand(fabricDevicesRemoveCmd)
}
