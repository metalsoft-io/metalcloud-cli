package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/user"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
)

var (
	userFlags = struct {
		configSource           string
		accountId              int
		sshKeyContent          string
		reason                 string
		archived               bool
		filterId               string
		filterDisplayName      string
		filterEmail            string
		filterAccountId        string
		filterInfrastructureId string
		sortBy                 string
		search                 string
		searchBy               string
		displayName            string
		email                  string
		password               string
		accessLevel            string
		emailVerified          bool
		createWithAccount      bool
	}{}

	userCmd = &cobra.Command{
		Use:     "user [command]",
		Aliases: []string{"users"},
		Short:   "User management",
		Long:    `User management commands.`,
	}

	userListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all users.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.List(
				cmd.Context(),
				userFlags.archived,
				userFlags.filterId,
				userFlags.filterDisplayName,
				userFlags.filterEmail,
				userFlags.filterAccountId,
				userFlags.filterInfrastructureId,
				userFlags.sortBy,
				userFlags.search,
				userFlags.searchBy,
			)
		},
	}

	userGetCmd = &cobra.Command{
		Use:          "get user_id",
		Aliases:      []string{"show"},
		Short:        "Get user details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.Get(cmd.Context(), args[0])
		},
	}

	userCreateCmd = &cobra.Command{
		Use:          "create",
		Aliases:      []string{"new"},
		Short:        "Create a new user.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			// If config source is provided, use it
			if userFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(userFlags.configSource)
				if err != nil {
					return err
				}
				return user.Create(cmd.Context(), config)
			}

			// Otherwise build config from command line parameters
			userConfig := sdk.CreateUser{
				Email:    userFlags.email,
				Password: sdk.PtrString(userFlags.password),
			}

			if cmd.Flags().Changed("email-verified") {
				userConfig.EmailVerified = sdk.PtrBool(userFlags.emailVerified)
			}

			if userFlags.accessLevel != "" {
				userConfig.AccessLevel = userFlags.accessLevel
			} else {
				userConfig.AccessLevel = "user"
			}

			if userFlags.displayName != "" {
				userConfig.DisplayName = userFlags.displayName
			} else {
				userConfig.DisplayName = userFlags.email
			}

			if userFlags.accountId != 0 {
				userConfig.AccountId = sdk.PtrFloat32(float32(userFlags.accountId))
			}

			if userFlags.createWithAccount {
				userConfig.CreateWithAccount = sdk.PtrBool(userFlags.createWithAccount)
			}

			configBytes, err := json.Marshal(userConfig)
			if err != nil {
				return fmt.Errorf("could not marshal user configuration: %s", err)
			}

			return user.Create(cmd.Context(), configBytes)
		},
	}

	userCreateBulkCmd = &cobra.Command{
		Use:          "create-bulk",
		Aliases:      []string{"bulk-create", "new-bulk"},
		Short:        "Create multiple users in a single operation.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Config source is required for bulk operations
			if userFlags.configSource == "" {
				return fmt.Errorf("config-source is required for bulk user creation")
			}

			config, err := utils.ReadConfigFromPipeOrFile(userFlags.configSource)
			if err != nil {
				return err
			}

			return user.CreateBulk(cmd.Context(), config)
		},
	}

	userArchiveCmd = &cobra.Command{
		Use:          "archive user_id",
		Aliases:      []string{"remove"},
		Short:        "Archive a user.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.Archive(cmd.Context(), args[0])
		},
	}

	userUnarchiveCmd = &cobra.Command{
		Use:          "unarchive user_id",
		Aliases:      []string{"restore"},
		Short:        "Unarchive a user.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.Unarchive(cmd.Context(), args[0])
		},
	}

	userLimitsGetCmd = &cobra.Command{
		Use:          "limits user_id",
		Aliases:      []string{"limits-get"},
		Short:        "Get user limits.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.GetLimits(cmd.Context(), args[0])
		},
	}

	userLimitsUpdateCmd = &cobra.Command{
		Use:          "limits-update user_id",
		Aliases:      []string{"update-limits"},
		Short:        "Update user limits.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(userFlags.configSource)
			if err != nil {
				return err
			}

			return user.UpdateLimits(cmd.Context(), args[0], config)
		},
	}

	userConfigUpdateCmd = &cobra.Command{
		Use:          "config-update user_id",
		Aliases:      []string{"update-config"},
		Short:        "Update user configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(userFlags.configSource)
			if err != nil {
				return err
			}

			return user.UpdateConfig(cmd.Context(), args[0], config)
		},
	}

	userChangeAccountCmd = &cobra.Command{
		Use:          "change-account user_id",
		Aliases:      []string{"move-account"},
		Short:        "Change user account.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.ChangeAccount(cmd.Context(), args[0], userFlags.accountId)
		},
	}

	userSshKeysGetCmd = &cobra.Command{
		Use:          "ssh-keys user_id",
		Aliases:      []string{"get-ssh-keys"},
		Short:        "Get user SSH keys.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.GetSSHKeys(cmd.Context(), args[0])
		},
	}

	userSshKeyAddCmd = &cobra.Command{
		Use:          "ssh-key-add user_id",
		Aliases:      []string{"add-ssh-key"},
		Short:        "Add an SSH key to a user.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.AddSSHKey(cmd.Context(), args[0], userFlags.sshKeyContent)
		},
	}

	userSshKeyDeleteCmd = &cobra.Command{
		Use:          "ssh-key-delete user_id key_id",
		Aliases:      []string{"delete-ssh-key", "remove-ssh-key"},
		Short:        "Delete an SSH key from a user.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.DeleteSSHKey(cmd.Context(), args[0], args[1])
		},
	}

	userSuspendCmd = &cobra.Command{
		Use:          "suspend user_id",
		Aliases:      []string{"disable"},
		Short:        "Suspend a user.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.Suspend(cmd.Context(), args[0], userFlags.reason)
		},
	}

	userUnsuspendCmd = &cobra.Command{
		Use:          "unsuspend user_id",
		Aliases:      []string{"enable"},
		Short:        "Unsuspend a user.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.Unsuspend(cmd.Context(), args[0])
		},
	}

	userPermissionsGetCmd = &cobra.Command{
		Use:          "permissions user_id",
		Aliases:      []string{"get-permissions"},
		Short:        "Get user permissions.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.GetPermissions(cmd.Context(), args[0])
		},
	}

	userPermissionsUpdateCmd = &cobra.Command{
		Use:          "permissions-update user_id",
		Aliases:      []string{"update-permissions"},
		Short:        "Update user permissions.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_AND_PERMISSIONS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(userFlags.configSource)
			if err != nil {
				return err
			}

			return user.UpdatePermissions(cmd.Context(), args[0], config)
		},
	}
)

