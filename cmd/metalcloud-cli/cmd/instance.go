package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/instance"
	"github.com/spf13/cobra"
)

// Instance Group management commands.
var (
	instanceGroupFlags = struct {
		label         string
		instanceCount int
		osTemplateId  int
	}{}

	instanceGroupCmd = &cobra.Command{
		Use:     "instance-group [command]",
		Aliases: []string{"ig", "instance-array", "ia"},
		Short:   "Instance Group management",
		Long:    `Instance Group management commands.`,
	}

	instanceGroupListCmd = &cobra.Command{
		Use:          "list infrastructure_id_or_label",
		Aliases:      []string{"ls"},
		Short:        "List all instance groups in an infrastructures.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_READ}, // TODO: Use specific permission
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return instance.InstanceGroupList(cmd.Context(), args[0])
		},
	}

	instanceGroupGetCmd = &cobra.Command{
		Use:          "get instance_group_id",
		Aliases:      []string{"show"},
		Short:        "Get instance group details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_READ}, // TODO: Use specific permission
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return instance.InstanceGroupGet(cmd.Context(), args[0])
		},
	}

	instanceGroupCreateCmd = &cobra.Command{
		Use:          "create infrastructure_id_or_label label server_type_id instance_count [os_template_id]",
		Aliases:      []string{"new"},
		Short:        "Create new instance group in an infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_READ}, // TODO: Use specific permission
		Args:         cobra.RangeArgs(4, 5),
		RunE: func(cmd *cobra.Command, args []string) error {
			os_template_id := ""
			if len(args) == 5 {
				os_template_id = args[4]
			}

			return instance.InstanceGroupCreate(cmd.Context(), args[0], args[1], args[2], args[3], os_template_id)
		},
	}

	instanceGroupUpdateCmd = &cobra.Command{
		Use:          "update instance_group_id",
		Aliases:      []string{"edit"},
		Short:        "Update instance group configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_READ}, // TODO: Use specific permission
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return instance.InstanceGroupUpdate(cmd.Context(), args[0], instanceGroupFlags.label, instanceGroupFlags.instanceCount, instanceGroupFlags.osTemplateId)
		},
	}

	instanceGroupDeleteCmd = &cobra.Command{
		Use:          "delete instance_group_id",
		Aliases:      []string{"rm"},
		Short:        "Update instance group configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_READ}, // TODO: Use specific permission
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return instance.InstanceGroupDelete(cmd.Context(), args[0])
		},
	}

	instanceGroupInstancesCmd = &cobra.Command{
		Use:          "instances instance_group_id",
		Aliases:      []string{"instances-list", "instances-ls"},
		Short:        "List instance group instances.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_READ}, // TODO: Use specific permission
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return instance.InstanceGroupInstances(cmd.Context(), args[0])
		},
	}
)

// Instance management commands.
var (
	instanceCmd = &cobra.Command{
		Use:     "instance [command]",
		Aliases: []string{"inst"},
		Short:   "Instance management",
		Long:    `Instance management commands.`,
	}

	instanceGetCmd = &cobra.Command{
		Use:          "get instance_id",
		Aliases:      []string{"show"},
		Short:        "Get instance details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_READ}, // TODO: Use specific permission
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return instance.InstanceGet(cmd.Context(), args[0])
		},
	}
)

func init() {
	// Instance Group management commands.
	rootCmd.AddCommand(instanceGroupCmd)

	instanceGroupCmd.AddCommand(instanceGroupListCmd)

	instanceGroupCmd.AddCommand(instanceGroupGetCmd)

	instanceGroupCmd.AddCommand(instanceGroupCreateCmd)

	instanceGroupCmd.AddCommand(instanceGroupUpdateCmd)
	instanceGroupUpdateCmd.Flags().StringVar(&instanceGroupFlags.label, "label", "", "Set the instance group label.")
	instanceGroupUpdateCmd.Flags().IntVar(&instanceGroupFlags.instanceCount, "instance-count", 0, "Set the count of instance group instances.")
	instanceGroupUpdateCmd.Flags().IntVar(&instanceGroupFlags.osTemplateId, "os-template-id", 0, "Set the instance group OS template Id.")
	instanceGroupUpdateCmd.MarkFlagsOneRequired("label", "instance-count", "os-template-id")

	instanceGroupCmd.AddCommand(instanceGroupDeleteCmd)

	instanceGroupCmd.AddCommand(instanceGroupInstancesCmd)

	// Instance management commands.
	rootCmd.AddCommand(instanceCmd)

	instanceCmd.AddCommand(instanceGetCmd)
}
