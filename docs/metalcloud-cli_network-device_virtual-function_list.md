## metalcloud-cli network-device virtual-function list

List all virtual functions of a network device

### Synopsis

List every virtual function of a network device, across all of its interfaces.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # List the virtual functions of device 12345
  metalcloud-cli network-device virtual-function list 12345

```
metalcloud-cli network-device virtual-function list <network_device_id> [flags]
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

* [metalcloud-cli network-device virtual-function](metalcloud-cli_network-device_virtual-function.md)	 - Inspect the virtual functions of a network device

