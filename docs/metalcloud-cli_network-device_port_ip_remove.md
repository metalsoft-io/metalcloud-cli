## metalcloud-cli network-device port ip remove

Remove one IP address from a network device port

### Synopsis

Remove one IP address from a network device port.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  port_id             The numeric interface id of the port
  family              Address family: ipv4 or ipv6
  ip_id               The numeric id of the address

Examples:
  # Remove IPv4 address 7 from interface 42
  metalcloud-cli network-device port ip remove 12345 42 ipv4 7

```
metalcloud-cli network-device port ip remove <network_device_id> <port_id> <family> <ip_id> [flags]
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

* [metalcloud-cli network-device port ip](metalcloud-cli_network-device_port_ip.md)	 - Manage the IP addresses of a network device port

