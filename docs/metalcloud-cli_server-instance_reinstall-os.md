## metalcloud-cli server-instance reinstall-os

Schedule OS reinstall for a server instance

### Synopsis

Schedule an OS reinstallation for a server instance.

This command marks the server instance for OS reinstallation. The reinstall
will take effect at the next infrastructure deploy.

Arguments:
  server_instance_id  The numeric ID of the server instance

Examples:
  # Schedule OS reinstall for server instance 5678
  metalcloud-cli server-instance reinstall-os 5678

  # Using alias
  metalcloud-cli inst reinstall 5678

```
metalcloud-cli server-instance reinstall-os <server_instance_id> [flags]
```

### Options

```
  -h, --help   help for reinstall-os
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

* [metalcloud-cli server-instance](metalcloud-cli_server-instance.md)	 - Manage individual server instances

