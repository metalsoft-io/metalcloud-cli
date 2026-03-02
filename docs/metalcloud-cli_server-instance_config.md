## metalcloud-cli server-instance config

Get configuration of a server instance

### Synopsis

Get the current configuration of a server instance.

This command retrieves the configuration details of a server instance including
server type, OS template, hostname, deploy status, and other settings.

Arguments:
  server_instance_id  The numeric ID of the server instance

Examples:
  # Get configuration of server instance 5678
  metalcloud-cli server-instance config 5678

  # Using alias
  metalcloud-cli inst get-config 5678

```
metalcloud-cli server-instance config <server_instance_id> [flags]
```

### Options

```
  -h, --help   help for config
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

