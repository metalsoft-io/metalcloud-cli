## metalcloud-cli infrastructure all-statistics

Get deployment statistics for all infrastructures

### Synopsis

Display aggregated deployment and job execution statistics for all infrastructures.

This command provides a global overview of infrastructure deployment health, including 
total infrastructure counts, active deployments, error rates, and ongoing operations
across the entire system.

The statistics include:
- Total infrastructure count and service status breakdown
- Number of ongoing deployments and their status
- Error counts and retry statistics for failed deployments

Examples:
  # Get statistics for all infrastructures
  metalcloud-cli infrastructure all-statistics

  # Using aliases
  metalcloud-cli infrastructure all-stats
  metalcloud-cli infrastructure get-all-statistics

```
metalcloud-cli infrastructure all-statistics [flags]
```

### Options

```
  -h, --help   help for all-statistics
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

* [metalcloud-cli infrastructure](metalcloud-cli_infrastructure.md)	 - Manage infrastructure resources and configurations

