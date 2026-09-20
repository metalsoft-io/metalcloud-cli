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
		archived     bool
		includeUsage bool
	}{}

	accountCmd = &cobra.Command{
		Use:     "account [command]",
		Aliases: []string{"accounts"},
		Short:   "Manage user accounts and account-related operations",
		Long: `Manage user accounts and account-related operations in the MetalCloud platform.

This command group provides functionality to:
- List all accounts in the system
- View detailed information about specific accounts
- Create new accounts with custom configurations
- Update existing account properties
- Archive accounts to disable them
- List users associated with an account

All account operations require appropriate permissions to perform user management tasks.`,
	}

	accountListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all accounts in the system",
		Long: `List all accounts in the MetalCloud platform.

This command displays a table of all accounts including their ID, name, status, 
creation date, and other relevant information. The output can be formatted as 
JSON, YAML, or table format.

Required Permissions:
  - users:read

Examples:
  # List all accounts in table format
  metalcloud-cli account list

  # List all accounts in JSON format  
  metalcloud-cli account list -o json

  # List all accounts using alias
  metalcloud-cli accounts ls

  # List archived accounts (excluded by default)
  metalcloud-cli account list --archived`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return account.AccountList(cmd.Context(), accountFlags.archived)
		},
	}

	accountGetCmd = &cobra.Command{
		Use:     "get account_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific account",
		Long: `Get detailed information about a specific account in the MetalCloud platform.

This command displays comprehensive details about an account including its 
configuration, status, creation date, and associated metadata. The account
is identified by its unique account ID.

Required Permissions:
  - users:read

Arguments:
  account_id    The unique identifier of the account to retrieve

Examples:
  # Get account details by ID
  metalcloud-cli account get 1234

  # Get account details in JSON format
  metalcloud-cli account get 1234 -o json

  # Using alias
  metalcloud-cli account show 1234`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return account.AccountGet(cmd.Context(), args[0])
		},
	}

	accountCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new account with custom configuration",
		Long: `Create a new account in the MetalCloud platform with custom configuration.

This command creates a new account using configuration data provided through 
a JSON file or standard input. The configuration must include all required 
account properties such as name, description, and any custom settings.

Required Permissions:
  - users:write

Required Flags:
  --config-source    Source of the new account configuration

Flag Details:
  --config-source string    Source of the new account configuration. 
                           Can be 'pipe' to read from stdin or path to a JSON file.
                           The JSON should contain account properties like name, 
                           description, and other account settings.

Examples:
  # Create account from JSON file
  metalcloud-cli account create --config-source /path/to/account.json

  # Create account from stdin
  echo '{"name":"test-account","description":"Test account"}' | metalcloud-cli account create --config-source pipe

  # Using alias
  metalcloud-cli account new --config-source account-config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(accountFlags.configSource)
			if err != nil {
				return err
			}

			return account.AccountCreate(cmd.Context(), config)
		},
	}

	accountUpdateCmd = &cobra.Command{
		Use:     "update account_id",
		Aliases: []string{"edit"},
		Short:   "Update an existing account configuration",
		Long: `Update an existing account in the MetalCloud platform with new configuration.

This command updates an existing account using configuration data provided through 
a JSON file or standard input. Only the properties specified in the configuration 
will be updated; other properties remain unchanged. The account is identified 
by its unique account ID.

Required Permissions:
  - users:write

Arguments:
  account_id    The unique identifier of the account to update

Required Flags:
  --config-source    Source of the account update configuration

Flag Details:
  --config-source string    Source of the account update configuration. 
                           Can be 'pipe' to read from stdin or path to a JSON file.
                           The JSON should contain the account properties to update.

