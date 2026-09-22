## metalcloud-cli fabric accept-deploy

Accept a pending fabric deployment

### Synopsis

Accept a fabric deployment that is waiting for confirmation.

A deploy started with confirmation required stays pending until it is accepted
or rejected. Accepting it lets the deploy job continue.

Required Arguments:
  fabric_id    The ID or name of the fabric

Examples:
  metalcloud-cli fabric accept-deploy 12345
  metalcloud-cli fabric accept-deploy my-fabric

```
metalcloud-cli fabric accept-deploy fabric_id [flags]
```

### Options

```
  -h, --help   help for accept-deploy
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

