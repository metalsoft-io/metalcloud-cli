## metalcloud-cli user initiate-password-reset

Send a password reset message to an e-mail address

### Synopsis

Start the self-service password reset flow for an e-mail address.

The owner of the address receives a link that lets them choose a new password.

Required Flags:
  --email                 The e-mail address of the account to reset

Optional Flags:
  --redirect-url          URL the user is redirected to after resetting the password

Examples:
  metalcloud-cli user initiate-password-reset --email user@company.com
  metalcloud-cli user initiate-password-reset --email user@company.com --redirect-url https://metalsoft.io/login

```
metalcloud-cli user initiate-password-reset [flags]
```

### Options

```
      --email string          The e-mail address of the account to reset.
  -h, --help                  help for initiate-password-reset
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

