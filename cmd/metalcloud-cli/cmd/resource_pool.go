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
		Short:   "Resource pool management",
		Long:    `Resource pool management commands.`,
	}

	resourcePoolListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "Lists resource pools.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolList(cmd.Context(), resourcePoolFlags.pageFlag, resourcePoolFlags.limitFlag, resourcePoolFlags.searchFlag)
		},
	}

	resourcePoolGetCmd = &cobra.Command{
		Use:          "get <pool_id>",
		Short:        "Get resource pool information.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolGet(cmd.Context(), args[0])
		},
	}

	resourcePoolCreateCmd = &cobra.Command{
		Use:          "create",
		Short:        "Create a resource pool.",
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
		Use:          "delete <pool_id>",
		Aliases:      []string{"rm"},
		Short:        "Delete a resource pool.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolDelete(cmd.Context(), args[0])
		},
	}

	resourcePoolGetUsersCmd = &cobra.Command{
		Use:          "get-users <pool_id>",
		Aliases:      []string{"users"},
		Short:        "Get users that have access to a resource pool.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolGetUsers(cmd.Context(), args[0])
		},
	}

	resourcePoolAddUserCmd = &cobra.Command{
		Use:          "add-user <pool_id> <user_id>",
		Short:        "Add a user to a resource pool.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolAddUser(cmd.Context(), args[0], args[1])
		},
	}

	resourcePoolRemoveUserCmd = &cobra.Command{
		Use:          "remove-user <pool_id> <user_id>",
		Aliases:      []string{"rm-user"},
		Short:        "Remove a user from a resource pool.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolRemoveUser(cmd.Context(), args[0], args[1])
		},
	}

	resourcePoolGetServersCmd = &cobra.Command{
		Use:          "get-servers <pool_id>",
		Aliases:      []string{"servers"},
		Short:        "Get servers that are part of a resource pool.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolGetServers(cmd.Context(), args[0])
		},
	}

	resourcePoolAddServerCmd = &cobra.Command{
		Use:          "add-server <pool_id> <server_id>",
		Short:        "Add a server to a resource pool.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolAddServer(cmd.Context(), args[0], args[1])
		},
	}

	resourcePoolRemoveServerCmd = &cobra.Command{
		Use:          "remove-server <pool_id> <server_id>",
		Aliases:      []string{"rm-server"},
		Short:        "Remove a server from a resource pool.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolRemoveServer(cmd.Context(), args[0], args[1])
		},
	}

	resourcePoolGetSubnetPoolsCmd = &cobra.Command{
		Use:          "get-subnet-pools <pool_id>",
		Aliases:      []string{"subnet-pools", "subnets"},
		Short:        "Get subnet pools that are part of a resource pool.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolGetSubnetPools(cmd.Context(), args[0])
		},
	}

	resourcePoolAddSubnetPoolCmd = &cobra.Command{
		Use:          "add-subnet-pool <pool_id> <subnet_pool_id>",
		Short:        "Add a subnet pool to a resource pool.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_RESOURCE_POOLS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return resource_pool.ResourcePoolAddSubnetPool(cmd.Context(), args[0], args[1])
		},
	}

	resourcePoolRemoveSubnetPoolCmd = &cobra.Command{
		Use:          "remove-subnet-pool <pool_id> <subnet_pool_id>",
		Aliases:      []string{"rm-subnet-pool"},
		Short:        "Remove a subnet pool from a resource pool.",
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
