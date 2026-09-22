## metalcloud-cli resource-pool update

Update an existing resource pool

### Synopsis

Update the label and/or the description of an existing resource pool.

The new values are read from a JSON or YAML configuration. Resource pools are not
under optimistic concurrency control, so no revision has to be supplied.

Required Arguments:
  pool_id           The numeric ID of the resource pool to update

Required Flags:
  --config-source   Source of the resource pool update configuration. Can be 'pipe' or path to a JSON/YAML file.

The configuration accepts the following properties:
  resourcePoolLabel        (optional) the new label
  resourcePoolDescription  (optional) the new description

Examples:
  # Update a resource pool from a file
  metalcloud-cli resource-pool update 123 --config-source ./pool.json

  # Rename a resource pool from piped input
  echo '{"resourcePoolLabel":"Production Pool"}' | metalcloud-cli resource-pool update 123 --config-source pipe

```
metalcloud-cli resource-pool update pool_id [flags]
```

### Options

```
      --config-source string   Source of the resource pool update configuration. Can be 'pipe' or path to a JSON file.
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

* [metalcloud-cli resource-pool](metalcloud-cli_resource-pool.md)	 - Manage resource pools and their associated resources

