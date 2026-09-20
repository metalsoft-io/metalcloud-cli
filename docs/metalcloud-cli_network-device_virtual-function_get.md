## metalcloud-cli network-device virtual-function get

Get one virtual function of a network device

### Synopsis

Display one virtual function of a network device.

Required Arguments:
  network_device_id     The numeric id or label of the network device
  virtual_function_id   The numeric id of the virtual function

Examples:
  # Show virtual function 3 of device 12345
  metalcloud-cli network-device virtual-function get 12345 3

```
metalcloud-cli network-device virtual-function get <network_device_id> <virtual_function_id> [flags]
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

* [metalcloud-cli network-device virtual-function](metalcloud-cli_network-device_virtual-function.md)	 - Inspect the virtual functions of a network device

