## metalcloud-cli endpoint interface

Endpoint interface management

### Synopsis

Manage the individual network interfaces of an endpoint.

Each endpoint interface binds the endpoint to one switch port
(networkDeviceInterfaceId), optionally with the MAC address seen on it.

Available Commands:
  get      Get one endpoint interface
  add      Add an interface to an endpoint
  update   Update an existing endpoint interface
  remove   Remove an interface from an endpoint

Use 'endpoint interfaces endpoint_id' to list all interfaces of an endpoint.

### Options

```
  -h, --help   help for interface
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

* [metalcloud-cli endpoint](metalcloud-cli_endpoint.md)	 - Endpoint management
* [metalcloud-cli endpoint interface add](metalcloud-cli_endpoint_interface_add.md)	 - Add an interface to an endpoint
* [metalcloud-cli endpoint interface get](metalcloud-cli_endpoint_interface_get.md)	 - Get one interface of an endpoint
* [metalcloud-cli endpoint interface remove](metalcloud-cli_endpoint_interface_remove.md)	 - Remove an interface from an endpoint
* [metalcloud-cli endpoint interface update](metalcloud-cli_endpoint_interface_update.md)	 - Update an interface of an endpoint

