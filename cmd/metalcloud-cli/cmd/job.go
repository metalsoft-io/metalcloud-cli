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
		Short: "Manage MetalCloud jobs and job execution",
		Long: `Manage MetalCloud jobs and job execution.

Jobs in MetalCloud represent asynchronous operations that are executed by the system.
These commands allow you to list, view, and monitor job execution status and details.

Available Commands:
  list    List jobs with optional filtering and sorting
  get     Get detailed information about a specific job

Use "metalcloud-cli job [command] --help" for more information about a command.`,
	}

	jobListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List jobs with optional filtering and sorting",
		Long: `List MetalCloud jobs with optional filtering and sorting capabilities.

This command displays all jobs in the system with their current status, function names,
creation timestamps, and group associations. You can filter results by job ID, status,
or job group ID, and sort the output by various fields.

Flags:
  --filter-job-id strings        Filter jobs by specific job IDs. Accepts multiple values.
                                 Example: --filter-job-id 123,456,789

  --filter-status strings        Filter jobs by status. Common statuses include:
                                 - pending: Job is waiting to be executed
                                 - running: Job is currently executing
                                 - completed: Job has finished successfully
                                 - failed: Job has failed
                                 - cancelled: Job was cancelled
                                 Example: --filter-status pending,running

  --filter-job-group-id strings  Filter jobs by job group ID. Useful for viewing
                                 jobs that belong to specific operation groups.
                                 Example: --filter-job-group-id 10,20

  --sort-by strings              Sort results by specified fields. Format: field:direction
                                 Available fields: jobId, status, functionName, createdTimestamp, jobGroupId
                                 Directions: ASC (ascending), DESC (descending)
                                 Example: --sort-by jobId:DESC,status:ASC

Examples:
  # List all jobs
  metalcloud-cli job list

  # List only pending and running jobs
  metalcloud-cli job list --filter-status pending,running

  # List jobs sorted by creation time (newest first)
  metalcloud-cli job list --sort-by createdTimestamp:DESC

  # List jobs for specific job group, sorted by job ID
  metalcloud-cli job list --filter-job-group-id 15 --sort-by jobId:ASC

  # List specific jobs by ID
  metalcloud-cli job list --filter-job-id 123,456`,
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
		Use:     "get job_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific job",
		Long: `Get detailed information about a specific job by its ID.

This command retrieves and displays comprehensive details about a single job,
including its current status, execution parameters, response data, retry information,
timestamps, and any associated metadata.

Arguments:
  job_id (required)    The numeric ID of the job to retrieve. Must be a valid
                       job identifier that exists in the system.

The output includes:
  - Job ID and basic information
  - Current status and function name
  - Creation, update, and start timestamps
  - Execution parameters and response data
  - Retry count and configuration
  - Duration and performance metrics
  - Associated infrastructure components
  - Job group membership
  - Error information (if applicable)

Examples:
  # Get details for job with ID 12345
  metalcloud-cli job get 12345

  # Get job details (using alias)
  metalcloud-cli job show 12345

Permissions:
  Requires job queue read permissions to execute this command.`,
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
		Short: "Manage MetalCloud job groups and group operations",
		Long: `Manage MetalCloud job groups and group operations.

Job groups in MetalCloud represent collections of related jobs that are executed together
as part of a larger operation. These commands allow you to list, view, and monitor
job group execution status and their constituent jobs.

Available Commands:
  list    List job groups with optional filtering and sorting
  get     Get detailed information about a specific job group

Use "metalcloud-cli job-group [command] --help" for more information about a command.`,
	}

	jobGroupListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List job groups with optional filtering and sorting",
		Long: `List MetalCloud job groups with optional filtering and sorting capabilities.

This command displays all job groups in the system with their current status, types,
creation timestamps, and execution progress. You can filter results by job group ID,
type, and sort the output by various fields.

Flags:
  --filter-job-group-id strings  Filter job groups by specific job group IDs. 
                                 Accepts multiple values for listing specific groups.
                                 Example: --filter-job-group-id 10,15,20

  --filter-type strings          Filter job groups by their type. Common types include:
                                 - infrastructure_deploy: Infrastructure deployment operations
                                 - infrastructure_edit: Infrastructure modification operations
                                 - server_provision: Server provisioning operations
                                 - server_decomission: Server decommissioning operations
                                 - storage_operation: Storage-related operations
                                 - network_operation: Network configuration operations
                                 Example: --filter-type infrastructure_deploy,server_provision

  --sort-by strings              Sort results by specified fields. Format: field:direction
                                 Available fields: jobGroupId, type, createdTimestamp, status
                                 Directions: ASC (ascending), DESC (descending)
                                 Example: --sort-by createdTimestamp:DESC,jobGroupId:ASC

Examples:
  # List all job groups
  metalcloud-cli job-group list

  # List job groups for infrastructure operations only
  metalcloud-cli job-group list --filter-type infrastructure_deploy,infrastructure_edit

  # List specific job groups by ID
  metalcloud-cli job-group list --filter-job-group-id 10,15

  # List job groups sorted by creation time (newest first)
  metalcloud-cli job-group list --sort-by createdTimestamp:DESC

  # List infrastructure deployment groups, sorted by ID
  metalcloud-cli job-group list --filter-type infrastructure_deploy --sort-by jobGroupId:ASC

Permissions:
  Requires job queue read permissions to execute this command.`,
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
		Use:     "get job_group_id",
		Aliases: []string{"show"},
		Short:   "Get detailed information about a specific job group",
		Long: `Get detailed information about a specific job group by its ID.

This command retrieves and displays comprehensive details about a single job group,
including its status, type, execution progress, and all jobs that belong to the group.

Arguments:
  job_group_id (required)    The numeric ID of the job group to retrieve. Must be a valid
                            job group identifier that exists in the system.

The output includes:
  - Job group ID and basic information
  - Current status and type
  - Creation and update timestamps
  - Total number of jobs in the group
  - Individual job details within the group
  - Overall execution progress
  - Associated infrastructure components
  - Error information (if applicable)

Examples:
  # Get details for job group with ID 15
  metalcloud-cli job-group get 15

  # Get job group details (using alias)
  metalcloud-cli job-group show 15

Permissions:
  Requires job queue read permissions to execute this command.`,
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
