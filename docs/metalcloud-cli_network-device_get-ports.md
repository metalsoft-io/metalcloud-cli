## metalcloud-cli network-device get-ports

Get real-time port statistics directly from the network device

### Synopsis

Retrieve real-time port statistics and status information directly from
the network device. This provides current operational data including:
- Port status (up/down)
- Traffic statistics (bytes, packets)
- Error counters
- Link speed and duplex settings

This data is fetched directly from the device rather than cached information.

Arguments:
  network_device_id   The unique identifier of the network device

Examples:
  # Get current port statistics
  metalcloud-cli network-device get-ports 12345

  # Using alias
  metalcloud-cli switch get-ports 12345

```
metalcloud-cli network-device get-ports <network_device_id> [flags]
```

### Options

```
  -h, --help   help for get-ports
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

