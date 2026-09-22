## metalcloud-cli network-device breakout create

Create a port breakout on a network device

### Synopsis

Create a port breakout on a network device, splitting one physical port into
several child interfaces.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Required Flags:
  --config-source     Source of the breakout configuration
                      Values: 'pipe' for stdin input, or path to a JSON/YAML file

Examples:
  # Create a breakout from a file
  metalcloud-cli network-device breakout create 12345 --config-source breakout.json

  # Split Ethernet0 into 4x25G from pipe input
  echo '{"portName":"Ethernet0","groups":[{"numberOfInterfaces":4,"speed":"25G"}]}' | metalcloud-cli network-device breakout create 12345 --config-source pipe

```
metalcloud-cli network-device breakout create <network_device_id> [flags]
```

### Options

```
      --config-source string   Source of the new breakout configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for create
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

