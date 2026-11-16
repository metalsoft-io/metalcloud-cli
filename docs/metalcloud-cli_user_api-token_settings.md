## metalcloud-cli user api-token settings

Show the API token limits of the controller

### Synopsis

Show the default validity, the maximum number of tokens per user and the
latest expiry a new or regenerated token may have.

Examples:
  metalcloud-cli user api-token settings

```
metalcloud-cli user api-token settings [flags]
```

### Options

```
  -h, --help   help for settings
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

* [metalcloud-cli user api-token](metalcloud-cli_user_api-token.md)	 - Manage the current user's API tokens

