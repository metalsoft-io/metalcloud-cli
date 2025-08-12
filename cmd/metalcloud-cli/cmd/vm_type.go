package cmd

import (
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/vm_type"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	vmTypeFlags = struct {
		configSource string
		limit        string
		page         string
	}{}

	vmTypeCmd = &cobra.Command{
		Use:     "vm-type [command]",
		Aliases: []string{"vmt"},
		Short:   "Manage VM types and configurations",
		Long: `Manage VM types and their configurations in the MetalCloud platform.

VM types define the resource specifications (CPU cores, RAM) for virtual machines. 
They can be experimental or production-ready, and can be restricted to unmanaged VMs only.

Available commands:
  list          List all VM types with pagination support
  get           Get detailed information about a specific VM type
  create        Create a new VM type from configuration
  update        Update an existing VM type configuration
  delete        Delete a VM type
  vms           List all VMs using a specific VM type
  config-example Show an example configuration for creating VM types

Use "metalcloud vm-type [command] --help" for more information about a command.`,
	}

	vmTypeListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all VM types with optional pagination",
		Long: `List all VM types available in the MetalCloud platform.

This command displays all VM types with their specifications including CPU cores, RAM, 
experimental status, and whether they are restricted to unmanaged VMs only.

FLAGS:
  --limit string    Number of records per page (optional)
  --page string     Page number to retrieve (optional, requires --limit)

EXAMPLES:
  # List all VM types
  metalcloud vm-type list
  
  # List VM types with pagination (10 per page, page 1)
  metalcloud vm-type list --limit 10 --page 1
  
  # List first 5 VM types
  metalcloud vm-type list --limit 5`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_TYPES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			var limit, page float32
			if vmTypeFlags.limit != "" {
				limitVal, err := strconv.ParseFloat(vmTypeFlags.limit, 32)
				if err != nil {
					return err
				}
				limit = float32(limitVal)
			}
			if vmTypeFlags.page != "" {
				pageVal, err := strconv.ParseFloat(vmTypeFlags.page, 32)
				if err != nil {
					return err
				}
				page = float32(pageVal)
			}
			return vm_type.VMTypeList(cmd.Context(), limit, page)
		},
	}

	vmTypeGetCmd = &cobra.Command{
		Use:     "get vm_type_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific VM type",
		Long: `Get detailed information about a specific VM type including its specifications, 
experimental status, and configuration details.

ARGUMENTS:
  vm_type_id    The numeric ID of the VM type to retrieve

EXAMPLES:
  # Get VM type with ID 1
  metalcloud vm-type get 1
  
  # Using alias
  metalcloud vm-type show 1`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_TYPES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_type.VMTypeGet(cmd.Context(), args[0])
		},
	}

	vmTypeConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Show an example configuration for creating VM types",
		Long: `Display an example JSON configuration that can be used to create a new VM type.

This command outputs a sample configuration showing all available fields and their expected values.
The generated configuration can be saved to a file and modified as needed for creating or updating VM types.

EXAMPLES:
  # Display configuration example
  metalcloud vm-type config-example
  
  # Save example to file for editing
  metalcloud vm-type config-example > vm-type-config.json
  
  # Use saved configuration to create a VM type
  metalcloud vm-type create --config-source vm-type-config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_TYPES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_type.VMTypeConfigExample(cmd.Context())
		},
	}

	vmTypeCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new VM type from configuration",
		Long: `Create a new VM type in the MetalCloud platform using a configuration file or piped input.

The configuration must be provided in JSON format and include all required fields:
- name: Unique name for the VM type
- cpuCores: Number of CPU cores
- ramGB: Amount of RAM in gigabytes

Optional fields include displayName, label, isExperimental, forUnmanagedVMsOnly, and tags.

REQUIRED FLAGS:
  --config-source string    Source of the VM type configuration (required)
                           Can be 'pipe' for stdin or path to a JSON file

