package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/firmware_catalog"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
)

var (
	firmwareCatalogFlags = struct {
		configSource      string
		name              string
		description       string
		vendor            string
		updateType        string
		vendorId          string
		vendorUrl         string
		vendorRelease     string
		msServerTypes     []string
		vendorServerTypes []string
		vendorConfig      string
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
		Aliases:      []string{"example"},
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
			// If config source is provided, use it
			if firmwareCatalogFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(firmwareCatalogFlags.configSource)
				if err != nil {
					return err
				}
				return firmware_catalog.FirmwareCatalogCreate(cmd.Context(), config)
			}

			// Otherwise build config from command line parameters
			catalogConfig := sdk.CreateFirmwareCatalog{
				Name:       firmwareCatalogFlags.name,
				Vendor:     sdk.FirmwareVendorType(firmwareCatalogFlags.vendor),
				UpdateType: sdk.CatalogUpdateType(firmwareCatalogFlags.updateType),
			}

			if firmwareCatalogFlags.description != "" {
				catalogConfig.Description = sdk.PtrString(firmwareCatalogFlags.description)
			}

			if firmwareCatalogFlags.vendorId != "" {
				catalogConfig.VendorId = sdk.PtrString(firmwareCatalogFlags.vendorId)
			}

			if firmwareCatalogFlags.vendorUrl != "" {
				catalogConfig.VendorUrl = sdk.PtrString(firmwareCatalogFlags.vendorUrl)
			}

			if firmwareCatalogFlags.vendorRelease != "" {
				catalogConfig.VendorReleaseTimestamp = sdk.PtrString(firmwareCatalogFlags.vendorRelease)
			}

			if len(firmwareCatalogFlags.msServerTypes) > 0 {
				catalogConfig.MetalsoftServerTypesSupported = firmwareCatalogFlags.msServerTypes
			}

			if len(firmwareCatalogFlags.vendorServerTypes) > 0 {
				catalogConfig.VendorServerTypesSupported = firmwareCatalogFlags.vendorServerTypes
			}

			if firmwareCatalogFlags.vendorConfig != "" {
				var vendorConfiguration map[string]interface{}
				if err := json.Unmarshal([]byte(firmwareCatalogFlags.vendorConfig), &vendorConfiguration); err != nil {
					return fmt.Errorf("invalid vendor configuration JSON: %s", err)
				}
				catalogConfig.VendorConfiguration = vendorConfiguration
			}

			configBytes, err := json.Marshal(catalogConfig)
			if err != nil {
				return err
			}

			return firmware_catalog.FirmwareCatalogCreate(cmd.Context(), configBytes)
		},
	}

	firmwareCatalogUpdateCmd = &cobra.Command{
		Use:          "update firmware_catalog_id",
		Aliases:      []string{"edit"},
		Short:        "Update a firmware catalog.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		Args:         cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(firmwareCatalogFlags.configSource)
			if err != nil {
				return err
			}

			firmwareCatalogId := ""
			if len(args) > 0 {
				firmwareCatalogId = args[0]
			}

			return firmware_catalog.FirmwareCatalogUpdate(cmd.Context(), firmwareCatalogId, config)
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
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.name, "name", "", "Name of the firmware catalog (required)")
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.description, "description", "", "Description of the firmware catalog")
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.vendor, "vendor", "", "Vendor type (e.g., 'dell', 'hp') (required)")
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.updateType, "update-type", "", "Update type (e.g., 'online', 'offline') (required)")
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.vendorId, "vendor-id", "", "Vendor ID")
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.vendorUrl, "vendor-url", "", "Vendor URL")
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.vendorRelease, "vendor-release", "", "Vendor release timestamp")
	firmwareCatalogCreateCmd.Flags().StringSliceVar(&firmwareCatalogFlags.msServerTypes, "ms-server-types", []string{}, "Metalsoft server types supported (comma-separated)")
	firmwareCatalogCreateCmd.Flags().StringSliceVar(&firmwareCatalogFlags.vendorServerTypes, "vendor-server-types", []string{}, "Vendor server types supported (comma-separated)")
	firmwareCatalogCreateCmd.Flags().StringVar(&firmwareCatalogFlags.vendorConfig, "vendor-config", "", "Vendor configuration as JSON string")
	firmwareCatalogCreateCmd.MarkFlagsMutuallyExclusive("config-source", "name")
	firmwareCatalogCreateCmd.MarkFlagsRequiredTogether("name", "vendor", "update-type")

	firmwareCatalogCmd.AddCommand(firmwareCatalogUpdateCmd)
	firmwareCatalogUpdateCmd.Flags().StringVar(&firmwareCatalogFlags.configSource, "config-source", "", "Source of the firmware catalog configuration updates. Can be 'pipe' or path to a JSON file.")
	firmwareCatalogUpdateCmd.MarkFlagsOneRequired("config-source")

	firmwareCatalogCmd.AddCommand(firmwareCatalogDeleteCmd)
}
