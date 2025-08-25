## metalcloud-cli job list

List jobs with optional filtering and sorting

### Synopsis

List MetalCloud jobs with optional filtering and sorting capabilities.

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
  metalcloud-cli job list --filter-job-id 123,456

```
metalcloud-cli job list [flags]
```

### Options

```
      --filter-job-group-id strings   Filter by job group ID.
      --filter-job-id strings         Filter by job ID.
      --filter-status strings         Filter by job status.
  -h, --help                          help for list
      --sort-by strings               Sort by fields (e.g., jobId:ASC, status:DESC).
```

### Options inherited from parent commands

```
  -k, --api_key string         MetalCloud API key
  -c, --config string          Config file path
  -d, --debug                  Set to enable debug logging
  -e, --endpoint string        MetalCloud API endpoint
  -f, --format string          Output format. Supported values are 'text','csv','md','json','yaml'. (default "text")
  -i, --insecure_skip_verify   Set to allow insecure transport
  -l, --log_file string        Log file path
  -v, --verbosity string       Log level verbosity (default "INFO")
```

### SEE ALSO

* [metalcloud-cli job](metalcloud-cli_job.md)	 - Manage MetalCloud jobs and job execution

