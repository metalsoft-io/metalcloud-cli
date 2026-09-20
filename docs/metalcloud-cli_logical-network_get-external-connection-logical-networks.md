## metalcloud-cli logical-network get-external-connection-logical-networks

List the external connection attachments of a logical network

### Synopsis

List the external connection logical networks of a logical network.

Each record is one attachment between the logical network and an external
connection, with the status of that attachment.

Required Arguments:
  logical_network_id  The ID of the logical network

Examples:
  metalcloud-cli logical-network get-external-connection-logical-networks 12

```
metalcloud-cli logical-network get-external-connection-logical-networks logical_network_id [flags]
```

### Options

```
  -h, --help   help for get-external-connection-logical-networks
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

