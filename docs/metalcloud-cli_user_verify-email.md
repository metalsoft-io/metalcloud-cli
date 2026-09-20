## metalcloud-cli user verify-email

Verify an e-mail address with a verification token

### Synopsis

Consume an e-mail verification token, as following the link from the e-mail would.

The token is the value of the 'token' query parameter of the verification link that
the platform sent by e-mail.

Required Flags:
  --token                 The e-mail verification token

Examples:
  metalcloud-cli user verify-email --token eyJhbGciOi...

```
metalcloud-cli user verify-email [flags]
```

### Options

```
  -h, --help           help for verify-email
      --token string   The e-mail verification token.
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

