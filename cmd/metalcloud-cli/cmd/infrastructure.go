package cmd

import (
	"time"

	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/infrastructure"
	"github.com/spf13/cobra"
)

var (
	infrastructureFlags = struct {
		showAll             bool
		showOrdered         bool
		showDeleted         bool
		customVariables     string
		allowDataLoss       bool
		attemptSoftShutdown bool
		attemptHardShutdown bool
		softShutdownTimeout int
		forceShutdown       bool
		userId              int
		startTime           time.Time
		endTime             time.Time
		siteIds             []int
		infrastructureIds   []int
	}{}

	infrastructureCmd = &cobra.Command{
		Use:     "infrastructure [command]",
		Aliases: []string{"infra"},
		Short:   "Infrastructure management",
		Long:    `Infrastructure management commands.`,
	}

	infrastructureListCmd = &cobra.Command{
		Use:          "list [flags...]",
		Aliases:      []string{"ls"},
		Short:        "List all infrastructures.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureList(cmd.Context(),
				infrastructureFlags.showAll,
				infrastructureFlags.showOrdered,
				infrastructureFlags.showDeleted)
		},
	}

	infrastructureGetCmd = &cobra.Command{
		Use:          "get infrastructure_id_or_label",
		Aliases:      []string{"show"},
		Short:        "Get infrastructure details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureGet(cmd.Context(), args[0])
		},
	}

	infrastructureCreateCmd = &cobra.Command{
		Use:          "create site_id label",
		Aliases:      []string{"new"},
		Short:        "Create new infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureCreate(cmd.Context(), args[0], args[1])
		},
	}

	infrastructureUpdateCmd = &cobra.Command{
		Use:          "update infrastructure_id_or_label [new_label]",
		Aliases:      []string{"edit"},
		Short:        "Update infrastructure configuration.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			label := ""
			if len(args) > 1 {
				label = args[1]
			}

			return infrastructure.InfrastructureUpdate(cmd.Context(), args[0], label,
				infrastructureFlags.customVariables)
		},
	}

	infrastructureDeleteCmd = &cobra.Command{
		Use:          "delete infrastructure_id_or_label",
		Aliases:      []string{"rm"},
		Short:        "Delete infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureDelete(cmd.Context(), args[0])
		},
	}

	infrastructureDeployCmd = &cobra.Command{
		Use:          "deploy infrastructure_id_or_label",
		Aliases:      []string{"apply"},
		Short:        "Deploy infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureDeploy(cmd.Context(), args[0],
				infrastructureFlags.allowDataLoss,
				infrastructureFlags.attemptSoftShutdown,
				infrastructureFlags.attemptHardShutdown,
				infrastructureFlags.softShutdownTimeout,
				infrastructureFlags.forceShutdown)
		},
	}

	infrastructureRevertCmd = &cobra.Command{
		Use:          "revert infrastructure_id_or_label",
		Aliases:      []string{"undo"},
		Short:        "Revert infrastructure changes.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureRevert(cmd.Context(), args[0])
		},
	}

	infrastructureGetUsersCmd = &cobra.Command{
		Use:          "users infrastructure_id_or_label",
		Aliases:      []string{"list-users", "get-users"},
		Short:        "Get users for an infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureGetUsers(cmd.Context(), args[0])
		},
	}

	infrastructureAddUserCmd = &cobra.Command{
		Use:          "add-user infrastructure_id_or_label user_email [create_if_not_exists]",
		Short:        "Add a user to an infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			createMissing := "false"
			if len(args) == 3 {
				createMissing = args[2]
			}
			return infrastructure.InfrastructureAddUser(cmd.Context(), args[0], args[1], createMissing)
		},
	}

	infrastructureRemoveUserCmd = &cobra.Command{
		Use:          "remove-user infrastructure_id_or_label user_id",
		Aliases:      []string{"delete-user"},
		Short:        "Remove a user from an infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_WRITE},
		Args:         cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureRemoveUser(cmd.Context(), args[0], args[1])
		},
	}

	infrastructureGetUserLimitsCmd = &cobra.Command{
		Use:          "user-limits infrastructure_id_or_label",
		Aliases:      []string{"get-user-limits"},
		Short:        "Get user limits for an infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureGetUserLimits(cmd.Context(), args[0])
		},
	}

	infrastructureGetStatisticsCmd = &cobra.Command{
		Use:          "statistics infrastructure_id_or_label",
		Aliases:      []string{"stats", "get-statistics"},
		Short:        "Get statistics for an infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureGetStatistics(cmd.Context(), args[0])
		},
	}

	infrastructureGetConfigInfoCmd = &cobra.Command{
		Use:          "config-info infrastructure_id_or_label",
		Aliases:      []string{"get-config-info"},
		Short:        "Get configuration information for an infrastructure.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureGetConfigInfo(cmd.Context(), args[0])
		},
	}

	infrastructureGetAllStatisticsCmd = &cobra.Command{
		Use:          "all-statistics",
		Aliases:      []string{"all-stats", "get-all-statistics"},
		Short:        "Get statistics for all infrastructures.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureGetAllStatistics(cmd.Context())
		},
	}

	infrastructureUtilizationCmd = &cobra.Command{
		Use:     "utilization",
		Aliases: []string{"get-utilization"},
		Short:   "Get resource utilization report for infrastructures.",
		Long: `Get detailed utilization report for infrastructure resources within a specified time range. 
The report provides insights into resource usage patterns and capacity planning for infrastructures.

Required flags:
  --user-id       ID of the user to include in the report
  --start-time    Start time for the report (RFC3339 or date format)
  --end-time      End time for the report (RFC3339 or date format)

Optional flags:
  --site-id            Site IDs to include in the report (can be specified multiple times)
  --infrastructure-id  Infrastructure IDs to include in the report (can be specified multiple times)

Examples:
  # Get utilization report for user 123 for the last 7 days
  metalcloud-cli infrastructure utilization --user-id 123 --start-time 2025-08-01 --end-time 2025-08-08

  # Get utilization for specific sites and infrastructures
  metalcloud-cli infrastructure utilization --user-id 123 --start-time 2025-08-01T00:00:00Z --end-time 2025-08-08T23:59:59Z --site-id 1 --site-id 2 --infrastructure-id 100 --infrastructure-id 101

  # Get utilization for all infrastructures of a user in specific sites
  metalcloud-cli infrastructure utilization --user-id 123 --start-time 2025-08-01 --end-time 2025-08-08 --site-id 1 --site-id 3`,
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_INFRASTRUCTURES_READ},
		RunE: func(cmd *cobra.Command, args []string) error {
			return infrastructure.InfrastructureGetUtilization(
				cmd.Context(),
				infrastructureFlags.userId,
				infrastructureFlags.startTime,
				infrastructureFlags.endTime,
				infrastructureFlags.siteIds,
				infrastructureFlags.infrastructureIds)
		},
	}
)

