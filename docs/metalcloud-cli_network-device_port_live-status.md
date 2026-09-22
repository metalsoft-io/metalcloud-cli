## metalcloud-cli network-device port live-status

Query the device for the operational state of its ports

### Synopsis

Query the network device itself for the operational state of its physical ports:
admin/link state, negotiated speed and duplex, and traffic utilization.

This differs from 'network-device get-ports', which lists the interface
inventory stored by MetalSoft without contacting the device.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # Show the live port state of device 12345
  metalcloud-cli network-device port live-status 12345

```
metalcloud-cli network-device port live-status <network_device_id> [flags]
```

### Options

```
  -h, --help   help for live-status
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

* [metalcloud-cli network-device port](metalcloud-cli_network-device_port.md)	 - Manage the interfaces of a network device

