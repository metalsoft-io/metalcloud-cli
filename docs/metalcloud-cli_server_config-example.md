## metalcloud-cli server config-example

Show a server configuration example

### Synopsis

Show an example of the configuration accepted by the server sub-commands that
take a --config-source flag.

Required Arguments:
  kind                   The configuration kind. One of:
                           register-production
                           import-unmanaged
                           connect-interface
                           set-interfaces-default-fabric
                           set-interfaces-redundancy-group

Examples:
  # Show the production server registration configuration
  metalcloud-cli server config-example register-production

  # Save the unmanaged server import configuration for editing
  metalcloud-cli server config-example import-unmanaged > unmanaged-server.json


```
metalcloud-cli server config-example kind [flags]
```

### Options

```
  -h, --help   help for config-example
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

