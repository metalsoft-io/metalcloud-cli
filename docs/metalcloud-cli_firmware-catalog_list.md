## metalcloud-cli firmware-catalog list

List all firmware catalogs

### Synopsis

List all firmware catalogs in the system.

This command displays all available firmware catalogs with their basic information including:
- Catalog ID and name
- Vendor type (Dell, HP, Lenovo)
- Update type (online/offline)
- Creation date and status

No additional flags are required for this command.

```
metalcloud-cli firmware-catalog list [flags]
```

### Options

```
  -h, --help   help for list
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

* [metalcloud-cli firmware-catalog](metalcloud-cli_firmware-catalog.md)	 - Manage firmware catalogs for server hardware updates

