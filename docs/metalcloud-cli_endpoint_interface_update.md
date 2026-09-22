## metalcloud-cli endpoint interface update

Update an interface of an endpoint

### Synopsis

Update an existing endpoint interface.

The updates can be described either by individual flags or by a JSON/YAML
configuration passed with --config-source. The interface's current revision is
sent as the If-Match entity tag.

Required Arguments:
  endpoint_id    The ID of the endpoint
  interface_id   The ID of the endpoint interface

Optional Flags:
  --network-device-interface-id  Move the interface to another switch port
  --mac-address                  The new MAC address of the interface
  --config-source                'pipe' to read from stdin, or a path to a JSON/YAML file

Examples:
  metalcloud-cli endpoint interface update 12 3 --mac-address AA:BB:CC:DD:EE:FF
  metalcloud-cli endpoint interface update 12 3 --config-source interface.json

```
metalcloud-cli endpoint interface update endpoint_id interface_id [flags]
```

### Options

```
      --config-source string              Source of the endpoint interface updates. Can be 'pipe' or path to a JSON file.
  -h, --help                              help for update
      --mac-address string                The new MAC address of the endpoint interface.
      --network-device-interface-id int   The new network device interface (switch port) ID.
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

* [metalcloud-cli endpoint interface](metalcloud-cli_endpoint_interface.md)	 - Endpoint interface management

