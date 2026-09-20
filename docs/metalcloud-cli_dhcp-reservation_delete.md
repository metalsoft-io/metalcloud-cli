## metalcloud-cli dhcp-reservation delete

Delete a DHCP reservation

### Synopsis

Delete a DHCP reservation from a site.

The reservation's current revision is sent as the If-Match entity tag, so a
concurrent change is rejected instead of being overwritten.

Required Arguments:
  site_id_or_label  The ID or label of the site
  ip_version        'ipv4' or 'ipv6'
  reservation_id    The ID of the reservation

Examples:
  metalcloud-cli dhcp-reservation delete 1 ipv4 12
  metalcloud-cli dhcp-reservation rm dc-1 ipv6 13

```
metalcloud-cli dhcp-reservation delete site_id_or_label ip_version reservation_id [flags]
```

### Options

```
  -h, --help   help for delete
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

