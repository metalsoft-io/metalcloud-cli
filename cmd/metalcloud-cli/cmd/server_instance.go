package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/server_instance"
	"github.com/spf13/cobra"
)

// Server Instance Group management commands.
var (
	serverInstanceGroupFlags = struct {
		label         string
		instanceCount int
		osTemplateId  int
		accessMode    string
		tagged        string
		redundancy    string
	}{}

	serverInstanceGroupCmd = &cobra.Command{
		Use:     "server-instance-group [command]",
		Aliases: []string{"ig", "instance-array", "ia"},
		Short:   "Manage server instance groups within infrastructures",
		Long: `Server Instance Group management commands.

Server Instance Groups are collections of server instances that share the same configuration
and can be managed as a single unit within an infrastructure. They provide scaling capabilities
and simplified management of multiple servers with identical specifications.

Available commands include:
- list: List all instance groups in an infrastructure
- get: View detailed configuration of a specific instance group
- create: Create a new instance group with specified parameters
- update: Modify existing instance group properties
- delete: Remove an instance group from the infrastructure
- instances: List all server instances within a group
- network: Manage network connections for instance groups

Use "metalcloud-cli server-instance-group [command] --help" for detailed information about each command.`,
	}

	serverInstanceGroupListCmd = &cobra.Command{
		Use:     "list infrastructure_id_or_label",
		Aliases: []string{"ls"},
		Short:   "List all server instance groups in an infrastructure",
		Long: `List all server instance groups in an infrastructure.

This command displays all server instance groups within a specified infrastructure,
showing their configuration details including ID, label, status, and timestamps.

Arguments:
  infrastructure_id_or_label  The infrastructure ID (numeric) or label (string) to list groups from

Examples:
  # List all instance groups in infrastructure with ID 1234
  metalcloud-cli server-instance-group list 1234

  # List all instance groups in infrastructure with label "prod-env"
  metalcloud-cli server-instance-group list prod-env

  # Using alias
  metalcloud-cli ig ls 1234`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupList(cmd.Context(), args[0])
		},
	}

	serverInstanceGroupGetCmd = &cobra.Command{
		Use:     "get server_instance_group_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a server instance group",
		Long: `Get detailed information about a server instance group.

This command retrieves and displays comprehensive information about a specific server
instance group, including its configuration, status, resource specifications, and
associated metadata.

Arguments:
  server_instance_group_id  The numeric ID of the server instance group to retrieve

Examples:
  # Get details of server instance group with ID 1234
  metalcloud-cli server-instance-group get 1234

  # Using alias
  metalcloud-cli ig show 1234`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupGet(cmd.Context(), args[0])
		},
	}

	serverInstanceGroupCreateCmd = &cobra.Command{
		Use:     "create infrastructure_id_or_label label server_type_id instance_count [os_template_id]",
		Aliases: []string{"new"},
		Short:   "Create a new server instance group in an infrastructure",
		Long: `Create a new server instance group in an infrastructure.

This command creates a new server instance group with the specified configuration.
The group will contain multiple server instances with identical specifications.

Arguments:
  infrastructure_id_or_label  The infrastructure ID (numeric) or label (string) to create the group in
  label                       Label for the new instance group
  server_type_id             Server type ID to use for instances in the group
  instance_count             Number of instances to create in the group
  os_template_id             (Optional) OS template ID to use for instances

Examples:
  # Create instance group with 3 instances using server type 100
  metalcloud-cli server-instance-group create 1234 web-servers 100 3

  # Create instance group with specific OS template
  metalcloud-cli server-instance-group create prod-env db-cluster 200 2 50

  # Using alias
  metalcloud-cli ig new 1234 app-servers 150 5`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_WRITE},
		Args:         cobra.RangeArgs(4, 5),
		RunE: func(cmd *cobra.Command, args []string) error {
			os_template_id := ""
			if len(args) == 5 {
				os_template_id = args[4]
			}

			return server_instance.ServerInstanceGroupCreate(cmd.Context(), args[0], args[1], args[2], args[3], os_template_id)
		},
	}

	serverInstanceGroupUpdateCmd = &cobra.Command{
		Use:     "update server_instance_group_id",
		Aliases: []string{"edit"},
		Short:   "Update server instance group configuration",
		Long: `Update server instance group configuration.

This command allows you to modify the configuration of an existing server instance group.
You can update the label, instance count, or OS template. At least one flag must be provided.

Arguments:
  server_instance_group_id  The numeric ID of the server instance group to update

Flags:
  --label string           Set the instance group label
  --instance-count int     Set the count of instance group instances (must be > 0)
  --os-template-id int     Set the instance group OS template ID (must be > 0)

Note: At least one of the flags (--label, --instance-count, --os-template-id) must be provided.

Examples:
  # Update the label of instance group 1234
  metalcloud-cli server-instance-group update 1234 --label "new-web-servers"

  # Scale instance group to 5 instances
  metalcloud-cli server-instance-group update 1234 --instance-count 5

  # Change OS template
  metalcloud-cli server-instance-group update 1234 --os-template-id 25

  # Update multiple properties at once
  metalcloud-cli server-instance-group update 1234 --label "updated-servers" --instance-count 3

  # Using alias
  metalcloud-cli ig edit 1234 --instance-count 10`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupUpdate(cmd.Context(), args[0], serverInstanceGroupFlags.label, serverInstanceGroupFlags.instanceCount, serverInstanceGroupFlags.osTemplateId)
		},
	}

	serverInstanceGroupDeleteCmd = &cobra.Command{
		Use:     "delete server_instance_group_id",
		Aliases: []string{"rm"},
		Short:   "Delete a server instance group from an infrastructure",
		Long: `Delete a server instance group from an infrastructure.

This command permanently removes a server instance group and all its associated instances
from the infrastructure. This action cannot be undone.

Arguments:
  server_instance_group_id  The numeric ID of the server instance group to delete

Warning: This operation will delete all server instances within the group and cannot be reversed.

Examples:
  # Delete server instance group with ID 1234
  metalcloud-cli server-instance-group delete 1234

  # Using alias
  metalcloud-cli ig rm 1234`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupDelete(cmd.Context(), args[0])
		},
	}

	serverInstanceGroupInstancesCmd = &cobra.Command{
		Use:     "instances server_instance_group_id",
		Aliases: []string{"instances-list", "instances-ls"},
		Short:   "List all server instances within a server instance group",
		Long: `List all server instances within a server instance group.

This command displays all server instances that belong to a specific server instance group,
showing their individual configurations, status, and resource assignments.

Arguments:
  server_instance_group_id  The numeric ID of the server instance group to list instances from

Examples:
  # List all instances in server instance group 1234
  metalcloud-cli server-instance-group instances 1234

  # Using alias
  metalcloud-cli ig instances-ls 1234`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupInstances(cmd.Context(), args[0])
		},
	}

	serverInstanceGroupNetworkCmd = &cobra.Command{
		Use:     "network [command]",
		Aliases: []string{"net"},
		Short:   "Manage network connections for server instance groups",
		Long: `Manage network connections for server instance groups.

This command group provides operations for managing network connections between
server instance groups and networks. You can list, view, create, update, and
delete network connections.

Available commands:
- list: List all network connections for a server instance group
- get: Get details of a specific network connection
- connect: Connect a server instance group to a network
- update: Update an existing network connection
- disconnect: Remove a network connection

Use "metalcloud-cli server-instance-group network [command] --help" for detailed information about each command.`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_READ},
	}

	serverInstanceGroupNetworkListCmd = &cobra.Command{
		Use:     "list server_instance_group_id",
		Aliases: []string{"ls"},
		Short:   "List all network connections for a server instance group",
		Long: `List all network connections for a server instance group.

This command displays all network connections associated with a specific server instance group,
showing connection details including network ID, subnet information, access mode, and redundancy settings.

Arguments:
  server_instance_group_id  The numeric ID of the server instance group to list network connections for

Examples:
  # List all network connections for server instance group 1234
  metalcloud-cli server-instance-group network list 1234

  # Using alias
  metalcloud-cli ig net ls 1234`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupNetworkList(cmd.Context(), args[0])
		},
	}

	serverInstanceGroupNetworkGetCmd = &cobra.Command{
		Use:     "get server_instance_group_id connection_id",
		Aliases: []string{"show"},
		Short:   "Get network connection details for a server instance group",
		Long: `Get network connection details for a server instance group.

This command retrieves and displays detailed information about a specific network connection
associated with a server instance group, including access mode, VLAN configuration, and
redundancy settings.

Arguments:
  server_instance_group_id  The numeric ID of the server instance group
  connection_id            The numeric ID of the network connection to retrieve

Examples:
  # Get details of network connection 5 for server instance group 1234
  metalcloud-cli server-instance-group network get 1234 5

  # Using alias
  metalcloud-cli ig net show 1234 5`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupNetworkGet(cmd.Context(), args[0], args[1])
		},
	}

	serverInstanceGroupNetworkConnectCmd = &cobra.Command{
		Use:     "connect server_instance_group_id network_id access_mode tagged [redundancy]",
		Aliases: []string{"new", "add"},
		Short:   "Connect a server instance group to a network",
		Long: `Connect a server instance group to a network.

This command creates a new network connection between a server instance group and a network,
configuring the access mode, VLAN tagging, and optionally redundancy settings.

Arguments:
  server_instance_group_id  The numeric ID of the server instance group to connect
  network_id               The ID of the network to connect to
  access_mode              Network access mode (e.g., "trunk", "access")
  tagged                   Whether VLAN tagging is enabled (true/false)
  redundancy               (Optional) Redundancy configuration (e.g., "active-backup", "load-balancing")

Examples:
  # Connect server instance group to network with trunk access
  metalcloud-cli server-instance-group network connect 1234 567 trunk true

  # Connect with access mode and no tagging
  metalcloud-cli server-instance-group network connect 1234 567 access false

  # Connect with redundancy configuration
  metalcloud-cli server-instance-group network connect 1234 567 trunk true active-backup

  # Using alias
  metalcloud-cli ig net add 1234 567 trunk true`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_WRITE},
		Args:         cobra.RangeArgs(4, 5),
		RunE: func(cmd *cobra.Command, args []string) error {
			redundancy := ""
			if len(args) == 5 {
				redundancy = args[4]
			}
			return server_instance.ServerInstanceGroupNetworkConnect(cmd.Context(), args[0], args[1], args[2], args[3], redundancy)
		},
	}

	serverInstanceGroupNetworkUpdateCmd = &cobra.Command{
		Use:     "update server_instance_group_id connection_id",
		Aliases: []string{"edit"},
		Short:   "Update network connection for a server instance group",
		Long: `Update network connection for a server instance group.

This command allows you to modify the configuration of an existing network connection
between a server instance group and a network. You can update the access mode, VLAN
tagging settings, or redundancy configuration.

Arguments:
  server_instance_group_id  The numeric ID of the server instance group
  connection_id            The numeric ID of the network connection to update

Flags:
  --access-mode string     Network connection access mode (e.g., "trunk", "access")
  --tagged string          Network connection VLAN tagging (true/false)
  --redundancy string      Network connection redundancy mode (e.g., "active-backup", "load-balancing")

Note: At least one of the flags (--access-mode, --tagged, --redundancy) must be provided.

Examples:
  # Update access mode to trunk
  metalcloud-cli server-instance-group network update 1234 5 --access-mode trunk

  # Enable VLAN tagging
  metalcloud-cli server-instance-group network update 1234 5 --tagged true

  # Set redundancy mode
  metalcloud-cli server-instance-group network update 1234 5 --redundancy active-backup

  # Update multiple properties at once
  metalcloud-cli server-instance-group network update 1234 5 --access-mode trunk --tagged true

  # Using alias
  metalcloud-cli ig net edit 1234 5 --access-mode access`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupNetworkUpdate(cmd.Context(), args[0], args[1], serverInstanceGroupFlags.accessMode, serverInstanceGroupFlags.tagged, serverInstanceGroupFlags.redundancy)
		},
	}

	serverInstanceGroupNetworkDisconnectCmd = &cobra.Command{
		Use:     "disconnect server_instance_group_id connection_id",
		Aliases: []string{"rm", "remove"},
		Short:   "Remove a network connection from a server instance group",
		Long: `Remove a network connection from a server instance group.

This command permanently removes a network connection between a server instance group
and a network. This action cannot be undone and will disconnect all instances in the
group from the specified network.

Arguments:
  server_instance_group_id  The numeric ID of the server instance group
  connection_id            The numeric ID of the network connection to remove

Warning: This operation will disconnect all instances in the group from the network and cannot be reversed.

Examples:
  # Remove network connection 5 from server instance group 1234
  metalcloud-cli server-instance-group network disconnect 1234 5

  # Using alias
  metalcloud-cli ig net rm 1234 5`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupNetworkDisconnect(cmd.Context(), args[0], args[1])
		},
	}
)

