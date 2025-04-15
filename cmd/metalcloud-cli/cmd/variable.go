package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/variable"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	variableFlags = struct {
		configSource string
	}{}

	variableCmd = &cobra.Command{
		Use:     "variable [command]",
		Aliases: []string{"var", "vars"},
		Short:   "Variable management",
		Long:    `Variable management commands.`,
	}

	variableListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all variables.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return variable.VariableList(cmd.Context())
		},
	}

	variableGetCmd = &cobra.Command{
		Use:          "get variable_id",
		Aliases:      []string{"show"},
		Short:        "Get variable details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return variable.VariableGet(cmd.Context(), args[0])
		},
	}

	variableConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Get variable configuration example.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return variable.VariableConfigExample(cmd.Context())
		},
	}

	variableCreateCmd = &cobra.Command{
		Use:          "create",
		Aliases:      []string{"new"},
		Short:        "Create a new variable.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(variableFlags.configSource)
			if err != nil {
				return err
			}

			return variable.VariableCreate(cmd.Context(), config)
		},
	}

	variableUpdateCmd = &cobra.Command{
		Use:          "update variable_id",
		Aliases:      []string{"edit"},
		Short:        "Update a variable.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(variableFlags.configSource)
			if err != nil {
				return err
			}

			return variable.VariableUpdate(cmd.Context(), args[0], config)
		},
	}

	variableDeleteCmd = &cobra.Command{
		Use:          "delete variable_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a variable.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VARIABLES_AND_SECRETS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return variable.VariableDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(variableCmd)

	variableCmd.AddCommand(variableListCmd)
	variableCmd.AddCommand(variableGetCmd)
	variableCmd.AddCommand(variableConfigExampleCmd)

	variableCmd.AddCommand(variableCreateCmd)
	variableCreateCmd.Flags().StringVar(&variableFlags.configSource, "config-source", "", "Source of the new variable configuration. Can be 'pipe' or path to a JSON file.")
	variableCreateCmd.MarkFlagsOneRequired("config-source")

	variableCmd.AddCommand(variableUpdateCmd)
	variableUpdateCmd.Flags().StringVar(&variableFlags.configSource, "config-source", "", "Source of the variable configuration updates. Can be 'pipe' or path to a JSON file.")
	variableUpdateCmd.MarkFlagsOneRequired("config-source")

	variableCmd.AddCommand(variableDeleteCmd)
}
