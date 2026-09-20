## metalcloud-cli network-fabric-interconnect

Network fabric interconnect management

### Synopsis

Manage network fabric interconnects (e.g. EVPN data center interconnects) that
join two or more network fabrics through BGP.

Command categories:
  Lifecycle:   list, get, create, update, delete, config-example, template
  Links:       get-links, get-link, add-link, remove-link, activate-links, deactivate-links
  Fabrics:     get-fabrics, get-available-fabrics
  Deployment:  deploy, deployment-check, deployment-info, accept-deploy, reject-deploy, detach

### Options

```
  -h, --help   help for network-fabric-interconnect
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
* [metalcloud-cli network-fabric-interconnect accept-deploy](metalcloud-cli_network-fabric-interconnect_accept-deploy.md)	 - Accept a pending interconnect deploy
* [metalcloud-cli network-fabric-interconnect activate-links](metalcloud-cli_network-fabric-interconnect_activate-links.md)	 - Activate interconnect links
* [metalcloud-cli network-fabric-interconnect add-link](metalcloud-cli_network-fabric-interconnect_add-link.md)	 - Add a link to an interconnect
* [metalcloud-cli network-fabric-interconnect config-example](metalcloud-cli_network-fabric-interconnect_config-example.md)	 - Print an example interconnect create configuration
* [metalcloud-cli network-fabric-interconnect create](metalcloud-cli_network-fabric-interconnect_create.md)	 - Create a network fabric interconnect
* [metalcloud-cli network-fabric-interconnect deactivate-links](metalcloud-cli_network-fabric-interconnect_deactivate-links.md)	 - Deactivate interconnect links
* [metalcloud-cli network-fabric-interconnect delete](metalcloud-cli_network-fabric-interconnect_delete.md)	 - Delete a network fabric interconnect
* [metalcloud-cli network-fabric-interconnect deploy](metalcloud-cli_network-fabric-interconnect_deploy.md)	 - Deploy a network fabric interconnect
* [metalcloud-cli network-fabric-interconnect deployment-check](metalcloud-cli_network-fabric-interconnect_deployment-check.md)	 - Validate whether the interconnect links can be activated
* [metalcloud-cli network-fabric-interconnect deployment-info](metalcloud-cli_network-fabric-interconnect_deployment-info.md)	 - Show the deployment status and preview of an interconnect
* [metalcloud-cli network-fabric-interconnect detach](metalcloud-cli_network-fabric-interconnect_detach.md)	 - Detach an interconnect, removing its configuration from the devices
* [metalcloud-cli network-fabric-interconnect get](metalcloud-cli_network-fabric-interconnect_get.md)	 - Get network fabric interconnect details
* [metalcloud-cli network-fabric-interconnect get-available-fabrics](metalcloud-cli_network-fabric-interconnect_get-available-fabrics.md)	 - List the fabrics that can be linked to an interconnect
* [metalcloud-cli network-fabric-interconnect get-fabrics](metalcloud-cli_network-fabric-interconnect_get-fabrics.md)	 - List the fabrics attached to an interconnect
* [metalcloud-cli network-fabric-interconnect get-link](metalcloud-cli_network-fabric-interconnect_get-link.md)	 - Get one link of an interconnect
* [metalcloud-cli network-fabric-interconnect get-links](metalcloud-cli_network-fabric-interconnect_get-links.md)	 - List the links of an interconnect
* [metalcloud-cli network-fabric-interconnect list](metalcloud-cli_network-fabric-interconnect_list.md)	 - List network fabric interconnects
* [metalcloud-cli network-fabric-interconnect reject-deploy](metalcloud-cli_network-fabric-interconnect_reject-deploy.md)	 - Reject a pending interconnect deploy
* [metalcloud-cli network-fabric-interconnect remove-link](metalcloud-cli_network-fabric-interconnect_remove-link.md)	 - Remove a link from an interconnect
* [metalcloud-cli network-fabric-interconnect template](metalcloud-cli_network-fabric-interconnect_template.md)	 - Show the BGP templates used for an interconnect type
* [metalcloud-cli network-fabric-interconnect update](metalcloud-cli_network-fabric-interconnect_update.md)	 - Update a network fabric interconnect

