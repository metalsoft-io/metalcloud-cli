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
		Short:   "VM Type management",
		Long:    `VM Type management commands.`,
	}

	vmTypeListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all VM types.",
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
		Use:          "get vm_type_id",
		Aliases:      []string{"show"},
		Short:        "Get VM type details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_TYPES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_type.VMTypeGet(cmd.Context(), args[0])
		},
	}

	vmTypeConfigExampleCmd = &cobra.Command{
		Use:          "config-example",
		Short:        "Get VM type configuration example.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_TYPES_WRITE},
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_type.VMTypeConfigExample(cmd.Context())
		},
	}

	vmTypeCreateCmd = &cobra.Command{
		Use:          "create",
		Short:        "Create a new VM type.",
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
		Use:          "update vm_type_id",
		Short:        "Update a VM type.",
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
		Use:          "delete vm_type_id",
		Short:        "Delete a VM type.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VM_TYPES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_type.VMTypeDelete(cmd.Context(), args[0])
		},
	}

	vmTypeGetVMsCmd = &cobra.Command{
		Use:          "vms vm_type_id",
		Short:        "Get VMs for a VM type.",
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
