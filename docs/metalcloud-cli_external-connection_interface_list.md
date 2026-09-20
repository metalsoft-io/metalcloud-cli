## metalcloud-cli external-connection interface list

List the interfaces of an external connection

### Synopsis

List the network device interfaces attached to an external connection.

Required Arguments:
  external_connection_id_or_label   The ID or label of the external connection

Examples:
  metalcloud external-connection interface list 12

```
metalcloud-cli external-connection interface list external_connection_id_or_label [flags]
```

### Options

```
  -h, --help   help for list
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

