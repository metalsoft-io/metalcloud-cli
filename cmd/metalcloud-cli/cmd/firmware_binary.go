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
- Reboot requirements and vendor details
- Release timestamps and download URLs
- External ID and vendor info URLs

The output includes both catalog-managed binaries and individually registered ones.

No additional flags are required for this command.

Examples:
  metalcloud-cli firmware-binary list
  metalcloud-cli fw-binary ls`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_binary.FirmwareBinaryList(cmd.Context())
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
- Vendor information and release timestamp
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
- Vendor information and external references

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
- Vendor details

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
  ],
  "vendor": {
    "name": "Dell Inc.",
    "contact": "support@dell.com"
  }
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
	firmwareBinaryCreateCmd.Flags().StringVar(&firmwareBinaryFlags.configSource, "config-source", "", "Source of the new firmware binary configuration. Can be 'pipe' or path to a JSON file.")
	firmwareBinaryCreateCmd.MarkFlagsOneRequired("config-source")

	firmwareBinaryCmd.AddCommand(firmwareBinaryDeleteCmd)
}
