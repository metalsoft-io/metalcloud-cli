## metalcloud-cli dns-zone record

Get a single DNS record set

### Synopsis

Get the details of a single DNS record set of a zone.

Required Arguments:
  zone_id          The ID of the DNS zone
  record_set_id    The ID of the DNS record set

Examples:
  # Get record set 456 of zone 123
  metalcloud-cli dns-zone record 123 456


```
metalcloud-cli dns-zone record zone_id record_set_id [flags]
```

### Options

```
  -h, --help   help for record
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

