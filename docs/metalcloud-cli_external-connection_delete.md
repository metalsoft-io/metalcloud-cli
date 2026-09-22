## metalcloud-cli external-connection delete

Delete an external connection

### Synopsis

Delete an external connection together with its interfaces.

Required Arguments:
  external_connection_id_or_label   The ID or label of the external connection

Examples:
  metalcloud external-connection delete 12

```
metalcloud-cli external-connection delete external_connection_id_or_label [flags]
```

### Options

```
  -h, --help   help for delete
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

