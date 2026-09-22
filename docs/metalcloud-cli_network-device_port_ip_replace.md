## metalcloud-cli network-device port ip replace

Replace the whole IP address set of a network device port

### Synopsis

Replace the complete IP address set of one address family on a network device
port with the supplied list. Addresses missing from the list are removed.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  port_id             The numeric interface id of the port
  family              Address family: ipv4 or ipv6

Required Flags:
  --config-source     Source of the address list
                      Values: 'pipe' for stdin input, or path to a JSON/YAML file

Examples:
  # Replace the IPv4 addresses of interface 42
  metalcloud-cli network-device port ip replace 12345 42 ipv4 --config-source ips.json

  # Replace them from pipe input
  echo '{"ips":[{"address":"10.0.0.1","prefixLength":32}]}' | metalcloud-cli network-device port ip replace 12345 42 ipv4 --config-source pipe

```
metalcloud-cli network-device port ip replace <network_device_id> <port_id> <family> [flags]
```

### Options

```
      --config-source string   Source of the address list. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for replace
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

