## metalcloud-cli site dhcp-oob-reservations

Manage DHCP Option82 OOB IP reservations for a site

### Synopsis

Manage DHCP Option82 to IP address mappings in the site's server policy configuration.

These reservations map MAC addresses (DHCP Option82) to static IP addresses for
out-of-band (OOB) management interfaces. The mappings are stored in the site
configuration under serverPolicy.dhcpOption82ToIPMapping.

Available Commands:
  list      List all DHCP OOB reservations for a site
  add       Add a MAC-to-IP reservation entry
  remove    Remove a reservation entry by MAC address
  replace   Replace all reservation entries from JSON input

### Options

```
  -h, --help   help for dhcp-oob-reservations
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

* [metalcloud-cli site](metalcloud-cli_site.md)	 - Manage sites (datacenters) and their configurations
* [metalcloud-cli site dhcp-oob-reservations add](metalcloud-cli_site_dhcp-oob-reservations_add.md)	 - Add a DHCP OOB reservation entry
* [metalcloud-cli site dhcp-oob-reservations list](metalcloud-cli_site_dhcp-oob-reservations_list.md)	 - List all DHCP OOB reservations for a site
* [metalcloud-cli site dhcp-oob-reservations remove](metalcloud-cli_site_dhcp-oob-reservations_remove.md)	 - Remove a DHCP OOB reservation entry by MAC address
* [metalcloud-cli site dhcp-oob-reservations replace](metalcloud-cli_site_dhcp-oob-reservations_replace.md)	 - Replace all DHCP OOB reservation entries

