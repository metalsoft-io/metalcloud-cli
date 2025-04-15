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
		filterType   string
		configSource string
		limit        string
		page         string
	}{}

	vmPoolCmd = &cobra.Command{
		Use:     "vm-pool [command]",
		Aliases: []string{"vmp"},
		Short:   "VM Pool management",
		Long:    `VM Pool management commands.`,
	}

	vmPoolListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all VM pools.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolList(cmd.Context(), vmPoolFlags.filterType)
		},
	}

	vmPoolGetCmd = &cobra.Command{
		Use:          "get vm_pool_id",
		Aliases:      []string{"show"},
		Short:        "Get VM pool details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolGet(cmd.Context(), args[0])
		},
	}

	vmPoolConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Get VM pool configuration example.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolConfigExample(cmd.Context())
		},
	}

	vmPoolCreateCmd = &cobra.Command{
		Use:          "create",
		Short:        "Create a new VM pool.",
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
		Use:          "delete vm_pool_id",
		Short:        "Delete a VM pool.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolDelete(cmd.Context(), args[0])
		},
	}

	vmPoolGetCredentialsCmd = &cobra.Command{
		Use:          "credentials vm_pool_id",
		Short:        "Get VM pool credentials.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_POOLS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_pool.VMPoolGetCredentials(cmd.Context(), args[0])
		},
	}

	vmPoolGetVMsCmd = &cobra.Command{
		Use:          "vms vm_pool_id",
		Short:        "Get VMs for a VM pool.",
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
		Use:          "cluster-hosts vm_pool_id",
		Short:        "Get cluster hosts for a VM pool.",
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
		Use:          "cluster-host-vms vm_pool_id host_id",
		Short:        "Get VMs for a cluster host in a VM pool.",
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
		Use:          "cluster-host-interfaces vm_pool_id host_id",
		Short:        "Get interfaces for a cluster host in a VM pool.",
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
	vmPoolListCmd.Flags().StringVar(&vmPoolFlags.filterType, "filter-type", "", "Filter the result by VM pool type (e.g., vmware, hyperv, etc).")

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
