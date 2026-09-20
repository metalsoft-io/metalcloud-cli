## metalcloud-cli logical-network get-config

Get the config of a logical network

### Synopsis

Display the config object of a logical network.

The config is a separate sub-resource holding the desired state of the network:
its kind, MTU, deploy type and status, and the allocation strategies that the
'logical-network allocation-strategy' commands manage. It carries its own
revision, distinct from the logical network's.

Required Arguments:
  logical_network_id  The ID of the logical network

Examples:
  metalcloud-cli logical-network get-config 12
  metalcloud-cli logical-network get-config 12 -f json

```
metalcloud-cli logical-network get-config logical_network_id [flags]
```

### Options

```
  -h, --help   help for get-config
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

