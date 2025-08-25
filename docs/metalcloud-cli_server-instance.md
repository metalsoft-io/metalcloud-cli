## metalcloud-cli server-instance

Manage individual server instances

### Synopsis

Server Instance management commands.

Server Instances are individual compute resources within server instance groups.
They represent physical or virtual servers with specific hardware configurations
and network connections. Each instance inherits properties from its parent
instance group but can have individual characteristics and status.

Available commands include:
- get: View detailed information about a specific server instance

Use "metalcloud-cli server-instance [command] --help" for detailed information about each command.

### Options

```
  -h, --help   help for server-instance
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

* [metalcloud-cli](metalcloud-cli.md)	 - MetalCloud CLI
* [metalcloud-cli server-instance get](metalcloud-cli_server-instance_get.md)	 - Get detailed information about a server instance

