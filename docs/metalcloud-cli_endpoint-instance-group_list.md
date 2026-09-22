## metalcloud-cli endpoint-instance-group list

List the endpoint instance groups of an infrastructure

### Synopsis

List all endpoint instance groups of an infrastructure.

Required Arguments:
  infrastructure_id_or_label   The ID or label of the infrastructure

Optional Flags:
  --filter-extension-instance-id strings  Filter by extension instance ID. Repeatable or comma-separated.
  --filter-service-status strings         Filter by service status.
  --filter-config-deploy-status strings   Filter by the deploy status of the pending configuration.
  --filter-config-deploy-type strings     Filter by the deploy type of the pending configuration.

Examples:
  metalcloud-cli endpoint-instance-group list 1234
  metalcloud-cli eig ls prod-env --filter-service-status active

```
metalcloud-cli endpoint-instance-group list infrastructure_id_or_label [flags]
```

### Options

```
      --filter-config-deploy-status strings    Filter by the deploy status of the pending configuration.
      --filter-config-deploy-type strings      Filter by the deploy type of the pending configuration.
      --filter-extension-instance-id strings   Filter by extension instance ID.
      --filter-service-status strings          Filter by service status.
  -h, --help                                   help for list
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

* [metalcloud-cli endpoint-instance-group](metalcloud-cli_endpoint-instance-group.md)	 - Endpoint instance group management

