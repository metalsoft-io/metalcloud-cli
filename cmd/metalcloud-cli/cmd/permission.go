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
		Short:   "Manage system permissions and access control",
		Long: `Manage system permissions and access control.

Permissions define what actions users and roles can perform within the MetalCloud platform.
This command group provides functionality to view and manage the available permissions
in the system.

Available Commands:
  list        List all available permissions in the system

Examples:
  # List all permissions
  metalcloud-cli permission list
  
  # List permissions with short alias
  metalcloud-cli permissions ls

Required Permissions:
  Most permission management operations require the 'roles:read' permission.`,
	}

	permissionListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all available system permissions",
		Long: `List all available system permissions.

This command displays all permissions available in the MetalCloud system, including
their IDs, names, labels, and descriptions. Permissions control access to various
system resources and operations.

The output includes:
  - ID: Unique identifier for the permission
  - Label: Human-readable label for the permission
  - Name: Technical name of the permission
  - Description: Detailed description of what the permission allows

Required Permissions:
  - roles:read: Required to view permission information

Examples:
  # List all permissions
  metalcloud-cli permission list
  
  # List permissions using alias
  metalcloud-cli permissions ls
  
  # List permissions with JSON output format
  metalcloud-cli permission list --format json`,
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
