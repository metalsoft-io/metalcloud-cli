## metalcloud-cli server drift list

List the configuration drift of a server

### Synopsis

List the configuration drift entries recorded for a server.

Required Arguments:
  server_id              The ID of the server

Examples:
  # List the drift history of server 123
  metalcloud-cli server drift list 123


```
metalcloud-cli server drift list server_id [flags]
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

* [metalcloud-cli server drift](metalcloud-cli_server_drift.md)	 - Inspect the configuration drift of a server

