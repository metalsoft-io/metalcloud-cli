## metalcloud-cli external-system

External system management

### Synopsis

Manage the external systems registered with the platform.

An external system is a third-party system the platform integrates with; it is
identified by a unique label and can carry free-form annotations.

Command categories:
  Read:    list, get
  Write:   create, update, delete
  Helper:  config-example

### Options

```
  -h, --help   help for external-system
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
* [metalcloud-cli external-system config-example](metalcloud-cli_external-system_config-example.md)	 - Print an example external system configuration
* [metalcloud-cli external-system create](metalcloud-cli_external-system_create.md)	 - Create an external system
* [metalcloud-cli external-system delete](metalcloud-cli_external-system_delete.md)	 - Delete an external system
* [metalcloud-cli external-system get](metalcloud-cli_external-system_get.md)	 - Get external system details
* [metalcloud-cli external-system list](metalcloud-cli_external-system_list.md)	 - List external systems
* [metalcloud-cli external-system update](metalcloud-cli_external-system_update.md)	 - Update an external system

