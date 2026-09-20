## metalcloud-cli network-device revert-failed-state

Take a network device out of the failed state

### Synopsis

Take a network device out of the failed state and back to its previous status,
the counterpart of 'set-failed'.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # Revert the failed state of device 12345
  metalcloud-cli network-device revert-failed-state 12345

```
metalcloud-cli network-device revert-failed-state <network_device_id> [flags]
```

### Options

```
  -h, --help   help for revert-failed-state
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

