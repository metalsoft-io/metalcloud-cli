## metalcloud-cli endpoint interface add

Add an interface to an endpoint

### Synopsis

Add a new interface to an existing endpoint.

The interface can be described either by individual flags or by a JSON/YAML
configuration passed with --config-source.

Required Arguments:
  endpoint_id  The ID of the endpoint

Required Flags (one of):
  --network-device-interface-id  The numeric ID of the switch port to attach
  --config-source                'pipe' to read from stdin, or a path to a JSON/YAML file

Optional Flags:
  --mac-address  The MAC address seen on the interface

Examples:
  metalcloud-cli endpoint interface add 12 --network-device-interface-id 4567
  metalcloud-cli endpoint interface add 12 --network-device-interface-id 4567 --mac-address AA:BB:CC:DD:EE:FF
  metalcloud-cli endpoint interface add 12 --config-source interface.json

```
metalcloud-cli endpoint interface add endpoint_id [flags]
```

### Options

```
      --config-source string              Source of the new endpoint interface configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                              help for add
      --mac-address string                The MAC address of the endpoint interface.
      --network-device-interface-id int   The network device interface (switch port) ID to attach.
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

