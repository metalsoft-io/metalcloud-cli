package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/extension"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	extensionFlags = struct {
		definitionSource string
		filterLabel      string
		filterName       string
		filterStatus     string
	}{}

	extensionCmd = &cobra.Command{
		Use:     "extension [command]",
		Aliases: []string{"ext", "extensions"},
		Short:   "Extension management",
		Long:    `Extension management commands.`,
	}

	extensionListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all extensions.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSIONS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return extension.ExtensionList(
				cmd.Context(),
				extensionFlags.filterLabel,
				extensionFlags.filterName,
				extensionFlags.filterStatus,
			)
		},
	}

	extensionGetCmd = &cobra.Command{
		Use:          "get extension_id_or_label",
		Aliases:      []string{"show"},
		Short:        "Get extension info.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSIONS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return extension.ExtensionGet(cmd.Context(), args[0])
		},
	}

	extensionCreateCmd = &cobra.Command{
		Use:          "create name kind description",
		Aliases:      []string{"new"},
		Short:        "Create new extension.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSIONS_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			definition, err := utils.ReadConfigFromPipeOrFile(extensionFlags.definitionSource)
			if err != nil {
				return err
			}

			return extension.ExtensionCreate(cmd.Context(), args[0], args[1], args[2], definition)
		},
	}

	extensionUpdateCmd = &cobra.Command{
		Use:          "update extension_id_or_label [name [description]]",
		Aliases:      []string{"edit"},
		Short:        "Update extension configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSIONS_WRITE},
		Args:         cobra.RangeArgs(1, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := ""
			if len(args) > 1 {
				name = args[1]
			}

			description := ""
			if len(args) > 2 {
				description = args[2]
			}

			var definition []byte
			var err error
			if extensionFlags.definitionSource != "" {
				definition, err = utils.ReadConfigFromPipeOrFile(extensionFlags.definitionSource)
				if err != nil {
					return err
				}
			}

			return extension.ExtensionUpdate(cmd.Context(), args[0], name, description, definition)
		},
	}

	extensionPublishCmd = &cobra.Command{
		Use:          "publish extension_id_or_label",
		Short:        "Publish draft extension.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSIONS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return extension.ExtensionPublish(cmd.Context(), args[0])
		},
	}

	extensionArchiveCmd = &cobra.Command{
		Use:          "archive extension_id_or_label",
		Short:        "Archive published extension.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EXTENSIONS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return extension.ExtensionArchive(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(extensionCmd)

	extensionCmd.AddCommand(extensionListCmd)
	extensionListCmd.Flags().StringVar(&extensionFlags.filterLabel, "filter-label", "", "Filter extensions by label")
	extensionListCmd.Flags().StringVar(&extensionFlags.filterName, "filter-name", "", "Filter extensions by name")
	extensionListCmd.Flags().StringVar(&extensionFlags.filterStatus, "filter-status", "", "Filter extensions by status (DRAFT, PUBLISHED, ARCHIVED)")

	extensionCmd.AddCommand(extensionGetCmd)

	extensionCmd.AddCommand(extensionCreateCmd)
	extensionCreateCmd.Flags().StringVar(&extensionFlags.definitionSource, "definition-source", "", "Source of the extension definition. Can be 'pipe' or path to a JSON file.")
	extensionCreateCmd.MarkFlagRequired("definition-source")

	extensionCmd.AddCommand(extensionUpdateCmd)
	extensionUpdateCmd.Flags().StringVar(&extensionFlags.definitionSource, "definition-source", "", "Source of the updated extension definition. Can be 'pipe' or path to a JSON file.")

	extensionCmd.AddCommand(extensionPublishCmd)
	extensionCmd.AddCommand(extensionArchiveCmd)
}
