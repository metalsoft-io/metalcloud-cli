## metalcloud-cli network-endpoint-group update

Update a network endpoint group

### Synopsis

Update the name or site of a network endpoint group.

Required Arguments:
  network_endpoint_group_id_or_name   The ID or name of the network endpoint group

Required Flags:
  --config-source string              Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud network-endpoint-group update 3 --config-source update.json
  echo '{"name":"new name"}' | metalcloud neg update dc1-endpoint-group --config-source pipe

```
metalcloud-cli network-endpoint-group update network_endpoint_group_id_or_name [flags]
```

### Options

```
      --config-source string   Source of the updated network endpoint group configuration. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for update
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

