## metalcloud-cli vm-pool update-config-example

Display a VM pool update configuration example

### Synopsis

Display a sample VM pool update payload showing the fields accepted by
'vm-pool update'. Every field is optional - delete the ones you do not want to
change.

Examples:
  # Display the example
  metalcloud-cli vm-pool update-config-example

  # Save it for editing
  metalcloud-cli vm-pool update-config-example > update.json

  # Change one field and apply it
  metalcloud-cli vm-pool update-config-example | jq '{description}' | metalcloud-cli vm-pool update 123 --config-source pipe

```
metalcloud-cli vm-pool update-config-example [flags]
```

### Options

```
  -h, --help   help for update-config-example
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

* [metalcloud-cli vm-pool](metalcloud-cli_vm-pool.md)	 - Manage virtual machine pools and their resources

