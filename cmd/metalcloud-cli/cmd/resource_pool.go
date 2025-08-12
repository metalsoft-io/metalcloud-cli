package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/resource_pool"
	"github.com/spf13/cobra"
)

var (
	resourcePoolFlags = struct {
		pageFlag        int
		limitFlag       int
		searchFlag      string
		labelFlag       string
		descriptionFlag string
	}{}

	resourcePoolCmd = &cobra.Command{
		Use:     "resource-pool [command]",
		Aliases: []string{"rpool", "rp"},
		Short:   "Manage resource pools and their associated resources",
		Long: `Manage resource pools and their associated resources including users, servers, and subnet pools.

Resource pools are logical groupings of infrastructure resources that can be assigned to users
for organizing and controlling access to compute, network, and storage resources.

Available commands:
  list              List all resource pools
  get               Get detailed information about a specific resource pool
  create            Create a new resource pool
  delete            Delete a resource pool
  get-users         List users with access to a resource pool
  add-user          Grant a user access to a resource pool
  remove-user       Revoke user access from a resource pool
  get-servers       List servers assigned to a resource pool
  add-server        Add a server to a resource pool
  remove-server     Remove a server from a resource pool
  get-subnet-pools  List subnet pools assigned to a resource pool
  add-subnet-pool   Add a subnet pool to a resource pool
  remove-subnet-pool Remove a subnet pool from a resource pool

Examples:
  # List all resource pools
  metalcloud-cli resource-pool list

  # Create a new resource pool
  metalcloud-cli resource-pool create --label "Production Pool" --description "Pool for production workloads"

  # Add a server to a resource pool
  metalcloud-cli resource-pool add-server 123 456`,
	}

	resourcePoolListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all resource pools with optional filtering and pagination",
		Long: `List all resource pools in the system with optional filtering and pagination support.

This command displays a table of resource pools showing their ID, label, and description.
You can filter results using search terms and control the output with pagination parameters.

Flags:
  --page int      Page number for pagination (default: 0, shows all results)
  --limit int     Maximum number of records per page (default: 0, max: 100)
  --search string Search term to filter results by label or description

The search parameter performs a case-insensitive substring match against both
the resource pool label and description fields.

Examples:
  # List all resource pools
  metalcloud-cli resource-pool list

  # List resource pools with pagination
  metalcloud-cli resource-pool list --page 1 --limit 10

  # Search for resource pools containing "production"
  metalcloud-cli resource-pool list --search "production"

  # Combine search with pagination
  metalcloud-cli resource-pool list --search "dev" --page 1 --limit 5`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolList(cmd.Context(), resourcePoolFlags.pageFlag, resourcePoolFlags.limitFlag, resourcePoolFlags.searchFlag)
		},
	}

	resourcePoolGetCmd = &cobra.Command{
		Use:   "get <pool_id>",
		Short: "Get detailed information about a specific resource pool",
		Long: `Get detailed information about a specific resource pool by its ID.

This command retrieves and displays comprehensive information about a resource pool
including its ID, label, description, and any associated metadata.

Arguments:
  pool_id    The numeric ID of the resource pool to retrieve

Examples:
  # Get information about resource pool with ID 123
  metalcloud-cli resource-pool get 123

  # Get resource pool details using alias
  metalcloud-cli rp get 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolGet(cmd.Context(), args[0])
		},
	}

	resourcePoolCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new resource pool",
		Long: `Create a new resource pool with the specified label and description.

This command creates a new resource pool that can be used to group and organize
infrastructure resources. After creation, you can assign users, servers, and 
subnet pools to this resource pool.

Required Flags:
  --label string        Label/name for the resource pool (required)
  --description string  Description of the resource pool's purpose (required)

Examples:
  # Create a production resource pool
  metalcloud-cli resource-pool create --label "Production Pool" --description "Pool for production workloads"

  # Create a development resource pool
  metalcloud-cli resource-pool create --label "Dev Environment" --description "Development and testing resources"

  # Create using aliases
  metalcloud-cli rp create --label "QA Pool" --description "Quality assurance testing pool"`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolCreate(
				cmd.Context(),
				resourcePoolFlags.labelFlag,
				resourcePoolFlags.descriptionFlag,
			)
		},
	}

	resourcePoolDeleteCmd = &cobra.Command{
		Use:     "delete <pool_id>",
		Aliases: []string{"rm"},
		Short:   "Delete a resource pool",
		Long: `Delete a resource pool by its ID.

This command permanently removes a resource pool from the system. The resource pool
must be empty (no assigned users, servers, or subnet pools) before it can be deleted.

Arguments:
  pool_id    The numeric ID of the resource pool to delete

Examples:
  # Delete resource pool with ID 123
  metalcloud-cli resource-pool delete 123

  # Delete using alias
  metalcloud-cli rp rm 456

Note: This operation is irreversible. Ensure the resource pool is no longer needed
before deletion.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolDelete(cmd.Context(), args[0])
		},
	}

	resourcePoolGetUsersCmd = &cobra.Command{
		Use:     "get-users <pool_id>",
		Aliases: []string{"users"},
		Short:   "List users with access to a resource pool",
		Long: `List all users that have access to a specific resource pool.

This command retrieves and displays a list of users who have been granted access
to the specified resource pool. The output includes user details and their
permissions within the resource pool.

Arguments:
  pool_id    The numeric ID of the resource pool

Examples:
  # List users for resource pool with ID 123
  metalcloud-cli resource-pool get-users 123

  # List users using alias
  metalcloud-cli rp users 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolGetUsers(cmd.Context(), args[0])
		},
	}

	resourcePoolAddUserCmd = &cobra.Command{
		Use:   "add-user <pool_id> <user_id>",
		Short: "Grant a user access to a resource pool",
		Long: `Grant a user access to a resource pool by specifying the resource pool ID and user ID.

This command adds a user to a resource pool, giving them access to the resources
within that pool according to their role permissions.

Arguments:
  pool_id    The numeric ID of the resource pool
  user_id    The numeric ID of the user to add

Examples:
  # Add user with ID 789 to resource pool with ID 123
  metalcloud-cli resource-pool add-user 123 789

  # Add user using alias
  metalcloud-cli rp add-user 456 789`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolAddUser(cmd.Context(), args[0], args[1])
		},
	}

	resourcePoolRemoveUserCmd = &cobra.Command{
		Use:     "remove-user <pool_id> <user_id>",
		Aliases: []string{"rm-user"},
		Short:   "Revoke user access from a resource pool",
		Long: `Revoke a user's access from a resource pool by specifying the resource pool ID and user ID.

This command removes a user from a resource pool, revoking their access to the resources
within that pool. The user will no longer be able to interact with resources in this pool.

Arguments:
  pool_id    The numeric ID of the resource pool
  user_id    The numeric ID of the user to remove

Examples:
  # Remove user with ID 789 from resource pool with ID 123
  metalcloud-cli resource-pool remove-user 123 789

  # Remove user using alias
  metalcloud-cli rp rm-user 456 789`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolRemoveUser(cmd.Context(), args[0], args[1])
		},
	}

	resourcePoolGetServersCmd = &cobra.Command{
		Use:     "get-servers <pool_id>",
		Aliases: []string{"servers"},
		Short:   "List servers assigned to a resource pool",
		Long: `List all servers that are assigned to a specific resource pool.

This command retrieves and displays a list of servers that have been assigned
to the specified resource pool. The output includes server details and their
current status within the resource pool.

Arguments:
  pool_id    The numeric ID of the resource pool

Examples:
  # List servers for resource pool with ID 123
  metalcloud-cli resource-pool get-servers 123

  # List servers using alias
  metalcloud-cli rp servers 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolGetServers(cmd.Context(), args[0])
		},
	}

	resourcePoolAddServerCmd = &cobra.Command{
		Use:   "add-server <pool_id> <server_id>",
		Short: "Add a server to a resource pool",
		Long: `Add a server to a resource pool by specifying the resource pool ID and server ID.

This command assigns a server to a resource pool, making it available to users
who have access to that pool. The server must exist and not be assigned to another
resource pool.

Arguments:
  pool_id     The numeric ID of the resource pool
  server_id   The numeric ID of the server to add

Examples:
  # Add server with ID 456 to resource pool with ID 123
  metalcloud-cli resource-pool add-server 123 456

  # Add server using alias
  metalcloud-cli rp add-server 789 101112`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolAddServer(cmd.Context(), args[0], args[1])
		},
	}

	resourcePoolRemoveServerCmd = &cobra.Command{
		Use:     "remove-server <pool_id> <server_id>",
		Aliases: []string{"rm-server"},
		Short:   "Remove a server from a resource pool",
		Long: `Remove a server from a resource pool by specifying the resource pool ID and server ID.

This command unassigns a server from a resource pool, removing its association
with that pool. The server will no longer be available to users of this resource pool.

Arguments:
  pool_id     The numeric ID of the resource pool
  server_id   The numeric ID of the server to remove

Examples:
  # Remove server with ID 456 from resource pool with ID 123
  metalcloud-cli resource-pool remove-server 123 456

  # Remove server using alias
  metalcloud-cli rp rm-server 789 101112`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolRemoveServer(cmd.Context(), args[0], args[1])
		},
	}

	resourcePoolGetSubnetPoolsCmd = &cobra.Command{
		Use:     "get-subnet-pools <pool_id>",
		Aliases: []string{"subnet-pools", "subnets"},
		Short:   "List subnet pools assigned to a resource pool",
		Long: `List all subnet pools that are assigned to a specific resource pool.

This command retrieves and displays a list of subnet pools that have been assigned
to the specified resource pool. The output includes subnet pool details and their
current configuration within the resource pool.

Arguments:
  pool_id    The numeric ID of the resource pool

Examples:
  # List subnet pools for resource pool with ID 123
  metalcloud-cli resource-pool get-subnet-pools 123

  # List subnet pools using alias
  metalcloud-cli rp subnets 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolGetSubnetPools(cmd.Context(), args[0])
		},
	}

	resourcePoolAddSubnetPoolCmd = &cobra.Command{
		Use:   "add-subnet-pool <pool_id> <subnet_pool_id>",
		Short: "Add a subnet pool to a resource pool",
		Long: `Add a subnet pool to a resource pool by specifying the resource pool ID and subnet pool ID.

This command assigns a subnet pool to a resource pool, making the subnet pool's
network resources available to users who have access to that resource pool.

Arguments:
  pool_id         The numeric ID of the resource pool
  subnet_pool_id  The numeric ID of the subnet pool to add

Examples:
  # Add subnet pool with ID 789 to resource pool with ID 123
  metalcloud-cli resource-pool add-subnet-pool 123 789

  # Add subnet pool using alias
  metalcloud-cli rp add-subnet-pool 456 101112`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolAddSubnetPool(cmd.Context(), args[0], args[1])
		},
	}

	resourcePoolRemoveSubnetPoolCmd = &cobra.Command{
		Use:     "remove-subnet-pool <pool_id> <subnet_pool_id>",
		Aliases: []string{"rm-subnet-pool"},
		Short:   "Remove a subnet pool from a resource pool",
		Long: `Remove a subnet pool from a resource pool by specifying the resource pool ID and subnet pool ID.

This command unassigns a subnet pool from a resource pool, removing its association
with that pool. The subnet pool's network resources will no longer be available
to users of this resource pool.

Arguments:
  pool_id         The numeric ID of the resource pool
  subnet_pool_id  The numeric ID of the subnet pool to remove

Examples:
  # Remove subnet pool with ID 789 from resource pool with ID 123
  metalcloud-cli resource-pool remove-subnet-pool 123 789

  # Remove subnet pool using alias
  metalcloud-cli rp rm-subnet-pool 456 101112`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolRemoveSubnetPool(cmd.Context(), args[0], args[1])
		},
	}
)

func init() {
	rootCmd.AddCommand(resourcePoolCmd)

	// Resource Pool commands
	resourcePoolCmd.AddCommand(resourcePoolListCmd)
	resourcePoolCmd.AddCommand(resourcePoolGetCmd)
	resourcePoolCmd.AddCommand(resourcePoolCreateCmd)
	resourcePoolCmd.AddCommand(resourcePoolDeleteCmd)
	resourcePoolCmd.AddCommand(resourcePoolGetUsersCmd)
	resourcePoolCmd.AddCommand(resourcePoolAddUserCmd)
	resourcePoolCmd.AddCommand(resourcePoolRemoveUserCmd)
	resourcePoolCmd.AddCommand(resourcePoolGetServersCmd)
	resourcePoolCmd.AddCommand(resourcePoolAddServerCmd)
	resourcePoolCmd.AddCommand(resourcePoolRemoveServerCmd)
	resourcePoolCmd.AddCommand(resourcePoolGetSubnetPoolsCmd)
	resourcePoolCmd.AddCommand(resourcePoolAddSubnetPoolCmd)
	resourcePoolCmd.AddCommand(resourcePoolRemoveSubnetPoolCmd)

	// Add flags for list command
	resourcePoolListCmd.Flags().IntVar(&resourcePoolFlags.pageFlag, "page", 0, "Page number")
	resourcePoolListCmd.Flags().IntVar(&resourcePoolFlags.limitFlag, "limit", 0, "Number of records per page (max 100)")
	resourcePoolListCmd.Flags().StringVar(&resourcePoolFlags.searchFlag, "search", "", "Search term to filter results")

	// Add flags for create command
	resourcePoolCreateCmd.Flags().StringVar(&resourcePoolFlags.labelFlag, "label", "", "Resource pool label")
	resourcePoolCreateCmd.Flags().StringVar(&resourcePoolFlags.descriptionFlag, "description", "", "Resource pool description")

	// Required flags
	resourcePoolCreateCmd.MarkFlagRequired("label")
	resourcePoolCreateCmd.MarkFlagRequired("description")
}
