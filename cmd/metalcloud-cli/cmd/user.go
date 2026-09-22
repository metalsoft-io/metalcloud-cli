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
		twoFAToken             string
		redirectUrl            string
		newPassword            string
		currentPassword        string
		callbackToken          string
	}{}

	userCmd = &cobra.Command{
		Use:     "user [command]",
		Aliases: []string{"users"},
		Short:   "Manage user accounts and their properties",
		Long: `Comprehensive user management commands for creating, modifying, and managing user accounts.
These commands allow you to perform various operations on user accounts including:
- Creating individual or bulk users
- Managing user permissions and limits
- Handling SSH keys and authentication
- User lifecycle operations (archive/unarchive, suspend/unsuspend)
- Account management and configuration updates

All commands require appropriate permissions and most modification commands require
the user ID as a parameter. Use 'metalcloud-cli user list' to find user IDs.`,
	}

	userListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List users with filtering and search options",
		Long: `List all users in the system with advanced filtering, searching, and sorting capabilities.

This command displays user information including ID, name, email, access level, archived status,
creation date, and last login timestamp. Users can be filtered by various criteria and results
can be sorted and searched through.

Filters:
  --archived              Return only archived users
  --filter-id             Filter by specific user ID
  --filter-display-name   Filter by user display name (partial matches)
  --filter-email          Filter by email address (partial matches)
  --filter-account-id     Filter by account ID
  --filter-infrastructure-id Filter by default infrastructure ID

Search and Sort:
  --search                Search term applied across multiple fields
  --search-by             Specify which fields to search in (comma-separated)
  --sort-by               Sort results by field and direction (e.g., "id:ASC", "email:DESC")

Examples:
  metalcloud-cli user list
  metalcloud-cli user list --archived
  metalcloud-cli user list --filter-email "@company.com"
  metalcloud-cli user list --search "john" --search-by "displayName,email"
  metalcloud-cli user list --sort-by "createdTimestamp:DESC"`,
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
		Use:     "get user_id",
		Aliases: []string{"show"},
		Short:   "Display detailed information for a specific user",
		Long: `Retrieve and display comprehensive information for a specific user account.

This command shows all available user details including personal information, account settings,
access levels, timestamps, and status flags. The user ID is required and can be found using
the 'user list' command.

Arguments:
  user_id                 The numeric ID of the user to display

Examples:
  metalcloud-cli user get 12345
  metalcloud-cli user show 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.Get(cmd.Context(), args[0])
		},
	}

	userCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new user account with specified properties",
		Long: `Create a new user account in the system with comprehensive configuration options.

This command allows creating users either through individual command-line flags or by providing
a JSON configuration file/pipe. The user can be associated with an existing account or a new
account can be created automatically.

Required Flags (when not using --config-source):
  --email                 User's email address (required, used as login)
  --password              User's password (required for CLI creation)

Optional Flags:
  --config-source         Source of user configuration (JSON file path or 'pipe')
  --display-name          User's display name (defaults to email if not provided)
  --access-level          User access level: admin, user, readonly (default: user)
  --email-verified        Mark user email as verified (default: false)
  --account-id            Associate user with existing account ID
  --create-with-account   Create a new account for the user (mutually exclusive with --account-id)

Dependencies:
  - --email and --password are required together when not using --config-source
  - --config-source is mutually exclusive with --email
  - --account-id and --create-with-account are mutually exclusive

Configuration File Format (JSON):
  {
    "displayName": "John Doe",
    "email": "john.doe@company.com",
    "password": "securePassword123",
    "accessLevel": "user",
    "emailVerified": true,
    "accountId": 12345
  }`,
		Example: `  # Create user with command-line flags
  metalcloud-cli user create --email test.user@metalsoft.io --password secret --access-level user
  
  # Create user with additional properties
  metalcloud-cli user create --email test.user@metalsoft.io --password secret --access-level user --display-name "Test User" --email-verified true --account-id 12345
  
  # Create user with new account
  metalcloud-cli user create --email admin@company.com --password admin123 --access-level admin --create-with-account
  
  # Create user from JSON file
  metalcloud-cli user create --config-source user1.json
  
  # Create user from pipe
  echo '{"email": "test.user@metalsoft.io", "password": "secret", "accessLevel": "user", "displayName": "Test User"}' | metalcloud-cli user create --config-source pipe`,
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
				userConfig.AccountId = sdk.PtrInt64(int64(userFlags.accountId))
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
		Use:     "create-bulk",
		Aliases: []string{"bulk-create", "new-bulk"},
		Short:   "Create multiple users in a single operation",
		Long: `Create multiple users at once from a JSON or YAML configuration file.

