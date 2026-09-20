## metalcloud-cli network-device breakout update-config

Stage a new breakout group set on a breakout

### Synopsis

Stage a new set of breakout groups on an existing breakout. The staged groups
are applied on the next deploy.

Required Arguments:
  network_device_id   The numeric id or label of the network device
  breakout_id         The numeric id of the breakout

Required Flags:
  --config-source     Source of the breakout group configuration
                      Values: 'pipe' for stdin input, or path to a JSON/YAML file

Examples:
  # Stage a new group set from a file
  metalcloud-cli network-device breakout update-config 12345 7 --config-source groups.json

  # Stage a 2x50G split from pipe input
  echo '{"breakoutGroups":[{"numberOfInterfaces":2,"speed":"50G"}]}' | metalcloud-cli network-device breakout update-config 12345 7 --config-source pipe

```
metalcloud-cli network-device breakout update-config <network_device_id> <breakout_id> [flags]
```

### Options

```
      --config-source string   Source of the breakout group configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for update-config
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

