## metalcloud-cli point-to-point-link static-route add

Stage a static route on a point-to-point link

### Synopsis

Stage a static route on one address family of a point-to-point link's
configuration. The destination prefix's address family must match the
collection it is added to.

Required Arguments:
  link_id               The ID of the point-to-point link
  family                One of: ipv4, ipv6
  destination_prefix    Destination prefix in CIDR notation

Examples:
  metalcloud-cli point-to-point-link static-route add 42 ipv4 10.0.0.0/24
  metalcloud-cli p2p route add 42 ipv6 2001:db8::/64

```
metalcloud-cli point-to-point-link static-route add link_id family destination_prefix [flags]
```

### Options

```
  -h, --help   help for add
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

* [metalcloud-cli point-to-point-link static-route](metalcloud-cli_point-to-point-link_static-route.md)	 - Manage the staged static routes of a point-to-point link

