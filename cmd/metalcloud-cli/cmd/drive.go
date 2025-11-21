package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/drive"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	driveFlags = struct {
		configSource string
		filterStatus []string
	}{}

	driveCmd = &cobra.Command{
		Use:     "drive [command]",
		Aliases: []string{"drives", "disk", "block-storage-volume"},
		Short:   "Manage drives within infrastructures",
		Long: `Manage drives within infrastructures including creation, configuration, metadata updates, and host assignments.

Drives are storage devices that can be attached to servers within an infrastructure. This command group
provides comprehensive drive management capabilities including listing, creating, updating configurations,
managing metadata, and controlling host assignments.

Available Commands:
  list          List all drives for an infrastructure
  get           Get detailed information about a specific drive
  create        Create a new drive with specified configuration
  delete        Remove a drive from the infrastructure
  update-config Update drive configuration settings
  update-meta   Update drive metadata
  get-hosts     Show hosts assigned to a drive
  update-hosts  Update host assignments for a drive
  config-info   Get configuration information for a drive

Examples:
  # List all drives in an infrastructure
  metalcloud-cli drive list my-infrastructure

  # Get details of a specific drive
  metalcloud-cli drive get my-infrastructure 12345

  # Create a new drive from JSON configuration
  metalcloud-cli drive create my-infrastructure --config-source drive-config.json`,
	}

	driveListCmd = &cobra.Command{
		Use:     "list infrastructure_id_or_label",
		Aliases: []string{"ls"},
		Short:   "List all drives within an infrastructure",
		Long: `List all drives within an infrastructure with optional status filtering.

This command displays a comprehensive list of all drives associated with the specified infrastructure,
including their IDs, configurations, current status, and metadata.

Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure

Flags:
  --filter-status strings       Filter drives by status (optional)
                               Multiple statuses can be provided as comma-separated values

Examples:
  # List all drives in an infrastructure
  metalcloud-cli drive list my-infrastructure

  # List drives with specific status
  metalcloud-cli drive list my-infrastructure --filter-status active

  # List drives with multiple statuses
  metalcloud-cli drive list my-infrastructure --filter-status active,pending`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_DRIVES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return drive.DriveList(cmd.Context(), args[0], driveFlags.filterStatus)
		},
	}

	driveGetCmd = &cobra.Command{
		Use:     "get infrastructure_id_or_label drive_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific drive",
		Long: `Get detailed information about a specific drive within an infrastructure.

This command retrieves comprehensive information about a single drive, including its configuration,
status, metadata, and any associated host assignments.

Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure
  drive_id                     The unique identifier of the drive

Examples:
  # Get details of a specific drive
  metalcloud-cli drive get my-infrastructure 12345

  # Get drive details using infrastructure ID
  metalcloud-cli drive get 1001 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_DRIVES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return drive.DriveGet(cmd.Context(), args[0], args[1])
		},
	}

	driveCreateCmd = &cobra.Command{
		Use:     "create infrastructure_id_or_label",
		Aliases: []string{"new"},
		Short:   "Create a new drive with specified configuration",
		Long: `Create a new drive within an infrastructure using JSON configuration.

This command creates a new drive with the specified configuration. The drive configuration
must be provided via JSON file or stdin pipe.

Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure

Required Flags:
  --config-source string        Source of the new drive configuration
                               Can be 'pipe' (for stdin input) or path to a JSON file

Examples:
  # Create a drive from JSON file
  metalcloud-cli drive create my-infrastructure --config-source drive-config.json

  # Create a drive from stdin
  echo '{"size_bytes": 1000000000, "type": "ssd"}' | metalcloud-cli drive create my-infrastructure --config-source pipe

  # Create a drive with configuration from file
  metalcloud-cli drive create 1001 --config-source /path/to/drive.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_DRIVES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(driveFlags.configSource)
			if err != nil {
				return err
			}

			return drive.DriveCreate(cmd.Context(), args[0], config)
		},
	}

	driveDeleteCmd = &cobra.Command{
		Use:     "delete infrastructure_id_or_label drive_id",
		Aliases: []string{"rm"},
		Short:   "Remove a drive from the infrastructure",
		Long: `Remove a drive from the infrastructure permanently.

This command deletes a drive from the specified infrastructure. The operation is irreversible
and will remove all drive data and configurations.

Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure
  drive_id                     The unique identifier of the drive to delete

Examples:
  # Delete a drive by ID
  metalcloud-cli drive delete my-infrastructure 12345

  # Delete a drive using infrastructure ID
  metalcloud-cli drive rm 1001 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_DRIVES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return drive.DriveDelete(cmd.Context(), args[0], args[1])
		},
	}

	driveUpdateConfigCmd = &cobra.Command{
		Use:     "update-config infrastructure_id_or_label drive_id",
		Aliases: []string{"config-update"},
		Short:   "Update drive configuration settings",
		Long: `Update drive configuration settings using JSON configuration.

This command updates the configuration of an existing drive. The new configuration
must be provided via JSON file or stdin pipe. Only the specified fields will be updated.

Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure
  drive_id                     The unique identifier of the drive

Required Flags:
  --config-source string        Source of the drive configuration updates
                               Can be 'pipe' (for stdin input) or path to a JSON file

Examples:
  # Update drive configuration from JSON file
  metalcloud-cli drive update-config my-infrastructure 12345 --config-source drive-updates.json

  # Update drive configuration from stdin
  echo '{"size_bytes": 2000000000}' | metalcloud-cli drive update-config my-infrastructure 12345 --config-source pipe

  # Update using configuration file path
  metalcloud-cli drive config-update 1001 67890 --config-source /path/to/config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_DRIVES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(driveFlags.configSource)
			if err != nil {
				return err
			}

			return drive.DriveUpdateConfig(cmd.Context(), args[0], args[1], config)
		},
	}

	driveUpdateMetaCmd = &cobra.Command{
		Use:     "update-meta infrastructure_id_or_label drive_id",
		Aliases: []string{"meta-update"},
		Short:   "Update drive metadata",
		Long: `Update drive metadata using JSON configuration.

This command updates the metadata of an existing drive. Metadata includes custom labels,
tags, and other descriptive information. The metadata updates must be provided via JSON file
or stdin pipe.

Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure
  drive_id                     The unique identifier of the drive

Required Flags:
  --config-source string        Source of the drive metadata updates
                               Can be 'pipe' (for stdin input) or path to a JSON file

Examples:
  # Update drive metadata from JSON file
  metalcloud-cli drive update-meta my-infrastructure 12345 --config-source metadata.json

  # Update metadata from stdin
  echo '{"labels": {"environment": "production"}}' | metalcloud-cli drive update-meta my-infrastructure 12345 --config-source pipe

  # Update using metadata file path
  metalcloud-cli drive meta-update 1001 67890 --config-source /path/to/metadata.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_DRIVES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(driveFlags.configSource)
			if err != nil {
				return err
			}

			return drive.DriveUpdateMeta(cmd.Context(), args[0], args[1], config)
		},
	}

	driveGetHostsCmd = &cobra.Command{
		Use:     "get-hosts infrastructure_id_or_label drive_id",
		Aliases: []string{"hosts"},
		Short:   "Show hosts assigned to a drive",
		Long: `Show hosts assigned to a drive within an infrastructure.

This command displays all hosts that are currently assigned to the specified drive,
including host IDs, names, and assignment details.

Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure
  drive_id                     The unique identifier of the drive

Examples:
  # Get hosts assigned to a drive
  metalcloud-cli drive get-hosts my-infrastructure 12345

  # Get hosts using infrastructure ID
  metalcloud-cli drive hosts 1001 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_DRIVES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return drive.DriveGetHosts(cmd.Context(), args[0], args[1])
		},
	}

	driveUpdateHostsCmd = &cobra.Command{
		Use:     "update-hosts infrastructure_id_or_label drive_id",
		Aliases: []string{"hosts-update"},
		Short:   "Update host assignments for a drive",
		Long: `Update host assignments for a drive within an infrastructure using JSON configuration.

This command updates the host assignments for a drive. You can assign or remove hosts from the drive.
The host configuration must be provided via JSON file or stdin pipe.

Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure
  drive_id                     The unique identifier of the drive

Required Flags:
  --config-source string        Source of the drive hosts configuration
                               Can be 'pipe' (for stdin input) or path to a JSON file

Examples:
  # Update drive host assignments from JSON file
  metalcloud-cli drive update-hosts my-infrastructure 12345 --config-source hosts.json

  # Update hosts from stdin
  echo '{"host_ids": [123, 456]}' | metalcloud-cli drive update-hosts my-infrastructure 12345 --config-source pipe

  # Update using hosts configuration file
  metalcloud-cli drive hosts-update 1001 67890 --config-source /path/to/hosts.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_DRIVES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(driveFlags.configSource)
			if err != nil {
				return err
			}

			return drive.DriveUpdateHosts(cmd.Context(), args[0], args[1], config)
		},
	}

	driveGetConfigInfoCmd = &cobra.Command{
		Use:     "config-info infrastructure_id_or_label drive_id",
		Aliases: []string{"get-config-info"},
		Short:   "Get configuration information for a drive",
		Long: `Get configuration information for a drive within an infrastructure.

This command retrieves the current configuration information for a specified drive,
including all configuration parameters, settings, and current values.

Arguments:
  infrastructure_id_or_label    The ID or label of the infrastructure
  drive_id                     The unique identifier of the drive

Examples:
  # Get configuration information for a drive
  metalcloud-cli drive config-info my-infrastructure 12345

  # Get config info using infrastructure ID
  metalcloud-cli drive get-config-info 1001 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_DRIVES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return drive.DriveGetConfigInfo(cmd.Context(), args[0], args[1])
		},
	}
)

