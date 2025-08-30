## metalcloud-cli dns-zone list

List DNS zones

### Synopsis

List all DNS zones in the MetalSoft infrastructure.

This command displays information about all DNS zones including their IDs, labels, 
zone names, zone types, status, and other configuration details.

Optional Flags:
  --filter-default    Filter zones by default status (true/false)

Examples:
  # List all DNS zones
  metalcloud-cli dns-zone list

  # Filter zones by status
  metalcloud-cli dns-zone list --filter-default true


```
metalcloud-cli dns-zone list [flags]
```

### Options

```
      --filter-default strings   Filter the result by default status.
  -h, --help                     help for list
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

* [metalcloud-cli dns-zone](metalcloud-cli_dns-zone.md)	 - DNS Zone management

