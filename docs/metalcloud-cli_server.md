## metalcloud-cli server

Server management

### Synopsis

Server management commands.

This command group provides comprehensive server management capabilities including
registration, power control, firmware management, and monitoring. Servers can be
managed individually or in bulk operations.

Available command categories:
  - Basic operations: list, get, register, update, delete
  - Power management: power (on, off, reset, cycle, soft, status)
  - Maintenance: re-register, factory-reset, archive
  - Security: update-ipmi-credentials, enable-snmp, enable-syslog
  - Remote access: vnc-info, console-info
  - Firmware: firmware subcommands for component management and upgrades
  - Information: capabilities, statistics
  - Drift detection: drift subcommands, snapshots, sync-target-snapshot
  - Hardware: hardware-rescan, connect-interface, set-interfaces-default-fabric,
    set-interfaces-redundancy-group
  - Onboarding: register-production, import-unmanaged, config-example

Use "metalcloud-cli server [command] --help" for detailed information about each command.


### Options

```
  -h, --help   help for server
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
* [metalcloud-cli server archive](metalcloud-cli_server_archive.md)	 - Archive a server
* [metalcloud-cli server capabilities](metalcloud-cli_server_capabilities.md)	 - Get server capabilities
* [metalcloud-cli server config-example](metalcloud-cli_server_config-example.md)	 - Show a server configuration example
* [metalcloud-cli server connect-interface](metalcloud-cli_server_connect-interface.md)	 - Record the network device port a server interface is cabled to
* [metalcloud-cli server console-info](metalcloud-cli_server_console-info.md)	 - Get server remote console information
* [metalcloud-cli server delete](metalcloud-cli_server_delete.md)	 - Delete a server
* [metalcloud-cli server drift](metalcloud-cli_server_drift.md)	 - Inspect the configuration drift of a server
* [metalcloud-cli server enable-snmp](metalcloud-cli_server_enable-snmp.md)	 - Enable SNMP on server
* [metalcloud-cli server enable-syslog](metalcloud-cli_server_enable-syslog.md)	 - Enable remote syslog for a server
* [metalcloud-cli server factory-reset](metalcloud-cli_server_factory-reset.md)	 - Reset a server to factory defaults
* [metalcloud-cli server firmware](metalcloud-cli_server_firmware.md)	 - Server firmware management
* [metalcloud-cli server get](metalcloud-cli_server_get.md)	 - Get detailed server information
* [metalcloud-cli server hardware-rescan](metalcloud-cli_server_hardware-rescan.md)	 - Re-read the hardware inventory of a server
* [metalcloud-cli server identify](metalcloud-cli_server_identify.md)	 - Identify a server by blinking its chassis LED
* [metalcloud-cli server import-unmanaged](metalcloud-cli_server_import-unmanaged.md)	 - Import a server whose lifecycle MetalSoft does not manage
* [metalcloud-cli server list](metalcloud-cli_server_list.md)	 - List servers
* [metalcloud-cli server power](metalcloud-cli_server_power.md)	 - Control server power state
* [metalcloud-cli server re-register](metalcloud-cli_server_re-register.md)	 - Re-register an existing server
* [metalcloud-cli server register](metalcloud-cli_server_register.md)	 - Register a new server in MetalSoft
* [metalcloud-cli server register-production](metalcloud-cli_server_register-production.md)	 - Register a server that is already running a production workload
* [metalcloud-cli server set-interfaces-default-fabric](metalcloud-cli_server_set-interfaces-default-fabric.md)	 - Set the default fabric of some server interfaces
* [metalcloud-cli server set-interfaces-redundancy-group](metalcloud-cli_server_set-interfaces-redundancy-group.md)	 - Set the redundancy group of some server interfaces
* [metalcloud-cli server snapshots](metalcloud-cli_server_snapshots.md)	 - List the configuration snapshots of a server
* [metalcloud-cli server statistics](metalcloud-cli_server_statistics.md)	 - Get aggregated server statistics
* [metalcloud-cli server sync-target-snapshot](metalcloud-cli_server_sync-target-snapshot.md)	 - Accept the current configuration of a server as the drift target
* [metalcloud-cli server update](metalcloud-cli_server_update.md)	 - Update server information
* [metalcloud-cli server update-ipmi-credentials](metalcloud-cli_server_update-ipmi-credentials.md)	 - Update server IPMI credentials
* [metalcloud-cli server vnc-info](metalcloud-cli_server_vnc-info.md)	 - Get server VNC information

