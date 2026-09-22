## metalcloud-cli container power-status

Get the power status of a container

### Synopsis

Get the current power status of a container.

Required Arguments:
  container_id  The numeric ID of the container

Examples:
  # Show the power status of container 100
  metalcloud-cli container power-status 100

  # Using the alias
  metalcloud-cli containers power-state 100

```
metalcloud-cli container power-status container_id [flags]
```

### Options

```
  -h, --help   help for power-status
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

* [metalcloud-cli container](metalcloud-cli_container.md)	 - Manage provisioned containers

