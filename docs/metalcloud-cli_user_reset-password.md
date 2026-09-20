## metalcloud-cli user reset-password

Consume a password reset token

### Synopsis

Consume a password reset token, as following the link from the e-mail would.

The token is the value of the 'token' query parameter of the password reset link that
the platform sent by e-mail. The API exposes no password parameter on this endpoint:
it only validates the token and redirects to the page where the new password is
chosen. Use 'user change-password' to set a password from the CLI, or
'user set-password' to set the password of another user as an administrator.

Required Flags:
  --token                 The password reset token

Examples:
  metalcloud-cli user reset-password --token eyJhbGciOi...

```
metalcloud-cli user reset-password [flags]
```

### Options

```
  -h, --help           help for reset-password
      --token string   The password reset token.
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

