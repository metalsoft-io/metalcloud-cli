## metalcloud-cli network-device snmp-monitoring enable

Subscribe a network device to SNMP monitoring

### Synopsis

Subscribe a network device to SNMP monitoring.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # Enable SNMP monitoring for device 12345
  metalcloud-cli network-device snmp-monitoring enable 12345

```
metalcloud-cli network-device snmp-monitoring enable <network_device_id> [flags]
```

### Options

```
  -h, --help   help for enable
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

* [metalcloud-cli network-device snmp-monitoring](metalcloud-cli_network-device_snmp-monitoring.md)	 - Manage the SNMP monitoring subscription of network devices

