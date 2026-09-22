## metalcloud-cli site statistics

Show aggregated statistics for all sites

### Synopsis

Show aggregated statistics for all sites: how many sites exist, how many are
active, and how many site controllers were seen online, may be offline or are offline.

Use --format json or yaml to also get the per-site resource counts.

Required Permissions:
  sites:read - Permission to view site information

Examples:
  # Show the site statistics
  metalcloud-cli site statistics

  # Show the full statistics as JSON
  metalcloud-cli site stats -f json

```
metalcloud-cli site statistics [flags]
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

* [metalcloud-cli site](metalcloud-cli_site.md)	 - Manage sites (datacenters) and their configurations

