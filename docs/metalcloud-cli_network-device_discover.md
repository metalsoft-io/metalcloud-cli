## metalcloud-cli network-device discover

Discover and inventory network device interfaces and configuration

### Synopsis

Initiate discovery process for a network device to automatically detect and
inventory its interfaces, hardware components, and software configuration.

This process connects to the device using its management interface and gathers
detailed information about:
- Physical interfaces and their status
- Hardware components and capabilities
- Software version and configuration
- VLAN and networking setup

Arguments:
  network_device_id   The unique identifier of the network device to discover

Examples:
  # Discover device interfaces and configuration
  metalcloud-cli network-device discover 12345

  # Using alias
  metalcloud-cli switch discover 12345

```
metalcloud-cli network-device discover <network_device_id> [flags]
```

### Options

```
  -h, --help   help for discover
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

