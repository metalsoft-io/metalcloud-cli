## metalcloud-cli network-device port virtual-function list

List the virtual functions of a network device port

### Synopsis

List the virtual functions (SR-IOV VFs) exposed by one network device port.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  port_id             The numeric interface id of the port

Examples:
  # List the virtual functions of interface 42
  metalcloud-cli network-device port virtual-function list 12345 42

```
metalcloud-cli network-device port virtual-function list <network_device_id> <port_id> [flags]
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

* [metalcloud-cli network-device port virtual-function](metalcloud-cli_network-device_port_virtual-function.md)	 - Inspect the virtual functions of a network device port

