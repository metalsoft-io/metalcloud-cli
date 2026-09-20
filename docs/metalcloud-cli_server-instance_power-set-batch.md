## metalcloud-cli server-instance power-set-batch

Set the power state of several server instances at once

### Synopsis

Set the power state of several server instances of one infrastructure in a
single call.

The instance IDs and the power command are sent in the request body; the
infrastructure is resolved by ID or by label.

Valid power actions:
  on     - Power on the server instances
  off    - Power off the server instances
  reset  - Hard reset the server instances
  soft   - Graceful shutdown of the server instances

Required Arguments:
  infrastructure_id_or_label  The ID or label of the infrastructure
  action                      Power action to perform (on, off, reset, soft)
  server_instance_id...       One or more numeric server instance IDs

Examples:
  metalcloud-cli server-instance power-set-batch 1234 off 5678 5679
  metalcloud-cli inst power-batch prod-env on 5678

```
metalcloud-cli server-instance power-set-batch infrastructure_id_or_label <on|off|reset|soft> server_instance_id... [flags]
```

### Options

```
  -h, --help   help for power-set-batch
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

