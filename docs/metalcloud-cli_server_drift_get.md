## metalcloud-cli server drift get

Get one configuration drift entry of a server

### Synopsis

Display one configuration drift entry of a server, including the configuration
difference that was detected.

Required Arguments:
  server_id              The ID of the server
  drift_id               The ID of the drift entry

Examples:
  # Show drift entry 9 of server 123
  metalcloud-cli server drift get 123 9


```
metalcloud-cli server drift get server_id drift_id [flags]
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

* [metalcloud-cli server drift](metalcloud-cli_server_drift.md)	 - Inspect the configuration drift of a server

