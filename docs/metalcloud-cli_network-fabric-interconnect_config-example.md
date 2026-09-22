## metalcloud-cli network-fabric-interconnect config-example

Print an example interconnect create configuration

### Synopsis

Print an example configuration that can be edited and passed to
'network-fabric-interconnect create --config-source'.

Examples:
  metalcloud network-fabric-interconnect config-example > interconnect.json
  metalcloud nfi config-example -f yaml

```
metalcloud-cli network-fabric-interconnect config-example [flags]
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

* [metalcloud-cli network-fabric-interconnect](metalcloud-cli_network-fabric-interconnect.md)	 - Network fabric interconnect management

