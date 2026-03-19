## metalcloud-cli network-configuration link-aggregation-template

Manage network devices link aggregation configuration templates

### Synopsis

Network device link aggregation configuration template commands.

Network device link aggregation configuration templates are used to deploy link aggregation configurations to network devices
Available commands:
  list                List all available Network device link aggregation configuration templates
  get                 Show detailed information about a specific template
  create              Create a new Network device link aggregation configuration template from JSON configuration
  update              Update an existing Network device link aggregation configuration template
  delete              Delete a Network device link aggregation configuration template

### Options

```
  -h, --help   help for link-aggregation-template
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
* [metalcloud-cli network-configuration link-aggregation-template config-example](metalcloud-cli_network-configuration_link-aggregation-template_config-example.md)	 - Generate example configuration template for network device link aggregation configuration template
* [metalcloud-cli network-configuration link-aggregation-template create](metalcloud-cli_network-configuration_link-aggregation-template_create.md)	 - Create a new network device link aggregation configuration template with specified configuration
* [metalcloud-cli network-configuration link-aggregation-template delete](metalcloud-cli_network-configuration_link-aggregation-template_delete.md)	 - Delete a network device link aggregation configuration template from the system
* [metalcloud-cli network-configuration link-aggregation-template get](metalcloud-cli_network-configuration_link-aggregation-template_get.md)	 - Get detailed information about a specific network device link aggregation configuration template
* [metalcloud-cli network-configuration link-aggregation-template list](metalcloud-cli_network-configuration_link-aggregation-template_list.md)	 - List network device link aggregation configuration templates with optional filtering
* [metalcloud-cli network-configuration link-aggregation-template update](metalcloud-cli_network-configuration_link-aggregation-template_update.md)	 - Update configuration of an existing network device link aggregation configuration template

