package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/role"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	roleFlags = struct {
		configSource string
	}{}

	roleCmd = &cobra.Command{
		Use:     "role [command]",
		Aliases: []string{"roles"},
		Short:   "Manage user roles and permissions",
		Long: `Manage user roles and permissions in the MetalCloud platform.

Roles define sets of permissions that can be assigned to users to control access
to various platform features and resources. This command group provides operations
to list, view, create, update, and delete roles.

Available operations:
  list    - List all available roles
  get     - View detailed information about a specific role
  create  - Create a new role with specified permissions
  update  - Update an existing role's permissions or metadata
  delete  - Remove a role from the system

Examples:
  # List all roles
  metalcloud-cli role list

  # View details of a specific role
  metalcloud-cli role get admin-role

  # Create a new role from a JSON file
  metalcloud-cli role create --config-source role-config.json

  # Update a role using piped JSON
  echo '{"permissions": ["read", "write"]}' | metalcloud-cli role update my-role --config-source pipe`,
	}

	roleListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all available roles",
		Long: `List all available roles in the MetalCloud platform.

This command displays a table of all roles with their ID, label, name, description,
and the number of users assigned to each role. The output includes both system
roles and custom roles created by administrators.

Examples:
  # List all roles
  metalcloud-cli role list

  # List roles using alias
  metalcloud-cli roles ls`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ROLES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return role.List(cmd.Context())
		},
	}

	roleGetCmd = &cobra.Command{
		Use:     "get <role_name>",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific role",
		Long: `Get detailed information about a specific role in the MetalCloud platform.

This command retrieves and displays comprehensive details about a role including
its ID, name, label, description, type, assigned permissions, and the number of
users that have this role assigned.

Arguments:
  role_name    Name of the role to retrieve information for

Examples:
  # Get details of a specific role
  metalcloud-cli role get admin-role

  # Using the alias
  metalcloud-cli role show editor-role

  # Get details of a system role
  metalcloud-cli roles get super-admin`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ROLES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return role.Get(cmd.Context(), args[0])
		},
	}

	roleCreateCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   "Create a new role with specified permissions",
		Long: `Create a new role with specified permissions in the MetalCloud platform.

This command creates a new role based on configuration provided through a JSON file
or piped input. The role configuration must include a label and permissions list.
Description is optional but recommended for clarity.

Required configuration fields:
  label         Human-readable name for the role
  permissions   Array of permission strings

Optional configuration fields:
  description   Detailed description of the role's purpose

Flags:
  --config-source string   Required. Source of the role configuration.
                          Can be 'pipe' for stdin input or path to a JSON file.

Configuration format (JSON):
{
  "label": "Custom Admin Role",
  "description": "Administrative role with full system access",
  "permissions": ["roles:read", "roles:write", "users:read", "users:write"]
}

Examples:
  # Create role from JSON file
  metalcloud-cli role create --config-source role-config.json

  # Create role using piped JSON
  echo '{"label": "Editor", "permissions": ["content:read", "content:write"]}' | metalcloud-cli role create --config-source pipe

  # Create role with multiple permissions from file
  metalcloud-cli role new --config-source /path/to/role.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ROLES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(roleFlags.configSource)
			if err != nil {
				return err
			}
			return role.Create(cmd.Context(), config)
		},
	}

	roleDeleteCmd = &cobra.Command{
		Use:     "delete <role_name>",
		Aliases: []string{"remove", "rm"},
		Short:   "Delete a role from the system",
		Long: `Delete a role from the MetalCloud platform.

This command permanently removes a role from the system. Once deleted, the role
cannot be recovered and any users assigned to this role will lose those permissions.
System roles cannot be deleted.

Arguments:
  role_name    Name of the role to delete

Warning:
  - This operation is irreversible
  - Users assigned to this role will lose the associated permissions
  - System roles cannot be deleted

Examples:
  # Delete a custom role
  metalcloud-cli role delete custom-editor

  # Using the alias
  metalcloud-cli role rm temp-role

  # Using the remove alias
  metalcloud-cli roles remove old-role`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ROLES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return role.Delete(cmd.Context(), args[0])
		},
	}

	roleUpdateCmd = &cobra.Command{
		Use:     "update <role_name>",
		Aliases: []string{"edit"},
		Short:   "Update an existing role's permissions or metadata",
		Long: `Update an existing role's permissions or metadata in the MetalCloud platform.

This command modifies an existing role based on configuration provided through a JSON file
or piped input. You can update the role's label, description, and permissions.
System roles cannot be updated.

Arguments:
  role_name    Name of the role to update

Configuration fields (all optional):
  label         New human-readable name for the role
  description   New description of the role's purpose
  permissions   New array of permission strings (replaces existing permissions)

Flags:
  --config-source string   Required. Source of the role configuration.
                          Can be 'pipe' for stdin input or path to a JSON file.

Configuration format (JSON):
{
  "label": "Updated Admin Role",
  "description": "Updated administrative role with modified access",
  "permissions": ["roles:read", "roles:write", "users:read"]
}

Note: When updating permissions, the provided array completely replaces the existing
permissions. To add permissions, include both existing and new permissions in the array.

Examples:
  # Update role from JSON file
  metalcloud-cli role update custom-admin --config-source role-update.json

  # Update role using piped JSON
  echo '{"description": "Updated role description"}' | metalcloud-cli role update my-role --config-source pipe

  # Update role permissions and label
  metalcloud-cli roles edit editor-role --config-source /path/to/updated-role.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_ROLES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(roleFlags.configSource)
			if err != nil {
				return err
			}
			return role.Update(cmd.Context(), args[0], config)
		},
	}
)

func init() {
	rootCmd.AddCommand(roleCmd)

	// Role list
	roleCmd.AddCommand(roleListCmd)

	// Role get
	roleCmd.AddCommand(roleGetCmd)

	// Role create
	roleCmd.AddCommand(roleCreateCmd)
	roleCreateCmd.Flags().StringVar(&roleFlags.configSource, "config-source", "", "Source of the role configuration. Can be 'pipe' or path to a JSON file.")
	roleCreateCmd.MarkFlagRequired("config-source")

	// Role delete
	roleCmd.AddCommand(roleDeleteCmd)

	// Role update
	roleCmd.AddCommand(roleUpdateCmd)
	roleUpdateCmd.Flags().StringVar(&roleFlags.configSource, "config-source", "", "Source of the role configuration. Can be 'pipe' or path to a JSON file.")
	roleUpdateCmd.MarkFlagRequired("config-source")
}
