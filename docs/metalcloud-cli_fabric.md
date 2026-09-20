## metalcloud-cli fabric

Manage network fabrics

### Synopsis

Manage network fabrics in MetalCloud.

Fabrics are logical network constructs that group network devices and define how they are interconnected.
This command provides operations to create, configure, activate, and manage fabric devices.

Available Commands:
  list              List all fabrics
  get               Get fabric details
  create            Create a new fabric
  update            Update fabric configuration
  delete            Delete a fabric
  activate          Activate a fabric
  deploy            Deploy a fabric
  accept-deploy     Accept a pending fabric deployment
  reject-deploy     Reject a pending fabric deployment
  config-example    Show configuration example
  get-devices       List fabric devices
  add-device        Add devices to fabric
  remove-device     Remove device from fabric
  get-links         List fabric links
  get-link          Get one fabric link
  add-link          Add fabric link
  remove-link       Remove fabric link
  get-interconnects List the fabric's network fabric interconnects
  bgp-session       Manage the fabric's BGP sessions
  link-aggregation  Manage the fabric's link aggregations

### Options

```
  -h, --help   help for fabric
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
* [metalcloud-cli fabric accept-deploy](metalcloud-cli_fabric_accept-deploy.md)	 - Accept a pending fabric deployment
* [metalcloud-cli fabric activate](metalcloud-cli_fabric_activate.md)	 - Activate a fabric
* [metalcloud-cli fabric add-device](metalcloud-cli_fabric_add-device.md)	 - Add network device(s) to a fabric
* [metalcloud-cli fabric add-link](metalcloud-cli_fabric_add-link.md)	 - Add a network fabric link
* [metalcloud-cli fabric bgp-session](metalcloud-cli_fabric_bgp-session.md)	 - Manage the BGP sessions of a fabric
* [metalcloud-cli fabric config-example](metalcloud-cli_fabric_config-example.md)	 - Show example fabric configuration
* [metalcloud-cli fabric configure-bgp](metalcloud-cli_fabric_configure-bgp.md)	 - Register the BGP underlay/overlay/PFC templates + profiles (step 8b)
* [metalcloud-cli fabric configure-bgp-example](metalcloud-cli_fabric_configure-bgp-example.md)	 - Show an example config for configure-bgp
* [metalcloud-cli fabric configure-freeform](metalcloud-cli_fabric_configure-freeform.md)	 - Register the base freeform template + per-switch profiles (step 8a)
* [metalcloud-cli fabric configure-freeform-example](metalcloud-cli_fabric_configure-freeform-example.md)	 - Show an example config for configure-freeform
* [metalcloud-cli fabric configure-switches](metalcloud-cli_fabric_configure-switches.md)	 - Configure all switches of a fabric from a declarative YAML/JSON
* [metalcloud-cli fabric configure-switches-example](metalcloud-cli_fabric_configure-switches-example.md)	 - Show an example switch configuration for configure-switches
* [metalcloud-cli fabric create](metalcloud-cli_fabric_create.md)	 - Create a new fabric
* [metalcloud-cli fabric delete](metalcloud-cli_fabric_delete.md)	 - Delete a fabric
* [metalcloud-cli fabric deploy](metalcloud-cli_fabric_deploy.md)	 - Deploy a fabric
* [metalcloud-cli fabric get](metalcloud-cli_fabric_get.md)	 - Get detailed fabric information
* [metalcloud-cli fabric get-devices](metalcloud-cli_fabric_get-devices.md)	 - List devices in a fabric
* [metalcloud-cli fabric get-interconnects](metalcloud-cli_fabric_get-interconnects.md)	 - List the network fabric interconnects of a fabric
* [metalcloud-cli fabric get-link](metalcloud-cli_fabric_get-link.md)	 - Get one fabric link
* [metalcloud-cli fabric get-links](metalcloud-cli_fabric_get-links.md)	 - List links in a fabric
* [metalcloud-cli fabric import-devices](metalcloud-cli_fabric_import-devices.md)	 - Bulk-import switches and attach them to a fabric (idempotent)
* [metalcloud-cli fabric link-aggregation](metalcloud-cli_fabric_link-aggregation.md)	 - Manage the link aggregations of a fabric
* [metalcloud-cli fabric list](metalcloud-cli_fabric_list.md)	 - List all network fabrics
* [metalcloud-cli fabric reject-deploy](metalcloud-cli_fabric_reject-deploy.md)	 - Reject a pending fabric deployment
* [metalcloud-cli fabric remove-device](metalcloud-cli_fabric_remove-device.md)	 - Remove network device from a fabric
* [metalcloud-cli fabric remove-link](metalcloud-cli_fabric_remove-link.md)	 - Remove a network fabric link
* [metalcloud-cli fabric rescan-links](metalcloud-cli_fabric_rescan-links.md)	 - Re-scan (discover) the fabric's links from LLDP
* [metalcloud-cli fabric update](metalcloud-cli_fabric_update.md)	 - Update fabric configuration

