## metalcloud-cli dns-zone nameservers

List the nameservers of a DNS zone

### Synopsis

List the nameservers configured for a DNS zone.

Required Arguments:
  dns_zone_id    The ID of the DNS zone

Examples:
  # List the nameservers of zone 123
  metalcloud-cli dns-zone nameservers 123

  # List them as JSON
  metalcloud-cli dns-zone nameservers 123 -f json


```
metalcloud-cli dns-zone nameservers dns_zone_id [flags]
```

### Options

```
  -h, --help   help for nameservers
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

* [metalcloud-cli dns-zone](metalcloud-cli_dns-zone.md)	 - DNS Zone management

