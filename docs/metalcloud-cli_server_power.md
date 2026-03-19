## metalcloud-cli server power

Control server power state

### Synopsis

Control server power state.

This command allows you to control the power state of a server by sending
power management commands to the server's BMC/IPMI interface.

Arguments:
  server_id              The ID of the server to control
  action                 Power action to perform (on, off, reset, cycle, soft, status)

Valid Actions:
  on                    Power on the server
  off                   Hard power off the server
  reset                 Hard reset the server
  cycle                 Power cycle the server (off then on)
  soft                  Soft power off the server (graceful shutdown)
  status                Get the current power status of the server

Subcommands:
  status                Get the current power status of a server

Examples:
  # Power on server
  metalcloud-cli server power 123 on

  # Hard power off server
  metalcloud-cli server power 123 off

  # Reset server
  metalcloud-cli server power 123 reset

  # Power cycle server
  metalcloud-cli server power 123 cycle

  # Soft power off (graceful shutdown)
  metalcloud-cli server power 123 soft

  # Get power status
  metalcloud-cli server power 123 status


```
metalcloud-cli server power <server_id> <on|off|reset|cycle|soft|status> [flags]
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

* [metalcloud-cli server](metalcloud-cli_server.md)	 - Server management
* [metalcloud-cli server power status](metalcloud-cli_server_power_status.md)	 - Get server power status

