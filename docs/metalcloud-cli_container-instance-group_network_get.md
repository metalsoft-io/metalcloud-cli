## metalcloud-cli container-instance-group network get

Get one network connection of a container instance group

### Synopsis

Get detailed information about one network connection of a container instance group.

Required Arguments:
  infrastructure_id_or_label   The ID or the label of the infrastructure
  container_instance_group_id  The numeric ID of the container instance group
  connection_id                The numeric ID of the network connection

Examples:
  # Get network connection 5 of group 77
  metalcloud-cli container-instance-group network get my-infra 77 5

  # Using the alias
  metalcloud-cli cig net show 1234 77 5

```
metalcloud-cli container-instance-group network get infrastructure_id_or_label container_instance_group_id connection_id [flags]
```

### Options

```
  -h, --help   help for get
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

* [metalcloud-cli container-instance-group network](metalcloud-cli_container-instance-group_network.md)	 - Manage the network connections of a container instance group

