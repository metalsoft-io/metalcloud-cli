## metalcloud-cli fabric activate

Activate a fabric

### Synopsis

Activate a network fabric to make it operational.

This command activates a fabric that has been created and configured. Once activated,
the fabric will begin managing the network connectivity according to its configuration.
Only fabrics in an inactive state can be activated.

Arguments:
  fabric_id    The ID or label of the fabric to activate

Examples:
  # Activate fabric by ID
  metalcloud fabric activate 12345
  
  # Activate fabric by label
  metalcloud fabric activate my-fabric-label
  
  # Using alias
  metalcloud fc start 12345

```
metalcloud-cli fabric activate fabric_id [flags]
```

### Options

```
  -h, --help   help for activate
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

