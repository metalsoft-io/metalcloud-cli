## metalcloud-cli vm-instance-group network

Manage the network connections of a VM instance group

### Synopsis

Manage the network connections of a VM instance group.

Available commands:
  list            Show the network configuration of a group
  connections     List the network connections of a group
  get             Get one network connection of a group
  connect         Connect a group to a logical network
  update          Update one network connection of a group
  disconnect      Remove one network connection of a group
  config-example  Print a network connection configuration example

Examples:
  metalcloud-cli vm-instance-group network list my-infra 67890
  metalcloud-cli vmg net connect my-infra 67890 --logical-network-id 5 --access-mode l2 --tagged true

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

* [metalcloud-cli vm-instance-group](metalcloud-cli_vm-instance-group.md)	 - Manage VM instance groups within infrastructures
* [metalcloud-cli vm-instance-group network config-example](metalcloud-cli_vm-instance-group_network_config-example.md)	 - Print a network connection configuration example
* [metalcloud-cli vm-instance-group network connect](metalcloud-cli_vm-instance-group_network_connect.md)	 - Connect a VM instance group to a logical network
* [metalcloud-cli vm-instance-group network connections](metalcloud-cli_vm-instance-group_network_connections.md)	 - List the network connections of a VM instance group
* [metalcloud-cli vm-instance-group network disconnect](metalcloud-cli_vm-instance-group_network_disconnect.md)	 - Remove one network connection of a VM instance group
* [metalcloud-cli vm-instance-group network get](metalcloud-cli_vm-instance-group_network_get.md)	 - Get one network connection of a VM instance group
* [metalcloud-cli vm-instance-group network list](metalcloud-cli_vm-instance-group_network_list.md)	 - Show the network configuration of a VM instance group
* [metalcloud-cli vm-instance-group network update](metalcloud-cli_vm-instance-group_network_update.md)	 - Update one network connection of a VM instance group

