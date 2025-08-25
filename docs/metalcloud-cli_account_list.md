## metalcloud-cli account list

List all accounts in the system

### Synopsis

List all accounts in the MetalCloud platform.

This command displays a table of all accounts including their ID, name, status, 
creation date, and other relevant information. The output can be formatted as 
JSON, YAML, or table format.

Required Permissions:
  - users:read

Examples:
  # List all accounts in table format
  metalcloud-cli account list

  # List all accounts in JSON format  
  metalcloud-cli account list -o json

  # List all accounts using alias
  metalcloud-cli accounts ls

```
metalcloud-cli account list [flags]
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

* [metalcloud-cli account](metalcloud-cli_account.md)	 - Manage user accounts and account-related operations

