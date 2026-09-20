## metalcloud-cli server snapshots

List the configuration snapshots of a server

### Synopsis

List the configuration snapshots stored for a server.

Required Arguments:
  server_id              The ID of the server

Optional Flags:
  --kind                 Restrict the listing to one snapshot class

Examples:
  # List all snapshots of server 123
  metalcloud-cli server snapshots 123

  # List only the backup snapshots of server 123
  metalcloud-cli server snapshots 123 --kind backup


```
metalcloud-cli server snapshots server_id [flags]
```

### Options

```
  -h, --help          help for snapshots
      --kind string   Restrict the listing to one snapshot class.
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

