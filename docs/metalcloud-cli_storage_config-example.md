## metalcloud-cli storage config-example

Display a configuration template for creating storage pools

### Synopsis

Display a configuration template for creating storage pools.

This command outputs a JSON template showing all available configuration options
for creating a storage pool. The template includes required fields, optional fields,
and example values to help you create valid storage configurations.

The output can be used as a starting point for creating storage pool configurations
that can be passed to the 'create' command via the --config-source flag.

Examples:
  # Display the configuration template
  metalcloud storage config-example

  # Save the template to a file for editing
  metalcloud storage config-example > storage-config.json

```
metalcloud-cli storage config-example [flags]
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

* [metalcloud-cli storage](metalcloud-cli_storage.md)	 - Manage storage pools and related resources

