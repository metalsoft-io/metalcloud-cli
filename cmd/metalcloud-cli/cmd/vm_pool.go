package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/vm_pool"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/cobra"
)

var (
	vmPoolFlags = struct {
		filterType       []string
		configSource     string
		limit            string
		page             string
		infrastructureId string
		vmNames          string
		// Each write command that takes --config-source keeps its own field:
		// runCLI resets the shared root command between test runs and flags
		// registered on different sub-commands must not share state.
		updateConfigSource             string
		clusterHostConfigSource        string
		hostInterfaceConfigSource      string
		networkDeviceConfigSource      string
		networkDeviceRef               string
		networkDeviceInterfaceNameFlag string
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
to ongoing monitoring and resource inspection.

Pool level commands:
  list, get, create, update, delete, config-example, update-config-example,
  credentials, statistics, containers, vms, sync, refresh, import-vms

Cluster host commands (flat, one command per operation):
  cluster-hosts, cluster-host, update-cluster-host, cluster-host-statistics,
  cluster-host-containers, cluster-host-vms, cluster-host-interfaces,
  cluster-host-interface, update-cluster-host-interface

Cluster host interface network device assignments live in their own sub-group,
because they are a full CRUD set of their own:
  cluster-host-interface-network-device (alias: chi-network-device, chind)
    list | get | add | remove | config-example`,
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

	vmPoolImportVMsCmd = &cobra.Command{
		Use:   "import-vms vm_pool_id",
		Short: "Import VMs from the hypervisor into a VM pool",
		Long: `Import virtual machines from the hypervisor management system into a VM pool.

This command imports VMs that exist in the hypervisor (e.g., VMware vCenter, Hyper-V)
but are not yet registered in the MetalCloud VM pool. The import configuration can be
specified either via a configuration file/pipe or by using flags.

ARGUMENTS:
  vm_pool_id       The numeric ID of the VM pool

CONFIGURATION VIA FLAGS:
  --vm-names            Comma-separated list of VM names to import
  --infrastructure-id   Infrastructure ID to associate with the imported VMs

CONFIGURATION VIA FILE/PIPE:
  --config-source      Source of the import configuration
                       Values: 'pipe' for stdin input, or path to JSON file

NOTE: Either use --config-source OR both --infrastructure-id and --vm-names flags.

CONFIGURATION FILE FIELDS:
  - vmNames            Array of VM names to import from the hypervisor
  - infrastructureId   Infrastructure ID to associate with the imported VMs

EXAMPLES:
  # Import VMs using flags
  metalcloud-cli vm-pool import-vms 123 --infrastructure-id 456 --vm-names "vm-prod-01,vm-prod-02,vm-test-01"

  # Import VMs using configuration from file
  metalcloud-cli vm-pool import-vms 123 --config-source import-config.json

  # Import VMs using piped configuration
  echo '{"infrastructureId": 456, "vmNames": ["vm-123", "vm-456"]}' | metalcloud-cli vm-pool import-vms 123 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var importVMs sdk.VMPoolImportVMs
			if vmPoolFlags.configSource != "" {
				config, err := utils.ReadConfigFromPipeOrFile(vmPoolFlags.configSource)
				if err != nil {
					return err
				}

				err = utils.UnmarshalContent(config, &importVMs)
				if err != nil {
					return err
				}
			} else {
				// Parse CSV list of VM names
				vmNamesList := strings.Split(vmPoolFlags.vmNames, ",")
				for i, name := range vmNamesList {
					vmNamesList[i] = strings.TrimSpace(name)
				}

				importVMs = sdk.VMPoolImportVMs{
					VmNames: vmNamesList,
				}

				if vmPoolFlags.infrastructureId != "" {
					infraId, err := strconv.ParseInt(vmPoolFlags.infrastructureId, 10, 64)
					if err != nil {
						return fmt.Errorf("invalid infrastructure ID: '%s'", vmPoolFlags.infrastructureId)
					}

					importVMs.InfrastructureId = sdk.PtrInt64(infraId)
				}
			}

			return vm_pool.VMPoolImportVMs(cmd.Context(), args[0], importVMs)
		},
	}

	vmPoolUpdateCmd = &cobra.Command{
		Use:   "update vm_pool_id",
		Short: "Update a VM pool from a configuration file or pipe",
		Long: `Update an existing virtual machine pool from a JSON or YAML configuration.

Only the fields present in the configuration are changed; everything else keeps
its current value. The endpoint does not use optimistic concurrency, so no
revision needs to be fetched first.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool to update

Required Flags:
  --config-source  Source of the VM pool update configuration.
                   Values: 'pipe' for stdin input, or path to a JSON/YAML file.

Configuration Fields (all optional):
  description, managementHost, managementPort, certificate, privateKey,
  username, password, inMaintenance, isExperimental, tags, options,
  networkFabricId

Examples:
  # Update from a file
  metalcloud-cli vm-pool update 123 --config-source update.json

  # Put a pool into maintenance mode
  echo '{"inMaintenance": 1}' | metalcloud-cli vm-pool update 123 --config-source pipe

  # Start from the generated example
  metalcloud-cli vm-pool update-config-example > update.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(vmPoolFlags.updateConfigSource)
			if err != nil {
				return err
			}

			return vm_pool.VMPoolUpdate(cmd.Context(), args[0], config)
		},
	}

	vmPoolUpdateConfigExampleCmd = &cobra.Command{
		Use:   "update-config-example",
		Short: "Display a VM pool update configuration example",
		Long: `Display a sample VM pool update payload showing the fields accepted by
'vm-pool update'. Every field is optional - delete the ones you do not want to
change.

Examples:
  # Display the example
  metalcloud-cli vm-pool update-config-example

  # Save it for editing
  metalcloud-cli vm-pool update-config-example > update.json

  # Change one field and apply it
  metalcloud-cli vm-pool update-config-example | jq '{description}' | metalcloud-cli vm-pool update 123 --config-source pipe`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolUpdateConfigExample(cmd.Context())
		},
	}

	vmPoolSyncCmd = &cobra.Command{
		Use:   "sync vm_pool_id",
		Short: "Sync a VM pool with its hypervisor",
		Long: `Start a synchronization job for a VM pool.

The sync discovers objects that were created directly on the hypervisor. On
VMware VCF, for example, it discovers new Virtual Distributed Switches. The
command returns the job information of the job that was started; use
'metalcloud-cli job get <job_id>' to follow it.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool to sync

Examples:
  # Sync VM pool 123
  metalcloud-cli vm-pool sync 123`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolSync(cmd.Context(), args[0])
		},
	}

	vmPoolRefreshCmd = &cobra.Command{
		Use:   "refresh vm_pool_id",
		Short: "Refresh the information of a VM pool",
		Long: `Refresh the cached information MetalSoft holds about a VM pool.

On VMware VCF this reports any new datastores. The refreshed VM pool is printed
when the operation completes.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool to refresh

Examples:
  # Refresh VM pool 123
  metalcloud-cli vm-pool refresh 123`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolRefresh(cmd.Context(), args[0])
		},
	}

	vmPoolStatisticsCmd = &cobra.Command{
		Use:     "statistics vm_pool_id",
		Aliases: []string{"stats"},
		Short:   "Show the resource usage statistics of a VM pool",
		Long: `Show the aggregated RAM, disk and GPU usage of a VM pool.

In the tabular formats (text, csv, md) the GPU list is flattened into a single
column; json and yaml return the full statistics object.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool

Examples:
  # Show statistics for VM pool 123
  metalcloud-cli vm-pool statistics 123

  # Get the raw statistics object
  metalcloud-cli vm-pool stats 123 -f json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolGetStatistics(cmd.Context(), args[0])
		},
	}

	vmPoolContainersCmd = &cobra.Command{
		Use:   "containers vm_pool_id",
		Short: "List the containers running in a VM pool",
		Long: `List all containers deployed on the hosts of a VM pool.

The listing walks every page transparently.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool

Examples:
  # List containers of VM pool 123
  metalcloud-cli vm-pool containers 123`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolGetContainers(cmd.Context(), args[0])
		},
	}

	vmPoolGetClusterHostCmd = &cobra.Command{
		Use:   "cluster-host vm_pool_id host_id",
		Short: "Get details of one cluster host of a VM pool",
		Long: `Get the full details of a single cluster host (ESXi host, Hyper-V server, ...)
that belongs to a VM pool, including its health status, roles and the flags that
control whether VMs and containers may be created on it.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool
  host_id          The numeric ID of the cluster host

Examples:
  # Get cluster host 456 of VM pool 123
  metalcloud-cli vm-pool cluster-host 123 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolGetClusterHost(cmd.Context(), args[0], args[1])
		},
	}

	vmPoolUpdateClusterHostCmd = &cobra.Command{
		Use:   "update-cluster-host vm_pool_id host_id",
		Short: "Update one cluster host of a VM pool",
		Long: `Update a cluster host of a VM pool from a JSON or YAML configuration.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool
  host_id          The numeric ID of the cluster host

Required Flags:
  --config-source  Source of the cluster host update configuration.
                   Values: 'pipe' for stdin input, or path to a JSON/YAML file.

Configuration Fields (all optional):
  allowVMsToBeCreated         Allow new VMs to be created on this host
  allowContainersToBeCreated  Allow new containers to be created on this host

Examples:
  # Stop scheduling new VMs on a host
  echo '{"allowVMsToBeCreated": false}' | metalcloud-cli vm-pool update-cluster-host 123 456 --config-source pipe

  # Update from a file
  metalcloud-cli vm-pool update-cluster-host 123 456 --config-source host.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(vmPoolFlags.clusterHostConfigSource)
			if err != nil {
				return err
			}

			return vm_pool.VMPoolUpdateClusterHost(cmd.Context(), args[0], args[1], config)
		},
	}

	vmPoolUpdateClusterHostConfigExampleCmd = &cobra.Command{
		Use:   "update-cluster-host-config-example",
		Short: "Display a cluster host update configuration example",
		Long: `Display a sample payload for 'vm-pool update-cluster-host'.

Examples:
  # Display the example
  metalcloud-cli vm-pool update-cluster-host-config-example

  # Save it for editing
  metalcloud-cli vm-pool update-cluster-host-config-example > host.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolUpdateClusterHostConfigExample(cmd.Context())
		},
	}

	vmPoolClusterHostStatisticsCmd = &cobra.Command{
		Use:     "cluster-host-statistics vm_pool_id host_id",
		Aliases: []string{"cluster-host-stats"},
		Short:   "Show the resource usage statistics of one cluster host",
		Long: `Show the RAM, disk and GPU usage of a single cluster host of a VM pool.

In the tabular formats (text, csv, md) the GPU list is flattened into a single
column; json and yaml return the full statistics object.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool
  host_id          The numeric ID of the cluster host

Examples:
  # Show statistics for cluster host 456 of VM pool 123
  metalcloud-cli vm-pool cluster-host-statistics 123 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolGetClusterHostStatistics(cmd.Context(), args[0], args[1])
		},
	}

	vmPoolClusterHostContainersCmd = &cobra.Command{
		Use:   "cluster-host-containers vm_pool_id host_id",
		Short: "List the containers running on one cluster host",
		Long: `List all containers deployed on a single cluster host of a VM pool.

The listing walks every page transparently.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool
  host_id          The numeric ID of the cluster host

Examples:
  # List containers on cluster host 456 of VM pool 123
  metalcloud-cli vm-pool cluster-host-containers 123 456`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolGetClusterHostContainers(cmd.Context(), args[0], args[1])
		},
	}

	vmPoolGetClusterHostInterfaceCmd = &cobra.Command{
		Use:   "cluster-host-interface vm_pool_id host_id interface_id",
		Short: "Get one network interface of a cluster host",
		Long: `Get the details of a single network interface of a cluster host, including its
MAC address, management status, fabric and network device assignments.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool
  host_id          The numeric ID of the cluster host
  interface_id     The numeric ID of the cluster host interface

Examples:
  # Get interface 789 of cluster host 456 in VM pool 123
  metalcloud-cli vm-pool cluster-host-interface 123 456 789`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_READ},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolGetClusterHostInterface(cmd.Context(), args[0], args[1], args[2])
		},
	}

	vmPoolUpdateClusterHostInterfaceCmd = &cobra.Command{
		Use:   "update-cluster-host-interface vm_pool_id host_id interface_id",
		Short: "Update one network interface of a cluster host",
		Long: `Update a cluster host network interface from a JSON or YAML configuration.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool
  host_id          The numeric ID of the cluster host
  interface_id     The numeric ID of the cluster host interface

Required Flags:
  --config-source  Source of the interface update configuration.
                   Values: 'pipe' for stdin input, or path to a JSON/YAML file.

Configuration Fields:
  status           Required. One of 'managed', 'unmanaged', 'inactive'.

Examples:
  # Take an interface under management
  echo '{"status":"managed"}' | metalcloud-cli vm-pool update-cluster-host-interface 123 456 789 --config-source pipe

  # Update from a file
  metalcloud-cli vm-pool update-cluster-host-interface 123 456 789 --config-source interface.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(vmPoolFlags.hostInterfaceConfigSource)
			if err != nil {
				return err
			}

			return vm_pool.VMPoolUpdateClusterHostInterface(cmd.Context(), args[0], args[1], args[2], config)
		},
	}

	vmPoolUpdateClusterHostInterfaceConfigExampleCmd = &cobra.Command{
		Use:   "update-cluster-host-interface-config-example",
		Short: "Display a cluster host interface update configuration example",
		Long: `Display a sample payload for 'vm-pool update-cluster-host-interface'.

Examples:
  # Display the example
  metalcloud-cli vm-pool update-cluster-host-interface-config-example

  # Save it for editing
  metalcloud-cli vm-pool update-cluster-host-interface-config-example > interface.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolUpdateClusterHostInterfaceConfigExample(cmd.Context())
		},
	}

	vmPoolHostInterfaceNetworkDeviceCmd = &cobra.Command{
		Use:     "cluster-host-interface-network-device [command]",
		Aliases: []string{"chi-network-device", "chind"},
		Short:   "Manage the network device assignments of a cluster host interface",
		Long: `Manage the links between a VM pool cluster host interface and the interfaces of
the network devices (switches) it is cabled to.

Every command takes the VM pool, cluster host and cluster host interface as
positional arguments; 'get' and 'remove' additionally take the numeric ID of the
assignment itself (not the network device ID).

Available Commands:
  list            List the network device assignments of an interface
  get             Get one network device assignment
  add             Link an interface to a network device interface
  remove          Delete one network device assignment
  config-example  Display an example 'add' payload`,
	}

	vmPoolHostInterfaceNetworkDeviceListCmd = &cobra.Command{
		Use:     "list vm_pool_id host_id interface_id",
		Aliases: []string{"ls"},
		Short:   "List the network device assignments of a cluster host interface",
		Long: `List every network device (switch) interface a VM pool cluster host interface is
linked to.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool
  host_id          The numeric ID of the cluster host
  interface_id     The numeric ID of the cluster host interface

Examples:
  # List the assignments of interface 789
  metalcloud-cli vm-pool cluster-host-interface-network-device list 123 456 789

  # Using the group alias
  metalcloud-cli vm-pool chind ls 123 456 789`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_READ},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolGetClusterHostInterfaceNetworkDevices(cmd.Context(), args[0], args[1], args[2])
		},
	}

	vmPoolHostInterfaceNetworkDeviceGetCmd = &cobra.Command{
		Use:     "get vm_pool_id host_id interface_id assignment_id",
		Aliases: []string{"show"},
		Short:   "Get one network device assignment of a cluster host interface",
		Long: `Get one link between a VM pool cluster host interface and a network device
(switch) interface.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool
  host_id          The numeric ID of the cluster host
  interface_id     The numeric ID of the cluster host interface
  assignment_id    The numeric ID of the network device assignment

Examples:
  # Get assignment 5 of interface 789
  metalcloud-cli vm-pool cluster-host-interface-network-device get 123 456 789 5`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_READ},
		Args:         cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolGetClusterHostInterfaceNetworkDevice(cmd.Context(), args[0], args[1], args[2], args[3])
		},
	}

	vmPoolHostInterfaceNetworkDeviceAddCmd = &cobra.Command{
		Use:     "add vm_pool_id host_id interface_id",
		Aliases: []string{"create", "new"},
		Short:   "Link a cluster host interface to a network device interface",
		Long: `Link a VM pool cluster host interface to an interface of a network device
(switch). The network device must be active and in a leaf position.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool
  host_id          The numeric ID of the cluster host
  interface_id     The numeric ID of the cluster host interface

Required Flags (either the pair of flags, or --config-source):
  --network-device            ID or label of the network device (switch)
  --network-device-interface  Name of the interface on the network device
  --config-source             Source of the assignment configuration.
                              Values: 'pipe' for stdin input, or path to a
                              JSON/YAML file. Mutually exclusive with the flags
                              above; the file carries a numeric networkDeviceId.

Examples:
  # Link by network device label
  metalcloud-cli vm-pool chind add 123 456 789 --network-device leaf-su00-r0 --network-device-interface Ethernet1/1

  # Link by network device ID
  metalcloud-cli vm-pool chind add 123 456 789 --network-device 42 --network-device-interface Ethernet1/1

  # Link from a configuration file
  metalcloud-cli vm-pool chind add 123 456 789 --config-source assignment.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_WRITE},
		Args:         cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var config []byte
			if vmPoolFlags.networkDeviceConfigSource != "" {
				var err error
				config, err = utils.ReadConfigFromPipeOrFile(vmPoolFlags.networkDeviceConfigSource)
				if err != nil {
					return err
				}
			}

			return vm_pool.VMPoolAddClusterHostInterfaceNetworkDevice(cmd.Context(), args[0], args[1], args[2],
				vmPoolFlags.networkDeviceRef, vmPoolFlags.networkDeviceInterfaceNameFlag, config)
		},
	}

	vmPoolHostInterfaceNetworkDeviceRemoveCmd = &cobra.Command{
		Use:     "remove vm_pool_id host_id interface_id assignment_id",
		Aliases: []string{"delete", "rm"},
		Short:   "Delete one network device assignment of a cluster host interface",
		Long: `Delete the link between a VM pool cluster host interface and a network device
(switch) interface.

Required Arguments:
  vm_pool_id       The numeric ID of the VM pool
  host_id          The numeric ID of the cluster host
  interface_id     The numeric ID of the cluster host interface
  assignment_id    The numeric ID of the network device assignment

Examples:
  # Remove assignment 5 of interface 789
  metalcloud-cli vm-pool cluster-host-interface-network-device remove 123 456 789 5`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_WRITE},
		Args:         cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolRemoveClusterHostInterfaceNetworkDevice(cmd.Context(), args[0], args[1], args[2], args[3])
		},
	}

	vmPoolHostInterfaceNetworkDeviceConfigExampleCmd = &cobra.Command{
		Use:   "config-example",
		Short: "Display a network device assignment configuration example",
		Long: `Display a sample payload for
'vm-pool cluster-host-interface-network-device add --config-source'.

Examples:
  # Display the example
  metalcloud-cli vm-pool chind config-example

  # Save it for editing
  metalcloud-cli vm-pool chind config-example > assignment.json`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolAddClusterHostInterfaceNetworkDeviceConfigExample(cmd.Context())
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

	// Import VMs command
	vmPoolCmd.AddCommand(vmPoolImportVMsCmd)
	vmPoolImportVMsCmd.Flags().StringVar(&vmPoolFlags.configSource, "config-source", "", "Source of the import configuration. Can be 'pipe' or path to a JSON file.")
	vmPoolImportVMsCmd.Flags().StringVar(&vmPoolFlags.infrastructureId, "infrastructure-id", "", "Infrastructure ID to associate with the imported VMs")
	vmPoolImportVMsCmd.Flags().StringVar(&vmPoolFlags.vmNames, "vm-names", "", "Comma-separated list of VM names to import")
	vmPoolImportVMsCmd.MarkFlagsOneRequired("config-source", "vm-names")
	vmPoolImportVMsCmd.MarkFlagsMutuallyExclusive("config-source", "vm-names")
	vmPoolImportVMsCmd.MarkFlagsMutuallyExclusive("config-source", "infrastructure-id")

	// Update command
	vmPoolCmd.AddCommand(vmPoolUpdateCmd)
	vmPoolUpdateCmd.Flags().StringVar(&vmPoolFlags.updateConfigSource, "config-source", "", "Source of the VM pool update configuration. Can be 'pipe' or path to a JSON file.")
	vmPoolUpdateCmd.MarkFlagsOneRequired("config-source")

	// Update config example command
	vmPoolCmd.AddCommand(vmPoolUpdateConfigExampleCmd)

	// Sync command
	vmPoolCmd.AddCommand(vmPoolSyncCmd)

	// Refresh command
	vmPoolCmd.AddCommand(vmPoolRefreshCmd)

	// Statistics command
	vmPoolCmd.AddCommand(vmPoolStatisticsCmd)

	// Containers command
	vmPoolCmd.AddCommand(vmPoolContainersCmd)

	// Get cluster host command
	vmPoolCmd.AddCommand(vmPoolGetClusterHostCmd)

	// Update cluster host command
	vmPoolCmd.AddCommand(vmPoolUpdateClusterHostCmd)
	vmPoolUpdateClusterHostCmd.Flags().StringVar(&vmPoolFlags.clusterHostConfigSource, "config-source", "", "Source of the cluster host update configuration. Can be 'pipe' or path to a JSON file.")
	vmPoolUpdateClusterHostCmd.MarkFlagsOneRequired("config-source")

	vmPoolCmd.AddCommand(vmPoolUpdateClusterHostConfigExampleCmd)

	// Cluster host statistics command
	vmPoolCmd.AddCommand(vmPoolClusterHostStatisticsCmd)

	// Cluster host containers command
	vmPoolCmd.AddCommand(vmPoolClusterHostContainersCmd)

	// Get cluster host interface command
	vmPoolCmd.AddCommand(vmPoolGetClusterHostInterfaceCmd)

	// Update cluster host interface command
	vmPoolCmd.AddCommand(vmPoolUpdateClusterHostInterfaceCmd)
	vmPoolUpdateClusterHostInterfaceCmd.Flags().StringVar(&vmPoolFlags.hostInterfaceConfigSource, "config-source", "", "Source of the cluster host interface update configuration. Can be 'pipe' or path to a JSON file.")
	vmPoolUpdateClusterHostInterfaceCmd.MarkFlagsOneRequired("config-source")

	vmPoolCmd.AddCommand(vmPoolUpdateClusterHostInterfaceConfigExampleCmd)

	// Cluster host interface network device sub-group
	vmPoolCmd.AddCommand(vmPoolHostInterfaceNetworkDeviceCmd)

	vmPoolHostInterfaceNetworkDeviceCmd.AddCommand(vmPoolHostInterfaceNetworkDeviceListCmd)

	vmPoolHostInterfaceNetworkDeviceCmd.AddCommand(vmPoolHostInterfaceNetworkDeviceGetCmd)

	vmPoolHostInterfaceNetworkDeviceCmd.AddCommand(vmPoolHostInterfaceNetworkDeviceAddCmd)
	vmPoolHostInterfaceNetworkDeviceAddCmd.Flags().StringVar(&vmPoolFlags.networkDeviceConfigSource, "config-source", "", "Source of the network device assignment configuration. Can be 'pipe' or path to a JSON file.")
	vmPoolHostInterfaceNetworkDeviceAddCmd.Flags().StringVar(&vmPoolFlags.networkDeviceRef, "network-device", "", "ID or label of the network device (switch) to link the interface to.")
	vmPoolHostInterfaceNetworkDeviceAddCmd.Flags().StringVar(&vmPoolFlags.networkDeviceInterfaceNameFlag, "network-device-interface", "", "Name of the interface on the network device.")
	vmPoolHostInterfaceNetworkDeviceAddCmd.MarkFlagsOneRequired("config-source", "network-device")
	vmPoolHostInterfaceNetworkDeviceAddCmd.MarkFlagsRequiredTogether("network-device", "network-device-interface")
	vmPoolHostInterfaceNetworkDeviceAddCmd.MarkFlagsMutuallyExclusive("config-source", "network-device")
	vmPoolHostInterfaceNetworkDeviceAddCmd.MarkFlagsMutuallyExclusive("config-source", "network-device-interface")

	vmPoolHostInterfaceNetworkDeviceCmd.AddCommand(vmPoolHostInterfaceNetworkDeviceRemoveCmd)

	vmPoolHostInterfaceNetworkDeviceCmd.AddCommand(vmPoolHostInterfaceNetworkDeviceConfigExampleCmd)
}
