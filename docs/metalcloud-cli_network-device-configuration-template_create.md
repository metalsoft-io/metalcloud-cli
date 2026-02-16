## metalcloud-cli network-device-configuration-template create

Create a new network device configuration template with specified configuration

### Synopsis

Create a new network device configuration template using configuration provided via JSON file or pipe.

Required Flags:
  --config-source   Source of configuration data (required)
                   Values: 'pipe' for stdin input, or path to JSON file

Use the 'config-example' command to generate an example configuration:

  {
    "networkType": "underlay",
    "networkDeviceDriver": "cisco_aci51",
    "networkDevicePosition": "all",
    "remoteNetworkDevicePosition": "all",
    "bgpNumbering": "numbered",
    "bgpLinkConfiguration": "disabled",
    "executionType": "cli",
    "libraryLabel": "string",
    "preparation": "string",
    "configuration": "string"
  }

Note: Preparation and configuration fields need to be base64 encoded when submitted.

Examples:
  # Create template from JSON file
  metalcloud-cli network-device-configuration-template create --config-source template.json

  # Create template from pipe input
  cat template.json | metalcloud-cli network-device-configuration-template create --config-source pipe

  # Create template with inline JSON
  echo '{"networkDevicePosition":"all","remoteNetworkDevicePosition":"all"}' | metalcloud-cli ndct create --config-source pipe

```
metalcloud-cli network-device-configuration-template create [flags]
```

### Options

```
      --config-source string   Source of the new network device configuration template. Can be 'pipe' or path to a JSON file.
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

* [metalcloud-cli network-device-configuration-template](metalcloud-cli_network-device-configuration-template.md)	 - Manage network devices configuration templates

