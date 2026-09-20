## metalcloud-cli fabric get-link

Get one fabric link

### Synopsis

Get the details of a single link of a fabric.

Use 'fabric get-links' to list the links and their IDs.

Required Arguments:
  fabric_id    The ID or name of the fabric
  link_id      The ID of the link

Examples:
  metalcloud-cli fabric get-link 12345 67890
  metalcloud-cli fabric show-link my-fabric 67890

```
metalcloud-cli fabric get-link fabric_id link_id [flags]
```

### Options

```
  -h, --help   help for get-link
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

* [metalcloud-cli fabric](metalcloud-cli_fabric.md)	 - Manage network fabrics

