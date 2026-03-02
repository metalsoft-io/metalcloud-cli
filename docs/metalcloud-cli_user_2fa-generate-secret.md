## metalcloud-cli user 2fa-generate-secret

Generate a new 2FA secret for setting up an authenticator app

### Synopsis

Generate a new two-factor authentication secret and QR code.

Use the generated secret or QR code to configure your authenticator app (e.g., Google
Authenticator, Authy), then call 'user 2fa-enable --token <code>' to activate 2FA.

Examples:
  metalcloud-cli user 2fa-generate-secret

```
metalcloud-cli user 2fa-generate-secret [flags]
```

### Options

```
  -h, --help   help for 2fa-generate-secret
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

