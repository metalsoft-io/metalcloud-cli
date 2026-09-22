## metalcloud-cli network-device replace

Replace a network device with another one

### Synopsis

Replace a network device with a replacement device, moving its configuration
and connections over to the new device.

Required Arguments:
  network_device_id       The numeric id or label of the device being replaced
  new_network_device_id   The numeric id or label of the replacement device

Examples:
  # Replace device 12345 with device 12399
  metalcloud-cli network-device replace 12345 12399

```
metalcloud-cli network-device replace <network_device_id> <new_network_device_id> [flags]
```

### Options

```
  -h, --help   help for replace
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

