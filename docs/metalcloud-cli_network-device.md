## metalcloud-cli network-device

Manage network devices (switches) in the infrastructure

### Synopsis

Network device management commands for switches and other network infrastructure.

Network devices are physical switches that connect servers and provide network connectivity
within the MetalSoft infrastructure. These commands allow you to manage, configure, and
monitor network devices.

### Options

```
  -h, --help   help for network-device
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

* [metalcloud-cli](metalcloud-cli.md)	 - MetalCloud CLI
* [metalcloud-cli network-device archive](metalcloud-cli_network-device_archive.md)	 - Archive a network device (soft delete with history preservation)
* [metalcloud-cli network-device change-status](metalcloud-cli_network-device_change-status.md)	 - Change the operational status of a network device
* [metalcloud-cli network-device config-example](metalcloud-cli_network-device_config-example.md)	 - Generate example configuration template for network devices
* [metalcloud-cli network-device create](metalcloud-cli_network-device_create.md)	 - Create a new network device with specified configuration
* [metalcloud-cli network-device delete](metalcloud-cli_network-device_delete.md)	 - Delete a network device from the infrastructure
* [metalcloud-cli network-device discover](metalcloud-cli_network-device_discover.md)	 - Discover and inventory network device interfaces and configuration
* [metalcloud-cli network-device enable-syslog](metalcloud-cli_network-device_enable-syslog.md)	 - Enable remote syslog forwarding on the network device
* [metalcloud-cli network-device get](metalcloud-cli_network-device_get.md)	 - Get detailed information about a specific network device
* [metalcloud-cli network-device get-credentials](metalcloud-cli_network-device_get-credentials.md)	 - Retrieve management credentials for a network device
* [metalcloud-cli network-device get-defaults](metalcloud-cli_network-device_get-defaults.md)	 - Get default network device configuration settings for a site
* [metalcloud-cli network-device get-ports](metalcloud-cli_network-device_get-ports.md)	 - Get real-time port statistics directly from the network device
* [metalcloud-cli network-device list](metalcloud-cli_network-device_list.md)	 - List network devices with optional status filtering
* [metalcloud-cli network-device reset](metalcloud-cli_network-device_reset.md)	 - Reset network device to factory defaults (destructive operation)
* [metalcloud-cli network-device set-port-status](metalcloud-cli_network-device_set-port-status.md)	 - Enable or disable a specific port on the network device
* [metalcloud-cli network-device update](metalcloud-cli_network-device_update.md)	 - Update configuration of an existing network device

