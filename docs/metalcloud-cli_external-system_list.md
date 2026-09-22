## metalcloud-cli external-system list

List external systems

### Synopsis

List all external systems registered with the platform.

Optional Flags:
  --filter-label strings   Filter by label. Repeatable or comma-separated.

Examples:
  metalcloud external-system list
  metalcloud ext-system ls --filter-label my-external-system

```
metalcloud-cli external-system list [flags]
```

### Options

```
      --filter-label strings   Filter by external system label.
  -h, --help                   help for list
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

* [metalcloud-cli external-system](metalcloud-cli_external-system.md)	 - External system management

