## metalcloud-cli server capabilities

Get server capabilities

### Synopsis

Get server capabilities.

This command retrieves information about the capabilities supported by the
specified server, including firmware upgrade support, VNC capabilities,
and other available features.

Required Arguments:
  server_id              The ID of the server to query

Examples:
  # Get capabilities for server with ID 123
  metalcloud-cli server capabilities 123


```
metalcloud-cli server capabilities server_id [flags]
```

### Options

```
  -h, --help   help for capabilities
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

