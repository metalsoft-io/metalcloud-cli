## metalcloud-cli server get

Get detailed server information

### Synopsis

Get detailed information for a specific server.

This command retrieves comprehensive information about a server including its configuration,
status, hardware details, and optionally IPMI credentials.

Required Arguments:
  server_id              The ID of the server to retrieve information for

Optional Flags:
  --show-credentials     Include IPMI credentials (username and password) in the output

Examples:
  # Get basic server information
  metalcloud-cli server get 123

  # Get server information including IPMI credentials
  metalcloud-cli server get 123 --show-credentials


```
metalcloud-cli server get server_id [flags]
```

### Options

```
  -h, --help               help for get
      --show-credentials   If set returns the server IPMI credentials.
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

