## metalcloud-cli user api-key-regenerate

Regenerate the current user's API key

### Synopsis

Regenerate the API key for the currently authenticated user.

WARNING: This will invalidate your current API key. You will need to update
any scripts or configurations that use the old key.

Examples:
  metalcloud-cli user api-key-regenerate

```
metalcloud-cli user api-key-regenerate [flags]
```

### Options

```
  -h, --help   help for api-key-regenerate
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

