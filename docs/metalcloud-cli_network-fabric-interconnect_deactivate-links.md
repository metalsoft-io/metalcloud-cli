## metalcloud-cli network-fabric-interconnect deactivate-links

Deactivate interconnect links

### Synopsis

Deactivate one or more links of a network fabric interconnect, removing the BGP
configuration from their network devices. Returns the job that performs the deactivation.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect
  link_id...                 One or more link IDs to deactivate

Optional Flags:
  --require-confirmation     Pause until the deploy is accepted with 'accept-deploy'

Examples:
  metalcloud network-fabric-interconnect deactivate-links 12 3

```
metalcloud-cli network-fabric-interconnect deactivate-links interconnect_id_or_label link_id... [flags]
```

### Options

```
  -h, --help                   help for deactivate-links
      --require-confirmation   Pause until the deploy is accepted with 'accept-deploy'.
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

* [metalcloud-cli network-fabric-interconnect](metalcloud-cli_network-fabric-interconnect.md)	 - Network fabric interconnect management

