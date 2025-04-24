package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/role"
	"github.com/spf13/cobra"
)

var (
	permissionCmd = &cobra.Command{
		Use:     "permission [command]",
		Aliases: []string{"permissions"},
		Short:   "Permission management",
		Long:    `Permission management commands.`,
	}

	permissionListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all permissions.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ROLES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return role.ListPermissions(cmd.Context())
		},
	}
)

func init() {
	rootCmd.AddCommand(permissionCmd)
	permissionCmd.AddCommand(permissionListCmd)
}
