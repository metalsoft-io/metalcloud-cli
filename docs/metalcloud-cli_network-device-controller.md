## metalcloud-cli network-device-controller

Network device controller management

### Synopsis

Manage network device controllers: the fabric controllers (Cisco NDFC/ACI,
NVIDIA UFM, Brocade) that MetalSoft drives instead of talking to each switch directly.

Command categories:
  Lifecycle:   list, get, create, update, delete, config-example
  Operations:  get-credentials, deploy-confirm

### Options

```
  -h, --help   help for network-device-controller
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
* [metalcloud-cli network-device-controller config-example](metalcloud-cli_network-device-controller_config-example.md)	 - Print an example controller create configuration
* [metalcloud-cli network-device-controller create](metalcloud-cli_network-device-controller_create.md)	 - Create a network device controller
* [metalcloud-cli network-device-controller delete](metalcloud-cli_network-device-controller_delete.md)	 - Delete a network device controller
* [metalcloud-cli network-device-controller deploy-confirm](metalcloud-cli_network-device-controller_deploy-confirm.md)	 - Confirm a pending controller deploy
* [metalcloud-cli network-device-controller get](metalcloud-cli_network-device-controller_get.md)	 - Get network device controller details
* [metalcloud-cli network-device-controller get-credentials](metalcloud-cli_network-device-controller_get-credentials.md)	 - Get the management credentials of a controller
* [metalcloud-cli network-device-controller list](metalcloud-cli_network-device-controller_list.md)	 - List network device controllers
* [metalcloud-cli network-device-controller update](metalcloud-cli_network-device-controller_update.md)	 - Update a network device controller

