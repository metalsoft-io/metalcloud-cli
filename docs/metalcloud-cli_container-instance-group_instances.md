## metalcloud-cli container-instance-group instances

List the container instances of a group

### Synopsis

List the container instances that belong to a container instance group.

Required Arguments:
  infrastructure_id_or_label   The ID or the label of the infrastructure
  container_instance_group_id  The numeric ID of the container instance group

Examples:
  # List the instances of group 77
  metalcloud-cli container-instance-group instances my-infra 77

  # Using the alias
  metalcloud-cli cig list-instances 1234 77

```
metalcloud-cli container-instance-group instances infrastructure_id_or_label container_instance_group_id [flags]
```

### Options

```
  -h, --help   help for instances
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

