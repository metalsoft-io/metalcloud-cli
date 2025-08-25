## metalcloud-cli vm-type config-example

Show an example configuration for creating VM types

### Synopsis

Display an example JSON configuration that can be used to create a new VM type.

This command outputs a sample configuration showing all available fields and their expected values.
The generated configuration can be saved to a file and modified as needed for creating or updating VM types.

EXAMPLES:
  # Display configuration example
  metalcloud vm-type config-example
  
  # Save example to file for editing
  metalcloud vm-type config-example > vm-type-config.json
  
  # Use saved configuration to create a VM type
  metalcloud vm-type create --config-source vm-type-config.json

```
metalcloud-cli vm-type config-example [flags]
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

* [metalcloud-cli vm-type](metalcloud-cli_vm-type.md)	 - Manage VM types and configurations

