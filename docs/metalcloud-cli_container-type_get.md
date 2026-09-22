## metalcloud-cli container-type get

Get container type details

### Synopsis

Get detailed information about a container type.

Required Arguments:
  container_type_id  The numeric ID of the container type

Examples:
  # Get the details of container type 42
  metalcloud-cli container-type get 42

  # Using the alias
  metalcloud-cli ct show 42

```
metalcloud-cli container-type get container_type_id [flags]
```

### Options

```
  -h, --help   help for get
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