This command processes an array of user configurations and creates all users in sequence.
Each user configuration follows the same format as the single user create command.

Required Flags:
  --config-source         Source of bulk user configuration (JSON/YAML file path or 'pipe')

Configuration File Format (JSON):
  [
    {
      "displayName": "John Doe", 
      "email": "john.doe@company.com",
      "password": "securePassword123",
      "accessLevel": "user",
      "emailVerified": true,
      "accountId": 12345
    },
    {
      "displayName": "Jane Smith",
      "email": "jane.smith@company.com", 
      "password": "anotherPassword456",
      "accessLevel": "admin",
      "createWithAccount": true
    }
  ]

The command will report success/failure for each user and provide a summary at the end.`,
		Example: `  # Create users from JSON file
  metalcloud-cli user create-bulk --config-source users.json
  
  # Create users from YAML file  
  metalcloud-cli user create-bulk --config-source users.yaml
  
  # Create users from pipe
  echo '[{"email": "user1@company.com", "password": "pass1", "accessLevel": "user"}, {"email": "user2@company.com", "password": "pass2", "accessLevel": "admin"}]' | metalcloud-cli user create-bulk --config-source pipe`,
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
		Use:     "archive user_id",
		Aliases: []string{"remove"},
		Short:   "Archive a user account to mark it as inactive",
		Long: `Archive a user account to mark it as inactive and prevent future logins.

Archiving a user preserves their data but prevents them from logging in or accessing
the system. This is a reversible action - archived users can be unarchived later.

Arguments:
  user_id                 The numeric ID of the user to archive

