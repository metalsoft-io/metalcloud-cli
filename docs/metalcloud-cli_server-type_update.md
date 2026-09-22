## metalcloud-cli server-type update

Update an existing server type

### Synopsis

Update an existing server type from a JSON or YAML configuration.

Only the descriptive attributes of a server type can be updated: its label,
name, description, experimental flag, tags and allowed vendor SKU ids. The
hardware specification itself is derived from the registered servers.

Required Arguments:
  server_type_id     The numeric ID of the server type to update

Required Flags:
  --config-source    Source of the server type update configuration. Can be 'pipe' or path to a JSON file.

Examples:
  # Update a server type from a JSON file
  metalcloud-cli server-type update 123 --config-source ./server-type-update.json

  # Update a server type from piped configuration
  echo '{"label":"m-32-128-2","description":"Updated"}' | metalcloud-cli server-type update 123 --config-source pipe


```
metalcloud-cli server-type update server_type_id [flags]
```

### Options

```
      --config-source string   Source of the server type update configuration. Can be 'pipe' or path to a JSON file.
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

* [metalcloud-cli server-type](metalcloud-cli_server-type.md)	 - Manage server types and hardware configurations

