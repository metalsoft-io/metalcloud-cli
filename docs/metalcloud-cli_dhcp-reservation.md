## metalcloud-cli dhcp-reservation

Manage site DHCP reservations

### Synopsis

Manage the DHCP reservations of a site.

A DHCP reservation pins how the site's DHCP service answers a matching request:
either with a fixed IP (a manual allocation) or with the next free IP of one or
more IPAM OOB subnet pools (an auto allocation). Requests are matched on MAC
address, Option 82 circuit id, relay gateway address, device type and vendor.

Reservations are scoped to a site and an IP version, so every command takes
both: the site by ID or label, and 'ipv4' or 'ipv6'.

Available Commands:
  list            List the reservations of a site and IP version
  get             Get one reservation
  create          Create a reservation
  update          Update a reservation
  delete          Delete a reservation
  config-example  Print example create configurations

Examples:
  metalcloud-cli dhcp-reservation list dc-1 ipv4
  metalcloud-cli dhcp-reservation get dc-1 ipv4 12

### Options

```
  -h, --help   help for dhcp-reservation
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
* [metalcloud-cli dhcp-reservation config-example](metalcloud-cli_dhcp-reservation_config-example.md)	 - Print example DHCP reservation configurations
* [metalcloud-cli dhcp-reservation create](metalcloud-cli_dhcp-reservation_create.md)	 - Create a DHCP reservation
* [metalcloud-cli dhcp-reservation delete](metalcloud-cli_dhcp-reservation_delete.md)	 - Delete a DHCP reservation
* [metalcloud-cli dhcp-reservation get](metalcloud-cli_dhcp-reservation_get.md)	 - Get one DHCP reservation
* [metalcloud-cli dhcp-reservation list](metalcloud-cli_dhcp-reservation_list.md)	 - List the DHCP reservations of a site
* [metalcloud-cli dhcp-reservation update](metalcloud-cli_dhcp-reservation_update.md)	 - Update a DHCP reservation

