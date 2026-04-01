## metalcloud-cli user set-password

Set the password for a user (admin only)

### Synopsis

Set the password for a specific user account as an administrator.

Arguments:
  user_id                 The numeric ID of the user whose password to set

Required Flags:
  --password              The new password to set for the user

Examples:
  metalcloud-cli user set-password 12345 --password newSecret123

```
metalcloud-cli user set-password user_id [flags]
```

### Options

```
  -h, --help              help for set-password
      --password string   The new password to set for the user.
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

