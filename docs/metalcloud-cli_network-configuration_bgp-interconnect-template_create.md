## metalcloud-cli network-configuration bgp-interconnect-template create

Create a new network device BGP interconnect configuration template with specified configuration

### Synopsis

Create a new network device BGP interconnect configuration template using configuration provided via JSON file or pipe.

Required Flags:
  --config-source   Source of configuration data (required)
                   Values: 'pipe' for stdin input, or path to JSON file

Use the 'config-example' command to generate an example configuration:

  {
    "label": "string",
    "name": "string",
    "networkDeviceDriver": "junos",
    "executionType": "cli",
    "addGlobalConfig": "string",
    "removeGlobalConfig": "string",
    "addNeighbor": "string",
    "removeNeighbor": "string"
  }

Note: addGlobalConfig, removeGlobalConfig, addNeighbor and removeNeighbor fields need to be base64 encoded when submitted.

Examples:
  # Create template from JSON file
  metalcloud-cli network-configuration bgp-interconnect-template create --config-source template.json

  # Create template from pipe input
  cat template.json | metalcloud-cli network-configuration bgp-interconnect-template create --config-source pipe

  # Create template with inline JSON
  echo '{"label":"l","name":"n","networkDeviceDriver":"junos","executionType":"cli"}' | metalcloud-cli nc bit create --config-source pipe

```
metalcloud-cli network-configuration bgp-interconnect-template create [flags]
```

### Options

```
      --config-source string   Source of the new network device BGP interconnect configuration template. Can be 'pipe' or path to a JSON file.
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

* [metalcloud-cli network-configuration bgp-interconnect-template](metalcloud-cli_network-configuration_bgp-interconnect-template.md)	 - Manage network devices BGP interconnect configuration templates

