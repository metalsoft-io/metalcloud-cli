package cmd

import (
	"strconv"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/vm_pool"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	vmPoolFlags = struct {
		filterType   []string
		configSource string
		limit        string
		page         string
	}{}

	vmPoolCmd = &cobra.Command{
		Use:     "vm-pool [command]",
		Aliases: []string{"vmp"},
		Short:   "Manage virtual machine pools and their resources",
		Long: `Manage virtual machine pools including VMware vSphere, Hyper-V, and other hypervisor environments.

VM pools provide centralized management of virtualization infrastructure, allowing you to:
- Create and configure connections to hypervisor management systems
- Monitor VM and cluster host resources
- Manage credentials and certificates for secure access
- Control maintenance and experimental modes

Available commands support full lifecycle management from initial configuration
to ongoing monitoring and resource inspection.`,
	}

	vmPoolListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all VM pools with optional filtering",
		Long: `List all virtual machine pools in the environment with optional filtering by type.

This command displays a table showing VM pool details including ID, site, name, type, 
management host and port, and current status.

FLAGS:
  --filter-type    Filter results by VM pool type (can be used multiple times)
                   Available types: vmware, hyperv, kvm, xen

EXAMPLES:
  # List all VM pools
  metalcloud-cli vm-pool list

  # List only VMware pools
  metalcloud-cli vm-pool list --filter-type vmware

  # List VMware and Hyper-V pools
  metalcloud-cli vm-pool list --filter-type vmware --filter-type hyperv`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolList(cmd.Context(), vmPoolFlags.filterType)
		},
	}

	vmPoolGetCmd = &cobra.Command{
		Use:     "get vm_pool_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific VM pool",
		Long: `Get comprehensive details about a virtual machine pool including configuration, 
status, and connection information.

ARGUMENTS:
  vm_pool_id       The numeric ID of the VM pool to retrieve

EXAMPLES:
  # Get details for VM pool with ID 123
  metalcloud-cli vm-pool get 123

  # Using alias
  metalcloud-cli vm-pool show 123`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolGet(cmd.Context(), args[0])
		},
	}

	vmPoolConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Display a complete VM pool configuration example",
		Long: `Display a sample VM pool configuration in JSON format showing all available fields.

This command outputs an example configuration that can be used as a template for creating
new VM pools. The example includes all required and optional fields with sample values.

The generated configuration can be:
- Saved to a file and modified as needed
- Used directly with the 'create' command via pipe

EXAMPLES:
  # Display configuration example
  metalcloud-cli vm-pool config-example

  # Save example to file for editing
  metalcloud-cli vm-pool config-example > vmpool-config.json

  # Create VM pool using example as template
  metalcloud-cli vm-pool config-example | jq '.name = "my-vmpool"' | metalcloud-cli vm-pool create --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolConfigExample(cmd.Context())
		},
	}

	vmPoolCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new VM pool from configuration file or pipe",
		Long: `Create a new virtual machine pool from a JSON configuration file or piped input.

This command creates a VM pool by reading configuration from either a file or standard input.
The configuration must include all required fields and may include optional fields for 
complete setup.

REQUIRED FLAGS:
  --config-source  Source of the VM pool configuration (required)
                   Values: 'pipe' for stdin input, or path to JSON file

CONFIGURATION FIELDS:
  Required:
  - siteId         Site ID where the VM pool will be created
  - managementHost Hostname or IP of the hypervisor management interface
  - managementPort Port for management interface (typically 443 for VMware)
  - name           Name for the VM pool
  - type           VM pool type (e.g., vmware, hyperv, kvm, xen)

  Optional:
  - description    Descriptive text for the VM pool
  - certificate    TLS certificate for secure connections
  - privateKey     Private key corresponding to the certificate
  - username       Username for authentication (alternative to certificates)
  - password       Password for authentication (alternative to certificates)
  - inMaintenance  Set to 1 to create in maintenance mode (default: 0)
  - isExperimental Set to 1 to mark as experimental (default: 0)
  - tags           Array of string tags for categorization
  - options        Additional configuration options specific to the pool type

