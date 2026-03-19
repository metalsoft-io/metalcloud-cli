## metalcloud-cli server-instance power status

Get power status of a server instance

### Synopsis

Get the current power status of a server instance.

Arguments:
  server_instance_id  The numeric ID of the server instance

Examples:
  # Get power status of server instance 5678
  metalcloud-cli server-instance power status 5678

  # Using alias
  metalcloud-cli inst power power-state 5678

```
metalcloud-cli server-instance power status <server_instance_id> [flags]
```

### Options

```
  -h, --help   help for status
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

* [metalcloud-cli server-instance power](metalcloud-cli_server-instance_power.md)	 - Set power state of a server instance

