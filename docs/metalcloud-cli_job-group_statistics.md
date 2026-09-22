## metalcloud-cli job-group statistics

Get the job counters of a specific job group

### Synopsis

Get the execution counters of a job group: total number of jobs, how many
completed and how many threw an error, together with the group type and its
creation and completion timestamps.

Required Arguments:
  job_group_id    The numeric ID of the job group

Examples:
  # Get the statistics of job group 15
  metalcloud-cli job-group statistics 15

  # Using the alias
  metalcloud-cli job-group stats 15

Permissions:
  Requires job queue read permissions to execute this command.

```
metalcloud-cli job-group statistics job_group_id [flags]
```

### Options

```
  -h, --help   help for statistics
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

* [metalcloud-cli job-group](metalcloud-cli_job-group.md)	 - Manage MetalCloud job groups and group operations

