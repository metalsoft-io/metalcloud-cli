## metalcloud-cli server-instance list

List server instances

### Synopsis

List server instances.

Without an argument every server instance visible to the user is listed. Pass an
infrastructure ID or label to restrict the listing to a single infrastructure.

Optional Arguments:
  infrastructure_id_or_label  List only the server instances of this infrastructure

Optional Flags:
  --filter-infrastructure-id strings     Filter by infrastructure ID. Repeatable or comma-separated.
  --filter-group-id strings              Filter by server instance group ID.
  --filter-server-id strings             Filter by server ID.
  --filter-service-status strings        Filter by service status.
  --filter-config-server-id strings      Filter by the server ID of the pending configuration.
  --filter-config-deploy-status strings  Filter by the deploy status of the pending configuration.
  --filter-config-deploy-type strings    Filter by the deploy type of the pending configuration.

Examples:
  metalcloud-cli server-instance list
  metalcloud-cli server-instance list 1234
  metalcloud-cli inst ls prod-env --filter-service-status active

```
metalcloud-cli server-instance list [infrastructure_id_or_label] [flags]
```

### Options

```
      --filter-config-deploy-status strings   Filter by the deploy status of the pending configuration.
      --filter-config-deploy-type strings     Filter by the deploy type of the pending configuration.
      --filter-config-server-id strings       Filter by the server ID of the pending configuration.
      --filter-group-id strings               Filter by server instance group ID.
      --filter-infrastructure-id strings      Filter by infrastructure ID.
      --filter-server-id strings              Filter by server ID.
      --filter-service-status strings         Filter by service status.
  -h, --help                                  help for list
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

