## metalcloud-cli firmware-policy global-config config-example

Show example global firmware configuration in JSON format

### Synopsis

Show an example global firmware configuration in JSON format.

This command displays a sample JSON configuration that can be used as a template
for updating the global firmware configuration. The example includes all available
fields with sample values that control system-wide firmware upgrade behavior.

The example output can be saved to a file and modified to update your global
firmware configuration settings.

No flags or arguments are required for this command.

Examples:
  # Show example global configuration
  metalcloud-cli firmware-policy global-config config-example
  
  # Save example to file for editing
  metalcloud-cli fw-policy global config-example > global-config.json
  
  # Use example as template for updating global config
  metalcloud-cli firmware-policy global-config config-example > config.json
  # Edit config.json with your values
  metalcloud-cli firmware-policy global-config update --config-source config.json

```
metalcloud-cli firmware-policy global-config config-example [flags]
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

* [metalcloud-cli firmware-policy global-config](metalcloud-cli_firmware-policy_global-config.md)	 - Manage global firmware configuration settings

