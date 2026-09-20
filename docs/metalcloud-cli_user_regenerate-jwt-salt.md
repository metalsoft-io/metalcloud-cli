## metalcloud-cli user regenerate-jwt-salt

Regenerate the JWT salt of the current user

### Synopsis

Regenerate the JWT salt of the user owning the API key in use.

WARNING: this invalidates every session and token issued so far for this user, on
every device and in every browser. You will have to log in again.

Examples:
  metalcloud-cli user regenerate-jwt-salt

```
metalcloud-cli user regenerate-jwt-salt [flags]
```

### Options

```
  -h, --help   help for regenerate-jwt-salt
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

