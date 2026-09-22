## metalcloud-cli network-device snmp-monitoring agent-info

Show which monitoring agent polls each network device

### Synopsis

Show the monitoring agent allocated to each of the selected network devices.
Devices are selected by id (positional arguments), by site or by fabric; with
no selection the whole allocation map is returned.

Optional Arguments:
  network_device_id...   Numeric ids of the network devices to look up

Optional Flags:
  --site-id     Restrict the lookup to these numeric site ids (repeatable)
  --fabric-id   Restrict the lookup to these numeric fabric ids (repeatable)

Examples:
  # Show the agents of two devices
  metalcloud-cli network-device snmp-monitoring agent-info 12345 12346

  # Show the agents of every device in site 1
  metalcloud-cli network-device snmp-monitoring agent-info --site-id 1

```
metalcloud-cli network-device snmp-monitoring agent-info [network_device_id...] [flags]
```

### Options

```
      --fabric-id strings   Restrict the lookup to these numeric fabric ids. Repeatable.
  -h, --help                help for agent-info
      --site-id strings     Restrict the lookup to these numeric site ids. Repeatable.
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

