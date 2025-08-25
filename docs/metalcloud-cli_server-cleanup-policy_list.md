## metalcloud-cli server-cleanup-policy list

List all server cleanup policies

### Synopsis

List all server cleanup policies configured in the system.

This command displays a table of all available server cleanup policies with their
key attributes including ID, name, status, and configuration summary.

Output Format:
  By default, output is formatted as a table. Use global flags to change output format:
  --format=json    JSON output
  --format=csv     CSV output
  --format=yaml    YAML output

Required Permissions:
  - server_cleanup_policies:read

Examples:
  # List all server cleanup policies in table format
  metalcloud-cli server-cleanup-policy list

  # List policies in JSON format
  metalcloud-cli server-cleanup-policy list --format=json

  # List policies with custom output format
  metalcloud-cli scp ls --format=csv

```
metalcloud-cli server-cleanup-policy list [flags]
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

* [metalcloud-cli server-cleanup-policy](metalcloud-cli_server-cleanup-policy.md)	 - Manage server cleanup policies for automated server maintenance

