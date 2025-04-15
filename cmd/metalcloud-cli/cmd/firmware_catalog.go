package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/firmware_catalog"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	firmwareCatalogFlags = struct {
		configSource string
	}{}

	firmwareCatalogCmd = &cobra.Command{
		Use:     "firmware-catalog [command]",
		Aliases: []string{"fw-catalog", "firmware"},
		Short:   "Firmware catalog management",
		Long:    `Firmware catalog management commands.`,
	}

	firmwareCatalogListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all firmware catalogs.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_catalog.FirmwareCatalogList(cmd.Context())
		},
	}

	firmwareCatalogGetCmd = &cobra.Command{
		Use:          "get firmware_catalog_id",
		Aliases:      []string{"show"},
		Short:        "Get firmware catalog details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_catalog.FirmwareCatalogGet(cmd.Context(), args[0])
		},
	}

	firmwareCatalogConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Get firmware catalog configuration example.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_catalog.FirmwareCatalogConfigExample(cmd.Context())
		},
	}

	firmwareCatalogCreateCmd = &cobra.Command{
		Use:          "create",
		Aliases:      []string{"new"},
		Short:        "Create a new firmware catalog.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(firmwareCatalogFlags.configSource)
			if err != nil {
				return err
			}

			return firmware_catalog.FirmwareCatalogCreate(cmd.Context(), config)
		},
	}

	firmwareCatalogUpdateCmd = &cobra.Command{
		Use:          "update firmware_catalog_id",
		Aliases:      []string{"edit"},
		Short:        "Update a firmware catalog.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(firmwareCatalogFlags.configSource)
			if err != nil {
				return err
			}

			return firmware_catalog.FirmwareCatalogUpdate(cmd.Context(), args[0], config)
		},
	}

	firmwareCatalogDeleteCmd = &cobra.Command{
		Use:          "delete firmware_catalog_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a firmware catalog.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_catalog.FirmwareCatalogDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(firmwareCatalogCmd)

	firmwareCatalogCmd.AddCommand(firmwareCatalogListCmd)
	firmwareCatalogCmd.AddCommand(firmwareCatalogGetCmd)
	firmwareCatalogCmd.AddCommand(firmwareCatalogConfigExampleCmd)

	firmwareCatalogCmd.AddCommand(firmwareCatalogCreateCmd)
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.configSource, "config-source", "", "Source of the new firmware catalog configuration. Can be 'pipe' or path to a JSON file.")
	firmwareCatalogCreateCmd.MarkFlagsOneRequired("config-source")

	firmwareCatalogCmd.AddCommand(firmwareCatalogUpdateCmd)
	firmwareCatalogUpdateCmd.Flags().StringVar(&firmwareCatalogFlags.configSource, "config-source", "", "Source of the firmware catalog configuration updates. Can be 'pipe' or path to a JSON file.")
	firmwareCatalogUpdateCmd.MarkFlagsOneRequired("config-source")

	firmwareCatalogCmd.AddCommand(firmwareCatalogDeleteCmd)
}
