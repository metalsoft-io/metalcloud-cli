package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/template_asset"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	templateAssetFlags = struct {
		configSource string
		templateId   string
		usage        string
		mimeType     string
	}{}

	templateAssetCmd = &cobra.Command{
		Use:     "template-asset [command]",
		Aliases: []string{"template-assets", "assets"},
		Short:   "Template asset management",
		Long:    `Template asset management commands.`,
	}

	templateAssetListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all template assets.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return template_asset.TemplateAssetList(
				cmd.Context(),
				templateAssetFlags.templateId,
				templateAssetFlags.usage,
				templateAssetFlags.mimeType)
		},
	}

	templateAssetGetCmd = &cobra.Command{
		Use:          "get template_asset_id",
		Aliases:      []string{"show"},
		Short:        "Get template asset details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return template_asset.TemplateAssetGet(cmd.Context(), args[0])
		},
	}

	templateAssetConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Get template asset configuration example.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return template_asset.TemplateAssetConfigExample(cmd.Context())
		},
	}

	templateAssetCreateCmd = &cobra.Command{
		Use:          "create",
		Aliases:      []string{"new"},
		Short:        "Create a template asset.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(templateAssetFlags.configSource)
			if err != nil {
				return err
			}

			return template_asset.TemplateAssetCreate(cmd.Context(), config)
		},
	}

	templateAssetUpdateCmd = &cobra.Command{
		Use:          "update template_asset_id",
		Aliases:      []string{"edit"},
		Short:        "Update a template asset.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(templateAssetFlags.configSource)
			if err != nil {
				return err
			}

			return template_asset.TemplateAssetUpdate(cmd.Context(), args[0], config)
		},
	}

	templateAssetDeleteCmd = &cobra.Command{
		Use:          "delete template_asset_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a template asset.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return template_asset.TemplateAssetDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(templateAssetCmd)

	// List command with filter options
	templateAssetCmd.AddCommand(templateAssetListCmd)
	templateAssetListCmd.Flags().StringVar(&templateAssetFlags.templateId, "template-id", "", "Filter assets by template ID.")
	templateAssetListCmd.Flags().StringVar(&templateAssetFlags.usage, "usage", "", "Filter assets by usage type (e.g., logo, icon, etc.).")
	templateAssetListCmd.Flags().StringVar(&templateAssetFlags.mimeType, "mime-type", "", "Filter assets by file MIME type (e.g., image/png, image/jpeg, etc.).")

	// Get command
	templateAssetCmd.AddCommand(templateAssetGetCmd)

	// Config example command
	templateAssetCmd.AddCommand(templateAssetConfigExampleCmd)

	// Create command
	templateAssetCmd.AddCommand(templateAssetCreateCmd)
	templateAssetCreateCmd.Flags().StringVar(&templateAssetFlags.configSource, "config-source", "", "Source of the new template asset configuration. Can be 'pipe' or path to a JSON file.")
	templateAssetCreateCmd.MarkFlagsOneRequired("config-source")

	// Update command
	templateAssetCmd.AddCommand(templateAssetUpdateCmd)
	templateAssetUpdateCmd.Flags().StringVar(&templateAssetFlags.configSource, "config-source", "", "Source of the template asset configuration updates. Can be 'pipe' or path to a JSON file.")
	templateAssetUpdateCmd.MarkFlagsOneRequired("config-source")

	// Delete command
	templateAssetCmd.AddCommand(templateAssetDeleteCmd)
}
