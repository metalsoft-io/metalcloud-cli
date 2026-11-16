## metalcloud-cli user api-token delete

Delete an API token of the current user

### Synopsis

Delete an API token of the currently authenticated user. The token stops
working immediately.

Examples:
  metalcloud-cli user api-token delete 12

```
metalcloud-cli user api-token delete api_token_id [flags]
```

### Options

```
  -h, --help   help for delete
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

