## metalcloud-cli job kill

Kill a running job

### Synopsis

Kill a running job by its ID.

This command issues a kill command for the specified job, terminating its
execution immediately.

Arguments:
  job_id (required)    The numeric ID of the job to kill.

Examples:
  metalcloud-cli job kill 12345

Permissions:
  Requires job queue write permissions to execute this command.

```
metalcloud-cli job kill job_id [flags]
```

### Options

```
  -h, --help   help for kill
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

