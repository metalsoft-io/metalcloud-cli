## metalcloud-cli container-instance-group apply-type

Apply a container type on a container instance group

### Synopsis

Apply a container type on a container instance group.

The current group revision is fetched automatically and sent as the If-Match
header.

Required Arguments:
  infrastructure_id_or_label   The ID or the label of the infrastructure
  container_instance_group_id  The numeric ID of the container instance group
  container_type_id            The numeric ID of the container type to apply

Examples:
  # Apply container type 42 on group 77
  metalcloud-cli container-instance-group apply-type my-infra 77 42

  # Using the alias
  metalcloud-cli cig set-type 1234 77 42

```
metalcloud-cli container-instance-group apply-type infrastructure_id_or_label container_instance_group_id container_type_id [flags]
```

### Options

```
  -h, --help   help for apply-type
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

* [metalcloud-cli container-instance-group](metalcloud-cli_container-instance-group.md)	 - Manage container instance groups within infrastructures

