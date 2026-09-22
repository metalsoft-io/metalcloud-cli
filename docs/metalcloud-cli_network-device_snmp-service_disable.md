## metalcloud-cli network-device snmp-service disable

Disable the SNMP agent on a network device

### Synopsis

Disable the SNMP agent running on the network device.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # Disable the SNMP agent of device 12345
  metalcloud-cli network-device snmp-service disable 12345

```
metalcloud-cli network-device snmp-service disable <network_device_id> [flags]
```

### Options

```
  -h, --help   help for disable
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

* [metalcloud-cli network-device snmp-service](metalcloud-cli_network-device_snmp-service.md)	 - Manage the SNMP agent running on a network device

