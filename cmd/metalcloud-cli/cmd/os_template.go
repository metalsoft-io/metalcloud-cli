package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/os_template"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	osTemplateFlags = struct {
		configSource string
		deviceType   string
		visibility   string
	}{}

	osTemplateCmd = &cobra.Command{
		Use:     "os-template [command]",
		Aliases: []string{"templates"},
		Short:   "OS template management",
		Long:    `OS template management commands.`,
	}

	osTemplateListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all OS templates.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return os_template.OsTemplateList(cmd.Context())
		},
	}

	osTemplateGetCmd = &cobra.Command{
		Use:          "get os_template_id",
		Aliases:      []string{"show"},
		Short:        "Get OS template details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.TEMPLATES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return os_template.OsTemplateGet(cmd.Context(), args[0])
		},
	}

	osTemplateCreateCmd = &cobra.Command{
		Use:          "create",
		Aliases:      []string{"new"},
		Short:        "Create a new OS template.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.TEMPLATES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(osTemplateFlags.configSource)
			if err != nil {
				return err
			}

			return os_template.OsTemplateCreate(cmd.Context(), config)
		},
	}

	osTemplateUpdateCmd = &cobra.Command{
		Use:          "update os_template_id",
		Aliases:      []string{"edit"},
		Short:        "Update an OS template.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(osTemplateFlags.configSource)
			if err != nil {
				return err
			}

			return os_template.OsTemplateUpdate(cmd.Context(), args[0], config)
		},
	}

	osTemplateDeleteCmd = &cobra.Command{
		Use:          "delete os_template_id",
		Aliases:      []string{"rm"},
		Short:        "Delete an OS template.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return os_template.OsTemplateDelete(cmd.Context(), args[0])
		},
	}

	osTemplateGetCredentialsCmd = &cobra.Command{
		Use:          "get-credentials os_template_id",
		Aliases:      []string{"creds"},
		Short:        "Get credentials for an OS template.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.TEMPLATES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return os_template.OsTemplateGetCredentials(cmd.Context(), args[0])
		},
	}

	osTemplateGetAssetsCmd = &cobra.Command{
		Use:          "get-assets os_template_id",
		Aliases:      []string{"assets"},
		Short:        "Get assets for an OS template.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.TEMPLATES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return os_template.OsTemplateGetAssets(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(osTemplateCmd)

	osTemplateCmd.AddCommand(osTemplateListCmd)
	osTemplateCmd.AddCommand(osTemplateGetCmd)

	osTemplateCmd.AddCommand(osTemplateCreateCmd)
	osTemplateCreateCmd.Flags().StringVar(&osTemplateFlags.configSource, "config-source", "", "Source of the new OS template configuration. Can be 'pipe' or path to a JSON file.")
	osTemplateCreateCmd.MarkFlagsOneRequired("config-source")

	osTemplateCmd.AddCommand(osTemplateUpdateCmd)
	osTemplateUpdateCmd.Flags().StringVar(&osTemplateFlags.configSource, "config-source", "", "Source of the OS template configuration updates. Can be 'pipe' or path to a JSON file.")
	osTemplateUpdateCmd.MarkFlagsOneRequired("config-source")

	osTemplateCmd.AddCommand(osTemplateDeleteCmd)
	osTemplateCmd.AddCommand(osTemplateGetCredentialsCmd)
	osTemplateCmd.AddCommand(osTemplateGetAssetsCmd)
}
