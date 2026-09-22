## metalcloud-cli fabric link-aggregation create

Create a link aggregation on a fabric

### Synopsis

Create a link aggregation on a network fabric from a JSON/YAML configuration.

The configuration contains:
- type: the link aggregation type (e.g. lag, mlag, mlag-peer-link)
- linkIds: the IDs of the fabric links to aggregate (see 'fabric get-links')
- mlagDomainIdentifier: only for the mlag-peer-link type
- customVariables / options: optional

Required Arguments:
  fabric_id    The ID or name of the fabric

Required Flags:
  --config-source   'pipe' to read from stdin, or a path to a JSON/YAML file.

Examples:
  metalcloud-cli fabric link-aggregation config-example > lag.yaml
  metalcloud-cli fabric link-aggregation create 12345 --config-source lag.yaml
  cat lag.yaml | metalcloud-cli fabric lag create 12345 --config-source pipe

```
metalcloud-cli fabric link-aggregation create fabric_id [flags]
```

### Options

```
      --config-source string   Source of the link aggregation configuration. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for create
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

* [metalcloud-cli fabric link-aggregation](metalcloud-cli_fabric_link-aggregation.md)	 - Manage the link aggregations of a fabric

