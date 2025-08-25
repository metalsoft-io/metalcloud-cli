## metalcloud-cli network-device change-status

Change the operational status of a network device

### Synopsis

Change the operational status of a network device in the management system.
This affects how the system treats the device for operations and monitoring.

Arguments:
  network_device_id   The unique identifier of the network device
  status             New operational status for the device
                     Values: active, inactive, maintenance, error

Status descriptions:
  active       - Device is operational and available for use
  inactive     - Device is present but not operational
  maintenance  - Device is under maintenance, avoid new allocations
  error        - Device has issues and requires attention

Examples:
  # Put device in maintenance mode
  metalcloud-cli network-device change-status 12345 maintenance

  # Activate device after maintenance
  metalcloud-cli network-device change-status 12345 active

  # Mark device as having errors
  metalcloud-cli switch change-status 12345 error

```
metalcloud-cli network-device change-status <network_device_id> <status> [flags]
```

### Options

```
  -h, --help   help for change-status
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

