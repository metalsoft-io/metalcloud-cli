## metalcloud-cli vm-instance config-example

Print a VM instance configuration example

### Synopsis

Print a VM instance configuration example.

The printed document lists every field accepted by the create command and can be
saved to a file, edited and passed back through the --config-source flag.

Examples:
  # Show the example configuration
  metalcloud-cli vm-instance config-example

  # Save the example to a file
  metalcloud-cli vmi config-example > vm-instance.json

```
metalcloud-cli vm-instance config-example [flags]
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

* [metalcloud-cli vm-instance](metalcloud-cli_vm-instance.md)	 - Manage individual VM instances within infrastructures

