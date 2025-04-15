package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/firmware_binary"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	firmwareBinaryFlags = struct {
		configSource string
	}{}

	firmwareBinaryCmd = &cobra.Command{
		Use:     "firmware-binary [command]",
		Aliases: []string{"fw-binary", "firmware-bin"},
		Short:   "Firmware binary management",
		Long:    `Firmware binary management commands.`,
	}

	firmwareBinaryListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all firmware binaries.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_binary.FirmwareBinaryList(cmd.Context())
		},
	}

	firmwareBinaryGetCmd = &cobra.Command{
		Use:          "get firmware_binary_id",
		Aliases:      []string{"show"},
		Short:        "Get firmware binary details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_binary.FirmwareBinaryGet(cmd.Context(), args[0])
		},
	}

	firmwareBinaryConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Get firmware binary configuration example.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_binary.FirmwareBinaryConfigExample(cmd.Context())
		},
	}

	firmwareBinaryCreateCmd = &cobra.Command{
		Use:          "create",
		Aliases:      []string{"new"},
		Short:        "Create a new firmware binary.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(firmwareBinaryFlags.configSource)
			if err != nil {
				return err
			}

			return firmware_binary.FirmwareBinaryCreate(cmd.Context(), config)
		},
	}

	firmwareBinaryDeleteCmd = &cobra.Command{
		Use:          "delete firmware_binary_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a firmware binary.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_binary.FirmwareBinaryDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(firmwareBinaryCmd)

	firmwareBinaryCmd.AddCommand(firmwareBinaryListCmd)
	firmwareBinaryCmd.AddCommand(firmwareBinaryGetCmd)
	firmwareBinaryCmd.AddCommand(firmwareBinaryConfigExampleCmd)

	firmwareBinaryCmd.AddCommand(firmwareBinaryCreateCmd)
	firmwareBinaryCreateCmd.Flags().StringVar(&firmwareBinaryFlags.configSource, "config-source", "", "Source of the new firmware binary configuration. Can be 'pipe' or path to a JSON file.")
	firmwareBinaryCreateCmd.MarkFlagsOneRequired("config-source")

	firmwareBinaryCmd.AddCommand(firmwareBinaryDeleteCmd)
}
