## metalcloud-cli account

Manage user accounts and account-related operations

### Synopsis

Manage user accounts and account-related operations in the MetalCloud platform.

This command group provides functionality to:
- List all accounts in the system
- View detailed information about specific accounts
- Create new accounts with custom configurations
- Update existing account properties
- Archive accounts to disable them
- List users associated with an account

All account operations require appropriate permissions to perform user management tasks.

### Options

```
  -h, --help   help for account
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

* [metalcloud-cli](metalcloud-cli.md)	 - MetalCloud CLI
* [metalcloud-cli account archive](metalcloud-cli_account_archive.md)	 - Archive an account to disable it
* [metalcloud-cli account create](metalcloud-cli_account_create.md)	 - Create a new account with custom configuration
* [metalcloud-cli account get](metalcloud-cli_account_get.md)	 - Get detailed information about a specific account
* [metalcloud-cli account list](metalcloud-cli_account_list.md)	 - List all accounts in the system
* [metalcloud-cli account update](metalcloud-cli_account_update.md)	 - Update an existing account configuration
* [metalcloud-cli account users](metalcloud-cli_account_users.md)	 - List all users associated with a specific account

