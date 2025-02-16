package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/user"
	"github.com/spf13/cobra"
)

var (
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
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.USERS_AND_PERMISSIONS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.List(cmd.Context())
		},
	}

	userGetCmd = &cobra.Command{
		Use:          "get",
		Aliases:      []string{"show"},
		Short:        "Get user details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.USERS_AND_PERMISSIONS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return user.Get(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(userCmd)

	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userGetCmd)
}
