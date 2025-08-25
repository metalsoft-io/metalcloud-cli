## metalcloud-cli storage list

List all storage pools

### Synopsis

List all storage pools with optional filtering.

This command displays information about storage pools including their ID, site, driver,
technology, type, name, and status. The output can be filtered by storage technology.

Flags:
  --filter-technology strings   Filter results by storage technology (e.g., block, file, object)

Examples:
  # List all storage pools
  metalcloud storage list

  # List only block storage pools
  metalcloud storage list --filter-technology block

  # List multiple storage types
  metalcloud storage list --filter-technology block,file

```
metalcloud-cli storage list [flags]
```

### Options

```
      --filter-technology strings   Filter the result by storage technology.
  -h, --help                        help for list
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

