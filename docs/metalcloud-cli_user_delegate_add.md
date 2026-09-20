## metalcloud-cli user delegate add

Grant a user delegate access to another user

### Synopsis

Grant a delegate user access to the resources of a user.

After this command the delegate user can operate on the resources owned by the user
identified by user_id.

Required Arguments:
  user_id                 The numeric ID of the user whose resources are delegated
  delegate_id             The numeric ID of the user receiving the delegate access

Examples:
  metalcloud-cli user delegate add 12345 67890

```
metalcloud-cli user delegate add user_id delegate_id [flags]
```

### Options

```
  -h, --help   help for add
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

