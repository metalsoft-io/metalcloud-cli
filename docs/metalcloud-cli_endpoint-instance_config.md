## metalcloud-cli endpoint-instance config

Get the pending configuration of an endpoint instance

### Synopsis

Get the pending (not yet deployed) configuration of an endpoint instance.

Required Arguments:
  endpoint_instance_id   The numeric ID of the endpoint instance

Examples:
  metalcloud-cli endpoint-instance config 42
  metalcloud-cli ei get-config 42

```
metalcloud-cli endpoint-instance config endpoint_instance_id [flags]
```

### Options

```
  -h, --help   help for config
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

