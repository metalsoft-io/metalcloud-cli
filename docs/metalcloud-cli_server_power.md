## metalcloud-cli server power

Control server power state

### Synopsis

Control server power state.

This command allows you to control the power state of a server by sending
power management commands to the server's BMC/IPMI interface.

Required Arguments:
  server_id              The ID of the server to control

Required Flags:
  --action              Power action to perform

Valid Actions:
  on                    Power on the server
  off                   Hard power off the server
  reset                 Hard reset the server
  cycle                 Power cycle the server (off then on)
  soft                  Soft power off the server (graceful shutdown)

Examples:
  # Power on server
  metalcloud-cli server power 123 --action on

  # Hard power off server
  metalcloud-cli server power 123 --action off

  # Reset server
  metalcloud-cli server power 123 --action reset

  # Power cycle server
  metalcloud-cli server power 123 --action cycle

  # Soft power off (graceful shutdown)
  metalcloud-cli server power 123 --action soft


```
metalcloud-cli server power server_id [flags]
```

### Options

```
      --action string   Power action: on, off, reset, cycle, soft
  -h, --help            help for power
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

