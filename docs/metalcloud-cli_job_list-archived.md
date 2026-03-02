## metalcloud-cli job list-archived

List archived jobs

### Synopsis

List jobs from the archive with optional filtering and sorting.

Archived jobs are completed jobs that have been moved to long-term storage.
This command supports the same filtering options as the regular job list.

Flags:
  --filter-job-id strings        Filter by job ID
  --filter-status strings        Filter by status
  --filter-job-group-id strings  Filter by job group ID
  --sort-by strings              Sort by fields (e.g., jobId:DESC)

Examples:
  # List all archived jobs
  metalcloud-cli job list-archived

  # List archived jobs filtered by status
  metalcloud-cli job list-archived --filter-status completed

  # List archived jobs sorted by job ID
  metalcloud-cli job list-archived --sort-by jobId:DESC

```
metalcloud-cli job list-archived [flags]
```

### Options

```
      --filter-job-group-id strings   Filter by job group ID.
      --filter-job-id strings         Filter by job ID.
      --filter-status strings         Filter by job status.
  -h, --help                          help for list-archived
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

