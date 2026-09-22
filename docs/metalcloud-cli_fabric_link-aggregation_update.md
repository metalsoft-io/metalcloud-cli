## metalcloud-cli fabric link-aggregation update

Update a link aggregation of a fabric

### Synopsis

Update a link aggregation of a network fabric from a JSON/YAML configuration.

The body replaces the aggregated links (linkIds), the options and the custom
variables; the aggregation type is fixed at creation time.

Required Arguments:
  fabric_id              The ID or name of the fabric
  link_aggregation_id    The ID of the link aggregation

Required Flags:
  --config-source   'pipe' to read from stdin, or a path to a JSON/YAML file.

Examples:
  metalcloud-cli fabric link-aggregation update 12345 3 --config-source lag.yaml
  echo '{"linkIds":[1,2,3]}' | metalcloud-cli fabric lag update 12345 3 --config-source pipe

```
metalcloud-cli fabric link-aggregation update fabric_id link_aggregation_id [flags]
```

### Options

```
      --config-source string   Source of the link aggregation configuration. Can be 'pipe' or path to a JSON/YAML file.
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

* [metalcloud-cli fabric link-aggregation](metalcloud-cli_fabric_link-aggregation.md)	 - Manage the link aggregations of a fabric

