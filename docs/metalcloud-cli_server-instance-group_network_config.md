## metalcloud-cli server-instance-group network config

Get the network configuration of a server instance group

### Synopsis

Get the network endpoint group that holds the network configuration of a server
instance group. Use 'network list' for the individual network connections.

Required Arguments:
  server_instance_group_id  The numeric ID of the server instance group

Examples:
  metalcloud-cli server-instance-group network config 1234
  metalcloud-cli ig net get-config 1234

```
metalcloud-cli server-instance-group network config server_instance_group_id [flags]
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

* [metalcloud-cli server-instance-group network](metalcloud-cli_server-instance-group_network.md)	 - Manage network connections for server instance groups

