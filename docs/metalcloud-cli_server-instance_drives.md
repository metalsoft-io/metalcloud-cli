## metalcloud-cli server-instance drives

List the drives of a server instance

### Synopsis

List all drives attached to a server instance.

Required Arguments:
  server_instance_id  The numeric ID of the server instance

Examples:
  metalcloud-cli server-instance drives 5678
  metalcloud-cli inst drives 5678

```
metalcloud-cli server-instance drives server_instance_id [flags]
```

### Options

```
  -h, --help   help for drives
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

* [metalcloud-cli server-instance](metalcloud-cli_server-instance.md)	 - Manage individual server instances

