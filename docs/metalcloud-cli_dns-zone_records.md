## metalcloud-cli dns-zone records

List DNS record sets

### Synopsis

List DNS record sets.

Called with a zone ID the command lists the record sets of that zone. Called
without an argument it lists every record set of every zone through the global
record set endpoint.

Optional Arguments:
  zone_id    The ID of the DNS zone. When omitted, the record sets of all zones
             are listed.

Examples:
  # List the record sets of zone 123
  metalcloud-cli dns-zone records 123

  # List the record sets of all zones
  metalcloud-cli dns-zone records


```
metalcloud-cli dns-zone records [zone_id] [flags]
```

### Options

```
  -h, --help   help for records
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

