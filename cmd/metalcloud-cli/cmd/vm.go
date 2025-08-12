package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/vm"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

// VM commands
var (
	vmFlags = struct {
		configSource string
		powerAction  string
	}{}

	vmCmd = &cobra.Command{
		Use:     "vm [command]",
		Aliases: []string{"vms", "virtual-machine"},
		Short:   "Manage virtual machines lifecycle and configuration",
		Long: `Comprehensive virtual machine management commands for controlling VM power states,
configuration updates, and monitoring. Supports operations like start, stop, reboot,
status checking, and configuration updates through JSON files or pipes.

Available Commands:
  get            Retrieve detailed VM information and configuration
  power-status   Check current power state of a VM
  start          Power on a VM
  shutdown       Gracefully shutdown or force stop a VM  
  reboot         Restart a VM
  update         Update VM configuration from JSON file or pipe
  console-info   Get remote console connection details

Examples:
  metalcloud-cli vm get 12345
  metalcloud-cli vm start 12345
  metalcloud-cli vm update 12345 --config-source config.json
  cat vm-config.json | metalcloud-cli vm update 12345 --config-source pipe`,
	}

	vmGetCmd = &cobra.Command{
		Use:     "get vm_id",
		Aliases: []string{"show"},
		Short:   "Retrieve detailed VM information and configuration",
		Long: `Retrieve comprehensive information about a virtual machine including its current
configuration, power state, network settings, storage details, and metadata.

Arguments:
  vm_id          Required. The unique identifier of the virtual machine.

Examples:
  # Get VM information by ID
  metalcloud-cli vm get 12345
  
  # Get VM info using alias
  metalcloud-cli vm show 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VMS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm.VMGet(cmd.Context(), args[0])
		},
	}

	vmPowerStatusCmd = &cobra.Command{
		Use:   "power-status vm_id",
		Short: "Check current power state of a VM",
		Long: `Check the current power status of a virtual machine. Returns the current 
power state such as 'running', 'stopped', 'suspended', or 'unknown'.

Arguments:
  vm_id          Required. The unique identifier of the virtual machine.

Examples:
  # Check power status of VM
  metalcloud-cli vm power-status 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VMS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm.VMPowerStatus(cmd.Context(), args[0])
		},
	}

	vmStartCmd = &cobra.Command{
		Use:   "start vm_id",
		Short: "Power on a virtual machine",
		Long: `Start (power on) a virtual machine. The VM will boot up and become available
for connections. The operation is asynchronous - the command returns immediately
while the VM powers on in the background.

Arguments:
  vm_id          Required. The unique identifier of the virtual machine to start.

Prerequisites:
  - VM must be in 'stopped' or 'suspended' state
  - Valid VM configuration must exist
  - Sufficient infrastructure resources must be available

Examples:
  # Start a VM by ID
  metalcloud-cli vm start 12345
  
  # Start multiple VMs (using shell loops)
  for vm in 12345 12346 12347; do metalcloud-cli vm start $vm; done`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VMS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm.VMStart(cmd.Context(), args[0])
		},
	}

	vmShutdownCmd = &cobra.Command{
		Use:     "shutdown vm_id",
		Aliases: []string{"stop"},
		Short:   "Gracefully shutdown or force stop a VM",
		Long: `Shutdown a virtual machine gracefully or force stop it. The command attempts
a graceful shutdown first, which allows the guest OS to properly close running
applications and services. If the VM doesn't respond, it can be force stopped.

Arguments:
  vm_id          Required. The unique identifier of the virtual machine to shutdown.

Prerequisites:
  - VM must be in 'running' state
  - User must have write permissions for the VM

Examples:
  # Gracefully shutdown a VM
  metalcloud-cli vm shutdown 12345
  
  # Using the alias 'stop'
  metalcloud-cli vm stop 12345
  
  # Shutdown multiple VMs
  for vm in 12345 12346 12347; do metalcloud-cli vm shutdown $vm; done`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VMS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm.VMShutdown(cmd.Context(), args[0])
		},
	}

	vmRebootCmd = &cobra.Command{
		Use:   "reboot vm_id",
		Short: "Restart a virtual machine",
		Long: `Restart a virtual machine by performing a graceful reboot. The VM will be
shutdown gracefully and then started again. This is equivalent to performing
a shutdown followed by a start operation.

Arguments:
  vm_id          Required. The unique identifier of the virtual machine to reboot.

Prerequisites:
  - VM must be in 'running' state
  - User must have write permissions for the VM
  - VM configuration must be valid

Examples:
  # Reboot a VM by ID
  metalcloud-cli vm reboot 12345
  
  # Reboot multiple VMs sequentially
  for vm in 12345 12346 12347; do metalcloud-cli vm reboot $vm; done`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VMS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm.VMReboot(cmd.Context(), args[0])
		},
	}

	vmUpdateCmd = &cobra.Command{
		Use:   "update vm_id",
		Short: "Update VM configuration from JSON file or pipe",
		Long: `Update virtual machine configuration using JSON data from a file or pipe.
This command allows you to modify VM settings like CPU, memory, network interfaces,
and other configuration parameters.

Arguments:
  vm_id          Required. The unique identifier of the virtual machine to update.

Required Flags:
  --config-source    Required. Source of the VM update configuration.
                     Can be 'pipe' (to read from stdin) or path to a JSON file.

Prerequisites:
  - VM must exist and be accessible
  - User must have write permissions for the VM
  - Configuration JSON must be valid and complete
  - VM may need to be stopped for certain configuration changes

Examples:
  # Update VM from JSON file
  metalcloud-cli vm update 12345 --config-source vm-config.json
  
  # Update VM from pipe (stdin)
  cat vm-config.json | metalcloud-cli vm update 12345 --config-source pipe
  
  # Update VM using inline JSON
  echo '{"cpu_count": 4, "memory_size_mb": 8192}' | metalcloud-cli vm update 12345 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VMS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(vmFlags.configSource)
			if err != nil {
				return err
			}
			return vm.VMUpdate(cmd.Context(), args[0], config)
		},
	}

	vmRemoteConsoleInfoCmd = &cobra.Command{
		Use:   "console-info vm_id",
		Short: "Get VM remote console connection details",
		Long: `Retrieve remote console connection information for a virtual machine.
This command provides the necessary details to establish a remote console connection
including connection URLs, protocols, and authentication credentials.

Arguments:
  vm_id          Required. The unique identifier of the virtual machine.

Prerequisites:
  - VM must exist and be accessible
  - User must have read permissions for the VM
  - VM must be in a state that supports console access

Examples:
  # Get console connection info for a VM
  metalcloud-cli vm console-info 12345`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VMS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm.VMRemoteConsoleInfo(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(vmCmd)

	// VM commands
	vmCmd.AddCommand(vmGetCmd)

	vmCmd.AddCommand(vmPowerStatusCmd)

	vmCmd.AddCommand(vmStartCmd)

	vmCmd.AddCommand(vmShutdownCmd)

	vmCmd.AddCommand(vmRebootCmd)

	vmCmd.AddCommand(vmUpdateCmd)
	vmUpdateCmd.Flags().StringVar(&vmFlags.configSource, "config-source", "", "Source of the VM update configuration. Can be 'pipe' or path to a JSON file.")
	vmUpdateCmd.MarkFlagsOneRequired("config-source")

	vmCmd.AddCommand(vmRemoteConsoleInfoCmd)
}
