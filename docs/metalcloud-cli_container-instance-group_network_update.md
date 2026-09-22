## metalcloud-cli container-instance-group network update

Update one network connection of a container instance group

### Synopsis

Update one network connection of a container instance group.

Required Arguments:
  infrastructure_id_or_label   The ID or the label of the infrastructure
  container_instance_group_id  The numeric ID of the container instance group
  connection_id                The numeric ID of the network connection

Required Flags:
  --access-mode string  The network connection access mode (e.g. 'trunk', 'access').
  --tagged string       Whether VLAN tagging is enabled (true/false).
  --redundancy string   The network connection redundancy mode.

Flag Dependencies:
  At least one of --access-mode, --tagged or --redundancy must be provided.

Examples:
  # Change the access mode of connection 5
  metalcloud-cli container-instance-group network update my-infra 77 5 --access-mode trunk

  # Enable VLAN tagging
  metalcloud-cli cig net edit 1234 77 5 --tagged true

```
metalcloud-cli container-instance-group network update infrastructure_id_or_label container_instance_group_id connection_id [flags]
```

### Options

```
      --access-mode string   Network connection access mode.
  -h, --help                 help for update
      --redundancy string    Network connection redundancy mode.
      --tagged string        Network connection VLAN tagging (true/false).
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

