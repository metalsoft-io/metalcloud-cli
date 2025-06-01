package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/file_share"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	fileShareFlags = struct {
		configSource string
		filterStatus []string
	}{}

	fileShareCmd = &cobra.Command{
		Use:     "file-share [command]",
		Aliases: []string{"fs"},
		Short:   "File Share management",
		Long:    `File Share management commands.`,
	}

	fileShareListCmd = &cobra.Command{
		Use:          "list infrastructure_id_or_label",
		Aliases:      []string{"ls"},
		Short:        "List all file shares for an infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FILE_SHARE_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return file_share.FileShareList(cmd.Context(), args[0], fileShareFlags.filterStatus)
		},
	}

	fileShareGetCmd = &cobra.Command{
		Use:          "get infrastructure_id_or_label file_share_id",
		Aliases:      []string{"show"},
		Short:        "Get file share details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FILE_SHARE_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return file_share.FileShareGet(cmd.Context(), args[0], args[1])
		},
	}

	fileShareCreateCmd = &cobra.Command{
		Use:          "create infrastructure_id_or_label",
		Aliases:      []string{"new"},
		Short:        "Create a new file share.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FILE_SHARE_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(fileShareFlags.configSource)
			if err != nil {
				return err
			}

			return file_share.FileShareCreate(cmd.Context(), args[0], config)
		},
	}

	fileShareDeleteCmd = &cobra.Command{
		Use:          "delete infrastructure_id_or_label file_share_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a file share.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FILE_SHARE_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return file_share.FileShareDelete(cmd.Context(), args[0], args[1])
		},
	}

	fileShareUpdateConfigCmd = &cobra.Command{
		Use:          "update-config infrastructure_id_or_label file_share_id",
		Aliases:      []string{"config-update"},
		Short:        "Update file share configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FILE_SHARE_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(fileShareFlags.configSource)
			if err != nil {
				return err
			}

			return file_share.FileShareUpdateConfig(cmd.Context(), args[0], args[1], config)
		},
	}

	fileShareUpdateMetaCmd = &cobra.Command{
		Use:          "update-meta infrastructure_id_or_label file_share_id",
		Aliases:      []string{"meta-update"},
		Short:        "Update file share metadata.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FILE_SHARE_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(fileShareFlags.configSource)
			if err != nil {
				return err
			}

			return file_share.FileShareUpdateMeta(cmd.Context(), args[0], args[1], config)
		},
	}

	fileShareGetHostsCmd = &cobra.Command{
		Use:          "get-hosts infrastructure_id_or_label file_share_id",
		Aliases:      []string{"hosts"},
		Short:        "Get hosts for a file share.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FILE_SHARE_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return file_share.FileShareGetHosts(cmd.Context(), args[0], args[1])
		},
	}

	fileShareUpdateHostsCmd = &cobra.Command{
		Use:          "update-hosts infrastructure_id_or_label file_share_id",
		Aliases:      []string{"hosts-update"},
		Short:        "Update hosts for a file share.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FILE_SHARE_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(fileShareFlags.configSource)
			if err != nil {
				return err
			}

			return file_share.FileShareUpdateHosts(cmd.Context(), args[0], args[1], config)
		},
	}

	fileShareGetConfigInfoCmd = &cobra.Command{
		Use:          "config-info infrastructure_id_or_label file_share_id",
		Aliases:      []string{"get-config-info"},
		Short:        "Get configuration information for a file share.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_FILE_SHARE_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return file_share.FileShareGetConfigInfo(cmd.Context(), args[0], args[1])
		},
	}
)

func init() {
	rootCmd.AddCommand(fileShareCmd)

	fileShareCmd.AddCommand(fileShareListCmd)
	fileShareListCmd.Flags().StringSliceVar(&fileShareFlags.filterStatus, "filter-status", nil, "Filter the result by file share status.")

	fileShareCmd.AddCommand(fileShareGetCmd)

	fileShareCmd.AddCommand(fileShareCreateCmd)
	fileShareCreateCmd.Flags().StringVar(&fileShareFlags.configSource, "config-source", "", "Source of the new file share configuration. Can be 'pipe' or path to a JSON file.")
	fileShareCreateCmd.MarkFlagsOneRequired("config-source")

	fileShareCmd.AddCommand(fileShareDeleteCmd)

	fileShareCmd.AddCommand(fileShareUpdateConfigCmd)
	fileShareUpdateConfigCmd.Flags().StringVar(&fileShareFlags.configSource, "config-source", "", "Source of the file share configuration updates. Can be 'pipe' or path to a JSON file.")
	fileShareUpdateConfigCmd.MarkFlagsOneRequired("config-source")

	fileShareCmd.AddCommand(fileShareUpdateMetaCmd)
	fileShareUpdateMetaCmd.Flags().StringVar(&fileShareFlags.configSource, "config-source", "", "Source of the file share metadata updates. Can be 'pipe' or path to a JSON file.")
	fileShareUpdateMetaCmd.MarkFlagsOneRequired("config-source")

	fileShareCmd.AddCommand(fileShareGetHostsCmd)

	fileShareCmd.AddCommand(fileShareUpdateHostsCmd)
	fileShareUpdateHostsCmd.Flags().StringVar(&fileShareFlags.configSource, "config-source", "", "Source of the file share hosts configuration. Can be 'pipe' or path to a JSON file.")
	fileShareUpdateHostsCmd.MarkFlagsOneRequired("config-source")

	fileShareCmd.AddCommand(fileShareGetConfigInfoCmd)
}
