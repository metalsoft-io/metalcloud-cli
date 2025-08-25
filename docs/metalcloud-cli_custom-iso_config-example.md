## metalcloud-cli custom-iso config-example

Show configuration example for creating custom ISOs

### Synopsis

Display a JSON configuration example that can be used as a template for creating
or updating custom ISO images.

The configuration example shows all available fields, their data types, and
expected values. You can save this output to a file and modify it according
to your requirements.

Required permissions:
  - custom_iso:write

Examples:
  # Show configuration example
  metalcloud-cli custom-iso config-example
  
  # Save example to file for editing
  metalcloud-cli custom-iso config-example > custom-iso-config.json

```
metalcloud-cli custom-iso config-example [flags]
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

* [metalcloud-cli custom-iso](metalcloud-cli_custom-iso.md)	 - Manage custom ISO images for server provisioning

