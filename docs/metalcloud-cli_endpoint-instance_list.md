## metalcloud-cli endpoint-instance list

List endpoint instances

### Synopsis

List endpoint instances.

Without flags every endpoint instance visible to the user is listed. Pass
--infrastructure to restrict the listing to a single infrastructure; the
infrastructure is resolved by ID or by label.

Optional Flags:
  --infrastructure string               List only the endpoint instances of this infrastructure (ID or label).
  --filter-infrastructure-id strings    Filter by infrastructure ID. Repeatable or comma-separated.
  --filter-group-id strings             Filter by endpoint instance group ID.
  --filter-endpoint-id strings          Filter by endpoint ID.
  --filter-service-status strings       Filter by service status.
  --filter-config-endpoint-id strings   Filter by the endpoint ID of the pending configuration.
  --filter-config-deploy-status strings Filter by the deploy status of the pending configuration.
  --filter-config-deploy-type strings   Filter by the deploy type of the pending configuration.

Examples:
  metalcloud-cli endpoint-instance list
  metalcloud-cli endpoint-instance list --infrastructure prod-env
  metalcloud-cli ei ls --filter-service-status active

```
metalcloud-cli endpoint-instance list [flags]
```

### Options

```
      --filter-config-deploy-status strings   Filter by the deploy status of the pending configuration.
      --filter-config-deploy-type strings     Filter by the deploy type of the pending configuration.
      --filter-config-endpoint-id strings     Filter by the endpoint ID of the pending configuration.
      --filter-endpoint-id strings            Filter by endpoint ID.
      --filter-group-id strings               Filter by endpoint instance group ID.
      --filter-infrastructure-id strings      Filter by infrastructure ID.
      --filter-service-status strings         Filter by service status.
  -h, --help                                  help for list
      --infrastructure string                 List only the endpoint instances of this infrastructure (ID or label).
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