EXAMPLES:
  # Create from file
  metalcloud-cli vm-pool create --config-source vmpool.json

  # Create from pipe using config example as template
  metalcloud-cli vm-pool config-example | jq '.siteId = 2 | .name = "Production-VMware"' | metalcloud-cli vm-pool create --config-source pipe

  # Create minimal VMware pool from pipe
  echo '{"siteId":1,"managementHost":"vcenter.company.com","managementPort":443,"name":"Test-Pool","type":"vmware"}' | metalcloud-cli vm-pool create --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(vmPoolFlags.configSource)
			if err != nil {
				return err
			}

			return vm_pool.VMPoolCreate(cmd.Context(), config)
		},
	}

	vmPoolDeleteCmd = &cobra.Command{
		Use:   "delete vm_pool_id",
		Short: "Delete a VM pool permanently",
		Long: `Delete a virtual machine pool permanently from the environment.

This command removes the VM pool configuration and disconnects it from the hypervisor
management system. This action cannot be undone.

WARNING: Deleting a VM pool does not affect the actual virtual machines or hosts
in the hypervisor environment - only the MetalCloud management connection is removed.

ARGUMENTS:
  vm_pool_id       The numeric ID of the VM pool to delete

EXAMPLES:
  # Delete VM pool with ID 123
  metalcloud-cli vm-pool delete 123`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolDelete(cmd.Context(), args[0])
		},
	}

	vmPoolGetCredentialsCmd = &cobra.Command{
		Use:   "credentials vm_pool_id",
		Short: "Retrieve authentication credentials for a VM pool",
		Long: `Retrieve and display the authentication credentials configured for a VM pool.

This command shows the credentials used by MetalCloud to connect to the hypervisor
management interface. Sensitive information like passwords and private keys are
typically masked or encrypted in the output.

ARGUMENTS:
  vm_pool_id       The numeric ID of the VM pool

EXAMPLES:
  # Get credentials for VM pool 123
  metalcloud-cli vm-pool credentials 123`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolGetCredentials(cmd.Context(), args[0])
		},
	}

	vmPoolGetVMsCmd = &cobra.Command{
		Use:   "vms vm_pool_id",
		Short: "List virtual machines in a VM pool with pagination",
		Long: `List all virtual machines present in a specific VM pool with optional pagination.

This command displays VMs that are currently managed by the specified VM pool,
including their status, resource allocation, and basic configuration details.

ARGUMENTS:
  vm_pool_id       The numeric ID of the VM pool

OPTIONAL FLAGS:
  --limit          Number of records to return per page (default: all)
  --page           Page number to retrieve (1-based, default: 1)
                   Only effective when --limit is specified

PAGINATION:
When using pagination, specify both --limit and --page for best results.
The --limit flag controls how many records are returned, while --page
specifies which page of results to retrieve.

EXAMPLES:
  # List all VMs in VM pool 123
  metalcloud-cli vm-pool vms 123

  # List first 10 VMs
  metalcloud-cli vm-pool vms 123 --limit 10

  # List second page of 10 VMs each
  metalcloud-cli vm-pool vms 123 --limit 10 --page 2`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var limit, page float32
			if vmPoolFlags.limit != "" {
				limitVal, err := strconv.ParseFloat(vmPoolFlags.limit, 32)
				if err != nil {
					return err
				}
				limit = float32(limitVal)
			}
			if vmPoolFlags.page != "" {
				pageVal, err := strconv.ParseFloat(vmPoolFlags.page, 32)
				if err != nil {
					return err
				}
				page = float32(pageVal)
			}
			return vm_pool.VMPoolGetVMs(cmd.Context(), args[0], limit, page)
		},
	}

	vmPoolGetClusterHostsCmd = &cobra.Command{
		Use:   "cluster-hosts vm_pool_id",
		Short: "List cluster hosts in a VM pool with pagination",
		Long: `List all cluster hosts (ESXi hosts, Hyper-V servers, etc.) in a specific VM pool with optional pagination.

This command displays the hypervisor hosts that are part of the specified VM pool,
including their status, resource utilization, and connection details.

ARGUMENTS:
  vm_pool_id       The numeric ID of the VM pool

OPTIONAL FLAGS:
  --limit          Number of records to return per page (default: all)
  --page           Page number to retrieve (1-based, default: 1)
                   Only effective when --limit is specified

PAGINATION:
When using pagination, specify both --limit and --page for best results.
The --limit flag controls how many records are returned, while --page
specifies which page of results to retrieve.

EXAMPLES:
  # List all cluster hosts in VM pool 123
  metalcloud-cli vm-pool cluster-hosts 123

  # List first 5 hosts
  metalcloud-cli vm-pool cluster-hosts 123 --limit 5

  # List second page of 5 hosts each
  metalcloud-cli vm-pool cluster-hosts 123 --limit 5 --page 2`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var limit, page float32
			if vmPoolFlags.limit != "" {
				limitVal, err := strconv.ParseFloat(vmPoolFlags.limit, 32)
				if err != nil {
					return err
				}
				limit = float32(limitVal)
			}
			if vmPoolFlags.page != "" {
				pageVal, err := strconv.ParseFloat(vmPoolFlags.page, 32)
				if err != nil {
					return err
				}
				page = float32(pageVal)
			}
			return vm_pool.VMPoolGetClusterHosts(cmd.Context(), args[0], limit, page)
		},
	}

	vmPoolGetClusterHostVMsCmd = &cobra.Command{
		Use:   "cluster-host-vms vm_pool_id host_id",
		Short: "List virtual machines on a specific cluster host with pagination",
		Long: `List all virtual machines running on a specific cluster host within a VM pool with optional pagination.

This command displays VMs that are currently running on the specified cluster host,
including their status, resource allocation, and configuration details.

ARGUMENTS:
  vm_pool_id       The numeric ID of the VM pool
  host_id          The numeric ID of the cluster host

OPTIONAL FLAGS:
  --limit          Number of records to return per page (default: all)
  --page           Page number to retrieve (1-based, default: 1)
                   Only effective when --limit is specified

PAGINATION:
When using pagination, specify both --limit and --page for best results.
The --limit flag controls how many records are returned, while --page
specifies which page of results to retrieve.

EXAMPLES:
  # List all VMs on cluster host 456 in VM pool 123
  metalcloud-cli vm-pool cluster-host-vms 123 456

  # List first 5 VMs on the host
  metalcloud-cli vm-pool cluster-host-vms 123 456 --limit 5

  # List second page of 5 VMs each
  metalcloud-cli vm-pool cluster-host-vms 123 456 --limit 5 --page 2`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var limit, page float32
			if vmPoolFlags.limit != "" {
				limitVal, err := strconv.ParseFloat(vmPoolFlags.limit, 32)
				if err != nil {
					return err
				}
				limit = float32(limitVal)
			}
			if vmPoolFlags.page != "" {
				pageVal, err := strconv.ParseFloat(vmPoolFlags.page, 32)
				if err != nil {
					return err
				}
				page = float32(pageVal)
			}
			return vm_pool.VMPoolGetClusterHostVMs(cmd.Context(), args[0], args[1], limit, page)
		},
	}

	vmPoolGetClusterHostInterfacesCmd = &cobra.Command{
		Use:   "cluster-host-interfaces vm_pool_id host_id",
		Short: "List network interfaces for a cluster host in a VM pool",
		Long: `List all network interfaces available on a specific cluster host within a VM pool.

This command displays the network interfaces that are configured on the specified 
cluster host, including their status, configuration, and network details.

ARGUMENTS:
  vm_pool_id       The numeric ID of the VM pool
  host_id          The numeric ID of the cluster host

EXAMPLES:
  # List interfaces for cluster host 456 in VM pool 123
  metalcloud-cli vm-pool cluster-host-interfaces 123 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolGetClusterHostInterfaces(cmd.Context(), args[0], args[1])
		},
	}
)

func init() {
	rootCmd.AddCommand(vmPoolCmd)

	// List command
	vmPoolCmd.AddCommand(vmPoolListCmd)
	vmPoolListCmd.Flags().StringSliceVar(&vmPoolFlags.filterType, "filter-type", nil, "Filter the result by VM pool type (e.g., vmware, hyperv, etc).")

	// Get command
	vmPoolCmd.AddCommand(vmPoolGetCmd)

	// Config example command
	vmPoolCmd.AddCommand(vmPoolConfigExampleCmd)

	// Create command
	vmPoolCmd.AddCommand(vmPoolCreateCmd)
	vmPoolCreateCmd.Flags().StringVar(&vmPoolFlags.configSource, "config-source", "", "Source of the new VM pool configuration. Can be 'pipe' or path to a JSON file.")
	vmPoolCreateCmd.MarkFlagsOneRequired("config-source")

	// Delete command
	vmPoolCmd.AddCommand(vmPoolDeleteCmd)

	// Get credentials command
	vmPoolCmd.AddCommand(vmPoolGetCredentialsCmd)

	// Get VMs command
	vmPoolCmd.AddCommand(vmPoolGetVMsCmd)
	vmPoolGetVMsCmd.Flags().StringVar(&vmPoolFlags.limit, "limit", "", "Number of records per page")
	vmPoolGetVMsCmd.Flags().StringVar(&vmPoolFlags.page, "page", "", "Page number")

	// Get cluster hosts command
	vmPoolCmd.AddCommand(vmPoolGetClusterHostsCmd)
	vmPoolGetClusterHostsCmd.Flags().StringVar(&vmPoolFlags.limit, "limit", "", "Number of records per page")
	vmPoolGetClusterHostsCmd.Flags().StringVar(&vmPoolFlags.page, "page", "", "Page number")

	// Get cluster host VMs command
	vmPoolCmd.AddCommand(vmPoolGetClusterHostVMsCmd)
	vmPoolGetClusterHostVMsCmd.Flags().StringVar(&vmPoolFlags.limit, "limit", "", "Number of records per page")
	vmPoolGetClusterHostVMsCmd.Flags().StringVar(&vmPoolFlags.page, "page", "", "Page number")

	// Get cluster host interfaces command
	vmPoolCmd.AddCommand(vmPoolGetClusterHostInterfacesCmd)
}
