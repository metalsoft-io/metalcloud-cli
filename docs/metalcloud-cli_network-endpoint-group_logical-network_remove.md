## metalcloud-cli network-endpoint-group logical-network remove

Detach a logical network from a network endpoint group

### Synopsis

Detach a logical network from a network endpoint group.

Required Arguments:
  network_endpoint_group_id_or_name   The ID or name of the network endpoint group
  logical_network_id                  The numeric ID of the logical network

Examples:
  metalcloud network-endpoint-group logical-network remove 3 44

```
metalcloud-cli network-endpoint-group logical-network remove network_endpoint_group_id_or_name logical_network_id [flags]
```

### Options

```
  -h, --help   help for remove
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

