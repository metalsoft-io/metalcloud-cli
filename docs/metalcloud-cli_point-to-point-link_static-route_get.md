## metalcloud-cli point-to-point-link static-route get

Get one staged static route

### Synopsis

Get one static route staged on a point-to-point link.

Required Arguments:
  link_id    The ID of the point-to-point link
  family     One of: ipv4, ipv6
  route_id   The ID of the static route

Examples:
  metalcloud-cli point-to-point-link static-route get 42 ipv4 3

```
metalcloud-cli point-to-point-link static-route get link_id family route_id [flags]
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

* [metalcloud-cli point-to-point-link static-route](metalcloud-cli_point-to-point-link_static-route.md)	 - Manage the staged static routes of a point-to-point link

