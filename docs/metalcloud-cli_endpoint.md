## metalcloud-cli endpoint

Endpoint management

### Synopsis

Endpoint management commands.

An endpoint represents a device attached to the fabric. Its interfaces bind the
endpoint to individual switch ports.

Available Commands:
  list, get, create, create-bulk, update, delete
  interfaces                      List the interfaces of an endpoint
  interface                       Get, add, update or remove one interface
  get-network-device-interfaces   List a switch's ports and their endpoints

### Options

```
  -h, --help   help for endpoint
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
* [metalcloud-cli endpoint create](metalcloud-cli_endpoint_create.md)	 - Create a new endpoint
* [metalcloud-cli endpoint create-bulk](metalcloud-cli_endpoint_create-bulk.md)	 - Create multiple endpoints in one call
* [metalcloud-cli endpoint delete](metalcloud-cli_endpoint_delete.md)	 - Delete an endpoint
* [metalcloud-cli endpoint get](metalcloud-cli_endpoint_get.md)	 - Display detailed information about a specific endpoint
* [metalcloud-cli endpoint get-network-device-interfaces](metalcloud-cli_endpoint_get-network-device-interfaces.md)	 - List a network device's interfaces and their endpoints
* [metalcloud-cli endpoint interface](metalcloud-cli_endpoint_interface.md)	 - Endpoint interface management
* [metalcloud-cli endpoint interfaces](metalcloud-cli_endpoint_interfaces.md)	 - List interfaces of an endpoint
* [metalcloud-cli endpoint list](metalcloud-cli_endpoint_list.md)	 - List endpoints
* [metalcloud-cli endpoint update](metalcloud-cli_endpoint_update.md)	 - Update an existing endpoint

