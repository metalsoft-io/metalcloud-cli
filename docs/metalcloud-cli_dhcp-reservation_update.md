## metalcloud-cli dhcp-reservation update

Update a DHCP reservation

### Synopsis

Update an existing DHCP reservation.

The update replaces the reservation (the endpoint is a PUT), so the
configuration must describe the reservation in full, allocation included. The
reservation's current revision is sent as the If-Match entity tag, so a
concurrent change is rejected instead of being overwritten.

Required Arguments:
  site_id_or_label  The ID or label of the site
  ip_version        'ipv4' or 'ipv6'
  reservation_id    The ID of the reservation

Required Flags:
  --config-source  'pipe' to read from stdin, or a path to a JSON/YAML file

Examples:
  metalcloud-cli dhcp-reservation update 1 ipv4 12 --config-source reservation.yaml
  metalcloud-cli dhcp-reservation get 1 ipv4 12 -f json > r.json
  metalcloud-cli dhcp-reservation update 1 ipv4 12 --config-source r.json

```
metalcloud-cli dhcp-reservation update site_id_or_label ip_version reservation_id [flags]
```

### Options

```
      --config-source string   Source of the DHCP reservation configuration. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for update
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

