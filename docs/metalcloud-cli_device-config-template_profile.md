## metalcloud-cli device-config-template profile

Manage device configuration template profiles

### Synopsis

Device configuration template profile commands.

Profiles bind a device configuration template to a specific network device or fabric,
with variables and lifecycle/apply settings.

Available commands:
  list                List profiles
  get                 Show details about a specific profile
  create              Create a new profile from JSON configuration
  update              Update an existing profile
  delete              Delete a profile
  render              Render a profile for a given device
  find-applicable     Find profiles applicable to a device/fabric
  render-applicable   Render profiles applicable to a device/fabric
  bulk-assign         Bulk-assign a template to multiple devices
  config-example      Generate an example profile configuration

### Options

```
  -h, --help   help for profile
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

* [metalcloud-cli device-config-template](metalcloud-cli_device-config-template.md)	 - Manage device configuration templates and profiles
* [metalcloud-cli device-config-template profile bulk-assign](metalcloud-cli_device-config-template_profile_bulk-assign.md)	 - Bulk-assign a device configuration template to multiple devices
* [metalcloud-cli device-config-template profile config-example](metalcloud-cli_device-config-template_profile_config-example.md)	 - Generate an example device configuration template profile
* [metalcloud-cli device-config-template profile create](metalcloud-cli_device-config-template_profile_create.md)	 - Create a new device configuration template profile
* [metalcloud-cli device-config-template profile delete](metalcloud-cli_device-config-template_profile_delete.md)	 - Delete a device configuration template profile
* [metalcloud-cli device-config-template profile find-applicable](metalcloud-cli_device-config-template_profile_find-applicable.md)	 - Find device configuration template profiles applicable to a device or fabric
* [metalcloud-cli device-config-template profile get](metalcloud-cli_device-config-template_profile_get.md)	 - Get detailed information about a specific profile
* [metalcloud-cli device-config-template profile list](metalcloud-cli_device-config-template_profile_list.md)	 - List device configuration template profiles with optional filtering
* [metalcloud-cli device-config-template profile render](metalcloud-cli_device-config-template_profile_render.md)	 - Render a device configuration template profile for a device
* [metalcloud-cli device-config-template profile render-applicable](metalcloud-cli_device-config-template_profile_render-applicable.md)	 - Render device configuration template profiles applicable to a device or fabric
* [metalcloud-cli device-config-template profile update](metalcloud-cli_device-config-template_profile_update.md)	 - Update an existing device configuration template profile

