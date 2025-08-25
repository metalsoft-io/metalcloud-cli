## metalcloud-cli server-type get

Get detailed information about a specific server type

### Synopsis

Get detailed information about a specific server type by its ID.

This command retrieves comprehensive hardware specifications for a specific server type,
including detailed processor information, memory configuration, storage details,
network interface specifications, GPU information (if applicable), and other
hardware characteristics.

The detailed output includes:
- Server type ID, name, and label
- Complete processor specifications (count, speed, core count, names)
- Memory configuration (RAM in GB)
- Storage information (disk count and disk groups)
- Network interface details (count, speeds, total capacity)
- GPU information (count and detailed GPU info)
- Server class and boot type
- Various flags (experimental, unmanaged servers only, etc.)
- Allowed vendor SKU IDs
- Tags associated with the server type

Arguments:
  server-type-id    The numeric ID of the server type to retrieve

Examples:
  # Get information about server type with ID 123
  metalcloud server-type get 123

  # Get server type information (using alias)
  metalcloud server-type show 456

  # Get server type info with JSON output
  metalcloud server-type get 789 --output json

Required Permissions:
  - Server Types Read

Output Format:
  The command outputs data in table format by default. Use global output flags
  to change the format (--output json, --output yaml, etc.).

```
metalcloud-cli server-type get <server-type-id> [flags]
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

* [metalcloud-cli server-type](metalcloud-cli_server-type.md)	 - Manage server types and hardware configurations

