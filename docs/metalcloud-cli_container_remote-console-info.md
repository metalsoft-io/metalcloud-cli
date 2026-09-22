## metalcloud-cli container remote-console-info

Get remote console information for a container

### Synopsis

Get remote console information for a container.

Shows how many remote console connections are currently active for the container.

Required Arguments:
  container_id  The numeric ID of the container

Examples:
  # Show the remote console info of container 100
  metalcloud-cli container remote-console-info 100

  # Using the alias
  metalcloud-cli containers console-info 100

```
metalcloud-cli container remote-console-info container_id [flags]
```

### Options

```
  -h, --help   help for remote-console-info
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

