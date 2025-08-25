## metalcloud-cli os-template example-create

Show example JSON for creating OS templates

### Synopsis

Display example JSON configuration for creating OS templates.

This command outputs a complete example JSON structure showing all available
fields and their expected values for creating OS templates. The example includes
both template configuration and associated assets.

The output can be used as a starting point for creating custom templates by
modifying the values to match your requirements.

Examples:
  # Show example JSON
  metalcloud-cli os-template example-create
  
  # Save example to file for editing
  metalcloud-cli os-template example-create > template.json

```
metalcloud-cli os-template example-create [flags]
```

### Options

```
  -h, --help   help for example-create
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

* [metalcloud-cli os-template](metalcloud-cli_os-template.md)	 - Manage OS templates for server deployments

