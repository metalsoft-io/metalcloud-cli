## metalcloud-cli variable list

List all variables

### Synopsis

List all variables available in the current account.

This command displays all variables that have been created, showing their ID, name, 
value (truncated if long), usage type, owner, and timestamps.

Required Permissions:
  VARIABLES_AND_SECRETS_READ

Examples:
  # List all variables
  metalcloud-cli variable list
  
  # List variables using alias
  metalcloud-cli var ls

```
metalcloud-cli variable list [flags]
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

* [metalcloud-cli variable](metalcloud-cli_variable.md)	 - Manage variables for infrastructure configuration

