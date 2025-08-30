## metalcloud-cli dns-zone

DNS Zone management

### Synopsis

DNS Zone management commands.

This command group provides comprehensive DNS zone management capabilities including
creation, retrieval, updating, and deletion of DNS zones. DNS zones can be
managed individually with their associated record sets.

Available command categories:
  - Basic operations: list, get, create, update, delete
  - Record management: list-records, get-record
  - Information: nameservers

Use "metalcloud-cli dns-zone [command] --help" for detailed information about each command.


### Options

```
  -h, --help   help for dns-zone
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
* [metalcloud-cli dns-zone create](metalcloud-cli_dns-zone_create.md)	 - Create a new DNS zone
* [metalcloud-cli dns-zone delete](metalcloud-cli_dns-zone_delete.md)	 - Delete a DNS zone
* [metalcloud-cli dns-zone get](metalcloud-cli_dns-zone_get.md)	 - Get detailed DNS zone information
* [metalcloud-cli dns-zone list](metalcloud-cli_dns-zone_list.md)	 - List DNS zones
* [metalcloud-cli dns-zone update](metalcloud-cli_dns-zone_update.md)	 - Update DNS zone information

