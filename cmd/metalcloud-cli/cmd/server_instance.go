package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/server_instance"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
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

		configSource             string
		filterInfrastructureId   []string
		filterServiceStatus      []string
		filterConfigDeployStatus []string
		filterConfigDeployType   []string

		aclRuleType           string
		aclDirection          string
		aclSequence           int32
		aclForwardingAction   string
		aclEnforcementPoint   string
		aclNetworkProtocol    string
		aclSourceAddress      string
		aclDestinationAddress string
		aclSourcePort         string
		aclDestinationPort    string
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
- update-meta: Update the metadata (tags) of an instance group
- instances: List all server instances within a group
- drive-groups: List the drive groups attached to an instance group
- interfaces / interface: List or inspect the network interfaces of an instance group
- network: Manage the network configuration, connections and security rules of instance groups

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
- config: Get the network configuration (network endpoint group) of an instance group
- replace: Create or replace the network configuration of an instance group
- list: List all network connections for a server instance group
- get: Get details of a specific network connection
- connect: Connect a server instance group to a network
- update: Update an existing network connection
- disconnect: Remove a network connection
- acl: Manage the security rules of a network connection

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

	serverInstanceGroupUpdateMetaCmd = &cobra.Command{
		Use:     "update-meta server_instance_group_id",
		Aliases: []string{"meta"},
		Short:   "Update the metadata of a server instance group",
		Long: `Update the metadata (tags) of a server instance group.

Required Arguments:
  server_instance_group_id  The numeric ID of the server instance group

Required Flags:
  --config-source string   Source of the metadata. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud-cli server-instance-group update-meta 1234 --config-source meta.json
  echo '{"tags":["prod"]}' | metalcloud-cli ig update-meta 1234 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(serverInstanceGroupFlags.configSource)
			if err != nil {
				return err
			}

			return server_instance.ServerInstanceGroupMetaUpdate(cmd.Context(), args[0], config)
		},
	}

	serverInstanceGroupDriveGroupsCmd = &cobra.Command{
		Use:     "drive-groups server_instance_group_id",
		Aliases: []string{"drives"},
		Short:   "List the drive groups of a server instance group",
		Long: `List all drive groups attached to a server instance group.

Required Arguments:
  server_instance_group_id  The numeric ID of the server instance group

Examples:
  metalcloud-cli server-instance-group drive-groups 1234
  metalcloud-cli ig drives 1234`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupDriveGroups(cmd.Context(), args[0])
		},
	}

	serverInstanceGroupInterfacesCmd = &cobra.Command{
		Use:     "interfaces server_instance_group_id",
		Aliases: []string{"ifaces"},
		Short:   "List the interfaces of a server instance group",
		Long: `List all network interfaces of a server instance group.

Required Arguments:
  server_instance_group_id  The numeric ID of the server instance group

Optional Flags:
  --filter-infrastructure-id strings     Filter by infrastructure ID. Repeatable or comma-separated.
  --filter-service-status strings        Filter by service status.
  --filter-config-deploy-status strings  Filter by the deploy status of the pending configuration.
  --filter-config-deploy-type strings    Filter by the deploy type of the pending configuration.

Examples:
  metalcloud-cli server-instance-group interfaces 1234
  metalcloud-cli ig ifaces 1234 --filter-service-status active`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupInterfaces(cmd.Context(), args[0], server_instance.ServerInstanceFilters{
				InfrastructureId:   serverInstanceGroupFlags.filterInfrastructureId,
				ServiceStatus:      serverInstanceGroupFlags.filterServiceStatus,
				ConfigDeployStatus: serverInstanceGroupFlags.filterConfigDeployStatus,
				ConfigDeployType:   serverInstanceGroupFlags.filterConfigDeployType,
			})
		},
	}

	serverInstanceGroupInterfaceGetCmd = &cobra.Command{
		Use:     "interface server_instance_group_id interface_id",
		Aliases: []string{"iface"},
		Short:   "Get one interface of a server instance group",
		Long: `Get the details of one network interface of a server instance group.

Required Arguments:
  server_instance_group_id  The numeric ID of the server instance group
  interface_id              The numeric ID of the interface

Examples:
  metalcloud-cli server-instance-group interface 1234 7
  metalcloud-cli ig iface 1234 7`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupInterfaceGet(cmd.Context(), args[0], args[1])
		},
	}

	serverInstanceGroupNetworkConfigCmd = &cobra.Command{
		Use:     "config server_instance_group_id",
		Aliases: []string{"get-config"},
		Short:   "Get the network configuration of a server instance group",
		Long: `Get the network endpoint group that holds the network configuration of a server
instance group. Use 'network list' for the individual network connections.

Required Arguments:
  server_instance_group_id  The numeric ID of the server instance group

Examples:
  metalcloud-cli server-instance-group network config 1234
  metalcloud-cli ig net get-config 1234`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupNetworkConfiguration(cmd.Context(), args[0])
		},
	}

	serverInstanceGroupNetworkReplaceCmd = &cobra.Command{
		Use:   "replace server_instance_group_id",
		Short: "Create or replace the network configuration of a server instance group",
		Long: `Create or replace (PUT) the network configuration of a server instance group.

Without --config-source the operation is sent with no body, which creates the network
configuration if it does not exist yet and returns it otherwise. With --config-source
the supplied document is sent as the request body together with the If-Match of the
group configuration.

Required Arguments:
  server_instance_group_id  The numeric ID of the server instance group

Optional Flags:
  --config-source string   Source of the network configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud-cli server-instance-group network replace 1234
  metalcloud-cli server-instance-group network replace 1234 --config-source networking.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(serverInstanceGroupFlags.configSource)
			if err != nil {
				return err
			}

			return server_instance.ServerInstanceGroupNetworkReplace(cmd.Context(), args[0], config)
		},
	}

	serverInstanceGroupACLCmd = &cobra.Command{
		Use:     "acl [command]",
		Aliases: []string{"security"},
		Short:   "Manage the security rules of a network connection",
		Long: `Manage the security rules (ACLs) of one network connection of a server instance
group.

Available commands:
- list: List the security rules of a network connection
- get: Get one security rule
- add: Add a security rule
- update: Update a security rule
- remove: Remove a security rule
- config-example: Print an example security rule configuration

Use "metalcloud-cli server-instance-group network acl [command] --help" for detailed information about each command.`,
	}

	serverInstanceGroupACLListCmd = &cobra.Command{
		Use:     "list server_instance_group_id connection_id",
		Aliases: []string{"ls"},
		Short:   "List the security rules of a network connection",
		Long: `List the security rules of one network connection of a server instance group.

Required Arguments:
  server_instance_group_id  The numeric ID of the server instance group
  connection_id             The numeric ID of the network connection

Examples:
  metalcloud-cli server-instance-group network acl list 1234 5
  metalcloud-cli ig net acl ls 1234 5`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupACLList(cmd.Context(), args[0], args[1])
		},
	}

	serverInstanceGroupACLGetCmd = &cobra.Command{
		Use:     "get server_instance_group_id connection_id rule_id",
		Aliases: []string{"show"},
		Short:   "Get a security rule of a network connection",
		Long: `Get the details of one security rule of a network connection.

Required Arguments:
  server_instance_group_id  The numeric ID of the server instance group
  connection_id             The numeric ID of the network connection
  rule_id                   The numeric ID of the security rule

Examples:
  metalcloud-cli server-instance-group network acl get 1234 5 3
  metalcloud-cli ig net acl show 1234 5 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupACLGet(cmd.Context(), args[0], args[1], args[2])
		},
	}

	serverInstanceGroupACLConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Print an example security rule configuration",
		Long: `Print an example configuration that can be edited and passed to
'server-instance-group network acl add --config-source'.

Examples:
  metalcloud-cli server-instance-group network acl config-example > rule.json
  metalcloud-cli ig net acl config-example -f yaml`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupACLConfigExample(cmd.Context())
		},
	}

	serverInstanceGroupACLAddCmd = &cobra.Command{
		Use:     "add server_instance_group_id connection_id",
		Aliases: []string{"create", "new"},
		Short:   "Add a security rule to a network connection",
		Long: `Add a security rule to one network connection of a server instance group.

The rule can be described either with a configuration file (--config-source) or with
individual flags (--rule-type, --direction, ...).

Required Arguments:
  server_instance_group_id  The numeric ID of the server instance group
  connection_id             The numeric ID of the network connection

Required Flags (one of):
  --config-source string       Source of the new rule. Can be 'pipe' or path to a JSON/YAML file.
  --rule-type string           Rule type: ipv4, ipv6 or mac

Optional Flags (when not using --config-source):
  --direction string           Rule direction: in or out (default "in")
  --sequence int               Evaluation order of the rule
  --forwarding-action string   Forwarding action: allow, deny, transit or discard (default "allow")
  --enforcement-point string   Enforcement point of the rule (default "svi")
  --network-protocol string    Network protocol of the rule
  --source-address string      Source address of the rule
  --destination-address string Destination address of the rule
  --source-port string         Source port of the rule
  --destination-port string    Destination port of the rule

Examples:
  metalcloud-cli server-instance-group network acl add 1234 5 --rule-type ipv4 --sequence 10 --source-address 10.0.0.0/24
  metalcloud-cli ig net acl add 1234 5 --config-source rule.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var create sdk.CreateLogicalNetworkACL

			if serverInstanceGroupFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(serverInstanceGroupFlags.configSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &create); err != nil {
					return err
				}
			} else {
				create = sdk.CreateLogicalNetworkACL{
					RuleType:         sdk.ACLType(serverInstanceGroupFlags.aclRuleType),
					Direction:        sdk.ACLDirection(serverInstanceGroupFlags.aclDirection),
					Sequence:         serverInstanceGroupFlags.aclSequence,
					ForwardingAction: sdk.ACLForwardingAction(serverInstanceGroupFlags.aclForwardingAction),
					EnforcementPoint: sdk.ACLEnforcementPoint(serverInstanceGroupFlags.aclEnforcementPoint),
				}
				if serverInstanceGroupFlags.aclNetworkProtocol != "" {
					create.NetworkProtocol = sdk.PtrString(serverInstanceGroupFlags.aclNetworkProtocol)
				}
				if serverInstanceGroupFlags.aclSourceAddress != "" {
					create.SourceAddress = sdk.PtrString(serverInstanceGroupFlags.aclSourceAddress)
				}
				if serverInstanceGroupFlags.aclDestinationAddress != "" {
					create.DestinationAddress = sdk.PtrString(serverInstanceGroupFlags.aclDestinationAddress)
				}
				if serverInstanceGroupFlags.aclSourcePort != "" {
					create.SourcePort = sdk.PtrString(serverInstanceGroupFlags.aclSourcePort)
				}
				if serverInstanceGroupFlags.aclDestinationPort != "" {
					create.DestinationPort = sdk.PtrString(serverInstanceGroupFlags.aclDestinationPort)
				}
			}

			return server_instance.ServerInstanceGroupACLAdd(cmd.Context(), args[0], args[1], create)
		},
	}

	serverInstanceGroupACLUpdateCmd = &cobra.Command{
		Use:     "update server_instance_group_id connection_id rule_id",
		Aliases: []string{"edit"},
		Short:   "Update a security rule of a network connection",
		Long: `Update a security rule of one network connection of a server instance group.

Required Arguments:
  server_instance_group_id  The numeric ID of the server instance group
  connection_id             The numeric ID of the network connection
  rule_id                   The numeric ID of the security rule

Required Flags:
  --config-source string   Source of the updated rule. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud-cli server-instance-group network acl update 1234 5 3 --config-source rule.json
  echo '{"forwardingAction":"deny"}' | metalcloud-cli ig net acl update 1234 5 3 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(serverInstanceGroupFlags.configSource)
			if err != nil {
				return err
			}

			return server_instance.ServerInstanceGroupACLUpdate(cmd.Context(), args[0], args[1], args[2], config)
		},
	}

	serverInstanceGroupACLRemoveCmd = &cobra.Command{
		Use:     "remove server_instance_group_id connection_id rule_id",
		Aliases: []string{"rm", "delete"},
		Short:   "Remove a security rule from a network connection",
		Long: `Remove a security rule from one network connection of a server instance group.

Required Arguments:
  server_instance_group_id  The numeric ID of the server instance group
  connection_id             The numeric ID of the network connection
  rule_id                   The numeric ID of the security rule

Examples:
  metalcloud-cli server-instance-group network acl remove 1234 5 3
  metalcloud-cli ig net acl rm 1234 5 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupACLRemove(cmd.Context(), args[0], args[1], args[2])
		},
	}
)

