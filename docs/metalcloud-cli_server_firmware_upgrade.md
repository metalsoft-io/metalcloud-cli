## metalcloud-cli server firmware upgrade

Upgrade firmware for all components on a server

### Synopsis

Upgrade firmware for all components on a server.

This command initiates a firmware upgrade process for all upgradeable components
on the specified server. The system will automatically determine which components
need updates and apply the latest available firmware versions.

Required Arguments:
  server_id              The ID of the server to upgrade

Examples:
  # Upgrade all firmware components for server with ID 123
  metalcloud-cli server firmware upgrade 123


```
metalcloud-cli server firmware upgrade server_id [flags]
```

### Options

```
  -h, --help   help for upgrade
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

