## metalcloud-cli server-instance interfaces

List the interfaces of a server instance

### Synopsis

List all network interfaces of a server instance.

Required Arguments:
  server_instance_id  The numeric ID of the server instance

Optional Flags:
  --filter-infrastructure-id strings     Filter by infrastructure ID. Repeatable or comma-separated.
  --filter-service-status strings        Filter by service status.
  --filter-config-deploy-status strings  Filter by the deploy status of the pending configuration.
  --filter-config-deploy-type strings    Filter by the deploy type of the pending configuration.

Examples:
  metalcloud-cli server-instance interfaces 5678
  metalcloud-cli inst ifaces 5678 --filter-service-status active

```
metalcloud-cli server-instance interfaces server_instance_id [flags]
```

### Options

```
      --filter-config-deploy-status strings   Filter by the deploy status of the pending configuration.
      --filter-config-deploy-type strings     Filter by the deploy type of the pending configuration.
      --filter-infrastructure-id strings      Filter by infrastructure ID.
      --filter-service-status strings         Filter by service status.
  -h, --help                                  help for interfaces
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

