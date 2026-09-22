## metalcloud-cli network-device snmp-monitoring disable-batch

Unsubscribe several network devices from SNMP monitoring

### Synopsis

Unsubscribe several network devices from SNMP monitoring in a single call.
Devices are selected by id (positional arguments), by site or by fabric; at
least one selection must be given.

Optional Arguments:
  network_device_id...   Numeric ids of the network devices to unsubscribe

Optional Flags:
  --site-id     Numeric site ids whose devices are unsubscribed (repeatable)
  --fabric-id   Numeric fabric ids whose devices are unsubscribed (repeatable)

Examples:
  # Unsubscribe three devices
  metalcloud-cli network-device snmp-monitoring disable-batch 12345 12346 12347

  # Unsubscribe every device of fabric 4
  metalcloud-cli network-device snmp-monitoring disable-batch --fabric-id 4

```
metalcloud-cli network-device snmp-monitoring disable-batch [network_device_id...] [flags]
```

### Options

```
      --fabric-id strings   Numeric fabric ids whose network devices are unsubscribed. Repeatable.
  -h, --help                help for disable-batch
      --site-id strings     Numeric site ids whose network devices are unsubscribed. Repeatable.
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

