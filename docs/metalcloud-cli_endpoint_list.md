## metalcloud-cli endpoint list

List endpoints

### Synopsis

List all endpoints in MetalSoft with optional filtering.

This command displays a table of all endpoints available in the system. You can filter the results 
by site or external ID to narrow down the output.

Flags:
  --filter-site strings           Filter results by site name(s). Can be specified multiple times.
  --filter-external-id strings    Filter results by external ID(s). Can be specified multiple times.

Examples:
  metalcloud-cli endpoint list
  metalcloud-cli endpoint ls --filter-site "site1" --filter-site "site2"
  metalcloud-cli endpoint list --filter-external-id "ext-001"
  metalcloud-cli endpoint list --filter-site "production" --filter-external-id "api-endpoint"

```
metalcloud-cli endpoint list [flags]
```

### Options

```
      --filter-external-id strings   Filter the result by endpoint external Id.
      --filter-site strings          Filter the result by site.
  -h, --help                         help for list
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

