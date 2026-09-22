## metalcloud-cli server-instance statistics

Get server instance statistics

### Synopsis

Get global server instance statistics: the number of instances per server status
and per site.

Examples:
  metalcloud-cli server-instance statistics
  metalcloud-cli inst stats -f json

```
metalcloud-cli server-instance statistics [flags]
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

* [metalcloud-cli server-instance](metalcloud-cli_server-instance.md)	 - Manage individual server instances

