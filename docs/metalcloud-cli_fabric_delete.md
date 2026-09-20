## metalcloud-cli fabric delete

Delete a fabric

### Synopsis

Delete a network fabric.

The fabric must no longer have devices attached to it. Deleting a fabric cannot
be undone.

Required Arguments:
  fabric_id    The ID or name of the fabric to delete

Examples:
  metalcloud-cli fabric delete 12345
  metalcloud-cli fabric rm my-fabric

```
metalcloud-cli fabric delete fabric_id [flags]
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

* [metalcloud-cli fabric](metalcloud-cli_fabric.md)	 - Manage network fabrics

