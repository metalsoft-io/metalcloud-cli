## metalcloud-cli container-instance-group network list

List the network connections of a container instance group

### Synopsis

List all network connections of a container instance group.

Required Arguments:
  infrastructure_id_or_label   The ID or the label of the infrastructure
  container_instance_group_id  The numeric ID of the container instance group

Optional Flags:
  --configuration  Show the network endpoint group configuration instead of the
                   list of connections.

Examples:
  # List the network connections of group 77
  metalcloud-cli container-instance-group network list my-infra 77

  # Show the network endpoint group configuration
  metalcloud-cli cig net ls 1234 77 --configuration

```
metalcloud-cli container-instance-group network list infrastructure_id_or_label container_instance_group_id [flags]
```

### Options

```
      --configuration   Show the network endpoint group configuration instead of the list of connections.
  -h, --help            help for list
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

