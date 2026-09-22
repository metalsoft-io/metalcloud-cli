## metalcloud-cli storage update-interface

Update one interface of a storage pool

### Synopsis

Update a storage pool interface from a JSON or YAML configuration.

The updatable properties are isUplink, useForDeploys and
networkEquipmentInterfaceId. The current revision of the interface is read
first and sent back as the If-Match entity tag.

Required Arguments:
  storage_id      The numeric ID of the storage pool
  interface_id    The numeric ID of the storage interface

Required Flags:
  --config-source   Source of the interface update configuration. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  # Mark an interface as the one used for deploys
  echo '{"useForDeploys": true}' | metalcloud-cli storage update-interface 123 7 --config-source pipe

  # Update an interface from a file
  metalcloud-cli storage update-interface 123 7 --config-source ./interface.json

```
metalcloud-cli storage update-interface storage_id interface_id [flags]
```

### Options

```
      --config-source string   Source of the storage interface update configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for update-interface
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

