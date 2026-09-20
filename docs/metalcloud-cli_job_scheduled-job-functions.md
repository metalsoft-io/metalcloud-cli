## metalcloud-cli job scheduled-job-functions

List the functions supported by scheduled jobs

### Synopsis

List the functions that can be referenced when creating a scheduled job,
together with their description and the names of their parameters.

This command lives under 'job' rather than under 'scheduled-job' because the
scheduled-job command group is wired in a separate file.

Examples:
  # List the supported scheduled job functions
  metalcloud-cli job scheduled-job-functions

  # Show the full parameter schemas
  metalcloud-cli job scheduled-job-functions -f json

Permissions:
  Requires job queue read permissions to execute this command.

```
metalcloud-cli job scheduled-job-functions [flags]
```

### Options

```
  -h, --help   help for scheduled-job-functions
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

