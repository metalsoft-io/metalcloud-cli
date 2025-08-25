## metalcloud-cli resource-pool get-users

List users with access to a resource pool

### Synopsis

List all users that have access to a specific resource pool.

This command retrieves and displays a list of users who have been granted access
to the specified resource pool. The output includes user details and their
permissions within the resource pool.

Arguments:
  pool_id    The numeric ID of the resource pool

Examples:
  # List users for resource pool with ID 123
  metalcloud-cli resource-pool get-users 123

  # List users using alias
  metalcloud-cli rp users 456

```
metalcloud-cli resource-pool get-users <pool_id> [flags]
```

### Options

```
  -h, --help   help for get-users
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

