## metalcloud-cli dns-zone get

Get detailed DNS zone information

### Synopsis

Get detailed information for a specific DNS zone.

This command retrieves comprehensive information about a DNS zone including its
configuration, status, name servers, and other metadata.

Required Arguments:
  dns_zone_id           The ID of the DNS zone to retrieve information for

Examples:
  # Get DNS zone information
  metalcloud-cli dns-zone get 123

  # Get DNS zone information using alias
  metalcloud-cli dns-zone show 123


```
metalcloud-cli dns-zone get dns_zone_id [flags]
```

### Options

```
  -h, --help   help for get
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

