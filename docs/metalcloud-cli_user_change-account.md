## metalcloud-cli user change-account

Move a user to a different account

### Synopsis

Move a user from their current account to a different account in the system.

This command transfers user ownership between accounts while preserving all user data,
settings, and permissions. The user will be associated with the new account immediately
after the command executes successfully.

Arguments:
  user_id                 The numeric ID of the user to move

Required Flags:
  --account-id            The ID of the destination account to move the user to

Examples:
  metalcloud-cli user change-account 12345 --account-id 67890
  metalcloud-cli user move-account 12345 --account-id 67890

```
metalcloud-cli user change-account user_id [flags]
```

### Options

```
      --account-id int   The ID of the account to move the user to.
  -h, --help             help for change-account
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

* [metalcloud-cli user](metalcloud-cli_user.md)	 - Manage user accounts and their properties

