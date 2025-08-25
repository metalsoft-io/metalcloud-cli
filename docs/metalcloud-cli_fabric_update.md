## metalcloud-cli fabric update

Update fabric configuration

### Synopsis

Update the configuration, name, or description of an existing fabric.

This command allows you to modify fabric properties and configuration. The fabric
configuration can be updated by providing a new configuration through the --config-source flag.
The name and description are optional and will only be updated if provided.

Arguments:
  fabric_id            The ID or label of the fabric to update
  fabric_name          Optional new name for the fabric
  fabric_description   Optional new description for the fabric

Required Flags:
  --config-source string   Source of the updated fabric configuration. Can be 'pipe' for piped input
                          or path to a JSON file containing the updated configuration

Examples:
  # Update fabric configuration from file
  metalcloud fabric update 12345 --config-source updated-config.json
  
  # Update name, description and configuration
  metalcloud fabric update my-fabric "New Name" "New Description" --config-source config.json
  
  # Update with piped configuration
  cat new-config.json | metalcloud fabric update 12345 --config-source pipe
  
  # Update only configuration, keeping existing name and description
  metalcloud fabric update my-fabric --config-source config.json

```
metalcloud-cli fabric update fabric_id [fabric_name [fabric_description]] [flags]
```

### Options

```
      --config-source string   Source of the updated fabric configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for update
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

