## metalcloud-cli server-instance-group network acl update

Update a security rule of a network connection

### Synopsis

Update a security rule of one network connection of a server instance group.

Required Arguments:
  server_instance_group_id  The numeric ID of the server instance group
  connection_id             The numeric ID of the network connection
  rule_id                   The numeric ID of the security rule

Required Flags:
  --config-source string   Source of the updated rule. Can be 'pipe' or path to a JSON/YAML file.

Examples:
  metalcloud-cli server-instance-group network acl update 1234 5 3 --config-source rule.json
  echo '{"forwardingAction":"deny"}' | metalcloud-cli ig net acl update 1234 5 3 --config-source pipe

```
metalcloud-cli server-instance-group network acl update server_instance_group_id connection_id rule_id [flags]
```

### Options

```
      --config-source string   Source of the updated rule. Can be 'pipe' or path to a JSON/YAML file.
  -h, --help                   help for update
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

