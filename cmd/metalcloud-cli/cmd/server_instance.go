package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/server_instance"
	"github.com/spf13/cobra"
)

// Server Instance Group management commands.
var (
	serverInstanceGroupFlags = struct {
		label         string
		instanceCount int
		osTemplateId  int
	}{}

	serverInstanceGroupCmd = &cobra.Command{
		Use:     "server-instance-group [command]",
		Aliases: []string{"ig", "instance-array", "ia"},
		Short:   "Server Instance Group management",
		Long:    `Server Instance Group management commands.`,
	}

	serverInstanceGroupListCmd = &cobra.Command{
		Use:          "list infrastructure_id_or_label",
		Aliases:      []string{"ls"},
		Short:        "List all server instance groups in an infrastructures.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_READ}, // TODO: Use specific permission
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupList(cmd.Context(), args[0])
		},
	}

	serverInstanceGroupGetCmd = &cobra.Command{
		Use:          "get server_instance_group_id",
		Aliases:      []string{"show"},
		Short:        "Get server instance group details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_READ}, // TODO: Use specific permission
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupGet(cmd.Context(), args[0])
		},
	}

	serverInstanceGroupCreateCmd = &cobra.Command{
		Use:          "create infrastructure_id_or_label label server_type_id instance_count [os_template_id]",
		Aliases:      []string{"new"},
		Short:        "Create new server instance group in an infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_READ}, // TODO: Use specific permission
		Args:         cobra.RangeArgs(4, 5),
		RunE: func(cmd *cobra.Command, args []string) error {
			os_template_id := ""
			if len(args) == 5 {
				os_template_id = args[4]
			}

			return server_instance.ServerInstanceGroupCreate(cmd.Context(), args[0], args[1], args[2], args[3], os_template_id)
		},
	}

	serverInstanceGroupUpdateCmd = &cobra.Command{
		Use:          "update server_instance_group_id",
		Aliases:      []string{"edit"},
		Short:        "Update server instance group configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_READ}, // TODO: Use specific permission
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupUpdate(cmd.Context(), args[0], serverInstanceGroupFlags.label, serverInstanceGroupFlags.instanceCount, serverInstanceGroupFlags.osTemplateId)
		},
	}

	serverInstanceGroupDeleteCmd = &cobra.Command{
		Use:          "delete server_instance_group_id",
		Aliases:      []string{"rm"},
		Short:        "Update server instance group configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_READ}, // TODO: Use specific permission
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupDelete(cmd.Context(), args[0])
		},
	}

	serverInstanceGroupInstancesCmd = &cobra.Command{
		Use:          "instances server_instance_group_id",
		Aliases:      []string{"instances-list", "instances-ls"},
		Short:        "List server instance group instances.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_READ}, // TODO: Use specific permission
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGroupInstances(cmd.Context(), args[0])
		},
	}
)

// Server Instance management commands.
var (
	serverInstanceCmd = &cobra.Command{
		Use:     "server-instance [command]",
		Aliases: []string{"inst"},
		Short:   "Server Instance management",
		Long:    `Server Instance management commands.`,
	}

	serverInstanceGetCmd = &cobra.Command{
		Use:          "get server_instance_id",
		Aliases:      []string{"show"},
		Short:        "Get server instance details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.SITE_READ}, // TODO: Use specific permission
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return server_instance.ServerInstanceGet(cmd.Context(), args[0])
		},
	}
)

func init() {
	// Server Instance Group management commands.
	rootCmd.AddCommand(serverInstanceGroupCmd)

	serverInstanceGroupCmd.AddCommand(serverInstanceGroupListCmd)

	serverInstanceGroupCmd.AddCommand(serverInstanceGroupGetCmd)

	serverInstanceGroupCmd.AddCommand(serverInstanceGroupCreateCmd)

	serverInstanceGroupCmd.AddCommand(serverInstanceGroupUpdateCmd)
	serverInstanceGroupUpdateCmd.Flags().StringVar(&serverInstanceGroupFlags.label, "label", "", "Set the instance group label.")
	serverInstanceGroupUpdateCmd.Flags().IntVar(&serverInstanceGroupFlags.instanceCount, "instance-count", 0, "Set the count of instance group instances.")
	serverInstanceGroupUpdateCmd.Flags().IntVar(&serverInstanceGroupFlags.osTemplateId, "os-template-id", 0, "Set the instance group OS template Id.")
	serverInstanceGroupUpdateCmd.MarkFlagsOneRequired("label", "instance-count", "os-template-id")

	serverInstanceGroupCmd.AddCommand(serverInstanceGroupDeleteCmd)

	serverInstanceGroupCmd.AddCommand(serverInstanceGroupInstancesCmd)

	// Server Instance management commands.
	rootCmd.AddCommand(serverInstanceCmd)

	serverInstanceCmd.AddCommand(serverInstanceGetCmd)
}
