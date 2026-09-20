## metalcloud-cli container-type create

Create a new container type

### Synopsis

Create a new container type.

The new container type is described by a JSON or YAML document. Use the
config-example command to obtain a template.

Required Flags:
  --config-source string  Source of the container type configuration.
                          Can be 'pipe' or path to a JSON/YAML file.

Examples:
  # Create a container type from a file
  metalcloud-cli container-type create --config-source container-type.json

  # Create a container type from stdin
  cat container-type.json | metalcloud-cli ct new --config-source pipe

```
metalcloud-cli container-type create [flags]
```

### Options

```
      --config-source string   Source of the new container type configuration. Can be 'pipe' or path to a JSON/YAML file.
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

* [metalcloud-cli container-type](metalcloud-cli_container-type.md)	 - Manage container types

