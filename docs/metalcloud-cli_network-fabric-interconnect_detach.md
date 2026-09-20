## metalcloud-cli network-fabric-interconnect detach

Detach an interconnect, removing its configuration from the devices

### Synopsis

Detach a network fabric interconnect. The BGP configuration is removed from all
linked network devices. Returns the job that performs the detach.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect

Examples:
  metalcloud network-fabric-interconnect detach 12

```
metalcloud-cli network-fabric-interconnect detach interconnect_id_or_label [flags]
```

### Options

```
  -h, --help   help for detach
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

