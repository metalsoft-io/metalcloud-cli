## metalcloud-cli external-connection interface remove

Remove an interface from an external connection

### Synopsis

Detach a network device interface from an external connection.

Required Arguments:
  external_connection_id_or_label   The ID or label of the external connection
  interface_id                      The ID of the external connection interface

Examples:
  metalcloud external-connection interface remove 12 5

```
metalcloud-cli external-connection interface remove external_connection_id_or_label interface_id [flags]
```

### Options

```
  -h, --help   help for remove
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

