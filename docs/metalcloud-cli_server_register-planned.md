## metalcloud-cli server register-planned

Start the registration of a planned server

### Synopsis

Start the registration of a server in the "planned" state.

The server must have a reachable BMC address and valid BMC credentials. Set them
first with "server update" if they were not provided when the server was created.

Required Arguments:
  server_id              The ID of the planned server

Examples:
  # Start the registration of planned server 123
  metalcloud-cli server register-planned 123


```
metalcloud-cli server register-planned server_id [flags]
```

### Options

```
  -h, --help   help for register-planned
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

