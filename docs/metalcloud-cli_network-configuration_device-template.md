## metalcloud-cli network-configuration device-template

Manage network devices configuration templates

### Synopsis

Network device configuration template commands.

Network device configuration templates are used to deploy configurations to network devices
Available commands:
  list                List all available Network device configuration templates
  get                 Show detailed information about a specific template
  create              Create a new Network device configuration template from JSON configuration
  update              Update an existing Network device configuration template
  delete              Delete a Network device configuration template

### Options

```
  -h, --help   help for device-template
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
* [metalcloud-cli network-configuration device-template config-example](metalcloud-cli_network-configuration_device-template_config-example.md)	 - Generate example configuration template for network device configuration template
* [metalcloud-cli network-configuration device-template create](metalcloud-cli_network-configuration_device-template_create.md)	 - Create a new network device configuration template with specified configuration
* [metalcloud-cli network-configuration device-template delete](metalcloud-cli_network-configuration_device-template_delete.md)	 - Delete a network device configuration template from the system
* [metalcloud-cli network-configuration device-template get](metalcloud-cli_network-configuration_device-template_get.md)	 - Get detailed information about a specific network device configuration template
* [metalcloud-cli network-configuration device-template list](metalcloud-cli_network-configuration_device-template_list.md)	 - List network device configuration templates with optional filtering
* [metalcloud-cli network-configuration device-template update](metalcloud-cli_network-configuration_device-template_update.md)	 - Update configuration of an existing network device configuration template

