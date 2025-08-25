## metalcloud-cli fabric create

Create a new fabric

### Synopsis

Create a new network fabric in MetalCloud.

This command creates a new fabric with the specified configuration. The fabric configuration
must be provided through the --config-source flag, which can be a JSON file or piped input.

Arguments:
  site_id_or_label       The ID or label of the site where the fabric will be created
  fabric_name           The name for the new fabric
  fabric_type           The type of fabric to create (e.g., "spine-leaf", "collapsed-core")
  fabric_description    Optional description for the fabric (defaults to fabric_name if not provided)

Required Flags:
  --config-source string   Source of the fabric configuration. Can be 'pipe' for piped input
                          or path to a JSON file containing the fabric configuration

Examples:
  # Create fabric with configuration from file
  metalcloud fabric create site1 my-fabric spine-leaf "Production fabric" --config-source fabric-config.json
  
  # Create fabric with piped configuration
  cat fabric-config.json | metalcloud fabric create site1 my-fabric spine-leaf --config-source pipe
  
  # Get example config and create fabric
  metalcloud fabric config-example spine-leaf > config.json
  # Edit config.json as needed
  metalcloud fabric create site1 my-fabric spine-leaf --config-source config.json

```
metalcloud-cli fabric create site_id_or_label fabric_name fabric_type [fabric_description] [flags]
```

### Options

```
      --config-source string   Source of the new fabric configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for create
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

