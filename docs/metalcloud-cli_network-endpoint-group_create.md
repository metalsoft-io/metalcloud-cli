## metalcloud-cli network-endpoint-group create

Create a network endpoint group

### Synopsis

Create a new network endpoint group.

Required Flags (one of):
  --config-source string   Source of the configuration. Can be 'pipe' or path to a JSON/YAML file.
  --name string            Name of the new network endpoint group

Optional Flags (when not using --config-source):
  --site-id int            ID of the site the network endpoint group belongs to

Examples:
  metalcloud network-endpoint-group create --config-source neg.json
  metalcloud neg create --name dc1-endpoint-group --site-id 1

```
metalcloud-cli network-endpoint-group create [flags]
```

### Options

```
      --config-source string   Source of the new network endpoint group configuration. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for create
      --name string            Name of the new network endpoint group.
      --site-id int            ID of the site the network endpoint group belongs to.
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

