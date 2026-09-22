## metalcloud-cli resource-pool list-for-user

List the resource pools a user has access to

### Synopsis

List all resource pools that a given user has been granted access to.

This is the user-centric counterpart of 'resource-pool get-users', which lists the
users of one pool.

Required Arguments:
  user_id    The numeric ID of the user

Examples:
  # List the resource pools of user 42
  metalcloud-cli resource-pool list-for-user 42

  # Using an alias
  metalcloud-cli rp user-pools 42

```
metalcloud-cli resource-pool list-for-user user_id [flags]
```

### Options

```
  -h, --help   help for list-for-user
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

* [metalcloud-cli resource-pool](metalcloud-cli_resource-pool.md)	 - Manage resource pools and their associated resources

