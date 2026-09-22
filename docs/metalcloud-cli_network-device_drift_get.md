## metalcloud-cli network-device drift get

Get one configuration drift entry

### Synopsis

Display one configuration drift entry of a network device, including the
configuration difference that was detected.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  drift_id            The numeric id of the drift entry

Examples:
  # Show drift entry 9 of device 12345
  metalcloud-cli network-device drift get 12345 9

```
metalcloud-cli network-device drift get <network_device_id> <drift_id> [flags]
```

### Options

```
  -h, --help   help for get
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

* [metalcloud-cli network-device drift](metalcloud-cli_network-device_drift.md)	 - Inspect the configuration drift of a network device

