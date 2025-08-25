## metalcloud-cli resource-pool get-subnet-pools

List subnet pools assigned to a resource pool

### Synopsis

List all subnet pools that are assigned to a specific resource pool.

This command retrieves and displays a list of subnet pools that have been assigned
to the specified resource pool. The output includes subnet pool details and their
current configuration within the resource pool.

Arguments:
  pool_id    The numeric ID of the resource pool

Examples:
  # List subnet pools for resource pool with ID 123
  metalcloud-cli resource-pool get-subnet-pools 123

  # List subnet pools using alias
  metalcloud-cli rp subnets 456

```
metalcloud-cli resource-pool get-subnet-pools <pool_id> [flags]
```

### Options

```
  -h, --help   help for get-subnet-pools
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

