## metalcloud-cli server firmware inventory

Get firmware inventory from redfish

### Synopsis

Get firmware inventory from redfish.

This command retrieves the firmware inventory for the specified server using
the Redfish API, providing detailed information about all firmware components
installed on the server.

Required Arguments:
  server_id              The ID of the server to query

Examples:
  # Get firmware inventory for server with ID 123
  metalcloud-cli server firmware inventory 123


```
metalcloud-cli server firmware inventory server_id [flags]
```

### Options

```
  -h, --help   help for inventory
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

* [metalcloud-cli server firmware](metalcloud-cli_server_firmware.md)	 - Server firmware management

