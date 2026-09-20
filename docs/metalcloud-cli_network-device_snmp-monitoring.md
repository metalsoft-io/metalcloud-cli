## metalcloud-cli network-device snmp-monitoring

Manage the SNMP monitoring subscription of network devices

### Synopsis

Subscribe network devices to, or unsubscribe them from, SNMP monitoring, and
inspect which monitoring agent polls each device.

### Options

```
  -h, --help   help for snmp-monitoring
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
* [metalcloud-cli network-device snmp-monitoring agent-info](metalcloud-cli_network-device_snmp-monitoring_agent-info.md)	 - Show which monitoring agent polls each network device
* [metalcloud-cli network-device snmp-monitoring disable](metalcloud-cli_network-device_snmp-monitoring_disable.md)	 - Unsubscribe a network device from SNMP monitoring
* [metalcloud-cli network-device snmp-monitoring disable-batch](metalcloud-cli_network-device_snmp-monitoring_disable-batch.md)	 - Unsubscribe several network devices from SNMP monitoring
* [metalcloud-cli network-device snmp-monitoring enable](metalcloud-cli_network-device_snmp-monitoring_enable.md)	 - Subscribe a network device to SNMP monitoring
* [metalcloud-cli network-device snmp-monitoring enable-batch](metalcloud-cli_network-device_snmp-monitoring_enable-batch.md)	 - Subscribe several network devices to SNMP monitoring

