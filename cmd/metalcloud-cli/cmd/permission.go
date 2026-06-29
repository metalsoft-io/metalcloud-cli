package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/permission"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

// Permission commands
var (
	permissionFlags = struct {
		configSource string
	}{}

	permissionCmd = &cobra.Command{
		Use:     "permission [command]",
		Aliases: []string{"permissions", "perm"},
		Short:   "Permission management",
		Long: `Permission management commands.

This command group provides management capabilities for permissions including
listing, creating, and deleting permissions.

Available commands:
  - list, create, delete, config-example

Use "metalcloud-cli permission [command] --help" for detailed information about each command.
`,
	}

	permissionListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List permissions",
		Long: `List all permissions.

This command displays information about all permissions including their IDs,
names, labels, types, and descriptions.

Examples:
  # List all permissions
  metalcloud-cli permission list
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_PERMISSIONS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return permission.PermissionList(cmd.Context())
		},
	}

	permissionCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new permission",
		Long: `Create a new permission.

You must provide the permission configuration using the --config-source flag.
The configuration source can be a path to a JSON file or 'pipe' to read from
standard input.

Required Flags:
  --config-source       Source of the permission configuration. Can be 'pipe' or path to a JSON file

Examples:
  # Create using JSON configuration file
  metalcloud-cli permission create --config-source ./permission.json

  # Create using piped JSON configuration
  cat permission.json | metalcloud-cli permission create --config-source pipe
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_PERMISSIONS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(permissionFlags.configSource)
			if err != nil {
				return err
			}
			return permission.PermissionCreate(cmd.Context(), config)
		},
	}

	permissionDeleteCmd = &cobra.Command{
		Use:     "delete permission_name",
		Aliases: []string{"rm"},
		Short:   "Delete a permission",
		Long: `Delete a permission.

This command permanently deletes a permission. This action cannot be undone,
so use with caution.

Required Arguments:
  permission_name       The name of the permission to delete

Examples:
  # Delete a permission
  metalcloud-cli permission delete custom_read

  # Delete a permission using alias
  metalcloud-cli permission rm custom_read
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_PERMISSIONS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return permission.PermissionDelete(cmd.Context(), args[0])
		},
	}

	permissionConfigExampleCmd = &cobra.Command{
		Use:     "config-example",
		Aliases: []string{"example"},
		Short:   "Show a permission configuration example",
		Long: `Show an example permission configuration that can be used as a
template for creating a new permission.

Examples:
  # Show a permission configuration example
  metalcloud-cli permission config-example
`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_PERMISSIONS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return permission.PermissionConfigExample(cmd.Context())
		},
	}
)

func init() {
	rootCmd.AddCommand(permissionCmd)

	permissionCmd.AddCommand(permissionListCmd)

	permissionCmd.AddCommand(permissionCreateCmd)
	permissionCreateCmd.Flags().StringVar(&permissionFlags.configSource, "config-source", "", "Source of the new permission configuration. Can be 'pipe' or path to a JSON file.")
	permissionCreateCmd.MarkFlagsOneRequired("config-source")

	permissionCmd.AddCommand(permissionDeleteCmd)

	permissionCmd.AddCommand(permissionConfigExampleCmd)
}
