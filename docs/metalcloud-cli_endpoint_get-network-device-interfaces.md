## metalcloud-cli endpoint get-network-device-interfaces

List a network device's interfaces and their endpoints

### Synopsis

List every interface of a network device together with the endpoint attached to it.

Interfaces with no endpoint are listed with empty endpoint columns, so this is
the command to use when looking for a free switch port.

Required Arguments:
  network_device_id  The ID or label (identifier string) of the network device

Examples:
  metalcloud-cli endpoint get-network-device-interfaces 45
  metalcloud-cli endpoint get-network-device-interfaces leaf-01

```
metalcloud-cli endpoint get-network-device-interfaces network_device_id [flags]
```

### Options

```
  -h, --help   help for get-network-device-interfaces
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

