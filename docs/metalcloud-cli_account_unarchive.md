## metalcloud-cli account unarchive

Restore a previously archived account

### Synopsis

Restore a previously archived account in the MetalCloud platform.

This command reverses 'account archive': the account becomes usable again while all
its data and configuration are preserved. The account is identified by its unique
account ID.

Required Permissions:
  - users:write

Arguments:
  account_id    The unique identifier of the account to restore

Examples:
  # Restore an account by ID
  metalcloud-cli account unarchive 1234

  # Using aliases
  metalcloud-cli account restore 1234
  metalcloud-cli account unar 1234

```
metalcloud-cli account unarchive account_id [flags]
```

### Options

```
  -h, --help   help for unarchive
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

