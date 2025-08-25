## metalcloud-cli firmware-policy config-example

Show example firmware policy configuration in JSON format

### Synopsis

Show an example firmware policy configuration in JSON format.

This command displays a sample JSON configuration that can be used as a template
for creating new firmware policies. The example includes all available fields
with sample values and explains the structure of policy rules and server
instance group associations.

The example output can be saved to a file and modified to create your own
firmware policy configurations.

No flags or arguments are required for this command.

Examples:
  # Show example configuration
  metalcloud-cli firmware-policy config-example
  
  # Save example to file for editing
  metalcloud-cli fw-policy config-example > my-policy.json
  
  # Use example as template for creating policy
  metalcloud-cli firmware-policy config-example > policy.json
  # Edit policy.json with your values
  metalcloud-cli firmware-policy create --config-source policy.json

```
metalcloud-cli firmware-policy config-example [flags]
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

* [metalcloud-cli firmware-policy](metalcloud-cli_firmware-policy.md)	 - Manage server firmware upgrade policies and global firmware configurations

