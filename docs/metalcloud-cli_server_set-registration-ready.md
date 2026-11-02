## metalcloud-cli server set-registration-ready

Set a server to the registration_ready status

### Synopsis

Set the status of a server to "registration_ready".

Required Arguments:
  server_id              The ID of the server

Examples:
  # Mark planned server 123 as ready to be registered
  metalcloud-cli server set-registration-ready 123


```
metalcloud-cli server set-registration-ready server_id [flags]
```

### Options

```
  -h, --help   help for set-registration-ready
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

* [metalcloud-cli server](metalcloud-cli_server.md)	 - Server management

