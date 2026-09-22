## metalcloud-cli network-endpoint-group logical-network get

Get one logical network of a network endpoint group

### Synopsis

Get the settings of one logical network attached to a network endpoint group.

Required Arguments:
  network_endpoint_group_id_or_name   The ID or name of the network endpoint group
  logical_network_id                  The numeric ID of the logical network

Examples:
  metalcloud network-endpoint-group logical-network get 3 44

```
metalcloud-cli network-endpoint-group logical-network get network_endpoint_group_id_or_name logical_network_id [flags]
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

* [metalcloud-cli network-endpoint-group logical-network](metalcloud-cli_network-endpoint-group_logical-network.md)	 - Network endpoint group logical network management

