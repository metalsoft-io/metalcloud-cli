package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/job"
	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	cronJobFlags = struct {
		configSource string
	}{}

	cronJobCmd = &cobra.Command{
		Use:     "cron-job [command]",
		Aliases: []string{"cron", "cronjob"},
		Short:   "Manage scheduled cron jobs",
		Long: `Manage scheduled cron jobs in MetalCloud.

Cron jobs allow you to schedule recurring operations that are executed
automatically on a defined schedule.

Available Commands:
  list    List all cron jobs
  get     Get detailed information about a specific cron job
  create  Create a new cron job
  update  Update an existing cron job
  delete  Delete a cron job

Use "metalcloud-cli cron-job [command] --help" for more information about a command.`,
	}

	cronJobListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List all cron jobs",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_JOB_QUEUE_READ},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return job.CronJobList(cmd.Context())
		},
	}

	cronJobGetCmd = &cobra.Command{
		Use:          "get cron_job_id",
		Aliases:      []string{"show"},
		Short:        "Get detailed information about a specific cron job",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_JOB_QUEUE_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return job.CronJobGet(cmd.Context(), args[0])
		},
	}

	cronJobCreateCmd = &cobra.Command{
		Use:          "create",
		Short:        "Create a new cron job from configuration",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_JOB_QUEUE_WRITE},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(cronJobFlags.configSource)
			if err != nil {
				return err
			}
			return job.CronJobCreate(cmd.Context(), config)
		},
	}

	cronJobUpdateCmd = &cobra.Command{
		Use:          "update cron_job_id",
		Short:        "Update an existing cron job",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_JOB_QUEUE_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := utils.ReadConfigFromPipeOrFile(cronJobFlags.configSource)
			if err != nil {
				return err
			}
			return job.CronJobUpdate(cmd.Context(), args[0], config)
		},
	}

	cronJobDeleteCmd = &cobra.Command{
		Use:          "delete cron_job_id",
		Aliases:      []string{"rm"},
		Short:        "Delete a cron job",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_JOB_QUEUE_WRITE},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return job.CronJobDelete(cmd.Context(), args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(cronJobCmd)

	cronJobCmd.AddCommand(cronJobListCmd)
	cronJobCmd.AddCommand(cronJobGetCmd)

	cronJobCmd.AddCommand(cronJobCreateCmd)
	cronJobCreateCmd.Flags().StringVar(&cronJobFlags.configSource, "config-source", "", "Path to JSON config file or 'pipe' for stdin.")
	cronJobCreateCmd.MarkFlagsOneRequired("config-source")

	cronJobCmd.AddCommand(cronJobUpdateCmd)
	cronJobUpdateCmd.Flags().StringVar(&cronJobFlags.configSource, "config-source", "", "Path to JSON config file or 'pipe' for stdin.")
	cronJobUpdateCmd.MarkFlagsOneRequired("config-source")

	cronJobCmd.AddCommand(cronJobDeleteCmd)
}
