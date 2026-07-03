## metalcloud-cli network-configuration bgp-interconnect-template update

Update configuration of an existing network device BGP interconnect configuration template

### Synopsis

Update the configuration of an existing network device BGP interconnect configuration template using JSON configuration
provided via file or pipe. Only the specified fields will be updated; other
configuration will remain unchanged.

Arguments:
  network_device_bgp_interconnect_configuration_template_id   The unique identifier of the network device BGP interconnect configuration template to update

Required Flags:
  --config-source   Source of configuration updates (required)
                   Values: 'pipe' for stdin input, or path to JSON file

Examples:
  # Update template from JSON file
  metalcloud-cli network-configuration bgp-interconnect-template update 12345 --config-source updates.json

  # Update template from pipe input
  cat updates.json | metalcloud-cli network-configuration bgp-interconnect-template update 12345 --config-source pipe

  # Update specific field
  echo '{"name":"new name"}' | metalcloud-cli nc bit update 12345 --config-source pipe

```
metalcloud-cli network-configuration bgp-interconnect-template update <network_device_bgp_interconnect_configuration_template_id> [flags]
```

### Options

```
      --config-source string   Source of the network device BGP interconnect configuration template updates. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for update
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

