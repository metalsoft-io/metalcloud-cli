## metalcloud-cli server firmware component-info

Get firmware component information

### Synopsis

Get detailed information for a specific firmware component.

This command retrieves comprehensive information about a firmware component
including its current version, available updates, and configuration options.

Required Arguments:
  server_id              The ID of the server
  component_id           The ID of the firmware component

Examples:
  # Get firmware component information
  metalcloud-cli server firmware component-info 123 456

  # Using alias
  metalcloud-cli server firmware get-component 123 456


```
metalcloud-cli server firmware component-info server_id component_id [flags]
```

### Options

```
  -h, --help   help for component-info
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

