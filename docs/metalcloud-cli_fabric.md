## metalcloud-cli fabric

Manage network fabrics

### Synopsis

Manage network fabrics in MetalCloud.

Fabrics are logical network constructs that group network devices and define how they are interconnected.
This command provides operations to create, configure, activate, and manage fabric devices.

Available Commands:
  list           List all fabrics
  get            Get fabric details
  create         Create a new fabric
  update         Update fabric configuration
  activate       Activate a fabric
  config-example Show configuration example
  get-devices    List fabric devices
  add-device     Add devices to fabric
  remove-device  Remove device from fabric

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
* [metalcloud-cli fabric activate](metalcloud-cli_fabric_activate.md)	 - Activate a fabric
* [metalcloud-cli fabric add-device](metalcloud-cli_fabric_add-device.md)	 - Add network device(s) to a fabric
* [metalcloud-cli fabric config-example](metalcloud-cli_fabric_config-example.md)	 - Show example fabric configuration
* [metalcloud-cli fabric create](metalcloud-cli_fabric_create.md)	 - Create a new fabric
* [metalcloud-cli fabric get](metalcloud-cli_fabric_get.md)	 - Get detailed fabric information
* [metalcloud-cli fabric get-devices](metalcloud-cli_fabric_get-devices.md)	 - List devices in a fabric
* [metalcloud-cli fabric list](metalcloud-cli_fabric_list.md)	 - List all network fabrics
* [metalcloud-cli fabric remove-device](metalcloud-cli_fabric_remove-device.md)	 - Remove network device from a fabric
* [metalcloud-cli fabric update](metalcloud-cli_fabric_update.md)	 - Update fabric configuration

