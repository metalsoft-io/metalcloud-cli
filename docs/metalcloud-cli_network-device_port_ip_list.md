## metalcloud-cli network-device port ip list

List the IP addresses of a network device port

### Synopsis

List the IP addresses of one address family staged on a network device port.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  port_id             The numeric interface id of the port
  family              Address family: ipv4 or ipv6

Examples:
  # List the IPv4 addresses of interface 42
  metalcloud-cli network-device port ip list 12345 42 ipv4

```
metalcloud-cli network-device port ip list <network_device_id> <port_id> <family> [flags]
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

* [metalcloud-cli network-device port ip](metalcloud-cli_network-device_port_ip.md)	 - Manage the IP addresses of a network device port