// Server Instance management commands.
var (
	serverInstanceCmd = &cobra.Command{
		Use:     "server-instance [command]",
		Aliases: []string{"inst"},
		Short:   "Manage individual server instances",
		Long: `Server Instance management commands.

Server Instances are individual compute resources within server instance groups.
They represent physical or virtual servers with specific hardware configurations
and network connections. Each instance inherits properties from its parent
instance group but can have individual characteristics and status.

Available commands include:
- get: View detailed information about a specific server instance

Use "metalcloud-cli server-instance [command] --help" for detailed information about each command.`,
	}

	serverInstanceGetCmd = &cobra.Command{
		Use:     "get server_instance_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a server instance",
		Long: `Get detailed information about a server instance.

This command retrieves and displays comprehensive information about a specific server
instance, including its configuration, status, hardware specifications, network
connections, and metadata. The instance may be part of a server instance group
or standalone.

Arguments:
  server_instance_id  The numeric ID of the server instance to retrieve

Examples:
  # Get details of server instance with ID 5678
  metalcloud-cli server-instance get 5678

  # Using alias
  metalcloud-cli inst show 5678`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGet(cmd.Context(), args[0])
		},
	}
)

func init() {
	// Server Instance Group management commands.
	rootCmd.AddCommand(serverInstanceGroupCmd)

	serverInstanceGroupCmd.AddCommand(serverInstanceGroupListCmd)

	serverInstanceGroupCmd.AddCommand(serverInstanceGroupGetCmd)

	serverInstanceGroupCmd.AddCommand(serverInstanceGroupCreateCmd)

	serverInstanceGroupCmd.AddCommand(serverInstanceGroupUpdateCmd)
	serverInstanceGroupUpdateCmd.Flags().StringVar(&serverInstanceGroupFlags.label, "label", "", "Set the instance group label.")
	serverInstanceGroupUpdateCmd.Flags().IntVar(&serverInstanceGroupFlags.instanceCount, "instance-count", 0, "Set the count of instance group instances.")
	serverInstanceGroupUpdateCmd.Flags().IntVar(&serverInstanceGroupFlags.osTemplateId, "os-template-id", 0, "Set the instance group OS template Id.")
	serverInstanceGroupUpdateCmd.MarkFlagsOneRequired("label", "instance-count", "os-template-id")

	serverInstanceGroupCmd.AddCommand(serverInstanceGroupDeleteCmd)

	serverInstanceGroupCmd.AddCommand(serverInstanceGroupInstancesCmd)

	serverInstanceGroupCmd.AddCommand(serverInstanceGroupNetworkCmd)

	serverInstanceGroupNetworkCmd.AddCommand(serverInstanceGroupNetworkListCmd)

	serverInstanceGroupNetworkCmd.AddCommand(serverInstanceGroupNetworkGetCmd)

	serverInstanceGroupNetworkCmd.AddCommand(serverInstanceGroupNetworkConnectCmd)

	serverInstanceGroupNetworkCmd.AddCommand(serverInstanceGroupNetworkUpdateCmd)
	serverInstanceGroupNetworkUpdateCmd.Flags().StringVar(&serverInstanceGroupFlags.accessMode, "access-mode", "", "Network connection access mode.")
	serverInstanceGroupNetworkUpdateCmd.Flags().StringVar(&serverInstanceGroupFlags.tagged, "tagged", "", "Network connection tagged.")
	serverInstanceGroupNetworkUpdateCmd.Flags().StringVar(&serverInstanceGroupFlags.redundancy, "redundancy", "", "Network connection redundancy.")
	serverInstanceGroupNetworkUpdateCmd.MarkFlagsOneRequired("access-mode", "tagged", "redundancy")

	serverInstanceGroupNetworkCmd.AddCommand(serverInstanceGroupNetworkDisconnectCmd)

	// Server Instance management commands.
	rootCmd.AddCommand(serverInstanceCmd)

	serverInstanceCmd.AddCommand(serverInstanceGetCmd)
}