Examples:
  metalcloud-cli user archive 12345
  metalcloud-cli user remove 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.Archive(cmd.Context(), args[0])
		},
	}

	userUnarchiveCmd = &cobra.Command{
		Use:     "unarchive user_id",
		Aliases: []string{"restore"},
		Short:   "Unarchive a user account to restore access",
		Long: `Unarchive a previously archived user account to restore their access to the system.

This command reverses the archive operation, allowing the user to log in and access
the system again. All user data and settings are preserved during archive/unarchive.

Arguments:
  user_id                 The numeric ID of the user to unarchive

Examples:
  metalcloud-cli user unarchive 12345
  metalcloud-cli user restore 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.Unarchive(cmd.Context(), args[0])
		},
	}

	userLimitsGetCmd = &cobra.Command{
		Use:     "limits user_id",
		Aliases: []string{"limits-get"},
		Short:   "Display resource limits for a specific user",
		Long: `Retrieve and display the resource limits configured for a specific user account.

This command shows limits for compute nodes, drives, and infrastructures that the user
can provision. These limits control resource allocation and prevent overuse.

Arguments:
  user_id                 The numeric ID of the user whose limits to display

Examples:
  metalcloud-cli user limits 12345
  metalcloud-cli user limits-get 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.GetLimits(cmd.Context(), args[0])
		},
	}

	userConfigUpdateCmd = &cobra.Command{
		Use:     "config-update user_id",
		Aliases: []string{"update-config"},
		Short:   "Update comprehensive user configuration settings",
		Long: `Update comprehensive configuration settings for a specific user account.

This command allows updating various user properties including display name, email, access level,
and other account settings. The configuration is provided through a JSON file or pipe.

Arguments:
  user_id                 The numeric ID of the user whose configuration to update

Required Flags:
  --config-source         Source of user configuration (JSON file path or 'pipe')

Configuration File Format (JSON):
  {
    "displayName": "Updated Name",
    "accessLevel": "admin",
    "emailVerified": true,
    "language": "en"
  }

Examples:
  metalcloud-cli user config-update 12345 --config-source config.json
  echo '{"displayName": "New Name", "accessLevel": "admin"}' | metalcloud-cli user config-update 12345 --config-source pipe`,
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
		Use:     "change-account user_id",
		Aliases: []string{"move-account"},
		Short:   "Move a user to a different account",
		Long: `Move a user from their current account to a different account in the system.

This command transfers user ownership between accounts while preserving all user data,
settings, and permissions. The user will be associated with the new account immediately
after the command executes successfully.

Arguments:
  user_id                 The numeric ID of the user to move

Required Flags:
  --account-id            The ID of the destination account to move the user to

Examples:
  metalcloud-cli user change-account 12345 --account-id 67890
  metalcloud-cli user move-account 12345 --account-id 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.ChangeAccount(cmd.Context(), args[0], userFlags.accountId)
		},
	}

	userSshKeysGetCmd = &cobra.Command{
		Use:     "ssh-keys user_id",
		Aliases: []string{"get-ssh-keys"},
		Short:   "Display SSH keys for a specific user",
		Long: `Retrieve and display all SSH keys associated with a specific user account.

This command shows SSH key information including key ID, name, fingerprint, and
creation timestamp. SSH keys are used for authentication when connecting to instances.

Arguments:
  user_id                 The numeric ID of the user whose SSH keys to display

Examples:
  metalcloud-cli user ssh-keys 12345
  metalcloud-cli user get-ssh-keys 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.GetSSHKeys(cmd.Context(), args[0])
		},
	}

	userSshKeyAddCmd = &cobra.Command{
		Use:     "ssh-key-add user_id",
		Aliases: []string{"add-ssh-key"},
		Short:   "Add an SSH key to a user account",
		Long: `Add a new SSH key to a specific user account for authentication purposes.

This command allows adding SSH public keys to user accounts which can then be used for
authentication when connecting to provisioned instances. The SSH key content should be
a valid public key in OpenSSH format.

Arguments:
  user_id                 The numeric ID of the user to add the SSH key to

Required Flags:
  --key                   The SSH public key content (OpenSSH format)

Examples:
  metalcloud-cli user ssh-key-add 12345 --key "ssh-rsa AAAAB3NzaC1yc2EAAAA..."
  metalcloud-cli user add-ssh-key 12345 --key "$(cat ~/.ssh/id_rsa.pub)"`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.AddSSHKey(cmd.Context(), args[0], userFlags.sshKeyContent)
		},
	}

	userSshKeyDeleteCmd = &cobra.Command{
		Use:     "ssh-key-delete user_id key_id",
		Aliases: []string{"delete-ssh-key", "remove-ssh-key"},
		Short:   "Delete an SSH key from a user account",
		Long: `Remove an existing SSH key from a specific user account.

This command permanently deletes an SSH key from the user's account. Once deleted,
the key can no longer be used for authentication to instances.

Arguments:
  user_id                 The numeric ID of the user whose SSH key to delete
  key_id                  The numeric ID of the SSH key to delete

Examples:
  metalcloud-cli user ssh-key-delete 12345 67890
  metalcloud-cli user delete-ssh-key 12345 67890
  metalcloud-cli user remove-ssh-key 12345 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.DeleteSSHKey(cmd.Context(), args[0], args[1])
		},
	}

	userSuspendCmd = &cobra.Command{
		Use:     "suspend user_id",
		Aliases: []string{"disable"},
		Short:   "Suspend a user account temporarily",
		Long: `Suspend a user account to temporarily prevent access while preserving all data.

Suspending a user prevents them from logging in and accessing the system, but unlike
archiving, this is typically used for temporary restrictions. A reason for suspension
is required for auditing purposes.

Arguments:
  user_id                 The numeric ID of the user to suspend

Required Flags:
  --reason                The reason for suspending the user (required for audit trail)

Examples:
  metalcloud-cli user suspend 12345 --reason "Policy violation"
  metalcloud-cli user disable 12345 --reason "Account under review"`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.Suspend(cmd.Context(), args[0], userFlags.reason)
		},
	}

	userUnsuspendCmd = &cobra.Command{
		Use:     "unsuspend user_id",
		Aliases: []string{"enable"},
		Short:   "Unsuspend a user account to restore access",
		Long: `Unsuspend a previously suspended user account to restore their access to the system.

This command reverses the suspend operation, allowing the user to log in and access
the system again. All user data and settings are preserved during suspend/unsuspend.

Arguments:
  user_id                 The numeric ID of the user to unsuspend

Examples:
  metalcloud-cli user unsuspend 12345
  metalcloud-cli user enable 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.Unsuspend(cmd.Context(), args[0])
		},
	}

	userApiKeyGetCmd = &cobra.Command{
		Use:   "api-key",
		Short: "Get the current user's API key",
		Long: `Retrieve the API key for the currently authenticated user.

Examples:
  metalcloud-cli user api-key`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.GetApiKey(cmd.Context())
		},
	}

	userApiKeyRegenerateCmd = &cobra.Command{
		Use:   "api-key-regenerate",
		Short: "Regenerate the current user's API key",
		Long: `Regenerate the API key for the currently authenticated user.

WARNING: This will invalidate your current API key. You will need to update
any scripts or configurations that use the old key.

Examples:
  metalcloud-cli user api-key-regenerate`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.RegenerateApiKey(cmd.Context())
		},
	}

	user2FAEnableCmd = &cobra.Command{
		Use:   "2fa-enable",
		Short: "Enable two-factor authentication",
		Long: `Enable two-factor authentication for the current user.

You must first generate a 2FA secret using 'user 2fa-generate-secret', configure
your authenticator app, then provide the TOTP token to verify and enable 2FA.

Required Flags:
  --token string    The TOTP code from your authenticator app

Examples:
  # First generate the secret
  metalcloud-cli user 2fa-generate-secret

  # Then enable 2FA with the token from your authenticator app
  metalcloud-cli user 2fa-enable --token 123456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.Enable2FA(cmd.Context(), userFlags.twoFAToken)
		},
	}

	user2FADisableCmd = &cobra.Command{
		Use:   "2fa-disable",
		Short: "Disable two-factor authentication",
		Long: `Disable two-factor authentication for the current user.

Examples:
  metalcloud-cli user 2fa-disable`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.Disable2FA(cmd.Context())
		},
	}

	userSetPasswordCmd = &cobra.Command{
		Use:   "set-password user_id",
		Short: "Set the password for a user (admin only)",
		Long: `Set the password for a specific user account as an administrator.

Arguments:
  user_id                 The numeric ID of the user whose password to set

Required Flags:
  --password              The new password to set for the user

Examples:
  metalcloud-cli user set-password 12345 --password newSecret123`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.SetPassword(cmd.Context(), args[0], userFlags.password)
		},
	}

	user2FAGenerateSecretCmd = &cobra.Command{
		Use:   "2fa-generate-secret",
		Short: "Generate a new 2FA secret for setting up an authenticator app",
		Long: `Generate a new two-factor authentication secret and QR code.

Use the generated secret or QR code to configure your authenticator app (e.g., Google
Authenticator, Authy), then call 'user 2fa-enable --token <code>' to activate 2FA.

Examples:
  metalcloud-cli user 2fa-generate-secret`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.GenerateUser2FASecret(cmd.Context())
		},
	}

	userConfigGetCmd = &cobra.Command{
		Use:     "config user_id",
		Aliases: []string{"get-config"},
		Short:   "Display the configuration of a user",
		Long: `Display the configuration object of a specific user account.

The configuration holds the settings written by 'user config-update': display name,
access level, language, brand, login state and password policy flags.

Required Arguments:
  user_id                 The numeric ID of the user whose configuration to display

Examples:
  metalcloud-cli user config 12345
  metalcloud-cli user get-config 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.GetConfiguration(cmd.Context(), args[0])
		},
	}

	userUpdateMetaCmd = &cobra.Command{
		Use:     "update-meta user_id",
		Aliases: []string{"meta-update"},
		Short:   "Update the metadata of a user",
		Long: `Update the metadata of a specific user account.

User metadata carries the GUI settings of the user. The payload replaces the stored
metadata, so provide the complete object.

Required Arguments:
  user_id                 The numeric ID of the user whose metadata to update

Required Flags:
  --config-source         Source of the metadata (JSON/YAML file path or 'pipe')

Configuration File Format (JSON):
  {
    "guiSettings": {
      "defaultPage": "infrastructures"
    }
  }

Examples:
  metalcloud-cli user update-meta 12345 --config-source meta.json
  echo '{"guiSettings":{}}' | metalcloud-cli user update-meta 12345 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(userFlags.configSource)
			if err != nil {
				return err
			}

			return user.UpdateMeta(cmd.Context(), args[0], config)
		},
	}

	userSshKeyGetCmd = &cobra.Command{
		Use:     "ssh-key user_id key_id",
		Aliases: []string{"get-ssh-key"},
		Short:   "Display a single SSH key of a user",
		Long: `Display one SSH key of a specific user account.

Use 'user ssh-keys' to list the SSH keys of the user and obtain their IDs.

Required Arguments:
  user_id                 The numeric ID of the user owning the SSH key
  key_id                  The numeric ID of the SSH key to display

Examples:
  metalcloud-cli user ssh-key 12345 67890
  metalcloud-cli user get-ssh-key 12345 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.GetSSHKey(cmd.Context(), args[0], args[1])
		},
	}

	userSuspendReasonsCmd = &cobra.Command{
		Use:     "suspend-reasons user_id",
		Aliases: []string{"get-suspend-reasons"},
		Short:   "List the suspend reasons recorded for a user",
		Long: `List the suspend reasons recorded for a specific user account.

Each entry shows the type of the suspension, the public and private comments left by
the administrator and the interval during which the suspension was active.

Required Arguments:
  user_id                 The numeric ID of the user whose suspend reasons to list

Examples:
  metalcloud-cli user suspend-reasons 12345
  metalcloud-cli user get-suspend-reasons 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.GetSuspendReasons(cmd.Context(), args[0])
		},
	}

	userDelegateCmd = &cobra.Command{
		Use:     "delegate [command]",
		Aliases: []string{"delegates"},
		Short:   "Manage user delegates",
		Long: `Manage the delegation relationships between user accounts.

A delegate user can act on the resources of the user that delegated access to them.
These commands allow you to:
- Grant delegate access to another user (add)
- Revoke delegate access (remove)
- List the users that delegated access to a user (parents)
- List the users a user delegated access to (children)`,
	}

	userDelegateAddCmd = &cobra.Command{
		Use:     "add user_id delegate_id",
		Aliases: []string{"new"},
		Short:   "Grant a user delegate access to another user",
		Long: `Grant a delegate user access to the resources of a user.

After this command the delegate user can operate on the resources owned by the user
identified by user_id.

Required Arguments:
  user_id                 The numeric ID of the user whose resources are delegated
  delegate_id             The numeric ID of the user receiving the delegate access

Examples:
  metalcloud-cli user delegate add 12345 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.AddDelegate(cmd.Context(), args[0], args[1])
		},
	}

	userDelegateRemoveCmd = &cobra.Command{
		Use:     "remove user_id delegate_id",
		Aliases: []string{"rm", "delete"},
		Short:   "Revoke the delegate access of a user",
		Long: `Revoke the delegate access previously granted to a user.

Required Arguments:
  user_id                 The numeric ID of the user whose resources were delegated
  delegate_id             The numeric ID of the user losing the delegate access

Examples:
  metalcloud-cli user delegate remove 12345 67890
  metalcloud-cli user delegate rm 12345 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.RemoveDelegate(cmd.Context(), args[0], args[1])
		},
	}

	userDelegateParentsCmd = &cobra.Command{
		Use:     "parents user_id",
		Aliases: []string{"parent-delegates"},
		Short:   "List the users that delegated access to a user",
		Long: `List the users that delegated their resources to a specific user.

Required Arguments:
  user_id                 The numeric ID of the delegate user

Examples:
  metalcloud-cli user delegate parents 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.GetParentDelegates(cmd.Context(), args[0])
		},
	}

	userDelegateChildrenCmd = &cobra.Command{
		Use:     "children user_id",
		Aliases: []string{"child-delegates"},
		Short:   "List the users a user delegated access to",
		Long: `List the delegate users that can act on the resources of a specific user.

Required Arguments:
  user_id                 The numeric ID of the user whose delegates to list

Examples:
  metalcloud-cli user delegate children 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.GetChildDelegates(cmd.Context(), args[0])
		},
	}

	userResendEmailVerificationCmd = &cobra.Command{
		Use:     "resend-email-verification user_id",
		Aliases: []string{"resend-verification"},
		Short:   "Resend the e-mail verification message to a user",
		Long: `Resend the e-mail address verification message to a specific user.

The user receives a new verification link. Use --redirect-url to control where the
user lands after following the link.

Required Arguments:
  user_id                 The numeric ID of the user to notify

Optional Flags:
  --redirect-url          URL the user is redirected to after verifying the address

Examples:
  metalcloud-cli user resend-email-verification 12345
  metalcloud-cli user resend-email-verification 12345 --redirect-url https://metalsoft.io/welcome`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.ResendEmailVerification(cmd.Context(), args[0], userFlags.redirectUrl)
		},
	}

	userResendInvitationCmd = &cobra.Command{
		Use:     "resend-invitation user_id",
		Aliases: []string{"resend-user-invitation"},
		Short:   "Resend the platform invitation to a user",
		Long: `Resend the platform invitation message to a specific user.

The user receives a new invitation link. Use --redirect-url to control where the user
lands after accepting the invitation.

Required Arguments:
  user_id                 The numeric ID of the user to invite again

Optional Flags:
  --redirect-url          URL the user is redirected to after accepting the invitation

Examples:
  metalcloud-cli user resend-invitation 12345
  metalcloud-cli user resend-invitation 12345 --redirect-url https://metalsoft.io/welcome`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.ResendInvitation(cmd.Context(), args[0], userFlags.redirectUrl)
		},
	}

	userSendPasswordResetCmd = &cobra.Command{
		Use:     "send-password-reset user_id",
		Aliases: []string{"password-reset-send"},
		Short:   "Send a password reset message to a user (admin)",
		Long: `Send a password reset message to a specific user as an administrator.

The user receives a link that lets them choose a new password. Use 'user set-password'
instead to set a password directly without involving the user.

Required Arguments:
  user_id                 The numeric ID of the user to notify

Optional Flags:
  --redirect-url          URL the user is redirected to after resetting the password

Examples:
  metalcloud-cli user send-password-reset 12345
  metalcloud-cli user send-password-reset 12345 --redirect-url https://metalsoft.io/login`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.SendPasswordReset(cmd.Context(), args[0], userFlags.redirectUrl)
		},
	}

	userDeleteCmd = &cobra.Command{
		Use:     "delete user_id",
		Aliases: []string{"rm"},
		Short:   "Delete a user and erase their personal information",
		Long: `Delete a user account and irreversibly erase their personal information.

WARNING: this is NOT the same as 'user archive'. Archiving only marks the account as
inactive and can be undone with 'user unarchive'. Deleting archives the user AND
permanently removes their personally identifiable information; it cannot be undone
and 'user unarchive' will not bring the information back.

Required Arguments:
  user_id                 The numeric ID of the user to delete

Examples:
  metalcloud-cli user delete 12345
  metalcloud-cli user rm 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.Delete(cmd.Context(), args[0])
		},
	}

	userChangePasswordCmd = &cobra.Command{
		Use:     "change-password",
		Aliases: []string{"password-change"},
		Short:   "Change the password of the current user",
		Long: `Change the password of the user owning the API key in use.

The new password can be supplied through flags or through a JSON/YAML payload. The
current password is required unless the platform policy allows changing it without.

Required Flags (when not using --config-source):
  --new-password          The new password

Optional Flags:
  --current-password      The password currently in use
  --config-source         Source of the password change payload (JSON/YAML file or 'pipe')

Configuration File Format (JSON):
  {
    "newPassword": "newSecret123",
    "oldPassword": "oldSecret123"
  }

Examples:
  metalcloud-cli user change-password --current-password oldSecret123 --new-password newSecret123
  echo '{"newPassword":"newSecret123","oldPassword":"oldSecret123"}' | metalcloud-cli user change-password --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if userFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(userFlags.configSource)
				if err != nil {
					return err
				}

				return user.ChangePassword(cmd.Context(), config)
			}

			passwordUpdate := sdk.UserUpdatePassword{
				NewPassword: userFlags.newPassword,
			}
			if userFlags.currentPassword != "" {
				passwordUpdate.OldPassword = sdk.PtrString(userFlags.currentPassword)
			}

			configBytes, err := json.Marshal(passwordUpdate)
			if err != nil {
				return fmt.Errorf("could not marshal password change payload: %s", err)
			}

			return user.ChangePassword(cmd.Context(), configBytes)
		},
	}

	userInitiatePasswordResetCmd = &cobra.Command{
		Use:     "initiate-password-reset",
		Aliases: []string{"password-reset"},
		Short:   "Send a password reset message to an e-mail address",
		Long: `Start the self-service password reset flow for an e-mail address.

The owner of the address receives a link that lets them choose a new password.

Required Flags:
  --email                 The e-mail address of the account to reset

Optional Flags:
  --redirect-url          URL the user is redirected to after resetting the password

Examples:
  metalcloud-cli user initiate-password-reset --email user@company.com
  metalcloud-cli user initiate-password-reset --email user@company.com --redirect-url https://metalsoft.io/login`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.InitiatePasswordReset(cmd.Context(), userFlags.email, userFlags.redirectUrl)
		},
	}

	userInitiateEmailChangeCmd = &cobra.Command{
		Use:     "initiate-email-change",
		Aliases: []string{"email-change"},
		Short:   "Start changing the e-mail address of the current user",
		Long: `Start the e-mail address change flow for the user owning the API key in use.

A verification message is sent to the new address; the change takes effect only after
the link in that message is followed.

Required Flags:
  --email                 The new e-mail address

Optional Flags:
  --redirect-url          URL the user is redirected to after verifying the address

Examples:
  metalcloud-cli user initiate-email-change --email new.address@company.com`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.InitiateEmailChange(cmd.Context(), userFlags.email, userFlags.redirectUrl)
		},
	}

	userRegenerateJwtSaltCmd = &cobra.Command{
		Use:     "regenerate-jwt-salt",
		Aliases: []string{"jwt-salt-regenerate"},
		Short:   "Regenerate the JWT salt of the current user",
		Long: `Regenerate the JWT salt of the user owning the API key in use.

WARNING: this invalidates every session and token issued so far for this user, on
every device and in every browser. You will have to log in again.

Examples:
  metalcloud-cli user regenerate-jwt-salt`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.RegenerateJwtSalt(cmd.Context())
		},
	}

	userPermissionsCmd = &cobra.Command{
		Use:     "permissions",
		Aliases: []string{"my-permissions"},
		Short:   "List the permissions of the current user",
		Long: `List the permissions of the user owning the API key in use.

The permissions come from the roles assigned to the user and determine which API
operations, and therefore which CLI commands, are available.

Examples:
  metalcloud-cli user permissions
  metalcloud-cli user permissions -f json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.GetPermissions(cmd.Context())
		},
	}

	userVerifyEmailCmd = &cobra.Command{
		Use:     "verify-email",
		Aliases: []string{"email-verify"},
		Short:   "Verify an e-mail address with a verification token",
		Long: `Consume an e-mail verification token, as following the link from the e-mail would.

The token is the value of the 'token' query parameter of the verification link that
the platform sent by e-mail.

Required Flags:
  --token                 The e-mail verification token

Examples:
  metalcloud-cli user verify-email --token eyJhbGciOi...`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.VerifyEmail(cmd.Context(), userFlags.callbackToken)
		},
	}

	userResetPasswordCmd = &cobra.Command{
		Use:     "reset-password",
		Aliases: []string{"handle-password-reset"},
		Short:   "Consume a password reset token",
		Long: `Consume a password reset token, as following the link from the e-mail would.

The token is the value of the 'token' query parameter of the password reset link that
the platform sent by e-mail. The API exposes no password parameter on this endpoint:
it only validates the token and redirects to the page where the new password is
chosen. Use 'user change-password' to set a password from the CLI, or
'user set-password' to set the password of another user as an administrator.

Required Flags:
  --token                 The password reset token

Examples:
  metalcloud-cli user reset-password --token eyJhbGciOi...`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.ResetPassword(cmd.Context(), userFlags.callbackToken)
		},
	}
)

func init() {
	rootCmd.AddCommand(userCmd)

	userCmd.AddCommand(userListCmd)
	userListCmd.Flags().BoolVar(&userFlags.archived, "archived", false, "Return only archived users.")
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

	// User limits (effective quota limits, read-only)
	userCmd.AddCommand(userLimitsGetCmd)

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

	// Set password
	userCmd.AddCommand(userSetPasswordCmd)
	userSetPasswordCmd.Flags().StringVar(&userFlags.password, "password", "", "The new password to set for the user.")
	userSetPasswordCmd.MarkFlagRequired("password")

	// User configuration and metadata
	userCmd.AddCommand(userConfigGetCmd)

	userCmd.AddCommand(userUpdateMetaCmd)
	userUpdateMetaCmd.Flags().StringVar(&userFlags.configSource, "config-source", "", "Source of the user metadata. Can be 'pipe' or path to a JSON/YAML file.")
	userUpdateMetaCmd.MarkFlagsOneRequired("config-source")

	// Single SSH key
	userCmd.AddCommand(userSshKeyGetCmd)

	// Suspend reasons
	userCmd.AddCommand(userSuspendReasonsCmd)

	// Delegates
	userCmd.AddCommand(userDelegateCmd)
	userDelegateCmd.AddCommand(userDelegateAddCmd)
	userDelegateCmd.AddCommand(userDelegateRemoveCmd)
	userDelegateCmd.AddCommand(userDelegateParentsCmd)
	userDelegateCmd.AddCommand(userDelegateChildrenCmd)

	// Notifications
	userCmd.AddCommand(userResendEmailVerificationCmd)
	userResendEmailVerificationCmd.Flags().StringVar(&userFlags.redirectUrl, "redirect-url", "", "URL the user is redirected to after verifying the e-mail address.")

	userCmd.AddCommand(userResendInvitationCmd)
	userResendInvitationCmd.Flags().StringVar(&userFlags.redirectUrl, "redirect-url", "", "URL the user is redirected to after accepting the invitation.")

	userCmd.AddCommand(userSendPasswordResetCmd)
	userSendPasswordResetCmd.Flags().StringVar(&userFlags.redirectUrl, "redirect-url", "", "URL the user is redirected to after resetting the password.")

	// Delete (archives the user and erases their personal information)
	userCmd.AddCommand(userDeleteCmd)

	// Self-service: password
	userCmd.AddCommand(userChangePasswordCmd)
	userChangePasswordCmd.Flags().StringVar(&userFlags.newPassword, "new-password", "", "The new password of the current user.")
	userChangePasswordCmd.Flags().StringVar(&userFlags.currentPassword, "current-password", "", "The password currently in use.")
	userChangePasswordCmd.Flags().StringVar(&userFlags.configSource, "config-source", "", "Source of the password change payload. Can be 'pipe' or path to a JSON/YAML file.")
	userChangePasswordCmd.MarkFlagsOneRequired("new-password", "config-source")
	userChangePasswordCmd.MarkFlagsMutuallyExclusive("new-password", "config-source")
	userChangePasswordCmd.MarkFlagsMutuallyExclusive("current-password", "config-source")

	userCmd.AddCommand(userInitiatePasswordResetCmd)
	userInitiatePasswordResetCmd.Flags().StringVar(&userFlags.email, "email", "", "The e-mail address of the account to reset.")
	userInitiatePasswordResetCmd.Flags().StringVar(&userFlags.redirectUrl, "redirect-url", "", "URL the user is redirected to after resetting the password.")
	userInitiatePasswordResetCmd.MarkFlagRequired("email")

	// Self-service: e-mail address
	userCmd.AddCommand(userInitiateEmailChangeCmd)
	userInitiateEmailChangeCmd.Flags().StringVar(&userFlags.email, "email", "", "The new e-mail address of the current user.")
	userInitiateEmailChangeCmd.Flags().StringVar(&userFlags.redirectUrl, "redirect-url", "", "URL the user is redirected to after verifying the e-mail address.")
	userInitiateEmailChangeCmd.MarkFlagRequired("email")

	// Self-service: sessions and permissions
	userCmd.AddCommand(userRegenerateJwtSaltCmd)
	userCmd.AddCommand(userPermissionsCmd)

	// Self-service: e-mail verification and password reset callbacks
	userCmd.AddCommand(userVerifyEmailCmd)
	userVerifyEmailCmd.Flags().StringVar(&userFlags.callbackToken, "token", "", "The e-mail verification token.")
	userVerifyEmailCmd.MarkFlagRequired("token")

	userCmd.AddCommand(userResetPasswordCmd)
	userResetPasswordCmd.Flags().StringVar(&userFlags.callbackToken, "token", "", "The password reset token.")
	userResetPasswordCmd.MarkFlagRequired("token")

	// Self-service: API key
	userCmd.AddCommand(userApiKeyGetCmd)
	userCmd.AddCommand(userApiKeyRegenerateCmd)

	// Self-service: 2FA
	userCmd.AddCommand(user2FAEnableCmd)
	user2FAEnableCmd.Flags().StringVar(&userFlags.twoFAToken, "token", "", "The TOTP code from your authenticator app.")
	user2FAEnableCmd.MarkFlagRequired("token")
	userCmd.AddCommand(user2FADisableCmd)
	userCmd.AddCommand(user2FAGenerateSecretCmd)
}
