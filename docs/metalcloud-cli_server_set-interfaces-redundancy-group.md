## metalcloud-cli server set-interfaces-redundancy-group

Set the redundancy group of some server interfaces

### Synopsis

Group the given server interfaces into a redundancy group, or remove them from one.

You can provide the configuration either via command-line flags or by specifying a
configuration source using the --config-source flag. The configuration source can
be a path to a JSON file or 'pipe' to read from standard input.

Required Arguments:
  server_id              The ID of the server

Required Flags (when not using --config-source):
  --interface-id         ID of a server interface; repeat or comma-separate for several
  --group-index          The redundancy group index (or use --clear-group)

Optional Flags:
  --config-source        Source of the configuration. Can be 'pipe' or path to a JSON file.
  --clear-group          Remove the given interfaces from their redundancy group

Flag Dependencies:
  --config-source and --interface-id are mutually exclusive
  --group-index and --clear-group are mutually exclusive

Examples:
  # Put interfaces 1 and 2 of server 123 in redundancy group 1
  metalcloud-cli server set-interfaces-redundancy-group 123 --interface-id 1,2 --group-index 1

  # Remove interface 1 of server 123 from its redundancy group
  metalcloud-cli server set-interfaces-redundancy-group 123 --interface-id 1 --clear-group

  # Use a JSON configuration file
  metalcloud-cli server set-interfaces-redundancy-group 123 --config-source ./redundancy.json


```
metalcloud-cli server set-interfaces-redundancy-group server_id [flags]
```

### Options

```
      --clear-group            Remove the given interfaces from their redundancy group.
      --config-source string   Source of the configuration. Can be 'pipe' or path to a JSON file.
      --group-index int        The redundancy group index.
  -h, --help                   help for set-interfaces-redundancy-group
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

