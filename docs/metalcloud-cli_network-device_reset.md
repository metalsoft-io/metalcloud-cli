## metalcloud-cli network-device reset

Reset network device to factory defaults (destructive operation)

### Synopsis

Reset a network device to its factory default state, destroying all
custom configurations, VLANs, and settings. This is a destructive operation
that will:
- Remove all VLANs and network configurations
- Reset interface configurations
- Clear all custom settings
- Restore factory default credentials

WARNING: This operation is irreversible and will cause network disruption.
Ensure all connected services are properly migrated before performing this reset.

Arguments:
  network_device_id   The unique identifier of the network device to reset

Examples:
  # Reset device to factory defaults
  metalcloud-cli network-device reset 12345

  # Confirm the operation is intentional
  metalcloud-cli switch reset 12345

```
metalcloud-cli network-device reset <network_device_id> [flags]
```

### Options

```
  -h, --help   help for reset
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

