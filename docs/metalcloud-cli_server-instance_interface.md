## metalcloud-cli server-instance interface

Get one interface of a server instance

### Synopsis

Get the details of one network interface of a server instance.

Required Arguments:
  server_instance_id  The numeric ID of the server instance
  interface_id        The numeric ID of the interface

Examples:
  metalcloud-cli server-instance interface 5678 7
  metalcloud-cli inst iface 5678 7

```
metalcloud-cli server-instance interface server_instance_id interface_id [flags]
```

### Options

```
  -h, --help   help for interface
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

