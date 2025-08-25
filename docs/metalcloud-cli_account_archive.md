## metalcloud-cli account archive

Archive an account to disable it

### Synopsis

Archive an account in the MetalCloud platform to disable it.

This command archives an account, effectively disabling it while preserving 
its data and configuration. Archived accounts cannot be used for new operations 
but their historical data remains accessible. The account is identified by its 
unique account ID.

Note: This operation is typically irreversible. Archived accounts may require 
administrator intervention to reactivate.

Required Permissions:
  - users:write

Arguments:
  account_id    The unique identifier of the account to archive

Examples:
  # Archive an account by ID
  metalcloud-cli account archive 1234

  # Using alias
  metalcloud-cli account ar 1234

```
metalcloud-cli account archive account_id [flags]
```

### Options

```
  -h, --help   help for archive
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

