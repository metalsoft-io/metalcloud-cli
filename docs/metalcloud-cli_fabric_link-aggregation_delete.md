## metalcloud-cli fabric link-aggregation delete

Delete a link aggregation of a fabric

### Synopsis

Remove a link aggregation from a network fabric.

Required Arguments:
  fabric_id              The ID or name of the fabric
  link_aggregation_id    The ID of the link aggregation

Examples:
  metalcloud-cli fabric link-aggregation delete 12345 3
  metalcloud-cli fabric lag rm my-fabric 3

```
metalcloud-cli fabric link-aggregation delete fabric_id link_aggregation_id [flags]
```

### Options

```
  -h, --help   help for delete
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

* [metalcloud-cli fabric link-aggregation](metalcloud-cli_fabric_link-aggregation.md)	 - Manage the link aggregations of a fabric