func init() {
	rootCmd.AddCommand(driveCmd)

	driveCmd.AddCommand(driveListCmd)
	driveListCmd.Flags().StringSliceVar(&driveFlags.filterStatus, "filter-status", nil, "Filter the result by drive status.")

	driveCmd.AddCommand(driveGetCmd)

	driveCmd.AddCommand(driveCreateCmd)
	driveCreateCmd.Flags().StringVar(&driveFlags.configSource, "config-source", "", "Source of the new drive configuration. Can be 'pipe' or path to a JSON file.")
	driveCreateCmd.MarkFlagsOneRequired("config-source")

	driveCmd.AddCommand(driveDeleteCmd)

	driveCmd.AddCommand(driveUpdateConfigCmd)
	driveUpdateConfigCmd.Flags().StringVar(&driveFlags.configSource, "config-source", "", "Source of the drive configuration updates. Can be 'pipe' or path to a JSON file.")
	driveUpdateConfigCmd.MarkFlagsOneRequired("config-source")

	driveCmd.AddCommand(driveUpdateMetaCmd)
	driveUpdateMetaCmd.Flags().StringVar(&driveFlags.configSource, "config-source", "", "Source of the drive metadata updates. Can be 'pipe' or path to a JSON file.")
	driveUpdateMetaCmd.MarkFlagsOneRequired("config-source")

	driveCmd.AddCommand(driveGetHostsCmd)

	driveCmd.AddCommand(driveUpdateHostsCmd)
	driveUpdateHostsCmd.Flags().StringVar(&driveFlags.configSource, "config-source", "", "Source of the drive hosts configuration. Can be 'pipe' or path to a JSON file.")
	driveUpdateHostsCmd.MarkFlagsOneRequired("config-source")

	driveCmd.AddCommand(driveGetConfigInfoCmd)
}
