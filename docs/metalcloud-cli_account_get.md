## metalcloud-cli account get

Get detailed information about a specific account

### Synopsis

Get detailed information about a specific account in the MetalCloud platform.

This command displays comprehensive details about an account including its 
configuration, status, creation date, and associated metadata. The account
is identified by its unique account ID.

Required Permissions:
  - users:read

Arguments:
  account_id    The unique identifier of the account to retrieve

Examples:
  # Get account details by ID
  metalcloud-cli account get 1234

  # Get account details in JSON format
  metalcloud-cli account get 1234 -o json

  # Using alias
  metalcloud-cli account show 1234

```
metalcloud-cli account get account_id [flags]
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

* [metalcloud-cli account](metalcloud-cli_account.md)	 - Manage user accounts and account-related operations

