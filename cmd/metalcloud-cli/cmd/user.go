package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/user"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	userFlags = struct {
		configSource  string
		accountId     float32
		sshKeyContent string
		reason        string
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
			return user.List(cmd.Context())
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

	userCreateCmd = &cobra.Command{
		Use:          "create",
		Aliases:      []string{"new"},
		Short:        "Create a new user.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_USERS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(userFlags.configSource)
			if err != nil {
				return err
			}

			return user.Create(cmd.Context(), config)
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
	userCmd.AddCommand(userGetCmd)
	userCmd.AddCommand(userLimitsGetCmd)

	// User create
	userCmd.AddCommand(userCreateCmd)
	userCreateCmd.Flags().StringVar(&userFlags.configSource, "config-source", "", "Source of the new user configuration. Can be 'pipe' or path to a JSON file.")
	userCreateCmd.MarkFlagsOneRequired("config-source")

	// User archive/unarchive
	userCmd.AddCommand(userArchiveCmd)
	userCmd.AddCommand(userUnarchiveCmd)

	// User limits update
	userCmd.AddCommand(userLimitsUpdateCmd)
	userLimitsUpdateCmd.Flags().StringVar(&userFlags.configSource, "config-source", "", "Source of the user limits configuration. Can be 'pipe' or path to a JSON file.")
	userLimitsUpdateCmd.MarkFlagsOneRequired("config-source")

	// User config update
	userCmd.AddCommand(userConfigUpdateCmd)
	userConfigUpdateCmd.Flags().StringVar(&userFlags.configSource, "config-source", "", "Source of the user configuration. Can be 'pipe' or path to a JSON file.")
	userConfigUpdateCmd.MarkFlagsOneRequired("config-source")

	// Change account
	userCmd.AddCommand(userChangeAccountCmd)
	userChangeAccountCmd.Flags().Float32Var(&userFlags.accountId, "account-id", 0, "The ID of the account to move the user to.")
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
