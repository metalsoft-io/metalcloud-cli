## metalcloud-cli network-configuration bgp-interconnect-template get

Get detailed information about a specific network device BGP interconnect configuration template

### Synopsis

Display detailed information about a specific network device BGP interconnect configuration template.

Arguments:
  network_device_bgp_interconnect_configuration_template_id   The unique identifier of the network device BGP interconnect configuration template

Examples:
  # Get details for network device BGP interconnect configuration template with ID 12345
  metalcloud-cli network-configuration bgp-interconnect-template get 12345
  # Using alias
  metalcloud-cli network-configuration bgp-interconnect-template show 12345

```
metalcloud-cli network-configuration bgp-interconnect-template get <network_device_bgp_interconnect_configuration_template_id> [flags]
```

### Options

```
  -h, --help   help for get
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

