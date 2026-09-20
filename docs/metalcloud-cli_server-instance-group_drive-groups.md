## metalcloud-cli server-instance-group drive-groups

List the drive groups of a server instance group

### Synopsis

List all drive groups attached to a server instance group.

Required Arguments:
  server_instance_group_id  The numeric ID of the server instance group

Examples:
  metalcloud-cli server-instance-group drive-groups 1234
  metalcloud-cli ig drives 1234

```
metalcloud-cli server-instance-group drive-groups server_instance_group_id [flags]
```

### Options

```
  -h, --help   help for drive-groups
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

* [metalcloud-cli server-instance-group](metalcloud-cli_server-instance-group.md)	 - Manage server instance groups within infrastructures