Examples:
  # Update account from JSON file
  metalcloud-cli account update 1234 --config-source /path/to/updates.json

  # Update account from stdin
  echo '{"description":"Updated description"}' | metalcloud-cli account update 1234 --config-source pipe

  # Using alias
  metalcloud-cli account edit 1234 --config-source account-updates.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
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
		Use:     "archive account_id",
		Aliases: []string{"ar"},
		Short:   "Archive an account to disable it",
		Long: `Archive an account in the MetalCloud platform to disable it.

This command archives an account, effectively disabling it while preserving 
its data and configuration. Archived accounts cannot be used for new operations 
but their historical data remains accessible. The account is identified by its 
unique account ID.

Note: This operation is typically irreversible. Archived accounts may require 
administrator intervention to reactivate.

Required Permissions:
  - users:write

Arguments:
  account_id    The unique identifier of the account to archive

Examples:
  # Archive an account by ID
  metalcloud-cli account archive 1234

  # Using alias
  metalcloud-cli account ar 1234`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return account.AccountArchive(cmd.Context(), args[0])
		},
	}

	accountGetUsersCmd = &cobra.Command{
		Use:     "users account_id",
		Aliases: []string{"get-users", "list-users"},
		Short:   "List all users associated with a specific account",
		Long: `List all users associated with a specific account in the MetalCloud platform.

This command displays a list of users that belong to the specified account, 
including their user details, roles, and permissions within that account. 
The account is identified by its unique account ID.

Required Permissions:
  - users:read

Arguments:
  account_id    The unique identifier of the account to list users for

Examples:
  # List users for a specific account
  metalcloud-cli account users 1234

  # List users in JSON format
  metalcloud-cli account users 1234 -o json

  # Using aliases
  metalcloud-cli account get-users 1234
  metalcloud-cli account list-users 1234`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return account.AccountGetUsers(cmd.Context(), args[0])
		},
	}

	accountUnarchiveCmd = &cobra.Command{
		Use:     "unarchive account_id",
		Aliases: []string{"restore", "unar"},
		Short:   "Restore a previously archived account",
		Long: `Restore a previously archived account in the MetalCloud platform.

This command reverses 'account archive': the account becomes usable again while all
its data and configuration are preserved. The account is identified by its unique
account ID.

Required Permissions:
  - users:write

Arguments:
  account_id    The unique identifier of the account to restore

Examples:
  # Restore an account by ID
  metalcloud-cli account unarchive 1234

  # Using aliases
  metalcloud-cli account restore 1234
  metalcloud-cli account unar 1234`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return account.AccountUnarchive(cmd.Context(), args[0])
		},
	}

	accountQuotaBreakdownCmd = &cobra.Command{
		Use:     "quota-breakdown account_id",
		Aliases: []string{"quota", "quota-limits"},
		Short:   "Show how the quota limits of an account are derived",
		Long: `Show the quota limits breakdown of an account in the MetalCloud platform.

The breakdown lists, for every quota limit, the effective value together with the
value contributed by the quota profile of the account and by the quota profile of its
parent account. The effective value is the most restrictive of the two; individual
users may still be further restricted by their role or group quota profiles.

Text, CSV and Markdown output flatten the breakdown into one row per limit. JSON and
YAML output keep the original nested object returned by the API.

Required Permissions:
  - users:read

Arguments:
  account_id    The unique identifier of the account

Optional Flags:
  --include-usage    Also report the aggregate resource counts of the account and of
                     the parent account (adds two columns to the table)

Examples:
  # Show the quota breakdown of an account
  metalcloud-cli account quota-breakdown 1234

  # Include the current usage
  metalcloud-cli account quota-breakdown 1234 --include-usage

  # Keep the raw object
  metalcloud-cli account quota 1234 -f json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return account.AccountQuotaBreakdown(cmd.Context(), args[0], accountFlags.includeUsage)
		},
	}

	accountConfigGetCmd = &cobra.Command{
		Use:     "config account_id",
		Aliases: []string{"get-config"},
		Short:   "Display the configuration of an account",
		Long: `Display the configuration object of a specific account.

The configuration holds the properties written by 'account update': name, code, fiscal
number, contacts, parent account, quota profile and API key validity duration.

Required Permissions:
  - users:read

Arguments:
  account_id    The unique identifier of the account

Examples:
  # Show the configuration of an account
  metalcloud-cli account config 1234

  # Using alias
  metalcloud-cli account get-config 1234`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return account.AccountGetConfig(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(accountCmd)

	accountCmd.AddCommand(accountListCmd)
	accountListCmd.Flags().BoolVar(&accountFlags.archived, "archived", false, "List archived accounts (excluded by default).")
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
	accountCmd.AddCommand(accountUnarchiveCmd)

	// Account users
	accountCmd.AddCommand(accountGetUsersCmd)

	// Account quota limits breakdown
	accountCmd.AddCommand(accountQuotaBreakdownCmd)
	accountQuotaBreakdownCmd.Flags().BoolVar(&accountFlags.includeUsage, "include-usage", false, "Also report the aggregate resource usage of the account and of its parent account.")

	// Account configuration
	accountCmd.AddCommand(accountConfigGetCmd)
}
