## metalcloud-cli network-device breakout list

List the breakouts of a network device

### Synopsis

List the port breakouts configured on a network device.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Examples:
  # List the breakouts of device 12345
  metalcloud-cli network-device breakout list 12345

```
metalcloud-cli network-device breakout list <network_device_id> [flags]
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

* [metalcloud-cli network-device breakout](metalcloud-cli_network-device_breakout.md)	 - Manage the port breakouts of a network device

