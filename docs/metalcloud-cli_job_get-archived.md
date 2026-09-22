## metalcloud-cli job get-archived

Get a single archived job

### Synopsis

Get the details of a single job that has been moved to the job archive.

This is the single-record counterpart of 'job list-archived': once a job has
been archived it is no longer returned by 'job get'.

Required Arguments:
  job_id    The numeric ID of the archived job

Examples:
  # Get archived job 12345
  metalcloud-cli job get-archived 12345

  # Get the full archived record as JSON
  metalcloud-cli job get-archived 12345 -f json

Permissions:
  Requires job queue read permissions to execute this command.

```
metalcloud-cli job get-archived job_id [flags]
```

### Options

```
  -h, --help   help for get-archived
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

