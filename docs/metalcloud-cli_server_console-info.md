## metalcloud-cli server console-info

Get server remote console information

### Synopsis

Get server remote console information.

This command retrieves remote console connection details for the specified server,
including active connections and console access information.

Required Arguments:
  server_id              The ID of the server to query

Examples:
  # Get remote console information for server with ID 123
  metalcloud-cli server console-info 123


```
metalcloud-cli server console-info server_id [flags]
```

### Options

```
  -h, --help   help for console-info
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

