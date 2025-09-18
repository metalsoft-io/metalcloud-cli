## metalcloud-cli network-device add-defaults

Add network device default configuration

### Synopsis

Add network device default configuration that will be applied to new
devices when they are added to sites. These defaults provide consistent
baseline configurations across your infrastructure.

Default configurations can include:
- Management network settings and credentials
- Standard VLAN configurations
- Security policies and access controls
- Monitoring and logging settings
- Device-specific operational parameters
- Network topology preferences

The configuration is provided via JSON file or pipe input and will be merged
with existing defaults, allowing for incremental updates.

Required Flags:
  --config-source   Source of default configuration data (required)
                   Values: 'pipe' for stdin input, or path to JSON file

Use the 'example-defaults' command to see the configuration format:

Examples:
  # Add defaults from JSON file
  metalcloud-cli network-device add-defaults --config-source defaults.json

  # Add defaults from pipe input
  cat site-defaults.json | metalcloud-cli network-device add-defaults --config-source pipe

  # Update specific default settings
  echo '{"syslogEnabled": true, "managementPort": 22}' | metalcloud-cli nd add-defaults --config-source pipe

```
metalcloud-cli network-device add-defaults [flags]
```

### Options

```
      --config-source string   Source of the network device default configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for add-defaults
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

* [metalcloud-cli network-device](metalcloud-cli_network-device.md)	 - Manage network devices (switches) in the infrastructure

