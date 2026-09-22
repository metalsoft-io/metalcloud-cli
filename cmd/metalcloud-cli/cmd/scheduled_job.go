package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/job"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	scheduledJobFlags = struct {
		configSource string
	}{}

	scheduledJobCmd = &cobra.Command{
		Use:     "scheduled-job [command]",
		Aliases: []string{"cron", "cronjob"},
		Short:   "Manage scheduled jobs",
		Long: `Manage scheduled jobs in MetalCloud.

Scheduled jobs allow you to schedule recurring operations that are executed
automatically on a defined schedule.

Available Commands:
  list    List all scheduled jobs
  get     Get detailed information about a specific scheduled job
  create  Create a new scheduled job
  update  Update an existing scheduled job
  delete  Delete a scheduled job

Use "metalcloud-cli scheduled-job [command] --help" for more information about a command.`,
	}

	scheduledJobListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all scheduled jobs",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_JOB_QUEUE_READ},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return job.ScheduledJobList(cmd.Context())
		},
	}

	scheduledJobGetCmd = &cobra.Command{
		Use:          "get scheduled_job_id",
		Aliases:      []string{"show"},
		Short:        "Get detailed information about a specific scheduled job",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_JOB_QUEUE_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return job.ScheduledJobGet(cmd.Context(), args[0])
		},
	}

	scheduledJobCreateCmd = &cobra.Command{
		Use:          "create",
		Short:        "Create a new scheduled job from configuration",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_JOB_QUEUE_WRITE},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(scheduledJobFlags.configSource)
			if err != nil {
				return err
			}
			return job.ScheduledJobCreate(cmd.Context(), config)
		},
	}

	scheduledJobUpdateCmd = &cobra.Command{
		Use:          "update scheduled_job_id",
		Short:        "Update an existing scheduled job",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_JOB_QUEUE_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(scheduledJobFlags.configSource)
			if err != nil {
				return err
			}
			return job.ScheduledJobUpdate(cmd.Context(), args[0], config)
		},
	}

	scheduledJobDeleteCmd = &cobra.Command{
		Use:          "delete scheduled_job_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a scheduled job",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_JOB_QUEUE_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return job.ScheduledJobDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(scheduledJobCmd)

	scheduledJobCmd.AddCommand(scheduledJobListCmd)
	scheduledJobCmd.AddCommand(scheduledJobGetCmd)

	scheduledJobCmd.AddCommand(scheduledJobCreateCmd)
	scheduledJobCreateCmd.Flags().StringVar(&scheduledJobFlags.configSource, "config-source", "", "Path to JSON config file or 'pipe' for stdin.")
	scheduledJobCreateCmd.MarkFlagsOneRequired("config-source")

	scheduledJobCmd.AddCommand(scheduledJobUpdateCmd)
	scheduledJobUpdateCmd.Flags().StringVar(&scheduledJobFlags.configSource, "config-source", "", "Path to JSON config file or 'pipe' for stdin.")
	scheduledJobUpdateCmd.MarkFlagsOneRequired("config-source")

	scheduledJobCmd.AddCommand(scheduledJobDeleteCmd)
}
