## metalcloud-cli network-device breakout get

Get one breakout of a network device

### Synopsis

Display one port breakout of a network device, including the applied breakout
groups and the staged configuration.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  breakout_id         The numeric id of the breakout

Examples:
  # Show breakout 7 of device 12345
  metalcloud-cli network-device breakout get 12345 7

```
metalcloud-cli network-device breakout get <network_device_id> <breakout_id> [flags]
```

### Options

```
  -h, --help   help for get
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

* [metalcloud-cli network-device breakout](metalcloud-cli_network-device_breakout.md)	 - Manage the port breakouts of a network device

