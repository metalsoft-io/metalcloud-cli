package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/vm_instance"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

// VM Instance Group management commands.
var (
	vmInstanceGroupFlags = struct {
		label                 string
		instanceCount         string
		diskSizeGB            string
		customVariablesSource string
		configSource          string
	}{}

	vmInstanceGroupCmd = &cobra.Command{
		Use:     "vm-instance-group [command]",
		Aliases: []string{"vmg", "vm-group"},
		Short:   "Manage VM instance groups within infrastructures",
		Long: `Manage VM instance groups within infrastructures.

VM instance groups allow you to create and manage collections of virtual machines
with similar configurations. This provides easier scaling and management of
multiple VM instances that serve the same purpose.

Available commands:
  list         List all VM instance groups in an infrastructure
  get          Get details of a specific VM instance group
  create       Create a new VM instance group
  update       Update VM instance group configuration
  delete       Delete a VM instance group
  instances    List instances within a VM instance group
  config       Show the pending configuration of a VM instance group
  update-meta  Update the metadata (tags) of a VM instance group
  apply-type   Apply a VM type on a VM instance group
  interfaces   List the network interfaces of a VM instance group
  interface    Get one network interface of a VM instance group
  network      Manage the network connections of a VM instance group

Examples:
  metalcloud-cli vm-instance-group list my-infra
  metalcloud-cli vmg get 12345 67890
  metalcloud-cli vm-group create 12345 5 100 3 7`,
	}

	vmInstanceGroupListCmd = &cobra.Command{
		Use:     "list infrastructure_id_or_label",
		Aliases: []string{"ls"},
		Short:   "List all VM instance groups in an infrastructure",
		Long: `List all VM instance groups in an infrastructure.

This command retrieves and displays all VM instance groups that exist within
the specified infrastructure. The output includes group details such as ID,
label, instance count, VM type, and current status.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure

Examples:
  # List all VM instance groups in infrastructure 12345
  metalcloud-cli vm-instance-group list 12345

  # List groups by infrastructure label
  metalcloud-cli vmg ls my-infra`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGroupList(cmd.Context(), args[0])
		},
	}

	vmInstanceGroupGetCmd = &cobra.Command{
		Use:     "get infrastructure_id_or_label vm_instance_group_id",
		Aliases: []string{"show"},
		Short:   "Get details of a specific VM instance group",
		Long: `Get detailed information about a specific VM instance group.

This command retrieves comprehensive information about a VM instance group
including its configuration, current status, instances, and associated metadata.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group

Examples:
  # Get details of VM instance group 67890 in infrastructure 12345
  metalcloud-cli vm-instance-group get 12345 67890

  # Get group details using alias
  metalcloud-cli vmg show my-infra 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGroupGet(cmd.Context(), args[0], args[1])
		},
	}

	vmInstanceGroupCreateCmd = &cobra.Command{
		Use:     "create infrastructure_id_or_label vm_type_id disk_size_gb instance_count [os_template_id]",
		Aliases: []string{"new"},
		Short:   "Create a new VM instance group in an infrastructure",
		Long: `Create a new VM instance group in an infrastructure.

This command creates a new VM instance group with the specified configuration.
The group will contain multiple VM instances of the same type and configuration,
making it easier to manage and scale similar workloads.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_type_id                  The VM type ID defining CPU, memory and other specs
  disk_size_gb                The disk size in GB of each VM instance of the group
  instance_count              The number of VM instances to create in the group
  os_template_id              Optional. The OS template ID for the VM instances

Examples:
  # Create a VM instance group with 3 instances using a specific OS template
  metalcloud-cli vm-instance-group create 12345 5 100 3 7

  # Create a VM instance group without specifying an OS template
  metalcloud-cli vmg new my-infra 5 50 2`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_WRITE},
		Args:         cobra.RangeArgs(4, 5),
		RunE: func(cmd *cobra.Command, args []string) error {
			osTemplateId := ""
			if len(args) == 5 {
				osTemplateId = args[4]
			}
			return vm_instance.VMInstanceGroupCreate(cmd.Context(), args[0], args[1], args[2], args[3], osTemplateId)
		},
	}

	vmInstanceGroupUpdateCmd = &cobra.Command{
		Use:     "update infrastructure_id_or_label vm_instance_group_id",
		Aliases: []string{"edit"},
		Short:   "Update VM instance group configuration",
		Long: `Update VM instance group configuration.

This command allows you to modify the configuration of an existing VM instance
group. You can update the label or custom variables associated with the group.
The current configuration revision is fetched automatically and sent as the
If-Match header, so concurrent modifications are rejected by the API.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group

Optional Flags:
  --label string                    Set or update the VM instance group label
  --custom-variables-source string  Source of the custom variables.
                                    Can be 'pipe' or path to a JSON file.

Flag Dependencies:
  At least one of --label or --custom-variables-source must be provided

Examples:
  # Update the label of a VM instance group
  metalcloud-cli vm-instance-group update 12345 67890 --label "Web Servers"

  # Update custom variables from a JSON file
  metalcloud-cli vmg edit my-infra 67890 --custom-variables-source /path/to/vars.json

  # Update custom variables from stdin
  echo '{"env": "production"}' | metalcloud-cli vm-group update 12345 67890 --custom-variables-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			customVariables, err := utils.ReadConfigFromPipeOrFile(vmInstanceGroupFlags.customVariablesSource)
			if err != nil {
				return err
			}

			return vm_instance.VMInstanceGroupUpdate(cmd.Context(), args[0], args[1],
				vmInstanceGroupFlags.label, customVariables)
		},
	}

	vmInstanceGroupDeleteCmd = &cobra.Command{
		Use:     "delete infrastructure_id_or_label vm_instance_group_id",
		Aliases: []string{"rm"},
		Short:   "Delete a VM instance group",
		Long: `Delete a VM instance group from an infrastructure.

This command permanently removes a VM instance group and all its instances
from the specified infrastructure. This action is irreversible and will
terminate all running instances within the group.

WARNING: This operation cannot be undone. All data on the VM instances
will be permanently lost unless backed up elsewhere.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group

Examples:
  # Delete VM instance group 67890 from infrastructure 12345
  metalcloud-cli vm-instance-group delete 12345 67890

  # Delete group using alias
  metalcloud-cli vmg rm my-infra 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGroupDelete(cmd.Context(), args[0], args[1])
		},
	}

	vmInstanceGroupInstancesCmd = &cobra.Command{
		Use:     "instances infrastructure_id_or_label vm_instance_group_id",
		Aliases: []string{"instances-list", "instances-ls"},
		Short:   "List VM instances within a VM instance group",
		Long: `List all VM instances within a specific VM instance group.

This command displays all individual VM instances that belong to the specified
VM instance group. The output includes instance details such as ID, status,
type and disk size for each instance in the group.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group

Examples:
  # List all instances in VM instance group 67890 from infrastructure 12345
  metalcloud-cli vm-instance-group instances 12345 67890

  # List instances using alias
  metalcloud-cli vmg instances-ls my-infra 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGroupInstances(cmd.Context(), args[0], args[1])
		},
	}

	vmInstanceGroupConfigCmd = &cobra.Command{
		Use:     "config infrastructure_id_or_label vm_instance_group_id",
		Aliases: []string{"get-config"},
		Short:   "Show the pending configuration of a VM instance group",
		Long: `Show the pending configuration of a VM instance group.

The configuration holds the changes that are applied when the infrastructure is
deployed, together with the revision that guards concurrent updates.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group

Examples:
  # Show the configuration of group 67890
  metalcloud-cli vm-instance-group config my-infra 67890

  # Using the alias
  metalcloud-cli vmg get-config 12345 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGroupGetConfig(cmd.Context(), args[0], args[1])
		},
	}

	vmInstanceGroupUpdateMetaCmd = &cobra.Command{
		Use:     "update-meta infrastructure_id_or_label vm_instance_group_id",
		Aliases: []string{"edit-meta"},
		Short:   "Update the metadata of a VM instance group",
		Long: `Update the metadata (tags) of a VM instance group.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group

Required Flags:
  --config-source string  Source of the group metadata updates.
                          Can be 'pipe' or path to a JSON/YAML file.

Examples:
  # Update the tags from a file
  metalcloud-cli vm-instance-group update-meta my-infra 67890 --config-source meta.json

  # Update the tags from stdin
  echo '{"tags":["prod"]}' | metalcloud-cli vmg edit-meta 12345 67890 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(vmInstanceGroupFlags.configSource)
			if err != nil {
				return err
			}

			return vm_instance.VMInstanceGroupUpdateMeta(cmd.Context(), args[0], args[1], config)
		},
	}

	vmInstanceGroupApplyTypeCmd = &cobra.Command{
		Use:     "apply-type infrastructure_id_or_label vm_instance_group_id vm_type_id",
		Aliases: []string{"set-type"},
		Short:   "Apply a VM type on a VM instance group",
		Long: `Apply a VM type on a VM instance group.

Every instance of the group is reconfigured with the CPU, memory and GPU
specification of the given VM type. The current group revision is fetched
automatically and sent as the If-Match header.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group
  vm_type_id                  The numeric ID of the VM type to apply

Examples:
  # Apply VM type 5 on group 67890
  metalcloud-cli vm-instance-group apply-type my-infra 67890 5

  # Using the alias
  metalcloud-cli vmg set-type 12345 67890 5`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGroupApplyType(cmd.Context(), args[0], args[1], args[2])
		},
	}

	vmInstanceGroupInterfacesCmd = &cobra.Command{
		Use:     "interfaces infrastructure_id_or_label vm_instance_group_id",
		Aliases: []string{"list-interfaces"},
		Short:   "List the network interfaces of a VM instance group",
		Long: `List the network interfaces of a VM instance group.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group

Examples:
  # List the interfaces of group 67890
  metalcloud-cli vm-instance-group interfaces my-infra 67890

  # Using the alias
  metalcloud-cli vmg list-interfaces 12345 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGroupInterfaces(cmd.Context(), args[0], args[1])
		},
	}

	vmInstanceGroupInterfaceCmd = &cobra.Command{
		Use:     "interface infrastructure_id_or_label vm_instance_group_id interface_id",
		Aliases: []string{"get-interface"},
		Short:   "Get one network interface of a VM instance group",
		Long: `Get detailed information about one network interface of a VM instance group.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group
  interface_id                The numeric ID of the interface

Examples:
  # Get interface 2 of group 67890
  metalcloud-cli vm-instance-group interface my-infra 67890 2

  # Using the alias
  metalcloud-cli vmg get-interface 12345 67890 2`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGroupInterfaceGet(cmd.Context(), args[0], args[1], args[2])
		},
	}

	vmInstanceGroupNetworkFlags = struct {
		configSource     string
		logicalNetworkId string
		accessMode       string
		tagged           string
		redundancy       string
		mtu              string
	}{}

	vmInstanceGroupNetworkCmd = &cobra.Command{
		Use:     "network [command]",
		Aliases: []string{"net"},
		Short:   "Manage the network connections of a VM instance group",
		Long: `Manage the network connections of a VM instance group.

Available commands:
  list            Show the network configuration of a group
  connections     List the network connections of a group
  get             Get one network connection of a group
  connect         Connect a group to a logical network
  update          Update one network connection of a group
  disconnect      Remove one network connection of a group
  config-example  Print a network connection configuration example

Examples:
  metalcloud-cli vm-instance-group network list my-infra 67890
  metalcloud-cli vmg net connect my-infra 67890 --logical-network-id 5 --access-mode l2 --tagged true`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_READ},
	}

	vmInstanceGroupNetworkListCmd = &cobra.Command{
		Use:     "list infrastructure_id_or_label vm_instance_group_id",
		Aliases: []string{"ls"},
		Short:   "Show the network configuration of a VM instance group",
		Long: `Show the network endpoint group that carries the network connections of a
VM instance group.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group

Examples:
  # Show the network configuration of group 67890
  metalcloud-cli vm-instance-group network list my-infra 67890

  # Using the alias
  metalcloud-cli vmg net ls 12345 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGroupNetworkConfiguration(cmd.Context(), args[0], args[1])
		},
	}

	vmInstanceGroupNetworkConnectionsCmd = &cobra.Command{
		Use:     "connections infrastructure_id_or_label vm_instance_group_id",
		Aliases: []string{"list-connections"},
		Short:   "List the network connections of a VM instance group",
		Long: `List all network connections of a VM instance group.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group

Examples:
  # List the network connections of group 67890
  metalcloud-cli vm-instance-group network connections my-infra 67890

  # Using the alias
  metalcloud-cli vmg net list-connections 12345 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGroupNetworkConnections(cmd.Context(), args[0], args[1])
		},
	}

	vmInstanceGroupNetworkGetCmd = &cobra.Command{
		Use:     "get infrastructure_id_or_label vm_instance_group_id connection_id",
		Aliases: []string{"show"},
		Short:   "Get one network connection of a VM instance group",
		Long: `Get detailed information about one network connection of a VM instance group.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group
  connection_id               The numeric ID of the network connection

Examples:
  # Get network connection 535 of group 67890
  metalcloud-cli vm-instance-group network get my-infra 67890 535

  # Using the alias
  metalcloud-cli vmg net show 12345 67890 535`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGroupNetworkGet(cmd.Context(), args[0], args[1], args[2])
		},
	}

	vmInstanceGroupNetworkConfigExampleCmd = &cobra.Command{
		Use:     "config-example",
		Aliases: []string{"example"},
		Short:   "Print a network connection configuration example",
		Long: `Print a network connection configuration example.

The printed document lists the fields accepted by the connect command and can be
saved to a file, edited and passed back through the --config-source flag.

Examples:
  # Show the example configuration
  metalcloud-cli vm-instance-group network config-example

  # Save the example to a file
  metalcloud-cli vmg net config-example > connection.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGroupNetworkConfigExample(cmd.Context())
		},
	}

	vmInstanceGroupNetworkConnectCmd = &cobra.Command{
		Use:     "connect infrastructure_id_or_label vm_instance_group_id",
		Aliases: []string{"new", "add"},
		Short:   "Connect a VM instance group to a logical network",
		Long: `Connect a VM instance group to a logical network.

The connection can be described either by a complete configuration document
(--config-source) or by individual flags.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group

Required Flags:
  --config-source string       Source of the new connection configuration.
                               Can be 'pipe' or path to a JSON/YAML file.
  --logical-network-id string  The ID of the logical network to connect to.
  --access-mode string         The network access mode (e.g. 'l2').

Optional Flags:
  --tagged string      Whether VLAN tagging is enabled (true/false).
  --redundancy string  The redundancy mode ('active-backup', 'active-active').
  --mtu string         The MTU of the network connection.

Flag Dependencies:
  --config-source is mutually exclusive with the individual flags. One of
  --config-source or --logical-network-id is required; --access-mode is
  required when --logical-network-id is used.

Examples:
  # Connect group 67890 to logical network 5 from flags
  metalcloud-cli vm-instance-group network connect my-infra 67890 --logical-network-id 5 --access-mode l2 --tagged true

  # Connect from a configuration file
  metalcloud-cli vmg net add 12345 67890 --config-source connection.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if vmInstanceGroupNetworkFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(vmInstanceGroupNetworkFlags.configSource)
				if err != nil {
					return err
				}

				return vm_instance.VMInstanceGroupNetworkConnect(cmd.Context(), args[0], args[1], config)
			}

			return vm_instance.VMInstanceGroupNetworkConnectFromFlags(cmd.Context(), args[0], args[1],
				vmInstanceGroupNetworkFlags.logicalNetworkId,
				vmInstanceGroupNetworkFlags.accessMode,
				vmInstanceGroupNetworkFlags.tagged,
				vmInstanceGroupNetworkFlags.redundancy,
				vmInstanceGroupNetworkFlags.mtu)
		},
	}

	vmInstanceGroupNetworkUpdateCmd = &cobra.Command{
		Use:     "update infrastructure_id_or_label vm_instance_group_id connection_id",
		Aliases: []string{"edit"},
		Short:   "Update one network connection of a VM instance group",
		Long: `Update one network connection of a VM instance group.

The updates can be described either by a complete configuration document
(--config-source) or by individual flags.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group
  connection_id               The numeric ID of the network connection

Required Flags:
  --config-source string  Source of the connection updates.
                          Can be 'pipe' or path to a JSON/YAML file.
  --access-mode string    The network access mode (e.g. 'l2').
  --tagged string         Whether VLAN tagging is enabled (true/false).
  --redundancy string     The redundancy mode ('active-backup', 'active-active').
  --mtu string            The MTU of the network connection.

Flag Dependencies:
  --config-source is mutually exclusive with the individual flags. At least one
  of --config-source, --access-mode, --tagged, --redundancy or --mtu is required.

Examples:
  # Change the access mode of connection 535
  metalcloud-cli vm-instance-group network update my-infra 67890 535 --access-mode l2

  # Enable VLAN tagging
  metalcloud-cli vmg net edit 12345 67890 535 --tagged true`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			if vmInstanceGroupNetworkFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(vmInstanceGroupNetworkFlags.configSource)
				if err != nil {
					return err
				}

				return vm_instance.VMInstanceGroupNetworkUpdate(cmd.Context(), args[0], args[1], args[2], config)
			}

			return vm_instance.VMInstanceGroupNetworkUpdateFromFlags(cmd.Context(), args[0], args[1], args[2],
				vmInstanceGroupNetworkFlags.accessMode,
				vmInstanceGroupNetworkFlags.tagged,
				vmInstanceGroupNetworkFlags.redundancy,
				vmInstanceGroupNetworkFlags.mtu)
		},
	}

	vmInstanceGroupNetworkDisconnectCmd = &cobra.Command{
		Use:     "disconnect infrastructure_id_or_label vm_instance_group_id connection_id",
		Aliases: []string{"rm", "remove"},
		Short:   "Remove one network connection of a VM instance group",
		Long: `Remove one network connection of a VM instance group.

This disconnects every instance of the group from the logical network.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_group_id        The numeric ID of the VM instance group
  connection_id               The numeric ID of the network connection to remove

Examples:
  # Remove network connection 535 of group 67890
  metalcloud-cli vm-instance-group network disconnect my-infra 67890 535

  # Using an alias
  metalcloud-cli vmg net rm 12345 67890 535`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGroupNetworkDisconnect(cmd.Context(), args[0], args[1], args[2])
		},
	}
)

// VM Instance management commands.
var (
	vmInstanceFlags = struct {
		configSource string
		vmTypeId     string
		groupId      string
		diskSizeGB   string
		tags         []string
		usage        string
		removeEmpty  bool
	}{}

	vmInstanceCmd = &cobra.Command{
		Use:     "vm-instance [command]",
		Aliases: []string{"vmi", "vm"},
		Short:   "Manage individual VM instances within infrastructures",
		Long: `Manage individual VM instances within infrastructures.

VM instances are individual virtual machines that can be created, managed,
and controlled independently. This includes operations like getting instance
details, listing instances, managing power states, and accessing configuration.

Available commands:
  list                  List all VM instances in an infrastructure
  get                   Get details of a specific VM instance
  config-example        Print an example VM instance configuration
  create                Create a new VM instance
  delete                Delete a VM instance
  config                Show the pending configuration of a VM instance
  update-config         Update the pending configuration of a VM instance
  update-meta           Update the metadata (tags) of a VM instance
  credentials           Show the credentials of a VM instance
  variables             Show the variables of a VM instance
  os-installation-data  Show the OS installation data of a VM instance
  power-status          Get VM instance power status
  start                 Start a VM instance
  shutdown              Shutdown a VM instance
  reboot                Reboot a VM instance
  apply-type            Apply a VM type on a VM instance

Examples:
  metalcloud-cli vm-instance list my-infra
  metalcloud-cli vmi get 12345 67890
  metalcloud-cli vm start 12345 67890`,
	}

	vmInstanceGetCmd = &cobra.Command{
		Use:     "get infrastructure_id_or_label vm_instance_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific VM instance",
		Long: `Get detailed information about a specific VM instance.

This command retrieves comprehensive information about a VM instance including
its current status, configuration, network details, disk information, and
associated metadata.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_id              The numeric ID of the VM instance

Examples:
  # Get details of VM instance 67890 in infrastructure 12345
  metalcloud-cli vm-instance get 12345 67890

  # Get instance details using alias
  metalcloud-cli vmi show my-infra 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGet(cmd.Context(), args[0], args[1])
		},
	}

	vmInstanceListCmd = &cobra.Command{
		Use:     "list infrastructure_id_or_label",
		Aliases: []string{"ls"},
		Short:   "List all VM instances in an infrastructure",
		Long: `List all VM instances in an infrastructure.

This command retrieves and displays all VM instances that exist within the
specified infrastructure. The output includes instance details such as ID,
status, VM type, disk size and other relevant information for each instance.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure

Examples:
  # List all VM instances in infrastructure 12345
  metalcloud-cli vm-instance list 12345

  # List instances by infrastructure label
  metalcloud-cli vmi ls my-infra`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceList(cmd.Context(), args[0])
		},
	}

	vmInstanceConfigExampleCmd = &cobra.Command{
		Use:     "config-example",
		Aliases: []string{"example"},
		Short:   "Print a VM instance configuration example",
		Long: `Print a VM instance configuration example.

The printed document lists every field accepted by the create command and can be
saved to a file, edited and passed back through the --config-source flag.

Examples:
  # Show the example configuration
  metalcloud-cli vm-instance config-example

  # Save the example to a file
  metalcloud-cli vmi config-example > vm-instance.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceConfigExample(cmd.Context())
		},
	}

	vmInstanceCreateCmd = &cobra.Command{
		Use:     "create infrastructure_id_or_label",
		Aliases: []string{"new"},
		Short:   "Create a new VM instance",
		Long: `Create a new VM instance in an infrastructure.

The instance can be described either by a complete configuration document
(--config-source) or by individual flags. The new instance joins an existing
VM instance group and is materialised when the infrastructure is deployed.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure

Required Flags:
  --config-source string  Source of the new VM instance configuration.
                          Can be 'pipe' or path to a JSON/YAML file.
  --vm-type-id string     The VM type of the new VM instance.
  --group-id string       The VM instance group of the new VM instance.

Optional Flags:
  --disk-size-gb string  Disk size in GB of the new VM instance.
  --tags strings         Tags of the new VM instance.

Flag Dependencies:
  --config-source is mutually exclusive with the individual flags. One of
  --config-source or --vm-type-id is required; --group-id is required when
  --vm-type-id is used.

Examples:
  # Create a VM instance from flags
  metalcloud-cli vm-instance create my-infra --vm-type-id 5 --group-id 67890 --disk-size-gb 40

  # Create a VM instance from a file
  metalcloud-cli vmi new 12345 --config-source vm-instance.json

  # Create a VM instance from stdin
  echo '{"typeId":5,"groupId":67890}' | metalcloud-cli vm create 12345 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if vmInstanceFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(vmInstanceFlags.configSource)
				if err != nil {
					return err
				}

				return vm_instance.VMInstanceCreate(cmd.Context(), args[0], config)
			}

			return vm_instance.VMInstanceCreateFromFlags(cmd.Context(), args[0],
				vmInstanceFlags.vmTypeId,
				vmInstanceFlags.groupId,
				vmInstanceFlags.diskSizeGB,
				vmInstanceFlags.tags)
		},
	}

	vmInstanceDeleteCmd = &cobra.Command{
		Use:     "delete infrastructure_id_or_label vm_instance_id",
		Aliases: []string{"rm"},
		Short:   "Delete a VM instance",
		Long: `Delete a VM instance from an infrastructure.

The instance is removed when the infrastructure is deployed. The current
instance revision is fetched automatically and sent as the If-Match header.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_id              The numeric ID of the VM instance to delete

Examples:
  # Delete VM instance 67890 from infrastructure 12345
  metalcloud-cli vm-instance delete 12345 67890

  # Using the alias
  metalcloud-cli vmi rm my-infra 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceDelete(cmd.Context(), args[0], args[1])
		},
	}

	vmInstanceGetConfigCmd = &cobra.Command{
		Use:     "config infrastructure_id_or_label vm_instance_id",
		Aliases: []string{"get-config"},
		Short:   "Show the pending configuration of a VM instance",
		Long: `Show the pending configuration of a VM instance.

This command retrieves the current configuration of a VM instance including
hardware specifications, disk configuration and other system parameters,
together with the revision that guards concurrent updates.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_id              The numeric ID of the VM instance

Examples:
  # Get configuration for VM instance 67890 in infrastructure 12345
  metalcloud-cli vm-instance config 12345 67890

  # Get configuration using alias
  metalcloud-cli vmi get-config my-infra 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGetConfig(cmd.Context(), args[0], args[1])
		},
	}

	vmInstanceUpdateConfigCmd = &cobra.Command{
		Use:     "update-config infrastructure_id_or_label vm_instance_id",
		Aliases: []string{"edit-config"},
		Short:   "Update the pending configuration of a VM instance",
		Long: `Update the pending configuration of a VM instance.

The current configuration revision is fetched automatically and sent as the
If-Match header, so concurrent modifications are rejected by the API.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_id              The numeric ID of the VM instance

Required Flags:
  --config-source string  Source of the VM instance configuration updates.
                          Can be 'pipe' or path to a JSON/YAML file.

Examples:
  # Update the configuration from a file
  metalcloud-cli vm-instance update-config my-infra 67890 --config-source updates.json

  # Update the label from stdin
  echo '{"label":"web-1"}' | metalcloud-cli vmi edit-config 12345 67890 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(vmInstanceFlags.configSource)
			if err != nil {
				return err
			}

			return vm_instance.VMInstanceUpdateConfig(cmd.Context(), args[0], args[1], config)
		},
	}

	vmInstanceUpdateMetaCmd = &cobra.Command{
		Use:     "update-meta infrastructure_id_or_label vm_instance_id",
		Aliases: []string{"edit-meta"},
		Short:   "Update the metadata of a VM instance",
		Long: `Update the metadata (tags) of a VM instance.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_id              The numeric ID of the VM instance

Required Flags:
  --config-source string  Source of the VM instance metadata updates.
                          Can be 'pipe' or path to a JSON/YAML file.

Examples:
  # Update the tags from a file
  metalcloud-cli vm-instance update-meta my-infra 67890 --config-source meta.json

  # Update the tags from stdin
  echo '{"tags":["prod"]}' | metalcloud-cli vmi edit-meta 12345 67890 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(vmInstanceFlags.configSource)
			if err != nil {
				return err
			}

			return vm_instance.VMInstanceUpdateMeta(cmd.Context(), args[0], args[1], config)
		},
	}

	vmInstanceApplyTypeCmd = &cobra.Command{
		Use:     "apply-type infrastructure_id_or_label vm_instance_id vm_type_id",
		Aliases: []string{"set-type"},
		Short:   "Apply a VM type on a VM instance",
		Long: `Apply a VM type on a VM instance.

The instance is reconfigured with the CPU, memory and GPU specification of the
given VM type. The current instance revision is fetched automatically and sent
as the If-Match header.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_id              The numeric ID of the VM instance
  vm_type_id                  The numeric ID of the VM type to apply

Examples:
  # Apply VM type 5 on VM instance 67890
  metalcloud-cli vm-instance apply-type my-infra 67890 5

  # Using the alias
  metalcloud-cli vmi set-type 12345 67890 5`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceApplyType(cmd.Context(), args[0], args[1], args[2])
		},
	}

	vmInstanceVariablesCmd = &cobra.Command{
		Use:     "variables infrastructure_id_or_label vm_instance_id",
		Aliases: []string{"vars"},
		Short:   "Show the variables of a VM instance",
		Long: `Show the variables a VM instance exposes to extensions and OS templates.

The response is a deeply nested document, so every non native output format is
rendered as YAML.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_id              The numeric ID of the VM instance

Optional Flags:
  --usage string  Restrict the variables to one usage type.

Examples:
  # Show the variables of VM instance 67890
  metalcloud-cli vm-instance variables my-infra 67890

  # Show only the variables used by Ansible bundles
  metalcloud-cli vmi vars 12345 67890 --usage AnsibleBundle`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceVariables(cmd.Context(), args[0], args[1], vmInstanceFlags.usage)
		},
	}

	vmInstanceOsInstallationDataCmd = &cobra.Command{
		Use:     "os-installation-data infrastructure_id_or_label vm_instance_id",
		Aliases: []string{"os-data"},
		Short:   "Show the OS installation data of a VM instance",
		Long: `Show the values the OS installer of a VM instance is rendered with.

The response is a deeply nested document, so every non native output format is
rendered as YAML.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_id              The numeric ID of the VM instance

Optional Flags:
  --usage string   Restrict the OS installation data to one usage type.
  --remove-empty   Omit the empty entries from the response.

Examples:
  # Show the OS installation data of VM instance 67890
  metalcloud-cli vm-instance os-installation-data my-infra 67890

  # Skip the empty entries
  metalcloud-cli vmi os-data 12345 67890 --remove-empty`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceOSInstallationData(cmd.Context(), args[0], args[1],
				vmInstanceFlags.usage, vmInstanceFlags.removeEmpty)
		},
	}

	vmInstanceStartCmd = &cobra.Command{
		Use:   "start infrastructure_id_or_label vm_instance_id",
		Short: "Start a VM instance",
		Long: `Start a VM instance.

This command initiates the startup process for a VM instance that is currently
powered off or stopped. The instance will be powered on and begin booting
according to its configured operating system and startup settings.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_id              The numeric ID of the VM instance to start

Examples:
  # Start VM instance 67890 in infrastructure 12345
  metalcloud-cli vm-instance start 12345 67890

  # Start instance using alias
  metalcloud-cli vm start my-infra 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstancePowerControl(cmd.Context(), args[0], args[1], "start")
		},
	}

	vmInstanceShutdownCmd = &cobra.Command{
		Use:   "shutdown infrastructure_id_or_label vm_instance_id",
		Short: "Shutdown a VM instance",
		Long: `Shutdown a VM instance gracefully.

This command initiates a graceful shutdown process for a running VM instance.
The instance will receive a shutdown signal and will attempt to properly
terminate all running processes before powering off. This is the recommended
way to stop a VM instance to prevent data loss.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_id              The numeric ID of the VM instance to shutdown

Examples:
  # Shutdown VM instance 67890 in infrastructure 12345
  metalcloud-cli vm-instance shutdown 12345 67890

  # Shutdown instance using alias
  metalcloud-cli vm shutdown my-infra 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstancePowerControl(cmd.Context(), args[0], args[1], "shutdown")
		},
	}

	vmInstanceRebootCmd = &cobra.Command{
		Use:   "reboot infrastructure_id_or_label vm_instance_id",
		Short: "Reboot a VM instance",
		Long: `Reboot a VM instance.

This command initiates a restart process for a running VM instance. The instance
will be gracefully shutdown and then automatically restarted. This is useful
for applying configuration changes or recovering from software issues.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_id              The numeric ID of the VM instance to reboot

Examples:
  # Reboot VM instance 67890 in infrastructure 12345
  metalcloud-cli vm-instance reboot 12345 67890

  # Reboot instance using alias
  metalcloud-cli vm reboot my-infra 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstancePowerControl(cmd.Context(), args[0], args[1], "reboot")
		},
	}

	vmInstanceCredentialsCmd = &cobra.Command{
		Use:     "credentials infrastructure_id_or_label vm_instance_id",
		Aliases: []string{"creds"},
		Short:   "Get login credentials for a VM instance",
		Long: `Get the login credentials of a VM instance.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_id              The numeric ID of the VM instance

Examples:
  # Show the credentials of VM instance 67890
  metalcloud-cli vm-instance credentials 12345 67890

  # Using the alias
  metalcloud-cli vmi creds my-infra 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGetCredentials(cmd.Context(), args[0], args[1])
		},
	}

	vmInstancePowerStatusCmd = &cobra.Command{
		Use:   "power-status infrastructure_id_or_label vm_instance_id",
		Short: "Get VM instance power status",
		Long: `Get VM instance power status.

This command retrieves the current power state of a VM instance, indicating
whether it is running, stopped, starting, stopping, or in another power state.

Required Arguments:
  infrastructure_id_or_label  The ID or the label of the infrastructure
  vm_instance_id              The numeric ID of the VM instance

Examples:
  # Get power status of VM instance 67890 in infrastructure 12345
  metalcloud-cli vm-instance power-status 12345 67890

  # Get power status using alias
  metalcloud-cli vm power-status my-infra 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGetPowerStatus(cmd.Context(), args[0], args[1])
		},
	}
)

func init() {
	// VM Instance Group management commands
	rootCmd.AddCommand(vmInstanceGroupCmd)

	vmInstanceGroupCmd.AddCommand(vmInstanceGroupListCmd)

	vmInstanceGroupCmd.AddCommand(vmInstanceGroupGetCmd)

	vmInstanceGroupCmd.AddCommand(vmInstanceGroupCreateCmd)

	vmInstanceGroupCmd.AddCommand(vmInstanceGroupUpdateCmd)
	vmInstanceGroupUpdateCmd.Flags().StringVar(&vmInstanceGroupFlags.label, "label", "", "Set the VM instance group label.")
	vmInstanceGroupUpdateCmd.Flags().StringVar(&vmInstanceGroupFlags.customVariablesSource, "custom-variables-source", "", "Source of the custom variables. Can be 'pipe' or path to a JSON file.")
	vmInstanceGroupUpdateCmd.MarkFlagsOneRequired("label", "custom-variables-source")

	vmInstanceGroupCmd.AddCommand(vmInstanceGroupDeleteCmd)

	vmInstanceGroupCmd.AddCommand(vmInstanceGroupInstancesCmd)

	vmInstanceGroupCmd.AddCommand(vmInstanceGroupConfigCmd)

	vmInstanceGroupCmd.AddCommand(vmInstanceGroupUpdateMetaCmd)
	vmInstanceGroupUpdateMetaCmd.Flags().StringVar(&vmInstanceGroupFlags.configSource, "config-source", "", "Source of the VM instance group metadata updates. Can be 'pipe' or path to a JSON/YAML file.")
	vmInstanceGroupUpdateMetaCmd.MarkFlagsOneRequired("config-source")

	vmInstanceGroupCmd.AddCommand(vmInstanceGroupApplyTypeCmd)

	vmInstanceGroupCmd.AddCommand(vmInstanceGroupInterfacesCmd)

	vmInstanceGroupCmd.AddCommand(vmInstanceGroupInterfaceCmd)

	vmInstanceGroupCmd.AddCommand(vmInstanceGroupNetworkCmd)

	vmInstanceGroupNetworkCmd.AddCommand(vmInstanceGroupNetworkListCmd)

	vmInstanceGroupNetworkCmd.AddCommand(vmInstanceGroupNetworkConnectionsCmd)

	vmInstanceGroupNetworkCmd.AddCommand(vmInstanceGroupNetworkGetCmd)

	vmInstanceGroupNetworkCmd.AddCommand(vmInstanceGroupNetworkConfigExampleCmd)

	vmInstanceGroupNetworkCmd.AddCommand(vmInstanceGroupNetworkConnectCmd)
	vmInstanceGroupNetworkConnectCmd.Flags().StringVar(&vmInstanceGroupNetworkFlags.configSource, "config-source", "", "Source of the new network connection configuration. Can be 'pipe' or path to a JSON/YAML file.")
	vmInstanceGroupNetworkConnectCmd.Flags().StringVar(&vmInstanceGroupNetworkFlags.logicalNetworkId, "logical-network-id", "", "The ID of the logical network to connect to.")
	vmInstanceGroupNetworkConnectCmd.Flags().StringVar(&vmInstanceGroupNetworkFlags.accessMode, "access-mode", "", "Network connection access mode.")
	vmInstanceGroupNetworkConnectCmd.Flags().StringVar(&vmInstanceGroupNetworkFlags.tagged, "tagged", "", "Network connection VLAN tagging (true/false).")
	vmInstanceGroupNetworkConnectCmd.Flags().StringVar(&vmInstanceGroupNetworkFlags.redundancy, "redundancy", "", "Network connection redundancy mode.")
	vmInstanceGroupNetworkConnectCmd.Flags().StringVar(&vmInstanceGroupNetworkFlags.mtu, "mtu", "", "Network connection MTU.")
	vmInstanceGroupNetworkConnectCmd.MarkFlagsOneRequired("config-source", "logical-network-id")
	vmInstanceGroupNetworkConnectCmd.MarkFlagsMutuallyExclusive("config-source", "logical-network-id")
	vmInstanceGroupNetworkConnectCmd.MarkFlagsMutuallyExclusive("config-source", "access-mode")
	vmInstanceGroupNetworkConnectCmd.MarkFlagsMutuallyExclusive("config-source", "tagged")
	vmInstanceGroupNetworkConnectCmd.MarkFlagsMutuallyExclusive("config-source", "redundancy")
	vmInstanceGroupNetworkConnectCmd.MarkFlagsMutuallyExclusive("config-source", "mtu")
	vmInstanceGroupNetworkConnectCmd.MarkFlagsRequiredTogether("logical-network-id", "access-mode")

	vmInstanceGroupNetworkCmd.AddCommand(vmInstanceGroupNetworkUpdateCmd)
	vmInstanceGroupNetworkUpdateCmd.Flags().StringVar(&vmInstanceGroupNetworkFlags.configSource, "config-source", "", "Source of the network connection updates. Can be 'pipe' or path to a JSON/YAML file.")
	vmInstanceGroupNetworkUpdateCmd.Flags().StringVar(&vmInstanceGroupNetworkFlags.accessMode, "access-mode", "", "Network connection access mode.")
	vmInstanceGroupNetworkUpdateCmd.Flags().StringVar(&vmInstanceGroupNetworkFlags.tagged, "tagged", "", "Network connection VLAN tagging (true/false).")
	vmInstanceGroupNetworkUpdateCmd.Flags().StringVar(&vmInstanceGroupNetworkFlags.redundancy, "redundancy", "", "Network connection redundancy mode.")
	vmInstanceGroupNetworkUpdateCmd.Flags().StringVar(&vmInstanceGroupNetworkFlags.mtu, "mtu", "", "Network connection MTU.")
	vmInstanceGroupNetworkUpdateCmd.MarkFlagsOneRequired("config-source", "access-mode", "tagged", "redundancy", "mtu")
	vmInstanceGroupNetworkUpdateCmd.MarkFlagsMutuallyExclusive("config-source", "access-mode")
	vmInstanceGroupNetworkUpdateCmd.MarkFlagsMutuallyExclusive("config-source", "tagged")
	vmInstanceGroupNetworkUpdateCmd.MarkFlagsMutuallyExclusive("config-source", "redundancy")
	vmInstanceGroupNetworkUpdateCmd.MarkFlagsMutuallyExclusive("config-source", "mtu")

	vmInstanceGroupNetworkCmd.AddCommand(vmInstanceGroupNetworkDisconnectCmd)

	// VM Instance management commands
	rootCmd.AddCommand(vmInstanceCmd)

	vmInstanceCmd.AddCommand(vmInstanceGetCmd)

	vmInstanceCmd.AddCommand(vmInstanceListCmd)

	vmInstanceCmd.AddCommand(vmInstanceConfigExampleCmd)

	vmInstanceCmd.AddCommand(vmInstanceCreateCmd)
	vmInstanceCreateCmd.Flags().StringVar(&vmInstanceFlags.configSource, "config-source", "", "Source of the new VM instance configuration. Can be 'pipe' or path to a JSON/YAML file.")
	vmInstanceCreateCmd.Flags().StringVar(&vmInstanceFlags.vmTypeId, "vm-type-id", "", "The VM type of the new VM instance.")
	vmInstanceCreateCmd.Flags().StringVar(&vmInstanceFlags.groupId, "group-id", "", "The VM instance group of the new VM instance.")
	vmInstanceCreateCmd.Flags().StringVar(&vmInstanceFlags.diskSizeGB, "disk-size-gb", "", "Disk size in GB of the new VM instance.")
	vmInstanceCreateCmd.Flags().StringSliceVar(&vmInstanceFlags.tags, "tags", nil, "Tags of the new VM instance.")
	vmInstanceCreateCmd.MarkFlagsOneRequired("config-source", "vm-type-id")
	vmInstanceCreateCmd.MarkFlagsMutuallyExclusive("config-source", "vm-type-id")
	vmInstanceCreateCmd.MarkFlagsMutuallyExclusive("config-source", "group-id")
	vmInstanceCreateCmd.MarkFlagsMutuallyExclusive("config-source", "disk-size-gb")
	vmInstanceCreateCmd.MarkFlagsMutuallyExclusive("config-source", "tags")
	vmInstanceCreateCmd.MarkFlagsRequiredTogether("vm-type-id", "group-id")

	vmInstanceCmd.AddCommand(vmInstanceDeleteCmd)

	vmInstanceCmd.AddCommand(vmInstanceGetConfigCmd)

	vmInstanceCmd.AddCommand(vmInstanceUpdateConfigCmd)
	vmInstanceUpdateConfigCmd.Flags().StringVar(&vmInstanceFlags.configSource, "config-source", "", "Source of the VM instance configuration updates. Can be 'pipe' or path to a JSON/YAML file.")
	vmInstanceUpdateConfigCmd.MarkFlagsOneRequired("config-source")

	vmInstanceCmd.AddCommand(vmInstanceUpdateMetaCmd)
	vmInstanceUpdateMetaCmd.Flags().StringVar(&vmInstanceFlags.configSource, "config-source", "", "Source of the VM instance metadata updates. Can be 'pipe' or path to a JSON/YAML file.")
	vmInstanceUpdateMetaCmd.MarkFlagsOneRequired("config-source")

	vmInstanceCmd.AddCommand(vmInstanceApplyTypeCmd)

	vmInstanceCmd.AddCommand(vmInstanceVariablesCmd)
	vmInstanceVariablesCmd.Flags().StringVar(&vmInstanceFlags.usage, "usage", "", "Restrict the variables to one usage type (HTTPRequest, JavaScript, APICall, AnsibleBundle, SSHExec, Copy, OSAsset).")

	vmInstanceCmd.AddCommand(vmInstanceOsInstallationDataCmd)
	vmInstanceOsInstallationDataCmd.Flags().StringVar(&vmInstanceFlags.usage, "usage", "", "Restrict the OS installation data to one usage type (HTTPRequest, JavaScript, APICall, AnsibleBundle, SSHExec, Copy, OSAsset).")
	vmInstanceOsInstallationDataCmd.Flags().BoolVar(&vmInstanceFlags.removeEmpty, "remove-empty", false, "Omit the empty entries from the response.")

	vmInstanceCmd.AddCommand(vmInstanceStartCmd)

	vmInstanceCmd.AddCommand(vmInstanceShutdownCmd)

	vmInstanceCmd.AddCommand(vmInstanceRebootCmd)

	vmInstanceCmd.AddCommand(vmInstancePowerStatusCmd)

	vmInstanceCmd.AddCommand(vmInstanceCredentialsCmd)
}
