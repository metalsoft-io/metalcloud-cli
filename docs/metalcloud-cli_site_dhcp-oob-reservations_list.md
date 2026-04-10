## metalcloud-cli site dhcp-oob-reservations list

List all DHCP OOB reservations for a site

### Synopsis

List all DHCP Option82 to IP address mappings configured for a site.

Displays a table of MAC address to IP address reservations from the site's
serverPolicy.dhcpOption82ToIPMapping configuration.

Examples:
  # List reservations sorted by MAC address (default)
  metalcloud-cli site dhcp-oob-reservations list site-01

  # List reservations sorted by IP address
  metalcloud-cli site dhcp-oob-reservations list site-01 --sort-by ip

  # List in JSON format
  metalcloud-cli site dhcp-oob-reservations list site-01 --format json

```
metalcloud-cli site dhcp-oob-reservations list site_id_or_name [flags]
```

### Options

```
  -h, --help             help for list
      --sort-by string   Sort results by field: mac or ip (default "mac")
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

* [metalcloud-cli site dhcp-oob-reservations](metalcloud-cli_site_dhcp-oob-reservations.md)	 - Manage DHCP Option82 OOB IP reservations for a site

