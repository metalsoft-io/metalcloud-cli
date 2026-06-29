## metalcloud-cli point-to-point-link list

List point-to-point links

### Synopsis

List point-to-point links, optionally filtered by a referenced interface id
or route domain id.

Examples:
  metalcloud-cli point-to-point-link list
  metalcloud-cli p2p ls --interface-id 1001
  metalcloud-cli p2p ls --route-domain-id 5

```
metalcloud-cli point-to-point-link list [flags]
```

### Options

```
  -h, --help                     help for list
      --interface-id string      Filter links by a referenced interface id.
      --route-domain-id string   Filter links by route domain id.
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

