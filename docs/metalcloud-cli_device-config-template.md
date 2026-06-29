## metalcloud-cli device-config-template

Manage device configuration templates and profiles

### Synopsis

Device configuration template commands.

Device configuration templates hold renderable configuration content for network
devices, and profiles bind those templates to specific devices or fabrics.

Available commands:
  list                List device configuration templates
  get                 Show details about a specific template
  create              Create a new template from JSON configuration
  update              Update an existing template
  delete              Delete a template
  render              Render arbitrary template content
  render-saved        Render a saved template by ID
  config-example      Generate an example template configuration
  profile             Manage device configuration template profiles

### Options

```
  -h, --help   help for device-config-template
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

* [metalcloud-cli](metalcloud-cli.md)	 - MetalCloud CLI
* [metalcloud-cli device-config-template config-example](metalcloud-cli_device-config-template_config-example.md)	 - Generate an example device configuration template
* [metalcloud-cli device-config-template create](metalcloud-cli_device-config-template_create.md)	 - Create a new device configuration template
* [metalcloud-cli device-config-template delete](metalcloud-cli_device-config-template_delete.md)	 - Delete a device configuration template
* [metalcloud-cli device-config-template get](metalcloud-cli_device-config-template_get.md)	 - Get detailed information about a specific device configuration template
* [metalcloud-cli device-config-template list](metalcloud-cli_device-config-template_list.md)	 - List device configuration templates with optional filtering
* [metalcloud-cli device-config-template profile](metalcloud-cli_device-config-template_profile.md)	 - Manage device configuration template profiles
* [metalcloud-cli device-config-template render](metalcloud-cli_device-config-template_render.md)	 - Render arbitrary device configuration template content
* [metalcloud-cli device-config-template render-saved](metalcloud-cli_device-config-template_render-saved.md)	 - Render a saved device configuration template by ID
* [metalcloud-cli device-config-template update](metalcloud-cli_device-config-template_update.md)	 - Update an existing device configuration template

