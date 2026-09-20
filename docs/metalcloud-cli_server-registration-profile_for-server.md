## metalcloud-cli server-registration-profile for-server

Show the registration profile of a server

### Synopsis

Show the server registration profile a server was, or would be, registered with.

Required Arguments:
  server_id    The numeric ID of the server

Examples:
  # Show the registration profile of server 123
  metalcloud-cli server-registration-profile for-server 123


```
metalcloud-cli server-registration-profile for-server server_id [flags]
```

### Options

```
  -h, --help   help for for-server
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

* [metalcloud-cli server-registration-profile](metalcloud-cli_server-registration-profile.md)	 - Manage server registration profiles

