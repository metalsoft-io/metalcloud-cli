## metalcloud-cli network-device get-ports

List the interface inventory of a network device

### Synopsis

List the interfaces of a network device as MetalSoft has them inventoried:
interface id, name, kind, description, MAC address, LAG membership and tags.

The interface ids reported here are the ones the 'network-device port'
sub-commands take. For the operational state read from the device itself
(link state, negotiated speed, utilization) use
'network-device port live-status' instead.

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

