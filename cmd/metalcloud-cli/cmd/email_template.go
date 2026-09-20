package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/email_template"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	emailTemplateFlags = struct {
		configSource string
		filterName   []string
	}{}

	emailTemplateCmd = &cobra.Command{
		Use:     "email-template [command]",
		Aliases: []string{"email-templates"},
		Short:   "Email template management",
		Long: `Manage the email templates the platform uses for its notifications.

Templates are addressed by NAME on every command - they have no addressable
numeric ID.

Command categories:
  Read:    list, get
  Write:   create, update, delete
  Helper:  config-example`,
	}

	emailTemplateListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List email templates",
		Long: `List all email templates defined on the platform.

The table omits the text and html bodies because they are multi-line documents.
Use 'email-template get <name> -f json' to read them.

Optional Flags:
  --filter-name strings   Filter by template name. Repeatable or comma-separated.

Examples:
  metalcloud email-template list
  metalcloud email-templates ls --filter-name infrastructure-deployed`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EMAIL_TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return email_template.EmailTemplateList(cmd.Context(), emailTemplateFlags.filterName)
		},
	}

	emailTemplateGetCmd = &cobra.Command{
		Use:     "get email_template_name",
		Aliases: []string{"show"},
		Short:   "Get email template details",
		Long: `Get the details of an email template, including its subject and its text and
html bodies.

Required Arguments:
  email_template_name   The name of the email template

Examples:
  metalcloud email-template get infrastructure-deployed
  metalcloud email-template get infrastructure-deployed -f yaml`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EMAIL_TEMPLATES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return email_template.EmailTemplateGet(cmd.Context(), args[0])
		},
	}

	emailTemplateConfigExampleCmd = &cobra.Command{
		Use:     "config-example",
		Aliases: []string{"example"},
		Short:   "Print an example email template configuration",
		Long: `Print an example configuration that can be edited and passed to
'email-template create --config-source'.

Examples:
  metalcloud email-template config-example > template.json
  metalcloud email-template config-example -f yaml > template.yaml`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EMAIL_TEMPLATES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return email_template.EmailTemplateConfigExample(cmd.Context())
		},
	}

	emailTemplateCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create an email template",
		Long: `Create a new email template.

The template is described by a configuration document with the name, subject, text
and html fields (description is optional). Run 'email-template config-example' for
a starting point.

Required Flags:
  --config-source string   Source of the new template configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud email-template create --config-source template.json
  cat template.yaml | metalcloud email-template create --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EMAIL_TEMPLATES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(emailTemplateFlags.configSource)
			if err != nil {
				return err
			}

			return email_template.EmailTemplateCreate(cmd.Context(), config)
		},
	}

	emailTemplateUpdateCmd = &cobra.Command{
		Use:     "update email_template_name",
		Aliases: []string{"edit"},
		Short:   "Update an email template",
		Long: `Update the subject, description, text or html of an email template. Fields absent
from the configuration document are left unchanged. The template name cannot be
changed.

Required Arguments:
  email_template_name      The name of the email template

Required Flags:
  --config-source string   Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud email-template update infrastructure-deployed --config-source update.json
  echo '{"subject":"New subject"}' | metalcloud email-template update infrastructure-deployed --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EMAIL_TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(emailTemplateFlags.configSource)
			if err != nil {
				return err
			}

			return email_template.EmailTemplateUpdate(cmd.Context(), args[0], config)
		},
	}

	emailTemplateDeleteCmd = &cobra.Command{
		Use:     "delete email_template_name",
		Aliases: []string{"rm"},
		Short:   "Delete an email template",
		Long: `Delete an email template. Notifications that rely on it stop being sent until a
template with the same name is created again.

Required Arguments:
  email_template_name   The name of the email template

Examples:
  metalcloud email-template delete infrastructure-deployed
  metalcloud email-template rm infrastructure-deployed`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_EMAIL_TEMPLATES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return email_template.EmailTemplateDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(emailTemplateCmd)

	emailTemplateCmd.AddCommand(emailTemplateListCmd)
	emailTemplateListCmd.Flags().StringSliceVar(&emailTemplateFlags.filterName, "filter-name", nil, "Filter by template name.")

	emailTemplateCmd.AddCommand(emailTemplateGetCmd)

	emailTemplateCmd.AddCommand(emailTemplateConfigExampleCmd)

	emailTemplateCmd.AddCommand(emailTemplateCreateCmd)
	emailTemplateCreateCmd.Flags().StringVar(&emailTemplateFlags.configSource, "config-source", "", "Source of the new template configuration. Can be 'pipe' or path to a JSON/YAML file.")
	emailTemplateCreateCmd.MarkFlagRequired("config-source")

	emailTemplateCmd.AddCommand(emailTemplateUpdateCmd)
	emailTemplateUpdateCmd.Flags().StringVar(&emailTemplateFlags.configSource, "config-source", "", "Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.")
	emailTemplateUpdateCmd.MarkFlagRequired("config-source")

	emailTemplateCmd.AddCommand(emailTemplateDeleteCmd)
}
