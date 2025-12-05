## metalcloud-cli server-registration-profile list

List all server registration profiles

### Synopsis

List all server registration profiles configured in the system.

This command displays a table of all available server registration profiles with their
key attributes including ID, name, status, and configuration summary.

Output Format:
  By default, output is formatted as a table. Use global flags to change output format:
  --format=json    JSON output
  --format=csv     CSV output
  --format=yaml    YAML output

Examples:
  # List all server registration profiles in table format
  metalcloud-cli server-registration-profile list

  # List policies in JSON format
  metalcloud-cli server-registration-profile list --format=json

```
metalcloud-cli server-registration-profile list [flags]
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

* [metalcloud-cli server-registration-profile](metalcloud-cli_server-registration-profile.md)	 - Manage server registration profiles

