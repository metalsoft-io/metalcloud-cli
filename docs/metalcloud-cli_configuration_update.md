## metalcloud-cli configuration update

Update part of a section of the global platform configuration

### Synopsis

Merge a partial document into one service section of the global platform
configuration (HTTP PATCH). Settings absent from the document are left untouched.

This changes platform-wide settings that affect every user of the installation.

Required Arguments:
  service   The configuration section to update. One of: auth, platform, notification, tunnel, gateway-api, image-builder

Required Flags:
  --config-source string   Source of the partial configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud configuration update notification --config-source notification-patch.json
  echo '{"smtp":{"port":587}}' | metalcloud configuration update notification --config-source pipe

```
metalcloud-cli configuration update service [flags]
```

### Options

```
      --config-source string   Source of the partial configuration. Can be 'pipe' or path to a JSON/YAML file.
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

* [metalcloud-cli configuration](metalcloud-cli_configuration.md)	 - Global platform configuration management

