package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/template"
	"github.com/spf13/cobra"
)

var (
	templateCmd = &cobra.Command{
		Use:     "template [command]",
		Aliases: []string{"templates"},
		Short:   "Template management",
		Long:    `Template management commands.`,
	}

	templateListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all templates.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return template.TemplateList(cmd.Context())
		},
	}

	templateGetCmd = &cobra.Command{
		Use:          "get template_id",
		Aliases:      []string{"show"},
		Short:        "Get template details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.TEMPLATES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return template.TemplateGet(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(templateCmd)

	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateGetCmd)
}
