## metalcloud-cli network-device drift list

List the configuration drift of a network device

### Synopsis

List the configuration drift entries recorded for a network device.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # List the drift history of device 12345
  metalcloud-cli network-device drift list 12345

```
metalcloud-cli network-device drift list <network_device_id> [flags]
```

### Options

```
  -h, --help   help for list
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

