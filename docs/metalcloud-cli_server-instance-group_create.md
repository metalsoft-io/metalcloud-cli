## metalcloud-cli server-instance-group create

Create a new server instance group in an infrastructure

### Synopsis

Create a new server instance group in an infrastructure.

This command creates a new server instance group with the specified configuration.
The group will contain multiple server instances with identical specifications.

Arguments:
  infrastructure_id_or_label  The infrastructure ID (numeric) or label (string) to create the group in
  label                       Label for the new instance group
  server_type_id             Server type ID to use for instances in the group
  instance_count             Number of instances to create in the group
  os_template_id             (Optional) OS template ID to use for instances

Flags:
  --server-id strings             Pin the group to a specific server (can be repeated). Shorthand for
                                  --server-filter-rule server_id:equals:<id>
  --server-filter-rule strings    Restrict the servers the group may use, as property:operator[:value]
                                  (can be repeated). Properties: tags, rack, server_id, inventory_ID.
                                  Operators: equals, contains, does_not_contain, starts_with,
                                  does_not_starts_with, matches_regexp, affinity, anti_affinity.
                                  The value is optional for affinity and anti_affinity.

Examples:
  # Create instance group with 3 instances using server type 100
  metalcloud-cli server-instance-group create 1234 web-servers 100 3

  # Create instance group with specific OS template
  metalcloud-cli server-instance-group create prod-env db-cluster 200 2 50

  # Deploy the group on a specific server
  metalcloud-cli server-instance-group create 1234 web-servers 100 1 --server-id 236

  # Only use servers in a given rack
  metalcloud-cli server-instance-group create 1234 web-servers 100 3 --server-filter-rule rack:equals:A12

  # Using alias
  metalcloud-cli ig new 1234 app-servers 150 5

```
metalcloud-cli server-instance-group create infrastructure_id_or_label label server_type_id instance_count [os_template_id] [flags]
```

### Options

```
  -h, --help                             help for create
      --server-filter-rule stringArray   Restrict the servers the group may use, as property:operator[:value]. Can be repeated.
      --server-id strings                Pin the group to a specific server ID. Can be repeated.
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

* [metalcloud-cli server-instance-group](metalcloud-cli_server-instance-group.md)	 - Manage server instance groups within infrastructures

