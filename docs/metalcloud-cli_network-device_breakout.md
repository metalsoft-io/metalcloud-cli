## metalcloud-cli network-device breakout

Manage the port breakouts of a network device

### Synopsis

Manage the port breakouts of a network device.

A breakout splits one physical port into several child interfaces. The applied
breakout groups are reported by 'get', while 'get-config' and 'update-config'
work on the groups staged for the next deploy.

### Options

```
  -h, --help   help for breakout
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
* [metalcloud-cli network-device breakout config-example](metalcloud-cli_network-device_breakout_config-example.md)	 - Example configuration for creating a breakout
* [metalcloud-cli network-device breakout create](metalcloud-cli_network-device_breakout_create.md)	 - Create a port breakout on a network device
* [metalcloud-cli network-device breakout delete](metalcloud-cli_network-device_breakout_delete.md)	 - Delete a breakout of a network device
* [metalcloud-cli network-device breakout get](metalcloud-cli_network-device_breakout_get.md)	 - Get one breakout of a network device
* [metalcloud-cli network-device breakout get-config](metalcloud-cli_network-device_breakout_get-config.md)	 - Get the staged configuration of a breakout
* [metalcloud-cli network-device breakout list](metalcloud-cli_network-device_breakout_list.md)	 - List the breakouts of a network device
* [metalcloud-cli network-device breakout update-config](metalcloud-cli_network-device_breakout_update-config.md)	 - Stage a new breakout group set on a breakout
* [metalcloud-cli network-device breakout update-config-example](metalcloud-cli_network-device_breakout_update-config-example.md)	 - Example configuration for staging breakout groups

