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
  --archived              Include archived users in the results (default: false)
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

	userLimitsUpdateCmd = &cobra.Command{
		Use:     "limits-update user_id",
		Aliases: []string{"update-limits"},
		Short:   "Update resource limits for a specific user",
		Long: `Update the resource limits for a specific user account to control their resource allocation.

This command allows you to modify compute node, drive, and infrastructure limits that restrict
how many resources the user can provision. Changes take effect immediately.

Arguments:
  user_id                 The numeric ID of the user whose limits to update

Required Flags:
  --config-source         Source of user limits configuration (JSON file path or 'pipe')

Configuration File Format (JSON):
  {
    "computeNodesInstancesToProvisionLimit": 100,
    "drivesAttachedToInstancesLimit": 200,
    "infrastructuresLimit": 10
  }

Examples:
  metalcloud-cli user limits-update 12345 --config-source limits.json
  echo '{"computeNodesInstancesToProvisionLimit": 50}' | metalcloud-cli user limits-update 12345 --config-source pipe`,
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

	userPermissionsGetCmd = &cobra.Command{
		Use:     "permissions user_id",
		Aliases: []string{"get-permissions"},
		Short:   "Display permissions for a specific user",
		Long: `Retrieve and display all permissions associated with a specific user account.

This command shows the user's permission configuration including resource types,
resource IDs, and permission levels. Permissions control what resources the user
can access and what operations they can perform.

Arguments:
  user_id                 The numeric ID of the user whose permissions to display

Examples:
  metalcloud-cli user permissions 12345
  metalcloud-cli user get-permissions 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.GetPermissions(cmd.Context(), args[0])
		},
	}

	userPermissionsUpdateCmd = &cobra.Command{
		Use:     "permissions-update user_id",
		Aliases: []string{"update-permissions"},
		Short:   "Update permissions for a specific user",
		Long: `Update the permissions configuration for a specific user account.

This command allows modifying user permissions including resource access levels and
operational capabilities. Permissions changes take effect immediately and control
what resources the user can access and what operations they can perform.

Arguments:
  user_id                 The numeric ID of the user whose permissions to update

Required Flags:
  --config-source         Source of user permissions configuration (JSON file path or 'pipe')

Configuration File Format (JSON):
  {
    "permissions": [
      {
        "resourceType": "infrastructure",
        "resourceId": "123",
        "permissionLevel": "read"
      },
      {
        "resourceType": "account",
        "resourceId": "456",
        "permissionLevel": "write"
      }
    ]
  }

Examples:
  metalcloud-cli user permissions-update 12345 --config-source permissions.json
  echo '{"permissions": [{"resourceType": "infrastructure", "resourceId": "123", "permissionLevel": "read"}]}' | metalcloud-cli user permissions-update 12345 --config-source pipe`,
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
