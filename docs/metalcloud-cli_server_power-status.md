## metalcloud-cli server power-status

Get server power status

### Synopsis

Get the current power status of a server.

This command retrieves the current power state of the specified server
from its BMC/IPMI interface.

Required Arguments:
  server_id              The ID of the server to query

Examples:
  # Get power status for server with ID 123
  metalcloud-cli server power-status 123


```
metalcloud-cli server power-status server_id [flags]
```

### Options

```
  -h, --help   help for power-status
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

