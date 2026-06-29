## metalcloud-cli network-device create-bulk

Create multiple network devices in a single operation

### Synopsis

Create multiple network devices at once from a JSON or YAML configuration file.

This command processes a list of network device configurations and creates all devices in sequence.
Each device configuration follows the same format as the single network-device create command.

Required Flags:
  --config-source   Source of the bulk network device configuration (JSON/YAML file path or 'pipe')

Configuration File Format (JSON):
  [
    {
      "siteId": 1,
      "driver": "sonic_enterprise",
      "identifierString": "leaf-01",
      "position": "leaf",
      "managementAddress": "10.0.1.100",
      "managementPort": 22,
      "username": "admin",
      "managementPassword": "password"
    },
    {
      "siteId": 1,
      "driver": "sonic_enterprise",
      "identifierString": "spine-01",
      "position": "spine",
      "managementAddress": "10.0.1.101",
      "managementPort": 22,
      "username": "admin",
      "managementPassword": "password"
    }
  ]

The command will report success/failure for each device and provide a summary at the end.

```
metalcloud-cli network-device create-bulk [flags]
```

### Examples

```
  # Create network devices from JSON file
  metalcloud-cli network-device create-bulk --config-source devices.json

  # Create network devices from YAML file
  metalcloud-cli network-device create-bulk --config-source devices.yaml

  # Create network devices from pipe
  cat devices.json | metalcloud-cli network-device create-bulk --config-source pipe
```

### Options

```
      --config-source string   Source of the bulk network device configuration. Can be 'pipe' or path to a JSON/YAML file containing a list of devices.
  -h, --help                   help for create-bulk
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

