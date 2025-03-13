package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/instance"
	"github.com/spf13/cobra"
)

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
)

func init() {
	rootCmd.AddCommand(instanceGroupCmd)

	instanceGroupCmd.AddCommand(instanceGroupListCmd)
}
