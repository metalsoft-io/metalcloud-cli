## metalcloud-cli network-fabric-interconnect add-link

Add a link to an interconnect

### Synopsis

Add a link to a network fabric interconnect, binding a fabric and one of its
network devices (the border device that will run the interconnect BGP session).

Required Arguments:
  interconnect_id_or_label      The ID or label of the interconnect
  fabric_id_or_label            The ID or name of the fabric to link
  network_device_id_or_label    The ID or identifier of the network device in that fabric

Examples:
  metalcloud network-fabric-interconnect add-link 12 dc1-fabric 45
  metalcloud nfi add-link dc1-dc2 7 border-leaf-01

```
metalcloud-cli network-fabric-interconnect add-link interconnect_id_or_label fabric_id_or_label network_device_id_or_label [flags]
```

### Options

```
  -h, --help   help for add-link
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

