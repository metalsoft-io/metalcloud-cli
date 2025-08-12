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
  list      List all VM instance groups in an infrastructure
  get       Get details of a specific VM instance group
  create    Create a new VM instance group
  update    Update VM instance group configuration
  delete    Delete a VM instance group
  instances List instances within a VM instance group

Examples:
  metalcloud-cli vm-instance-group list 12345
  metalcloud-cli vmg get 12345 67890
  metalcloud-cli vm-group create 12345 vm-type-1 100 3 ubuntu-20.04`,
	}

	vmInstanceGroupListCmd = &cobra.Command{
		Use:     "list infrastructure_id",
		Aliases: []string{"ls"},
		Short:   "List all VM instance groups in an infrastructure",
		Long: `List all VM instance groups in an infrastructure.

This command retrieves and displays all VM instance groups that exist within
the specified infrastructure. The output includes group details such as ID,
label, instance count, VM type, and current status.

ARGUMENTS:
  infrastructure_id  The ID of the infrastructure to list groups from

EXAMPLES:
  # List all VM instance groups in infrastructure 12345
  metalcloud-cli vm-instance-group list 12345
  
  # List groups using alias
  metalcloud-cli vmg ls 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGroupList(cmd.Context(), args[0])
		},
	}

	vmInstanceGroupGetCmd = &cobra.Command{
		Use:     "get infrastructure_id vm_instance_group_id",
		Aliases: []string{"show"},
		Short:   "Get details of a specific VM instance group",
		Long: `Get detailed information about a specific VM instance group.

This command retrieves comprehensive information about a VM instance group
including its configuration, current status, instances, and associated metadata.

ARGUMENTS:
  infrastructure_id     The ID of the infrastructure containing the group
  vm_instance_group_id  The ID of the VM instance group to retrieve

EXAMPLES:
  # Get details of VM instance group 67890 in infrastructure 12345
  metalcloud-cli vm-instance-group get 12345 67890
  
  # Get group details using alias
  metalcloud-cli vmg show 12345 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGroupGet(cmd.Context(), args[0], args[1])
		},
	}

	vmInstanceGroupCreateCmd = &cobra.Command{
		Use:     "create infrastructure_id vm_type_id disk_size_gb instance_count [os_template_id]",
		Aliases: []string{"new"},
		Short:   "Create a new VM instance group in an infrastructure",
		Long: `Create a new VM instance group in an infrastructure.

This command creates a new VM instance group with the specified configuration.
The group will contain multiple VM instances of the same type and configuration,
making it easier to manage and scale similar workloads.

ARGUMENTS:
  infrastructure_id  The ID of the infrastructure where the group will be created
  vm_type_id         The VM type ID defining CPU, memory, and other hardware specs
  disk_size_gb       The disk size in GB for each VM instance in the group
  instance_count     The number of VM instances to create in the group
  os_template_id     Optional. The OS template ID for the VM instances

EXAMPLES:
  # Create a VM instance group with 3 instances using a specific OS template
  metalcloud-cli vm-instance-group create 12345 vm-type-1 100 3 ubuntu-20.04
  
  # Create a VM instance group without specifying OS template
  metalcloud-cli vmg new 12345 vm-type-small 50 2
  
  # Create a larger group for high-performance workloads
  metalcloud-cli vm-group create 12345 vm-type-performance 500 10 centos-8`,
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
		Use:     "update infrastructure_id vm_instance_group_id",
		Aliases: []string{"edit"},
		Short:   "Update VM instance group configuration",
		Long: `Update VM instance group configuration.

This command allows you to modify the configuration of an existing VM instance
group. You can update the label or custom variables associated with the group.
At least one of the available flags must be specified.

ARGUMENTS:
  infrastructure_id     The ID of the infrastructure containing the group
  vm_instance_group_id  The ID of the VM instance group to update

FLAGS:
  --label string                      Set or update the VM instance group label
  --custom-variables-source string    Source of custom variables to apply
                                     Can be 'pipe' for stdin or path to a JSON file

FLAG DEPENDENCIES:
  At least one of --label or --custom-variables-source must be provided

EXAMPLES:
  # Update the label of a VM instance group
  metalcloud-cli vm-instance-group update 12345 67890 --label "Web Servers"
  
  # Update custom variables from a JSON file
  metalcloud-cli vmg edit 12345 67890 --custom-variables-source /path/to/vars.json
  
  # Update custom variables from stdin
  echo '{"env": "production"}' | metalcloud-cli vm-group update 12345 67890 --custom-variables-source pipe
  
  # Update both label and custom variables
  metalcloud-cli vm-instance-group update 12345 67890 --label "Production Web" --custom-variables-source vars.json`,
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
		Use:     "delete infrastructure_id vm_instance_group_id",
		Aliases: []string{"rm"},
		Short:   "Delete a VM instance group",
		Long: `Delete a VM instance group from an infrastructure.

This command permanently removes a VM instance group and all its instances
from the specified infrastructure. This action is irreversible and will
terminate all running instances within the group.

WARNING: This operation cannot be undone. All data on the VM instances
will be permanently lost unless backed up elsewhere.

ARGUMENTS:
  infrastructure_id     The ID of the infrastructure containing the group
  vm_instance_group_id  The ID of the VM instance group to delete

EXAMPLES:
  # Delete VM instance group 67890 from infrastructure 12345
  metalcloud-cli vm-instance-group delete 12345 67890
  
  # Delete group using alias
  metalcloud-cli vmg rm 12345 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGroupDelete(cmd.Context(), args[0], args[1])
		},
	}

	vmInstanceGroupInstancesCmd = &cobra.Command{
		Use:     "instances infrastructure_id vm_instance_group_id",
		Aliases: []string{"instances-list", "instances-ls"},
		Short:   "List VM instances within a VM instance group",
		Long: `List all VM instances within a specific VM instance group.

This command displays all individual VM instances that belong to the specified
VM instance group. The output includes instance details such as ID, status,
IP addresses, and other relevant information for each instance in the group.

ARGUMENTS:
  infrastructure_id     The ID of the infrastructure containing the group
  vm_instance_group_id  The ID of the VM instance group to list instances from

EXAMPLES:
  # List all instances in VM instance group 67890 from infrastructure 12345
  metalcloud-cli vm-instance-group instances 12345 67890
  
  # List instances using alias
  metalcloud-cli vmg instances-ls 12345 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_INSTANCE_GROUPS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGroupInstances(cmd.Context(), args[0], args[1])
		},
	}
)

// VM Instance management commands.
var (
	vmInstanceCmd = &cobra.Command{
		Use:     "vm-instance [command]",
		Aliases: []string{"vmi", "vm"},
		Short:   "Manage individual VM instances within infrastructures",
		Long: `Manage individual VM instances within infrastructures.

VM instances are individual virtual machines that can be created, managed,
and controlled independently. This includes operations like getting instance
details, listing instances, managing power states, and accessing configuration.

Available commands:
  get          Get details of a specific VM instance
  list         List all VM instances in an infrastructure
  config       Get VM instance configuration
  start        Start a VM instance
  shutdown     Shutdown a VM instance
  reboot       Reboot a VM instance
  power-status Get VM instance power status

Examples:
  metalcloud-cli vm-instance list 12345
  metalcloud-cli vmi get 12345 67890
  metalcloud-cli vm start 12345 67890`,
	}

	vmInstanceGetCmd = &cobra.Command{
		Use:     "get infrastructure_id vm_instance_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific VM instance",
		Long: `Get detailed information about a specific VM instance.

This command retrieves comprehensive information about a VM instance including
its current status, configuration, network details, disk information, and
associated metadata. This is useful for debugging, monitoring, and understanding
the current state of a virtual machine.

ARGUMENTS:
  infrastructure_id  The ID of the infrastructure containing the VM instance
  vm_instance_id     The ID of the VM instance to retrieve details for

EXAMPLES:
  # Get details of VM instance 67890 in infrastructure 12345
  metalcloud-cli vm-instance get 12345 67890
  
  # Get instance details using alias
  metalcloud-cli vmi show 12345 67890
  
  # Get instance details using short alias
  metalcloud-cli vm get 12345 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGet(cmd.Context(), args[0], args[1])
		},
	}

	vmInstanceListCmd = &cobra.Command{
		Use:     "list infrastructure_id",
		Aliases: []string{"ls"},
		Short:   "List all VM instances in an infrastructure",
		Long: `List all VM instances in an infrastructure.

This command retrieves and displays all VM instances that exist within the
specified infrastructure. The output includes instance details such as ID,
status, VM type, IP addresses, and other relevant information for each instance.

ARGUMENTS:
  infrastructure_id  The ID of the infrastructure to list instances from

EXAMPLES:
  # List all VM instances in infrastructure 12345
  metalcloud-cli vm-instance list 12345
  
  # List instances using alias
  metalcloud-cli vmi ls 12345
  
  # List instances using short alias
  metalcloud-cli vm list 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceList(cmd.Context(), args[0])
		},
	}

	vmInstanceGetConfigCmd = &cobra.Command{
		Use:     "config infrastructure_id vm_instance_id",
		Aliases: []string{"get-config"},
		Short:   "Get VM instance configuration",
		Long: `Get VM instance configuration details.

This command retrieves the current configuration of a VM instance including
hardware specifications, network settings, disk configuration, and other
system parameters. This is useful for auditing, troubleshooting, and
understanding the VM's current setup.

ARGUMENTS:
  infrastructure_id  The ID of the infrastructure containing the VM instance
  vm_instance_id     The ID of the VM instance to retrieve configuration for

EXAMPLES:
  # Get configuration for VM instance 67890 in infrastructure 12345
  metalcloud-cli vm-instance config 12345 67890
  
  # Get configuration using alias
  metalcloud-cli vmi get-config 12345 67890
  
  # Get configuration using short alias
  metalcloud-cli vm config 12345 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGetConfig(cmd.Context(), args[0], args[1])
		},
	}

	vmInstanceStartCmd = &cobra.Command{
		Use:   "start infrastructure_id vm_instance_id",
		Short: "Start a VM instance",
		Long: `Start a VM instance.

This command initiates the startup process for a VM instance that is currently
powered off or stopped. The instance will be powered on and begin booting
according to its configured operating system and startup settings.

ARGUMENTS:
  infrastructure_id  The ID of the infrastructure containing the VM instance
  vm_instance_id     The ID of the VM instance to start

EXAMPLES:
  # Start VM instance 67890 in infrastructure 12345
  metalcloud-cli vm-instance start 12345 67890
  
  # Start instance using alias
  metalcloud-cli vm start 12345 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstancePowerControl(cmd.Context(), args[0], args[1], "start")
		},
	}

	vmInstanceShutdownCmd = &cobra.Command{
		Use:   "shutdown infrastructure_id vm_instance_id",
		Short: "Shutdown a VM instance",
		Long: `Shutdown a VM instance gracefully.

This command initiates a graceful shutdown process for a running VM instance.
The instance will receive a shutdown signal and will attempt to properly
terminate all running processes before powering off. This is the recommended
way to stop a VM instance to prevent data loss.

ARGUMENTS:
  infrastructure_id  The ID of the infrastructure containing the VM instance
  vm_instance_id     The ID of the VM instance to shutdown

EXAMPLES:
  # Shutdown VM instance 67890 in infrastructure 12345
  metalcloud-cli vm-instance shutdown 12345 67890
  
  # Shutdown instance using alias
  metalcloud-cli vm shutdown 12345 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstancePowerControl(cmd.Context(), args[0], args[1], "shutdown")
		},
	}

	vmInstanceRebootCmd = &cobra.Command{
		Use:   "reboot infrastructure_id vm_instance_id",
		Short: "Reboot a VM instance",
		Long: `Reboot a VM instance.

This command initiates a restart process for a running VM instance. The instance
will be gracefully shutdown and then automatically restarted. This is useful
for applying configuration changes or recovering from software issues.

ARGUMENTS:
  infrastructure_id  The ID of the infrastructure containing the VM instance
  vm_instance_id     The ID of the VM instance to reboot

EXAMPLES:
  # Reboot VM instance 67890 in infrastructure 12345
  metalcloud-cli vm-instance reboot 12345 67890
  
  # Reboot instance using alias
  metalcloud-cli vm reboot 12345 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstancePowerControl(cmd.Context(), args[0], args[1], "reboot")
		},
	}

	vmInstancePowerStatusCmd = &cobra.Command{
		Use:   "power-status infrastructure_id vm_instance_id",
		Short: "Get VM instance power status",
		Long: `Get VM instance power status.

This command retrieves the current power state of a VM instance, indicating
whether it is running, stopped, starting, stopping, or in another power state.
This is useful for monitoring and understanding the operational status of
virtual machines.

ARGUMENTS:
  infrastructure_id  The ID of the infrastructure containing the VM instance
  vm_instance_id     The ID of the VM instance to check power status for

EXAMPLES:
  # Get power status of VM instance 67890 in infrastructure 12345
  metalcloud-cli vm-instance power-status 12345 67890
  
  # Get power status using alias
  metalcloud-cli vm power-status 12345 67890`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
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

	// VM Instance management commands
	rootCmd.AddCommand(vmInstanceCmd)

	vmInstanceCmd.AddCommand(vmInstanceGetCmd)

	vmInstanceCmd.AddCommand(vmInstanceListCmd)

	vmInstanceCmd.AddCommand(vmInstanceGetConfigCmd)

	vmInstanceCmd.AddCommand(vmInstanceStartCmd)

	vmInstanceCmd.AddCommand(vmInstanceShutdownCmd)

	vmInstanceCmd.AddCommand(vmInstanceRebootCmd)

	vmInstanceCmd.AddCommand(vmInstancePowerStatusCmd)
}
