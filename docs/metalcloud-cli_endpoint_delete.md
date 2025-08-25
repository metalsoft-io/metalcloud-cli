## metalcloud-cli endpoint delete

Delete an endpoint

### Synopsis

Delete an endpoint from MetalSoft.

This command permanently removes an endpoint from the system. This action cannot be undone.

Arguments:
  endpoint_id    The unique identifier of the endpoint to delete (required)

Examples:
  metalcloud-cli endpoint delete 123
  metalcloud-cli endpoint rm 456
  metalcloud-cli ep del endpoint-uuid-123

Warning: This operation is irreversible. Make sure you have the correct endpoint ID before proceeding.

```
metalcloud-cli endpoint delete endpoint_id [flags]
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

* [metalcloud-cli endpoint](metalcloud-cli_endpoint.md)	 - Endpoint management

