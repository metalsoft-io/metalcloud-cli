## metalcloud-cli container-type

Manage container types

### Synopsis

Manage container types.

Container types describe the hardware envelope (CPU cores, RAM, GPUs) that a
container instance receives when it is provisioned. They are global objects and
are referenced by container instances and container instance groups.

Available commands:
  list            List all container types
  get             Get details of a container type
  config-example  Print an example container type configuration
  create          Create a new container type
  update          Update an existing container type
  delete          Delete a container type
  containers      List the containers provisioned with a container type

Examples:
  metalcloud-cli container-type list
  metalcloud-cli ct get 42
  metalcloud-cli ct containers 42

### Options

```
  -h, --help   help for container-type
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
* [metalcloud-cli container-type config-example](metalcloud-cli_container-type_config-example.md)	 - Print a container type configuration example
* [metalcloud-cli container-type containers](metalcloud-cli_container-type_containers.md)	 - List the containers of a container type
* [metalcloud-cli container-type create](metalcloud-cli_container-type_create.md)	 - Create a new container type
* [metalcloud-cli container-type delete](metalcloud-cli_container-type_delete.md)	 - Delete a container type
* [metalcloud-cli container-type get](metalcloud-cli_container-type_get.md)	 - Get container type details
* [metalcloud-cli container-type list](metalcloud-cli_container-type_list.md)	 - List all container types
* [metalcloud-cli container-type update](metalcloud-cli_container-type_update.md)	 - Update a container type

