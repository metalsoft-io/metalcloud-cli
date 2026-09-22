## metalcloud-cli logical-network update-config

Update the global settings of a logical network config

### Synopsis

Update the global settings of a logical network's config.

Only the config's global settings are updated here (for example the MTU and the
VXLAN properties); the allocation strategies are managed with the
'logical-network allocation-strategy' commands. The config's own revision is
sent as the If-Match entity tag, so a concurrent change is rejected instead of
being overwritten.

Required Arguments:
  logical_network_id  The ID of the logical network

Required Flags:
  --config-source  'pipe' to read from stdin, or a path to a JSON/YAML file

Examples:
  metalcloud-cli logical-network update-config 12 --config-source settings.json
  echo '{"mtu":9000}' | metalcloud-cli logical-network update-config 12 --config-source pipe

```
metalcloud-cli logical-network update-config logical_network_id [flags]
```

### Options

```
      --config-source string   Source of the logical network config updates. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for update-config
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

* [metalcloud-cli logical-network](metalcloud-cli_logical-network.md)	 - Manage logical networks within fabrics

