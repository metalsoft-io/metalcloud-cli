## metalcloud-cli network-configuration bgp-interconnect-template

Manage network devices BGP interconnect configuration templates

### Synopsis

Network device BGP interconnect configuration template commands.

Network device BGP interconnect configuration templates are used to deploy BGP interconnect configurations to network devices
Available commands:
  list                List all available Network device BGP interconnect configuration templates
  get                 Show detailed information about a specific template
  create              Create a new Network device BGP interconnect configuration template from JSON configuration
  update              Update an existing Network device BGP interconnect configuration template
  delete              Delete a Network device BGP interconnect configuration template

### Options

```
  -h, --help   help for bgp-interconnect-template
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

* [metalcloud-cli network-configuration](metalcloud-cli_network-configuration.md)	 - Manage network configuration templates
* [metalcloud-cli network-configuration bgp-interconnect-template config-example](metalcloud-cli_network-configuration_bgp-interconnect-template_config-example.md)	 - Generate example configuration template for network device BGP interconnect configuration template
* [metalcloud-cli network-configuration bgp-interconnect-template create](metalcloud-cli_network-configuration_bgp-interconnect-template_create.md)	 - Create a new network device BGP interconnect configuration template with specified configuration
* [metalcloud-cli network-configuration bgp-interconnect-template delete](metalcloud-cli_network-configuration_bgp-interconnect-template_delete.md)	 - Delete a network device BGP interconnect configuration template from the system
* [metalcloud-cli network-configuration bgp-interconnect-template get](metalcloud-cli_network-configuration_bgp-interconnect-template_get.md)	 - Get detailed information about a specific network device BGP interconnect configuration template
* [metalcloud-cli network-configuration bgp-interconnect-template list](metalcloud-cli_network-configuration_bgp-interconnect-template_list.md)	 - List network device BGP interconnect configuration templates with optional filtering
* [metalcloud-cli network-configuration bgp-interconnect-template update](metalcloud-cli_network-configuration_bgp-interconnect-template_update.md)	 - Update configuration of an existing network device BGP interconnect configuration template

