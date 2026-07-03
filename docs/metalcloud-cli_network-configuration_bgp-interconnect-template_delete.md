## metalcloud-cli network-configuration bgp-interconnect-template delete

Delete a network device BGP interconnect configuration template from the system

### Synopsis

Delete a network device BGP interconnect configuration template from the system.

Arguments:
  network_device_bgp_interconnect_configuration_template_id   The unique identifier of the network device BGP interconnect configuration template to delete

Examples:
  # Delete network device BGP interconnect configuration template
  metalcloud-cli network-configuration bgp-interconnect-template delete 12345

  # Using alias
  metalcloud-cli nc bit rm 12345

```
metalcloud-cli network-configuration bgp-interconnect-template delete <network_device_bgp_interconnect_configuration_template_id> [flags]
```

### Options

```
  -h, --help   help for delete
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

