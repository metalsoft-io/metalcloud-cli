package cmd

import (
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
}
