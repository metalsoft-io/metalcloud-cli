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
		Short:   "Manage firmware baselines for consistent hardware configurations",
		Long: `Manage firmware baselines for consistent hardware configurations.

Firmware baselines define standardized firmware configurations for specific hardware
types and deployment scenarios. They specify the firmware level and filtering criteria
for consistent hardware management across your infrastructure.

A firmware baseline includes:
  • Name and description for identification
  • Level specification (e.g., PRODUCTION, DEVELOPMENT)
  • Level filter for targeting specific hardware types
  • Catalog associations for firmware sources

Use cases:
  • Standardizing firmware configurations across server fleets
  • Defining deployment levels for different environments
  • Managing firmware catalog associations
  • Creating hardware-specific configuration templates`,
	}

	firmwareBaselineListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all firmware baselines",
		Long: `List all firmware baselines in the system.

This command displays all available firmware baselines with their basic information including:
- Baseline ID and name
- Level and level filter specifications
- Description and catalog associations
- Creation timestamp

The output provides an overview of all standardized firmware configurations
available for deployment across your infrastructure.

No additional flags are required for this command.

Examples:
  metalcloud-cli firmware-baseline list
  metalcloud-cli fw-baseline ls
  metalcloud-cli baseline list`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_baseline.FirmwareBaselineList(cmd.Context())
		},
	}

	firmwareBaselineGetCmd = &cobra.Command{
		Use:     "get firmware_baseline_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific firmware baseline",
		Long: `Get detailed information about a specific firmware baseline.

This command displays comprehensive information about a firmware baseline including:
- Basic metadata (name, description, level, level filter)
- Catalog associations
- Creation timestamp
- Unique baseline identifier

The firmware baseline contains the essential configuration for standardized firmware
deployments across compatible hardware.

Arguments:
  firmware_baseline_id    The ID of the firmware baseline to retrieve

Examples:
  metalcloud-cli firmware-baseline get 54321
  metalcloud-cli fw-baseline show dell-r640-standard
  metalcloud-cli baseline get production-baseline-v2.1`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_baseline.FirmwareBaselineGet(cmd.Context(), args[0])
		},
	}

	firmwareBaselineConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Display configuration file template for creating firmware baselines",
		Long: `Display configuration file template for creating firmware baselines.

This command outputs a comprehensive example configuration file that shows all available
options for creating firmware baselines. The example includes all required and optional
fields with their descriptions and sample values.

The configuration template covers:
- Basic metadata (name, description)
- Target hardware specifications (level, level filter)
- Catalog associations

Use this template as a starting point for creating your own firmware baseline configurations.
Copy the output to a file, modify the values as needed, and use it with the create command.

Examples:
  metalcloud-cli firmware-baseline config-example > baseline-template.json
  metalcloud-cli fw-baseline config-example | grep -A 5 "level"
  metalcloud-cli baseline config-example | jq '.name'`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_baseline.FirmwareBaselineConfigExample(cmd.Context())
		},
	}

	firmwareBaselineCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new firmware baseline from configuration file",
		Long: `Create a new firmware baseline from configuration file.

This command creates a new firmware baseline definition using a configuration file that
specifies all the baseline properties and target hardware specifications.

The firmware baseline will be validated for:
- Required field completeness
- Level and level filter consistency
- Catalog reference validity

Use the 'config-example' command to generate a template configuration file with all
available options and their descriptions.

Required Flags:
  --config-source    Source of the firmware baseline configuration (JSON/YAML file path or 'pipe')

The configuration file must include:
- Basic metadata (name, level, levelFilter)

Optional configuration includes:
- Description and catalog associations

Examples:
  metalcloud-cli firmware-baseline create --config-source ./production-baseline.json
  cat baseline-config.json | metalcloud-cli fw-baseline create --config-source pipe
  metalcloud-cli baseline new --config-source ./dell-r640-baseline.yaml

Configuration file example (production-baseline.json):
{
  "name": "Production Dell R640 Baseline",
  "description": "Standard firmware configuration for Dell PowerEdge R640 production servers",
  "level": "PRODUCTION",
  "levelFilter": ["dell_r640", "dell_r640_gen2"],
  "catalog": ["dell-catalog-r640", "dell-catalog-common"]
}`,
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
		Use:     "update firmware_baseline_id",
		Aliases: []string{"edit"},
		Short:   "Update an existing firmware baseline",
		Long: `Update an existing firmware baseline.

This command allows you to update the configuration of an existing firmware baseline.
Updates are provided through a configuration file (JSON or YAML format) that contains
the new settings to apply.

The configuration file should contain only the fields you want to update. Common
updates include:
- Basic metadata (name, description)
- Level and level filter specifications  
- Catalog associations

The baseline will be revalidated after updates to ensure consistency and compatibility.

Required Flags:
  --config-source    Source of the configuration updates (JSON/YAML file path or 'pipe')

Arguments:
  firmware_baseline_id    The ID of the firmware baseline to update

Examples:
  metalcloud-cli firmware-baseline update 54321 --config-source ./baseline-updates.json
  cat updates.json | metalcloud-cli fw-baseline update production-baseline --config-source pipe
  metalcloud-cli baseline edit dell-r640-standard --config-source ./version-update.yaml

Configuration file example (baseline-updates.json):
{
  "name": "Updated Production Dell R640 Baseline",
  "description": "Updated with latest security patches",
  "level": "PRODUCTION",
  "levelFilter": ["dell_r640", "dell_r640_gen2"],
  "catalog": ["dell-catalog-r640-v2", "dell-catalog-common"]
}`,
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
		Use:     "delete firmware_baseline_id",
		Aliases: []string{"rm"},
		Short:   "Delete a firmware baseline permanently",
		Long: `Delete a firmware baseline permanently.

This command removes a firmware baseline definition from the system. This action is irreversible
and will delete all associated configuration data, component specifications, and deployment policies.

Note: Deleting a firmware baseline does not affect the underlying firmware catalogs or binaries.
Only the baseline definition and its configuration are removed.

Before deletion, ensure that:
- No active deployments are using this baseline
- No automated processes reference this baseline
- You have backups of the configuration if needed for future reference

Arguments:
  firmware_baseline_id    The ID of the firmware baseline to delete

Examples:
  metalcloud-cli firmware-baseline delete 54321
  metalcloud-cli fw-baseline rm production-baseline-v2.1
  metalcloud-cli baseline delete dell-r640-standard`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FIRMWARE_BASELINES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return firmware_baseline.FirmwareBaselineDelete(cmd.Context(), args[0])
		},
	}

	firmwareBaselineSearchCmd = &cobra.Command{
		Use:   "search",
		Short: "Search firmware baselines by criteria",
		Long: `Search firmware baselines by criteria.

This command allows you to search for firmware baselines using specific criteria
such as vendor, datacenter, server type, and component filters.
Search criteria are provided through a configuration file (JSON or YAML format).

The search can filter baselines by:
- Vendor (e.g., DELL)
- Datacenter locations
- Server types and OS templates
- Component filters for specific hardware

Use the 'search-example' command to see available search criteria and their format.

Required Flags:
  --search-source    Source of the search criteria (JSON/YAML file path or 'pipe')

Examples:
  metalcloud-cli firmware-baseline search --search-source ./search-criteria.json
  cat search.json | metalcloud-cli fw-baseline search --search-source pipe
  metalcloud-cli baseline search --search-source ./dell-baselines.yaml

Search criteria example (search-criteria.json):
{
  "vendor": "DELL",
  "baselineFilter": {
    "datacenter": ["datacenter-1"],
    "serverType": ["dell_r740", "dell_r640"],
    "osTemplate": ["ubuntu-20.04"],
    "baselineId": ["baseline-1"]
  },
  "serverComponentFilter": {
    "dellComponentFilter": {
      "componentId": "component-1"
    }
  }
}`,
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
		Use:   "search-example",
		Short: "Display search criteria template for firmware baseline search",
		Long: `Display search criteria template for firmware baseline search.

This command outputs a comprehensive example search criteria file that shows all available
search options for finding firmware baselines. The example includes all searchable fields
with their descriptions and sample values.

The search criteria template covers:
- Vendor filtering (DELL, etc.)
- Baseline filtering (datacenter, server type, OS template, baseline ID)
- Component filtering for specific hardware components

Use this template as a starting point for creating your own search criteria.
Copy the output to a file, modify the values as needed, and use it with the search command.

Examples:
  metalcloud-cli firmware-baseline search-example > search-template.json
  metalcloud-cli fw-baseline search-example | grep -A 10 "vendor"
  metalcloud-cli baseline search-example | jq '.baselineFilter'`,
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
