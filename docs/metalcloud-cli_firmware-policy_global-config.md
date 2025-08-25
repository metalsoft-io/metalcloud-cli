## metalcloud-cli firmware-policy global-config

Manage global firmware configuration settings

### Synopsis

Manage global firmware configuration settings that control system-wide firmware upgrade behavior.

The global firmware configuration defines when firmware upgrades can be executed,
whether they are enabled system-wide, and other global constraints that affect
all firmware policies. This configuration acts as a master control for the
entire firmware upgrade system.

Available subcommands:
  get                     Get current global firmware configuration
  update                  Update global firmware configuration
  config-example          Show example global configuration

Examples:
  # Get current global configuration
  metalcloud-cli firmware-policy global-config get
  
  # Update global configuration from file
  metalcloud-cli firmware-policy global-config update --config-source config.json

### Options

```
  -h, --help   help for global-config
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
* [metalcloud-cli firmware-policy global-config config-example](metalcloud-cli_firmware-policy_global-config_config-example.md)	 - Show example global firmware configuration in JSON format
* [metalcloud-cli firmware-policy global-config get](metalcloud-cli_firmware-policy_global-config_get.md)	 - Get current global firmware configuration settings
* [metalcloud-cli firmware-policy global-config update](metalcloud-cli_firmware-policy_global-config_update.md)	 - Update global firmware configuration settings

