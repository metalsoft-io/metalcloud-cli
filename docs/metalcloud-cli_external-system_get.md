## metalcloud-cli external-system get

Get external system details

### Synopsis

Get the details of an external system, including its annotations.

Required Arguments:
  external_system_id   The numeric ID of the external system

Examples:
  metalcloud external-system get 12
  metalcloud ext-system show 12 -f json

```
metalcloud-cli external-system get external_system_id [flags]
```

### Options

```
  -h, --help   help for get
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

* [metalcloud-cli external-system](metalcloud-cli_external-system.md)	 - External system management

