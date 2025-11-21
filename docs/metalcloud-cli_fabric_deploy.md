## metalcloud-cli fabric deploy

Deploy a fabric

### Synopsis

Deploy a network fabric underlay.

This command deploys fabric underlay using the configured links and templates.

Arguments:
  fabric_id    The ID or label of the fabric to deploy

Examples:
  # Deploy fabric by ID
  metalcloud fabric deploy 12345
  
  # Deploy fabric by label
  metalcloud fabric deploy my-fabric-label

```
metalcloud-cli fabric deploy fabric_id [flags]
```

### Options

```
  -h, --help   help for deploy
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

