## metalcloud-cli job statistics

Get job queue statistics

### Synopsis

Get statistics about the job queue including counts by status and total archived jobs.

Examples:
  metalcloud-cli job statistics
  metalcloud-cli job stats

```
metalcloud-cli job statistics [flags]
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

* [metalcloud-cli job](metalcloud-cli_job.md)	 - Manage MetalCloud jobs and job execution

