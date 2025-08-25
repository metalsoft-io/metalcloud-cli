## metalcloud-cli resource-pool add-user

Grant a user access to a resource pool

### Synopsis

Grant a user access to a resource pool by specifying the resource pool ID and user ID.

This command adds a user to a resource pool, giving them access to the resources
within that pool according to their role permissions.

Arguments:
  pool_id    The numeric ID of the resource pool
  user_id    The numeric ID of the user to add

Examples:
  # Add user with ID 789 to resource pool with ID 123
  metalcloud-cli resource-pool add-user 123 789

  # Add user using alias
  metalcloud-cli rp add-user 456 789

```
metalcloud-cli resource-pool add-user <pool_id> <user_id> [flags]
```

### Options

```
  -h, --help   help for add-user
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

