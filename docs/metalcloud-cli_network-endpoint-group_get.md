## metalcloud-cli network-endpoint-group get

Get network endpoint group details

### Synopsis

Get the details of a network endpoint group.

Required Arguments:
  network_endpoint_group_id_or_name   The ID or name of the network endpoint group

Examples:
  metalcloud network-endpoint-group get 3
  metalcloud neg get dc1-endpoint-group

```
metalcloud-cli network-endpoint-group get network_endpoint_group_id_or_name [flags]
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

* [metalcloud-cli network-endpoint-group](metalcloud-cli_network-endpoint-group.md)	 - Network endpoint group management

