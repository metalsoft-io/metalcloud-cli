package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/instance"
	"github.com/spf13/cobra"
)

// Instance Group management commands.
var (
	instanceGroupCmd = &cobra.Command{
		Use:     "instance-group [command]",
		Aliases: []string{"ig"},
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
)

func init() {
	rootCmd.AddCommand(instanceGroupCmd)

	instanceGroupCmd.AddCommand(instanceGroupListCmd)

	instanceGroupCmd.AddCommand(instanceGroupGetCmd)

	instanceGroupCmd.AddCommand(instanceGroupCreateCmd)
}
