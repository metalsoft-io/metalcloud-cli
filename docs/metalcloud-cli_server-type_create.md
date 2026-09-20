## metalcloud-cli server-type create

Create a new server type

### Synopsis

Create a new server type from a JSON or YAML configuration.

The configuration describes the hardware the server type stands for: the
processor, memory, disk and network interface specifications together with the
server class and the optional GPU and disk group details.

Required Flags:
  --config-source    Source of the new server type configuration. Can be 'pipe' or path to a JSON file.

Examples:
  # Create a server type from a JSON file
  metalcloud-cli server-type create --config-source ./server-type.json

  # Create a server type from piped configuration
  metalcloud-cli server-type config-example | metalcloud-cli server-type create --config-source pipe


```
metalcloud-cli server-type create [flags]
```

### Options

```
      --config-source string   Source of the new server type configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for create
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

* [metalcloud-cli server-type](metalcloud-cli_server-type.md)	 - Manage server types and hardware configurations

