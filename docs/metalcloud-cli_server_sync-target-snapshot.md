## metalcloud-cli server sync-target-snapshot

Accept the current configuration of a server as the drift target

### Synopsis

Point the drift detection target snapshot of a server at its latest snapshot,
clearing the drift currently reported for the server.

Required Arguments:
  server_id              The ID of the server

Examples:
  # Accept the current configuration of server 123 as the new target
  metalcloud-cli server sync-target-snapshot 123


```
metalcloud-cli server sync-target-snapshot server_id [flags]
```

### Options

```
  -h, --help   help for sync-target-snapshot
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

