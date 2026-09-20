## metalcloud-cli server connect-interface

Record the network device port a server interface is cabled to

### Synopsis

Record the network device port a server interface is cabled to.

You can provide the connection either via command-line flags or by specifying a
configuration source using the --config-source flag. The configuration source can
be a path to a JSON file or 'pipe' to read from standard input.

Required Arguments:
  server_id              The ID of the server

Required Flags (when not using --config-source):
  --interface-id         The ID of the server interface to connect
  --port-id              The network device port name (e.g. Ethernet0)
  --hostname             The network device hostname

Optional Flags:
  --config-source        Source of the connection configuration. Can be 'pipe' or path to a JSON file.

Flag Dependencies:
  --config-source and --interface-id are mutually exclusive
  --interface-id, --port-id and --hostname must be used together

Examples:
  # Connect interface 1 of server 123 to port Ethernet0 of leaf-01
  metalcloud-cli server connect-interface 123 --interface-id 1 --port-id Ethernet0 --hostname leaf-01

  # Connect using a JSON configuration file
  metalcloud-cli server connect-interface 123 --config-source ./connect.json

  # Show the expected configuration
  metalcloud-cli server config-example connect-interface


```
metalcloud-cli server connect-interface server_id [flags]
```

### Options

```
      --config-source string   Source of the connection configuration. Can be 'pipe' or path to a JSON file.
  -h, --help                   help for connect-interface
      --hostname string        The network device hostname.
      --interface-id int       The ID of the server interface to connect.
      --port-id string         The network device port name.
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

