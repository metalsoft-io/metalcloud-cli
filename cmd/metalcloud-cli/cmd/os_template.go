package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/os_template"
	"github.com/spf13/cobra"
)

var (
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
)

func init() {
	rootCmd.AddCommand(osTemplateCmd)

	osTemplateCmd.AddCommand(osTemplateListCmd)
	osTemplateCmd.AddCommand(osTemplateGetCmd)
}
