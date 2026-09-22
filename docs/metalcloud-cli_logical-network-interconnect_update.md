## metalcloud-cli logical-network-interconnect update

Update a logical network interconnect

### Synopsis

Update the label, name, kind or annotations of a logical network interconnect.

Required Arguments:
  interconnect_id_or_label   The ID or label of the logical network interconnect

Required Flags:
  --config-source string     Source of the updated configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud logical-network-interconnect update 4 --config-source update.json
  echo '{"name":"new name"}' | metalcloud lni update dc1-dc2-ln --config-source pipe

```
metalcloud-cli logical-network-interconnect update interconnect_id_or_label [flags]
```

### Options

```
      --config-source string   Source of the updated logical network interconnect configuration. Can be 'pipe' or path to a JSON/YAML file.
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

* [metalcloud-cli logical-network-interconnect](metalcloud-cli_logical-network-interconnect.md)	 - Logical network interconnect management

