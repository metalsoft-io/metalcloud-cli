## metalcloud-cli network-device drift acknowledge

Acknowledge a configuration drift entry

### Synopsis

Mark one configuration drift entry of a network device as reviewed, recording
who acknowledged it and when.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  drift_id            The numeric id of the drift entry

Examples:
  # Acknowledge drift entry 9 of device 12345
  metalcloud-cli network-device drift acknowledge 12345 9

```
metalcloud-cli network-device drift acknowledge <network_device_id> <drift_id> [flags]
```

### Options

```
  -h, --help   help for acknowledge
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

