## metalcloud-cli fabric link-aggregation config-example

Show a link aggregation configuration example

### Synopsis

Print an example link aggregation create body. Edit it and pass it to
'fabric link-aggregation create' via --config-source.

Examples:
  metalcloud-cli fabric link-aggregation config-example
  metalcloud-cli fabric link-aggregation config-example -f yaml > lag.yaml

```
metalcloud-cli fabric link-aggregation config-example [flags]
```

### Options

```
  -h, --help   help for config-example
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

