## metalcloud-cli network-endpoint-group logical-network update

Update a logical network of a network endpoint group

### Synopsis

Update the settings of one logical network attached to a network endpoint group.

Required Arguments:
  network_endpoint_group_id_or_name   The ID or name of the network endpoint group
  logical_network_id                  The numeric ID of the logical network

Optional Flags:
  --config-source string        Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.
                                Mutually exclusive with the flags below.
  --access-mode string          Access mode of the connection
  --tagged                      Attach the logical network tagged
  --mtu int                     MTU of the connection
  --provides-default-route      The logical network provides the default route
  --disable-auto-ip-allocation  Disable automatic IPv4 allocation on this connection

Examples:
  metalcloud network-endpoint-group logical-network update 3 44 --mtu 9000
  metalcloud neg logical-network update dc1-endpoint-group 44 --config-source connection.json

```
metalcloud-cli network-endpoint-group logical-network update network_endpoint_group_id_or_name logical_network_id [flags]
```

### Options

```
      --access-mode string           Access mode of the connection. (default "l2")
      --config-source string         Source of the connection configuration. Can be 'pipe' or path to a JSON/YAML file.
      --disable-auto-ip-allocation   Disable automatic IPv4 allocation on this connection.
  -h, --help                         help for update
      --mtu int32                    MTU of the connection.
      --provides-default-route       The logical network provides the default route.
      --tagged                       Attach the logical network tagged.
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

* [metalcloud-cli network-endpoint-group logical-network](metalcloud-cli_network-endpoint-group_logical-network.md)	 - Network endpoint group logical network management

