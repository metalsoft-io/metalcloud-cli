## metalcloud-cli network-device create

Create a new network device with specified configuration

### Synopsis

Create a new network device using configuration provided via JSON file or pipe.

The configuration must include device details such as management IP, credentials,
device type, and other operational parameters.

Required Flags:
  --config-source   Source of configuration data (required)
                   Values: 'pipe' for stdin input, or path to JSON file

Use the 'config-example' command to generate an example configuration:

  {
    "siteId": 1,
    "driver": "sonic_enterprise",
    "identifierString": "example",
    "serialNumber": "1234567890",
    "chassisIdentifier": "example",
    "chassisRackId": 1,
    "position": "leaf",
    "isGateway": false,
    "isStorageSwitch": false,
    "isBorderDevice": false,
    "managementMAC": "AA:BB:CC:DD:EE:FF",
    "managementAddress": "1.1.1.1",
    "managementAddressGateway": "1.1.1.1",
    "managementAddressMask": "255.255.255.0",
    "loopbackAddress": "127.0.0.1",
    "vtepAddress": null,
    "asn": 65000,
    "managementPort": 22,
    "username": "admin",
    "managementPassword": "password",
    "syslogEnabled": true
  }

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

