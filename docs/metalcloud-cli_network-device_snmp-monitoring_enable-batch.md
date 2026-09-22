## metalcloud-cli network-device snmp-monitoring enable-batch

Subscribe several network devices to SNMP monitoring

### Synopsis

Subscribe several network devices to SNMP monitoring in a single call. Devices
are selected by id (positional arguments), by site or by fabric; at least one
selection must be given.

Optional Arguments:
  network_device_id...   Numeric ids of the network devices to subscribe

Optional Flags:
  --site-id     Numeric site ids whose devices are subscribed (repeatable)
  --fabric-id   Numeric fabric ids whose devices are subscribed (repeatable)

Examples:
  # Subscribe three devices
  metalcloud-cli network-device snmp-monitoring enable-batch 12345 12346 12347

  # Subscribe every device of site 1
  metalcloud-cli network-device snmp-monitoring enable-batch --site-id 1

```
metalcloud-cli network-device snmp-monitoring enable-batch [network_device_id...] [flags]
```

### Options

```
      --fabric-id strings   Numeric fabric ids whose network devices are subscribed. Repeatable.
  -h, --help                help for enable-batch
      --site-id strings     Numeric site ids whose network devices are subscribed. Repeatable.
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

