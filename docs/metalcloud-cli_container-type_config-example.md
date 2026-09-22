## metalcloud-cli container-type config-example

Print a container type configuration example

### Synopsis

Print a container type configuration example.

The printed document lists every field accepted by the create command and can be
saved to a file, edited and passed back through the --config-source flag.

Examples:
  # Show the example configuration
  metalcloud-cli container-type config-example

  # Save the example to a file
  metalcloud-cli ct config-example > container-type.json

```
metalcloud-cli container-type config-example [flags]
```

### Options

```
  -h, --help   help for config-example
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

