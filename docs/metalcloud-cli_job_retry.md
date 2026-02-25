## metalcloud-cli job retry

Retry a specific job

### Synopsis

Retry a failed job by its ID.

This command re-queues the specified job for execution. By default it only retries
jobs that have failed. Use --retry-even-if-successful to retry a job regardless
of its previous outcome.

Arguments:
  job_id (required)    The numeric ID of the job to retry.

Flags:
  --retry-even-if-successful    Retry the job even if it previously succeeded.

Examples:
  # Retry a failed job
  metalcloud-cli job retry 12345

  # Retry a job regardless of its previous status
  metalcloud-cli job retry 12345 --retry-even-if-successful

Permissions:
  Requires job queue write permissions to execute this command.

```
metalcloud-cli job retry job_id [flags]
```

### Options

```
  -h, --help                       help for retry
      --retry-even-if-successful   Retry even if the job was successful.
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

