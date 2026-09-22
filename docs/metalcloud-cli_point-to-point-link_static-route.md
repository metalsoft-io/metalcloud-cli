## metalcloud-cli point-to-point-link static-route

Manage the staged static routes of a point-to-point link

### Synopsis

Manage the static routes staged on a point-to-point link's configuration.

Static routes live under the link's config, one collection per address family
(ipv4 / ipv6), and are guarded by the config object's revision.

Commands:
  list, get, add, remove

### Options

```
  -h, --help   help for static-route
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

* [metalcloud-cli point-to-point-link](metalcloud-cli_point-to-point-link.md)	 - Manage point-to-point links between network interfaces
* [metalcloud-cli point-to-point-link static-route add](metalcloud-cli_point-to-point-link_static-route_add.md)	 - Stage a static route on a point-to-point link
* [metalcloud-cli point-to-point-link static-route get](metalcloud-cli_point-to-point-link_static-route_get.md)	 - Get one staged static route
* [metalcloud-cli point-to-point-link static-route list](metalcloud-cli_point-to-point-link_static-route_list.md)	 - List the staged static routes of one address family
* [metalcloud-cli point-to-point-link static-route remove](metalcloud-cli_point-to-point-link_static-route_remove.md)	 - Remove a staged static route

