## metalcloud-cli network-fabric-interconnect update

Update a network fabric interconnect

### Synopsis

Update the label, name, description or BGP template of a network fabric interconnect.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect

Required Flags:
  --config-source string     Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud network-fabric-interconnect update 12 --config-source update.json
  echo '{"description":"new description"}' | metalcloud nfi update dc1-dc2 --config-source pipe

```
metalcloud-cli network-fabric-interconnect update interconnect_id_or_label [flags]
```

### Options

```
      --config-source string   Source of the updated interconnect configuration. Can be 'pipe' or path to a JSON/YAML file.
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

* [metalcloud-cli network-fabric-interconnect](metalcloud-cli_network-fabric-interconnect.md)	 - Network fabric interconnect management

