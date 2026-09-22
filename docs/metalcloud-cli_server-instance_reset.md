## metalcloud-cli server-instance reset

Reset the deployed server of a server instance

### Synopsis

Reset the deployed server of a server instance. The operation is executed
immediately.

This is a different operation from 'server-instance power <id> reset': the power
command goes to the power-set endpoint and only cycles the power of the server,
while 'reset' asks the orchestration layer to reset the deployed server.

Required Arguments:
  server_instance_id  The numeric ID of the server instance

Examples:
  metalcloud-cli server-instance reset 5678
  metalcloud-cli inst reset 5678

```
metalcloud-cli server-instance reset server_instance_id [flags]
```

### Options

```
  -h, --help   help for reset
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

* [metalcloud-cli server-instance](metalcloud-cli_server-instance.md)	 - Manage individual server instances