func init() {
	rootCmd.AddCommand(infrastructureCmd)

	infrastructureCmd.AddCommand(infrastructureListCmd)
	infrastructureListCmd.Flags().BoolVar(&infrastructureFlags.showAll, "show-all", false, "If set will return all infrastructures.")
	infrastructureListCmd.Flags().BoolVar(&infrastructureFlags.showOrdered, "show-ordered", false, "If set will also return ordered (created but not deployed) infrastructures.")
	infrastructureListCmd.Flags().BoolVar(&infrastructureFlags.showDeleted, "show-deleted", false, "If set will also return deleted infrastructures.")

	infrastructureCmd.AddCommand(infrastructureGetCmd)

	infrastructureCmd.AddCommand(infrastructureCreateCmd)

	infrastructureCmd.AddCommand(infrastructureUpdateCmd)
	infrastructureUpdateCmd.Flags().StringVar(&infrastructureFlags.customVariables, "custom-variables", "", "Set of infrastructure custom variables.")

	infrastructureCmd.AddCommand(infrastructureDeleteCmd)

	infrastructureCmd.AddCommand(infrastructureDeployCmd)
	infrastructureDeployCmd.Flags().BoolVar(&infrastructureFlags.allowDataLoss, "allow-data-loss", false, "If set, deploy will not throw error if data loss is expected.")
	infrastructureDeployCmd.Flags().BoolVar(&infrastructureFlags.attemptSoftShutdown, "attempt-soft-shutdown", true, "If set, attempt a soft (ACPI) power off of all the servers in the infrastructure before the deploy.")
	infrastructureDeployCmd.Flags().BoolVar(&infrastructureFlags.attemptHardShutdown, "attempt-hard-shutdown", true, "If set, force a hard power off after timeout expired and the server is not powered off")
	infrastructureDeployCmd.Flags().IntVar(&infrastructureFlags.softShutdownTimeout, "soft-shutdown-timeout", 180, "Timeout to wait for soft shutdown before forcing hard shutdown.")
	infrastructureDeployCmd.Flags().BoolVar(&infrastructureFlags.forceShutdown, "force-shutdown", false, "If set, deploy will force shutdown of all servers in the infrastructure.")

	infrastructureCmd.AddCommand(infrastructureRevertCmd)

	infrastructureCmd.AddCommand(infrastructureGetUsersCmd)

	infrastructureCmd.AddCommand(infrastructureAddUserCmd)

	infrastructureCmd.AddCommand(infrastructureRemoveUserCmd)

	infrastructureCmd.AddCommand(infrastructureGetUserLimitsCmd)

	infrastructureCmd.AddCommand(infrastructureGetStatisticsCmd)

	infrastructureCmd.AddCommand(infrastructureGetConfigInfoCmd)

	infrastructureCmd.AddCommand(infrastructureGetAllStatisticsCmd)

	infrastructureCmd.AddCommand(infrastructureUtilizationCmd)
	infrastructureUtilizationCmd.Flags().IntVar(&infrastructureFlags.userId, "user-id", 0, "ID of the user to include in the report.")
	infrastructureUtilizationCmd.Flags().TimeVar(&infrastructureFlags.startTime, "start-time", time.Now().Add(-time.Duration(time.Now().Day())), []string{time.RFC3339, time.DateOnly}, "Start time for the report.")
	infrastructureUtilizationCmd.Flags().TimeVar(&infrastructureFlags.endTime, "end-time", time.Now(), []string{time.RFC3339, time.DateOnly}, "End time for the report.")
	infrastructureUtilizationCmd.Flags().IntSliceVar(&infrastructureFlags.siteIds, "site-id", []int{}, "Site IDs to include in the report.")
	infrastructureUtilizationCmd.Flags().IntSliceVar(&infrastructureFlags.infrastructureIds, "infrastructure-id", []int{}, "Infrastructure IDs to include in the report.")
	infrastructureUtilizationCmd.MarkFlagRequired("user-id")
	infrastructureUtilizationCmd.MarkFlagRequired("start-time")
	infrastructureUtilizationCmd.MarkFlagRequired("end-time")
}
