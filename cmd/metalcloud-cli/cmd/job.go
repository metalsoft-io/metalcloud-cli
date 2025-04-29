package cmd

import (
	"github.com/metalsoft-io/metalcloud-cli/cmd/metalcloud-cli/system"
	"github.com/metalsoft-io/metalcloud-cli/internal/job"
	"github.com/spf13/cobra"
)

var (
	jobFlags = struct {
		filterJobId      []string
		filterStatus     []string
		filterJobGroupId []string
		sortBy           []string
	}{}

	jobGroupFlags = struct {
		filterJobGroupId []string
		filterType       []string
		sortBy           []string
	}{}

	jobCmd = &cobra.Command{
		Use:   "job [command]",
		Short: "Job management",
		Long:  `Job management commands.`,
	}

	jobListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List jobs.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_JOB_QUEUE_READ},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return job.JobList(cmd.Context(), job.ListFlags{
				FilterJobId:      jobFlags.filterJobId,
				FilterStatus:     jobFlags.filterStatus,
				FilterJobGroupId: jobFlags.filterJobGroupId,
				SortBy:           jobFlags.sortBy,
			})
		},
	}

	jobGetCmd = &cobra.Command{
		Use:          "get job_id",
		Aliases:      []string{"show"},
		Short:        "Get job details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_JOB_QUEUE_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return job.JobGet(cmd.Context(), args[0])
		},
	}
)

var (
	jobGroupCmd = &cobra.Command{
		Use:   "job-group [command]",
		Short: "Job group management",
		Long:  `Job group management commands.`,
	}

	jobGroupListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "List job groups.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_JOB_QUEUE_READ},
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return job.JobGroupList(cmd.Context(), job.GroupListFlags{
				FilterJobGroupId: jobGroupFlags.filterJobGroupId,
				FilterType:       jobGroupFlags.filterType,
				SortBy:           jobGroupFlags.sortBy,
			})
		},
	}

	jobGroupGetCmd = &cobra.Command{
		Use:          "get job_group_id",
		Aliases:      []string{"show"},
		Short:        "Get job group details.",
		SilenceUsage: true,
		Annotations:  map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_JOB_QUEUE_READ},
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return job.JobGroupGet(cmd.Context(), args[0])
		},
	}
)

func init() {
	// Job commands
	rootCmd.AddCommand(jobCmd)

	jobCmd.AddCommand(jobListCmd)
	jobListCmd.Flags().StringSliceVar(&jobFlags.filterJobId, "filter-job-id", nil, "Filter by job ID.")
	jobListCmd.Flags().StringSliceVar(&jobFlags.filterStatus, "filter-status", nil, "Filter by job status.")
	jobListCmd.Flags().StringSliceVar(&jobFlags.filterJobGroupId, "filter-job-group-id", nil, "Filter by job group ID.")
	jobListCmd.Flags().StringSliceVar(&jobFlags.sortBy, "sort-by", nil, "Sort by fields (e.g., jobId:ASC, status:DESC).")

	jobCmd.AddCommand(jobGetCmd)

	// Job group commands
	rootCmd.AddCommand(jobGroupCmd)

	jobGroupCmd.AddCommand(jobGroupListCmd)
	jobGroupListCmd.Flags().StringSliceVar(&jobGroupFlags.filterJobGroupId, "filter-job-group-id", nil, "Filter by job group ID.")
	jobGroupListCmd.Flags().StringSliceVar(&jobGroupFlags.filterType, "filter-type", nil, "Filter by job group type.")
	jobGroupListCmd.Flags().StringSliceVar(&jobGroupFlags.sortBy, "sort-by", nil, "Sort by fields (e.g., jobGroupId:ASC).")

	jobGroupCmd.AddCommand(jobGroupGetCmd)
}
