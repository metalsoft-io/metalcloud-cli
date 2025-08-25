## metalcloud-cli infrastructure user-limits

Display resource limits for an infrastructure

### Synopsis

Show the resource limits configured for a specific infrastructure including compute nodes,
drives, and infrastructure count limits.

This information helps understand the maximum resources that can be provisioned within
the infrastructure and plan capacity accordingly.

Arguments:
  infrastructure_id_or_label  The ID (numeric) or label (string) of the infrastructure

Examples:
  # Get user limits for infrastructure by ID
  metalcloud-cli infrastructure user-limits 123

  # Get user limits for infrastructure by label
  metalcloud-cli infrastructure user-limits "web-cluster"

  # Using the alias
  metalcloud-cli infrastructure get-user-limits production-env

```
metalcloud-cli infrastructure user-limits infrastructure_id_or_label [flags]
```

### Options

```
  -h, --help   help for user-limits
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

