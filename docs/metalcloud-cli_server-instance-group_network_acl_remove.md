## metalcloud-cli server-instance-group network acl remove

Remove a security rule from a network connection

### Synopsis

Remove a security rule from one network connection of a server instance group.

Required Arguments:
  server_instance_group_id  The numeric ID of the server instance group
  connection_id             The numeric ID of the network connection
  rule_id                   The numeric ID of the security rule

Examples:
  metalcloud-cli server-instance-group network acl remove 1234 5 3
  metalcloud-cli ig net acl rm 1234 5 3

```
metalcloud-cli server-instance-group network acl remove server_instance_group_id connection_id rule_id [flags]
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

* [metalcloud-cli server-instance-group network acl](metalcloud-cli_server-instance-group_network_acl.md)	 - Manage the security rules of a network connection

