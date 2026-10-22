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

		// list
		filterId                     []string
		filterCatalogId              []string
		filterExternalId             []string
		filterPackageId              []string
		filterPackageVersion         []string
		filterUpdateSeverity         []string
		filterRebootRequired         []string
		filterCreatedTimestamp       []string
		filterVendorReleaseTimestamp []string
		page                         int
		limit                        int
	}{}

	firmwareBinaryCmd = &cobra.Command{
		Use:     "firmware-binary [command]",
		Aliases: []string{"fw-binary", "firmware-bin"},
		Short:   "Manage individual firmware binary files and packages",
		Long: `Manage individual firmware binary files and packages.

Firmware binaries are the actual firmware files that can be applied to hardware components.
They are typically part of firmware catalogs but can also be managed individually for
custom firmware deployment scenarios.

Each firmware binary contains:
  • Binary file data and metadata
  • Compatibility information (hardware models, component types)
  • Version and release information
  • Installation instructions and dependencies
  • Checksums and verification data

Use cases:
  • Manual firmware binary registration
  • Custom firmware package creation
  • Individual binary management outside of vendor catalogs
  • Firmware binary inspection and validation`,
	}

	firmwareBinaryListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all firmware binaries",
		Long: `List all firmware binaries in the system.

This command displays all available firmware binaries with their basic information including:
- Binary ID and name
- Catalog ID and package information
- Target hardware components and models
- Version information and update severity
- Reboot requirements and vendor release details
- Release timestamps and download URLs
- External ID and vendor info URLs

The output includes both catalog-managed binaries and individually registered ones.

Filters are applied server-side. Every --filter-* flag accepts a comma-separated list or can be
repeated; multiple values are combined with OR. A single value is sent as-is, so it may also use the
API operator syntax (e.g. '$not:$eq:value', '$ilike:value' or '$gt:2025-01-01').
Without --page/--limit all matching binaries are returned.

Optional Flags:
  --filter-id                        Filter by binary ID
  --filter-catalog-id                Filter by catalog ID
  --filter-update-severity           Filter by update severity (e.g. critical)
  --filter-reboot-required           Filter by reboot requirement (true or false)
  --filter-package-version           Filter by package version
  --filter-external-id               Filter by the vendor's artefact ID
  --filter-package-id                Filter by package ID
  --filter-created-timestamp         Filter by creation timestamp
  --filter-vendor-release-timestamp  Filter by vendor release timestamp
  --page                             Page number to retrieve (default: return all records)
  --limit                            Number of records per page, or maximum records when used without --page

Examples:
  # List all firmware binaries
  metalcloud-cli firmware-binary list

  # Critical binaries in catalog 8813
  metalcloud-cli firmware-binary list --filter-catalog-id 8813 --filter-update-severity critical

  # Binaries that do not require a reboot
  metalcloud-cli fw-binary ls --filter-reboot-required false

  # Second page of 50 binaries
  metalcloud-cli firmware-binary list --page 2 --limit 50`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_binary.FirmwareBinaryList(cmd.Context(), firmware_binary.ListFlags{
				FilterId:                     firmwareBinaryFlags.filterId,
				FilterCatalogId:              firmwareBinaryFlags.filterCatalogId,
				FilterExternalId:             firmwareBinaryFlags.filterExternalId,
				FilterPackageId:              firmwareBinaryFlags.filterPackageId,
				FilterPackageVersion:         firmwareBinaryFlags.filterPackageVersion,
				FilterUpdateSeverity:         firmwareBinaryFlags.filterUpdateSeverity,
				FilterRebootRequired:         firmwareBinaryFlags.filterRebootRequired,
				FilterCreatedTimestamp:       firmwareBinaryFlags.filterCreatedTimestamp,
				FilterVendorReleaseTimestamp: firmwareBinaryFlags.filterVendorReleaseTimestamp,
				Page:                         firmwareBinaryFlags.page,
				Limit:                        firmwareBinaryFlags.limit,
			})
		},
	}

	firmwareBinaryGetCmd = &cobra.Command{
		Use:     "get firmware_binary_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific firmware binary",
		Long: `Get detailed information about a specific firmware binary.

This command displays comprehensive information about a firmware binary including:
- Basic metadata (name, packageVersion, catalogId)
- Hardware compatibility information (vendorSupportedDevices, vendorSupportedSystems)
- Binary file details (vendorDownloadUrl, cacheDownloadUrl, externalId)
- Installation requirements (rebootRequired, updateSeverity)
- Vendor release timestamp
- Creation and modification timestamps
- Associated links and resources

Arguments:
  firmware_binary_id    The numeric ID of the firmware binary to retrieve

Examples:
  metalcloud-cli firmware-binary get 67890
  metalcloud-cli fw-binary show 12345
  metalcloud-cli firmware-bin get 98765`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_binary.FirmwareBinaryGet(cmd.Context(), args[0])
		},
	}

	firmwareBinaryConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Display configuration file template for creating firmware binaries",
		Long: `Display configuration file template for creating firmware binaries.

This command outputs a comprehensive example configuration file that shows all available
options for creating firmware binaries. The example includes all required and optional
fields with their descriptions and sample values.

The configuration template covers:
- Basic metadata (name, catalogId, packageId, packageVersion)
- File information (vendorDownloadUrl, vendorInfoUrl, cacheDownloadUrl)
- Hardware compatibility (vendorSupportedDevices, vendorSupportedSystems)
- Installation requirements (rebootRequired, updateSeverity)
- External references (externalId)

Use this template as a starting point for creating your own firmware binary configurations.
Copy the output to a file, modify the values as needed, and use it with the create command.

Examples:
  metalcloud-cli firmware-binary config-example > my-firmware.json
  metalcloud-cli fw-binary config-example | grep -A 50 "dell"`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_binary.FirmwareBinaryConfigExample(cmd.Context())
		},
	}

	firmwareBinaryCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new firmware binary from configuration file",
		Long: `Create a new firmware binary from configuration file.

This command creates a new firmware binary registration in the system using a configuration
file that defines all the binary's properties, compatibility information, and metadata.

The firmware binary will be validated for:
- Hardware compatibility specifications
- Version information consistency
- Required metadata completeness
- Update severity classification

Use the 'config-example' command to generate a template configuration file with all
available options and their descriptions.

Required Flags:
  --config-source    Source of the firmware binary configuration (JSON/YAML file path or 'pipe')

The configuration file must include:
- Basic metadata (name, catalogId)
- Vendor download URL
- Hardware compatibility information (vendorSupportedDevices, vendorSupportedSystems)
- Update requirements (rebootRequired, updateSeverity)

Optional configuration includes:
- Package identification (packageId, packageVersion)
- External references (externalId, vendorInfoUrl, cacheDownloadUrl)
- Release information (vendorReleaseTimestamp)

Examples:
  metalcloud-cli firmware-binary create --config-source ./bios-update.json
  cat firmware-config.json | metalcloud-cli fw-binary create --config-source pipe
  metalcloud-cli firmware-bin new --config-source ./dell-r640-bios.yaml

Configuration file example (bios-update.json):
{
  "name": "BIOS-R740-2.15.0",
  "catalogId": 1,
  "vendorDownloadUrl": "https://dell.com/downloads/firmware/R740/BIOS-2.15.0.bin",
  "vendorInfoUrl": "https://dell.com/support/firmware/R740/BIOS/2.15.0",
  "externalId": "DELL-R740-BIOS-2.15.0",
  "packageId": "BIOS",
  "packageVersion": "2.15.0",
  "rebootRequired": true,
  "updateSeverity": "recommended",
  "vendorReleaseTimestamp": "2024-04-01T12:00:00Z",
  "vendorSupportedDevices": [
    {
      "model": "PowerEdge R740",
      "type": "server"
    }
  ],
  "vendorSupportedSystems": [
    {
      "os": "any",
      "version": "any"
    }
  ]
}`,
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
		Use:     "delete firmware_binary_id",
		Aliases: []string{"rm"},
		Short:   "Delete a firmware binary permanently",
		Long: `Delete a firmware binary permanently.

This command removes a firmware binary registration from the system. This action is irreversible
and will delete all associated metadata and references to the binary file.

Note: This command only removes the binary registration from the system database. If the firmware
binary file is stored externally (e.g., on a remote repository), the actual file will not be
deleted and must be removed separately if needed.

Arguments:
  firmware_binary_id    The numeric ID of the firmware binary to delete

Examples:
  metalcloud-cli firmware-binary delete 67890
  metalcloud-cli fw-binary rm 12345
  metalcloud-cli firmware-bin delete 98765`,
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
	firmwareBinaryListCmd.Flags().StringSliceVar(&firmwareBinaryFlags.filterId, "filter-id", nil, "Filter the result by binary ID.")
	firmwareBinaryListCmd.Flags().StringSliceVar(&firmwareBinaryFlags.filterCatalogId, "filter-catalog-id", nil, "Filter the result by catalog ID.")
	firmwareBinaryListCmd.Flags().StringSliceVar(&firmwareBinaryFlags.filterUpdateSeverity, "filter-update-severity", nil, "Filter the result by update severity.")
	firmwareBinaryListCmd.Flags().StringSliceVar(&firmwareBinaryFlags.filterRebootRequired, "filter-reboot-required", nil, "Filter the result by reboot requirement (true or false).")
	firmwareBinaryListCmd.Flags().StringSliceVar(&firmwareBinaryFlags.filterPackageVersion, "filter-package-version", nil, "Filter the result by package version.")
	firmwareBinaryListCmd.Flags().StringSliceVar(&firmwareBinaryFlags.filterExternalId, "filter-external-id", nil, "Filter the result by the vendor's artefact ID.")
	firmwareBinaryListCmd.Flags().StringSliceVar(&firmwareBinaryFlags.filterPackageId, "filter-package-id", nil, "Filter the result by package ID.")
	firmwareBinaryListCmd.Flags().StringSliceVar(&firmwareBinaryFlags.filterCreatedTimestamp, "filter-created-timestamp", nil, "Filter the result by creation timestamp.")
	firmwareBinaryListCmd.Flags().StringSliceVar(&firmwareBinaryFlags.filterVendorReleaseTimestamp, "filter-vendor-release-timestamp", nil, "Filter the result by vendor release timestamp.")
	firmwareBinaryListCmd.Flags().IntVar(&firmwareBinaryFlags.page, "page", 0, "Page number to retrieve (default: return all records).")
	firmwareBinaryListCmd.Flags().IntVar(&firmwareBinaryFlags.limit, "limit", 0, "Number of records per page, or maximum records without --page (default: return all records).")

	firmwareBinaryCreateCmd.Flags().StringVar(&firmwareBinaryFlags.configSource, "config-source", "", "Source of the new firmware binary configuration. Can be 'pipe' or path to a JSON file.")
	firmwareBinaryCreateCmd.MarkFlagsOneRequired("config-source")

	firmwareBinaryCmd.AddCommand(firmwareBinaryDeleteCmd)
}
