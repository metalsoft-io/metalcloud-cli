## metalcloud-cli quota-profile

Quota Profile management

### Synopsis

Quota Profile management commands.

This command group provides comprehensive quota profile management capabilities
including creation, retrieval, updating, and deletion of quota profiles. Quota
profiles define resource limits that can be applied to users.

Available commands:
  - Basic operations: list, get, create, update, delete
  - Configuration: config-example

Use "metalcloud-cli quota-profile [command] --help" for detailed information about each command.


### Options

```
  -h, --help   help for quota-profile
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
* [metalcloud-cli quota-profile config-example](metalcloud-cli_quota-profile_config-example.md)	 - Show quota profile configuration example
* [metalcloud-cli quota-profile create](metalcloud-cli_quota-profile_create.md)	 - Create a new quota profile
* [metalcloud-cli quota-profile delete](metalcloud-cli_quota-profile_delete.md)	 - Delete a quota profile
* [metalcloud-cli quota-profile get](metalcloud-cli_quota-profile_get.md)	 - Get detailed quota profile information
* [metalcloud-cli quota-profile list](metalcloud-cli_quota-profile_list.md)	 - List quota profiles
* [metalcloud-cli quota-profile update](metalcloud-cli_quota-profile_update.md)	 - Update quota profile information

