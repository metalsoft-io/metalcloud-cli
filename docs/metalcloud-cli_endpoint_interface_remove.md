## metalcloud-cli endpoint interface remove

Remove an interface from an endpoint

### Synopsis

Remove an interface from an endpoint.

The interface's current revision is sent as the If-Match entity tag, so a
concurrent change is rejected instead of being overwritten.

Required Arguments:
  endpoint_id    The ID of the endpoint
  interface_id   The ID of the endpoint interface

Examples:
  metalcloud-cli endpoint interface remove 12 3

```
metalcloud-cli endpoint interface remove endpoint_id interface_id [flags]
```

### Options

```
  -h, --help   help for remove
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

* [metalcloud-cli endpoint interface](metalcloud-cli_endpoint_interface.md)	 - Endpoint interface management

