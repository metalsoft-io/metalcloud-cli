## metalcloud-cli network-device port add

Create a logical interface on a network device

### Synopsis

Create a logical interface (e.g. a loopback or a sub-interface) on a network device.

Required Arguments:
  network_device_id   The numeric id or label of the network device

Required Flags:
  --config-source     Source of the interface configuration
                      Values: 'pipe' for stdin input, or path to a JSON/YAML file

Examples:
  # Create a loopback interface from a file
  metalcloud-cli network-device port add 12345 --config-source loopback.json

  # Create an interface from pipe input
  echo '{"kind":"loopback","name":"Loopback1"}' | metalcloud-cli network-device port add 12345 --config-source pipe

```
metalcloud-cli network-device port add <network_device_id> [flags]
```

### Options

```
      --config-source string   Source of the new interface configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for add
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

