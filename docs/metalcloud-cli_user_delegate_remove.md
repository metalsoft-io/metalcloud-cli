## metalcloud-cli user delegate remove

Revoke the delegate access of a user

### Synopsis

Revoke the delegate access previously granted to a user.

Required Arguments:
  user_id                 The numeric ID of the user whose resources were delegated
  delegate_id             The numeric ID of the user losing the delegate access

Examples:
  metalcloud-cli user delegate remove 12345 67890
  metalcloud-cli user delegate rm 12345 67890

```
metalcloud-cli user delegate remove user_id delegate_id [flags]
```

### Options

```
  -h, --help   help for remove
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

* [metalcloud-cli user delegate](metalcloud-cli_user_delegate.md)	 - Manage user delegates

