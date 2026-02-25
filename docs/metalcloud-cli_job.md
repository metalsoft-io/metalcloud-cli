## metalcloud-cli job

Manage MetalCloud jobs and job execution

### Synopsis

Manage MetalCloud jobs and job execution.

Jobs in MetalCloud represent asynchronous operations that are executed by the system.
These commands allow you to list, view, and monitor job execution status and details.

Available Commands:
  list    List jobs with optional filtering and sorting
  get     Get detailed information about a specific job
  skip    Skip a pending or running job
  retry   Retry a failed job
  kill    Kill a running job

Use "metalcloud-cli job [command] --help" for more information about a command.

### Options

```
  -h, --help   help for job
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

* [metalcloud-cli](metalcloud-cli.md)	 - MetalCloud CLI
* [metalcloud-cli job get](metalcloud-cli_job_get.md)	 - Get detailed information about a specific job
* [metalcloud-cli job kill](metalcloud-cli_job_kill.md)	 - Kill a running job
* [metalcloud-cli job list](metalcloud-cli_job_list.md)	 - List jobs with optional filtering and sorting
* [metalcloud-cli job retry](metalcloud-cli_job_retry.md)	 - Retry a specific job
* [metalcloud-cli job skip](metalcloud-cli_job_skip.md)	 - Skip a specific job

