package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/secret"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	secretFlags = struct {
		configSource string
	}{}

	secretCmd = &cobra.Command{
		Use:     "secret [command]",
		Aliases: []string{"sec", "secrets"},
		Short:   "Secret management",
		Long:    `Secret management commands.`,
	}

	secretListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all secrets.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return secret.SecretList(cmd.Context())
		},
	}

	secretGetCmd = &cobra.Command{
		Use:          "get secret_id",
		Aliases:      []string{"show"},
		Short:        "Get secret details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return secret.SecretGet(cmd.Context(), args[0])
		},
	}

	secretConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Get secret configuration example.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return secret.SecretConfigExample(cmd.Context())
		},
	}

	secretCreateCmd = &cobra.Command{
		Use:          "create",
		Aliases:      []string{"new"},
		Short:        "Create a new secret.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(secretFlags.configSource)
			if err != nil {
				return err
			}

			return secret.SecretCreate(cmd.Context(), config)
		},
	}

	secretUpdateCmd = &cobra.Command{
		Use:          "update secret_id",
		Aliases:      []string{"edit"},
		Short:        "Update a secret.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(secretFlags.configSource)
			if err != nil {
				return err
			}

			return secret.SecretUpdate(cmd.Context(), args[0], config)
		},
	}

	secretDeleteCmd = &cobra.Command{
		Use:          "delete secret_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a secret.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return secret.SecretDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(secretCmd)

	secretCmd.AddCommand(secretListCmd)
	secretCmd.AddCommand(secretGetCmd)
	secretCmd.AddCommand(secretConfigExampleCmd)

	secretCmd.AddCommand(secretCreateCmd)
	secretCreateCmd.Flags().StringVar(&secretFlags.configSource, "config-source", "", "Source of the new secret configuration. Can be 'pipe' or path to a JSON file.")
	secretCreateCmd.MarkFlagsOneRequired("config-source")

	secretCmd.AddCommand(secretUpdateCmd)
	secretUpdateCmd.Flags().StringVar(&secretFlags.configSource, "config-source", "", "Source of the secret configuration updates. Can be 'pipe' or path to a JSON file.")
	secretUpdateCmd.MarkFlagsOneRequired("config-source")

	secretCmd.AddCommand(secretDeleteCmd)
}
