## metalcloud-cli endpoint interfaces

List interfaces of an endpoint

### Synopsis

List all network interfaces of a specific endpoint in MetalSoft.

This command displays detailed information about all network interfaces associated 
with the specified endpoint, including their configuration and status.

Arguments:
  endpoint_id    The unique identifier of the endpoint whose interfaces to list (required)

Examples:
  metalcloud-cli endpoint interfaces 123
  metalcloud-cli endpoint ifaces 456
  metalcloud-cli ep ifs endpoint-uuid-123

```
metalcloud-cli endpoint interfaces endpoint_id [flags]
```

### Options

```
  -h, --help   help for interfaces
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

