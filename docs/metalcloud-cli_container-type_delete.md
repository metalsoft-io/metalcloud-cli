## metalcloud-cli container-type delete

Delete a container type

### Synopsis

Delete a container type.

The container type must no longer be referenced by any container instance.

Required Arguments:
  container_type_id  The numeric ID of the container type to delete

Examples:
  # Delete container type 42
  metalcloud-cli container-type delete 42

  # Using the alias
  metalcloud-cli ct rm 42

```
metalcloud-cli container-type delete container_type_id [flags]
```

### Options

```
  -h, --help   help for delete
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

