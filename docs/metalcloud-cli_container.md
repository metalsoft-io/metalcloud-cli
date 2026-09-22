## metalcloud-cli container

Manage provisioned containers

### Synopsis

Manage provisioned containers.

Containers are the runtime objects allocated for container instances. They are
created by the platform when an infrastructure is deployed; this command group
inspects them and controls their power state.

Available commands:
  list                 List all containers
  get                  Get details of a container
  update               Update container comments and tags
  power-status         Show the current power state of a container
  start                Power on a container
  shutdown             Power off a container
  reboot               Reboot a container
  remote-console-info  Show remote console information for a container

Examples:
  metalcloud-cli container list
  metalcloud-cli containers get 100
  metalcloud-cli container power-status 100

### Options

```
  -h, --help   help for container
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
* [metalcloud-cli container get](metalcloud-cli_container_get.md)	 - Get container details
* [metalcloud-cli container list](metalcloud-cli_container_list.md)	 - List all containers
* [metalcloud-cli container power-status](metalcloud-cli_container_power-status.md)	 - Get the power status of a container
* [metalcloud-cli container reboot](metalcloud-cli_container_reboot.md)	 - Reboot a container
* [metalcloud-cli container remote-console-info](metalcloud-cli_container_remote-console-info.md)	 - Get remote console information for a container
* [metalcloud-cli container shutdown](metalcloud-cli_container_shutdown.md)	 - Power off a container
* [metalcloud-cli container start](metalcloud-cli_container_start.md)	 - Power on a container
* [metalcloud-cli container update](metalcloud-cli_container_update.md)	 - Update a container

