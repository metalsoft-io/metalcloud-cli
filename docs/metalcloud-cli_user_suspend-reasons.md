## metalcloud-cli user suspend-reasons

List the suspend reasons recorded for a user

### Synopsis

List the suspend reasons recorded for a specific user account.

Each entry shows the type of the suspension, the public and private comments left by
the administrator and the interval during which the suspension was active.

Required Arguments:
  user_id                 The numeric ID of the user whose suspend reasons to list

Examples:
  metalcloud-cli user suspend-reasons 12345
  metalcloud-cli user get-suspend-reasons 12345

```
metalcloud-cli user suspend-reasons user_id [flags]
```

### Options

```
  -h, --help   help for suspend-reasons
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

