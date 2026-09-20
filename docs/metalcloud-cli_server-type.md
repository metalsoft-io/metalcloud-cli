## metalcloud-cli server-type

Manage server types and hardware configurations

### Synopsis

Manage server types and view detailed hardware specifications.

Server types define the hardware configurations available for provisioning,
including CPU, memory, storage, and network interface specifications.

Available Commands:
  list            List all available server types
  get             Get detailed information about a specific server type
  create          Create a new server type
  update          Update an existing server type
  delete          Delete a server type
  clean-unused    Remove the server types that are no longer used
  statistics      Get the server availability statistics of a site
  config-example  Show a server type creation configuration example

Use "metalcloud server-type [command] --help" for more information about a command.

### Options

```
  -h, --help   help for server-type
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
* [metalcloud-cli server-type clean-unused](metalcloud-cli_server-type_clean-unused.md)	 - Remove the server types that are no longer used
* [metalcloud-cli server-type config-example](metalcloud-cli_server-type_config-example.md)	 - Show a server type creation configuration example
* [metalcloud-cli server-type create](metalcloud-cli_server-type_create.md)	 - Create a new server type
* [metalcloud-cli server-type delete](metalcloud-cli_server-type_delete.md)	 - Delete a server type
* [metalcloud-cli server-type get](metalcloud-cli_server-type_get.md)	 - Get detailed information about a specific server type
* [metalcloud-cli server-type list](metalcloud-cli_server-type_list.md)	 - List all available server types
* [metalcloud-cli server-type statistics](metalcloud-cli_server-type_statistics.md)	 - Get the server availability statistics of a site
* [metalcloud-cli server-type update](metalcloud-cli_server-type_update.md)	 - Update an existing server type
* [metalcloud-cli server-type update-config-example](metalcloud-cli_server-type_update-config-example.md)	 - Show a server type update configuration example

