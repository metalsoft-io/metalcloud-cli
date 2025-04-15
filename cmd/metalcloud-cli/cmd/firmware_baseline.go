package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/firmware_baseline"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	firmwareBaselineFlags = struct {
		configSource string
		searchSource string
	}{}

	firmwareBaselineCmd = &cobra.Command{
		Use:     "firmware-baseline [command]",
		Aliases: []string{"fw-baseline", "baseline"},
		Short:   "Firmware baseline management",
		Long:    `Firmware baseline management commands.`,
	}

	firmwareBaselineListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all firmware baselines.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_baseline.FirmwareBaselineList(cmd.Context())
		},
	}

	firmwareBaselineGetCmd = &cobra.Command{
		Use:          "get firmware_baseline_id",
		Aliases:      []string{"show"},
		Short:        "Get firmware baseline details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_baseline.FirmwareBaselineGet(cmd.Context(), args[0])
		},
	}

	firmwareBaselineConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Get firmware baseline configuration example.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_baseline.FirmwareBaselineConfigExample(cmd.Context())
		},
	}

	firmwareBaselineCreateCmd = &cobra.Command{
		Use:          "create",
		Aliases:      []string{"new"},
		Short:        "Create a new firmware baseline.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(firmwareBaselineFlags.configSource)
			if err != nil {
				return err
			}

			return firmware_baseline.FirmwareBaselineCreate(cmd.Context(), config)
		},
	}

	firmwareBaselineUpdateCmd = &cobra.Command{
		Use:          "update firmware_baseline_id",
		Aliases:      []string{"edit"},
		Short:        "Update a firmware baseline.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(firmwareBaselineFlags.configSource)
			if err != nil {
				return err
			}

			return firmware_baseline.FirmwareBaselineUpdate(cmd.Context(), args[0], config)
		},
	}

	firmwareBaselineDeleteCmd = &cobra.Command{
		Use:          "delete firmware_baseline_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a firmware baseline.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_baseline.FirmwareBaselineDelete(cmd.Context(), args[0])
		},
	}

	firmwareBaselineSearchCmd = &cobra.Command{
		Use:          "search",
		Short:        "Search firmware baselines.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(firmwareBaselineFlags.searchSource)
			if err != nil {
				return err
			}

			return firmware_baseline.FirmwareBaselineSearch(cmd.Context(), config)
		},
	}

	firmwareBaselineSearchExampleCmd = &cobra.Command{
		Use:          "search-example",
		Short:        "Get firmware baseline search criteria example.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_baseline.FirmwareBaselineSearchExample(cmd.Context())
		},
	}
)

func init() {
	rootCmd.AddCommand(firmwareBaselineCmd)

	firmwareBaselineCmd.AddCommand(firmwareBaselineListCmd)
	firmwareBaselineCmd.AddCommand(firmwareBaselineGetCmd)
	firmwareBaselineCmd.AddCommand(firmwareBaselineConfigExampleCmd)

	firmwareBaselineCmd.AddCommand(firmwareBaselineCreateCmd)
	firmwareBaselineCreateCmd.Flags().StringVar(&firmwareBaselineFlags.configSource, "config-source", "", "Source of the new firmware baseline configuration. Can be 'pipe' or path to a JSON file.")
	firmwareBaselineCreateCmd.MarkFlagsOneRequired("config-source")

	firmwareBaselineCmd.AddCommand(firmwareBaselineUpdateCmd)
	firmwareBaselineUpdateCmd.Flags().StringVar(&firmwareBaselineFlags.configSource, "config-source", "", "Source of the firmware baseline configuration updates. Can be 'pipe' or path to a JSON file.")
	firmwareBaselineUpdateCmd.MarkFlagsOneRequired("config-source")

	firmwareBaselineCmd.AddCommand(firmwareBaselineDeleteCmd)

	firmwareBaselineCmd.AddCommand(firmwareBaselineSearchCmd)
	firmwareBaselineSearchCmd.Flags().StringVar(&firmwareBaselineFlags.searchSource, "search-source", "", "Source of the search criteria. Can be 'pipe' or path to a JSON file.")
	firmwareBaselineSearchCmd.MarkFlagsOneRequired("search-source")

	firmwareBaselineCmd.AddCommand(firmwareBaselineSearchExampleCmd)
}
