## metalcloud-cli container-instance-group network connect

Connect a container instance group to a logical network

### Synopsis

Connect a container instance group to a logical network.

Required Arguments:
  infrastructure_id_or_label   The ID or the label of the infrastructure
  container_instance_group_id  The numeric ID of the container instance group
  logical_network_id           The ID of the logical network to connect to
  access_mode                  The network access mode (e.g. 'trunk', 'access')
  tagged                       Whether VLAN tagging is enabled (true/false)
  redundancy                   Optional. The redundancy mode of the connection

Examples:
  # Connect group 77 to logical network 5 in trunk mode
  metalcloud-cli container-instance-group network connect my-infra 77 5 trunk true

  # Connect with a redundancy mode
  metalcloud-cli cig net add 1234 77 5 trunk true active-backup

```
metalcloud-cli container-instance-group network connect infrastructure_id_or_label container_instance_group_id logical_network_id access_mode tagged [redundancy] [flags]
```

### Options

```
  -h, --help   help for connect
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

