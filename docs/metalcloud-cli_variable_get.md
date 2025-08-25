## metalcloud-cli variable get

Get details of a specific variable

### Synopsis

Get detailed information about a specific variable by its ID.

This command retrieves and displays all details of a variable including its ID, name,
complete value, usage type, owner information, and timestamps.

Required Arguments:
  variable_id    Numeric ID of the variable to retrieve

Required Permissions:
  VARIABLES_AND_SECRETS_READ

Examples:
  # Get variable details by ID
  metalcloud-cli variable get 123
  
  # Get variable details using alias
  metalcloud-cli var show 456

```
metalcloud-cli variable get variable_id [flags]
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

* [metalcloud-cli variable](metalcloud-cli_variable.md)	 - Manage variables for infrastructure configuration

