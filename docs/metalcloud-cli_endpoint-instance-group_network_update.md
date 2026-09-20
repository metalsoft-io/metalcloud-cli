## metalcloud-cli endpoint-instance-group network update

Update a network connection of an endpoint instance group

### Synopsis

Update an existing network connection of an endpoint instance group.

The change can be described either with a configuration file (--config-source) or with
individual flags (--access-mode, --tagged, ...).

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group
  connection_id                The numeric ID of the network connection

Required Flags (one of):
  --config-source string          Source of the updated connection. Can be 'pipe' or path to a JSON/YAML file.
  --access-mode string            New access mode of the connection
  --tagged string                 Whether the logical network is tagged (true/false)
  --mtu int                       New MTU of the connection
  --redundancy string             New redundancy mode of the connection
  --provides-default-route string Whether the connection provides the default route (true/false)

Examples:
  metalcloud-cli endpoint-instance-group network update 12 5 --access-mode l2
  metalcloud-cli eig net edit 12 5 --config-source connection.json

```
metalcloud-cli endpoint-instance-group network update endpoint_instance_group_id connection_id [flags]
```

### Options

```
      --access-mode string              New access mode of the connection.
      --config-source string            Source of the updated connection. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                            help for update
      --mtu int32                       New MTU of the connection.
      --provides-default-route string   Whether the connection provides the default route (true/false).
      --redundancy string               New redundancy mode of the connection.
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

