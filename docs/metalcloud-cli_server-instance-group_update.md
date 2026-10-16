## metalcloud-cli server-instance-group update

Update server instance group configuration

### Synopsis

Update server instance group configuration.

This command allows you to modify the configuration of an existing server instance group.
You can update the label, instance count, or OS template. At least one flag must be provided.

Arguments:
  server_instance_group_id  The numeric ID of the server instance group to update

Flags:
  --label string           Set the instance group label
  --instance-count int     Set the count of instance group instances (must be > 0)
  --os-template-id int     Set the instance group OS template ID (must be > 0)
  --server-id strings            Pin the group to a specific server (can be repeated). Shorthand for
                                 --server-filter-rule server_id:equals:<id>
  --server-filter-rule strings   Restrict the servers the group may use, as property:operator[:value]
                                 (can be repeated). Properties: tags, rack, server_id, inventory_ID.
                                 Operators: equals, contains, does_not_contain, starts_with,
                                 does_not_starts_with, matches_regexp, affinity, anti_affinity.

Note: At least one of the flags (--label, --instance-count, --os-template-id, --server-id,
--server-filter-rule) must be provided. When given, the server filter rules replace the group's
existing rules.

Examples:
  # Update the label of instance group 1234
  metalcloud-cli server-instance-group update 1234 --label "new-web-servers"

  # Scale instance group to 5 instances
  metalcloud-cli server-instance-group update 1234 --instance-count 5

  # Change OS template
  metalcloud-cli server-instance-group update 1234 --os-template-id 25

  # Pin the group to server 236
  metalcloud-cli server-instance-group update 1234 --server-id 236

  # Update multiple properties at once
  metalcloud-cli server-instance-group update 1234 --label "updated-servers" --instance-count 3

  # Using alias
  metalcloud-cli ig edit 1234 --instance-count 10

```
metalcloud-cli server-instance-group update server_instance_group_id [flags]
```

### Options

```
  -h, --help                             help for update
      --instance-count int               Set the count of instance group instances.
      --label string                     Set the instance group label.
      --os-template-id int               Set the instance group OS template Id.
      --server-filter-rule stringArray   Restrict the servers the group may use, as property:operator[:value]. Replaces the existing rules. Can be repeated.
      --server-id strings                Pin the group to a specific server ID. Replaces the existing rules. Can be repeated.
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

