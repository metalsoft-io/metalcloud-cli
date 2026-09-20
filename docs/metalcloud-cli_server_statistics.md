## metalcloud-cli server statistics

Get aggregated server statistics

### Synopsis

Get the aggregated server counts, grouped by server status and by site.

Examples:
  # Show the server statistics
  metalcloud-cli server statistics

  # Show the server statistics as JSON
  metalcloud-cli server statistics -f json


```
metalcloud-cli server statistics [flags]
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

* [metalcloud-cli server](metalcloud-cli_server.md)	 - Server management

