## metalcloud-cli configuration replace

Replace a section of the global platform configuration

### Synopsis

Replace the ENTIRE configuration of one service section of the platform (HTTP PUT).

Any setting missing from the supplied document is reset by the platform: use
'configuration update' to change individual settings. This affects every user of
the installation - take a copy with 'configuration get <service>' first.

Required Arguments:
  service   The configuration section to replace. One of: auth, platform, notification, tunnel, gateway-api, image-builder

Required Flags:
  --config-source string   Source of the new section configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud configuration get platform -f json > platform.json
  metalcloud configuration replace platform --config-source platform.json
  cat platform.yaml | metalcloud configuration replace platform --config-source pipe

```
metalcloud-cli configuration replace service [flags]
```

### Options

```
      --config-source string   Source of the new section configuration. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for replace
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

