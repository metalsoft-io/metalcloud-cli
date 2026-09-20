## metalcloud-cli account config

Display the configuration of an account

### Synopsis

Display the configuration object of a specific account.

The configuration holds the properties written by 'account update': name, code, fiscal
number, contacts, parent account, quota profile and API key validity duration.

Required Permissions:
  - users:read

Arguments:
  account_id    The unique identifier of the account

Examples:
  # Show the configuration of an account
  metalcloud-cli account config 1234

  # Using alias
  metalcloud-cli account get-config 1234

```
metalcloud-cli account config account_id [flags]
```

### Options

```
  -h, --help   help for config
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

