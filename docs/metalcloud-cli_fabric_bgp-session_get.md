## metalcloud-cli fabric bgp-session get

Get one BGP session of a fabric

### Synopsis

Get the details of one BGP session of a network fabric.

Required Arguments:
  fabric_id         The ID or name of the fabric
  bgp_session_id    The ID of the BGP session

Examples:
  metalcloud-cli fabric bgp-session get 12345 7
  metalcloud-cli fabric bgp show my-fabric 7

```
metalcloud-cli fabric bgp-session get fabric_id bgp_session_id [flags]
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

* [metalcloud-cli fabric bgp-session](metalcloud-cli_fabric_bgp-session.md)	 - Manage the BGP sessions of a fabric

