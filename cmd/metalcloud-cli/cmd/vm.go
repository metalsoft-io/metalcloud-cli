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
		Short:   "VM management",
		Long:    `Virtual Machine management commands.`,
	}

	vmGetCmd = &cobra.Command{
		Use:          "get vm_id",
		Aliases:      []string{"show"},
		Short:        "Get VM info.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VMS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm.VMGet(cmd.Context(), args[0])
		},
	}

	vmPowerStatusCmd = &cobra.Command{
		Use:          "power-status vm_id",
		Short:        "Get VM power status.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VMS_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm.VMPowerStatus(cmd.Context(), args[0])
		},
	}

	vmStartCmd = &cobra.Command{
		Use:          "start vm_id",
		Short:        "Start a VM.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VMS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm.VMStart(cmd.Context(), args[0])
		},
	}

	vmShutdownCmd = &cobra.Command{
		Use:          "shutdown vm_id",
		Aliases:      []string{"stop"},
		Short:        "Shutdown a VM.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VMS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm.VMShutdown(cmd.Context(), args[0])
		},
	}

	vmRebootCmd = &cobra.Command{
		Use:          "reboot vm_id",
		Short:        "Reboot a VM.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_VMS_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm.VMReboot(cmd.Context(), args[0])
		},
	}

	vmUpdateCmd = &cobra.Command{
		Use:          "update vm_id",
		Short:        "Update VM information.",
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
		Use:          "console-info vm_id",
		Short:        "Get VM remote console information.",
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
