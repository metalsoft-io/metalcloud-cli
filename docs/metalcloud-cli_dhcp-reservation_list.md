## metalcloud-cli dhcp-reservation list

List the DHCP reservations of a site

### Synopsis

List all DHCP reservations of one site and IP version.

Required Arguments:
  site_id_or_label  The ID or label of the site
  ip_version        'ipv4' or 'ipv6'

Examples:
  metalcloud-cli dhcp-reservation list 1 ipv4
  metalcloud-cli dhcp-reservation ls dc-1 ipv6

```
metalcloud-cli dhcp-reservation list site_id_or_label ip_version [flags]
```

### Options

```
  -h, --help   help for list
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

* [metalcloud-cli dhcp-reservation](metalcloud-cli_dhcp-reservation.md)	 - Manage site DHCP reservations