func init() {
	rootCmd.AddCommand(userCmd)

	userCmd.AddCommand(userListCmd)
	userListCmd.Flags().BoolVar(&userFlags.archived, "archived", false, "Include archived users in the list.")
	userListCmd.Flags().StringVar(&userFlags.filterId, "filter-id", "", "Filter by user ID.")
	userListCmd.Flags().StringVar(&userFlags.filterDisplayName, "filter-display-name", "", "Filter by display name.")
	userListCmd.Flags().StringVar(&userFlags.filterEmail, "filter-email", "", "Filter by email.")
	userListCmd.Flags().StringVar(&userFlags.filterAccountId, "filter-account-id", "", "Filter by account ID.")
	userListCmd.Flags().StringVar(&userFlags.filterInfrastructureId, "filter-infrastructure-id", "", "Filter by default infrastructure ID.")
	userListCmd.Flags().StringVar(&userFlags.sortBy, "sort-by", "id:ASC", "Sort by field (e.g., 'id:ASC').")
	userListCmd.Flags().StringVar(&userFlags.search, "search", "", "Search term to filter results.")
	userListCmd.Flags().StringVar(&userFlags.searchBy, "search-by", "", "Fields to search by (e.g., 'displayName,email').")

	userCmd.AddCommand(userGetCmd)

	// User create
	userCmd.AddCommand(userCreateCmd)
	userCreateCmd.Flags().StringVar(&userFlags.configSource, "config-source", "", "Source of the new user configuration. Can be 'pipe' or path to a JSON file.")

	// Individual fields for user creation
	userCreateCmd.Flags().StringVar(&userFlags.email, "email", "", "User's email address")
	userCreateCmd.Flags().BoolVar(&userFlags.emailVerified, "email-verified", false, "Set the user email as verified")
	userCreateCmd.Flags().StringVar(&userFlags.password, "password", "", "User's password (if not provided, a random password will be generated)")
	userCreateCmd.Flags().StringVar(&userFlags.displayName, "display-name", "", "User's display name")
	userCreateCmd.Flags().StringVar(&userFlags.accessLevel, "access-level", "", "Access level (e.g., 'admin', 'user')")
	userCreateCmd.Flags().IntVar(&userFlags.accountId, "account-id", 0, "Account ID to associate the user with")
	userCreateCmd.Flags().BoolVar(&userFlags.createWithAccount, "create-with-account", false, "Create new account for the user")

	// Mark required fields that are mutually exclusive with config-source
	userCreateCmd.MarkFlagsMutuallyExclusive("config-source", "email")
	userCreateCmd.MarkFlagsRequiredTogether("email", "password")
	userCreateCmd.MarkFlagsMutuallyExclusive("account-id", "create-with-account")

	// User create bulk
	userCmd.AddCommand(userCreateBulkCmd)
	userCreateBulkCmd.Flags().StringVar(&userFlags.configSource, "config-source", "", "Source of the bulk user configuration. Can be 'pipe' or path to a JSON/YAML file with an array of user configs.")
	userCreateBulkCmd.MarkFlagRequired("config-source")

	// User archive/unarchive
	userCmd.AddCommand(userArchiveCmd)
	userCmd.AddCommand(userUnarchiveCmd)

	// User limits update
	userCmd.AddCommand(userLimitsGetCmd)
	userCmd.AddCommand(userLimitsUpdateCmd)
	userLimitsUpdateCmd.Flags().StringVar(&userFlags.configSource, "config-source", "", "Source of the user limits configuration. Can be 'pipe' or path to a JSON file.")
	userLimitsUpdateCmd.MarkFlagsOneRequired("config-source")

	// User config update
	userCmd.AddCommand(userConfigUpdateCmd)
	userConfigUpdateCmd.Flags().StringVar(&userFlags.configSource, "config-source", "", "Source of the user configuration. Can be 'pipe' or path to a JSON file.")
	userConfigUpdateCmd.MarkFlagsOneRequired("config-source")

	// Change account
	userCmd.AddCommand(userChangeAccountCmd)
	userChangeAccountCmd.Flags().IntVar(&userFlags.accountId, "account-id", 0, "The ID of the account to move the user to.")
	userChangeAccountCmd.MarkFlagRequired("account-id")

	// SSH Keys
	userCmd.AddCommand(userSshKeysGetCmd)
	userCmd.AddCommand(userSshKeyAddCmd)
	userSshKeyAddCmd.Flags().StringVar(&userFlags.sshKeyContent, "key", "", "The content of the SSH key.")
	userSshKeyAddCmd.MarkFlagRequired("key")
	userCmd.AddCommand(userSshKeyDeleteCmd)

	// Suspend/Unsuspend
	userCmd.AddCommand(userSuspendCmd)
	userSuspendCmd.Flags().StringVar(&userFlags.reason, "reason", "", "The reason for suspending the user.")
	userSuspendCmd.MarkFlagRequired("reason")
	userCmd.AddCommand(userUnsuspendCmd)

	// Permissions
	userCmd.AddCommand(userPermissionsGetCmd)
	userCmd.AddCommand(userPermissionsUpdateCmd)
	userPermissionsUpdateCmd.Flags().StringVar(&userFlags.configSource, "config-source", "", "Source of the user permissions configuration. Can be 'pipe' or path to a JSON file.")
	userPermissionsUpdateCmd.MarkFlagsOneRequired("config-source")
}
