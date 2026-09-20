## metalcloud-cli external-connection logical-network add

Attach a logical network to an external connection

### Synopsis

Attach a logical network to an external connection.

Required Arguments:
  external_connection_id_or_label   The ID or label of the external connection
  logical_network_id                The numeric ID of the logical network

Examples:
  metalcloud external-connection logical-network add 12 44

```
metalcloud-cli external-connection logical-network add external_connection_id_or_label logical_network_id [flags]
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

* [metalcloud-cli external-connection logical-network](metalcloud-cli_external-connection_logical-network.md)	 - External connection logical network management

