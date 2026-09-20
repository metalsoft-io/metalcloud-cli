## metalcloud-cli server set-interfaces-default-fabric

Set the default fabric of some server interfaces

### Synopsis

Assign the default fabric of the given server interfaces, or clear it.

You can provide the configuration either via command-line flags or by specifying a
configuration source using the --config-source flag. The configuration source can
be a path to a JSON file or 'pipe' to read from standard input.

Required Arguments:
  server_id              The ID of the server

Required Flags (when not using --config-source):
  --interface-id         ID of a server interface; repeat or comma-separate for several
  --fabric-id            The ID of the fabric to set as default (or use --clear-fabric)

Optional Flags:
  --config-source        Source of the configuration. Can be 'pipe' or path to a JSON file.
  --clear-fabric         Clear the default fabric of the given interfaces

Flag Dependencies:
  --config-source and --interface-id are mutually exclusive
  --fabric-id and --clear-fabric are mutually exclusive

Examples:
  # Set fabric 10 as the default fabric of interfaces 1 and 2 of server 123
  metalcloud-cli server set-interfaces-default-fabric 123 --interface-id 1,2 --fabric-id 10

  # Clear the default fabric of interface 1 of server 123
  metalcloud-cli server set-interfaces-default-fabric 123 --interface-id 1 --clear-fabric

  # Use a JSON configuration file
  metalcloud-cli server set-interfaces-default-fabric 123 --config-source ./fabric.json


```
metalcloud-cli server set-interfaces-default-fabric server_id [flags]
```

### Options

```
      --clear-fabric           Clear the default fabric of the given interfaces.
      --config-source string   Source of the configuration. Can be 'pipe' or path to a JSON file.
      --fabric-id int          The ID of the fabric to set as default.
  -h, --help                   help for set-interfaces-default-fabric
      --interface-id strings   IDs of the server interfaces to change.
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

