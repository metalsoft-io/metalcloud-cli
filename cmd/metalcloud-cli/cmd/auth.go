package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/auth"
	"github.com/spf13/cobra"
)

var (
	authFlags = struct {
		roleName               string
		priority               int32
		userExternalIdentifier string
		username               string
		email                  string
	}{}

	authCmd = &cobra.Command{
		Use:     "auth [command]",
		Aliases: []string{"authentication"},
		Short:   "Manage authentication settings",
		Long: `Manage authentication settings for the MetalCloud platform.

This command provides subcommands to configure and manage various authentication
methods including LDAP integration and user role mappings.

Available authentication methods:
- LDAP: Configure LDAP server integration and group-to-role mappings

Examples:
  # List all available auth subcommands
  metalcloud-cli auth --help
  
  # Work with LDAP authentication
  metalcloud-cli auth ldap --help`,
	}

	authLdapCmd = &cobra.Command{
		Use:     "ldap [command]",
		Aliases: []string{"LDAP"},
		Short:   "Manage LDAP authentication settings",
		Long: `Manage LDAP authentication settings for the MetalCloud platform.

Configure LDAP server integration and manage group-to-role mappings that determine
user permissions based on their LDAP group memberships.

Available operations:
- List existing LDAP group mappings
- Add new LDAP group mappings
- Update existing LDAP group mappings  
- Remove LDAP group mappings

Examples:
  # List all LDAP group mappings
  metalcloud-cli auth ldap mapping-list
  
  # Add a new LDAP group mapping
  metalcloud-cli auth ldap mapping-add "Domain Admins" --role-name admin --priority 1
  
  # Update an existing mapping
  metalcloud-cli auth ldap mapping-update "Power Users" --role-name power-user --priority 5
  
  # Remove a mapping
  metalcloud-cli auth ldap mapping-remove "Guests"`,
	}

	authLdapMappingListCmd = &cobra.Command{
		Use:     "mapping-list",
		Aliases: []string{"mapping-ls", "map-list", "map-ls"},
		Short:   "List all LDAP group mappings",
		Long: `List all configured LDAP group mappings and their associated roles.

This command displays all LDAP group-to-role mappings that are currently configured
in the system. Each mapping shows the LDAP group name, the MetalCloud role it maps to,
and the priority of the mapping.

The priority determines which role is assigned when a user belongs to multiple LDAP
groups with different mappings. Lower priority numbers take precedence.

Examples:
  # List all LDAP group mappings
  metalcloud-cli auth ldap mapping-list
  
  # List mappings with short alias
  metalcloud-cli auth ldap map-ls`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_GLOBAL_CONFIGURATIONS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return auth.AuthLdapMappingList(cmd.Context())
		},
	}

	authLdapMappingAddCmd = &cobra.Command{
		Use:     "mapping-add group_name",
		Aliases: []string{"mapping-create", "map-add", "map-create"},
		Short:   "Add LDAP group mapping",
		Long: `Add a new LDAP group-to-role mapping for authentication.

This command creates a new mapping between an LDAP group and a MetalCloud role.
Users who belong to the specified LDAP group will be assigned the specified role
when they authenticate via LDAP.

The priority value determines which role is assigned when a user belongs to multiple
LDAP groups with different mappings. Lower priority numbers take precedence over
higher ones. If two mappings have the same priority, the behavior is undefined.

Examples:
  # Map "Domain Admins" group to admin role with highest priority
  metalcloud-cli auth ldap mapping-add "Domain Admins" --role-name admin --priority 1
  
  # Map "Power Users" group to power-user role  
  metalcloud-cli auth ldap mapping-add "Power Users" --role-name power-user --priority 5
  
  # Map "Developers" group to developer role
  metalcloud-cli auth ldap mapping-add "Developers" --role-name developer --priority 10`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_GLOBAL_CONFIGURATIONS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return auth.AuthLdapMappingAdd(cmd.Context(), args[0], auth.AuthLdapMappingOptions{
				RoleName:               authFlags.roleName,
				Priority:               authFlags.priority,
				UserExternalIdentifier: optionalStringFlag(cmd, "user-external-identifier", authFlags.userExternalIdentifier),
				Username:               optionalStringFlag(cmd, "username", authFlags.username),
				Email:                  optionalStringFlag(cmd, "email", authFlags.email),
			})
		},
	}

	authLdapMappingUpdateCmd = &cobra.Command{
		Use:     "mapping-update group_name",
		Aliases: []string{"map-update"},
		Short:   "Update LDAP group mapping",
		Long: `Update an existing LDAP group-to-role mapping.

This command modifies an existing mapping between an LDAP group and a MetalCloud role.
You can update the role name, the priority, or any of the optional attributes. At
least one of these flags must be provided.

Passing one of the optional string flags (--user-external-identifier, --username,
--email) with an empty value resets the attribute to its default (objectGUID,
sAMAccountName, mail respectively).

The priority value determines which role is assigned when a user belongs to multiple
LDAP groups with different mappings. Lower priority numbers take precedence over
higher ones.

Examples:
  # Update the role name for "Power Users" group
  metalcloud-cli auth ldap mapping-update "Power Users" --role-name senior-developer

  # Update the priority for "Developers" group
  metalcloud-cli auth ldap mapping-update "Developers" --priority 15

  # Update both role name and priority
  metalcloud-cli auth ldap mapping-update "Guests" --role-name read-only --priority 20

  # Reset the email attribute to its default value
  metalcloud-cli auth ldap mapping-update "Developers" --email ""`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_GLOBAL_CONFIGURATIONS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return auth.AuthLdapMappingUpdate(cmd.Context(), args[0], auth.AuthLdapMappingOptions{
				RoleName:               authFlags.roleName,
				Priority:               authFlags.priority,
				UserExternalIdentifier: optionalStringFlag(cmd, "user-external-identifier", authFlags.userExternalIdentifier),
				Username:               optionalStringFlag(cmd, "username", authFlags.username),
				Email:                  optionalStringFlag(cmd, "email", authFlags.email),
			})
		},
	}

	authLdapMappingRemoveCmd = &cobra.Command{
		Use:     "mapping-remove group_name",
		Aliases: []string{"mapping-rm", "map-remove", "map-rm"},
		Short:   "Delete LDAP group mapping",
		Long: `Delete an existing LDAP group-to-role mapping.

This command removes a mapping between an LDAP group and a MetalCloud role.
Once removed, users who belong only to the specified LDAP group will no longer
receive the associated role when authenticating via LDAP.

This operation is irreversible. If you need the mapping again, you will need
to recreate it using the mapping-add command.

Examples:
  # Remove mapping for "Guests" group
  metalcloud-cli auth ldap mapping-remove "Guests"
  
  # Remove mapping for "Contractors" group using short alias
  metalcloud-cli auth ldap map-rm "Contractors"
  
  # Remove mapping for group with spaces in name
  metalcloud-cli auth ldap mapping-remove "External Users"`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_GLOBAL_CONFIGURATIONS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return auth.AuthLdapMappingRemove(cmd.Context(), args[0])
		},
	}
)

