## metalcloud-cli endpoint-instance-group delete

Delete an endpoint instance group

### Synopsis

Delete an endpoint instance group.

The current revision of the group is fetched first and sent as If-Match, so the
delete fails if the group changed in the meantime.

Required Arguments:
  endpoint_instance_group_id   The numeric ID of the endpoint instance group

Examples:
  metalcloud-cli endpoint-instance-group delete 12
  metalcloud-cli eig rm 12

```
metalcloud-cli endpoint-instance-group delete endpoint_instance_group_id [flags]
```

### Options

```
  -h, --help   help for delete
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

