## metalcloud-cli network-device sync-target-snapshot

Accept the current configuration as the drift target

### Synopsis

Point the drift detection target snapshot of a network device at its latest
snapshot, clearing the drift currently reported for the device.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # Accept the running configuration of device 12345 as the new target
  metalcloud-cli network-device sync-target-snapshot 12345

```
metalcloud-cli network-device sync-target-snapshot <network_device_id> [flags]
```

### Options

```
  -h, --help   help for sync-target-snapshot
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

