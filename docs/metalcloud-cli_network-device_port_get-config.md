## metalcloud-cli network-device port get-config

Get the staged configuration of a network device port

### Synopsis

Display the staged (desired) configuration of a network device port, including
the admin overrides for description, MTU, enabled state and speed, plus the
optimistic-lock revision of the configuration buffer.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  port_id             The numeric interface id of the port

Examples:
  # Show the staged config of interface 42
  metalcloud-cli network-device port get-config 12345 42

```
metalcloud-cli network-device port get-config <network_device_id> <port_id> [flags]
```

### Options

```
  -h, --help   help for get-config
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

