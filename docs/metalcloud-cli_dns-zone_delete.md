## metalcloud-cli dns-zone delete

Delete a DNS zone

### Synopsis

Delete a DNS zone from MetalSoft infrastructure.

This command permanently deletes a DNS zone and all its associated DNS records.
This action cannot be undone, so use with caution.

Required Arguments:
  dns_zone_id           The ID of the DNS zone to delete

Examples:
  # Delete a DNS zone
  metalcloud-cli dns-zone delete 123

  # Delete a DNS zone using alias
  metalcloud-cli dns-zone rm 123


```
metalcloud-cli dns-zone delete dns_zone_id [flags]
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

* [metalcloud-cli dns-zone](metalcloud-cli_dns-zone.md)	 - DNS Zone management

