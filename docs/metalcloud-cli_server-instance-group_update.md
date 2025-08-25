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

Note: At least one of the flags (--label, --instance-count, --os-template-id) must be provided.

Examples:
  # Update the label of instance group 1234
  metalcloud-cli server-instance-group update 1234 --label "new-web-servers"

  # Scale instance group to 5 instances
  metalcloud-cli server-instance-group update 1234 --instance-count 5

  # Change OS template
  metalcloud-cli server-instance-group update 1234 --os-template-id 25

  # Update multiple properties at once
  metalcloud-cli server-instance-group update 1234 --label "updated-servers" --instance-count 3

  # Using alias
  metalcloud-cli ig edit 1234 --instance-count 10

```
metalcloud-cli server-instance-group update server_instance_group_id [flags]
```

### Options

```
  -h, --help                 help for update
      --instance-count int   Set the count of instance group instances.
      --label string         Set the instance group label.
      --os-template-id int   Set the instance group OS template Id.
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

