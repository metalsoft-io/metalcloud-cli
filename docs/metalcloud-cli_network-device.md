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
* [metalcloud-cli network-device add-port-ip](metalcloud-cli_network-device_add-port-ip.md)	 - Add an IP address to a network device port
* [metalcloud-cli network-device archive](metalcloud-cli_network-device_archive.md)	 - Archive a network device (soft delete with history preservation)
* [metalcloud-cli network-device breakout](metalcloud-cli_network-device_breakout.md)	 - Manage the port breakouts of a network device
* [metalcloud-cli network-device config-example](metalcloud-cli_network-device_config-example.md)	 - Generate example configuration template for network devices
* [metalcloud-cli network-device create](metalcloud-cli_network-device_create.md)	 - Create a new network device with specified configuration
* [metalcloud-cli network-device create-bulk](metalcloud-cli_network-device_create-bulk.md)	 - Create multiple network devices in a single operation
* [metalcloud-cli network-device delete](metalcloud-cli_network-device_delete.md)	 - Delete a network device from the infrastructure
* [metalcloud-cli network-device disable-syslog](metalcloud-cli_network-device_disable-syslog.md)	 - Unsubscribe a network device from remote syslog
* [metalcloud-cli network-device discover](metalcloud-cli_network-device_discover.md)	 - Discover and inventory network device interfaces and configuration
* [metalcloud-cli network-device drift](metalcloud-cli_network-device_drift.md)	 - Inspect the configuration drift of a network device
* [metalcloud-cli network-device driver-capabilities](metalcloud-cli_network-device_driver-capabilities.md)	 - List the onboarding actions supported by each driver
* [metalcloud-cli network-device enable-syslog](metalcloud-cli_network-device_enable-syslog.md)	 - Enable remote syslog forwarding on the network device
* [metalcloud-cli network-device get](metalcloud-cli_network-device_get.md)	 - Get detailed information about a specific network device
* [metalcloud-cli network-device get-credentials](metalcloud-cli_network-device_get-credentials.md)	 - Retrieve management credentials for a network device
* [metalcloud-cli network-device get-ports](metalcloud-cli_network-device_get-ports.md)	 - List the interface inventory of a network device
* [metalcloud-cli network-device health-summary](metalcloud-cli_network-device_health-summary.md)	 - Show the health assessment of a network device
* [metalcloud-cli network-device list](metalcloud-cli_network-device_list.md)	 - List network devices with optional status filtering
* [metalcloud-cli network-device mark-installation-ready](metalcloud-cli_network-device_mark-installation-ready.md)	 - Mark the physical installation of a network device as done
* [metalcloud-cli network-device port](metalcloud-cli_network-device_port.md)	 - Manage the interfaces of a network device
* [metalcloud-cli network-device re-provision](metalcloud-cli_network-device_re-provision.md)	 - Re-run provisioning on a network device
* [metalcloud-cli network-device replace](metalcloud-cli_network-device_replace.md)	 - Replace a network device with another one
* [metalcloud-cli network-device reset](metalcloud-cli_network-device_reset.md)	 - Reset network device to factory defaults (destructive operation)
* [metalcloud-cli network-device return-to-planned](metalcloud-cli_network-device_return-to-planned.md)	 - Move an archived network device back to planned
* [metalcloud-cli network-device return-to-planned-config-example](metalcloud-cli_network-device_return-to-planned-config-example.md)	 - Example configuration for the return-to-planned command
* [metalcloud-cli network-device revert-failed-state](metalcloud-cli_network-device_revert-failed-state.md)	 - Take a network device out of the failed state
* [metalcloud-cli network-device run-extension](metalcloud-cli_network-device_run-extension.md)	 - Run an extension against a network device
* [metalcloud-cli network-device run-extension-config-example](metalcloud-cli_network-device_run-extension-config-example.md)	 - Example configuration for the run-extension command
* [metalcloud-cli network-device secret](metalcloud-cli_network-device_secret.md)	 - Manage the secrets of a network device
* [metalcloud-cli network-device set-failed](metalcloud-cli_network-device_set-failed.md)	 - Set the network device as failed
* [metalcloud-cli network-device set-health-monitoring-filter](metalcloud-cli_network-device_set-health-monitoring-filter.md)	 - Set the network device filter of a health monitoring socket
* [metalcloud-cli network-device set-port-status](metalcloud-cli_network-device_set-port-status.md)	 - Enable or disable a specific port on the network device
* [metalcloud-cli network-device snapshot](metalcloud-cli_network-device_snapshot.md)	 - Inspect the configuration snapshots of a network device
* [metalcloud-cli network-device snmp-monitoring](metalcloud-cli_network-device_snmp-monitoring.md)	 - Manage the SNMP monitoring subscription of network devices
* [metalcloud-cli network-device snmp-service](metalcloud-cli_network-device_snmp-service.md)	 - Manage the SNMP agent running on a network device
* [metalcloud-cli network-device start-registration](metalcloud-cli_network-device_start-registration.md)	 - Start the onboarding registration of a network device
* [metalcloud-cli network-device statistics](metalcloud-cli_network-device_statistics.md)	 - Show global network device counters
* [metalcloud-cli network-device sync-target-snapshot](metalcloud-cli_network-device_sync-target-snapshot.md)	 - Accept the current configuration as the drift target
* [metalcloud-cli network-device update](metalcloud-cli_network-device_update.md)	 - Update configuration of an existing network device
* [metalcloud-cli network-device update-port-config](metalcloud-cli_network-device_update-port-config.md)	 - Update the staged config (enable/description) of a port
* [metalcloud-cli network-device vendor](metalcloud-cli_network-device_vendor.md)	 - Manage the network device vendor profiles
* [metalcloud-cli network-device virtual-function](metalcloud-cli_network-device_virtual-function.md)	 - Inspect the virtual functions of a network device

