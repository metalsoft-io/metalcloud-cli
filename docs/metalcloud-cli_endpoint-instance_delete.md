## metalcloud-cli endpoint-instance delete

Delete an endpoint instance

### Synopsis

Delete an endpoint instance.

The current revision of the instance is fetched first and sent as If-Match, so the
delete fails if the instance changed in the meantime.

Required Arguments:
  endpoint_instance_id   The numeric ID of the endpoint instance

Examples:
  metalcloud-cli endpoint-instance delete 42
  metalcloud-cli ei rm 42

```
metalcloud-cli endpoint-instance delete endpoint_instance_id [flags]
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

* [metalcloud-cli endpoint-instance](metalcloud-cli_endpoint-instance.md)	 - Endpoint instance management

