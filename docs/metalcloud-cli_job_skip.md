## metalcloud-cli job skip

Skip a specific job

### Synopsis

Skip a pending or running job by its ID.

This command signals the system to skip the specified job, preventing it from
executing or stopping its current execution.

Arguments:
  job_id (required)    The numeric ID of the job to skip.

Examples:
  metalcloud-cli job skip 12345

Permissions:
  Requires job queue write permissions to execute this command.

```
metalcloud-cli job skip job_id [flags]
```

### Options

```
  -h, --help   help for skip
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

