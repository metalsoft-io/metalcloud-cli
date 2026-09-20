## metalcloud-cli configuration get

Get the global platform configuration

### Synopsis

Show the global platform configuration, optionally restricted to one service section.

Optional Arguments:
  service   One of: auth, platform, notification, tunnel, gateway-api, image-builder. When omitted, every section is returned.

The document is deeply nested, so text output is rendered as YAML. Use -f json or
-f yaml for machine-readable output.

Examples:
  metalcloud configuration get
  metalcloud configuration get platform
  metalcloud configuration get auth -f json

```
metalcloud-cli configuration get [service] [flags]
```

### Options

```
  -h, --help   help for get
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

* [metalcloud-cli configuration](metalcloud-cli_configuration.md)	 - Global platform configuration management

