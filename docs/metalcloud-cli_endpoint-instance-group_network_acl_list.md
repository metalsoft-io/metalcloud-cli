## metalcloud-cli endpoint-instance-group network acl list

List the security rules of a network connection

### Synopsis

List the security rules of one network connection of an endpoint instance group.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group
  connection_id                The numeric ID of the network connection

Examples:
  metalcloud-cli endpoint-instance-group network acl list 12 5
  metalcloud-cli eig net acl ls 12 5

```
metalcloud-cli endpoint-instance-group network acl list endpoint_instance_group_id connection_id [flags]
```

### Options

```
  -h, --help   help for list
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

* [metalcloud-cli endpoint-instance-group network acl](metalcloud-cli_endpoint-instance-group_network_acl.md)	 - Manage the security rules of a network connection

