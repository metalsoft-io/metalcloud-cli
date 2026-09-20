## metalcloud-cli fabric bgp-session config-example

Show a BGP session configuration example

### Synopsis

Print an example BGP session create body. Edit it and pass it to
'fabric bgp-session create' via --config-source.

Examples:
  metalcloud-cli fabric bgp-session config-example
  metalcloud-cli fabric bgp-session config-example -f yaml > session.yaml

```
metalcloud-cli fabric bgp-session config-example [flags]
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

* [metalcloud-cli fabric bgp-session](metalcloud-cli_fabric_bgp-session.md)	 - Manage the BGP sessions of a fabric

