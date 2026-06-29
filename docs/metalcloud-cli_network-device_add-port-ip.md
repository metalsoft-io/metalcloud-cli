## metalcloud-cli network-device add-port-ip

Add an IP address to a network device port

### Synopsis

Stage a new IP address on a network device port, addressed by its numeric
interface id. Typically used to assign a /32 loopback address to the loopback
interface.

Arguments:
  network_device_id   The unique identifier of the network device

Required Flags:
  --port-id           Numeric interface id of the port
  --address           IP address to add (without prefix, e.g. "10.253.128.1")
  --prefix            Prefix length (e.g. 32 for a loopback /32)

Optional Flags:
  --family            Address family: ipv4 (default) or ipv6

Examples:
  # Add a /32 loopback address
  metalcloud-cli network-device add-port-ip 12345 --port-id 67890 --address 10.253.128.1 --prefix 32

  # Add an IPv6 address
  metalcloud-cli nd add-port-ip 12345 --port-id 67890 --family ipv6 --address 2001:db8::1 --prefix 128

```
metalcloud-cli network-device add-port-ip <network_device_id> [flags]
```

### Options

```
      --address string   IP address to add (without prefix).
      --family string    Address family: ipv4 or ipv6. (default "ipv4")
  -h, --help             help for add-port-ip
      --port-id string   Numeric interface id of the port.
      --prefix int32     Prefix length (e.g. 32 for a loopback /32).
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

