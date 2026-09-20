## metalcloud-cli endpoint-instance-group network get

Get a network connection of an endpoint instance group

### Synopsis

Get the details of one network connection of an endpoint instance group.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group
  connection_id                The numeric ID of the network connection

Examples:
  metalcloud-cli endpoint-instance-group network get 12 5
  metalcloud-cli eig net show 12 5

```
metalcloud-cli endpoint-instance-group network get endpoint_instance_group_id connection_id [flags]
```

### Options

```
  -h, --help   help for get
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

* [metalcloud-cli endpoint-instance-group network](metalcloud-cli_endpoint-instance-group_network.md)	 - Manage the network configuration of an endpoint instance group

