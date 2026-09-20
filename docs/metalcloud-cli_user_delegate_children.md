## metalcloud-cli user delegate children

List the users a user delegated access to

### Synopsis

List the delegate users that can act on the resources of a specific user.

Required Arguments:
  user_id                 The numeric ID of the user whose delegates to list

Examples:
  metalcloud-cli user delegate children 12345

```
metalcloud-cli user delegate children user_id [flags]
```

### Options

```
  -h, --help   help for children
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

