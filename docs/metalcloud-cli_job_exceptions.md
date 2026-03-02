## metalcloud-cli job exceptions

Get exceptions for a specific job

### Synopsis

Get the list of exceptions thrown during execution of a specific job.

This command shows all errors and exceptions that occurred while the job was running,
which is useful for debugging failed or errored jobs.

Arguments:
  job_id (required)    The numeric ID of the job.

Examples:
  metalcloud-cli job exceptions 12345

```
metalcloud-cli job exceptions job_id [flags]
```

### Options

```
  -h, --help   help for exceptions
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

