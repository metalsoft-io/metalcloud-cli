## metalcloud-cli external-connection interface update

Update an interface of an external connection

### Synopsis

Point an external connection interface at a different network device interface.

Required Arguments:
  external_connection_id_or_label   The ID or label of the external connection
  interface_id                      The ID of the external connection interface
  network_device_interface_id       The new network device interface ID

Examples:
  metalcloud external-connection interface update 12 5 102

```
metalcloud-cli external-connection interface update external_connection_id_or_label interface_id network_device_interface_id [flags]
```

### Options

```
  -h, --help   help for update
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

