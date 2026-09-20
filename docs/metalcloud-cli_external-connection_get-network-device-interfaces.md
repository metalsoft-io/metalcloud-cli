## metalcloud-cli external-connection get-network-device-interfaces

List a network device's interfaces and their external connections

### Synopsis

List all interfaces of a network device together with the external connection
each interface belongs to, if any. Use it to find the network device interface IDs
accepted by 'external-connection interface add'.

Required Arguments:
  network_device_id   The ID or identifier of the network device

Examples:
  metalcloud external-connection get-network-device-interfaces 45
  metalcloud ext-conn get-network-device-interfaces border-leaf-01

```
metalcloud-cli external-connection get-network-device-interfaces network_device_id [flags]
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

* [metalcloud-cli external-connection](metalcloud-cli_external-connection.md)	 - External connection management

