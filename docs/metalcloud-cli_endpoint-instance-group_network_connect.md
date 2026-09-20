## metalcloud-cli endpoint-instance-group network connect

Connect an endpoint instance group to a logical network

### Synopsis

Create a network connection between an endpoint instance group and a logical network.

The connection can be described either with a configuration file (--config-source) or
with individual flags (--logical-network-id, --access-mode, ...).

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group

Required Flags (one of):
  --config-source string       Source of the new connection configuration. Can be 'pipe' or path to a JSON/YAML file.
  --logical-network-id string  ID of the logical network to connect to

Optional Flags (when not using --config-source):
  --access-mode string            Access mode of the connection (default "l2")
  --tagged string                 Whether the logical network is tagged (true/false)
  --mtu int                       MTU of the connection
  --redundancy string             Redundancy mode of the connection
  --provides-default-route string Whether the connection provides the default route (true/false)

Examples:
  metalcloud-cli endpoint-instance-group network connect 12 --logical-network-id 7 --tagged true
  metalcloud-cli endpoint-instance-group network connect 12 --config-source connection.json

```
metalcloud-cli endpoint-instance-group network connect endpoint_instance_group_id [flags]
```

### Options

```
      --access-mode string              Access mode of the connection. (default "l2")
      --config-source string            Source of the new connection configuration. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                            help for connect
      --logical-network-id string       ID of the logical network to connect to.
      --mtu int32                       MTU of the connection.
      --provides-default-route string   Whether the connection provides the default route (true/false).
      --redundancy string               Redundancy mode of the connection.
      --tagged string                   Whether the logical network is tagged (true/false).
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

* [metalcloud-cli endpoint-instance-group network](metalcloud-cli_endpoint-instance-group_network.md)	 - Manage the network configuration of an endpoint instance group

