## metalcloud-cli server-type delete

Delete a server type

### Synopsis

Delete a server type.

A server type can only be deleted while no server references it.

Required Arguments:
  server_type_id     The numeric ID of the server type to delete

Examples:
  # Delete the server type with ID 123
  metalcloud-cli server-type delete 123

  # Delete using the alias
  metalcloud-cli server-type rm 123


```
metalcloud-cli server-type delete server_type_id [flags]
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

* [metalcloud-cli server-type](metalcloud-cli_server-type.md)	 - Manage server types and hardware configurations

