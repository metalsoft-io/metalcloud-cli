## metalcloud-cli server-instance-group interface

Get one interface of a server instance group

### Synopsis

Get the details of one network interface of a server instance group.

Required Arguments:
  server_instance_group_id  The numeric ID of the server instance group
  interface_id              The numeric ID of the interface

Examples:
  metalcloud-cli server-instance-group interface 1234 7
  metalcloud-cli ig iface 1234 7

```
metalcloud-cli server-instance-group interface server_instance_group_id interface_id [flags]
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

* [metalcloud-cli server-instance-group](metalcloud-cli_server-instance-group.md)	 - Manage server instance groups within infrastructures

