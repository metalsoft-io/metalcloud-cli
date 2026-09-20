## metalcloud-cli network-device port remove

Delete an interface of a network device

### Synopsis

Delete a logical interface of a network device.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  port_id             The numeric interface id of the port

Examples:
  # Delete interface 42 of device 12345
  metalcloud-cli network-device port remove 12345 42

```
metalcloud-cli network-device port remove <network_device_id> <port_id> [flags]
```

### Options

```
  -h, --help   help for remove
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

* [metalcloud-cli network-device port](metalcloud-cli_network-device_port.md)	 - Manage the interfaces of a network device

