## metalcloud-cli server drift

Inspect the configuration drift of a server

### Synopsis

Inspect and acknowledge the configuration drift detected between the running
configuration of a server and its target snapshot.

Use "metalcloud-cli server drift [command] --help" for detailed information about each command.


### Options

```
  -h, --help   help for drift
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
* [metalcloud-cli server drift acknowledge](metalcloud-cli_server_drift_acknowledge.md)	 - Acknowledge a configuration drift entry of a server
* [metalcloud-cli server drift get](metalcloud-cli_server_drift_get.md)	 - Get one configuration drift entry of a server
* [metalcloud-cli server drift list](metalcloud-cli_server_drift_list.md)	 - List the configuration drift of a server

