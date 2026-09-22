## metalcloud-cli logical-network detach-external-connection

Detach an external connection from a logical network

### Synopsis

Detach an external connection from a logical network.

Required Arguments:
  logical_network_id      The ID of the logical network
  external_connection_id  The ID of the external connection to detach

Examples:
  metalcloud-cli logical-network detach-external-connection 12 4

```
metalcloud-cli logical-network detach-external-connection logical_network_id external_connection_id [flags]
```

### Options

```
  -h, --help   help for detach-external-connection
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

* [metalcloud-cli logical-network](metalcloud-cli_logical-network.md)	 - Manage logical networks within fabrics

