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
		filterStatus string
	}{}

	driveCmd = &cobra.Command{
		Use:     "drive [command]",
		Aliases: []string{"drives", "disk"},
		Short:   "Drive management",
		Long:    `Drive management commands.`,
	}

	driveListCmd = &cobra.Command{
		Use:          "list infrastructure_id_or_label",
		Aliases:      []string{"ls"},
		Short:        "List all drives for an infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return drive.DriveList(cmd.Context(), args[0], driveFlags.filterStatus)
		},
	}

	driveGetCmd = &cobra.Command{
		Use:          "get infrastructure_id_or_label drive_id",
		Aliases:      []string{"show"},
		Short:        "Get drive details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return drive.DriveGet(cmd.Context(), args[0], args[1])
		},
	}

	driveCreateCmd = &cobra.Command{
		Use:          "create infrastructure_id_or_label",
		Aliases:      []string{"new"},
		Short:        "Create a new drive.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
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
		Use:          "delete infrastructure_id_or_label drive_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a drive.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return drive.DriveDelete(cmd.Context(), args[0], args[1])
		},
	}

	driveUpdateConfigCmd = &cobra.Command{
		Use:          "update-config infrastructure_id_or_label drive_id",
		Aliases:      []string{"config-update"},
		Short:        "Update drive configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
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
		Use:          "update-meta infrastructure_id_or_label drive_id",
		Aliases:      []string{"meta-update"},
		Short:        "Update drive metadata.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
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
		Use:          "get-hosts infrastructure_id_or_label drive_id",
		Aliases:      []string{"hosts"},
		Short:        "Get hosts for a drive.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return drive.DriveGetHosts(cmd.Context(), args[0], args[1])
		},
	}

	driveUpdateHostsCmd = &cobra.Command{
		Use:          "update-hosts infrastructure_id_or_label drive_id",
		Aliases:      []string{"hosts-update"},
		Short:        "Update hosts for a drive.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
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
		Use:          "config-info infrastructure_id_or_label drive_id",
		Aliases:      []string{"get-config-info"},
		Short:        "Get configuration information for a drive.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_STORAGE_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return drive.DriveGetConfigInfo(cmd.Context(), args[0], args[1])
		},
	}
)

func init() {
	rootCmd.AddCommand(driveCmd)

	driveCmd.AddCommand(driveListCmd)
	driveListCmd.Flags().StringVar(&driveFlags.filterStatus, "filter-status", "", "Filter the result by drive status.")

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
