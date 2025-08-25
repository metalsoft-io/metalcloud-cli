## metalcloud-cli firmware-binary list

List all firmware binaries

### Synopsis

List all firmware binaries in the system.

This command displays all available firmware binaries with their basic information including:
- Binary ID and name
- Catalog ID and package information
- Target hardware components and models
- Version information and update severity
- Reboot requirements and vendor details
- Release timestamps and download URLs
- External ID and vendor info URLs

The output includes both catalog-managed binaries and individually registered ones.

No additional flags are required for this command.

Examples:
  metalcloud-cli firmware-binary list
  metalcloud-cli fw-binary ls

```
metalcloud-cli firmware-binary list [flags]
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

* [metalcloud-cli firmware-binary](metalcloud-cli_firmware-binary.md)	 - Manage individual firmware binary files and packages

