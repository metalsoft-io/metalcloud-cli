## metalcloud-cli network-device create

Create a new network device with specified configuration

### Synopsis

Create a new network device using configuration provided via JSON file or pipe.

The configuration must include device details such as management IP, credentials,
device type, and other operational parameters.

Required Flags:
  --config-source   Source of configuration data (required)
                   Values: 'pipe' for stdin input, or path to JSON file

Examples:
  # Create device from JSON file
  metalcloud-cli network-device create --config-source device-config.json

  # Create device from pipe input
  cat device-config.json | metalcloud-cli network-device create --config-source pipe

  # Create device with inline JSON
  echo '{"management_ip":"10.0.1.100","type":"cisco"}' | metalcloud-cli nd create --config-source pipe

```
metalcloud-cli network-device create [flags]
```

### Options

```
      --config-source string   Source of the new network device configuration. Can be 'pipe' or path to a JSON file.
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

* [metalcloud-cli network-device](metalcloud-cli_network-device.md)	 - Manage network devices (switches) in the infrastructure

