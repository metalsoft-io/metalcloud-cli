## metalcloud-cli point-to-point-link static-route list

List the staged static routes of one address family

### Synopsis

List the static routes staged on a point-to-point link for one address family.

Required Arguments:
  link_id    The ID of the point-to-point link
  family     One of: ipv4, ipv6

Examples:
  metalcloud-cli point-to-point-link static-route list 42 ipv4
  metalcloud-cli p2p route ls 42 ipv6

```
metalcloud-cli point-to-point-link static-route list link_id family [flags]
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

* [metalcloud-cli point-to-point-link static-route](metalcloud-cli_point-to-point-link_static-route.md)	 - Manage the staged static routes of a point-to-point link

