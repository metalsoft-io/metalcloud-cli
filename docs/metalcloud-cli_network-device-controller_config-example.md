## metalcloud-cli network-device-controller config-example

Print an example controller create configuration

### Synopsis

Print an example configuration that can be edited and passed to
'network-device-controller create --config-source'.

Examples:
  metalcloud network-device-controller config-example > controller.json
  metalcloud ndc config-example -f yaml

```
metalcloud-cli network-device-controller config-example [flags]
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

* [metalcloud-cli network-device-controller](metalcloud-cli_network-device-controller.md)	 - Network device controller management

