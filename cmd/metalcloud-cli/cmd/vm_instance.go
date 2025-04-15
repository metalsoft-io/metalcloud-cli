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
		Short:   "VM Instance Group management",
		Long:    `VM Instance Group management commands.`,
	}

	vmInstanceGroupListCmd = &cobra.Command{
		Use:          "list infrastructure_id",
		Aliases:      []string{"ls"},
		Short:        "List all VM instance groups in an infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGroupList(cmd.Context(), args[0])
		},
	}

	vmInstanceGroupGetCmd = &cobra.Command{
		Use:          "get infrastructure_id vm_instance_group_id",
		Aliases:      []string{"show"},
		Short:        "Get VM instance group details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGroupGet(cmd.Context(), args[0], args[1])
		},
	}

	vmInstanceGroupCreateCmd = &cobra.Command{
		Use:          "create infrastructure_id vm_type_id disk_size_gb instance_count [os_template_id]",
		Aliases:      []string{"new"},
		Short:        "Create new VM instance group in an infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
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
		Use:          "update infrastructure_id vm_instance_group_id",
		Aliases:      []string{"edit"},
		Short:        "Update VM instance group configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
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
		Use:          "delete infrastructure_id vm_instance_group_id",
		Aliases:      []string{"rm"},
		Short:        "Delete VM instance group.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGroupDelete(cmd.Context(), args[0], args[1])
		},
	}

	vmInstanceGroupInstancesCmd = &cobra.Command{
		Use:          "instances infrastructure_id vm_instance_group_id",
		Aliases:      []string{"instances-list", "instances-ls"},
		Short:        "List VM instance group instances.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
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
		Short:   "VM Instance management",
		Long:    `VM Instance management commands.`,
	}

	vmInstanceGetCmd = &cobra.Command{
		Use:          "get infrastructure_id vm_instance_id",
		Aliases:      []string{"show"},
		Short:        "Get VM instance details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGet(cmd.Context(), args[0], args[1])
		},
	}

	vmInstanceListCmd = &cobra.Command{
		Use:          "list infrastructure_id",
		Aliases:      []string{"ls"},
		Short:        "List all VM instances in an infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceList(cmd.Context(), args[0])
		},
	}

	vmInstanceGetConfigCmd = &cobra.Command{
		Use:          "config infrastructure_id vm_instance_id",
		Aliases:      []string{"get-config"},
		Short:        "Get VM instance configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstanceGetConfig(cmd.Context(), args[0], args[1])
		},
	}

	vmInstanceStartCmd = &cobra.Command{
		Use:          "start infrastructure_id vm_instance_id",
		Short:        "Start VM instance.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstancePowerControl(cmd.Context(), args[0], args[1], "start")
		},
	}

	vmInstanceShutdownCmd = &cobra.Command{
		Use:          "shutdown infrastructure_id vm_instance_id",
		Short:        "Shutdown VM instance.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstancePowerControl(cmd.Context(), args[0], args[1], "shutdown")
		},
	}

	vmInstanceRebootCmd = &cobra.Command{
		Use:          "reboot infrastructure_id vm_instance_id",
		Short:        "Reboot VM instance.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm_instance.VMInstancePowerControl(cmd.Context(), args[0], args[1], "reboot")
		},
	}

	vmInstancePowerStatusCmd = &cobra.Command{
		Use:          "power-status infrastructure_id vm_instance_id",
		Short:        "Get VM instance power status.",
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
