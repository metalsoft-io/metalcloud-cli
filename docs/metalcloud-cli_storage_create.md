## metalcloud-cli storage create

Create a new storage pool

### Synopsis

Create a new storage pool using a configuration file or piped input.

This command creates a new storage pool in the MetalCloud infrastructure. The storage
configuration must be provided as JSON either through a file or piped input.

Required flags:
  --config-source    Source of the storage configuration. Can be 'pipe' for piped input
                     or a path to a JSON file containing the storage configuration.

The configuration must include required fields such as siteId, driver, technologies,
type, name, managementHost, username, password, and subnetType. Use the 'config-example'
command to see a complete template with all available options.

Examples:
  # Create storage from a JSON file
  metalcloud storage create --config-source ./storage-config.json

  # Create storage from piped input
  cat storage-config.json | metalcloud storage create --config-source pipe

  # Generate template, edit, and create
  metalcloud storage config-example > config.json
  # Edit config.json with your storage details
  metalcloud storage create --config-source config.json

```
metalcloud-cli storage create [flags]
```

### Options

```
      --config-source string   Source of the new storage configuration. Can be 'pipe' or path to a JSON file.
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

* [metalcloud-cli storage](metalcloud-cli_storage.md)	 - Manage storage pools and related resources

