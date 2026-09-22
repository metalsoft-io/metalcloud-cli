## metalcloud-cli server-instance power-status-batch

Get the power status of several server instances at once

### Synopsis

Get the power status of several server instances of one infrastructure in a
single call.

The instance IDs are sent in the request body; the infrastructure is resolved by ID
or by label.

Required Arguments:
  infrastructure_id_or_label  The ID or label of the infrastructure
  server_instance_id...       One or more numeric server instance IDs

Examples:
  metalcloud-cli server-instance power-status-batch 1234 5678 5679
  metalcloud-cli inst power-get-batch prod-env 5678

```
metalcloud-cli server-instance power-status-batch infrastructure_id_or_label server_instance_id... [flags]
```

### Options

```
  -h, --help   help for power-status-batch
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

