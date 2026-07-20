## metalcloud-cli job-group wait

Wait for a job group to finish executing

### Synopsis

Wait for a job group to finish executing.

This command first displays the current status of the specified job group. If the
job group has not finished yet, it polls the job group status every second until
it finishes, then prints the final status.

Arguments:
  job_group_id (required)    The numeric ID of the job group to wait for. Must be
                             a valid job group identifier that exists in the system.

The command exits when:
  - The job group has finished (final status is printed)
  - An API error occurs while polling
  - The command is interrupted (e.g., Ctrl+C)

Examples:
  # Wait for job group with ID 15 to finish
  metalcloud-cli job-group wait 15

Permissions:
  Requires job queue read permissions to execute this command.

```
metalcloud-cli job-group wait job_group_id [flags]
```

### Options

```
  -h, --help   help for wait
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

