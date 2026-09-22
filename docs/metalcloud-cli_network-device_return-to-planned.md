## metalcloud-cli network-device return-to-planned

Move an archived network device back to planned

### Synopsis

Move an archived network device back to the planned state so that it can be
installed again.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Optional Flags:
  --config-source     Source of the return configuration (management MAC, serial
                      number, OS template override)
                      Values: 'pipe' for stdin input, or path to a JSON/YAML file

Examples:
  # Return device 12345 to planned
  metalcloud-cli network-device return-to-planned 12345

  # Return it supplying the identifying MAC address
  echo '{"managementMAC":"AA:BB:CC:DD:EE:FF"}' | metalcloud-cli network-device return-to-planned 12345 --config-source pipe

```
metalcloud-cli network-device return-to-planned <network_device_id> [flags]
```

### Options

```
      --config-source string   Source of the return configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for return-to-planned
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

