## metalcloud-cli fabric bgp-session create

Create a BGP session on a fabric

### Synopsis

Create a BGP session on a network fabric from a JSON/YAML configuration.

The configuration contains:
- bgpNumbering: inherited, numbered or unnumbered
- bgpLinkConfiguration: disabled, active or passive
- linkId or linkAggregationId: the link (or link aggregation) the session runs on
- customVariables: optional map exposed to the configuration templates

Required Arguments:
  fabric_id    The ID or name of the fabric

Required Flags:
  --config-source   'pipe' to read from stdin, or a path to a JSON/YAML file.

Examples:
  metalcloud-cli fabric bgp-session config-example > session.yaml
  metalcloud-cli fabric bgp-session create 12345 --config-source session.yaml
  cat session.yaml | metalcloud-cli fabric bgp create 12345 --config-source pipe

```
metalcloud-cli fabric bgp-session create fabric_id [flags]
```

### Options

```
      --config-source string   Source of the BGP session configuration. Can be 'pipe' or path to a JSON/YAML file.
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

* [metalcloud-cli fabric bgp-session](metalcloud-cli_fabric_bgp-session.md)	 - Manage the BGP sessions of a fabric

