## metalcloud-cli container-instance-group network

Manage the network connections of a container instance group

### Synopsis

Manage the network connections of a container instance group.

Available commands:
  list        Show the network configuration of a group
  get         Get one network connection of a group
  connect     Connect a group to a logical network
  update      Update one network connection of a group
  disconnect  Remove one network connection of a group

Examples:
  metalcloud-cli container-instance-group network list my-infra 77
  metalcloud-cli cig net connect my-infra 77 5 trunk true

### Options

```
  -h, --help   help for network
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
* [metalcloud-cli container-instance-group network connect](metalcloud-cli_container-instance-group_network_connect.md)	 - Connect a container instance group to a logical network
* [metalcloud-cli container-instance-group network disconnect](metalcloud-cli_container-instance-group_network_disconnect.md)	 - Remove one network connection of a container instance group
* [metalcloud-cli container-instance-group network get](metalcloud-cli_container-instance-group_network_get.md)	 - Get one network connection of a container instance group
* [metalcloud-cli container-instance-group network list](metalcloud-cli_container-instance-group_network_list.md)	 - List the network connections of a container instance group
* [metalcloud-cli container-instance-group network update](metalcloud-cli_container-instance-group_network_update.md)	 - Update one network connection of a container instance group