func optionalStringFlag(cmd *cobra.Command, name string, value string) *string {
	if !cmd.Flags().Changed(name) {
		return nil
	}
	return &value
}

func init() {
	rootCmd.AddCommand(authCmd)

	authCmd.AddCommand(authLdapCmd)

	authLdapCmd.AddCommand(authLdapMappingListCmd)

	authLdapCmd.AddCommand(authLdapMappingAddCmd)
	authLdapMappingAddCmd.Flags().StringVar(&authFlags.roleName, "role-name", "", "Role name to map to the LDAP group.")
	authLdapMappingAddCmd.Flags().Int32Var(&authFlags.priority, "priority", 10, "Mapping priority.")
	authLdapMappingAddCmd.Flags().StringVar(&authFlags.userExternalIdentifier, "user-external-identifier", "", "LDAP attribute used as the user external identifier.")
	authLdapMappingAddCmd.Flags().StringVar(&authFlags.username, "username", "", "LDAP attribute used as the username.")
	authLdapMappingAddCmd.Flags().StringVar(&authFlags.email, "email", "", "LDAP attribute used as the email.")
	authLdapMappingAddCmd.MarkFlagRequired("role-name")
	authLdapMappingAddCmd.MarkFlagRequired("priority")

	authLdapCmd.AddCommand(authLdapMappingUpdateCmd)
	authLdapMappingUpdateCmd.Flags().StringVar(&authFlags.roleName, "role-name", "", "Role name to map to the LDAP group.")
	authLdapMappingUpdateCmd.Flags().Int32Var(&authFlags.priority, "priority", 10, "Mapping priority.")
	authLdapMappingUpdateCmd.Flags().StringVar(&authFlags.userExternalIdentifier, "user-external-identifier", "", "LDAP attribute used as the user external identifier.")
	authLdapMappingUpdateCmd.Flags().StringVar(&authFlags.username, "username", "", "LDAP attribute used as the username.")
	authLdapMappingUpdateCmd.Flags().StringVar(&authFlags.email, "email", "", "LDAP attribute used as the email.")
	authLdapMappingUpdateCmd.MarkFlagsOneRequired("role-name", "priority", "user-external-identifier", "username", "email")

	authLdapCmd.AddCommand(authLdapMappingRemoveCmd)
}
