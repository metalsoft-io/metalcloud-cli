## metalcloud-cli user send-password-reset

Send a password reset message to a user (admin)

### Synopsis

Send a password reset message to a specific user as an administrator.

The user receives a link that lets them choose a new password. Use 'user set-password'
instead to set a password directly without involving the user.

Required Arguments:
  user_id                 The numeric ID of the user to notify

Optional Flags:
  --redirect-url          URL the user is redirected to after resetting the password

Examples:
  metalcloud-cli user send-password-reset 12345
  metalcloud-cli user send-password-reset 12345 --redirect-url https://metalsoft.io/login

```
metalcloud-cli user send-password-reset user_id [flags]
```

### Options

```
  -h, --help                  help for send-password-reset
      --redirect-url string   URL the user is redirected to after resetting the password.
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

