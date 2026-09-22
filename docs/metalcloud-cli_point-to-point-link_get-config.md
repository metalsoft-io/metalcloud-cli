## metalcloud-cli point-to-point-link get-config

Get the staged configuration of a point-to-point link

### Synopsis

Get the staged configuration object of a point-to-point link: its MTU, the
IPv4/IPv6 subnet allocation strategies and the staged static routes.

The configuration object carries its own revision, which is the one that guards
the config sub-resources (static routes, allocation strategies).

Required Arguments:
  link_id            The ID of the point-to-point link

Examples:
  metalcloud-cli point-to-point-link get-config 42
  metalcloud-cli p2p show-config 42 -f yaml

```
metalcloud-cli point-to-point-link get-config link_id [flags]
```

### Options

```
  -h, --help   help for get-config
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

