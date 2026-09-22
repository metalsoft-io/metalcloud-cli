## metalcloud-cli container-type update

Update a container type

### Synopsis

Update an existing container type.

Only the display name, label, tags and the experimental / unmanaged-only flags
can be changed; CPU and RAM sizing is immutable.

Required Arguments:
  container_type_id  The numeric ID of the container type to update

Required Flags:
  --config-source string  Source of the container type updates.
                          Can be 'pipe' or path to a JSON/YAML file.

Examples:
  # Update container type 42 from a file
  metalcloud-cli container-type update 42 --config-source updates.json

  # Update container type 42 from stdin
  echo '{"label":"new-label"}' | metalcloud-cli ct edit 42 --config-source pipe

```
metalcloud-cli container-type update container_type_id [flags]
```

### Options

```
      --config-source string   Source of the container type updates. Can be 'pipe' or path to a JSON/YAML file.
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

* [metalcloud-cli container-type](metalcloud-cli_container-type.md)	 - Manage container types

