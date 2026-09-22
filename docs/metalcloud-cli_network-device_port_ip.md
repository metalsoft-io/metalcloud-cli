## metalcloud-cli network-device port ip

Manage the IP addresses of a network device port

### Synopsis

Manage the IP addresses staged on a network device port.

Addresses are grouped by address family; the family argument is either
'ipv4' or 'ipv6'.

### Options

```
  -h, --help   help for ip
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
* [metalcloud-cli network-device port ip config-example](metalcloud-cli_network-device_port_ip_config-example.md)	 - Example address list for the port IP replace command
* [metalcloud-cli network-device port ip get](metalcloud-cli_network-device_port_ip_get.md)	 - Get one IP address of a network device port
* [metalcloud-cli network-device port ip list](metalcloud-cli_network-device_port_ip_list.md)	 - List the IP addresses of a network device port
* [metalcloud-cli network-device port ip remove](metalcloud-cli_network-device_port_ip_remove.md)	 - Remove one IP address from a network device port
* [metalcloud-cli network-device port ip replace](metalcloud-cli_network-device_port_ip_replace.md)	 - Replace the whole IP address set of a network device port

