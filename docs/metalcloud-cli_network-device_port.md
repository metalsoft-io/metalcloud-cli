## metalcloud-cli network-device port

Manage the interfaces of a network device

### Synopsis

Manage the interfaces (ports) of a network device.

Ports are addressed by their numeric interface id, as shown by
'metalcloud-cli network-device get-ports <network_device_id>'.

### Options

```
  -h, --help   help for port
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
* [metalcloud-cli network-device port add](metalcloud-cli_network-device_port_add.md)	 - Create a logical interface on a network device
* [metalcloud-cli network-device port config-example](metalcloud-cli_network-device_port_config-example.md)	 - Example configuration for creating a network device interface
* [metalcloud-cli network-device port get-config](metalcloud-cli_network-device_port_get-config.md)	 - Get the staged configuration of a network device port
* [metalcloud-cli network-device port ip](metalcloud-cli_network-device_port_ip.md)	 - Manage the IP addresses of a network device port
* [metalcloud-cli network-device port live-status](metalcloud-cli_network-device_port_live-status.md)	 - Query the device for the operational state of its ports
* [metalcloud-cli network-device port remove](metalcloud-cli_network-device_port_remove.md)	 - Delete an interface of a network device
* [metalcloud-cli network-device port virtual-function](metalcloud-cli_network-device_port_virtual-function.md)	 - Inspect the virtual functions of a network device port

