## metalcloud-cli network-fabric-interconnect deploy

Deploy a network fabric interconnect

### Synopsis

Deploy a network fabric interconnect, pushing the BGP configuration to the
network devices referenced by its links. Returns the job that performs the deploy.

Required Arguments:
  interconnect_id_or_label   The ID or label of the interconnect

Optional Flags:
  --require-confirmation     Pause the deploy until it is accepted with 'accept-deploy'
                             (or discarded with 'reject-deploy')

Examples:
  metalcloud network-fabric-interconnect deploy 12
  metalcloud nfi deploy dc1-dc2 --require-confirmation

```
metalcloud-cli network-fabric-interconnect deploy interconnect_id_or_label [flags]
```

### Options

```
  -h, --help                   help for deploy
      --require-confirmation   Pause the deploy until it is accepted with 'accept-deploy'.
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

