## metalcloud-cli server firmware

Server firmware management

### Synopsis

Server firmware management commands.

This command group provides comprehensive firmware management capabilities for servers
including component listing, upgrades, scheduling, and auditing. Firmware operations
can be performed on individual components or entire servers.

Available commands:
  - Information: components, component-info, inventory, fetch-versions
  - Updates: update-info, update-component
  - Upgrades: upgrade, upgrade-component, schedule-upgrade
  - Auditing: generate-audit

Use "metalcloud-cli server firmware [command] --help" for detailed information about each command.


### Options

```
  -h, --help   help for firmware
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

* [metalcloud-cli server](metalcloud-cli_server.md)	 - Server management
* [metalcloud-cli server firmware component-info](metalcloud-cli_server_firmware_component-info.md)	 - Get firmware component information
* [metalcloud-cli server firmware components](metalcloud-cli_server_firmware_components.md)	 - List firmware components for a server
* [metalcloud-cli server firmware fetch-versions](metalcloud-cli_server_firmware_fetch-versions.md)	 - Fetch available firmware versions for a server
* [metalcloud-cli server firmware generate-audit](metalcloud-cli_server_firmware_generate-audit.md)	 - Generate firmware upgrade audit for servers
* [metalcloud-cli server firmware inventory](metalcloud-cli_server_firmware_inventory.md)	 - Get firmware inventory from redfish
* [metalcloud-cli server firmware schedule-upgrade](metalcloud-cli_server_firmware_schedule-upgrade.md)	 - Schedule a firmware upgrade for a server
* [metalcloud-cli server firmware update-component](metalcloud-cli_server_firmware_update-component.md)	 - Update firmware component settings
* [metalcloud-cli server firmware update-info](metalcloud-cli_server_firmware_update-info.md)	 - Update firmware information for a server
* [metalcloud-cli server firmware upgrade](metalcloud-cli_server_firmware_upgrade.md)	 - Upgrade firmware for all components on a server
* [metalcloud-cli server firmware upgrade-component](metalcloud-cli_server_firmware_upgrade-component.md)	 - Upgrade firmware for a specific component