// Server Instance management commands.
var (
	serverInstanceFlags = struct {
		configSource             string
		filterInfrastructureId   []string
		filterGroupId            []string
		filterServerId           []string
		filterServiceStatus      []string
		filterConfigServerId     []string
		filterConfigDeployStatus []string
		filterConfigDeployType   []string

		// The interface listing has its own filter fields so that the
		// values set on 'list' can never leak into 'interfaces'.
		ifaceFilterInfrastructureId   []string
		ifaceFilterServiceStatus      []string
		ifaceFilterConfigDeployStatus []string
		ifaceFilterConfigDeployType   []string

		label        string
		groupId      int64
		serverTypeId int64
		hostname     string
		osTemplateId int64
		tags         []string
		usage        string
	}{}

	serverInstanceCmd = &cobra.Command{
		Use:     "server-instance [command]",
		Aliases: []string{"inst"},
		Short:   "Manage individual server instances",
		Long: `Server Instance management commands.

Server Instances are individual compute resources within server instance groups.
They represent physical or virtual servers with specific hardware configurations
and network connections. Each instance inherits properties from its parent
instance group but can have individual characteristics and status.

Command categories:
  Lifecycle:      list, get, create, delete, config-example
  Configuration:  config, update-config, update-meta, variables, os-installation-data
  Power:          power, power-status-batch, power-set-batch
  Operations:     reset, reinstall-os, credentials
  Hardware:       drives, interfaces, interface, update-interface-config
  Reporting:      statistics

Use "metalcloud-cli server-instance [command] --help" for detailed information about each command.`,
	}

	serverInstancePowerCmd = &cobra.Command{
		Use:   "power <server_instance_id> <on|off|reset|soft|status>",
		Short: "Set power state of a server instance",
		Long: `Set the power state of a server instance.

This command sends a power command to a server instance via the orchestration layer.
For direct BMC/IPMI power control, use the 'server power' command instead.

Valid power actions:
  on     - Power on the server instance
  off    - Power off the server instance
  reset  - Hard reset the server instance
  soft   - Graceful shutdown of the server instance
  status - Get the current power status of the server instance

Arguments:
  server_instance_id  The numeric ID of the server instance
  action              Power action to perform (on, off, reset, soft, status)

Subcommands:
  status              Get the current power status of a server instance

Examples:
  # Power on server instance 5678
  metalcloud-cli server-instance power 5678 on

  # Gracefully shutdown server instance
  metalcloud-cli server-instance power 5678 soft

  # Hard reset server instance
  metalcloud-cli inst power 5678 reset

  # Get power status
  metalcloud-cli server-instance power 5678 status`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[1] == "status" {
				return server_instance.ServerInstancePowerStatus(cmd.Context(), args[0])
			}
			return server_instance.ServerInstancePower(cmd.Context(), args[0], args[1])
		},
	}

	serverInstancePowerStatusCmd = &cobra.Command{
		Use:     "status <server_instance_id>",
		Aliases: []string{"power-state"},
		Short:   "Get power status of a server instance",
		Long: `Get the current power status of a server instance.

Arguments:
  server_instance_id  The numeric ID of the server instance

Examples:
  # Get power status of server instance 5678
  metalcloud-cli server-instance power status 5678

  # Using alias
  metalcloud-cli inst power power-state 5678`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstancePowerStatus(cmd.Context(), args[0])
		},
	}

	serverInstanceCredentialsCmd = &cobra.Command{
		Use:     "credentials <server_instance_id>",
		Aliases: []string{"creds"},
		Short:   "Get login credentials for a server instance",
		Long: `Get the login credentials for a server instance.

This command retrieves the initial credentials configured for a server instance,
including username, password, and SSH public key if available.

Arguments:
  server_instance_id  The numeric ID of the server instance

Examples:
  # Get credentials for server instance 5678
  metalcloud-cli server-instance credentials 5678

  # Using alias
  metalcloud-cli inst creds 5678`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceCredentials(cmd.Context(), args[0])
		},
	}

	serverInstanceReinstallOSCmd = &cobra.Command{
		Use:     "reinstall-os <server_instance_id>",
		Aliases: []string{"reinstall"},
		Short:   "Schedule OS reinstall for a server instance",
		Long: `Schedule an OS reinstallation for a server instance.

This command marks the server instance for OS reinstallation. The reinstall
will take effect at the next infrastructure deploy.

Arguments:
  server_instance_id  The numeric ID of the server instance

Examples:
  # Schedule OS reinstall for server instance 5678
  metalcloud-cli server-instance reinstall-os 5678

  # Using alias
  metalcloud-cli inst reinstall 5678`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceReinstallOS(cmd.Context(), args[0])
		},
	}

	serverInstanceConfigCmd = &cobra.Command{
		Use:     "config <server_instance_id>",
		Aliases: []string{"get-config"},
		Short:   "Get configuration of a server instance",
		Long: `Get the current configuration of a server instance.

This command retrieves the configuration details of a server instance including
server type, OS template, hostname, deploy status, and other settings.

Arguments:
  server_instance_id  The numeric ID of the server instance

Examples:
  # Get configuration of server instance 5678
  metalcloud-cli server-instance config 5678

  # Using alias
  metalcloud-cli inst get-config 5678`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceConfig(cmd.Context(), args[0])
		},
	}

	serverInstanceListCmd = &cobra.Command{
		Use:     "list [infrastructure_id_or_label]",
		Aliases: []string{"ls"},
		Short:   "List server instances",
		Long: `List server instances.

Without an argument every server instance visible to the user is listed. Pass an
infrastructure ID or label to restrict the listing to a single infrastructure.

Optional Arguments:
  infrastructure_id_or_label  List only the server instances of this infrastructure

Optional Flags:
  --filter-infrastructure-id strings     Filter by infrastructure ID. Repeatable or comma-separated.
  --filter-group-id strings              Filter by server instance group ID.
  --filter-server-id strings             Filter by server ID.
  --filter-service-status strings        Filter by service status.
  --filter-config-server-id strings      Filter by the server ID of the pending configuration.
  --filter-config-deploy-status strings  Filter by the deploy status of the pending configuration.
  --filter-config-deploy-type strings    Filter by the deploy type of the pending configuration.

Examples:
  metalcloud-cli server-instance list
  metalcloud-cli server-instance list 1234
  metalcloud-cli inst ls prod-env --filter-service-status active`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_READ},
		Args:         cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			infrastructureIdOrLabel := ""
			if len(args) == 1 {
				infrastructureIdOrLabel = args[0]
			}

			return server_instance.ServerInstanceList(cmd.Context(), infrastructureIdOrLabel, server_instance.ServerInstanceFilters{
				InfrastructureId:   serverInstanceFlags.filterInfrastructureId,
				GroupId:            serverInstanceFlags.filterGroupId,
				ServerId:           serverInstanceFlags.filterServerId,
				ServiceStatus:      serverInstanceFlags.filterServiceStatus,
				ConfigServerId:     serverInstanceFlags.filterConfigServerId,
				ConfigDeployStatus: serverInstanceFlags.filterConfigDeployStatus,
				ConfigDeployType:   serverInstanceFlags.filterConfigDeployType,
			})
		},
	}

	serverInstanceConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Print an example server instance create configuration",
		Long: `Print an example configuration that can be edited and passed to
'server-instance create --config-source'.

Examples:
  metalcloud-cli server-instance config-example > server-instance.json
  metalcloud-cli inst config-example -f yaml`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceConfigExample(cmd.Context())
		},
	}

	serverInstanceCreateCmd = &cobra.Command{
		Use:     "create infrastructure_id_or_label",
		Aliases: []string{"new"},
		Short:   "Create a server instance in an infrastructure",
		Long: `Create a new server instance in an infrastructure.

The instance can be described either with a configuration file (--config-source) or
with individual flags (--label, --group-id, ...).

Required Arguments:
  infrastructure_id_or_label  The ID or label of the infrastructure

Required Flags (one of):
  --config-source string   Source of the new server instance configuration. Can be 'pipe' or path to a JSON/YAML file.
  --label string           Label of the new server instance

Optional Flags (when not using --config-source):
  --group-id int           ID of the server instance group the instance belongs to
  --server-type-id int     ID of the server type to allocate
  --hostname string        Custom hostname (subdomain) of the instance
  --os-template-id int     ID of the OS template to deploy
  --tags strings           Tags of the new server instance

Examples:
  metalcloud-cli server-instance create 1234 --label web-01 --server-type-id 100
  metalcloud-cli server-instance create prod-env --config-source server-instance.json
  cat server-instance.yaml | metalcloud-cli inst new prod-env --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var create sdk.ServerInstanceCreate

			if serverInstanceFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(serverInstanceFlags.configSource)
				if err != nil {
					return err
				}
				if err := utils.UnmarshalContent(config, &create); err != nil {
					return err
				}
			} else {
				create = sdk.ServerInstanceCreate{
					Label: sdk.PtrString(serverInstanceFlags.label),
				}
				if serverInstanceFlags.groupId != 0 {
					create.GroupId = sdk.PtrInt64(serverInstanceFlags.groupId)
				}
				if serverInstanceFlags.serverTypeId != 0 {
					create.ServerTypeId = sdk.PtrInt64(serverInstanceFlags.serverTypeId)
				}
				if serverInstanceFlags.hostname != "" {
					create.Hostname = sdk.PtrString(serverInstanceFlags.hostname)
				}
				if serverInstanceFlags.osTemplateId != 0 {
					create.OsTemplateId = sdk.PtrInt64(serverInstanceFlags.osTemplateId)
				}
				if len(serverInstanceFlags.tags) > 0 {
					create.Tags = serverInstanceFlags.tags
				}
			}

			return server_instance.ServerInstanceCreate(cmd.Context(), args[0], create)
		},
	}

	serverInstanceDeleteCmd = &cobra.Command{
		Use:     "delete server_instance_id",
		Aliases: []string{"rm"},
		Short:   "Delete a server instance",
		Long: `Delete a server instance.

The current revision of the instance is fetched first and sent as If-Match, so the
delete fails if the instance changed in the meantime.

Required Arguments:
  server_instance_id  The numeric ID of the server instance

Examples:
  metalcloud-cli server-instance delete 5678
  metalcloud-cli inst rm 5678`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceDelete(cmd.Context(), args[0])
		},
	}

	serverInstanceUpdateConfigCmd = &cobra.Command{
		Use:     "update-config server_instance_id",
		Aliases: []string{"update"},
		Short:   "Update the configuration of a server instance",
		Long: `Update the pending configuration of a server instance.

The current configuration is fetched first and its revision is sent as If-Match,
because writes below the /config path are guarded by the configuration revision.
Use 'server-instance config' to read the current configuration.

Required Arguments:
  server_instance_id  The numeric ID of the server instance

Required Flags:
  --config-source string   Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud-cli server-instance update-config 5678 --config-source update.json
  echo '{"label":"new-label"}' | metalcloud-cli inst update-config 5678 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(serverInstanceFlags.configSource)
			if err != nil {
				return err
			}

			return server_instance.ServerInstanceConfigUpdate(cmd.Context(), args[0], config)
		},
	}

	serverInstanceUpdateMetaCmd = &cobra.Command{
		Use:     "update-meta server_instance_id",
		Aliases: []string{"meta"},
		Short:   "Update the metadata of a server instance",
		Long: `Update the metadata (tags) of a server instance.

Required Arguments:
  server_instance_id  The numeric ID of the server instance

Required Flags:
  --config-source string   Source of the metadata. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud-cli server-instance update-meta 5678 --config-source meta.json
  echo '{"tags":["prod"]}' | metalcloud-cli inst update-meta 5678 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(serverInstanceFlags.configSource)
			if err != nil {
				return err
			}

			return server_instance.ServerInstanceMetaUpdate(cmd.Context(), args[0], config)
		},
	}

	serverInstanceResetCmd = &cobra.Command{
		Use:   "reset server_instance_id",
		Short: "Reset the deployed server of a server instance",
		Long: `Reset the deployed server of a server instance. The operation is executed
immediately.

This is a different operation from 'server-instance power <id> reset': the power
command goes to the power-set endpoint and only cycles the power of the server,
while 'reset' asks the orchestration layer to reset the deployed server.

Required Arguments:
  server_instance_id  The numeric ID of the server instance

Examples:
  metalcloud-cli server-instance reset 5678
  metalcloud-cli inst reset 5678`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceReset(cmd.Context(), args[0])
		},
	}

	serverInstanceDrivesCmd = &cobra.Command{
		Use:   "drives server_instance_id",
		Short: "List the drives of a server instance",
		Long: `List all drives attached to a server instance.

Required Arguments:
  server_instance_id  The numeric ID of the server instance

Examples:
  metalcloud-cli server-instance drives 5678
  metalcloud-cli inst drives 5678`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceDrives(cmd.Context(), args[0])
		},
	}

	serverInstanceInterfacesCmd = &cobra.Command{
		Use:     "interfaces server_instance_id",
		Aliases: []string{"ifaces"},
		Short:   "List the interfaces of a server instance",
		Long: `List all network interfaces of a server instance.

Required Arguments:
  server_instance_id  The numeric ID of the server instance

Optional Flags:
  --filter-infrastructure-id strings     Filter by infrastructure ID. Repeatable or comma-separated.
  --filter-service-status strings        Filter by service status.
  --filter-config-deploy-status strings  Filter by the deploy status of the pending configuration.
  --filter-config-deploy-type strings    Filter by the deploy type of the pending configuration.

Examples:
  metalcloud-cli server-instance interfaces 5678
  metalcloud-cli inst ifaces 5678 --filter-service-status active`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceInterfaces(cmd.Context(), args[0], server_instance.ServerInstanceFilters{
				InfrastructureId:   serverInstanceFlags.ifaceFilterInfrastructureId,
				ServiceStatus:      serverInstanceFlags.ifaceFilterServiceStatus,
				ConfigDeployStatus: serverInstanceFlags.ifaceFilterConfigDeployStatus,
				ConfigDeployType:   serverInstanceFlags.ifaceFilterConfigDeployType,
			})
		},
	}

	serverInstanceInterfaceGetCmd = &cobra.Command{
		Use:     "interface server_instance_id interface_id",
		Aliases: []string{"iface"},
		Short:   "Get one interface of a server instance",
		Long: `Get the details of one network interface of a server instance.

Required Arguments:
  server_instance_id  The numeric ID of the server instance
  interface_id        The numeric ID of the interface

Examples:
  metalcloud-cli server-instance interface 5678 7
  metalcloud-cli inst iface 5678 7`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceInterfaceGet(cmd.Context(), args[0], args[1])
		},
	}

	serverInstanceUpdateInterfaceConfigCmd = &cobra.Command{
		Use:     "update-interface-config server_instance_id interface_id",
		Aliases: []string{"update-iface"},
		Short:   "Update the configuration of a server instance interface",
		Long: `Update the pending configuration of one network interface of a server instance.

The interface is fetched first and the revision of its configuration is sent as
If-Match, because writes below the /config path are guarded by the configuration
revision.

Required Arguments:
  server_instance_id  The numeric ID of the server instance
  interface_id        The numeric ID of the interface

Required Flags:
  --config-source string   Source of the updated interface configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud-cli server-instance update-interface-config 5678 7 --config-source iface.json
  echo '{"networkId":42}' | metalcloud-cli inst update-iface 5678 7 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(serverInstanceFlags.configSource)
			if err != nil {
				return err
			}

			return server_instance.ServerInstanceInterfaceConfigUpdate(cmd.Context(), args[0], args[1], config)
		},
	}

	serverInstanceOSInstallationDataCmd = &cobra.Command{
		Use:     "os-installation-data server_instance_id",
		Aliases: []string{"os-data"},
		Short:   "Get the OS installation data of a server instance",
		Long: `Get the OS installation data of a server instance: the site, server, instance,
group, infrastructure, drive and network values the OS installer is rendered with.

Required Arguments:
  server_instance_id  The numeric ID of the server instance

Optional Flags:
  --usage string   Restrict the returned values to one usage type (HTTPRequest, JavaScript,
                   APICall, AnsibleBundle, SSHExec, Copy, OSAsset).

Examples:
  metalcloud-cli server-instance os-installation-data 5678
  metalcloud-cli inst os-data 5678 --usage AnsibleBundle`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceOSInstallationData(cmd.Context(), args[0], serverInstanceFlags.usage)
		},
	}

	serverInstanceVariablesCmd = &cobra.Command{
		Use:     "variables server_instance_id",
		Aliases: []string{"vars"},
		Short:   "Get the variables of a server instance",
		Long: `Get the variables of a server instance: the site, server, instance, group,
infrastructure, drive and network values available to extensions and templates.

Required Arguments:
  server_instance_id  The numeric ID of the server instance

Optional Flags:
  --usage string   Restrict the returned values to one usage type (HTTPRequest, JavaScript,
                   APICall, AnsibleBundle, SSHExec, Copy, OSAsset).

Examples:
  metalcloud-cli server-instance variables 5678
  metalcloud-cli inst vars 5678 --usage APICall`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceVariables(cmd.Context(), args[0], serverInstanceFlags.usage)
		},
	}

	serverInstanceStatisticsCmd = &cobra.Command{
		Use:     "statistics",
		Aliases: []string{"stats"},
		Short:   "Get server instance statistics",
		Long: `Get global server instance statistics: the number of instances per server status
and per site.

Examples:
  metalcloud-cli server-instance statistics
  metalcloud-cli inst stats -f json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_READ},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceStatistics(cmd.Context())
		},
	}

	serverInstancePowerStatusBatchCmd = &cobra.Command{
		Use:     "power-status-batch infrastructure_id_or_label server_instance_id...",
		Aliases: []string{"power-get-batch"},
		Short:   "Get the power status of several server instances at once",
		Long: `Get the power status of several server instances of one infrastructure in a
single call.

The instance IDs are sent in the request body; the infrastructure is resolved by ID
or by label.

Required Arguments:
  infrastructure_id_or_label  The ID or label of the infrastructure
  server_instance_id...       One or more numeric server instance IDs

Examples:
  metalcloud-cli server-instance power-status-batch 1234 5678 5679
  metalcloud-cli inst power-get-batch prod-env 5678`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_READ},
		Args:         cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstancePowerStatusBatch(cmd.Context(), args[0], args[1:])
		},
	}

	serverInstancePowerSetBatchCmd = &cobra.Command{
		Use:     "power-set-batch infrastructure_id_or_label <on|off|reset|soft> server_instance_id...",
		Aliases: []string{"power-batch"},
		Short:   "Set the power state of several server instances at once",
		Long: `Set the power state of several server instances of one infrastructure in a
single call.

The instance IDs and the power command are sent in the request body; the
infrastructure is resolved by ID or by label.

Valid power actions:
  on     - Power on the server instances
  off    - Power off the server instances
  reset  - Hard reset the server instances
  soft   - Graceful shutdown of the server instances

Required Arguments:
  infrastructure_id_or_label  The ID or label of the infrastructure
  action                      Power action to perform (on, off, reset, soft)
  server_instance_id...       One or more numeric server instance IDs

Examples:
  metalcloud-cli server-instance power-set-batch 1234 off 5678 5679
  metalcloud-cli inst power-batch prod-env on 5678`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_SERVER_INSTANCES_WRITE},
		Args:         cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstancePowerSetBatch(cmd.Context(), args[0], args[1], args[2:])
		},
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

	serverInstanceGroupCmd.AddCommand(serverInstanceGroupUpdateMetaCmd)
	serverInstanceGroupUpdateMetaCmd.Flags().StringVar(&serverInstanceGroupFlags.configSource, "config-source", "", "Source of the metadata. Can be 'pipe' or path to a JSON/YAML file.")
	serverInstanceGroupUpdateMetaCmd.MarkFlagRequired("config-source")

	serverInstanceGroupCmd.AddCommand(serverInstanceGroupDeleteCmd)

	serverInstanceGroupCmd.AddCommand(serverInstanceGroupInstancesCmd)

	serverInstanceGroupCmd.AddCommand(serverInstanceGroupDriveGroupsCmd)

	serverInstanceGroupCmd.AddCommand(serverInstanceGroupInterfacesCmd)
	serverInstanceGroupInterfacesCmd.Flags().StringSliceVar(&serverInstanceGroupFlags.filterInfrastructureId, "filter-infrastructure-id", nil, "Filter by infrastructure ID.")
	serverInstanceGroupInterfacesCmd.Flags().StringSliceVar(&serverInstanceGroupFlags.filterServiceStatus, "filter-service-status", nil, "Filter by service status.")
	serverInstanceGroupInterfacesCmd.Flags().StringSliceVar(&serverInstanceGroupFlags.filterConfigDeployStatus, "filter-config-deploy-status", nil, "Filter by the deploy status of the pending configuration.")
	serverInstanceGroupInterfacesCmd.Flags().StringSliceVar(&serverInstanceGroupFlags.filterConfigDeployType, "filter-config-deploy-type", nil, "Filter by the deploy type of the pending configuration.")

	serverInstanceGroupCmd.AddCommand(serverInstanceGroupInterfaceGetCmd)

	// Server Instance Group network commands.
	serverInstanceGroupCmd.AddCommand(serverInstanceGroupNetworkCmd)

	serverInstanceGroupNetworkCmd.AddCommand(serverInstanceGroupNetworkConfigCmd)

	serverInstanceGroupNetworkCmd.AddCommand(serverInstanceGroupNetworkReplaceCmd)
	serverInstanceGroupNetworkReplaceCmd.Flags().StringVar(&serverInstanceGroupFlags.configSource, "config-source", "", "Source of the network configuration. Can be 'pipe' or path to a JSON/YAML file.")

	serverInstanceGroupNetworkCmd.AddCommand(serverInstanceGroupNetworkListCmd)

	serverInstanceGroupNetworkCmd.AddCommand(serverInstanceGroupNetworkGetCmd)

	serverInstanceGroupNetworkCmd.AddCommand(serverInstanceGroupNetworkConnectCmd)

	serverInstanceGroupNetworkCmd.AddCommand(serverInstanceGroupNetworkUpdateCmd)
	serverInstanceGroupNetworkUpdateCmd.Flags().StringVar(&serverInstanceGroupFlags.accessMode, "access-mode", "", "Network connection access mode.")
	serverInstanceGroupNetworkUpdateCmd.Flags().StringVar(&serverInstanceGroupFlags.tagged, "tagged", "", "Network connection tagged.")
	serverInstanceGroupNetworkUpdateCmd.Flags().StringVar(&serverInstanceGroupFlags.redundancy, "redundancy", "", "Network connection redundancy.")
	serverInstanceGroupNetworkUpdateCmd.MarkFlagsOneRequired("access-mode", "tagged", "redundancy")

	serverInstanceGroupNetworkCmd.AddCommand(serverInstanceGroupNetworkDisconnectCmd)

	// Server Instance Group network security rule commands.
	serverInstanceGroupNetworkCmd.AddCommand(serverInstanceGroupACLCmd)

	serverInstanceGroupACLCmd.AddCommand(serverInstanceGroupACLListCmd)
	serverInstanceGroupACLCmd.AddCommand(serverInstanceGroupACLGetCmd)
	serverInstanceGroupACLCmd.AddCommand(serverInstanceGroupACLConfigExampleCmd)

	serverInstanceGroupACLCmd.AddCommand(serverInstanceGroupACLAddCmd)
	serverInstanceGroupACLAddCmd.Flags().StringVar(&serverInstanceGroupFlags.configSource, "config-source", "", "Source of the new rule. Can be 'pipe' or path to a JSON/YAML file.")
	serverInstanceGroupACLAddCmd.Flags().StringVar(&serverInstanceGroupFlags.aclRuleType, "rule-type", "", "Rule type: ipv4, ipv6 or mac.")
	serverInstanceGroupACLAddCmd.Flags().StringVar(&serverInstanceGroupFlags.aclDirection, "direction", string(sdk.ACLDIRECTION_IN), "Rule direction: in or out.")
	serverInstanceGroupACLAddCmd.Flags().Int32Var(&serverInstanceGroupFlags.aclSequence, "sequence", 0, "Evaluation order of the rule.")
	serverInstanceGroupACLAddCmd.Flags().StringVar(&serverInstanceGroupFlags.aclForwardingAction, "forwarding-action", string(sdk.ACLFORWARDINGACTION_ALLOW), "Forwarding action: allow, deny, transit or discard.")
	serverInstanceGroupACLAddCmd.Flags().StringVar(&serverInstanceGroupFlags.aclEnforcementPoint, "enforcement-point", string(sdk.ACLENFORCEMENTPOINT_SVI), "Enforcement point of the rule.")
	serverInstanceGroupACLAddCmd.Flags().StringVar(&serverInstanceGroupFlags.aclNetworkProtocol, "network-protocol", "", "Network protocol of the rule.")
	serverInstanceGroupACLAddCmd.Flags().StringVar(&serverInstanceGroupFlags.aclSourceAddress, "source-address", "", "Source address of the rule.")
	serverInstanceGroupACLAddCmd.Flags().StringVar(&serverInstanceGroupFlags.aclDestinationAddress, "destination-address", "", "Destination address of the rule.")
	serverInstanceGroupACLAddCmd.Flags().StringVar(&serverInstanceGroupFlags.aclSourcePort, "source-port", "", "Source port of the rule.")
	serverInstanceGroupACLAddCmd.Flags().StringVar(&serverInstanceGroupFlags.aclDestinationPort, "destination-port", "", "Destination port of the rule.")
	serverInstanceGroupACLAddCmd.MarkFlagsOneRequired("config-source", "rule-type")
	serverInstanceGroupACLAddCmd.MarkFlagsMutuallyExclusive("config-source", "rule-type")

	serverInstanceGroupACLCmd.AddCommand(serverInstanceGroupACLUpdateCmd)
	serverInstanceGroupACLUpdateCmd.Flags().StringVar(&serverInstanceGroupFlags.configSource, "config-source", "", "Source of the updated rule. Can be 'pipe' or path to a JSON/YAML file.")
	serverInstanceGroupACLUpdateCmd.MarkFlagRequired("config-source")

	serverInstanceGroupACLCmd.AddCommand(serverInstanceGroupACLRemoveCmd)

	// Server Instance management commands.
	rootCmd.AddCommand(serverInstanceCmd)

	serverInstanceCmd.AddCommand(serverInstanceGetCmd)

	serverInstanceCmd.AddCommand(serverInstanceListCmd)
	serverInstanceListCmd.Flags().StringSliceVar(&serverInstanceFlags.filterInfrastructureId, "filter-infrastructure-id", nil, "Filter by infrastructure ID.")
	serverInstanceListCmd.Flags().StringSliceVar(&serverInstanceFlags.filterGroupId, "filter-group-id", nil, "Filter by server instance group ID.")
	serverInstanceListCmd.Flags().StringSliceVar(&serverInstanceFlags.filterServerId, "filter-server-id", nil, "Filter by server ID.")
	serverInstanceListCmd.Flags().StringSliceVar(&serverInstanceFlags.filterServiceStatus, "filter-service-status", nil, "Filter by service status.")
	serverInstanceListCmd.Flags().StringSliceVar(&serverInstanceFlags.filterConfigServerId, "filter-config-server-id", nil, "Filter by the server ID of the pending configuration.")
	serverInstanceListCmd.Flags().StringSliceVar(&serverInstanceFlags.filterConfigDeployStatus, "filter-config-deploy-status", nil, "Filter by the deploy status of the pending configuration.")
	serverInstanceListCmd.Flags().StringSliceVar(&serverInstanceFlags.filterConfigDeployType, "filter-config-deploy-type", nil, "Filter by the deploy type of the pending configuration.")

	serverInstanceCmd.AddCommand(serverInstanceConfigExampleCmd)

	serverInstanceCmd.AddCommand(serverInstanceCreateCmd)
	serverInstanceCreateCmd.Flags().StringVar(&serverInstanceFlags.configSource, "config-source", "", "Source of the new server instance configuration. Can be 'pipe' or path to a JSON/YAML file.")
	serverInstanceCreateCmd.Flags().StringVar(&serverInstanceFlags.label, "label", "", "Label of the new server instance.")
	serverInstanceCreateCmd.Flags().Int64Var(&serverInstanceFlags.groupId, "group-id", 0, "ID of the server instance group the instance belongs to.")
	serverInstanceCreateCmd.Flags().Int64Var(&serverInstanceFlags.serverTypeId, "server-type-id", 0, "ID of the server type to allocate.")
	serverInstanceCreateCmd.Flags().StringVar(&serverInstanceFlags.hostname, "hostname", "", "Custom hostname (subdomain) of the instance.")
	serverInstanceCreateCmd.Flags().Int64Var(&serverInstanceFlags.osTemplateId, "os-template-id", 0, "ID of the OS template to deploy.")
	serverInstanceCreateCmd.Flags().StringSliceVar(&serverInstanceFlags.tags, "tags", nil, "Tags of the new server instance.")
	serverInstanceCreateCmd.MarkFlagsOneRequired("config-source", "label")
	serverInstanceCreateCmd.MarkFlagsMutuallyExclusive("config-source", "label")
	serverInstanceCreateCmd.MarkFlagsMutuallyExclusive("config-source", "group-id")
	serverInstanceCreateCmd.MarkFlagsMutuallyExclusive("config-source", "server-type-id")
	serverInstanceCreateCmd.MarkFlagsMutuallyExclusive("config-source", "hostname")
	serverInstanceCreateCmd.MarkFlagsMutuallyExclusive("config-source", "os-template-id")
	serverInstanceCreateCmd.MarkFlagsMutuallyExclusive("config-source", "tags")

	serverInstanceCmd.AddCommand(serverInstanceDeleteCmd)

	serverInstanceCmd.AddCommand(serverInstanceConfigCmd)

	serverInstanceCmd.AddCommand(serverInstanceUpdateConfigCmd)
	serverInstanceUpdateConfigCmd.Flags().StringVar(&serverInstanceFlags.configSource, "config-source", "", "Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.")
	serverInstanceUpdateConfigCmd.MarkFlagRequired("config-source")

	serverInstanceCmd.AddCommand(serverInstanceUpdateMetaCmd)
	serverInstanceUpdateMetaCmd.Flags().StringVar(&serverInstanceFlags.configSource, "config-source", "", "Source of the metadata. Can be 'pipe' or path to a JSON/YAML file.")
	serverInstanceUpdateMetaCmd.MarkFlagRequired("config-source")

	serverInstanceCmd.AddCommand(serverInstancePowerCmd)
	serverInstancePowerCmd.AddCommand(serverInstancePowerStatusCmd)

	serverInstanceCmd.AddCommand(serverInstancePowerStatusBatchCmd)
	serverInstanceCmd.AddCommand(serverInstancePowerSetBatchCmd)

	serverInstanceCmd.AddCommand(serverInstanceResetCmd)
	serverInstanceCmd.AddCommand(serverInstanceReinstallOSCmd)
	serverInstanceCmd.AddCommand(serverInstanceCredentialsCmd)

	serverInstanceCmd.AddCommand(serverInstanceDrivesCmd)

	serverInstanceCmd.AddCommand(serverInstanceInterfacesCmd)
	serverInstanceInterfacesCmd.Flags().StringSliceVar(&serverInstanceFlags.ifaceFilterInfrastructureId, "filter-infrastructure-id", nil, "Filter by infrastructure ID.")
	serverInstanceInterfacesCmd.Flags().StringSliceVar(&serverInstanceFlags.ifaceFilterServiceStatus, "filter-service-status", nil, "Filter by service status.")
	serverInstanceInterfacesCmd.Flags().StringSliceVar(&serverInstanceFlags.ifaceFilterConfigDeployStatus, "filter-config-deploy-status", nil, "Filter by the deploy status of the pending configuration.")
	serverInstanceInterfacesCmd.Flags().StringSliceVar(&serverInstanceFlags.ifaceFilterConfigDeployType, "filter-config-deploy-type", nil, "Filter by the deploy type of the pending configuration.")

	serverInstanceCmd.AddCommand(serverInstanceInterfaceGetCmd)

	serverInstanceCmd.AddCommand(serverInstanceUpdateInterfaceConfigCmd)
	serverInstanceUpdateInterfaceConfigCmd.Flags().StringVar(&serverInstanceFlags.configSource, "config-source", "", "Source of the updated interface configuration. Can be 'pipe' or path to a JSON/YAML file.")
	serverInstanceUpdateInterfaceConfigCmd.MarkFlagRequired("config-source")

	serverInstanceCmd.AddCommand(serverInstanceOSInstallationDataCmd)
	serverInstanceOSInstallationDataCmd.Flags().StringVar(&serverInstanceFlags.usage, "usage", "", "Restrict the returned values to one usage type.")

	serverInstanceCmd.AddCommand(serverInstanceVariablesCmd)
	serverInstanceVariablesCmd.Flags().StringVar(&serverInstanceFlags.usage, "usage", "", "Restrict the returned values to one usage type.")

	serverInstanceCmd.AddCommand(serverInstanceStatisticsCmd)
}
