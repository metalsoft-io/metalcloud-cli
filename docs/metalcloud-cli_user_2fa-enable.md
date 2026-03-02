## metalcloud-cli user 2fa-enable

Enable two-factor authentication

### Synopsis

Enable two-factor authentication for the current user.

You must first generate a 2FA secret using 'user 2fa-generate-secret', configure
your authenticator app, then provide the TOTP token to verify and enable 2FA.

Required Flags:
  --token string    The TOTP code from your authenticator app

Examples:
  # First generate the secret
  metalcloud-cli user 2fa-generate-secret

  # Then enable 2FA with the token from your authenticator app
  metalcloud-cli user 2fa-enable --token 123456

```
metalcloud-cli user 2fa-enable [flags]
```

### Options

```
  -h, --help           help for 2fa-enable
      --token string   The TOTP code from your authenticator app.
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

