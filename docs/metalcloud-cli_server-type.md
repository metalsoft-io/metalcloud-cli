## metalcloud-cli server-type

Manage server types and hardware configurations

### Synopsis

Manage server types and view detailed hardware specifications.

Server types define the hardware configurations available for provisioning,
including CPU, memory, storage, and network interface specifications.

Available Commands:
  list    List all available server types
  get     Get detailed information about a specific server type

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
* [metalcloud-cli server-type get](metalcloud-cli_server-type_get.md)	 - Get detailed information about a specific server type
* [metalcloud-cli server-type list](metalcloud-cli_server-type_list.md)	 - List all available server types