EXAMPLES:
  # Create VM type from file
  metalcloud vm-type create --config-source vm-type-config.json
  
  # Create VM type from stdin
  echo '{"name":"test-vm","cpuCores":2,"ramGB":4}' | metalcloud vm-type create --config-source pipe
  
  # Generate example config and create VM type
  metalcloud vm-type config-example > config.json
  # Edit config.json as needed
  metalcloud vm-type create --config-source config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_TYPES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(vmTypeFlags.configSource)
			if err != nil {
				return err
			}

			return vm_type.VMTypeCreate(cmd.Context(), config)
		},
	}

	vmTypeUpdateCmd = &cobra.Command{
		Use:   "update vm_type_id",
		Short: "Update an existing VM type configuration",
		Long: `Update an existing VM type in the MetalCloud platform using a configuration file or piped input.

The configuration must be provided in JSON format. Only the fields you want to update need to be included.
You can modify any of the following fields:
- name: VM type name
- displayName: Display name for the VM type
- label: Label for the VM type
- cpuCores: Number of CPU cores
- ramGB: Amount of RAM in gigabytes
- isExperimental: Whether the VM type is experimental (0 or 1)
- forUnmanagedVMsOnly: Whether restricted to unmanaged VMs (0 or 1)
- tags: Array of tags

ARGUMENTS:
  vm_type_id    The numeric ID of the VM type to update

REQUIRED FLAGS:
  --config-source string    Source of the VM type configuration (required)
                           Can be 'pipe' for stdin or path to a JSON file

EXAMPLES:
  # Update VM type from file
  metalcloud vm-type update 123 --config-source vm-type-update.json
  
  # Update VM type from stdin
  echo '{"cpuCores":8,"ramGB":16}' | metalcloud vm-type update 123 --config-source pipe
  
  # Generate example config and update VM type
  metalcloud vm-type config-example > config.json
  # Edit config.json as needed
  metalcloud vm-type update 123 --config-source config.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_TYPES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(vmTypeFlags.configSource)
			if err != nil {
				return err
			}

			return vm_type.VMTypeUpdate(cmd.Context(), args[0], config)
		},
	}

	vmTypeDeleteCmd = &cobra.Command{
		Use:   "delete vm_type_id",
		Short: "Delete a VM type",
		Long: `Delete a VM type from the MetalCloud platform.

WARNING: This action is irreversible. Ensure that no VMs are currently using this VM type
before deletion, as this could cause issues with existing deployments.

ARGUMENTS:
  vm_type_id    The numeric ID of the VM type to delete

EXAMPLES:
  # Delete VM type with ID 123
  metalcloud vm-type delete 123
  
  # Check VMs using a VM type before deletion
  metalcloud vm-type vms 123
  metalcloud vm-type delete 123`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_TYPES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_type.VMTypeDelete(cmd.Context(), args[0])
		},
	}

	vmTypeGetVMsCmd = &cobra.Command{
		Use:   "vms vm_type_id",
		Short: "List all VMs using a specific VM type",
		Long: `List all virtual machines that are currently using a specific VM type.

This command helps you understand which VMs are dependent on a particular VM type,
which is useful before making changes or deleting the VM type.

ARGUMENTS:
  vm_type_id    The numeric ID of the VM type to check

OPTIONAL FLAGS:
  --limit string    Number of records per page (optional)
  --page string     Page number to retrieve (optional, requires --limit)

EXAMPLES:
  # List all VMs using VM type 123
  metalcloud vm-type vms 123
  
  # List VMs with pagination (10 per page, page 1)
  metalcloud vm-type vms 123 --limit 10 --page 1
  
  # Check VM usage before deletion
  metalcloud vm-type vms 123
  metalcloud vm-type delete 123`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_TYPES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var limit, page float32
			if vmTypeFlags.limit != "" {
				limitVal, err := strconv.ParseFloat(vmTypeFlags.limit, 32)
				if err != nil {
					return err
				}
				limit = float32(limitVal)
			}
			if vmTypeFlags.page != "" {
				pageVal, err := strconv.ParseFloat(vmTypeFlags.page, 32)
				if err != nil {
					return err
				}
				page = float32(pageVal)
			}
			return vm_type.VMTypeGetVMs(cmd.Context(), args[0], limit, page)
		},
	}
)

func init() {
	rootCmd.AddCommand(vmTypeCmd)

	// List command
	vmTypeCmd.AddCommand(vmTypeListCmd)
	vmTypeListCmd.Flags().StringVar(&vmTypeFlags.limit, "limit", "", "Number of records per page")
	vmTypeListCmd.Flags().StringVar(&vmTypeFlags.page, "page", "", "Page number")

	// Get command
	vmTypeCmd.AddCommand(vmTypeGetCmd)

	// Config example command
	vmTypeCmd.AddCommand(vmTypeConfigExampleCmd)

	// Create command
	vmTypeCmd.AddCommand(vmTypeCreateCmd)
	vmTypeCreateCmd.Flags().StringVar(&vmTypeFlags.configSource, "config-source", "", "Source of the new VM type configuration. Can be 'pipe' or path to a JSON file.")
	vmTypeCreateCmd.MarkFlagsOneRequired("config-source")

	// Update command
	vmTypeCmd.AddCommand(vmTypeUpdateCmd)
	vmTypeUpdateCmd.Flags().StringVar(&vmTypeFlags.configSource, "config-source", "", "Source of the VM type update configuration. Can be 'pipe' or path to a JSON file.")
	vmTypeUpdateCmd.MarkFlagsOneRequired("config-source")

	// Delete command
	vmTypeCmd.AddCommand(vmTypeDeleteCmd)

	// Get VMs command
	vmTypeCmd.AddCommand(vmTypeGetVMsCmd)
	vmTypeGetVMsCmd.Flags().StringVar(&vmTypeFlags.limit, "limit", "", "Number of records per page")
	vmTypeGetVMsCmd.Flags().StringVar(&vmTypeFlags.page, "page", "", "Page number")
}
