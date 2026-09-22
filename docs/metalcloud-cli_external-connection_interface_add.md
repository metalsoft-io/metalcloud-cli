## metalcloud-cli external-connection interface add

Add an interface to an external connection

### Synopsis

Attach a network device interface to an external connection.

Required Arguments:
  external_connection_id_or_label   The ID or label of the external connection
  network_device_interface_id       The ID of the network device interface (see 'get-network-device-interfaces')

Examples:
  metalcloud external-connection interface add 12 101

```
metalcloud-cli external-connection interface add external_connection_id_or_label network_device_interface_id [flags]
```

### Options

```
  -h, --help   help for add
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

* [metalcloud-cli external-connection interface](metalcloud-cli_external-connection_interface.md)	 - External connection interface management

