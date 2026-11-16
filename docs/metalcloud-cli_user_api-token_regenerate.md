## metalcloud-cli user api-token regenerate

Regenerate an API token of the current user

### Synopsis

Regenerate an API token of the currently authenticated user.

WARNING: The previous token value stops working. Update any scripts or
configurations that use it.

Flags:
  --expires               New expiry as an RFC 3339 timestamp or a YYYY-MM-DD date (UTC);
                          defaults to the controller's default token validity

Examples:
  metalcloud-cli user api-token regenerate 12
  metalcloud-cli user api-token regenerate 12 --expires 2027-06-30

```
metalcloud-cli user api-token regenerate api_token_id [flags]
```

### Options

```
      --expires string   New expiry as an RFC 3339 timestamp or YYYY-MM-DD date (UTC).
  -h, --help             help for regenerate
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

