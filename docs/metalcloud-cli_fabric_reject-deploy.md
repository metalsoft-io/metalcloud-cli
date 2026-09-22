## metalcloud-cli fabric reject-deploy

Reject a pending fabric deployment

### Synopsis

Reject a fabric deployment that is waiting for confirmation.

Required Arguments:
  fabric_id    The ID or name of the fabric

Examples:
  metalcloud-cli fabric reject-deploy 12345
  metalcloud-cli fabric reject-deploy my-fabric

```
metalcloud-cli fabric reject-deploy fabric_id [flags]
```

### Options

```
  -h, --help   help for reject-deploy
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

