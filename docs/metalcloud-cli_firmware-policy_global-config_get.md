## metalcloud-cli firmware-policy global-config get

Get current global firmware configuration settings

### Synopsis

Get current global firmware configuration settings that control system-wide firmware upgrade behavior.

This command retrieves and displays the global firmware configuration which includes
settings such as whether firmware upgrades are enabled globally, upgrade time windows,
scheduling constraints, and other system-wide policies that affect all firmware
upgrade operations.

The global configuration acts as a master switch and constraint system for all
firmware policies. Even if individual policies are active, they must comply with
the global configuration settings.

No flags or arguments are required for this command.

Examples:
  # Get current global firmware configuration
  metalcloud-cli firmware-policy global-config get
  
  # Get global config using alias
  metalcloud-cli fw-policy global get

```
metalcloud-cli firmware-policy global-config get [flags]
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

* [metalcloud-cli firmware-policy global-config](metalcloud-cli_firmware-policy_global-config.md)	 - Manage global firmware configuration settings

