## metalcloud-cli network-endpoint-group logical-network list

List the logical networks of a network endpoint group

### Synopsis

List the logical networks attached to a network endpoint group.

Required Arguments:
  network_endpoint_group_id_or_name   The ID or name of the network endpoint group

Examples:
  metalcloud network-endpoint-group logical-network list 3

```
metalcloud-cli network-endpoint-group logical-network list network_endpoint_group_id_or_name [flags]
```

### Options

```
  -h, --help   help for list
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

