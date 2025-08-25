## metalcloud-cli network-device list

List network devices with optional status filtering

### Synopsis

List all network devices in the infrastructure with optional status filtering.

This command displays all network devices (switches) that are registered in the system.
You can filter the results by device status to focus on specific operational states.

Flags:
  --filter-status   Filter devices by operational status (default: ["active"])
                   Available statuses: active, inactive, maintenance, error, unknown

Examples:
  # List all active network devices (default)
  metalcloud-cli network-device list

  # List devices in maintenance mode
  metalcloud-cli network-device list --filter-status maintenance

  # List devices with multiple statuses
  metalcloud-cli network-device list --filter-status active,maintenance

  # List all devices regardless of status
  metalcloud-cli network-device list --filter-status active,inactive,maintenance,error,unknown

```
metalcloud-cli network-device list [flags]
```

### Options

```
      --filter-status strings   Filter the result by network device status. (default [active])
  -h, --help                    help for list
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

