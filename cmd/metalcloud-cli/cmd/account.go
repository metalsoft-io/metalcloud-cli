package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/account"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	accountFlags = struct {
		configSource string
		userEmail    string
		roleLevel    string
		reason       string
	}{}

	accountCmd = &cobra.Command{
		Use:     "account [command]",
		Aliases: []string{"accounts"},
		Short:   "Account management",
		Long:    `Account management commands.`,
	}

	accountListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all accounts.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.USERS_AND_PERMISSIONS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return account.AccountList(cmd.Context())
		},
	}

	accountGetCmd = &cobra.Command{
		Use:          "get account_id",
		Aliases:      []string{"show"},
		Short:        "Get account details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.USERS_AND_PERMISSIONS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return account.AccountGet(cmd.Context(), args[0])
		},
	}

	accountCreateCmd = &cobra.Command{
		Use:          "create",
		Aliases:      []string{"new"},
		Short:        "Create a new account.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.USERS_AND_PERMISSIONS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(accountFlags.configSource)
			if err != nil {
				return err
			}

			return account.AccountCreate(cmd.Context(), config)
		},
	}

	accountUpdateCmd = &cobra.Command{
		Use:          "update account_id",
		Aliases:      []string{"edit"},
		Short:        "Update an account.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.USERS_AND_PERMISSIONS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(accountFlags.configSource)
			if err != nil {
				return err
			}

			return account.AccountUpdate(cmd.Context(), args[0], config)
		},
	}

	accountArchiveCmd = &cobra.Command{
		Use:          "archive account_id",
		Aliases:      []string{"ar"},
		Short:        "Archive an account.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.USERS_AND_PERMISSIONS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return account.AccountArchive(cmd.Context(), args[0])
		},
	}

	accountGetUsersCmd = &cobra.Command{
		Use:          "users account_id",
		Aliases:      []string{"get-users", "list-users"},
		Short:        "Get users for an account.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.USERS_AND_PERMISSIONS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return account.AccountGetUsers(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(accountCmd)

	accountCmd.AddCommand(accountListCmd)
	accountCmd.AddCommand(accountGetCmd)

	// Account create
	accountCmd.AddCommand(accountCreateCmd)
	accountCreateCmd.Flags().StringVar(&accountFlags.configSource, "config-source", "", "Source of the new account configuration. Can be 'pipe' or path to a JSON file.")
	accountCreateCmd.MarkFlagsOneRequired("config-source")

	// Account update
	accountCmd.AddCommand(accountUpdateCmd)
	accountUpdateCmd.Flags().StringVar(&accountFlags.configSource, "config-source", "", "Source of the account updates. Can be 'pipe' or path to a JSON file.")
	accountUpdateCmd.MarkFlagsOneRequired("config-source")

	accountCmd.AddCommand(accountArchiveCmd)

	// Account users
	accountCmd.AddCommand(accountGetUsersCmd)
}
