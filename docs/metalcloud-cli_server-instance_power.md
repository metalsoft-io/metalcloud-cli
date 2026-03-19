## metalcloud-cli server-instance power

Set power state of a server instance

### Synopsis

Set the power state of a server instance.

This command sends a power command to a server instance via the orchestration layer.
For direct BMC/IPMI power control, use the 'server power' command instead.

Valid power actions:
  on     - Power on the server instance
  off    - Power off the server instance
  reset  - Hard reset the server instance
  soft   - Graceful shutdown of the server instance
  status - Get the current power status of the server instance

Arguments:
  server_instance_id  The numeric ID of the server instance
  action              Power action to perform (on, off, reset, soft, status)

Subcommands:
  status              Get the current power status of a server instance

Examples:
  # Power on server instance 5678
  metalcloud-cli server-instance power 5678 on

  # Gracefully shutdown server instance
  metalcloud-cli server-instance power 5678 soft

  # Hard reset server instance
  metalcloud-cli inst power 5678 reset

  # Get power status
  metalcloud-cli server-instance power 5678 status

```
metalcloud-cli server-instance power <server_instance_id> <on|off|reset|soft|status> [flags]
```

### Options

```
  -h, --help   help for power
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
* [metalcloud-cli server-instance power status](metalcloud-cli_server-instance_power_status.md)	 - Get power status of a server instance

