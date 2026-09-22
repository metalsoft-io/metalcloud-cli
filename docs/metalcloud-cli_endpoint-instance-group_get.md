## metalcloud-cli endpoint-instance-group get

Get endpoint instance group details

### Synopsis

Get the details of an endpoint instance group.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group

Examples:
  metalcloud-cli endpoint-instance-group get 12
  metalcloud-cli eig show 12

```
metalcloud-cli endpoint-instance-group get endpoint_instance_group_id [flags]
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

* [metalcloud-cli endpoint-instance-group](metalcloud-cli_endpoint-instance-group.md)	 - Endpoint instance group management

