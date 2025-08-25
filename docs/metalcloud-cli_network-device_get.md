## metalcloud-cli network-device get

Get detailed information about a specific network device

### Synopsis

Display detailed information about a specific network device including its
configuration, status, interfaces, and operational details.

Arguments:
  network_device_id   The unique identifier of the network device

Examples:
  # Get details for network device with ID 12345
  metalcloud-cli network-device get 12345

  # Using alias
  metalcloud-cli switch show 12345

```
metalcloud-cli network-device get <network_device_id> [flags]
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

* [metalcloud-cli network-device](metalcloud-cli_network-device.md)	 - Manage network devices (switches) in the infrastructure

