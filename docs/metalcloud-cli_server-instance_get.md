## metalcloud-cli server-instance get

Get detailed information about a server instance

### Synopsis

Get detailed information about a server instance.

This command retrieves and displays comprehensive information about a specific server
instance, including its configuration, status, hardware specifications, network
connections, and metadata. The instance may be part of a server instance group
or standalone.

Arguments:
  server_instance_id  The numeric ID of the server instance to retrieve

Examples:
  # Get details of server instance with ID 5678
  metalcloud-cli server-instance get 5678

  # Using alias
  metalcloud-cli inst show 5678

```
metalcloud-cli server-instance get server_instance_id [flags]
```

### Options

```
  -h, --help   help for get
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

