## metalcloud-cli logical-network

Manage logical networks within fabrics

### Synopsis

Manage logical networks within fabrics for network segmentation and isolation.

Logical networks provide Layer 2 network isolation within a fabric, allowing you to create
separate broadcast domains for different applications or tenants. Each logical network is
associated with a fabric and can have specific configurations based on its kind (vlan, vxlan, etc.).

Available Commands:
  list                 List logical networks with optional filtering
  get                  Get detailed information about a specific logical network
  create               Create a new logical network from configuration
  create-from-profile  Create a logical network from a logical network profile
  update               Update an existing logical network
  delete               Delete a logical network
  config-example       Get example configuration for a specific network kind
  get-config           Get the config sub-resource of a logical network
  update-config        Update the global settings of a logical network config
  apply-profiles       Apply a logical network profile to a logical network config
  allocation-strategy  Manage the config's allocation strategies
  get-external-connections                  List the attached external connections
  get-external-connection-logical-networks  List the external connection attachments
  get-interconnects                         List the attached logical network interconnects
  detach-external-connection                Detach an external connection

Examples:
  # List all logical networks
  metalcloud-cli logical-network list

  # List logical networks in a specific fabric
  metalcloud-cli logical-network list fabric-1

  # Create a VLAN logical network
  metalcloud-cli logical-network create vlan --config-source config.json

### Options

```
  -h, --help   help for logical-network
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

* [metalcloud-cli](metalcloud-cli.md)	 - MetalCloud CLI
* [metalcloud-cli logical-network allocation-strategy](metalcloud-cli_logical-network_allocation-strategy.md)	 - Manage logical network allocation strategies
* [metalcloud-cli logical-network apply-profiles](metalcloud-cli_logical-network_apply-profiles.md)	 - Apply a logical network profile to a logical network config
* [metalcloud-cli logical-network config-example](metalcloud-cli_logical-network_config-example.md)	 - Generate example configuration for a logical network kind
* [metalcloud-cli logical-network create](metalcloud-cli_logical-network_create.md)	 - Create a new logical network from configuration
* [metalcloud-cli logical-network create-from-profile](metalcloud-cli_logical-network_create-from-profile.md)	 - Create a logical network from a logical network profile
* [metalcloud-cli logical-network delete](metalcloud-cli_logical-network_delete.md)	 - Delete a logical network
* [metalcloud-cli logical-network detach-external-connection](metalcloud-cli_logical-network_detach-external-connection.md)	 - Detach an external connection from a logical network
* [metalcloud-cli logical-network get](metalcloud-cli_logical-network_get.md)	 - Get detailed information about a logical network
* [metalcloud-cli logical-network get-config](metalcloud-cli_logical-network_get-config.md)	 - Get the config of a logical network
* [metalcloud-cli logical-network get-external-connection-logical-networks](metalcloud-cli_logical-network_get-external-connection-logical-networks.md)	 - List the external connection attachments of a logical network
* [metalcloud-cli logical-network get-external-connections](metalcloud-cli_logical-network_get-external-connections.md)	 - List the external connections attached to a logical network
* [metalcloud-cli logical-network get-interconnects](metalcloud-cli_logical-network_get-interconnects.md)	 - List the logical network interconnects of a logical network
* [metalcloud-cli logical-network list](metalcloud-cli_logical-network_list.md)	 - List logical networks with optional filtering and sorting
* [metalcloud-cli logical-network update](metalcloud-cli_logical-network_update.md)	 - Update an existing logical network configuration
* [metalcloud-cli logical-network update-config](metalcloud-cli_logical-network_update-config.md)	 - Update the global settings of a logical network config

